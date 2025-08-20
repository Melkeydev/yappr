# Go Chat App - Build & Style Guide

## Build Commands
```bash
# Backend (from /server)
go run main.go                    # Run server locally
go build -o app main.go           # Build binary
go test ./...                     # Run all tests
go test -v ./internal/service/user -run TestCreateUser  # Run single test
go test -cover ./...              # Run tests with coverage
go mod tidy                       # Clean up dependencies
go fmt ./...                      # Format all Go files

# Frontend (from /client)
npm run dev                       # Start dev server (Vite on port 3000)
npm run build                     # Build for production
npm run lint                      # Run ESLint
npm run preview                   # Preview production build
npx tsc --noEmit                  # TypeScript type check

# Full Stack
docker-compose up                 # Run all services (dev)
docker-compose -f docker-compose-prod.yml up  # Production
docker-compose down -v            # Stop and remove volumes
```

## Go Code Style
- **Imports**: stdlib → blank line → external → blank line → internal (aliased: `coreHandler "path"`)
- **Error handling**: `if err != nil { return nil, fmt.Errorf("action failed: %w", err) }`
- **PostgreSQL errors**: Check `pgErr.Code == "23505"` for unique violations
- **Naming**: PascalCase exports, camelCase private, single-letter receivers (h, s, r)
- **Constructors**: `func NewXxx(deps) *Xxx { return &Xxx{field: deps} }`
- **Context**: `ctx, cancel := context.WithTimeout(ctx, 5*time.Second); defer cancel()`
- **Layers**: Repository (DB ops) → Service (business logic, JWT) → Handler (HTTP)
- **Logging**: `log.Printf("Component.Method - Action: %v", details)`

## TypeScript/React Style
- **Imports**: third-party → contexts → components → hooks → api → utils (no blank lines)
- **Components**: Function components, default export `export default function Name()`
- **Types**: `type User = {...} | null`, inline for useState: `useState<Room[]>([])`
- **Props**: Destructure in params: `function Component({ prop1, prop2 }: Props)`
- **Hooks**: `use` prefix, early returns for guards: `if (!user) return;`
- **Async**: async/await with try/catch, axios with `withCredentials: true`
- **WebSocket**: Reconnection with exponential backoff, cleanup in useEffect return
- **Styling**: Tailwind only, `clsx(baseClasses, { 'conditional': isTrue })`