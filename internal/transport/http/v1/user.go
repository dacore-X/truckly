package v1

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/dacore-x/truckly/internal/dto"
	"github.com/dacore-x/truckly/internal/transport/http/v1/middleware"
	"github.com/dacore-x/truckly/internal/usecase"
)

// userHandlers is a non-exportable struct
// that provides user-related handlers
type userHandlers struct {
	usecase.User
}

// newUserHandlers initializes a group of user's routes
func newUserHandlers(superGroup *gin.RouterGroup, u usecase.User, m *middleware.Middlewares) {
	handler := &userHandlers{u}

	userGroup := superGroup.Group("/user")
	{
		userGroup.GET("/me", m.RequireAuth, handler.me)
		userGroup.POST("/signup", handler.signUp)
		userGroup.POST("/login", handler.login)
		userGroup.POST("/:id/ban", m.RequireAuth, handler.ban)
		userGroup.POST("/:id/unban", m.RequireAuth, handler.unban)
	}

}

// me handler gets user's account data based on
// private user's data from "user" context key
func (h *userHandlers) me(c *gin.Context) {
	// Check user authorization
	authUserKey, _ := c.Get("user")
	if _, ok := authUserKey.(*dto.UserInfoResponse); !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user must be authorized",
		})
	}
	req := authUserKey.(*dto.UserInfoResponse)

	// Look up user in DB
	user, err := h.GetMe(context.Background(), req.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	// Fill response body with user's account data
	resp := dto.UserMeResponse{
		ID:          user.ID,
		Surname:     user.Surname,
		Name:        user.Name,
		Patronymic:  user.Patronymic,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		CreatedAt:   user.CreatedAt,
	}
	c.JSON(http.StatusOK, resp)
}

// signUp handler creates new user account
// based on request body data with password hashing
func (h *userHandlers) signUp(c *gin.Context) {
	// Get params from req body
	var body dto.UserSignUpRequestBody
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
	err = h.CreateTx(context.Background(), &body)
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

// login handler checks if user has an account
// based on request body data, creates new jwt token
// and stores it in cookie
func (h *userHandlers) login(c *gin.Context) {
	// Get params from req body
	var body dto.UserLoginRequestBody
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

// ban handler gets user's id from URI and bans him
func (h *userHandlers) ban(c *gin.Context) {
	// Check user authorization
	authUserKey, _ := c.Get("user")
	if _, ok := authUserKey.(*dto.UserInfoResponse); !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user must be authorized",
		})
	}

	// Get params from request
	var req dto.UserBanParams
	if c.ShouldBindUri(&req) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read uri",
		})
		return
	}

	// Ban user
	err := h.Ban(context.Background(), req.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to ban user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("user %v has been successfully banned", req.ID),
	})
}

// unban handler gets user's id from URI and unbans him
func (h *userHandlers) unban(c *gin.Context) {
	// Check user authorization
	authUserKey, _ := c.Get("user")
	if _, ok := authUserKey.(*dto.UserInfoResponse); !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user must be authorized",
		})
	}

	// Get params from request
	var req dto.UserBanParams
	if c.ShouldBindUri(&req) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read uri",
		})
		return
	}

	// Ban user
	err := h.Unban(context.Background(), req.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to unban user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("user %v has been successfully unbanned", req.ID),
	})
}
