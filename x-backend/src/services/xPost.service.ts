import axios from 'axios'
import { XPostRequest, XPostResponse } from '../types'

export class XPostService {
  // Create a post (tweet)
  async createPost(
    accessToken: string,
    postData: XPostRequest
  ): Promise<XPostResponse> {
    try {
      const requestBody: any = { text: postData.text }

      // Add media if provided
      if (postData.media_ids && postData.media_ids.length > 0) {
        requestBody.media = {
          media_ids: postData.media_ids,
        }
      }

      const response = await axios.post(
        'https://api.twitter.com/2/tweets',
        requestBody,
        {
          headers: {
            Authorization: `Bearer ${accessToken}`,
            'Content-Type': 'application/json',
          },
        }
      )

      return response.data
    } catch (error: any) {
      console.error('Create post error:', error.response?.data)
      throw new Error(
        `Failed to create post: ${error.response?.data?.detail || error.message}`
      )
    }
  }

  // Get user's tweets
  async getUserTweets(accessToken: string, userId: string) {
    try {
      const response = await axios.get(
        `https://api.twitter.com/2/users/${userId}/tweets`,
        {
          headers: {
            Authorization: `Bearer ${accessToken}`,
          },
          params: {
            max_results: 10,
            'tweet.fields': 'created_at',
          },
        }
      )

      return response.data
    } catch (error: any) {
      console.error('Get user tweets error:', error.response?.data)
      throw new Error(`Failed to get tweets: ${error.message}`)
    }
  }
}
