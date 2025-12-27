import { useState } from 'react'
import { postService } from '../../services/postService'
import { useAuth } from '../../context/AuthContext'
import { Platform } from '../../types/user'

interface PostFormProps {
  onPostCreated: () => void
}

const PostForm = ({ onPostCreated }: PostFormProps) => {
  const { connectedPlatforms } = useAuth()
  const [videoUrl, setVideoUrl] = useState('')
  const [caption, setCaption] = useState('')
  const [selectedPlatforms, setSelectedPlatforms] = useState<Platform[]>([])
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [error, setError] = useState<string | null>(null)

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

    if (!videoUrl.trim()) {
      setError('Please enter a video URL')
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
        media_url: videoUrl,
        caption: caption,
      })

      // Reset form
      setVideoUrl('')
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
          <label htmlFor="videoUrl">Video URL *</label>
          <input
            type="url"
            id="videoUrl"
            value={videoUrl}
            onChange={(e) => setVideoUrl(e.target.value)}
            placeholder="https://example.com/video.mp4"
            disabled={isSubmitting}
            required
          />
          <small>Enter a publicly accessible video URL (MP4, MOV, or WEBM)</small>
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
