package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// redisMiddlewares is a non-exportable struct
// that provides redis-related middlewares
type redisMiddlewares struct {
	redisClient *redis.Client
}

// RateLimit middleware checks if user's request rate
// is exceeded the number of requests limit
func (m *redisMiddlewares) RateLimit(c *gin.Context) {
	// Limit the number of requests to 5 requests per 10 minutes
	var (
		maxLimit  int           = 5
		limitTime time.Duration = 10 * time.Minute
	)

	// Get user id
	// Convert it to string since it is used as redis key
	userKey := fmt.Sprintf("%v", c.GetInt("user"))

	// Set key if not exists
	m.redisClient.SetNX(context.Background(), userKey, 1, limitTime).Val()

	counter, err := m.redisClient.Get(context.Background(), userKey).Int()
	if err != nil {
		err := fmt.Errorf("bad value received by key %v in redis", userKey)
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
	}

	// Check if counter exceeded the limit of requests
	if counter > maxLimit {
		err := fmt.Errorf("requests limit reached")
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
	}

	// Increment counter by one
	m.redisClient.Incr(context.Background(), userKey)

	// Continue
	c.Next()
}
