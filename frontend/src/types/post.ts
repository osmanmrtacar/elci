import { Platform } from './user'

export type PostStatus = 'pending' | 'processing' | 'published' | 'failed'
export type MediaType = 'video' | 'image'

// TikTok Privacy Level options
export type TikTokPrivacyLevel = 'PUBLIC_TO_EVERYONE' | 'MUTUAL_FOLLOW_FRIENDS' | 'FOLLOWER_OF_CREATOR' | 'SELF_ONLY'

// TikTok-specific settings (required by TikTok UX Guidelines)
export interface TikTokSettings {
  title?: string                    // Required: Editable title for the video
  privacy_level?: TikTokPrivacyLevel // Required: User must select (no default)
  allow_comment?: boolean           // Default: false (unchecked)
  allow_duet?: boolean              // Default: false (unchecked)
  allow_stitch?: boolean            // Default: false (unchecked)
  // Commercial content disclosure
  is_brand_content?: boolean        // "Your Brand" - promoting yourself
  is_brand_organic?: boolean        // "Branded Content" - paid partnership
  auto_add_music?: boolean          // Auto-add music (for TikTok photo posts)
  direct_post?: boolean             // Direct Post (true) vs Send to Inbox (false)
}

export interface Post {
  id: number
  video_url: string
  caption: string
  status: PostStatus
  platform: Platform
  platform_post_id?: string
  share_url?: string
  media_type?: MediaType
  // TikTok-specific fields
  title?: string
  privacy_level?: TikTokPrivacyLevel
  // Legacy fields (for backward compatibility)
  tiktok_post_id?: string
  tiktok_url?: string
  error_message?: string
  created_at: string
  published_at?: string
}

export interface CreatePostRequest {
  platforms: Platform[]
  media_url?: string                // Primary media URL (for single media)
  media_urls?: string[]             // Multiple media URLs (for carousel/multi-image)
  caption: string
  // TikTok-specific settings
  tiktok_settings?: TikTokSettings
}

export interface PostsResponse {
  posts: Post[]
  count: number
  limit: number
  offset: number
}
