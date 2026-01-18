import { Platform } from './user'

export type PostStatus = 'pending' | 'processing' | 'published' | 'failed'
export type MediaType = 'video' | 'image'

export interface Post {
  id: number
  video_url: string
  caption: string
  status: PostStatus
  platform: Platform
  platform_post_id?: string
  share_url?: string
  media_type?: MediaType
  // Legacy fields (for backward compatibility)
  tiktok_post_id?: string
  tiktok_url?: string
  error_message?: string
  created_at: string
  published_at?: string
}

export interface CreatePostRequest {
  platforms: Platform[]
  media_url: string
  caption: string
}

export interface PostsResponse {
  posts: Post[]
  count: number
  limit: number
  offset: number
}
