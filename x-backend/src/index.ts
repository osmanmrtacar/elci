// Load environment variables FIRST before any other imports
import dotenv from 'dotenv'
dotenv.config()

import express from 'express'
import cors from 'cors'
import cookieParser from 'cookie-parser'
import authRoutes from './routes/auth.routes'
import postRoutes from './routes/post.routes'

const app = express()
const PORT = process.env.PORT || 8888

// Middleware
app.use(
  cors({
    origin: process.env.CORS_ORIGIN?.split(',') || ['http://localhost:3000'],
    credentials: true,
  })
)
app.use(express.json())
app.use(express.urlencoded({ extended: true }))
app.use(cookieParser())

// Health check
app.get('/health', (req, res) => {
  res.json({ status: 'ok', service: 'x-backend' })
})

// Routes
app.use('/api/auth', authRoutes)
app.use('/api/posts', postRoutes)

// Start server
app.listen(PORT, () => {
  console.log(`🚀 X Backend running on port ${PORT}`)
  console.log(`📍 Health check: http://localhost:${PORT}/health`)
  console.log(`🔐 Auth endpoint: http://localhost:${PORT}/api/auth/x/login`)
})
