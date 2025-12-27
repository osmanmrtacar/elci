export interface XTokenResponse {
  access_token: string
  token_type: string
  expires_in: number
  refresh_token: string
  scope: string
}

export interface XUserInfo {
  id: string
  name: string
  username: string
  profile_image_url?: string
}

export interface XPostRequest {
  text: string
  media_ids?: string[]
  media_urls?: string[]
}

export interface XPostResponse {
  data: {
    id: string
    text: string
  }
}

export interface StoredToken {
  userId: string
  accessToken: string
  refreshToken: string
  expiresAt: Date
  oauth1a?: {
    accessToken: string
    accessTokenSecret: string
  }
}
