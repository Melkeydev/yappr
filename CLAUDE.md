# Go Chat App - Project Overview

## Architecture

### Backend (Go)
- **Framework**: Chi router with Gorilla WebSocket
- **Database**: PostgreSQL with automatic migrations
- **Authentication**: JWT tokens stored in HTTP-only cookies
- **Real-time**: WebSocket-based messaging with in-memory message history

### Frontend (React + TypeScript)
- **Build**: Vite
- **Styling**: Tailwind CSS
- **Forms**: react-hook-form
- **HTTP**: Axios with credentials

## Key Features

### 1. Room Management
- **24-hour expiration**: ALL rooms expire 24 hours after creation
- **Room limits**: 
  - Global limit: configurable via `MAX_ROOMS` env var (default: 50)
  - Per user limit: Logged-in users can only create 1 active room at a time
  - Anonymous users cannot create rooms
- **Database persistence**: Rooms and messages stored in PostgreSQL

### 2. User System
- **Authenticated users**: 
  - Can create one room at a time
  - Can change username via profile page
  - JWT authentication with 24-hour tokens
- **Anonymous users**:
  - Auto-generated username: `anonymousUser_[6-chars]`
  - Username persisted in localStorage
  - Cannot create rooms
  - Can join and chat in existing rooms

### 3. Pinned Topic Rooms
- **3 pinned rooms** with daily topics:
  1. "Tech Talk" - HackerNews top story
  2. "World News" - Reddit r/worldnews top post
  3. "Fun Facts" - Reddit r/todayilearned top post
- **Synchronized expiration**: All pinned rooms expire at midnight UTC
- **Auto-refresh**: New topics fetched from APIs when rooms recreate
- **Visual distinction**: Gradient background, pinned icon, topic preview
- **Topic display**: Banner in chat room with link to source

## Database Schema

### Tables
1. **users**: id, username, email, password_hash, created_at, updated_at
2. **rooms**: id, name, creator_id, is_pinned, topic_*, created_at, expires_at
3. **messages**: id, room_id, user_id, username, content, is_system, created_at

### Key Indexes
- `idx_rooms_expires_at`: For cleanup job
- `idx_rooms_creator_id`: For user room limits
- `idx_rooms_is_pinned`: For pinned room queries

## Background Jobs

### Room Cleanup Job
- Runs every 5 minutes
- Deletes expired rooms (CASCADE deletes messages)
- Checks and recreates pinned rooms if needed

### Pinned Room Management
- On startup: Creates pinned rooms if missing
- On cleanup: Recreates pinned rooms with fresh topics
- Topics fetched from:
  - HackerNews API: `https://hacker-news.firebaseio.com/v0/`
  - Reddit API: `https://www.reddit.com/r/[subreddit]/top.json`

## API Endpoints

### Authentication
- `POST /api/users/signup` - Register new user
- `POST /api/users/login` - Login (sets JWT cookie)
- `GET /api/users/logout` - Logout (clears cookie)
- `PUT /api/users/username` - Change username (auth required)

### WebSocket/Rooms
- `POST /ws/createRoom` - Create room (auth affects limits)
- `GET /ws/joinRoom/{roomId}` - WebSocket upgrade
- `GET /ws/getRooms` - List all active rooms
- `GET /ws/getClients/{roomId}` - List room participants

## Environment Variables
- `secretKey`: JWT signing key (required)
- `MAX_ROOMS`: Maximum total rooms (default: 50)
- Database connection handled by Docker Compose

## Common Issues & Solutions

### JWT Cookie Not Setting
- Removed restrictive `Domain: "localhost"`
- Added `SameSite: http.SameSiteLaxMode`
- Ensure `withCredentials: true` in axios

### Room Topics Not Showing
- Check `RoomRes` struct includes all topic fields
- Ensure WebSocket `Room` struct has topic fields
- Verify database migration added topic columns

### User Can't Create Room
- Check if user already has active room
- Verify JWT middleware is passing user ID
- Check global room limit hasn't been reached

## Development Commands

### Backend
```bash
cd server
go run main.go
```

### Frontend
```bash
cd client
npm run dev
```

### Database Migrations
Migrations run automatically on server startup. Manual migration:
```bash
cd server
go run db/migrations/migrate.go up
```

## Testing Checklist
- [ ] Anonymous users can join but not create rooms
- [ ] Logged-in users can create only 1 room
- [ ] Rooms expire after 24 hours
- [ ] Pinned rooms refresh at midnight UTC
- [ ] Topic banners show in pinned room chats
- [ ] Profile page allows username changes
- [ ] JWT authentication persists across refreshes