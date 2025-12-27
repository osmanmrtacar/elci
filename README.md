# TikTok Content Publisher

A social media content publishing platform that allows you to authenticate with TikTok and automatically post videos to your TikTok account.

## Features

- TikTok OAuth authentication
- Download videos from URLs
- Add captions and hashtags
- Automatic video posting to TikTok
- Post history and status tracking
- Secure token management with automatic refresh

## Technology Stack

**Backend:**
- Go (Golang) with Gin framework
- SQLite database
- JWT authentication
- TikTok Content Posting API

**Frontend:**
- React 18 with TypeScript
- Vite for build tooling
- React Router for navigation
- Axios for API calls
- React Hook Form for form handling

## Prerequisites

Before you begin, ensure you have the following installed:
- Go 1.21 or higher
- Node.js 18 or higher
- npm or yarn

You'll also need:
- A TikTok Developer account
- TikTok App with Login Kit and Content Posting API enabled

## TikTok Developer Setup

1. **Create TikTok Developer Account**
   - Visit https://developers.tiktok.com/
   - Sign up and verify your email

2. **Create New App**
   - Go to "My Apps" → "Create App"
   - Fill in app details (name, description, icon)
   - Select "Login Kit" and "Content Posting API"

3. **Configure OAuth**
   - Add redirect URI: `http://localhost:8080/api/v1/auth/tiktok/callback`
   - For production, add your production domain
   - Select scopes: `user.info.basic`, `video.publish`

4. **Get Credentials**
   - Copy your **Client Key** (client_key)
   - Copy your **Client Secret** (client_secret)
   - Save these for the backend configuration

## Installation

### 1. Clone the repository

```bash
git clone <repository-url>
cd sosyal
```

### 2. Backend Setup

```bash
# Navigate to backend directory
cd backend

# Install Go dependencies
make install

# Create .env file from example
cp .env.example .env

# Edit .env and add your TikTok credentials
nano .env  # or use your preferred editor
```

**Required .env configuration:**
```bash
TIKTOK_CLIENT_KEY=your_client_key_here
TIKTOK_CLIENT_SECRET=your_client_secret_here
JWT_SECRET=your_random_jwt_secret_here
```

### 3. Frontend Setup

```bash
# Navigate to frontend directory
cd ../frontend

# Install dependencies
npm install

# Create .env file from example
cp .env.example .env
```

## Running the Application

### Development Mode

**Terminal 1 - Backend:**
```bash
cd backend
make run
```

The backend server will start on `http://localhost:8080`

**Terminal 2 - Frontend:**
```bash
cd frontend
npm run dev
```

The frontend will start on `http://localhost:3000`

### Production Build

**Backend:**
```bash
cd backend
make build
./bin/server
```

**Frontend:**
```bash
cd frontend
npm run build
npm run preview
```

## Project Structure

```
sosyal/
├── backend/
│   ├── cmd/server/          # Application entry point
│   ├── internal/
│   │   ├── api/             # HTTP handlers and routes
│   │   ├── config/          # Configuration management
│   │   ├── database/        # Database and models
│   │   └── services/        # Business logic
│   ├── data/                # SQLite database
│   ├── go.mod               # Go dependencies
│   ├── .env.example         # Environment variables template
│   └── Makefile             # Build commands
│
└── frontend/
    ├── src/
    │   ├── components/      # React components
    │   ├── context/         # React context providers
    │   ├── hooks/           # Custom React hooks
    │   ├── pages/           # Page components
    │   ├── services/        # API client services
    │   └── types/           # TypeScript type definitions
    ├── package.json         # npm dependencies
    ├── vite.config.ts       # Vite configuration
    └── .env.example         # Environment variables template
```

## API Endpoints

### Authentication
- `GET /api/v1/auth/tiktok/login` - Initiate TikTok OAuth flow
- `GET /api/v1/auth/tiktok/callback` - Handle OAuth callback
- `POST /api/v1/auth/logout` - Logout user
- `GET /api/v1/auth/me` - Get current user info

### Posts
- `POST /api/v1/posts` - Create new post (requires auth)
- `GET /api/v1/posts` - Get user's post history (requires auth)
- `GET /api/v1/posts/:id` - Get specific post details (requires auth)
- `GET /api/v1/posts/:id/status` - Poll post status (requires auth)

## Usage

1. **Login with TikTok**
   - Open `http://localhost:3000`
   - Click "Login with TikTok"
   - Authorize the app on TikTok
   - You'll be redirected back to the dashboard

2. **Post a Video**
   - Enter a video URL (must be accessible)
   - Add a caption with optional hashtags
   - Click "Post to TikTok"
   - Monitor the status (downloading → uploading → published)

3. **View Post History**
   - See all your past posts
   - Check their status (pending, published, failed)
   - View timestamps and TikTok post links

## Video Requirements

- **Formats**: MP4, WEBM, MOV
- **Max Size**: 287.6 MB
- **Duration**: 3 seconds - 10 minutes
- **Resolution**: Minimum 360p, recommended 720p or 1080p
- **Aspect Ratio**: 9:16 (vertical), 16:9 (horizontal), or 1:1 (square)

## Development Commands

### Backend

```bash
make help          # Show available commands
make run           # Run the server
make build         # Build binary
make test          # Run tests
make clean         # Clean build artifacts
make format        # Format code
make db-reset      # Reset database
```

### Frontend

```bash
npm run dev        # Start development server
npm run build      # Build for production
npm run preview    # Preview production build
npm run lint       # Run linter
```

## Troubleshooting

### OAuth Issues
- Verify redirect_uri matches exactly in TikTok Developer Portal
- Check client_key and client_secret are correct
- Ensure CORS is properly configured

### Video Upload Issues
- Verify the video URL is publicly accessible
- Check video meets TikTok's requirements (size, format, duration)
- Check server logs for detailed error messages

### Token Issues
- Tokens automatically refresh when expired
- If you see authentication errors, try logging out and back in
- Check JWT_SECRET is set in backend .env

### Database Issues
- Reset database: `make db-reset` (backend directory)
- Database will be recreated on next server start

## Security Notes

- Never commit `.env` files to version control
- Keep your TikTok Client Secret secure
- Use HTTPS in production
- Implement rate limiting for production deployments
- Regularly update dependencies

## Future Enhancements

- Instagram and X (Twitter) integration
- Scheduled posting
- Video editing capabilities
- Direct file upload
- Analytics dashboard
- Multiple account management
- Bulk posting

## License

MIT

## Support

For issues and questions:
- Check the troubleshooting section above
- Review TikTok Developer documentation
- Check application logs for detailed error messages
