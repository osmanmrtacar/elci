export type Platform = 'tiktok' | 'x' | 'instagram' | 'youtube'

export interface PlatformConnection {
  platform: Platform
  username: string
  display_name: string
  avatar_url: string
  is_active: boolean
  connected_at: string
  last_used_at: string
}

export interface User {
  id: number
  username: string
  display_name: string
  avatar_url: string
  created_at?: string
  connected_platforms?: string[] // Array of platform names
}

export interface AuthResponse {
  token: string
  user: User
}
