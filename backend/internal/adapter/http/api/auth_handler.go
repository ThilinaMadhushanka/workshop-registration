package api

import (
	"errors"
	"net/http"

	"workshop-registration/backend/internal/adapter/http/dto"
	"workshop-registration/backend/internal/core/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Service *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{Service: s}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Valid email and password are required",
		})
		return
	}

	token, err := h.Service.Login(req.Email, req.Password)

	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid email or password",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Login failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":     token,
		"tokenType": "Bearer",
	})
}

func (h *AuthHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := h.Service.CreateUser(
		req.Name,
		req.Email,
		req.Password,
		req.Role,
	)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmailExists):
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
		case errors.Is(err, service.ErrInvalidUserData):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create user",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	})
}
