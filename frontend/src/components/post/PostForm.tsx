import { useState, useMemo } from 'react'
import { postService } from '../../services/postService'
import { useAuth } from '../../context/AuthContext'
import { Platform } from '../../types/user'

interface PostFormProps {
  onPostCreated: () => void
}

type MediaType = 'video' | 'image' | 'unknown'

const imageExtensions = ['.jpg', '.jpeg', '.png', '.gif', '.webp', '.bmp', '.heic', '.heif']
const videoExtensions = ['.mp4', '.mov', '.webm', '.avi', '.mkv', '.m4v', '.ts', '.3gp']

const detectMediaType = (url: string): MediaType => {
  if (!url) return 'unknown'
  const lowerUrl = url.toLowerCase()

  // Extract path from URL
  try {
    const parsedUrl = new URL(url)
    const path = parsedUrl.pathname.toLowerCase()

    if (imageExtensions.some(ext => path.endsWith(ext))) {
      return 'image'
    }
    if (videoExtensions.some(ext => path.endsWith(ext))) {
      return 'video'
    }
  } catch {
    // If URL parsing fails, try simple extension check
    if (imageExtensions.some(ext => lowerUrl.endsWith(ext))) {
      return 'image'
    }
    if (videoExtensions.some(ext => lowerUrl.endsWith(ext))) {
      return 'video'
    }
  }

  return 'unknown'
}

const PostForm = ({ onPostCreated }: PostFormProps) => {
  const { connectedPlatforms } = useAuth()
  const [mediaUrl, setMediaUrl] = useState('')
  const [caption, setCaption] = useState('')
  const [selectedPlatforms, setSelectedPlatforms] = useState<Platform[]>([])
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const detectedMediaType = useMemo(() => detectMediaType(mediaUrl), [mediaUrl])

  const handlePlatformToggle = (platform: Platform) => {
    setSelectedPlatforms(prev =>
      prev.includes(platform)
        ? prev.filter(p => p !== platform)
        : [...prev, platform]
    )
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)

    if (!mediaUrl.trim()) {
      setError('Please enter a media URL')
      return
    }

    if (selectedPlatforms.length === 0) {
      setError('Please select at least one platform')
      return
    }

    setIsSubmitting(true)

    try {
      await postService.createPost({
        platforms: selectedPlatforms,
        media_url: mediaUrl,
        caption: caption,
      })

      // Reset form
      setMediaUrl('')
      setCaption('')
      setSelectedPlatforms([])

      // Notify parent
      onPostCreated()
    } catch (err: any) {
      console.error('Post creation error:', err)
      setError(err.response?.data?.error || 'Failed to create post')
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <div className="post-form-container">
      <h2>Create New Post</h2>
      <form onSubmit={handleSubmit} className="post-form">
        <div className="form-group">
          <label htmlFor="mediaUrl">Media URL *</label>
          <input
            type="url"
            id="mediaUrl"
            value={mediaUrl}
            onChange={(e) => setMediaUrl(e.target.value)}
            placeholder="https://example.com/media.jpg or video.mp4"
            disabled={isSubmitting}
            required
          />
          <small>
            Enter a publicly accessible media URL
            {detectedMediaType === 'image' && (
              <span className="media-type-badge image"> - Photo detected</span>
            )}
            {detectedMediaType === 'video' && (
              <span className="media-type-badge video"> - Video detected</span>
            )}
          </small>
          <small className="supported-formats">
            Supported: JPG, PNG, GIF, WEBP (photos) | MP4, MOV, WEBM (videos)
          </small>
        </div>

        <div className="form-group">
          <label htmlFor="caption">Caption</label>
          <textarea
            id="caption"
            value={caption}
            onChange={(e) => setCaption(e.target.value)}
            placeholder="Add a caption with hashtags..."
            rows={4}
            disabled={isSubmitting}
            maxLength={2200}
          />
          <small>{caption.length} / 2200 characters</small>
        </div>

        <div className="form-group">
          <label>Select Platforms *</label>
          <div className="platform-checkboxes">
            {connectedPlatforms.filter(conn => conn.is_active).length === 0 ? (
              <p className="no-platforms-message">
                No platforms connected. Please connect at least one platform to create posts.
              </p>
            ) : (
              connectedPlatforms
                .filter(conn => conn.is_active)
                .map(connection => (
                  <label key={connection.platform} className="platform-checkbox">
                    <input
                      type="checkbox"
                      checked={selectedPlatforms.includes(connection.platform)}
                      onChange={() => handlePlatformToggle(connection.platform)}
                      disabled={isSubmitting}
                    />
                    <span className="platform-name">
                      {connection.platform === 'tiktok' && '🎵'}
                      {connection.platform === 'x' && '𝕏'}
                      {connection.platform === 'instagram' && '📷'}
                      {connection.platform === 'youtube' && '▶️'}
                      {' '}
                      {connection.platform.charAt(0).toUpperCase() + connection.platform.slice(1)}
                    </span>
                  </label>
                ))
            )}
          </div>
          <small>
            {selectedPlatforms.length > 0
              ? `Selected: ${selectedPlatforms.map(p => p.toUpperCase()).join(', ')}`
              : 'Select one or more platforms to post to'}
          </small>
        </div>

        {error && (
          <div className="error-message">
            {error}
          </div>
        )}

        <button
          type="submit"
          className="submit-button"
          disabled={isSubmitting || connectedPlatforms.filter(conn => conn.is_active).length === 0}
        >
          {isSubmitting
            ? 'Posting...'
            : selectedPlatforms.length > 1
            ? `Post to ${selectedPlatforms.length} Platforms`
            : selectedPlatforms.length === 1
            ? `Post to ${selectedPlatforms[0].toUpperCase()}`
            : 'Select Platforms to Post'}
        </button>
      </form>

      <style>{`
        .media-type-badge {
          font-weight: 600;
          padding: 2px 6px;
          border-radius: 4px;
          margin-left: 4px;
        }

        .media-type-badge.image {
          background: #e3f2fd;
          color: #1565c0;
        }

        .media-type-badge.video {
          background: #fce4ec;
          color: #c62828;
        }

        .supported-formats {
          display: block;
          color: #888;
          font-size: 11px;
          margin-top: 4px;
        }

        .platform-checkboxes {
          display: flex;
          flex-direction: column;
          gap: 12px;
          padding: 12px;
          background: #f8f9fa;
          border-radius: 8px;
          margin-top: 8px;
        }

        .platform-checkbox {
          display: flex;
          align-items: center;
          gap: 10px;
          padding: 10px 12px;
          background: white;
          border: 2px solid #e0e0e0;
          border-radius: 6px;
          cursor: pointer;
          transition: all 0.2s ease;
          user-select: none;
        }

        .platform-checkbox:hover {
          border-color: #4CAF50;
          background: #f1f8f4;
        }

        .platform-checkbox input[type="checkbox"] {
          width: 18px;
          height: 18px;
          cursor: pointer;
        }

        .platform-checkbox input[type="checkbox"]:checked + .platform-name {
          font-weight: 600;
          color: #2d7a3d;
        }

        .platform-name {
          font-size: 15px;
          color: #333;
          transition: all 0.2s ease;
        }

        .no-platforms-message {
          color: #666;
          font-style: italic;
          padding: 12px;
          text-align: center;
          background: #fff3cd;
          border-radius: 6px;
          margin: 0;
        }

        @media (min-width: 768px) {
          .platform-checkboxes {
            flex-direction: row;
            flex-wrap: wrap;
          }

          .platform-checkbox {
            flex: 1;
            min-width: 150px;
          }
        }
      `}</style>
    </div>
  )
}

export default PostForm
