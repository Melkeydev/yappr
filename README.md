<p align="center">
  <h1 align="center"><b>Yappr</b></h1>
</p>
<p align="center">Real-time chat rooms that disappear after 24 hours.</p>
<p align="center">
  <a href="https://github.com/melkeydev/go-chat-app/blob/main/LICENSE"><img alt="License" src="https://img.shields.io/github/license/melkeydev/go-chat-app?style=flat-square" /></a>
  <a href="https://go.dev"><img alt="Go" src="https://img.shields.io/badge/Go-1.24-00ADD8?style=flat-square&logo=go" /></a>
  <a href="https://react.dev"><img alt="React" src="https://img.shields.io/badge/React-18-61DAFB?style=flat-square&logo=react" /></a>
</p>

---

### Features

- **Ephemeral Rooms** - All rooms automatically expire after 24 hours
- **Daily Topics** - Three pinned rooms with fresh topics from HackerNews and Reddit
- **Real-time Chat** - WebSocket-based messaging with instant updates
- **User Accounts** - Optional registration with room creation limits
- **Anonymous Access** - Join and chat without creating an account

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

### Tech Stack

- **Backend**: Go, Chi Router, Gorilla WebSocket, PostgreSQL
- **Frontend**: React, TypeScript, Vite, Tailwind CSS
- **Infrastructure**: Docker, Docker Compose

### Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

**Built with ❤️ by [Melkey](https://github.com/melkeydev)**