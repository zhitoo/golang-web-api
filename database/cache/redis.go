package cache

import (
	"carsale/config"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

func InitRedis(cfg *config.Config) error {
	redisClient = redis.NewClient(&redis.Options{
		Addr:               fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password:           cfg.Redis.Password,
		DB:                 cfg.Redis.DB,
		DialTimeout:        cfg.Redis.DialTimeout * time.Second,
		ReadTimeout:        cfg.Redis.ReadTimeout * time.Second,
		WriteTimeout:       cfg.Redis.WriteTimeout * time.Second,
		PoolSize:           cfg.Redis.PoolSize,
		PoolTimeout:        cfg.Redis.PoolTimeout,
		IdleTimeout:        cfg.Redis.IdleTimeout * time.Millisecond,
		IdleCheckFrequency: cfg.Redis.IdleCheckFrequency * time.Millisecond,
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func GetRedis() *redis.Client {
	return redisClient
}

func CloseRedis() {
	redisClient.Close()
}
