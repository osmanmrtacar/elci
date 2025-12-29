import api from "./api";
import {
  AuthResponse,
  Platform,
  PlatformConnection,
  User,
} from "../types/user";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ||
  "http://localhost:8080";

export const authService = {
  // Initiate TikTok OAuth login
  initiateTikTokLogin: async () => {
    // If user is logged in, make API call with Authorization header
    // Backend will detect logged-in user and connect platform to existing account
    const token = localStorage.getItem("auth_token");
    if (token) {
      // Call login endpoint with auth header - backend will store user_id in oauth session
      try {
        const response = await api.get("/api/v1/auth/tiktok/login");
        // Check for URL in JSON response
        if (response.data && response.data.url) {
            window.location.href = response.data.url;
        }
      } catch (error: any) {
        throw error;
      }
    } else {
      // Not logged in - simple redirect for new user registration
      window.location.href = `${API_BASE_URL}/api/v1/auth/tiktok/login`;
    }
  },

  // Initiate X (Twitter) OAuth login
  initiateXLogin: async () => {
    const token = localStorage.getItem("auth_token");
    if (token) {
      try {
        const response = await api.get("/api/v1/auth/x/login");
        if (response.data && response.data.url) {
          window.location.href = response.data.url;
        }
      } catch (error: any) {
        throw error;
      }
    } else {
      window.location.href = `${API_BASE_URL}/api/v1/auth/x/login`;
    }
  },

  // Initiate Instagram OAuth login (via Facebook)
  initiateInstagramLogin: async () => {
    const token = localStorage.getItem("auth_token");
    if (token) {
      try {
        const response = await api.get("/api/v1/auth/instagram/login");
        if (response.data && response.data.url) {
          window.location.href = response.data.url;
        }
      } catch (error: any) {
        throw error;
      }
    } else {
      window.location.href = `${API_BASE_URL}/api/v1/auth/instagram/login`;
    }
  },

  // Generic platform login (backward compatibility)
  initiateLogin: (platform: Platform = "tiktok") => {
    if (platform === "x") {
      authService.initiateXLogin();
    } else {
      authService.initiateTikTokLogin();
    }
  },

  // Handle OAuth callback (called from CallbackPage)
  handleCallback: async (
    code: string,
    state: string,
  ): Promise<AuthResponse> => {
    const response = await api.get(`/api/v1/auth/tiktok/callback`, {
      params: { code, state },
    });
    return response.data;
  },

  // Get current user
  getCurrentUser: async (): Promise<User> => {
    const response = await api.get("/api/v1/auth/me");
    return response.data.user;
  },

  // Get connected platforms
  getConnectedPlatforms: async (): Promise<PlatformConnection[]> => {
    const response = await api.get("/api/v1/auth/platforms");
    return response.data.platforms;
  },

  // Disconnect a platform
  disconnectPlatform: async (platform: Platform): Promise<void> => {
    await api.delete(`/api/v1/auth/platforms/${platform}`);
  },

  // Logout
  logout: async (): Promise<void> => {
    await api.post("/api/v1/auth/logout");
    localStorage.removeItem("auth_token");
    localStorage.removeItem("user");
    localStorage.removeItem("connected_platforms");
  },

  // Save auth data to localStorage
  saveAuth: (token: string, user: User): void => {
    localStorage.setItem("auth_token", token);
    localStorage.setItem("user", JSON.stringify(user));
  },

  // Save connected platforms to localStorage
  saveConnectedPlatforms: (platforms: PlatformConnection[]): void => {
    localStorage.setItem("connected_platforms", JSON.stringify(platforms));
  },

  // Get saved connected platforms
  getStoredConnectedPlatforms: (): PlatformConnection[] => {
    const platformsStr = localStorage.getItem("connected_platforms");
    return platformsStr ? JSON.parse(platformsStr) : [];
  },

  // Get saved auth data
  getStoredAuth: (): { token: string | null; user: User | null } => {
    const token = localStorage.getItem("auth_token");
    const userStr = localStorage.getItem("user");
    const user = userStr ? JSON.parse(userStr) : null;
    return { token, user };
  },

  // Check if user is authenticated
  isAuthenticated: (): boolean => {
    return !!localStorage.getItem("auth_token");
  },

  // Check if a specific platform is connected
  isPlatformConnected: (platform: Platform): boolean => {
    const platforms = authService.getStoredConnectedPlatforms();
    return platforms.some((p) => p.platform === platform && p.is_active);
  },
};
