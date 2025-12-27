import { Router, Request, Response } from 'express'
import { XPostService } from '../services/xPost.service'
import { XAuthService } from '../services/xAuth.service'
import { XMediaService } from '../services/xMedia.service'
import { tokenStore } from '../services/tokenStore.service'

const router = Router()
const xPostService = new XPostService()
const xAuthService = new XAuthService()
const xMediaService = new XMediaService()

// Create a post (tweet)
router.post('/', async (req: Request, res: Response) => {
  try {
    const { userId, text, media_urls } = req.body

    if (!userId || !text) {
      return res.status(400).json({ error: 'userId and text are required' })
    }

    // Get stored token
    let storedToken = tokenStore.getToken(userId)
    if (!storedToken) {
      return res.status(401).json({ error: 'Not authenticated' })
    }

    // Check if token expired and refresh if needed
    if (new Date() > storedToken.expiresAt) {
      console.log('Token expired, refreshing...')
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

      storedToken = tokenStore.getToken(userId)!
    }

    // Upload media if URLs provided (using v2 API with OAuth 2.0)
    let media_ids: string[] = []
    if (media_urls && Array.isArray(media_urls) && media_urls.length > 0) {
      console.log(`Uploading ${media_urls.length} media file(s)...`)
      for (const mediaUrl of media_urls) {
        try {
          const mediaId = await xMediaService.uploadFromUrl(
            storedToken.accessToken,
            mediaUrl
          )
          media_ids.push(mediaId)
        } catch (uploadError: any) {
          console.error(`Failed to upload media from ${mediaUrl}:`, uploadError)
          return res.status(500).json({
            error: 'Failed to upload media',
            details: uploadError.message,
          })
        }
      }
    }

    // Create post
    const result = await xPostService.createPost(storedToken.accessToken, {
      text,
      media_ids,
    })

    console.log('Post created:', result.data.id)

    res.json({
      success: true,
      post: result.data,
      media_ids,
    })
  } catch (error: any) {
    console.error('Create post error:', error)
    res.status(500).json({
      error: 'Failed to create post',
      details: error.message,
    })
  }
})

// Get user's tweets
router.get('/:userId', async (req: Request, res: Response) => {
  try {
    const { userId } = req.params

    // Get stored token
    let storedToken = tokenStore.getToken(userId)
    if (!storedToken) {
      return res.status(401).json({ error: 'Not authenticated' })
    }

    // Check if token expired
    if (new Date() > storedToken.expiresAt) {
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

      storedToken = tokenStore.getToken(userId)!
    }

    const tweets = await xPostService.getUserTweets(
      storedToken.accessToken,
      userId
    )

    res.json({ tweets: tweets.data || [] })
  } catch (error: any) {
    console.error('Get tweets error:', error)
    res.status(500).json({ error: 'Failed to get tweets' })
  }
})

export default router
