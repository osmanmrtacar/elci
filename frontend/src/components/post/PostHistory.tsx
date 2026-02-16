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
        return (
          <svg className="w-4 h-4" viewBox="0 0 24 24" fill="currentColor">
            <path d="M19.59 6.69a4.83 4.83 0 0 1-3.77-4.25V2h-3.45v13.67a2.89 2.89 0 0 1-5.2 1.74 2.89 2.89 0 0 1 2.31-4.64 2.93 2.93 0 0 1 .88.13V9.4a6.84 6.84 0 0 0-1-.05A6.33 6.33 0 0 0 5 20.1a6.34 6.34 0 0 0 10.86-4.43v-7a8.16 8.16 0 0 0 4.77 1.52v-3.4a4.85 4.85 0 0 1-1-.1z" />
          </svg>
        )
      case 'x':
        return (
          <svg className="w-4 h-4" viewBox="0 0 24 24" fill="currentColor">
            <path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z" />
          </svg>
        )
      case 'instagram':
        return (
          <svg className="w-4 h-4" viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 2.163c3.204 0 3.584.012 4.85.07 3.252.148 4.771 1.691 4.919 4.919.058 1.265.069 1.645.069 4.849 0 3.205-.012 3.584-.069 4.849-.149 3.225-1.664 4.771-4.919 4.919-1.266.058-1.644.07-4.85.07-3.204 0-3.584-.012-4.849-.07-3.26-.149-4.771-1.699-4.919-4.92-.058-1.265-.07-1.644-.07-4.849 0-3.204.013-3.583.07-4.849.149-3.227 1.664-4.771 4.919-4.919 1.266-.057 1.645-.069 4.849-.069zM12 0C8.741 0 8.333.014 7.053.072 2.695.272.273 2.69.073 7.052.014 8.333 0 8.741 0 12c0 3.259.014 3.668.072 4.948.2 4.358 2.618 6.78 6.98 6.98C8.333 23.986 8.741 24 12 24c3.259 0 3.668-.014 4.948-.072 4.354-.2 6.782-2.618 6.979-6.98.059-1.28.073-1.689.073-4.948 0-3.259-.014-3.667-.072-4.947-.196-4.354-2.617-6.78-6.979-6.98C15.668.014 15.259 0 12 0zm0 5.838a6.162 6.162 0 100 12.324 6.162 6.162 0 000-12.324zM12 16a4 4 0 110-8 4 4 0 010 8zm6.406-11.845a1.44 1.44 0 100 2.881 1.44 1.44 0 000-2.881z" />
          </svg>
        )
      case 'youtube':
        return (
          <svg className="w-4 h-4" viewBox="0 0 24 24" fill="currentColor">
            <path d="M23.498 6.186a3.016 3.016 0 00-2.122-2.136C19.505 3.545 12 3.545 12 3.545s-7.505 0-9.377.505A3.017 3.017 0 00.502 6.186C0 8.07 0 12 0 12s0 3.93.502 5.814a3.016 3.016 0 002.122 2.136c1.871.505 9.376.505 9.376.505s7.505 0 9.377-.505a3.015 3.015 0 002.122-2.136C24 15.93 24 12 24 12s0-3.93-.502-5.814zM9.545 15.568V8.432L15.818 12l-6.273 3.568z" />
          </svg>
        )
      default:
        return null
    }
  }

  const getPlatformColor = (platform: Platform) => {
    switch (platform) {
      case 'tiktok':
        return '#fe2c55'
      case 'x':
        return '#000000'
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
