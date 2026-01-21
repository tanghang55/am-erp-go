# am-erp-go

Go backend for Amazon Local ERP system.

## Tech Stack

- Go 1.22+
- Gin (HTTP framework)
- GORM (ORM)
- MySQL
- JWT (golang-jwt/jwt/v5)
- bcrypt (password hashing)

## Project Structure

```
am-erp-go/
├── cmd/server/main.go           # Entry point
├── internal/
│   ├── infrastructure/
│   │   ├── config/              # Configuration
│   │   ├── db/                  # Database connection
│   │   ├── auth/                # JWT & middleware
│   │   └── router/              # Route registration
│   └── module/
│       └── identity/            # Identity module
│           ├── domain/          # Entities & interfaces
│           ├── repository/      # Data access (GORM)
│           ├── usecase/         # Business logic
│           └── delivery/http/   # HTTP handlers
└── migrations/                  # SQL migrations
```

## Build & Run

```bash
# Build
go build -o am-erp-go ./cmd/server

# Run (set env vars first)
export DB_HOST=localhost
export DB_PASSWORD=321456
export DB_DATABASE=am_erp
go run ./cmd/server
```

## API Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | /health | No | Health check |
| POST | /api/v1/auth/login | No | User login |
| GET | /api/v1/auth/me | Yes | Get current user |
| GET | /api/v1/menus/tree | Yes | Get menu tree |

## Database Migration

Run `migrations/001_init.sql` to create tables and seed admin user.

Default admin credentials:
- Username: `admin`
- Password: `admin123`

## Adding New Modules

Follow DDD pattern:
1. Create `internal/module/<name>/domain/` - entities, interfaces
2. Create `internal/module/<name>/repository/` - GORM implementations
3. Create `internal/module/<name>/usecase/` - business logic
4. Create `internal/module/<name>/delivery/http/` - handlers
5. Register routes in `internal/infrastructure/router/router.go`
