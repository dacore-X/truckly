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
		c.AbortWithStatus(http.StatusUnauthorized)
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
	} else {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Check the exp
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				c.AbortWithStatus(http.StatusUnauthorized)
			}

			// Find the user with token sub
			sub := claims["sub"].(float64)
			user, err := m.GetByID(context.Background(), int(sub))
			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
			}

			// Attach to request
			c.Set("user", user)
		}
	}

	// Continue
	c.Next()
}
