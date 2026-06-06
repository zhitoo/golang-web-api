# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

### Run the app
```bash
go run cmd/main.go
```

### Build & use the artisan CLI (migrations)
```bash
go build -o artisan artisan.go

./artisan migrate:up
./artisan migrate:down 1          # rolls back N steps (default 1)
./artisan migrate:create <name>   # creates timestamped up/down SQL files
./artisan migrate:force <version> # force-set schema version (fix dirty state)
./artisan serve                   # also starts the HTTP server
```

### Generate Swagger docs
```bash
swag init -g cmd/main.go
```

### Run tests
```bash
go test ./tests/unit/...
go test ./tests/integration/...
go test ./...
```

### Start infrastructure (Postgres, Redis, ELK)
```bash
docker-compose up -d
```

## Architecture

### Entry points
- `cmd/main.go` — boots config → Redis → Postgres → Gin HTTP server
- `artisan.go` — CLI tool for DB migrations and optionally `serve`

### Config system
Config is loaded by `config.GetConfig()`, which:
1. Reads `.env` from cwd (via godotenv)
2. Picks `config/config-{APP_ENV}.yml` (defaults to `config-local.yml`)
3. Unmarshals into a typed `Config` struct via Viper

To target a different environment set `APP_ENV=production` in `.env` or the shell.

### Module structure
Each feature lives under `app/modules/<name>/` with these files:
- `route.go` — registers Gin routes on a `*gin.RouterGroup`
- `handler.go` — thin Gin handlers: bind request → call service → return response
- `service.go` — business logic
- `request.go` — request structs with `binding:` validator tags
- `models/` — GORM model structs (user and location modules)

All modules are wired into `app/server.go → RegisterRoutes()`.

### HTTP response pattern
Always use the builder in `app/response`:
```go
response.NewReponse().SetResult(data).Json(c)
response.NewReponse().SetStatus(false).SetError(err).SetHttpStatusCode(http.StatusBadRequest).Json(c)
```

### Database
- **GORM** with Postgres driver (`database/db/postgres.go`)
- **Migrations** are raw SQL files under `database/migrations/` managed by `golang-migrate`
- `BaseModel` (`app/modules/user/models/base_model.go`) provides `CreatedAt/By`, `ModifiedAt/By`, `DeletedAt/By` via GORM hooks; it reads `UserId` from the gin context

### Cache (Redis)
`database/cache/redis.go` exposes generic helpers:
```go
cache.SetValue("key", value, ttl)   // JSON-serializes any type
cache.GetValue[T]("key")            // deserializes back to T
cache.DeleteValue("key")
```

### Logging
`pkg/logging` wraps zerolog. Every structured log call requires a **Category** and **SubCategory**:
```go
logger.Info(logging.General, logging.Startup, "message", map[logging.ExtraKey]any{...})
```
Categories and subcategories are typed constants in `pkg/logging/category.go`.  
Logs are written to a file (path from config) and shipped to Elasticsearch via Filebeat.

> **Known typo:** The type is spelled `Categoty` (not `Category`) throughout the logging package — do not "fix" this silently as it would be a breaking rename.

### i18n / translations
`lang/lang.go` exposes `lang.Trans(scope, key, ...values)`. Language is set from `cfg.App.Lang`. Translation maps are in `lang/fa.go` and `lang/en.go`.

### Validation
Custom validators are registered in `app/server.go → RegisterValidatores()`.  
`app/validations/` holds the validator functions (e.g. `irmobile` for Iranian mobile numbers) and `GetValidationErrors()` which maps `validator.ValidationErrors` to a translated `[]ValidationError` slice used by the response builder.

### Swagger
Swagger annotations live on handler functions. Run `swag init` after changing them to regenerate `docs/`.
