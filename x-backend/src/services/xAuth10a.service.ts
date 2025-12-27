import OAuth from 'oauth-1.0a'
import crypto from 'crypto'
import axios from 'axios'

export class XAuth10aService {
  private oauth: OAuth
  private apiKey: string
  private apiSecret: string
  private callbackUrl: string

  constructor() {
    this.apiKey = process.env.X_API_KEY || ''
    this.apiSecret = process.env.X_API_SECRET || ''
    this.callbackUrl =
      process.env.X_REDIRECT_URI?.replace('/x/callback', '/x/callback-media') ||
      ''

    this.oauth = new OAuth({
      consumer: {
        key: this.apiKey,
        secret: this.apiSecret,
      },
      signature_method: 'HMAC-SHA1',
      hash_function(base_string, key) {
        return crypto
          .createHmac('sha1', key)
          .update(base_string)
          .digest('base64')
      },
    })
  }

  // Step 1: Get request token
  async getRequestToken(): Promise<{
    oauth_token: string
    oauth_token_secret: string
  }> {
    const requestData = {
      url: 'https://api.twitter.com/oauth/request_token',
      method: 'POST',
      data: { oauth_callback: this.callbackUrl },
    }

    const authorized = this.oauth.authorize(requestData as any)
    const authHeader = this.oauth.toHeader(authorized)

    console.log('Request token URL:', requestData.url)
    console.log('Callback URL:', this.callbackUrl)

    try {
      const formData = new URLSearchParams()
      formData.append('oauth_callback', this.callbackUrl)

      const response = await axios.post(requestData.url, formData.toString(), {
        headers: {
          ...authHeader,
          'Content-Type': 'application/x-www-form-urlencoded',
        },
      })

      const params = new URLSearchParams(response.data)
      return {
        oauth_token: params.get('oauth_token') || '',
        oauth_token_secret: params.get('oauth_token_secret') || '',
      }
    } catch (error: any) {
      console.error('Request token error:', error.response?.data || error.message)
      throw new Error(`Failed to get request token: ${error.message}`)
    }
  }

  // Step 2: Generate authorization URL
  getAuthorizationUrl(oauthToken: string): string {
    return `https://api.twitter.com/oauth/authorize?oauth_token=${oauthToken}`
  }

  // Step 3: Exchange for access token
  async getAccessToken(
    oauthToken: string,
    oauthTokenSecret: string,
    oauthVerifier: string
  ): Promise<{
    oauth_token: string
    oauth_token_secret: string
    user_id: string
    screen_name: string
  }> {
    const requestData = {
      url: 'https://api.twitter.com/oauth/access_token',
      method: 'POST',
      data: { oauth_verifier: oauthVerifier },
    }

    const authHeader = this.oauth.toHeader(
      this.oauth.authorize(requestData as any, {
        key: oauthToken,
        secret: oauthTokenSecret,
      })
    )

    try {
      const formData = new URLSearchParams()
      formData.append('oauth_verifier', oauthVerifier)

      const response = await axios.post(requestData.url, formData.toString(), {
        headers: {
          ...authHeader,
          'Content-Type': 'application/x-www-form-urlencoded',
        },
      })

      const params = new URLSearchParams(response.data)
      return {
        oauth_token: params.get('oauth_token') || '',
        oauth_token_secret: params.get('oauth_token_secret') || '',
        user_id: params.get('user_id') || '',
        screen_name: params.get('screen_name') || '',
      }
    } catch (error: any) {
      console.error('Access token error:', error.response?.data || error.message)
      throw new Error(`Failed to get access token: ${error.message}`)
    }
  }

  // Generate OAuth 1.0a authorization header for media upload
  getMediaUploadHeaders(
    url: string,
    method: string,
    accessToken: string,
    accessTokenSecret: string,
    additionalParams?: any
  ): any {
    const requestData = {
      url,
      method,
      data: additionalParams || {},
    }

    return this.oauth.toHeader(
      this.oauth.authorize(requestData as any, {
        key: accessToken,
        secret: accessTokenSecret,
      })
    )
  }
}
