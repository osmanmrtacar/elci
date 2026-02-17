import { useEffect, useState } from 'react'
import { Post } from '../../types/post'
import { Platform } from '../../types/user'
import { postService } from '../../services/postService'
import PostStatus from './PostStatus'

const PostHistory = ({ refreshTrigger }: { refreshTrigger: number }) => {
  const [posts, setPosts] = useState<Post[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

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
    const now = new Date()
    const diffInSeconds = Math.floor((now.getTime() - date.getTime()) / 1000)

    if (diffInSeconds < 60) return 'Just now'
    if (diffInSeconds < 3600) return `${Math.floor(diffInSeconds / 60)}m ago`
    if (diffInSeconds < 86400) return `${Math.floor(diffInSeconds / 3600)}h ago`
    if (diffInSeconds < 604800) return `${Math.floor(diffInSeconds / 86400)}d ago`

    return date.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
    })
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
        return '#111827'
      case 'x':
        return '#111827'
      case 'instagram':
        return 'linear-gradient(135deg, #a855f7 0%, #ec4899 45%, #fb923c 100%)'
      case 'youtube':
        return '#dc2626'
      default:
        return '#4b5563'
    }
  }

  if (isLoading) {
    return (
      <div className="post-history-container">
        <div className="history-header">
          <h2>Recent Posts</h2>
        </div>
        <div className="loading-spinner"></div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="post-history-container">
        <div className="history-header">
          <h2>Recent Posts</h2>
        </div>
        <div className="error-message">{error}</div>
      </div>
    )
  }

  if (posts.length === 0) {
    return (
      <div className="post-history-container">
        <div className="history-header">
          <h2>Recent Posts</h2>
        </div>
        <div className="empty-state">
          <div className="empty-icon"></div>
          <p>No posts yet</p>
        </div>
      </div>
    )
  }

  return (
    <div className="post-history-container">
      <div className="history-header">
        <h2>Recent Posts</h2>
      </div>

      <div className="posts-list">
        {posts.map((post) => (
          <div key={post.id} className="post-card">
            <div className="post-header">
              <div className="header-left">
                <span
                  className="platform-badge"
                  style={{ background: getPlatformColor(post.platform) }}
                >
                  {getPlatformIcon(post.platform)}
                  <span className="platform-text">{post.platform.toUpperCase()}</span>
                </span>
              </div>
              <div className="header-center">
                <PostStatus status={post.status} postId={post.id} />
              </div>
              <span className="post-date">{formatDate(post.created_at)}</span>
            </div>

            <p className="post-caption">
              {post.caption || <em>No caption</em>}
            </p>

            {(post.share_url || post.tiktok_url) && (
              <a
                href={post.share_url || post.tiktok_url}
                target="_blank"
                rel="noopener noreferrer"
                className="platform-link"
              >
                <span>View post</span>
                <span className="link-arrow">↗</span>
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
          margin-bottom: 16px;
        }

        .history-header h2 {
          margin: 0;
          font-size: 14px;
          font-weight: 600;
          color: #111827;
        }

        .posts-list {
          display: flex;
          flex-direction: column;
          gap: 12px;
        }

        .post-card {
          background: #f9fafb;
          border: 1px solid #f3f4f6;
          border-radius: 10px;
          padding: 14px 16px;
          transition: border-color 0.2s ease, box-shadow 0.2s ease;
        }

        .post-card:hover {
          border-color: #e5e7eb;
          box-shadow: 0 4px 12px rgba(17, 24, 39, 0.06);
        }

        .post-header {
          display: grid;
          grid-template-columns: auto 1fr auto;
          align-items: center;
          margin-bottom: 8px;
          gap: 12px;
        }

        .header-left {
          display: flex;
          align-items: center;
          gap: 10px;
          flex-wrap: wrap;
        }

        .header-center {
          display: flex;
          align-items: center;
          justify-content: center;
        }

        .platform-badge {
          display: inline-flex;
          align-items: center;
          gap: 6px;
          padding: 4px 8px;
          border-radius: 6px;
          font-size: 11px;
          font-weight: 600;
          color: #ffffff;
          white-space: nowrap;
        }

        .platform-badge svg {
          width: 14px;
          height: 14px;
          display: block;
          flex-shrink: 0;
        }

        .platform-text {
          letter-spacing: 0.02em;
        }

        .post-date {
          font-size: 12px;
          color: #6b7280;
          white-space: nowrap;
        }

        .post-caption {
          font-size: 14px;
          color: #374151;
          margin: 0 0 10px;
          display: -webkit-box;
          -webkit-line-clamp: 2;
          -webkit-box-orient: vertical;
          overflow: hidden;
        }

        .platform-link {
          display: inline-flex;
          align-items: center;
          gap: 6px;
          font-size: 12px;
          color: #9ca3af;
          text-decoration: none;
          font-weight: 500;
          opacity: 0;
          transition: color 0.2s ease, opacity 0.2s ease;
        }

        .post-card:hover .platform-link {
          opacity: 1;
        }

        .platform-link:hover {
          color: #111827;
        }

        .link-arrow {
          font-size: 12px;
        }

        .empty-state {
          display: flex;
          flex-direction: column;
          align-items: center;
          gap: 8px;
          padding: 32px 0;
          color: #6b7280;
          font-size: 14px;
        }

        .empty-icon {
          width: 44px;
          height: 44px;
          border-radius: 999px;
          background: #f3f4f6;
        }
      `}</style>
    </div>
  )
}

export default PostHistory
