import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// Get allowed hosts from environment variable
// Format: comma-separated list like "localhost,mydomain.com,*.example.com"
const allowedHosts = process.env.VITE_ALLOWED_HOSTS
  ? process.env.VITE_ALLOWED_HOSTS.split(',').map(h => h.trim())
  : ['localhost', '.localhost']

export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000,
    host: true, // Listen on all addresses (0.0.0.0)
    allowedHosts: allowedHosts,
    proxy: {
      '/api': {
        // Use VITE_API_BASE_URL for proxy target (dev mode only)
        // In production, proxy is not used - frontend calls API directly
        target: process.env.VITE_API_BASE_URL || 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})
