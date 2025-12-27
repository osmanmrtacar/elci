import { useEffect, useState } from 'react'
import { PostStatus as PostStatusType } from '../../types/post'
import { postService } from '../../services/postService'

interface PostStatusProps {
  status: PostStatusType
  postId: number
}

const PostStatus = ({ status: initialStatus, postId }: PostStatusProps) => {
  const [status, setStatus] = useState(initialStatus)

  useEffect(() => {
    // Poll for status updates if post is in progress
    if (status === 'pending' || status === 'processing') {
      const interval = setInterval(async () => {
        try {
          const data = await postService.getPostStatus(postId)
          setStatus(data.status)

          // Stop polling if post is complete or failed
          if (data.status === 'published' || data.status === 'failed') {
            clearInterval(interval)
          }
        } catch (error) {
          console.error('Failed to fetch post status:', error)
        }
      }, 5000) // Poll every 5 seconds

      return () => clearInterval(interval)
    }
  }, [status, postId])

  const getStatusConfig = () => {
    switch (status) {
      case 'pending':
        return { label: 'Pending', color: '#ffa500', icon: '⏳' }
      case 'processing':
        return { label: 'Processing', color: '#2196f3', icon: '🔄' }
      case 'published':
        return { label: 'Published', color: '#4caf50', icon: '✅' }
      case 'failed':
        return { label: 'Failed', color: '#f44336', icon: '❌' }
      default:
        return { label: status, color: '#999', icon: '❓' }
    }
  }

  const config = getStatusConfig()

  return (
    <div className="post-status">
      <span
        className="status-badge"
        style={{
          backgroundColor: config.color,
          color: 'white',
          padding: '4px 12px',
          borderRadius: '12px',
          fontSize: '12px',
          fontWeight: 'bold',
          display: 'inline-flex',
          alignItems: 'center',
          gap: '4px',
        }}
      >
        <span>{config.icon}</span>
        <span>{config.label}</span>
      </span>
    </div>
  )
}

export default PostStatus
