package main

import (
	"github.com/zhitoo/golang-web-api/app"
	"github.com/zhitoo/golang-web-api/config"
	"github.com/zhitoo/golang-web-api/database/cache"
	"github.com/zhitoo/golang-web-api/database/db"
	"github.com/zhitoo/golang-web-api/pkg/logging"
)

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
