package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/money-advice/receipt-backend/internal/services"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// GoogleAuthRequest represents the request body for Google auth
type GoogleAuthRequest struct {
	AccessToken string `json:"access_token" binding:"required"`
}

// AuthResponse represents the response for successful authentication
type AuthResponse struct {
	Token string `json:"token"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// GoogleAuth handles Google OAuth authentication
func (h *AuthHandler) GoogleAuth(c *gin.Context) {
	var req GoogleAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	token, err := h.authService.ProcessGoogleAuth(req.AccessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Authentication failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token: token,
	})
}

// ValidateToken validates a JWT token
func (h *AuthHandler) ValidateToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Authorization header required",
		})
		return
	}

	// Remove "Bearer " prefix
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Invalid authorization header format",
		})
		return
	}

	token := authHeader[7:]
	userID, err := h.authService.ValidateJWT(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Invalid token: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"valid":   true,
	})
}