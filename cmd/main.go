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

	startupLog := logger.With(logging.General, logging.Startup)
	startupLog.Info("app starting", nil)

	err := cache.InitRedis(cfg)
	defer cache.CloseRedis()
	if err != nil {
		logger.With(logging.Redis, logging.Startup).Fatal(err.Error(), nil)
	}

	err = db.InitDb(cfg)
	defer db.CloseDb()
	if err != nil {
		logger.With(logging.Postgres, logging.Startup).Fatal(err.Error(), nil)
	}

	app.InitServer(cfg)
}
