import { useState, useMemo } from 'react'
import { postService } from '../../services/postService'
import { useAuth } from '../../context/AuthContext'
import { Platform } from '../../types/user'
import { TikTokPrivacyLevel, TikTokSettings } from '../../types/post'

interface PostFormProps {
  onPostCreated: () => void
}

type MediaType = 'video' | 'image' | 'unknown'

const imageExtensions = ['.jpg', '.jpeg', '.png', '.gif', '.webp', '.bmp', '.heic', '.heif']
const videoExtensions = ['.mp4', '.mov', '.webm', '.avi', '.mkv', '.m4v', '.ts', '.3gp']

const detectMediaType = (url: string): MediaType => {
  if (!url) return 'unknown'
  const lowerUrl = url.toLowerCase()

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
    if (imageExtensions.some(ext => lowerUrl.endsWith(ext))) {
      return 'image'
    }
    if (videoExtensions.some(ext => lowerUrl.endsWith(ext))) {
      return 'video'
    }
  }

  return 'unknown'
}

const PRIVACY_OPTIONS: { value: TikTokPrivacyLevel; label: string; description: string }[] = [
  { value: 'PUBLIC_TO_EVERYONE', label: 'Public', description: 'Anyone can view this video' },
  { value: 'MUTUAL_FOLLOW_FRIENDS', label: 'Friends', description: 'Only mutual followers can view' },
  { value: 'FOLLOWER_OF_CREATOR', label: 'Followers', description: 'Only your followers can view' },
  { value: 'SELF_ONLY', label: 'Only Me', description: 'Only you can view this video' },
]

const PostForm = ({ onPostCreated }: PostFormProps) => {
  const { connectedPlatforms } = useAuth()
  const [mediaUrl, setMediaUrl] = useState('')
  const [caption, setCaption] = useState('')
  const [title, setTitle] = useState('')
  const [selectedPlatforms, setSelectedPlatforms] = useState<Platform[]>([])
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [error, setError] = useState<string | null>(null)

  // TikTok-specific settings (required by TikTok UX Guidelines)
  const [privacyLevel, setPrivacyLevel] = useState<TikTokPrivacyLevel | ''>('')
  const [allowComment, setAllowComment] = useState(false)
  const [allowDuet, setAllowDuet] = useState(false)
  const [allowStitch, setAllowStitch] = useState(false)
  const [isBrandContent, setIsBrandContent] = useState(false)
  const [isBrandOrganic, setIsBrandOrganic] = useState(false)
  const [agreedToTerms, setAgreedToTerms] = useState(false)

  const detectedMediaType = useMemo(() => detectMediaType(mediaUrl), [mediaUrl])

  // Get TikTok connection info to display user's nickname
  const tiktokConnection = connectedPlatforms.find(c => c.platform === 'tiktok' && c.is_active)
  const isTikTokSelected = selectedPlatforms.includes('tiktok')

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

    // TikTok-specific validation
    if (isTikTokSelected) {
      if (!privacyLevel) {
        setError('Please select a privacy level for TikTok')
        return
      }
      if (!agreedToTerms) {
        setError('Please agree to the terms before posting to TikTok')
        return
      }
      // Branded content cannot be private
      if (isBrandOrganic && privacyLevel === 'SELF_ONLY') {
        setError('Branded content cannot be set to private')
        return
      }
    }

    setIsSubmitting(true)

    try {
      // Build TikTok settings if TikTok is selected
      const tiktokSettings: TikTokSettings | undefined = isTikTokSelected
        ? {
            title: title || undefined,
            privacy_level: privacyLevel as TikTokPrivacyLevel,
            allow_comment: allowComment,
            allow_duet: allowDuet,
            allow_stitch: allowStitch,
            is_brand_content: isBrandContent,
            is_brand_organic: isBrandOrganic,
          }
        : undefined

      await postService.createPost({
        platforms: selectedPlatforms,
        media_url: mediaUrl,
        caption: caption,
        tiktok_settings: tiktokSettings,
      })

      // Reset form
      setMediaUrl('')
      setCaption('')
      setTitle('')
      setSelectedPlatforms([])
      setPrivacyLevel('')
      setAllowComment(false)
      setAllowDuet(false)
      setAllowStitch(false)
      setIsBrandContent(false)
      setIsBrandOrganic(false)
      setAgreedToTerms(false)

      onPostCreated()
    } catch (err: any) {
      console.error('Post creation error:', err)
      setError(err.response?.data?.error || 'Failed to create post')
    } finally {
      setIsSubmitting(false)
    }
  }

  const activePlatforms = connectedPlatforms.filter(conn => conn.is_active)

  return (
    <div className="bg-white rounded-2xl border border-gray-200 overflow-hidden">
      {/* Header */}
      <div className="px-6 py-4 border-b border-gray-100 bg-gradient-to-r from-indigo-50 to-purple-50">
        <h2 className="text-xl font-semibold text-gray-900">Create New Post</h2>
        <p className="text-sm text-gray-600 mt-1">Share your content across multiple platforms</p>
      </div>

      <form onSubmit={handleSubmit} className="p-6 space-y-6">
        {/* Media URL */}
        <div className="space-y-2">
          <label htmlFor="mediaUrl" className="block text-sm font-medium text-gray-700">
            Media URL <span className="text-red-500">*</span>
          </label>
          <input
            type="url"
            id="mediaUrl"
            value={mediaUrl}
            onChange={(e) => setMediaUrl(e.target.value)}
            placeholder="https://example.com/video.mp4"
            disabled={isSubmitting}
            required
            className="w-full px-4 py-3 border-2 border-gray-200 rounded-xl focus:border-indigo-500 focus:ring-0 transition-colors"
          />
          <div className="flex items-center gap-2 text-sm text-gray-500">
            <span>Supported: MP4, MOV, WEBM (videos) | JPG, PNG, GIF (images)</span>
            {detectedMediaType === 'video' && (
              <span className="px-2 py-0.5 bg-purple-100 text-purple-700 rounded-full text-xs font-medium">
                Video detected
              </span>
            )}
            {detectedMediaType === 'image' && (
              <span className="px-2 py-0.5 bg-blue-100 text-blue-700 rounded-full text-xs font-medium">
                Image detected
              </span>
            )}
          </div>
        </div>

        {/* Title (for TikTok) */}
        <div className="space-y-2">
          <label htmlFor="title" className="block text-sm font-medium text-gray-700">
            Title {isTikTokSelected && <span className="text-gray-400">(for TikTok)</span>}
          </label>
          <input
            type="text"
            id="title"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="Give your video a title..."
            disabled={isSubmitting}
            maxLength={150}
            className="w-full px-4 py-3 border-2 border-gray-200 rounded-xl focus:border-indigo-500 focus:ring-0 transition-colors"
          />
          <p className="text-xs text-gray-400">{title.length} / 150 characters</p>
        </div>

        {/* Caption */}
        <div className="space-y-2">
          <label htmlFor="caption" className="block text-sm font-medium text-gray-700">
            Caption
          </label>
          <textarea
            id="caption"
            value={caption}
            onChange={(e) => setCaption(e.target.value)}
            placeholder="Write your caption with hashtags..."
            rows={3}
            disabled={isSubmitting}
            maxLength={2200}
            className="w-full px-4 py-3 border-2 border-gray-200 rounded-xl focus:border-indigo-500 focus:ring-0 transition-colors resize-none"
          />
          <p className="text-xs text-gray-400">{caption.length} / 2200 characters</p>
        </div>

        {/* Platform Selection */}
        <div className="space-y-3">
          <label className="block text-sm font-medium text-gray-700">
            Select Platforms <span className="text-red-500">*</span>
          </label>

          {activePlatforms.length === 0 ? (
            <div className="bg-amber-50 border-2 border-amber-200 rounded-xl p-4 text-center">
              <p className="text-amber-800">No platforms connected. Please connect at least one platform.</p>
            </div>
          ) : (
            <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
              {activePlatforms.map(connection => (
                <button
                  key={connection.platform}
                  type="button"
                  onClick={() => handlePlatformToggle(connection.platform)}
                  disabled={isSubmitting}
                  className={`flex flex-col items-center gap-2 p-4 rounded-xl border-2 transition-all ${
                    selectedPlatforms.includes(connection.platform)
                      ? 'border-indigo-500 bg-indigo-50'
                      : 'border-gray-200 hover:border-gray-300'
                  }`}
                >
                  <div className={`w-10 h-10 rounded-lg flex items-center justify-center text-xl ${
                    selectedPlatforms.includes(connection.platform)
                      ? 'bg-gradient-to-br from-indigo-500 to-purple-500 text-white'
                      : 'bg-gray-100'
                  }`}>
                    {connection.platform === 'tiktok' && '🎵'}
                    {connection.platform === 'x' && '𝕏'}
                    {connection.platform === 'instagram' && '📷'}
                    {connection.platform === 'youtube' && '▶️'}
                  </div>
                  <span className="text-sm font-medium text-gray-700 capitalize">
                    {connection.platform}
                  </span>
                  {selectedPlatforms.includes(connection.platform) && (
                    <svg className="w-5 h-5 text-indigo-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                    </svg>
                  )}
                </button>
              ))}
            </div>
          )}
        </div>

        {/* TikTok-specific settings (shown when TikTok is selected) */}
        {isTikTokSelected && tiktokConnection && (
          <div className="bg-gray-50 rounded-xl p-5 space-y-5 border border-gray-200">
            <div className="flex items-center gap-3 pb-4 border-b border-gray-200">
              <div className="w-10 h-10 rounded-full bg-gradient-to-br from-indigo-500 to-purple-500 flex items-center justify-center text-white">
                🎵
              </div>
              <div>
                <p className="font-medium text-gray-900">TikTok Settings</p>
                <p className="text-sm text-gray-500">
                  Posting as <span className="font-medium text-indigo-600">@{tiktokConnection.username || tiktokConnection.display_name}</span>
                </p>
              </div>
            </div>

            {/* Privacy Level (Required - no default) */}
            <div className="space-y-2">
              <label className="block text-sm font-medium text-gray-700">
                Privacy Level <span className="text-red-500">*</span>
              </label>
              <select
                value={privacyLevel}
                onChange={(e) => setPrivacyLevel(e.target.value as TikTokPrivacyLevel | '')}
                required={isTikTokSelected}
                className="w-full px-4 py-3 border-2 border-gray-200 rounded-xl focus:border-indigo-500 focus:ring-0 transition-colors bg-white"
              >
                <option value="">-- Select privacy level --</option>
                {PRIVACY_OPTIONS.map(option => (
                  <option key={option.value} value={option.value}>
                    {option.label} - {option.description}
                  </option>
                ))}
              </select>
            </div>

            {/* Interaction Settings (all unchecked by default) */}
            <div className="space-y-3">
              <label className="block text-sm font-medium text-gray-700">Interaction Settings</label>
              <div className="space-y-2">
                <label className="flex items-center gap-3 cursor-pointer">
                  <input
                    type="checkbox"
                    checked={allowComment}
                    onChange={(e) => setAllowComment(e.target.checked)}
                    className="w-5 h-5 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500"
                  />
                  <span className="text-sm text-gray-700">Allow Comments</span>
                </label>
                <label className="flex items-center gap-3 cursor-pointer">
                  <input
                    type="checkbox"
                    checked={allowDuet}
                    onChange={(e) => setAllowDuet(e.target.checked)}
                    className="w-5 h-5 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500"
                  />
                  <span className="text-sm text-gray-700">Allow Duet</span>
                </label>
                <label className="flex items-center gap-3 cursor-pointer">
                  <input
                    type="checkbox"
                    checked={allowStitch}
                    onChange={(e) => setAllowStitch(e.target.checked)}
                    className="w-5 h-5 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500"
                  />
                  <span className="text-sm text-gray-700">Allow Stitch</span>
                </label>
              </div>
            </div>

            {/* Commercial Content Disclosure */}
            <div className="space-y-3">
              <label className="block text-sm font-medium text-gray-700">Commercial Content Disclosure</label>
              <div className="space-y-2">
                <label className="flex items-center gap-3 cursor-pointer">
                  <input
                    type="checkbox"
                    checked={isBrandContent}
                    onChange={(e) => setIsBrandContent(e.target.checked)}
                    className="w-5 h-5 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500"
                  />
                  <div>
                    <span className="text-sm text-gray-700 font-medium">Your Brand</span>
                    <p className="text-xs text-gray-500">This content promotes yourself or your own business</p>
                  </div>
                </label>
                <label className="flex items-center gap-3 cursor-pointer">
                  <input
                    type="checkbox"
                    checked={isBrandOrganic}
                    onChange={(e) => setIsBrandOrganic(e.target.checked)}
                    className="w-5 h-5 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500"
                  />
                  <div>
                    <span className="text-sm text-gray-700 font-medium">Branded Content</span>
                    <p className="text-xs text-gray-500">This is a paid partnership with a brand</p>
                  </div>
                </label>
              </div>
              {isBrandOrganic && privacyLevel === 'SELF_ONLY' && (
                <p className="text-sm text-red-500 mt-2">
                  ⚠️ Branded content cannot be set to private
                </p>
              )}
            </div>

            {/* Consent Statement */}
            <div className="pt-3 border-t border-gray-200">
              <label className="flex items-start gap-3 cursor-pointer">
                <input
                  type="checkbox"
                  checked={agreedToTerms}
                  onChange={(e) => setAgreedToTerms(e.target.checked)}
                  required={isTikTokSelected}
                  className="w-5 h-5 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500 mt-0.5"
                />
                <span className="text-sm text-gray-600">
                  By posting, you agree to TikTok's{' '}
                  <a
                    href="https://www.tiktok.com/legal/terms-of-service"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-indigo-600 hover:underline"
                  >
                    Terms of Service
                  </a>{' '}
                  and{' '}
                  <a
                    href="https://www.tiktok.com/legal/music-usage-confirmation"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-indigo-600 hover:underline"
                  >
                    Music Usage Confirmation
                  </a>
                  {isBrandOrganic && (
                    <>
                      {' '}and{' '}
                      <a
                        href="https://www.tiktok.com/creators/creator-portal/en-us/getting-paid-to-create/branded-content-policy/"
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-indigo-600 hover:underline"
                      >
                        Branded Content Policy
                      </a>
                    </>
                  )}
                  .
                </span>
              </label>
            </div>
          </div>
        )}

        {/* Error Message */}
        {error && (
          <div className="bg-red-50 border-2 border-red-200 rounded-xl p-4">
            <p className="text-red-700 text-sm">{error}</p>
          </div>
        )}

        {/* Submit Button */}
        <button
          type="submit"
          disabled={isSubmitting || activePlatforms.length === 0}
          className={`w-full py-4 rounded-xl font-semibold text-white transition-all ${
            isSubmitting || activePlatforms.length === 0
              ? 'bg-gray-300 cursor-not-allowed'
              : 'bg-gradient-to-r from-indigo-600 to-purple-600 hover:shadow-lg hover:scale-[1.02]'
          }`}
        >
          {isSubmitting ? (
            <span className="flex items-center justify-center gap-2">
              <svg className="animate-spin h-5 w-5" viewBox="0 0 24 24">
                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" />
                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
              </svg>
              Posting...
            </span>
          ) : selectedPlatforms.length > 1 ? (
            `Publish to ${selectedPlatforms.length} Platforms`
          ) : selectedPlatforms.length === 1 ? (
            `Publish to ${selectedPlatforms[0].charAt(0).toUpperCase() + selectedPlatforms[0].slice(1)}`
          ) : (
            'Select Platforms to Publish'
          )}
        </button>
      </form>
    </div>
  )
}

export default PostForm
