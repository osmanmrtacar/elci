import { useState } from 'react'
import { Link } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import PostForm from '../components/post/PostForm'
import PostHistory from '../components/post/PostHistory'
import PlatformConnections from '../components/platform/PlatformConnections'
import { Footer } from '../components/landing'

const DashboardPage = () => {
  const { user, logout } = useAuth()
  const [refreshTrigger, setRefreshTrigger] = useState(0)

  const handlePostCreated = () => {
    setRefreshTrigger((prev) => prev + 1)
  }

  const handleLogout = async () => {
    await logout()
    window.location.href = '/'
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="sticky top-0 z-50 bg-white border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-6 py-4 flex items-center justify-between">
          <Link to="/" className="flex items-center gap-2">
            <img src="/logo.svg" alt="elci.io" className="h-10" />
            <span className="text-xl font-semibold text-gray-900">elci.io</span>
          </Link>

          <div className="flex items-center gap-4">
            {user?.avatar_url && (
              <img
                src={user.avatar_url}
                alt={user.display_name}
                className="w-10 h-10 rounded-full border-2 border-indigo-100"
              />
            )}
            <div className="hidden sm:block text-right">
              <div className="text-sm font-medium text-gray-900">{user?.display_name}</div>
              <div className="text-xs text-gray-500">@{user?.username}</div>
            </div>
            <button
              onClick={handleLogout}
              className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors"
            >
              Logout
            </button>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-6 py-8">
        {/* Page Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
          <p className="text-gray-600 mt-1">Manage your content and connected platforms</p>
        </div>

        {/* Platform Connections */}
        <div className="mb-8">
          <PlatformConnections />
        </div>

        {/* Two Column Layout */}
        <div className="grid lg:grid-cols-2 gap-8">
          {/* Post Form */}
          <div>
            <PostForm onPostCreated={handlePostCreated} />
          </div>

          {/* Post History */}
          <div>
            <PostHistory refreshTrigger={refreshTrigger} />
          </div>
        </div>
      </main>

      <Footer />
    </div>
  )
}

export default DashboardPage
