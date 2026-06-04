package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/zhitoo/golang-web-api/config"

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

func SetValue[T any](key string, value T, expiration time.Duration) error {
	ctx := context.Background()
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return redisClient.Set(ctx, key, jsonValue, expiration).Err()
}

func GetValue[T any](key string) (T, error) {
	ctx := context.Background()
	value, err := redisClient.Get(ctx, key).Result()
	var result T
	if err != nil {
		return result, err
	}
	err = json.Unmarshal([]byte(value), &result)
	return result, err
}

func DeleteValue(key string) error {
	ctx := context.Background()
	return redisClient.Del(ctx, key).Err()
}
