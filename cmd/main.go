package main

import (
	"carsale/app"
	"carsale/config"
	"carsale/database/cache"
	"carsale/database/db"
	"carsale/pkg/logging"
)

// @securityDefinitaions.apiKey AuthBearer
// @in header
// @name Authorization
func main() {
	cfg := config.GetConfig()
	logger := logging.NewLogger(cfg)

	logger.Info(logging.General, logging.Startup, "app starting", nil)

	err := cache.InitRedis(cfg)
	defer cache.CloseRedis()
	if err != nil {
		logger.Fatal(logging.Redis, logging.Startup, err.Error(), nil)
	}

	err = db.InitDb(cfg)
	defer db.CloseDb()
	if err != nil {
		logger.Fatal(logging.Postgres, logging.Startup, err.Error(), nil)
	}

	app.InitServer(cfg)
}
