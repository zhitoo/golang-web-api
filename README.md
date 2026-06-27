# Golang Web API

A REST API built with Go, Gin, GORM, and OTP-based authentication.

## Prerequisites

- Go 1.26+
- Docker & Docker Compose

## Setup

```bash
cp .env.sample .env
# fill in values in .env

docker-compose up -d        # start Postgres, Redis, ELK
./artisan migrate:up        # run migrations
go run cmd/main.go          # start the server
```

For hot reload during development:
```bash
air
```

## Artisan CLI

Build once, then use `./artisan <command>`:

```bash
go build -o artisan artisan.go
```

| Command | Description |
|---|---|
| `./artisan serve` | Start the HTTP server |
| `./artisan serve:dev` | Start with hot reload (requires `air`) |
| `./artisan migrate:up` | Run all pending migrations |
| `./artisan migrate:down [steps]` | Roll back N migration steps (default: 1) |
| `./artisan migrate:create <name>` | Create a new timestamped migration file |
| `./artisan migrate:force <version>` | Force-set schema version (fixes dirty state) |
| `./artisan swagger:generate` | Regenerate Swagger docs from annotations |

## Running Tests

```bash
go test ./...
go test ./tests/unit/...
go test ./tests/integration/...
go test -run TestFooBar ./tests/unit/...   # single test
```

## API Docs

After starting the server, Swagger UI is available at:

```
http://localhost:5050/swagger/index.html
```

To regenerate docs after changing handler annotations:

```bash
./artisan swagger:generate
```

## Auth Flow

1. `POST /api/v1/otp/send` — send OTP to mobile
2. `POST /api/v1/otp/verify` — verify OTP
3. `POST /api/v1/users/login` — login with mobile/email + OTP or password → returns `access_token` + `refresh_token`
4. `POST /api/v1/users/refresh-token` — get a new token pair

Protected endpoints require `Authorization: Bearer <access_token>`.
