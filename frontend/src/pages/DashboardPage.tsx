import { useState } from 'react'
import { useAuth } from '../context/AuthContext'
import PostForm from '../components/post/PostForm'
import PostHistory from '../components/post/PostHistory'
import Footer from '../components/Footer'

const DashboardPage = () => {
  const { user, logout } = useAuth()
  const [refreshTrigger, setRefreshTrigger] = useState(0)

  const handlePostCreated = () => {
    // Trigger refresh of post history
    setRefreshTrigger((prev) => prev + 1)
  }

  const handleLogout = async () => {
    await logout()
    window.location.href = '/'
  }

  return (
    <div className="dashboard-page">
      <header className="dashboard-header">
        <div className="header-content">
          <h1>Dashboard</h1>
          <div className="user-section">
            {user?.avatar_url && (
              <img
                src={user.avatar_url}
                alt={user.display_name}
                className="user-avatar"
              />
            )}
            <div className="user-info">
              <div className="user-name">{user?.display_name}</div>
              <div className="user-username">@{user?.username}</div>
            </div>
            <button onClick={handleLogout} className="logout-button">
              Logout
            </button>
          </div>
        </div>
      </header>

      <main className="dashboard-main">
        <div className="dashboard-grid">
          <div className="dashboard-section">
            <PostForm onPostCreated={handlePostCreated} />
          </div>
          <div className="dashboard-section">
            <PostHistory refreshTrigger={refreshTrigger} />
          </div>
        </div>
      </main>

      <Footer />
    </div>
  )
}

export default DashboardPage
