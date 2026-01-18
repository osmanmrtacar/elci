import { useEffect, useState } from 'react'
import { useAuth } from '../context/AuthContext'
import { useNavigate } from 'react-router-dom'

const SuccessPage = () => {
  const { user, isAuthenticated, refreshPlatforms } = useAuth()
  const navigate = useNavigate()
  const [countdown, setCountdown] = useState(3)

  useEffect(() => {
    // Check if authenticated
    if (!isAuthenticated) {
      navigate('/')
      return
    }

    // Refresh platforms to ensure connected state is up to date
    refreshPlatforms()

    // Auto-redirect countdown
    const timer = setInterval(() => {
      setCountdown((prev) => {
        if (prev <= 1) {
          clearInterval(timer)
          navigate('/dashboard')
          return 0
        }
        return prev - 1
      })
    }, 1000)

    return () => clearInterval(timer)
  }, [isAuthenticated, navigate, refreshPlatforms])

  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center p-4">
      <div className="bg-white rounded-2xl shadow-lg border border-gray-200 p-8 max-w-md w-full text-center">
        {/* Success Icon */}
        <div className="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-6">
          <svg className="w-8 h-8 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
          </svg>
        </div>

        <h1 className="text-2xl font-bold text-gray-900 mb-2">
          Connection Successful!
        </h1>

        <p className="text-gray-600 mb-6">
          Your account has been connected to elci.io
        </p>

        {/* User Info */}
        {user && (
          <div className="flex items-center justify-center gap-3 mb-6 p-4 bg-gray-50 rounded-xl">
            {user.avatar_url && (
              <img
                src={user.avatar_url}
                alt={user.display_name}
                className="w-12 h-12 rounded-full border-2 border-indigo-100"
              />
            )}
            <div className="text-left">
              <div className="font-medium text-gray-900">{user.display_name}</div>
              <div className="text-sm text-gray-500">@{user.username}</div>
            </div>
          </div>
        )}

        {/* Redirect Notice */}
        <div className="mb-6">
          <div className="flex items-center justify-center gap-2 text-gray-600">
            <div className="w-5 h-5 border-2 border-indigo-600 border-t-transparent rounded-full animate-spin"></div>
            <span>Redirecting to dashboard in {countdown}...</span>
          </div>
        </div>

        {/* Manual Button */}
        <button
          onClick={() => navigate('/dashboard')}
          className="w-full px-6 py-3 bg-indigo-600 text-white font-medium rounded-lg hover:bg-indigo-700 transition-colors"
        >
          Go to Dashboard Now
        </button>
      </div>
    </div>
  )
}

export default SuccessPage
