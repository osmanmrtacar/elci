import axios from 'axios'
import crypto from 'crypto'
import { XTokenResponse, XUserInfo } from '../types'

export class XAuthService {
  private clientId: string
  private clientSecret: string
  private redirectUri: string
  private codeVerifier: string = ''

  constructor() {
    this.clientId = process.env.X_CLIENT_ID || ''
    this.clientSecret = process.env.X_CLIENT_SECRET || ''
    this.redirectUri = process.env.X_REDIRECT_URI || ''
  }

  // Generate OAuth 2.0 authorization URL with PKCE
  generateAuthUrl(): { url: string; codeVerifier: string; state: string } {
    // Generate code verifier for PKCE
    this.codeVerifier = this.generateCodeVerifier()
    const codeChallenge = this.generateCodeChallenge(this.codeVerifier)

    // Generate state for CSRF protection
    const state = crypto.randomBytes(32).toString('hex')

    // X OAuth 2.0 scopes (including media.write for uploads)
    const scopes = ['tweet.read', 'tweet.write', 'users.read', 'offline.access', 'media.write']

    const params = new URLSearchParams({
      response_type: 'code',
      client_id: this.clientId,
      redirect_uri: this.redirectUri,
      scope: scopes.join(' '),
      state: state,
      code_challenge: codeChallenge,
      code_challenge_method: 'S256',
    })

    const authUrl = `https://twitter.com/i/oauth2/authorize?${params.toString()}`

    return { url: authUrl, codeVerifier: this.codeVerifier, state }
  }

  // Exchange authorization code for access token
  async exchangeCodeForToken(
    code: string,
    codeVerifier: string
  ): Promise<XTokenResponse> {
    const params = new URLSearchParams({
      code: code,
      grant_type: 'authorization_code',
      client_id: this.clientId,
      redirect_uri: this.redirectUri,
      code_verifier: codeVerifier,
    })

    // Create Basic Auth header
    const credentials = Buffer.from(
      `${this.clientId}:${this.clientSecret}`
    ).toString('base64')

    try {
      const response = await axios.post(
        'https://api.twitter.com/2/oauth2/token',
        params.toString(),
        {
          headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
            Authorization: `Basic ${credentials}`,
          },
        }
      )

      return response.data
    } catch (error: any) {
      console.error('Token exchange error:', error.response?.data)
      throw new Error(`Failed to exchange code for token: ${error.message}`)
    }
  }

  // Refresh access token
  async refreshAccessToken(refreshToken: string): Promise<XTokenResponse> {
    const params = new URLSearchParams({
      refresh_token: refreshToken,
      grant_type: 'refresh_token',
      client_id: this.clientId,
    })

    const credentials = Buffer.from(
      `${this.clientId}:${this.clientSecret}`
    ).toString('base64')

    try {
      const response = await axios.post(
        'https://api.twitter.com/2/oauth2/token',
        params.toString(),
        {
          headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
            Authorization: `Basic ${credentials}`,
          },
        }
      )

      return response.data
    } catch (error: any) {
      console.error('Token refresh error:', error.response?.data)
      throw new Error(`Failed to refresh token: ${error.message}`)
    }
  }

  // Get user info
  async getUserInfo(accessToken: string): Promise<XUserInfo> {
    try {
      const response = await axios.get('https://api.twitter.com/2/users/me', {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
        params: {
          'user.fields': 'profile_image_url',
        },
      })

      return response.data.data
    } catch (error: any) {
      console.error('Get user info error:', error.response?.data)
      throw new Error(`Failed to get user info: ${error.message}`)
    }
  }

  // PKCE helper methods
  private generateCodeVerifier(): string {
    return crypto.randomBytes(32).toString('base64url')
  }

  private generateCodeChallenge(verifier: string): string {
    return crypto.createHash('sha256').update(verifier).digest('base64url')
  }
}
