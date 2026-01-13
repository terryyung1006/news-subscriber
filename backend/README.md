# News Subscriber Core

The core backend service that handles all business logic, authentication, and database access.

## Tech Stack

- **Language**: Go 1.25+
- **API**: gRPC
- **Databases**: PostgreSQL (with GORM), ChromaDB
- **ORM**: GORM for database migrations and queries
- **Proto**: Definitions in `news-subscriber-grpc-proto` (submodule)

## Quick Setup

### 1. Prerequisites
```bash
# Install Go 1.25+
# Install Docker Desktop
# Install protoc (for proto generation)
```

### 2. Start Database
```bash
make db-up
# Or manually:
docker-compose up -d
```

### 3. Run Migrations
Migrations run automatically when the server starts. To run manually:
```bash
make db-migrate
```

### 4. Configure

Edit `config.yaml`:
```yaml
google:
  client_id: "YOUR_GOOGLE_CLIENT_ID"  # Get from Google Cloud Console
```

### 5. Install Dependencies
```bash
go mod download
```

### 6. Run Server
```bash
make run
```

Server will start on `localhost:8080`

## Development Commands

```bash
make build              # Build the server binary
make run                # Build and run the server
make gen-proto          # Generate Go code from proto files
make update-proto       # Update submodule and regenerate proto
make generate-invites   # Generate new invitation codes
```

## Database Management

### Quick Commands
```bash
make db-up              # Start PostgreSQL container
make db-down            # Stop PostgreSQL container
make db-migrate         # Run database migrations (auto-runs on server start)
make db-reset           # Drop all tables and recreate them
make db-seed            # Seed initial invite codes
make db-fresh           # Complete reset: down, up, reset, and seed
```

### Manual Database Access
```bash
# Access database shell
docker exec -it news_subscriber_postgres psql -U user -d news_subscriber

# View users
SELECT * FROM users;

# View invite codes
SELECT * FROM invite_codes;
```

### Migration Details
- Migrations use **GORM AutoMigrate** - automatically creates/updates tables based on model definitions
- Migrations run automatically when the server starts
- Use `make db-reset` to drop and recreate all tables (⚠️ deletes all data)
- Use `make db-fresh` for a complete clean setup


## Project Structure

```
.
├── api/proto/v1/           # Generated proto code
├── migrations/             # Legacy SQL migrations (now using GORM)
├── spec/                   # Feature specifications
│   └── google-login/       # Google OAuth + Invite codes
├── src/
│   ├── main.go            # Entry point
│   ├── config/            # Configuration loading
│   ├── models/            # Data models
│   ├── repository/        # Data access layer
│   │   ├── postgres/      # PostgreSQL repos
│   │   └── chroma/        # ChromaDB client
│   └── service/           # gRPC service implementations
├── tools/                 # Utility scripts
│   ├── migrate/          # Database migration tool
│   └── generate-invites/ # Invite code generator
├── config.yaml            # App configuration
├── docker-compose.yaml    # PostgreSQL container
└── Makefile              # Build commands
```

## Features

See `spec/` folder for detailed feature documentation:
- **Google Login** → [`spec/google-login/`](./spec/google-login/)

## Environment

**Required**:
- PostgreSQL (via Docker)
- Google OAuth Client ID

**Optional**:
- ChromaDB (for future news embedding features)

## Testing

Default test invite codes:
- `WELCOME2024`
- `BETA_USER_001`
- `EARLY_ACCESS`

Generate more: `make generate-invites`
