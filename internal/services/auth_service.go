package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/money-advice/receipt-backend/internal/database"
	"github.com/money-advice/receipt-backend/internal/models"
)

// GoogleUserInfo represents user info from Google OAuth
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

// AuthService handles authentication operations
type AuthService struct{}

// NewAuthService creates a new auth service instance
func NewAuthService() *AuthService {
	return &AuthService{}
}

// ValidateGoogleToken validates Google OAuth token and returns user info
func (s *AuthService) ValidateGoogleToken(accessToken string) (*GoogleUserInfo, error) {
	url := fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s", accessToken)
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to validate Google token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid Google token, status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var userInfo GoogleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	return &userInfo, nil
}

// ProcessGoogleAuth processes Google OAuth authentication
func (s *AuthService) ProcessGoogleAuth(accessToken string) (string, error) {
	// Validate Google token and get user info
	userInfo, err := s.ValidateGoogleToken(accessToken)
	if err != nil {
		return "", err
	}

	// Convert Google ID to uint64
	googleID, err := strconv.ParseUint(userInfo.ID, 10, 64)
	if err != nil {
		return "", fmt.Errorf("invalid Google ID format: %w", err)
	}

	// Check if user exists
	var user models.User
	result := database.DB.Where("google_id = ?", googleID).First(&user)
	
	if result.Error != nil {
		// User doesn't exist, create new one
		user = models.User{
			GoogleID: &googleID,
		}
		
		if err := database.DB.Create(&user).Error; err != nil {
			return "", fmt.Errorf("failed to create user: %w", err)
		}
	}

	// Generate JWT token
	token, err := s.GenerateJWT(user.GetUUID())
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: %w", err)
	}

	return token, nil
}

// GenerateJWT generates a JWT token for the user
func (s *AuthService) GenerateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET not set")
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateJWT validates a JWT token and returns user ID
func (s *AuthService) ValidateJWT(tokenString string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET not set")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, ok := claims["user_id"].(string); ok {
			return userID, nil
		}
		return "", fmt.Errorf("user_id not found in token")
	}

	return "", fmt.Errorf("invalid token")
}