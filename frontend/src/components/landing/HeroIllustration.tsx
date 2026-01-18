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
          {/* TikTok icon */}
          <g transform="translate(465, 105) scale(0.12)">
            <path d="M224,72a48.05,48.05,0,0,1-48-48,8,8,0,0,0-8-8H128a8,8,0,0,0-8,8V156a20,20,0,1,1-28.57-18.08A8,8,0,0,0,96,130.69V88a8,8,0,0,0-9.4-7.88C50.91,86.48,24,119.1,24,156a76,76,0,0,0,152,0V116.29A103.25,103.25,0,0,0,224,128a8,8,0,0,0,8-8V80A8,8,0,0,0,224,72Zm-8,39.64a87.19,87.19,0,0,1-43.33-16.15A8,8,0,0,0,160,102v54a60,60,0,0,1-120,0c0-25.9,16.64-49.13,40-57.6v27.67A36,36,0,1,0,136,156V32h24.5A64.14,64.14,0,0,0,216,87.5Z" fill="url(#gradient1)" />
          </g>
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
          {/* YouTube icon */}
          <g transform="translate(465, 315) scale(0.12)">
            <path d="M164.44,121.34l-48-32A8,8,0,0,0,104,96v64a8,8,0,0,0,12.44,6.66l48-32a8,8,0,0,0,0-13.32ZM120,145.05V111l25.58,17ZM234.33,69.52a24,24,0,0,0-14.49-16.4C185.56,39.88,131,40,128,40s-57.56-.12-91.84,13.12a24,24,0,0,0-14.49,16.4C19.08,79.5,16,97.74,16,128s3.08,48.5,5.67,58.48a24,24,0,0,0,14.49,16.41C69,215.56,120.4,216,127.34,216h1.32c6.94,0,58.37-.44,91.18-13.11a24,24,0,0,0,14.49-16.41c2.59-10,5.67-28.22,5.67-58.48S236.92,79.5,234.33,69.52Zm-15.49,113a8,8,0,0,1-4.77,5.49c-31.65,12.22-85.48,12-86,12H128c-.54,0-54.33.2-86-12a8,8,0,0,1-4.77-5.49C34.8,173.39,32,156.57,32,128s2.8-45.39,5.16-54.47A8,8,0,0,1,41.93,68c30.52-11.79,81.66-12,85.85-12h.27c.54,0,54.38-.18,86,12a8,8,0,0,1,4.77,5.49C221.2,82.61,224,99.43,224,128S221.2,173.39,218.84,182.47Z" fill="url(#gradient1)" />
          </g>
          {/* Coming Soon label */}
          <text x="480" y="380" fontSize="10" fill="url(#gradient1)" textAnchor="middle" fontFamily="Arial, sans-serif" fontWeight="500">Coming Soon</text>
        </g>

        {/* Twitter/X - Top Left */}
        <g className="animate-pulse" style={{ animationDelay: '0.9s', animationDuration: '2s' }}>
          <path d="M 225 210 Q 150 150 120 120" stroke="url(#gradient2)" strokeWidth="2" strokeDasharray="5,5" opacity="0.6">
            <animate attributeName="stroke-dashoffset" from="10" to="0" dur="1s" repeatCount="indefinite" />
          </path>
          <circle cx="120" cy="120" r="35" fill="white" stroke="url(#gradient2)" strokeWidth="3" />
          {/* X icon */}
          <g transform="translate(105, 105) scale(0.12)">
            <path d="M214.75,211.71l-62.6-98.38,61.77-67.95a8,8,0,0,0-11.84-10.76L143.24,99.34,102.75,35.71A8,8,0,0,0,96,32H48a8,8,0,0,0-6.75,12.3l62.6,98.37-61.77,68a8,8,0,1,0,11.84,10.76l58.84-64.72,40.49,63.63A8,8,0,0,0,160,224h48a8,8,0,0,0,6.75-12.29ZM164.39,208,62.57,48h29L193.43,208Z" fill="url(#gradient2)" />
          </g>
        </g>

        {/* LinkedIn - Left */}
        <g className="animate-pulse" style={{ animationDelay: '1.2s', animationDuration: '2s' }}>
          <path d="M 225 210 Q 150 210 100 240" stroke="url(#gradient1)" strokeWidth="2" strokeDasharray="5,5" opacity="0.6">
            <animate attributeName="stroke-dashoffset" from="10" to="0" dur="1s" repeatCount="indefinite" />
          </path>
          <circle cx="100" cy="240" r="35" fill="white" stroke="url(#gradient1)" strokeWidth="3" />
          {/* LinkedIn icon simplified */}
          <rect x="88" y="228" width="24" height="24" rx="3" fill="url(#gradient1)" />
          <text x="100" y="244" fontSize="14" fontWeight="bold" fill="white" textAnchor="middle" fontFamily="Arial, sans-serif">in</text>
          {/* Coming Soon label */}
          <text x="100" y="290" fontSize="10" fill="url(#gradient1)" textAnchor="middle" fontFamily="Arial, sans-serif" fontWeight="500">Coming Soon</text>
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
