import { useEffect, useState } from 'react'
import { Post } from '../../types/post'
import { Platform } from '../../types/user'
import { postService } from '../../services/postService'
import PostStatus from './PostStatus'

const PostHistory = ({ refreshTrigger }: { refreshTrigger: number }) => {
  const [posts, setPosts] = useState<Post[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [platformFilter, setPlatformFilter] = useState<Platform | 'all'>('all')

  useEffect(() => {
    fetchPosts()
  }, [refreshTrigger])

  const fetchPosts = async () => {
    try {
      setIsLoading(true)
      const data = await postService.getPosts()
      setPosts(data.posts || [])
    } catch (err: any) {
      console.error('Failed to fetch posts:', err)
      setError('Failed to load posts')
    } finally {
      setIsLoading(false)
    }
  }

  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleString()
  }

  const getPlatformIcon = (platform: Platform) => {
    switch (platform) {
      case 'tiktok':
        return '🎵'
      case 'x':
        return '𝕏'
      case 'instagram':
        return '📷'
      case 'youtube':
        return '▶️'
      default:
        return '📱'
    }
  }

  const getPlatformColor = (platform: Platform) => {
    switch (platform) {
      case 'tiktok':
        return '#fe2c55'
      case 'x':
        return '#1DA1F2'
      case 'instagram':
        return '#E1306C'
      case 'youtube':
        return '#FF0000'
      default:
        return '#666'
    }
  }

  const filteredPosts = platformFilter === 'all'
    ? posts
    : posts.filter(post => post.platform === platformFilter)

  const availablePlatforms = Array.from(new Set(posts.map(post => post.platform)))

  if (isLoading) {
    return (
      <div className="post-history-container">
        <h2>Post History</h2>
        <div className="loading-spinner"></div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="post-history-container">
        <h2>Post History</h2>
        <div className="error-message">{error}</div>
      </div>
    )
  }

  if (posts.length === 0) {
    return (
      <div className="post-history-container">
        <h2>Post History</h2>
        <p className="empty-state">No posts yet. Create your first post above!</p>
      </div>
    )
  }

  return (
    <div className="post-history-container">
      <div className="history-header">
        <h2>Post History</h2>
        {availablePlatforms.length > 1 && (
          <div className="platform-filter">
            <label htmlFor="platform-filter">Filter by platform:</label>
            <select
              id="platform-filter"
              value={platformFilter}
              onChange={(e) => setPlatformFilter(e.target.value as Platform | 'all')}
              className="filter-select"
            >
              <option value="all">All Platforms ({posts.length})</option>
              {availablePlatforms.map(platform => (
                <option key={platform} value={platform}>
                  {getPlatformIcon(platform)} {platform.toUpperCase()} ({posts.filter(p => p.platform === platform).length})
                </option>
              ))}
            </select>
          </div>
        )}
      </div>

      <div className="posts-list">
        {filteredPosts.map((post) => (
          <div key={post.id} className="post-card">
            <div className="post-header">
              <div className="header-left">
                <span
                  className="platform-badge"
                  style={{ backgroundColor: getPlatformColor(post.platform) }}
                >
                  {getPlatformIcon(post.platform)} {post.platform.toUpperCase()}
                </span>
                <PostStatus status={post.status} postId={post.id} />
              </div>
              <span className="post-date">{formatDate(post.created_at)}</span>
            </div>

            <div className="post-content">
              <p className="post-caption">
                {post.caption || <em>No caption</em>}
              </p>
              <a
                href={post.video_url}
                target="_blank"
                rel="noopener noreferrer"
                className="video-link"
              >
                View Video
              </a>
            </div>

            {(post.share_url || post.tiktok_url) && (
              <a
                href={post.share_url || post.tiktok_url}
                target="_blank"
                rel="noopener noreferrer"
                className="platform-link"
                style={{ borderColor: getPlatformColor(post.platform) }}
              >
                View on {post.platform.toUpperCase()} →
              </a>
            )}

            {post.error_message && (
              <div className="error-message">
                {post.error_message}
              </div>
            )}
          </div>
        ))}
      </div>

      <style>{`
        .history-header {
          display: flex;
          justify-content: space-between;
          align-items: center;
          margin-bottom: 20px;
          flex-wrap: wrap;
          gap: 12px;
        }

        .history-header h2 {
          margin: 0;
        }

        .platform-filter {
          display: flex;
          align-items: center;
          gap: 8px;
        }

        .platform-filter label {
          font-size: 14px;
          color: #666;
        }

        .filter-select {
          padding: 8px 12px;
          border: 2px solid #e0e0e0;
          border-radius: 6px;
          font-size: 14px;
          background: white;
          cursor: pointer;
          transition: border-color 0.2s ease;
        }

        .filter-select:hover {
          border-color: #4CAF50;
        }

        .filter-select:focus {
          outline: none;
          border-color: #4CAF50;
        }

        .post-header {
          display: flex;
          justify-content: space-between;
          align-items: flex-start;
          margin-bottom: 12px;
          gap: 12px;
        }

        .header-left {
          display: flex;
          align-items: center;
          gap: 8px;
          flex-wrap: wrap;
        }

        .platform-badge {
          display: inline-flex;
          align-items: center;
          gap: 4px;
          padding: 4px 10px;
          border-radius: 4px;
          font-size: 12px;
          font-weight: 600;
          color: white;
          white-space: nowrap;
        }

        .platform-link {
          display: inline-block;
          margin-top: 12px;
          padding: 8px 16px;
          color: #1976d2;
          text-decoration: none;
          font-weight: 500;
          border: 2px solid;
          border-radius: 6px;
          transition: all 0.2s ease;
        }

        .platform-link:hover {
          background: #f0f7ff;
          transform: translateX(4px);
        }
      `}</style>
    </div>
  )
}

export default PostHistory
