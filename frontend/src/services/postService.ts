import api from './api'
import { Post, CreatePostRequest, PostsResponse, TikTokCreatorInfo } from '../types/post'

export const postService = {
  // Create a new post
  createPost: async (data: CreatePostRequest): Promise<Post> => {
    const response = await api.post('/api/v1/posts', data)
    return response.data.post
  },

  // Get all posts
  getPosts: async (limit = 20, offset = 0): Promise<PostsResponse> => {
    const response = await api.get('/api/v1/posts', {
      params: { limit, offset },
    })
    return response.data
  },

  // Get a specific post
  getPost: async (id: number): Promise<Post> => {
    const response = await api.get(`/api/v1/posts/${id}`)
    return response.data.post
  },

  // Get post status (for polling)
  getPostStatus: async (id: number): Promise<Post> => {
    const response = await api.get(`/api/v1/posts/${id}/status`)
    return response.data
  },

  // Get TikTok creator info (privacy levels, interaction permissions, video duration limits)
  getTikTokCreatorInfo: async (): Promise<TikTokCreatorInfo> => {
    const response = await api.get('/api/v1/tiktok/creator-info')
    return response.data.creator_info
  },
}
