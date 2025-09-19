package controllers

import (
	"chatapp/models"
	"chatapp/service"
	"chatapp/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService service.AuthService
}

// NewAuthController creates a new auth controller
func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

// Login handles user login
func (ctrl *AuthController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	user, token, err := ctrl.authService.Login(req.Username, req.Password)
	if err != nil {
		utils.UnauthorizedResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, LoginResponse{
		Token: token,
		User:  *user,
	})
}

// GetProfile returns user profile
func (ctrl *AuthController) GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	user, err := ctrl.authService.GetUserProfile(userID.(uint))
	if err != nil {
		utils.NotFoundResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, user)
}

// Logout handles user logout
func (ctrl *AuthController) Logout(c *gin.Context) {
	// In a stateless JWT system, logout is typically handled client-side
	// by removing the token from storage. However, we can provide a
	// server-side endpoint for consistency and future token blacklisting.

	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	// Optional: Add any logout-specific business logic here
	// For example: logging the logout event, clearing user sessions, etc.

	utils.SuccessResponseWithMessage(c, "Logout successful", gin.H{
		"user_id": userID,
	})
}
