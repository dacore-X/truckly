package redishelper

import (
	"fmt"

	"github.com/dacore-x/truckly/config"
	"github.com/redis/go-redis/v9"
)

// getAddr returns the Redis host:port address
func getAddr(cfg *config.REDIS) string {
	return fmt.Sprintf("%v:%v", cfg.RedisHost, cfg.RedisPort)
}

// GetOptions returns the Redis configuration options
func GetOptions(cfg *config.REDIS) *redis.Options {
	return &redis.Options{
		Addr:     getAddr(cfg),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	}
}
