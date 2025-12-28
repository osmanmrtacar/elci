import { createContext, useContext, useState, useEffect, ReactNode, useCallback } from 'react'
import { User, Platform, PlatformConnection } from '../types/user'
import { authService } from '../services/authService'

interface AuthContextType {
  user: User | null
  connectedPlatforms: PlatformConnection[]
  isAuthenticated: boolean
  isLoading: boolean
  loginTikTok: () => Promise<void>
  loginX: () => Promise<void>
  loginInstagram: () => Promise<void>
  logout: () => void
  disconnectPlatform: (platform: Platform) => Promise<void>
  refreshPlatforms: () => Promise<void>
  setUser: (user: User) => void
  isPlatformConnected: (platform: Platform) => boolean
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export const useAuth = () => {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}

interface AuthProviderProps {
  children: ReactNode
}

export const AuthProvider = ({ children }: AuthProviderProps) => {
  const [user, setUserState] = useState<User | null>(null)
  const [connectedPlatforms, setConnectedPlatforms] = useState<PlatformConnection[]>([])
  const [isLoading, setIsLoading] = useState(true)

  const refreshPlatforms = useCallback(async () => {
    if (!authService.isAuthenticated()) {
      setConnectedPlatforms([])
      return
    }

    try {
      const platforms = await authService.getConnectedPlatforms()
      setConnectedPlatforms(platforms)
      authService.saveConnectedPlatforms(platforms)
    } catch (error) {
      console.error('Failed to fetch connected platforms:', error)
      // Load from cache if API call fails
      setConnectedPlatforms(authService.getStoredConnectedPlatforms())
    }
  }, [])

  useEffect(() => {
    // Check for stored auth on mount
    const { token, user: storedUser } = authService.getStoredAuth()
    if (token && storedUser) {
      setUserState(storedUser)
      // Load stored platforms initially
      setConnectedPlatforms(authService.getStoredConnectedPlatforms())
      // Then refresh from server
      refreshPlatforms()
    }
    setIsLoading(false)
  }, [refreshPlatforms])

  const loginTikTok = async () => {
    await authService.initiateTikTokLogin()
  }

  const loginX = async () => {
    await authService.initiateXLogin()
  }

  const loginInstagram = async () => {
    await authService.initiateInstagramLogin()
  }

  const logout = async () => {
    try {
      await authService.logout()
    } catch (error) {
      console.error('Logout error:', error)
    } finally {
      setUserState(null)
      setConnectedPlatforms([])
    }
  }

  const disconnectPlatform = async (platform: Platform) => {
    try {
      await authService.disconnectPlatform(platform)
      await refreshPlatforms()
    } catch (error) {
      console.error(`Failed to disconnect ${platform}:`, error)
      throw error
    }
  }

  const setUser = (user: User) => {
    setUserState(user)
  }

  const isPlatformConnected = (platform: Platform): boolean => {
    return connectedPlatforms.some(p => p.platform === platform && p.is_active)
  }

  const value: AuthContextType = {
    user,
    connectedPlatforms,
    isAuthenticated: !!user,
    isLoading,
    loginTikTok,
    loginX,
    loginInstagram,
    logout,
    disconnectPlatform,
    refreshPlatforms,
    setUser,
    isPlatformConnected,
  }

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}
