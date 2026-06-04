package db

import (
	"carsale/config"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbClient *gorm.DB

func InitDb(cfg *config.Config) error {
	var err error
	con := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DB,
		cfg.Postgres.SslMode,
		cfg.Postgres.TimeZone,
	)
	dbClient, err = gorm.Open(postgres.Open(con), &gorm.Config{})

	if err != nil {
		return err
	}

	sqlDb, _ := dbClient.DB()

	err = sqlDb.Ping()

	if err != nil {
		return err
	}

	sqlDb.SetMaxIdleConns(10)                 //todo: add to config
	sqlDb.SetMaxOpenConns(150)                //todo: add to config
	sqlDb.SetConnMaxLifetime(5 * time.Minute) //todo: add to config

	log.Println("db connection successfull")

	return nil

}

func GetDb() *gorm.DB {
	return dbClient
}

func CloseDb() {
	connection, _ := dbClient.DB()

	connection.Close()
}
