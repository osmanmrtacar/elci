import { useState, useEffect, useMemo, useCallback } from 'react'
import { postService } from '../../services/postService'
import { useAuth } from '../../context/AuthContext'
import { Platform } from '../../types/user'
import { TikTokPrivacyLevel, TikTokSettings, TikTokCreatorInfo } from '../../types/post'

interface PostFormProps {
  onPostCreated: () => void
}

type MediaType = 'video' | 'image' | 'unknown'

const imageExtensions = ['.jpg', '.jpeg', '.png', '.gif', '.webp', '.bmp', '.heic', '.heif']
const videoExtensions = ['.mp4', '.mov', '.webm', '.avi', '.mkv', '.m4v', '.ts', '.3gp']
const allMediaExtensions = [...imageExtensions, ...videoExtensions]

const validateMediaUrl = (url: string): string | null => {
  if (!url.trim()) return null // empty is ok, filtered out later
  try {
    const parsed = new URL(url)
    if (parsed.protocol !== 'https:' && parsed.protocol !== 'http:') {
      return 'URL must start with http:// or https://'
    }
    const path = parsed.pathname.toLowerCase()
    if (!allMediaExtensions.some(ext => path.endsWith(ext))) {
      return 'URL must point to a supported media file (MP4, MOV, JPG, PNG, GIF, etc.)'
    }
    return null
  } catch {
    return 'Please enter a valid URL'
  }
}

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

const PRIVACY_LABEL_MAP: Record<string, { label: string; description: string }> = {
  PUBLIC_TO_EVERYONE: { label: 'Public', description: 'Anyone can view this video' },
  MUTUAL_FOLLOW_FRIENDS: { label: 'Friends', description: 'Only mutual followers can view' },
  FOLLOWER_OF_CREATOR: { label: 'Followers', description: 'Only your followers can view' },
  SELF_ONLY: { label: 'Only Me', description: 'Only you can view this video' },
}

const PostForm = ({ onPostCreated }: PostFormProps) => {
  const { connectedPlatforms } = useAuth()
  const [mediaUrls, setMediaUrls] = useState<string[]>([''])
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
  const [tiktokSettingsOpen, setTiktokSettingsOpen] = useState(true)
  const [discloseContent, setDiscloseContent] = useState(false)
  const [isBrandContent, setIsBrandContent] = useState(false)
  const [isBrandOrganic, setIsBrandOrganic] = useState(false)
  const [agreedToTerms, setAgreedToTerms] = useState(false)
  const [autoAddMusic, setAutoAddMusic] = useState(false)
  const [directPost, setDirectPost] = useState(true)

  // TikTok Creator Info state
  const [creatorInfo, setCreatorInfo] = useState<TikTokCreatorInfo | null>(null)
  const [creatorInfoLoading, setCreatorInfoLoading] = useState(false)
  const [creatorInfoError, setCreatorInfoError] = useState<string | null>(null)

  // Media URL validation and preview error tracking
  const [urlErrors, setUrlErrors] = useState<Record<number, string>>({})
  const [previewErrors, setPreviewErrors] = useState<Record<number, boolean>>({})
  const [videoDuration, setVideoDuration] = useState<number | null>(null)

  const isTikTokSelected = selectedPlatforms.includes('tiktok')

  // Fetch creator info when TikTok is selected
  const fetchCreatorInfo = useCallback(async () => {
    setCreatorInfoLoading(true)
    setCreatorInfoError(null)
    try {
      const info = await postService.getTikTokCreatorInfo()
      setCreatorInfo(info)
      // Force-disable interactions that are disabled by creator settings
      if (info.comment_disabled) setAllowComment(false)
      if (info.duet_disabled) setAllowDuet(false)
      if (info.stitch_disabled) setAllowStitch(false)
    } catch {
      setCreatorInfoError('Failed to load TikTok creator settings')
    } finally {
      setCreatorInfoLoading(false)
    }
  }, [])

  useEffect(() => {
    if (isTikTokSelected) {
      fetchCreatorInfo()
    } else {
      setCreatorInfo(null)
      setCreatorInfoError(null)
    }
  }, [isTikTokSelected, fetchCreatorInfo])

  // Build dynamic privacy options from creator info
  const privacyOptions = useMemo(() => {
    if (!creatorInfo?.privacy_level_options) return []
    return creatorInfo.privacy_level_options.map(value => ({
      value: value as TikTokPrivacyLevel,
      label: PRIVACY_LABEL_MAP[value]?.label ?? value,
      description: PRIVACY_LABEL_MAP[value]?.description ?? '',
    }))
  }, [creatorInfo])

  // Detect media type from first URL
  const detectedMediaType = useMemo(() => detectMediaType(mediaUrls[0] || ''), [mediaUrls])

  // Check if this is a carousel (multiple images)
  const isCarousel = mediaUrls.filter(url => url.trim()).length >= 2

  // Get TikTok connection info to display user's nickname
  const tiktokConnection = connectedPlatforms.find(c => c.platform === 'tiktok' && c.is_active)

  const handlePlatformToggle = (platform: Platform) => {
    setSelectedPlatforms(prev =>
      prev.includes(platform)
        ? prev.filter(p => p !== platform)
        : [...prev, platform]
    )
  }

  // Handle adding a new media URL input
  const handleAddMediaUrl = () => {
    if (mediaUrls.length < 10) { // Instagram carousel max is 10
      setMediaUrls([...mediaUrls, ''])
    }
  }

  // Handle removing a media URL input
  const handleRemoveMediaUrl = (index: number) => {
    if (mediaUrls.length > 1) {
      setMediaUrls(mediaUrls.filter((_, i) => i !== index))
    }
  }

  // Handle updating a media URL
  const handleMediaUrlChange = (index: number, value: string) => {
    const newUrls = [...mediaUrls]
    newUrls[index] = value
    setMediaUrls(newUrls)
    // Clear errors when URL changes
    setUrlErrors(prev => {
      const next = { ...prev }
      delete next[index]
      return next
    })
    setPreviewErrors(prev => {
      const next = { ...prev }
      delete next[index]
      return next
    })
    if (index === 0) setVideoDuration(null)
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)

    // Filter out empty URLs
    const validUrls = mediaUrls.filter(url => url.trim())

    if (validUrls.length === 0) {
      setError('Please enter at least one media URL')
      return
    }

    // Validate all media URLs
    const errors: Record<number, string> = {}
    mediaUrls.forEach((url, index) => {
      const err = validateMediaUrl(url)
      if (err) errors[index] = err
    })
    if (Object.keys(errors).length > 0) {
      setUrlErrors(errors)
      setError('Please fix the invalid media URLs')
      return
    }

    if (selectedPlatforms.length === 0) {
      setError('Please select at least one platform')
      return
    }

    // Platform-specific validation for multiple media
    if (validUrls.length > 1) {
      if (selectedPlatforms.includes('x') && validUrls.length > 4) {
        setError('X (Twitter) only supports up to 4 images per post')
        return
      }
      if (selectedPlatforms.includes('instagram') && validUrls.length > 10) {
        setError('Instagram only supports up to 10 items per carousel')
        return
      }
      if (selectedPlatforms.includes('instagram') && validUrls.length < 2) {
        setError('Instagram carousel requires at least 2 items')
        return
      }
      if (selectedPlatforms.includes('tiktok') && validUrls.length > 35) {
        setError('TikTok only supports up to 35 photos per post')
        return
      }
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
      // Video duration must not exceed creator's max
      if (creatorInfo && videoDuration && detectedMediaType === 'video' && videoDuration > creatorInfo.max_video_post_duration_sec) {
        setError(`Video duration (${videoDuration}s) exceeds the maximum allowed (${creatorInfo.max_video_post_duration_sec}s)`)
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
            auto_add_music: autoAddMusic,
            direct_post: directPost,
          }
        : undefined

      await postService.createPost({
        platforms: selectedPlatforms,
        media_url: validUrls[0], // Primary URL for backwards compatibility
        media_urls: validUrls,   // All URLs for carousel/multi-image
        caption: caption,
        tiktok_settings: tiktokSettings,
      })

      // Reset form
      setMediaUrls([''])
      setCaption('')
      setTitle('')
      setSelectedPlatforms([])
      setPrivacyLevel('')
      setAllowComment(false)
      setAllowDuet(false)
      setAllowStitch(false)
      setDiscloseContent(false)
      setIsBrandContent(false)
      setIsBrandOrganic(false)
      setAgreedToTerms(false)
      setAutoAddMusic(false)
      setDirectPost(true)
      setUrlErrors({})
      setPreviewErrors({})
      setVideoDuration(null)

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
        {/* Media URLs */}
        <div className="space-y-3">
          <div className="flex items-center justify-between">
            <label className="block text-sm font-medium text-gray-700">
              Media URLs <span className="text-red-500">*</span>
            </label>
            <div className="flex items-center gap-2">
              {isCarousel && (
                <span className="px-2 py-0.5 bg-green-100 text-green-700 rounded-full text-xs font-medium">
                  Carousel ({mediaUrls.filter(u => u.trim()).length} items)
                </span>
              )}
              {detectedMediaType === 'video' && (
                <span className="px-2 py-0.5 bg-purple-100 text-purple-700 rounded-full text-xs font-medium">
                  Video
                </span>
              )}
              {detectedMediaType === 'image' && (
                <span className="px-2 py-0.5 bg-blue-100 text-blue-700 rounded-full text-xs font-medium">
                  Image
                </span>
              )}
            </div>
          </div>

          {mediaUrls.map((url, index) => (
            <div key={index} className="space-y-1">
              <div className="flex gap-2">
              <input
                type="url"
                value={url}
                onChange={(e) => handleMediaUrlChange(index, e.target.value)}
                onBlur={() => {
                  const err = validateMediaUrl(url)
                  setUrlErrors(prev => {
                    if (err) return { ...prev, [index]: err }
                    const next = { ...prev }
                    delete next[index]
                    return next
                  })
                }}
                placeholder={`https://example.com/${index === 0 ? 'video.mp4' : `image${index + 1}.jpg`}`}
                disabled={isSubmitting}
                className={`flex-1 px-4 py-3 border-2 rounded-xl focus:ring-0 transition-colors ${
                  urlErrors[index] ? 'border-red-300 focus:border-red-500' : 'border-gray-200 focus:border-indigo-500'
                }`}
              />
              {mediaUrls.length > 1 && (
                <button
                  type="button"
                  onClick={() => handleRemoveMediaUrl(index)}
                  disabled={isSubmitting}
                  className="px-3 py-3 text-red-500 hover:bg-red-50 rounded-xl transition-colors"
                  title="Remove"
                >
                  <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              )}
              </div>
              {urlErrors[index] && (
                <p className="text-xs text-red-500 pl-1">{urlErrors[index]}</p>
              )}
            </div>
          ))}

          <div className="flex items-center justify-between">
            <button
              type="button"
              onClick={handleAddMediaUrl}
              disabled={isSubmitting || mediaUrls.length >= 10}
              className="flex items-center gap-1 text-sm text-indigo-600 hover:text-indigo-800 disabled:text-gray-400"
            >
              <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
              </svg>
              Add another media
            </button>
            <span className="text-xs text-gray-400">
              {mediaUrls.length}/10 • Supports: MP4, MOV, JPG, PNG, GIF
            </span>
          </div>
        </div>

        {/* Media Preview */}
        {mediaUrls.some(url => url.trim()) && (
          <div className="space-y-2">
            <label className="block text-sm font-medium text-gray-700">Media Preview</label>
            <div className={`${isCarousel ? 'grid grid-cols-3 sm:grid-cols-4 gap-2' : ''}`}>
              {mediaUrls.map((url, index) => {
                if (!url.trim() || previewErrors[index]) return null
                const type = detectMediaType(url)
                if (type === 'video') {
                  return (
                    <video
                      key={index}
                      src={url}
                      controls
                      className="w-full max-h-64 rounded-xl border border-gray-200 object-contain bg-black"
                      onLoadedMetadata={(e) => {
                        if (index === 0) setVideoDuration(Math.round(e.currentTarget.duration))
                      }}
                      onError={() => setPreviewErrors(prev => ({ ...prev, [index]: true }))}
                    />
                  )
                }
                return (
                  <img
                    key={index}
                    src={url}
                    alt={`Media ${index + 1}`}
                    className={`rounded-xl border border-gray-200 object-cover ${isCarousel ? 'w-full h-24' : 'max-h-48'}`}
                    onError={() => setPreviewErrors(prev => ({ ...prev, [index]: true }))}
                  />
                )
              })}
            </div>
          </div>
        )}

        {/* Max video duration info from Creator Info API */}
        {isTikTokSelected && creatorInfo && detectedMediaType === 'video' && (
          <div className={`flex items-center gap-2 px-4 py-2 rounded-xl text-sm ${
            videoDuration && videoDuration > creatorInfo.max_video_post_duration_sec
              ? 'bg-red-50 border border-red-200 text-red-700'
              : 'bg-blue-50 border border-blue-200 text-blue-700'
          }`}>
            <svg className="w-4 h-4 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span>
              Maximum video duration: {creatorInfo.max_video_post_duration_sec}s
              {videoDuration != null && (
                <> — Your video: {videoDuration}s
                  {videoDuration > creatorInfo.max_video_post_duration_sec && ' (exceeds limit)'}
                </>
              )}
            </span>
          </div>
        )}

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
                    {connection.platform === 'tiktok' && (
                      <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                        <path d="M19.59 6.69a4.83 4.83 0 0 1-3.77-4.25V2h-3.45v13.67a2.89 2.89 0 0 1-5.2 1.74 2.89 2.89 0 0 1 2.31-4.64 2.93 2.93 0 0 1 .88.13V9.4a6.84 6.84 0 0 0-1-.05A6.33 6.33 0 0 0 5 20.1a6.34 6.34 0 0 0 10.86-4.43v-7a8.16 8.16 0 0 0 4.77 1.52v-3.4a4.85 4.85 0 0 1-1-.1z" />
                      </svg>
                    )}
                    {connection.platform === 'x' && '𝕏'}
                    {connection.platform === 'instagram' && (
                      <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                        <path d="M7.8 2h8.4C19.4 2 22 4.6 22 7.8v8.4a5.8 5.8 0 0 1-5.8 5.8H7.8C4.6 22 2 19.4 2 16.2V7.8A5.8 5.8 0 0 1 7.8 2m-.2 2A3.6 3.6 0 0 0 4 7.6v8.8C4 18.39 5.61 20 7.6 20h8.8a3.6 3.6 0 0 0 3.6-3.6V7.6C20 5.61 18.39 4 16.4 4H7.6m9.65 1.5a1.25 1.25 0 0 1 1.25 1.25A1.25 1.25 0 0 1 17.25 8 1.25 1.25 0 0 1 16 6.75a1.25 1.25 0 0 1 1.25-1.25M12 7a5 5 0 0 1 5 5 5 5 0 0 1-5 5 5 5 0 0 1-5-5 5 5 0 0 1 5-5m0 2a3 3 0 0 0-3 3 3 3 0 0 0 3 3 3 3 0 0 0 3-3 3 3 0 0 0-3-3z" />
                      </svg>
                    )}
                    {connection.platform === 'youtube' && '▶️'}
                  </div>
                  <div className="text-center">
                    <span className="text-sm font-medium text-gray-700 capitalize block">
                      {connection.platform}
                    </span>
                    {(connection.username || connection.display_name) && (
                      <span className="text-xs text-gray-400 truncate block max-w-[100px]">
                        @{connection.username || connection.display_name}
                      </span>
                    )}
                  </div>
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
          <div className="bg-gray-50 rounded-xl border border-gray-200">
            <button
              type="button"
              onClick={() => setTiktokSettingsOpen(prev => !prev)}
              className="w-full flex items-center gap-3 p-5 cursor-pointer"
            >
              <div className="w-10 h-10 rounded-full bg-gradient-to-br from-indigo-500 to-purple-500 flex items-center justify-center text-white flex-shrink-0">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M19.59 6.69a4.83 4.83 0 0 1-3.77-4.25V2h-3.45v13.67a2.89 2.89 0 0 1-5.2 1.74 2.89 2.89 0 0 1 2.31-4.64 2.93 2.93 0 0 1 .88.13V9.4a6.84 6.84 0 0 0-1-.05A6.33 6.33 0 0 0 5 20.1a6.34 6.34 0 0 0 10.86-4.43v-7a8.16 8.16 0 0 0 4.77 1.52v-3.4a4.85 4.85 0 0 1-1-.1z" />
                </svg>
              </div>
              <div className="flex-1 text-left">
                <p className="font-medium text-gray-900">TikTok Settings</p>
                <p className="text-sm text-gray-500">
                  Posting as <span className="font-medium text-indigo-600">@{tiktokConnection.username || tiktokConnection.display_name}</span>
                </p>
              </div>
              <svg className={`w-5 h-5 text-gray-400 transition-transform ${tiktokSettingsOpen ? 'rotate-180' : ''}`} fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
              </svg>
            </button>

            {tiktokSettingsOpen && (
            <div className="px-5 pb-5 space-y-5 border-t border-gray-200 pt-5">

            {/* Direct Post vs Send to Inbox (only for video posts) */}
            {detectedMediaType !== 'image' && (
              <div className="space-y-2">
                <label className="block text-sm font-medium text-gray-700">Publish Mode</label>
                <div className="flex gap-3">
                  <button
                    type="button"
                    onClick={() => setDirectPost(true)}
                    className={`flex-1 px-4 py-3 rounded-xl border-2 text-sm font-medium transition-all ${
                      directPost
                        ? 'border-indigo-500 bg-indigo-50 text-indigo-700'
                        : 'border-gray-200 text-gray-600 hover:border-gray-300'
                    }`}
                  >
                    Direct Post
                  </button>
                  <button
                    type="button"
                    onClick={() => setDirectPost(false)}
                    className={`flex-1 px-4 py-3 rounded-xl border-2 text-sm font-medium transition-all ${
                      !directPost
                        ? 'border-indigo-500 bg-indigo-50 text-indigo-700'
                        : 'border-gray-200 text-gray-600 hover:border-gray-300'
                    }`}
                  >
                    Send to Inbox
                  </button>
                </div>
                <p className="text-xs text-gray-500">
                  {directPost
                    ? 'Video will be published automatically to your TikTok account.'
                    : 'Video will be sent to your TikTok inbox for manual review before publishing.'}
                </p>
              </div>
            )}

            {/* Privacy Level (Required - dynamically populated from Creator Info API) */}
            <div className="space-y-2">
              <label className="block text-sm font-medium text-gray-700">
                Privacy Level <span className="text-red-500">*</span>
              </label>
              {creatorInfoLoading ? (
                <div className="flex items-center gap-2 px-4 py-3 border-2 border-gray-200 rounded-xl bg-gray-50 text-gray-500 text-sm">
                  <svg className="animate-spin h-4 w-4" viewBox="0 0 24 24">
                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" />
                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                  </svg>
                  Loading privacy options...
                </div>
              ) : creatorInfoError ? (
                <div className="px-4 py-3 border-2 border-red-200 rounded-xl bg-red-50 text-red-600 text-sm">
                  {creatorInfoError}
                  <button type="button" onClick={fetchCreatorInfo} className="ml-2 underline">Retry</button>
                </div>
              ) : (
                <select
                  value={privacyLevel}
                  onChange={(e) => setPrivacyLevel(e.target.value as TikTokPrivacyLevel | '')}
                  required={isTikTokSelected}
                  className="w-full px-4 py-3 border-2 border-gray-200 rounded-xl focus:border-indigo-500 focus:ring-0 transition-colors bg-white"
                >
                  <option value="">-- Select privacy level --</option>
                  {privacyOptions.map(option => (
                    <option key={option.value} value={option.value}>
                      {option.label}{option.description ? ` - ${option.description}` : ''}
                    </option>
                  ))}
                </select>
              )}
            </div>

            {/* Interaction Settings (disabled states from Creator Info API) */}
            <div className="space-y-3">
              <label className="block text-sm font-medium text-gray-700">Interaction Settings</label>
              <div className="space-y-2">
                <label className={`flex items-center gap-3 ${creatorInfo?.comment_disabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'}`}>
                  <input
                    type="checkbox"
                    checked={allowComment}
                    onChange={(e) => setAllowComment(e.target.checked)}
                    disabled={creatorInfo?.comment_disabled}
                    className="w-5 h-5 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500 disabled:opacity-50"
                  />
                  <span className="text-sm text-gray-700">
                    Allow Comments
                    {creatorInfo?.comment_disabled && <span className="text-xs text-gray-400 ml-1">(disabled by creator settings)</span>}
                  </span>
                </label>
                <label className={`flex items-center gap-3 ${creatorInfo?.duet_disabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'}`}>
                  <input
                    type="checkbox"
                    checked={allowDuet}
                    onChange={(e) => setAllowDuet(e.target.checked)}
                    disabled={creatorInfo?.duet_disabled}
                    className="w-5 h-5 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500 disabled:opacity-50"
                  />
                  <span className="text-sm text-gray-700">
                    Allow Duet
                    {creatorInfo?.duet_disabled && <span className="text-xs text-gray-400 ml-1">(disabled by creator settings)</span>}
                  </span>
                </label>
                <label className={`flex items-center gap-3 ${creatorInfo?.stitch_disabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'}`}>
                  <input
                    type="checkbox"
                    checked={allowStitch}
                    onChange={(e) => setAllowStitch(e.target.checked)}
                    disabled={creatorInfo?.stitch_disabled}
                    className="w-5 h-5 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500 disabled:opacity-50"
                  />
                  <span className="text-sm text-gray-700">
                    Allow Stitch
                    {creatorInfo?.stitch_disabled && <span className="text-xs text-gray-400 ml-1">(disabled by creator settings)</span>}
                  </span>
                </label>
              </div>
            </div>

            {/* Disclose Video Content */}
            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <label className="block text-sm font-medium text-gray-700">Disclose video content</label>
                <button
                  type="button"
                  role="switch"
                  aria-checked={discloseContent}
                  onClick={() => {
                    const next = !discloseContent
                    setDiscloseContent(next)
                    if (!next) {
                      setIsBrandContent(false)
                      setIsBrandOrganic(false)
                    }
                  }}
                  className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                    discloseContent ? 'bg-indigo-600' : 'bg-gray-300'
                  }`}
                >
                  <span className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                    discloseContent ? 'translate-x-6' : 'translate-x-1'
                  }`} />
                </button>
              </div>
              <p className="text-xs text-gray-500">
                Turn on to disclose that this video promotes goods or services in exchange for something of value. Your video could promote yourself, a third party, or both.
              </p>

              {discloseContent && (
                <div className="space-y-3 pt-2">
                  {/* Promotional content warning */}
                  <div className="flex items-start gap-2 px-3 py-2 bg-amber-50 border border-amber-200 rounded-lg">
                    <svg className="w-4 h-4 text-amber-500 mt-0.5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                      <path fillRule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
                    </svg>
                    <p className="text-xs text-amber-700">
                      Your video will be labeled "Promotional content". This cannot be changed once your video is posted.
                    </p>
                  </div>

                  {/* Your Brand */}
                  <label className="flex items-start gap-3 cursor-pointer">
                    <input
                      type="checkbox"
                      checked={isBrandContent}
                      onChange={(e) => setIsBrandContent(e.target.checked)}
                      className="w-5 h-5 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500 mt-0.5"
                    />
                    <div>
                      <span className="text-sm text-gray-700 font-medium">Your brand</span>
                      <p className="text-xs text-gray-500">You are promoting yourself or your own business. This video will be classified as Brand Organic.</p>
                    </div>
                  </label>

                  {/* Branded Content */}
                  <label className="flex items-start gap-3 cursor-pointer">
                    <input
                      type="checkbox"
                      checked={isBrandOrganic}
                      onChange={(e) => setIsBrandOrganic(e.target.checked)}
                      className="w-5 h-5 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500 mt-0.5"
                    />
                    <div>
                      <span className="text-sm text-gray-700 font-medium">Branded content</span>
                      <p className="text-xs text-gray-500">You are promoting another brand or a third party. This video will be classified as Branded Content.</p>
                    </div>
                  </label>

                  {isBrandOrganic && privacyLevel === 'SELF_ONLY' && (
                    <p className="text-xs text-red-500">
                      Branded content cannot be set to private.
                    </p>
                  )}
                </div>
              )}
            </div>

            {/* Auto-Add Music (for photo posts) */}
            {detectedMediaType === 'image' && (
              <div className="space-y-3">
                <label className="block text-sm font-medium text-gray-700">Music Settings</label>
                <label className="flex items-center gap-3 cursor-pointer">
                  <input
                    type="checkbox"
                    checked={autoAddMusic}
                    onChange={(e) => setAutoAddMusic(e.target.checked)}
                    className="w-5 h-5 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500"
                  />
                  <div>
                    <span className="text-sm text-gray-700 font-medium">Auto-Add Music</span>
                    <p className="text-xs text-gray-500">TikTok will automatically add trending music to your photo post</p>
                  </div>
                </label>
              </div>
            )}

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
