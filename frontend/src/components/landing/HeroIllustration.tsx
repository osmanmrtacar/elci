export function HeroIllustration() {
  return (
    <div className="relative w-full h-[500px]">
      <svg viewBox="0 0 600 500" fill="none" xmlns="http://www.w3.org/2000/svg" className="w-full h-full">
        {/* Background decorative elements */}
        <circle cx="500" cy="100" r="120" fill="url(#gradient1)" opacity="0.1" />
        <circle cx="100" cy="400" r="80" fill="url(#gradient2)" opacity="0.1" />

        {/* Central video/content card */}
        <g transform="translate(225, 150)">
          <rect width="150" height="120" rx="12" fill="white" stroke="url(#gradient1)" strokeWidth="3" />
          <rect x="15" y="15" width="120" height="75" rx="6" fill="url(#gradient1)" opacity="0.2" />
          {/* Play button */}
          <circle cx="75" cy="52.5" r="18" fill="url(#gradient1)" />
          <path d="M70 45 L85 52.5 L70 60 Z" fill="white" />
          {/* Progress bar */}
          <rect x="15" y="100" width="120" height="4" rx="2" fill="#E5E7EB" />
          <rect x="15" y="100" width="60" height="4" rx="2" fill="url(#gradient1)" />
        </g>

        {/* Animated connection lines and platform icons */}
        {/* TikTok - Top Right */}
        <g className="animate-pulse" style={{ animationDelay: '0s', animationDuration: '2s' }}>
          <path d="M 375 210 Q 450 150 480 120" stroke="url(#gradient1)" strokeWidth="2" strokeDasharray="5,5" opacity="0.6">
            <animate attributeName="stroke-dashoffset" from="10" to="0" dur="1s" repeatCount="indefinite" />
          </path>
          <circle cx="480" cy="120" r="35" fill="white" stroke="url(#gradient1)" strokeWidth="3" />
          {/* TikTok icon simplified */}
          <path d="M475 110 L475 130 M485 110 L485 125 M480 125 Q485 128 490 125" stroke="url(#gradient1)" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round" />
        </g>

        {/* Instagram - Right */}
        <g className="animate-pulse" style={{ animationDelay: '0.3s', animationDuration: '2s' }}>
          <path d="M 375 210 Q 450 210 500 240" stroke="url(#gradient2)" strokeWidth="2" strokeDasharray="5,5" opacity="0.6">
            <animate attributeName="stroke-dashoffset" from="10" to="0" dur="1s" repeatCount="indefinite" />
          </path>
          <circle cx="500" cy="240" r="35" fill="white" stroke="url(#gradient2)" strokeWidth="3" />
          {/* Instagram icon simplified */}
          <rect x="485" y="225" width="30" height="30" rx="8" stroke="url(#gradient2)" strokeWidth="2.5" fill="none" />
          <circle cx="500" cy="240" r="8" stroke="url(#gradient2)" strokeWidth="2.5" fill="none" />
          <circle cx="510" cy="230" r="2" fill="url(#gradient2)" />
        </g>

        {/* YouTube - Bottom Right */}
        <g className="animate-pulse" style={{ animationDelay: '0.6s', animationDuration: '2s' }}>
          <path d="M 375 210 Q 450 280 480 330" stroke="url(#gradient1)" strokeWidth="2" strokeDasharray="5,5" opacity="0.6">
            <animate attributeName="stroke-dashoffset" from="10" to="0" dur="1s" repeatCount="indefinite" />
          </path>
          <circle cx="480" cy="330" r="35" fill="white" stroke="url(#gradient1)" strokeWidth="3" />
          {/* YouTube icon simplified */}
          <path d="M 500 330 L 500 330 Q 502 325 495 325 L 465 325 Q 458 325 460 330 L 460 330 Q 458 335 465 335 L 495 335 Q 502 335 500 330 Z" fill="url(#gradient1)" />
          <path d="M 475 327 L 485 330 L 475 333 Z" fill="white" />
        </g>

        {/* Twitter/X - Top Left */}
        <g className="animate-pulse" style={{ animationDelay: '0.9s', animationDuration: '2s' }}>
          <path d="M 225 210 Q 150 150 120 120" stroke="url(#gradient2)" strokeWidth="2" strokeDasharray="5,5" opacity="0.6">
            <animate attributeName="stroke-dashoffset" from="10" to="0" dur="1s" repeatCount="indefinite" />
          </path>
          <circle cx="120" cy="120" r="35" fill="white" stroke="url(#gradient2)" strokeWidth="3" />
          {/* X icon simplified */}
          <path d="M 110 110 L 130 130 M 130 110 L 110 130" stroke="url(#gradient2)" strokeWidth="2.5" strokeLinecap="round" />
        </g>

        {/* LinkedIn - Left */}
        <g className="animate-pulse" style={{ animationDelay: '1.2s', animationDuration: '2s' }}>
          <path d="M 225 210 Q 150 210 100 240" stroke="url(#gradient1)" strokeWidth="2" strokeDasharray="5,5" opacity="0.6">
            <animate attributeName="stroke-dashoffset" from="10" to="0" dur="1s" repeatCount="indefinite" />
          </path>
          <circle cx="100" cy="240" r="35" fill="white" stroke="url(#gradient1)" strokeWidth="3" />
          {/* LinkedIn icon simplified */}
          <rect x="90" y="235" width="5" height="15" rx="1" fill="url(#gradient1)" />
          <rect x="100" y="230" width="5" height="20" rx="1" fill="url(#gradient1)" />
          <rect x="110" y="235" width="5" height="15" rx="1" fill="url(#gradient1)" />
        </g>

        {/* Gradients */}
        <defs>
          <linearGradient id="gradient1" x1="0%" y1="0%" x2="100%" y2="100%">
            <stop offset="0%" stopColor="#4F46E5" />
            <stop offset="100%" stopColor="#7C3AED" />
          </linearGradient>
          <linearGradient id="gradient2" x1="0%" y1="0%" x2="100%" y2="100%">
            <stop offset="0%" stopColor="#7C3AED" />
            <stop offset="100%" stopColor="#A855F7" />
          </linearGradient>
        </defs>
      </svg>
    </div>
  )
}
