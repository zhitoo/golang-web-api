# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

### Run the app
```bash
go run cmd/main.go
```

### Hot reload (development)
```bash
air   # uses .air.toml; builds to tmp/main
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
$(go env GOPATH)/bin/swag init -g cmd/main.go
# or via artisan (after go build -o artisan artisan.go):
./artisan swagger:generate
```

### Run tests
```bash
go test ./tests/unit/...
go test ./tests/integration/...
go test ./...
go test -run TestFooBar ./tests/unit/...   # single test
```

### Start infrastructure (Postgres, Redis, ELK)
```bash
docker-compose up -d
```

## First-time setup

Copy `.env.sample` to `.env` and fill in values. Config is then loaded from `config/config-{APP_ENV}.yml` (defaults to `config-local.yml`).

## Architecture

### Entry points
- `cmd/main.go` — boots config → Redis → Postgres → Gin HTTP server
- `artisan.go` — CLI tool for DB migrations and optionally `serve`

### Config system
`config.GetConfig()` reads `.env` (via godotenv), then picks `config/config-{APP_ENV}.yml`, and unmarshals into a typed `Config` struct via Viper. Set `APP_ENV=production` in `.env` or the shell to switch environments.

JWT durations (`AccessTokenExpireDuration`, `RefreshTokenExpireDuration`) are stored as plain integers in the YAML and treated as seconds (multiplied by `time.Second` in `auth/service.go`).

### Module structure
Each feature lives under `app/modules/<name>/`. The typical files are:
- `route.go` — registers Gin routes on a `*gin.RouterGroup`
- `handler.go` — thin Gin handlers: bind request → call service → return response
- `service.go` — business logic
- `request.go` — request structs with `binding:` validator tags
- `models/` — GORM model structs

Not all modules carry all files: `auth` is service-only (JWT generation/verification, used by the `user` module); `health` adds a middleware but has no handler file; `location` has only models and a route stub with no endpoints registered yet.

All modules are wired into `app/server.go → RegisterRoutes()`.

### Module logger pattern
Each module declares a package-level logger and initializes it in `init()`:
```go
var log logging.ScopedLogger

func init() {
    cfg := config.GetConfig()
    log = logging.NewLogger(cfg).With(logging.Internal, logging.Api)
}
```
Pass the scoped logger's category/subcategory pair that best describes the module (see `pkg/logging/category.go`).

### OTP-based auth flow
1. `POST /api/v1/otp/send` — generates a 6-digit OTP, stores it in Redis under `otp:<mobile>` for 2 minutes.
2. `POST /api/v1/otp/verify` — checks OTP and deletes it from Redis on success.
3. `POST /api/v1/users/login` — supports four login modes:
   - mobile + OTP: verifies OTP, auto-creates user if not found
   - mobile + password: password auth only, user must exist
   - email + password: password auth only, user must exist
   - email + OTP: verifies OTP against the mobile linked to the email account

### Authentication middleware
`app/middlewares/auth_middleware.go → Authentication(cfg)` validates the `Bearer` token and injects these keys into the Gin context for downstream handlers:
- `"UserId"` (string)
- `"Username"` (string)
- `"Roles"` ([]string)

### HTTP response pattern
Always use the builder in `app/response`:
```go
response.NewResponse().SetResult(data).Json(c)
response.NewResponse().SetStatus(false).SetError(err).SetHttpStatusCode(http.StatusBadRequest).Json(c)
```
`SetError` automatically sets `success: false`, populates `validation_errors` if the error is a `validator.ValidationErrors`, and defaults the HTTP status to 400 if it is still 200.

### Database
- **GORM** with Postgres driver (`database/db/postgres.go`)
- **Migrations** are raw SQL files under `database/migrations/` managed by `golang-migrate`
- `BaseModel` (`app/modules/user/models/base_model.go`) provides `CreatedAt/By`, `ModifiedAt/By`, `DeletedAt/By` via GORM hooks; it reads `UserId` from the Gin context

> **Known typo:** The timestamp fields (`CreatedAt`, `ModifiedAt`, `DeletedAt`) use the tag prefix `gotm:` instead of `gorm:`, so their timezone struct tags are silently ignored by GORM. Do not "fix" this without updating the migration SQL as well.

### Cache (Redis)
`database/cache/redis.go` exposes generic helpers:
```go
cache.SetValue("key", value, ttl)   // JSON-serializes any type
cache.GetValue[T]("key")            // deserializes back to T
cache.DeleteValue("key")
```

### Utilities (`pkg/utils`)
- `HashPassword(password)` / `VerifyPassword(password, hash)` — bcrypt wrappers
- `RandomString(length)` — random alphanumeric+symbol string

### Rate limiting
Two implementations exist:
- `app/middlewares/limitter_middleware.go` — global tollbooth-based limiter (`LimitByRequest`)
- `app/modules/otp/middleware.go` — per-IP token-bucket limiter specifically for OTP endpoints (`LimitOTP`); relaxed to 1 req/sec in local env

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
Custom validators are registered in `app/server.go → RegisterValidators()`.  
`app/validations/` holds the validator functions (e.g. `irmobile` for Iranian mobile numbers) and `GetValidationErrors()` which maps `validator.ValidationErrors` to a translated `[]ValidationError` slice used by the response builder.

### API documentation
- **Swagger** — annotations live on handler functions; run `swag init -g cmd/main.go` after changing them to regenerate `docs/`.
- **Bruno collection** — `docs/apiDocs/` contains a Bruno API collection with folders for Health, OTP, and Users. Use Bruno to manually test endpoints against the local environment (`docs/apiDocs/environments/`).
