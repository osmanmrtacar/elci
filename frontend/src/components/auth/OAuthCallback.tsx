import { useEffect, useState } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { useAuth } from '../../context/AuthContext'
import { authService } from '../../services/authService'

const OAuthCallback = () => {
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const { setUser } = useAuth()
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const handleCallback = async () => {
      // Check if we have a token from backend redirect
      const token = searchParams.get('token')

      if (token) {
        // Backend already handled the OAuth flow and gave us a token
        try {
          // Save token first
          localStorage.setItem('auth_token', token)

          // Fetch user info
          const currentUser = await authService.getCurrentUser()

          // Save user to localStorage and update context
          authService.saveAuth(token, currentUser)
          setUser(currentUser)

          navigate('/success')
        } catch (err: any) {
          console.error('Failed to get user info:', err)
          setError('Failed to complete authentication')
          setTimeout(() => navigate('/'), 3000)
        }
        return
      }

      // Legacy flow: handle code/state (shouldn't happen with new redirect flow)
      const code = searchParams.get('code')
      const state = searchParams.get('state')
      const errorParam = searchParams.get('error')

      if (errorParam) {
        setError(`Authentication failed: ${searchParams.get('error_description') || errorParam}`)
        setTimeout(() => navigate('/'), 3000)
        return
      }

      if (!code || !state) {
        setError('Missing authorization parameters')
        setTimeout(() => navigate('/'), 3000)
        return
      }

      try {
        const data = await authService.handleCallback(code, state)
        authService.saveAuth(data.token, data.user)
        setUser(data.user)
        navigate('/success')
      } catch (err: any) {
        console.error('Callback error:', err)
        setError(err.response?.data?.error || 'Authentication failed')
        setTimeout(() => navigate('/'), 3000)
      }
    }

    handleCallback()
  }, [searchParams, navigate, setUser])

  return (
    <div className="callback-container">
      <div className="callback-card">
        {error ? (
          <>
            <div className="error-icon">❌</div>
            <h2>Authentication Failed</h2>
            <p>{error}</p>
            <p style={{ fontSize: '14px', color: '#666', marginTop: '10px' }}>
              Redirecting to home...
            </p>
          </>
        ) : (
          <>
            <div className="loading-spinner"></div>
            <h2>Completing authentication...</h2>
            <p>Please wait while we log you in.</p>
          </>
        )}
      </div>
    </div>
  )
}

export default OAuthCallback
