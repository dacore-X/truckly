package v1

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/dacore-x/truckly/internal/dto"
	"github.com/dacore-x/truckly/internal/transport/http/v1/middleware"
)

type UserUseCase interface {
	Create(context.Context, dto.UserRequestSignUpBody) error
	GetMe(context.Context, int64) (*dto.UserResponseMeBody, error)
	GetByID(context.Context, int64) (*dto.UserResponseInfoBody, error)
	GetByEmail(context.Context, string) (*dto.UserResponseInfoBody, error)
}

type userHandlers struct {
	UserUseCase
}

func newUserHandlers(superGroup *gin.RouterGroup, u UserUseCase, m middleware.Middlewares) {
	handler := &userHandlers{u}

	userGroup := superGroup.Group("/user")
	{
		userGroup.POST("/signup", handler.signUp)
		userGroup.POST("/login", handler.login)
	}
}

func (h *userHandlers) signUp(c *gin.Context) {
	// Get params from req body
	var body dto.UserRequestSignUpBody
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})
		return
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to hash password",
		})
		return
	}
	body.Password = string(hash)

	// Create user
	err = h.Create(context.Background(), body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "account is successfully created",
	})
}

func (h *userHandlers) login(c *gin.Context) {
	// Get params from req body
	var body dto.UserRequestLoginBody
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})
		return
	}

	// Look up requested user in DB
	user, err := h.GetByEmail(context.Background(), body.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid email or password",
		})
		return
	}

	// Compare sent in password wwith saved user password hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid email or password",
		})
		return
	}

	// Generate a jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create token",
		})
		return
	}

	// Send it back
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"msg": "authorization is complete",
	})
}
