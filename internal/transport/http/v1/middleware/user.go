package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/net/context"

	"github.com/dacore-x/truckly/internal/usecase"
)

// userMiddlewares is a non-exportable struct
// that provides user-related middlewares
type userMiddlewares struct {
	usecase.User
}

// RequireAuth middleware checks if user is authenticated
// by decoding and validating user's jwt token and attaches
// private user's data to the request
func (m *userMiddlewares) RequireAuth(c *gin.Context) {
	// Get cookie from request
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "user must be authorized",
		})
		return
	}

	// Decode token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	// Check errors from decoding and validate
	if errors.Is(err, jwt.ErrTokenMalformed) ||
		errors.Is(err, jwt.ErrTokenExpired) ||
		errors.Is(err, jwt.ErrTokenNotValidYet) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "user must be authorized",
		})
		return
	} else {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Check the exp
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "user must be authorized",
				})
				return
			}

			// Attach to request
			sub := claims["sub"].(float64)
			c.Set("user", int(sub))
		}
	}

	// Continue
	c.Next()
}

// RequireNoBan middleware checks if user is not banned
func (m *userMiddlewares) RequireNoBan(c *gin.Context) {
	// Check user authorization
	userKey := c.GetInt("user")

	// Check ban status
	resp, err := m.GetUserMeta(context.Background(), userKey)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "user is banned",
		})
	}

	// Attach to request
	c.Set("ban status", resp.IsBanned)

	// continue
	c.Next()
}

// RequireAdmin middleware checks if user has admin privileges
func (m *userMiddlewares) RequireAdmin(c *gin.Context) {
	// Get user from keys
	userKey := c.GetInt("user")

	// Check for admin privileges
	resp, err := m.GetUserMeta(context.Background(), userKey)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "user is not admin",
		})
	}

	// Attach to request
	c.Set("admin", resp.IsAdmin)

	// continue
	c.Next()
}

// RequireCourier middleware checks if user's role is courier
func (m *userMiddlewares) RequireCourier(c *gin.Context) {
	// Get user from keys
	userKey := c.GetInt("user")

	// Check for courier role
	resp, err := m.GetUserMeta(context.Background(), userKey)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "user is not courier",
		})
	}

	// Attach to request
	c.Set("courier", resp.IsCourier)

	// continue
	c.Next()
}
