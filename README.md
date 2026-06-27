# Golang Web API

A REST API built with Go, Gin, GORM, and OTP-based authentication.

## Prerequisites

- Go 1.26+
- Docker & Docker Compose
- [air](https://github.com/air-verse/air) — optional, for hot reload
- [Bruno](https://www.usebruno.com) — optional, for API testing

## Setup

### Local (Go installed on machine)

```bash
cp .env.sample .env

docker-compose up -d postgres redis   # start only DB + cache
make migrate-up                       # run migrations
make serve                            # start server
```

For hot reload: `make serve-dev` (requires `air`)

### Fully Dockerized (no Go needed locally)

```bash
cp .env.sample .env

docker-compose up -d postgres redis   # start DB + cache first
make docker-dev                       # build image + start app with hot reload
```

On first run, Docker downloads the Go image and installs dependencies — subsequent starts are fast. Code changes trigger automatic reload inside the container.

### Environment variables (`.env`)

| Variable | Default | Description |
|---|---|---|
| `APP_ENV` | `local` | Config profile: `local`, `docker`, `production` |
| `APP_PORT` | `5050` | Host port exposed by the app container |
| `POSTGRES_PORT` | `5434` | Host port for Postgres container |
| `POSTGRES_USER` | `postgres` | Postgres username |
| `POSTGRES_PASSWORD` | `secret` | Postgres password |
| `POSTGRES_DB` | `golang_web_api` | Postgres database name |
| `REDIS_PORT` | `6385` | Host port for Redis container |
| `ELASTICSEARCH_PORT` | `9200` | Host port for Elasticsearch HTTP |
| `ELASTICSEARCH_TRANSPORT_PORT` | `9300` | Host port for Elasticsearch transport |
| `KIBANA_PORT` | `5601` | Host port for Kibana |
| `PGADMIN_PORT` | `8589` | Host port for pgAdmin |
| `PGADMIN_EMAIL` | `admin@main.com` | pgAdmin login email |
| `PGADMIN_PASSWORD` | `secret` | pgAdmin login password |
| `ELASTIC_PASSWORD` | — | Elasticsearch password |
| `KIBANA_SYSTEM_PASSWORD` | — | Kibana system password |
| `FILEBEAT_INTERNAL_PASSWORD` | — | Filebeat internal password |

App config (JWT secrets, log level, etc.) lives in `config/config-{APP_ENV}.yml`.

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
