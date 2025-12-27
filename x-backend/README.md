# X (Twitter) Backend API

TypeScript backend for X (Twitter) OAuth 2.0 authentication and posting.

## Setup

### 1. Install Dependencies

```bash
npm install
```

### 2. Configure X Developer Account

1. Go to [X Developer Portal](https://developer.x.com/en/portal/dashboard)
2. Create a new app or use existing
3. Enable OAuth 2.0
4. Set **Callback URI**: `https://adsl-cabinet-acceptance-investigate.trycloudflare.com/api/auth/x/callback`
5. Enable these permissions:
   - Read and write Tweets
   - Read users
6. Copy your **Client ID** and **Client Secret**

### 3. Create .env File

```bash
cp .env.example .env
```

Edit `.env` and add your X credentials:

```
X_CLIENT_ID=your_client_id_here
X_CLIENT_SECRET=your_client_secret_here
```

### 4. Start Server

```bash
npm run dev
```

Server runs on port 8888.

### 5. Start Cloudflare Tunnel

In another terminal:

```bash
cloudflared tunnel --url http://localhost:8888
```

Use the tunnel URL (e.g., `https://adsl-cabinet-acceptance-investigate.trycloudflare.com`) as your X callback URI.

## API Endpoints

### Authentication

**Login with X**
```
GET /api/auth/x/login
```
Redirects to X authorization page.

**OAuth Callback**
```
GET /api/auth/x/callback?code=xxx&state=xxx
```
Handles X OAuth callback, exchanges code for tokens.

**Get User Info**
```
GET /api/auth/x/me/:userId
```
Returns authenticated user info.

### Posts

**Create Post (Tweet)**
```
POST /api/posts
Content-Type: application/json

{
  "userId": "1234567890",
  "text": "Hello from the API!"
}
```

**Get User Tweets**
```
GET /api/posts/:userId
```

## How OAuth 2.0 Works

1. **User clicks "Login with X"** → Redirects to `/api/auth/x/login`
2. **Backend generates OAuth URL** with PKCE (code challenge)
3. **User authorizes** on X
4. **X redirects back** with authorization code
5. **Backend exchanges code** for access + refresh tokens
6. **Tokens stored** in memory (use database in production)
7. **User redirected** to frontend with user info

## PKCE Flow

X uses **PKCE** (Proof Key for Code Exchange) for security:
- Generate random `code_verifier`
- Create `code_challenge` = SHA256(code_verifier)
- Send challenge to X
- Exchange code + verifier for tokens

## Token Refresh

Tokens expire! The backend automatically refreshes them:
- Check if token expired before each request
- Use refresh token to get new access token
- Update stored tokens

## Testing with cURL

```bash
# 1. Login (in browser)
open https://adsl-cabinet-acceptance-investigate.trycloudflare.com/api/auth/x/login

# 2. After login, get your userId from the callback URL

# 3. Create a post
curl -X POST https://adsl-cabinet-acceptance-investigate.trycloudflare.com/api/posts \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "YOUR_USER_ID",
    "text": "Hello from the API!"
  }'
```

## Notes

- **In-memory storage**: Tokens are stored in memory (lost on restart). Use a database for production.
- **HTTPS required**: X requires HTTPS redirect URIs. Use Cloudflare tunnel for local development.
- **Rate limits**: X has rate limits. Be mindful when testing.
