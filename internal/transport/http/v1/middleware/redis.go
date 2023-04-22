package middleware

import "github.com/redis/go-redis/v9"

// redisMiddlewares is a non-exportable struct
// that provides redis-related middlewares
type redisMiddlewares struct {
	client *redis.Client
}
