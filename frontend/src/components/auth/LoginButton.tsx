import { useAuth } from '../../context/AuthContext'

const LoginButton = () => {
  const { loginTikTok, loginX, loginInstagram, isPlatformConnected } = useAuth()

  const isTikTokConnected = isPlatformConnected('tiktok')
  const isXConnected = isPlatformConnected('x')
  const isInstagramConnected = isPlatformConnected('instagram')

  return (
    <div className="login-buttons-container">
      <button
        onClick={loginTikTok}
        disabled={isTikTokConnected}
        className={`login-button tiktok-button ${isTikTokConnected ? 'connected' : ''}`}
        title={isTikTokConnected ? 'TikTok Connected' : 'Connect TikTok'}
      >
        <svg
          width="24"
          height="24"
          viewBox="0 0 24 24"
          fill="currentColor"
          style={{ marginRight: '8px' }}
        >
          <path d="M19.59 6.69a4.83 4.83 0 0 1-3.77-4.25V2h-3.45v13.67a2.89 2.89 0 0 1-5.2 1.74 2.89 2.89 0 0 1 2.31-4.64 2.93 2.93 0 0 1 .88.13V9.4a6.84 6.84 0 0 0-1-.05A6.33 6.33 0 0 0 5 20.1a6.34 6.34 0 0 0 10.86-4.43v-7a8.16 8.16 0 0 0 4.77 1.52v-3.4a4.85 4.85 0 0 1-1-.1z" />
        </svg>
        {isTikTokConnected ? (
          <>
            <span>TikTok Connected</span>
            <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor" style={{ marginLeft: '8px' }}>
              <path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"/>
            </svg>
          </>
        ) : (
          'Connect TikTok'
        )}
      </button>

      <button
        onClick={loginX}
        disabled={isXConnected}
        className={`login-button x-button ${isXConnected ? 'connected' : ''}`}
        title={isXConnected ? 'X Connected' : 'Connect X'}
      >
        <svg
          width="24"
          height="24"
          viewBox="0 0 24 24"
          fill="currentColor"
          style={{ marginRight: '8px' }}
        >
          <path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z"/>
        </svg>
        {isXConnected ? (
          <>
            <span>X Connected</span>
            <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor" style={{ marginLeft: '8px' }}>
              <path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"/>
            </svg>
          </>
        ) : (
          'Connect X'
        )}
      </button>

      <button
        onClick={loginInstagram}
        disabled={isInstagramConnected}
        className={`login-button instagram-button ${isInstagramConnected ? 'connected' : ''}`}
        title={isInstagramConnected ? 'Instagram Connected' : 'Connect Instagram'}
      >
        <svg
          width="24"
          height="24"
          viewBox="0 0 24 24"
          fill="currentColor"
          style={{ marginRight: '8px' }}
        >
          <path d="M7.8 2h8.4C19.4 2 22 4.6 22 7.8v8.4a5.8 5.8 0 0 1-5.8 5.8H7.8C4.6 22 2 19.4 2 16.2V7.8A5.8 5.8 0 0 1 7.8 2m-.2 2A3.6 3.6 0 0 0 4 7.6v8.8C4 18.39 5.61 20 7.6 20h8.8a3.6 3.6 0 0 0 3.6-3.6V7.6C20 5.61 18.39 4 16.4 4H7.6m9.65 1.5a1.25 1.25 0 0 1 1.25 1.25A1.25 1.25 0 0 1 17.25 8 1.25 1.25 0 0 1 16 6.75a1.25 1.25 0 0 1 1.25-1.25M12 7a5 5 0 0 1 5 5 5 5 0 0 1-5 5 5 5 0 0 1-5-5 5 5 0 0 1 5-5m0 2a3 3 0 0 0-3 3 3 3 0 0 0 3 3 3 3 0 0 0 3-3 3 3 0 0 0-3-3z"/>
        </svg>
        {isInstagramConnected ? (
          <>
            <span>Instagram Connected</span>
            <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor" style={{ marginLeft: '8px' }}>
              <path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"/>
            </svg>
          </>
        ) : (
          'Connect Instagram'
        )}
      </button>

      <style>{`
        .login-buttons-container {
          display: flex;
          flex-direction: column;
          gap: 12px;
          width: 100%;
          max-width: 400px;
        }

        .login-button {
          display: flex;
          align-items: center;
          justify-content: center;
          padding: 12px 24px;
          border: none;
          border-radius: 8px;
          font-size: 16px;
          font-weight: 600;
          cursor: pointer;
          transition: all 0.3s ease;
          color: white;
          width: 100%;
        }

        .tiktok-button {
          background: linear-gradient(135deg, #fe2c55 0%, #000000 100%);
        }

        .tiktok-button:hover:not(:disabled) {
          transform: translateY(-2px);
          box-shadow: 0 4px 12px rgba(254, 44, 85, 0.4);
        }

        .x-button {
          background: linear-gradient(135deg, #1DA1F2 0%, #000000 100%);
        }

        .x-button:hover:not(:disabled) {
          transform: translateY(-2px);
          box-shadow: 0 4px 12px rgba(29, 161, 242, 0.4);
        }

        .instagram-button {
          background: linear-gradient(135deg, #E1306C 0%, #FD1D1D 50%, #F77737 100%);
        }

        .instagram-button:hover:not(:disabled) {
          transform: translateY(-2px);
          box-shadow: 0 4px 12px rgba(225, 48, 108, 0.4);
        }

        .login-button:disabled {
          opacity: 0.7;
          cursor: not-allowed;
        }

        .login-button.connected {
          background: linear-gradient(135deg, #10b981 0%, #059669 100%);
        }

        .login-button.connected:hover {
          transform: none;
          box-shadow: none;
        }

        @media (min-width: 768px) {
          .login-buttons-container {
            flex-direction: row;
          }
        }
      `}</style>
    </div>
  )
}

export default LoginButton
