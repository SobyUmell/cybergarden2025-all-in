# Backend

Microservices-based backend built with Go for a personal finance management system.

## Architecture

The backend consists of three independent microservices communicating via gRPC:

### 1. Bot Service (`/bot`)
Telegram Bot integration service that provides the user interface through Telegram.

**Key Features:**
- Telegram Bot API integration using webhook
- Telegram Mini App launcher with inline keyboard
- User commands handling (`/start`, `/help`)
- gRPC client for database communication

**Tech Stack:**
- `go-telegram/bot` - Telegram Bot API
- `google.golang.org/grpc` - gRPC communication
- `telegram-mini-apps/init-data-golang` - Telegram WebApp data validation

**Configuration:** `bot/config/config.yaml`

### 2. Database Service (`/database`)
Dedicated database management service with migration support.

**Key Features:**
- PostgreSQL database management
- User and transaction data storage
- Database migrations using golang-migrate
- gRPC server exposing database operations

**Schema:**
- `users` - User profiles (id)
- `transactions` - Financial transactions with user relationships

**Tech Stack:**
- `lib/pq` - PostgreSQL driver
- `golang-migrate/migrate` - Database migrations
- `google.golang.org/grpc` - gRPC server

**Migrations:** `database/internal/migrations/`

### 3. Manager Service (`/manager`)
HTTP REST API gateway that orchestrates bot and database services.

**Key Features:**
- RESTful API for web frontend
- Telegram WebApp authentication
- Transaction management (CRUD operations)
- CORS-enabled for frontend integration
- gRPC clients for bot and database services

**API Endpoints:**
- `GET /webapp/datametrics` - Financial metrics
- `GET /webapp/datahistory` - Transaction history
- `POST /webapp/addt` - Add transaction
- `POST /webapp/deletet` - Delete transaction
- `POST /webapp/updatet` - Update transaction

**Tech Stack:**
- `gin-gonic/gin` - HTTP web framework
- `gin-contrib/cors` - CORS middleware
- `google.golang.org/grpc` - gRPC clients

**Configuration:** `manager/config/config.yaml`

## Common Structure

Each service follows clean architecture principles:

```
service/
├── cmd/main.go              # Entry point
├── config/config.yaml       # Configuration
├── internal/
│   ├── app/                 # Application initialization
│   ├── config/              # Config parser
│   ├── handlers/            # Request handlers
│   ├── models/              # Data models
│   ├── repository/          # Data access layer
│   └── services/            # Business logic
├── Dockerfile               # Container image
├── go.mod                   # Dependencies
└── go.sum                   # Dependency checksums
```

## Dependencies

All services use:
- `PrototypeSirius/protos_service` - Shared protobuf definitions
- `PrototypeSirius/ruglogger` - Logging with error handling
- `ilyakaznacheev/cleanenv` - Configuration management
- `sirupsen/logrus` - Structured logging

## Docker Support

Each service includes a Dockerfile for containerization. Use `docker-compose.yml` in the root directory to run all services together.

## Development

Each service can be run independently:

```bash
cd backend/<service>
go run cmd/main.go
```

**Go Version:** 1.24.2
