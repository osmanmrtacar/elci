import { useEffect, useState } from 'react'
import { useAuth } from '../context/AuthContext'
import { useNavigate } from 'react-router-dom'

const SuccessPage = () => {
  const { user, isAuthenticated, logout } = useAuth()
  const navigate = useNavigate()
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    // Check if authenticated
    if (!isAuthenticated) {
      navigate('/')
      return
    }
    setLoading(false)
  }, [isAuthenticated, navigate])

  const handleLogout = async () => {
    await logout()
    navigate('/')
  }

  if (loading) {
    return (
      <div style={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        minHeight: '100vh',
        fontFamily: 'system-ui, -apple-system, sans-serif'
      }}>
        <div>Loading...</div>
      </div>
    )
  }

  return (
    <div style={{
      minHeight: '100vh',
      background: 'linear-gradient(135deg, #E0E7D7 0%, #EDECEC 100%)',
      display: 'flex',
      justifyContent: 'center',
      alignItems: 'center',
      padding: '20px',
      fontFamily: 'system-ui, -apple-system, sans-serif'
    }}>
      <div style={{
        background: '#FEFEFE',
        borderRadius: '16px',
        padding: '40px',
        maxWidth: '500px',
        width: '100%',
        boxShadow: '0 8px 32px rgba(0,0,0,0.1)',
        border: '2px solid #CCCCCC',
        textAlign: 'center'
      }}>
        <div style={{
          fontSize: '48px',
          marginBottom: '20px'
        }}>
          ✅
        </div>

        <h1 style={{
          margin: '0 0 10px 0',
          fontSize: '32px',
          color: '#1a1a1a'
        }}>
          Login Successful!
        </h1>

        <p style={{
          color: '#666',
          marginBottom: '30px'
        }}>
          Welcome to your TikTok Content Publisher
        </p>

        {user?.avatar_url && (
          <img
            src={user.avatar_url}
            alt={user.display_name}
            style={{
              width: '100px',
              height: '100px',
              borderRadius: '50%',
              marginBottom: '20px',
              border: '4px solid #B7C396'
            }}
          />
        )}

        <div style={{
          background: '#EDECEC',
          padding: '20px',
          borderRadius: '12px',
          marginBottom: '30px',
          textAlign: 'left',
          border: '1px solid #CCCCCC'
        }}>
          <h3 style={{
            margin: '0 0 15px 0',
            fontSize: '16px',
            color: '#666',
            textAlign: 'center'
          }}>
            Account Information
          </h3>

          <div style={{ marginBottom: '10px' }}>
            <strong style={{ color: '#B7C396' }}>Display Name:</strong>{' '}
            <span style={{ color: '#1a1a1a' }}>{user?.display_name}</span>
          </div>

          <div style={{ marginBottom: '10px' }}>
            <strong style={{ color: '#B7C396' }}>Username:</strong>{' '}
            <span style={{ color: '#1a1a1a' }}>@{user?.username}</span>
          </div>

          <div>
            <strong style={{ color: '#B7C396' }}>User ID:</strong>{' '}
            <span style={{ color: '#1a1a1a' }}>{user?.id}</span>
          </div>
        </div>

        <div style={{
          background: '#E0E7D7',
          padding: '15px',
          borderRadius: '8px',
          marginBottom: '30px',
          fontSize: '14px',
          color: '#5a6650',
          textAlign: 'left',
          border: '1px solid #B7C396'
        }}>
          <strong>🎉 Your TikTok account is connected!</strong>
          <br />
          You can now post videos to TikTok via the API.
        </div>

        <div style={{
          display: 'flex',
          gap: '10px',
          flexDirection: 'column'
        }}>
          <button
            onClick={() => navigate('/dashboard')}
            style={{
              background: '#B7C396',
              color: 'white',
              border: 'none',
              padding: '14px 28px',
              borderRadius: '8px',
              fontSize: '16px',
              fontWeight: 'bold',
              cursor: 'pointer',
              transition: 'all 0.2s'
            }}
            onMouseOver={(e) => {
              e.currentTarget.style.background = '#BA9A91'
              e.currentTarget.style.transform = 'translateY(-2px)'
            }}
            onMouseOut={(e) => {
              e.currentTarget.style.background = '#B7C396'
              e.currentTarget.style.transform = 'translateY(0)'
            }}
          >
            Go to Dashboard
          </button>

          <button
            onClick={handleLogout}
            style={{
              background: '#FEFEFE',
              color: '#BA9A91',
              border: '2px solid #BA9A91',
              padding: '12px 28px',
              borderRadius: '8px',
              fontSize: '16px',
              fontWeight: 'bold',
              cursor: 'pointer',
              transition: 'all 0.2s'
            }}
            onMouseOver={(e) => {
              e.currentTarget.style.background = '#BA9A91'
              e.currentTarget.style.color = 'white'
            }}
            onMouseOut={(e) => {
              e.currentTarget.style.background = '#FEFEFE'
              e.currentTarget.style.color = '#BA9A91'
            }}
          >
            Logout
          </button>
        </div>

        <div style={{
          marginTop: '30px',
          padding: '15px',
          background: '#EDECEC',
          borderRadius: '8px',
          fontSize: '13px',
          color: '#666',
          border: '1px solid #CCCCCC'
        }}>
          <strong>⚠️ Sandbox Mode:</strong> Your app is in sandbox mode.
          Make sure you're added as a test user in the TikTok Developer Portal.
        </div>
      </div>
    </div>
  )
}

export default SuccessPage
