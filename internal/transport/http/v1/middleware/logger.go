package middleware

import (
	"fmt"
	"time"

	"github.com/dacore-x/truckly/pkg/logger"
	"github.com/gin-gonic/gin"
)

// loggerMiddlewares is a non-exportable struct
// that provides logger middlewares
type loggerMiddlewares struct {
	*logger.Logger
}

// Log middleware adds information to logs about transport level
func (m *loggerMiddlewares) Log(c *gin.Context) {
	start := time.Now() // Start timer
	path := c.Request.URL.Path
	raw := c.Request.URL.RawQuery
	queryParams := c.Request.URL.Query()

	// Process request
	c.Next()

	// Fill the params
	param := gin.LogFormatterParams{}

	param.TimeStamp = time.Now() // Stop timer
	param.Latency = param.TimeStamp.Sub(start)
	if param.Latency > time.Minute {
		param.Latency = param.Latency.Truncate(time.Second)
	}

	param.ClientIP = c.ClientIP()
	param.Method = c.Request.Method
	param.StatusCode = c.Writer.Status()
	param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()

	// Concat request path
	if raw != "" {
		path = path + "?"
	}

	queryNum := 0
	for key, values := range queryParams {
		if queryNum != 0 {
			path = path + "&"
		}
		path = path + key + "="
		for _, x := range values {
			path = path + x
		}
		queryNum++
	}
	param.Path = path

	// Turn on console color output for some params
	gin.ForceConsoleColor()
	statusColor := param.StatusCodeColor()
	methodColor := param.MethodColor()
	resetColor := param.ResetColor()

	// Logger message using the params
	msg := fmt.Sprintf("|%s %3d %s| %13v | %8s |%s %-7s %s %#v\n%s",
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		param.ClientIP,
		methodColor, param.Method, resetColor,
		param.Path,
		param.ErrorMessage,
	)

	// Log using the function that matches the status of the request
	if c.Writer.Status() >= 500 {
		m.Logger.Error(msg)
	} else {
		m.Logger.Info(msg)
	}
}

// DefaultLogger logs a gin HTTP request using Log middleware.
// Uses formatting configuration for logger.Logger.
func (m *loggerMiddlewares) DefaultLogger() gin.HandlerFunc {
	return m.Log
}
