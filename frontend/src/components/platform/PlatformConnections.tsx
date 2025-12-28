import { useAuth } from '../../context/AuthContext'

const PlatformConnections = () => {
  const {
    connectedPlatforms,
    loginTikTok,
    loginX,
    loginInstagram,
    disconnectPlatform,
    isPlatformConnected,
  } = useAuth()

  const platforms = [
    {
      id: 'tiktok' as const,
      name: 'TikTok',
      color: 'linear-gradient(135deg, #fe2c55 0%, #000000 100%)',
      hoverShadow: 'rgba(254, 44, 85, 0.4)',
      icon: (
        <svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
          <path d="M19.59 6.69a4.83 4.83 0 0 1-3.77-4.25V2h-3.45v13.67a2.89 2.89 0 0 1-5.2 1.74 2.89 2.89 0 0 1 2.31-4.64 2.93 2.93 0 0 1 .88.13V9.4a6.84 6.84 0 0 0-1-.05A6.33 6.33 0 0 0 5 20.1a6.34 6.34 0 0 0 10.86-4.43v-7a8.16 8.16 0 0 0 4.77 1.52v-3.4a4.85 4.85 0 0 1-1-.1z" />
        </svg>
      ),
      loginFn: loginTikTok,
    },
    {
      id: 'x' as const,
      name: 'X',
      color: 'linear-gradient(135deg, #1DA1F2 0%, #000000 100%)',
      hoverShadow: 'rgba(29, 161, 242, 0.4)',
      icon: (
        <svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
          <path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z" />
        </svg>
      ),
      loginFn: loginX,
    },
    {
      id: 'instagram' as const,
      name: 'Instagram',
      color: 'linear-gradient(135deg, #E1306C 0%, #FD1D1D 50%, #F77737 100%)',
      hoverShadow: 'rgba(225, 48, 108, 0.4)',
      icon: (
        <svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
          <path d="M7.8 2h8.4C19.4 2 22 4.6 22 7.8v8.4a5.8 5.8 0 0 1-5.8 5.8H7.8C4.6 22 2 19.4 2 16.2V7.8A5.8 5.8 0 0 1 7.8 2m-.2 2A3.6 3.6 0 0 0 4 7.6v8.8C4 18.39 5.61 20 7.6 20h8.8a3.6 3.6 0 0 0 3.6-3.6V7.6C20 5.61 18.39 4 16.4 4H7.6m9.65 1.5a1.25 1.25 0 0 1 1.25 1.25A1.25 1.25 0 0 1 17.25 8 1.25 1.25 0 0 1 16 6.75a1.25 1.25 0 0 1 1.25-1.25M12 7a5 5 0 0 1 5 5 5 5 0 0 1-5 5 5 5 0 0 1-5-5 5 5 0 0 1 5-5m0 2a3 3 0 0 0-3 3 3 3 0 0 0 3 3 3 3 0 0 0 3-3 3 3 0 0 0-3-3z" />
        </svg>
      ),
      loginFn: loginInstagram,
    },
  ]

  const handleDisconnect = async (platformId: 'tiktok' | 'x' | 'instagram') => {
    if (
      window.confirm(
        `Are you sure you want to disconnect your ${platformId.toUpperCase()} account?`
      )
    ) {
      try {
        await disconnectPlatform(platformId)
      } catch (error) {
        console.error(`Failed to disconnect ${platformId}:`, error)
        alert(`Failed to disconnect ${platformId}. Please try again.`)
      }
    }
  }

  return (
    <div className="platform-connections">
      <h2 className="section-title">Connected Platforms</h2>
      <div className="platforms-grid">
        {platforms.map((platform) => {
          const isConnected = isPlatformConnected(platform.id)
          const connection = connectedPlatforms.find(
            (p) => p.platform === platform.id
          )

          return (
            <div
              key={platform.id}
              className={`platform-card ${isConnected ? 'connected' : ''}`}
            >
              <div className="platform-header">
                <div className="platform-icon">{platform.icon}</div>
                <div className="platform-info">
                  <h3 className="platform-name">{platform.name}</h3>
                  {isConnected && connection && (
                    <div className="platform-details">
                      <p className="platform-username">
                        @{connection.username}
                      </p>
                      <p className="platform-display-name">
                        {connection.display_name}
                      </p>
                    </div>
                  )}
                </div>
              </div>

              <div className="platform-actions">
                {isConnected ? (
                  <button
                    onClick={() => handleDisconnect(platform.id)}
                    className="disconnect-button"
                  >
                    Disconnect
                  </button>
                ) : (
                  <button
                    onClick={platform.loginFn}
                    className="connect-button"
                    style={{ background: platform.color }}
                  >
                    Connect {platform.name}
                  </button>
                )}
              </div>
            </div>
          )
        })}
      </div>

      <style>{`
        .platform-connections {
          padding: 24px;
          background: white;
          border-radius: 12px;
          box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
        }

        .section-title {
          margin: 0 0 20px 0;
          font-size: 24px;
          font-weight: 600;
          color: #1f2937;
        }

        .platforms-grid {
          display: grid;
          grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
          gap: 16px;
        }

        .platform-card {
          padding: 20px;
          border: 2px solid #e5e7eb;
          border-radius: 12px;
          transition: all 0.3s ease;
          background: #f9fafb;
        }

        .platform-card.connected {
          border-color: #10b981;
          background: #f0fdf4;
        }

        .platform-card:hover {
          box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
        }

        .platform-header {
          display: flex;
          gap: 12px;
          margin-bottom: 16px;
        }

        .platform-icon {
          width: 48px;
          height: 48px;
          display: flex;
          align-items: center;
          justify-content: center;
          background: white;
          border-radius: 12px;
          box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
        }

        .platform-icon svg {
          width: 28px;
          height: 28px;
        }

        .platform-info {
          flex: 1;
        }

        .platform-name {
          margin: 0 0 4px 0;
          font-size: 18px;
          font-weight: 600;
          color: #1f2937;
        }

        .platform-details {
          margin-top: 4px;
        }

        .platform-username {
          margin: 0;
          font-size: 14px;
          font-weight: 500;
          color: #6b7280;
        }

        .platform-display-name {
          margin: 2px 0 0 0;
          font-size: 12px;
          color: #9ca3af;
        }

        .platform-actions {
          display: flex;
          gap: 8px;
        }

        .connect-button,
        .disconnect-button {
          flex: 1;
          padding: 10px 16px;
          border: none;
          border-radius: 8px;
          font-size: 14px;
          font-weight: 600;
          cursor: pointer;
          transition: all 0.3s ease;
          color: white;
        }

        .connect-button {
          background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
        }

        .connect-button:hover {
          transform: translateY(-2px);
          box-shadow: 0 4px 12px rgba(59, 130, 246, 0.4);
        }

        .disconnect-button {
          background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
        }

        .disconnect-button:hover {
          transform: translateY(-2px);
          box-shadow: 0 4px 12px rgba(239, 68, 68, 0.4);
        }

        @media (max-width: 768px) {
          .platform-connections {
            padding: 16px;
          }

          .platforms-grid {
            grid-template-columns: 1fr;
          }

          .section-title {
            font-size: 20px;
          }
        }
      `}</style>
    </div>
  )
}

export default PlatformConnections
