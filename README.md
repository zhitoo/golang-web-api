# Golang Web API

A REST API built with Go, Gin, GORM, and OTP-based authentication.

## Prerequisites

- Go 1.26+
- Docker & Docker Compose
- [air](https://github.com/air-verse/air) — optional, for hot reload
- [Bruno](https://www.usebruno.com) — optional, for API testing

## Setup

```bash
cp .env.sample .env
# fill in values in .env

docker-compose up -d   # start Postgres, Redis, ELK
make migrate-up        # build artisan + run migrations
make serve             # build artisan + start server
```

### Environment variables (`.env`)

| Variable | Description |
|---|---|
| `APP_ENV` | Config profile to load: `local` (default), `docker`, `production` |
| `ELASTIC_PASSWORD` | Elasticsearch password |
| `KIBANA_SYSTEM_PASSWORD` | Kibana system password |
| `FILEBEAT_INTERNAL_PASSWORD` | Filebeat internal password |

Config values (ports, DB credentials, JWT secrets, etc.) live in `config/config-{APP_ENV}.yml`.

## Artisan CLI

All `make` targets build artisan automatically for the current OS (`artisan` on Linux/macOS, `artisan.exe` on Windows):

| Make target | Equivalent artisan command | Description |
|---|---|---|
| `make build` | — | Build the artisan binary |
| `make serve` | `./artisan serve` | Start the HTTP server |
| `make serve-dev` | `./artisan serve:dev` | Start with hot reload (requires `air`) |
| `make swagger` | `./artisan swagger:generate` | Regenerate Swagger docs |
| `make migrate-up` | `./artisan migrate:up` | Run all pending migrations |
| `make migrate-down` | `./artisan migrate:down 1` | Roll back 1 migration step |

Additional artisan commands (run directly after `make build`):

```bash
./artisan migrate:down [steps]      # roll back N steps
./artisan migrate:create <name>     # create a new migration file
./artisan migrate:force <version>   # force-set schema version (fix dirty state)
```

**Windows (without make):** run `build.bat` to build `artisan.exe`, then use `.\artisan.exe <command>`.

## Running Tests

```bash
go test ./...
go test ./tests/unit/...
go test ./tests/integration/...
go test -run TestFooBar ./tests/unit/...   # single test
```

## API Docs

Two options for exploring the API:

**Swagger UI** — available at `http://localhost:5050/swagger/index.html` after starting the server. Regenerate after changing handler annotations:

```bash
make swagger
```

**Bruno collection** — open `docs/apiDocs/` in Bruno and select the `Local` environment (`http://localhost:5050/api/v1`). Folders: `Health`, `OTP`, `Users`.

## Auth Flow

1. `POST /api/v1/otp/send` — send OTP to mobile
2. `POST /api/v1/otp/verify` — verify OTP
3. `POST /api/v1/users/login` — login with mobile/email + OTP or password → returns `access_token` + `refresh_token`
4. `POST /api/v1/users/refresh-token` — get a new token pair

Protected endpoints require `Authorization: Bearer <access_token>`.
