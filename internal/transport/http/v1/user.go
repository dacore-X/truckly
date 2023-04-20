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
		userGroup.POST("/logout", m.RequireAuth, handler.logout)
		userGroup.POST("/:id/ban", m.RequireAuth, m.RequireNoBan, m.RequireAdmin, handler.ban)
		userGroup.POST("/:id/unban", m.RequireAuth, m.RequireNoBan, m.RequireAdmin, handler.unban)
	}
}

// me handler gets user's account data based on
// private user's data from "user" context key
func (h *userHandlers) me(c *gin.Context) {
	// Check user authorization
	userKey := c.GetInt("user")

	// Look up user in DB
	user, err := h.GetUserByID(context.Background(), userKey)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Fill response body with user's account data
	resp := dto.UserMeResponse{
		ID:          user.ID,
		Surname:     user.Surname,
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		CreatedAt:   user.CreatedAt,
		Meta:        user.Meta,
	}
	c.JSON(http.StatusOK, resp)
}

// signUp handler creates new user account
// based on request body data with password hashing
func (h *userHandlers) signUp(c *gin.Context) {
	// Get params from req body
	var body dto.UserSignUpRequestBody
	if c.BindJSON(&body) != nil {
		err := fmt.Errorf("failed to read body")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Check if the user with specified email from req body already exists in database
	record, _ := h.GetUserPrivateByEmail(context.Background(), body.Email)
	if record != nil {
		err := fmt.Errorf("user with this email already exists")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		err := fmt.Errorf("failed to hash password")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	body.Password = string(hash)

	// Create user
	err = h.CreateUser(context.Background(), &body)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
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
	if c.BindJSON(&body) != nil {
		err := fmt.Errorf("failed to read body")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Look up requested user in DB
	user, err := h.GetUserPrivateByEmail(context.Background(), body.Email)
	if err != nil {
		err := fmt.Errorf("invalid email or password")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Compare sent in password wwith saved user password hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		err := fmt.Errorf("invalid email or password")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
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
		err := fmt.Errorf("failed to create token")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Send it back
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 60*60*24*7, "", "", false, true) // cookie expires in a week

	c.JSON(http.StatusOK, gin.H{
		"msg": "authorization is complete",
	})
}

// ban handler gets user's id from URI and bans him
func (h *userHandlers) ban(c *gin.Context) {
	// Get params from request
	var req dto.UserBanParams
	if c.ShouldBindUri(&req) != nil {
		err := fmt.Errorf("failed to read uri")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Ban user
	err := h.BanUser(context.Background(), req.ID)
	if err != nil {
		err := fmt.Errorf("failed to ban user")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("user %v has been successfully banned", req.ID),
	})
}

// unban handler gets user's id from URI and unbans him
func (h *userHandlers) unban(c *gin.Context) {
	// Get params from request
	var req dto.UserBanParams
	if c.ShouldBindUri(&req) != nil {
		err := fmt.Errorf("failed to read uri")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Ban user
	err := h.UnbanUser(context.Background(), req.ID)
	if err != nil {
		err := fmt.Errorf("failed to unban user")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("user %v has been successfully unbanned", req.ID),
	})
}

func (h *userHandlers) logout(c *gin.Context) {
	// Send it back
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", "", -1, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"msg": "logout is complete",
	})
}
