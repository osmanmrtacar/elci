import { useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import LoginButton from '../components/auth/LoginButton'
import Footer from '../components/Footer'

const HomePage = () => {
  const { isAuthenticated } = useAuth()
  const navigate = useNavigate()

  useEffect(() => {
    if (isAuthenticated) {
      navigate('/dashboard')
    }
  }, [isAuthenticated, navigate])

  return (
    <div className="home-page">
      <div className="hero-container">
        <div className="hero-content">
          <h1>Social Media Publisher</h1>
          <p className="hero-subtitle">
            Publish your content to TikTok, Instagram, and X (Twitter)
          </p>
          <p className="hero-description">
            Connect your social media accounts and start posting content across multiple platforms.
            Simple, fast, and powerful.
          </p>
          <div className="hero-buttons">
            <LoginButton />
          </div>
          <div className="features">
            <div className="feature">
              <div className="feature-icon">🔐</div>
              <h3>Secure OAuth</h3>
              <p>Login safely with official platform authentication</p>
            </div>
            <div className="feature">
              <div className="feature-icon">🎬</div>
              <h3>Multi-Platform</h3>
              <p>Post to TikTok, Instagram, and X simultaneously</p>
            </div>
            <div className="feature">
              <div className="feature-icon">📊</div>
              <h3>Track History</h3>
              <p>View all your posts and their status in one place</p>
            </div>
          </div>
        </div>
      </div>
      <Footer />
    </div>
  )
}

export default HomePage
