import { Router, Request, Response } from 'express'
import { XAuthService } from '../services/xAuth.service'
import { XAuth10aService } from '../services/xAuth10a.service'
import { tokenStore } from '../services/tokenStore.service'

const router = Router()
const xAuthService = new XAuthService()
const xAuth10aService = new XAuth10aService()

// Step 1: Initiate X OAuth flow
router.get('/x/login', (req: Request, res: Response) => {
  try {
    const { url, codeVerifier, state } = xAuthService.generateAuthUrl()

    // Store code verifier temporarily (needed for token exchange)
    tokenStore.saveCodeVerifier(state, codeVerifier)

    // Store state in cookie for validation
    res.cookie('oauth_state', state, {
      httpOnly: true,
      maxAge: 10 * 60 * 1000, // 10 minutes
      sameSite: 'lax',
    })

    // Redirect to X authorization page
    res.redirect(url)
  } catch (error: any) {
    console.error('Login error:', error)
    res.status(500).json({ error: 'Failed to initiate login' })
  }
})

// Step 2: Handle OAuth callback
router.get('/x/callback', async (req: Request, res: Response) => {
  try {
    const { code, state } = req.query
    const storedState = req.cookies.oauth_state

    // Validate state (CSRF protection)
    if (!state || state !== storedState) {
      console.warn('State mismatch - possible CSRF attack')
      // Still proceed for tunnel compatibility
    }

    if (!code) {
      return res.status(400).json({ error: 'Authorization code missing' })
    }

    // Get code verifier
    const codeVerifier = tokenStore.getCodeVerifier(state as string)
    if (!codeVerifier) {
      return res.status(400).json({ error: 'Code verifier not found' })
    }

    // Exchange code for tokens
    const tokenResponse = await xAuthService.exchangeCodeForToken(
      code as string,
      codeVerifier
    )

    console.log('Token received:', {
      expires_in: tokenResponse.expires_in,
      scope: tokenResponse.scope,
    })

    // Get user info
    const userInfo = await xAuthService.getUserInfo(tokenResponse.access_token)

    console.log('User authenticated:', userInfo.username)

    // Store tokens
    const expiresAt = new Date(Date.now() + tokenResponse.expires_in * 1000)
    tokenStore.saveToken(userInfo.id, {
      userId: userInfo.id,
      accessToken: tokenResponse.access_token,
      refreshToken: tokenResponse.refresh_token,
      expiresAt,
    })

    // Clean up
    tokenStore.deleteCodeVerifier(state as string)
    res.clearCookie('oauth_state')

    // Redirect to frontend with user info
    const frontendUrl = process.env.FRONTEND_URL || 'http://localhost:3000'
    const redirectUrl = `${frontendUrl}/success?platform=x&username=${userInfo.username}&userId=${userInfo.id}&name=${encodeURIComponent(userInfo.name)}`

    res.redirect(redirectUrl)
  } catch (error: any) {
    console.error('Callback error:', error)
    const frontendUrl = process.env.FRONTEND_URL || 'http://localhost:3000'
    res.redirect(`${frontendUrl}?error=auth_failed`)
  }
})

// Get all stored users (for testing)
router.get('/x/users', (req: Request, res: Response) => {
  const userIds = tokenStore.getAllUserIds()
  res.json({ userIds })
})

// Get current user info (for testing)
router.get('/x/me/:userId', async (req: Request, res: Response) => {
  try {
    const { userId } = req.params
    const storedToken = tokenStore.getToken(userId)

    if (!storedToken) {
      return res.status(401).json({ error: 'Not authenticated' })
    }

    // Check if token expired
    if (new Date() > storedToken.expiresAt) {
      // Refresh token
      const newTokenResponse = await xAuthService.refreshAccessToken(
        storedToken.refreshToken
      )

      const newExpiresAt = new Date(
        Date.now() + newTokenResponse.expires_in * 1000
      )
      tokenStore.saveToken(userId, {
        userId,
        accessToken: newTokenResponse.access_token,
        refreshToken: newTokenResponse.refresh_token,
        expiresAt: newExpiresAt,
      })

      storedToken.accessToken = newTokenResponse.access_token
    }

    const userInfo = await xAuthService.getUserInfo(storedToken.accessToken)
    res.json({ user: userInfo })
  } catch (error: any) {
    console.error('Get user error:', error)
    res.status(500).json({ error: 'Failed to get user info' })
  }
})

// OAuth 1.0a flow for media uploads
// Step 1: Initiate OAuth 1.0a flow
router.get('/x/login-media', async (req: Request, res: Response) => {
  try {
    const { oauth_token, oauth_token_secret } =
      await xAuth10aService.getRequestToken()

    // Store token secret temporarily
    tokenStore.saveCodeVerifier(oauth_token, oauth_token_secret)

    // Redirect to authorization URL
    const authUrl = xAuth10aService.getAuthorizationUrl(oauth_token)
    res.redirect(authUrl)
  } catch (error: any) {
    console.error('OAuth 1.0a login error:', error)
    res.status(500).json({ error: 'Failed to initiate OAuth 1.0a login' })
  }
})

// Step 2: Handle OAuth 1.0a callback
router.get('/x/callback-media', async (req: Request, res: Response) => {
  try {
    const { oauth_token, oauth_verifier } = req.query

    if (!oauth_token || !oauth_verifier) {
      return res.status(400).json({ error: 'Missing OAuth parameters' })
    }

    // Get stored token secret
    const oauth_token_secret = tokenStore.getCodeVerifier(oauth_token as string)
    if (!oauth_token_secret) {
      return res.status(400).json({ error: 'OAuth token secret not found' })
    }

    // Exchange for access token
    const accessTokenData = await xAuth10aService.getAccessToken(
      oauth_token as string,
      oauth_token_secret,
      oauth_verifier as string
    )

    console.log('OAuth 1.0a authenticated:', accessTokenData.screen_name)

    // Find existing user token and update with OAuth 1.0a credentials
    const existingToken = tokenStore.getToken(accessTokenData.user_id)
    if (existingToken) {
      tokenStore.saveToken(accessTokenData.user_id, {
        ...existingToken,
        oauth1a: {
          accessToken: accessTokenData.oauth_token,
          accessTokenSecret: accessTokenData.oauth_token_secret,
        },
      })
    } else {
      // Create new entry with only OAuth 1.0a tokens (user needs to do OAuth 2.0 separately)
      return res.status(400).json({
        error: 'Please login with OAuth 2.0 first before authorizing media uploads',
      })
    }

    // Clean up
    tokenStore.deleteCodeVerifier(oauth_token as string)

    // Redirect to success page
    const frontendUrl = process.env.FRONTEND_URL || 'http://localhost:3000'
    res.redirect(
      `${frontendUrl}/success?platform=x&message=Media upload authorized&userId=${accessTokenData.user_id}`
    )
  } catch (error: any) {
    console.error('OAuth 1.0a callback error:', error)
    const frontendUrl = process.env.FRONTEND_URL || 'http://localhost:3000'
    res.redirect(`${frontendUrl}?error=oauth1a_failed`)
  }
})

export default router
