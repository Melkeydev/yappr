<p align="center">
  <img src="client/public/go-chat-logo.png" alt="Yappr Logo" width="80" height="80">
  <h1 align="center"><b>Yappr</b></h1>
</p>
<p align="center">Real-time chat rooms that disappear after 24 hours.</p>

> **This project is archived.** Yappr is no longer actively maintained or hosted, but the source code remains available here for reference and learning.

---

### Features

- **Ephemeral Rooms** - All rooms automatically expire after 24 hours
- **Daily Topics** - Three pinned rooms with fresh topics from HackerNews and Reddit
- **Real-time Chat** - WebSocket-based messaging with instant updates
- **User Accounts** - Optional registration with room creation limits
- **Anonymous Access** - Join and chat without creating an account

### Tech Stack

- **Backend**: Go, Chi Router, Gorilla WebSocket, PostgreSQL
- **Frontend**: React, TypeScript, Vite, Tailwind CSS
- **Infrastructure**: Docker, Docker Compose

### Quick Start

```bash
# Clone the repository
git clone https://github.com/melkeydev/go-chat-app.git
cd go-chat-app

# Start with Docker Compose
docker-compose up

# Visit http://localhost:3000
```

### Development

#### Backend (Go)

```bash
cd server
go mod download
go run main.go
```

#### Frontend (React)

```bash
cd client
npm install
npm run dev
```

### Environment Variables

#### Server

```env
secretKey=your-jwt-secret
MAX_ROOMS=50
REDDIT_CLIENT_ID=your-reddit-client-id
REDDIT_CLIENT_SECRET=your-reddit-client-secret
```

#### Client

```env
VITE_API_BASE_URL=http://localhost:8080
VITE_WEBSOCKET_URL=ws://localhost:8080
```

### License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

**Built by [Melkey](https://github.com/melkeydev)**
