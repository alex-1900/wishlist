package action

import (
	"net/http"
	"time"

	"github.com/alex-1900/wishlist/src/app"
	"github.com/alex-1900/wishlist/src/auth"
	"github.com/alex-1900/wishlist/src/model"
	"github.com/gin-gonic/gin"
)

// UserLoginRequest represents the request structure for user login
type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserLoginResponse represents the response structure for user login
type UserLoginResponse struct {
	Token     string              `json:"token"`
	User      *model.UserResponse `json:"user"`
	ExpiresIn int64               `json:"expires_in"`
	TokenType string              `json:"token_type"`
}

// ActionLogin handles user authentication
func ActionLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req UserLoginRequest

		// Bind JSON request to struct
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request format",
				"details": err.Error(),
			})
			return
		}

		// Get repository
		userRepo := app.GetRepository().User()

		// Find user by email
		user, err := userRepo.GetByEmail(req.Email)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid email or password",
			})
			return
		}

		// Check password
		if err := auth.CheckPassword(req.Password, user.PasswordHash); err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid email or password",
			})
			return
		}

		// Get JWT manager from app config
		config := app.GetConfig()
		jwtManager := auth.NewJWTManager(config.JWTSecret, time.Duration(config.JWTExpiration)*time.Hour)

		// Generate JWT token
		token, err := jwtManager.GenerateToken(user.ID, user.Username, user.Email)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to generate authentication token",
				"details": err.Error(),
			})
			return
		}

		// Return login response
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
			"data": UserLoginResponse{
				Token:     token,
				User:      user.ToResponse(),
				ExpiresIn: int64(config.JWTExpiration * 3600), // Convert hours to seconds
				TokenType: "Bearer",
			},
		})
	}
}

// ActionLogout handles user logout (placeholder for token blacklisting)
func ActionLogout() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// In a real implementation, you would:
		// 1. Add the token to a blacklist/revocation list
		// 2. Or use a refresh token pattern
		// 3. Or implement token rotation

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Logout successful",
		})
	}
}

// ActionRefreshToken handles token refresh
func ActionRefreshToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get user from context (should be authenticated)
		userID, exists := auth.GetUserID(ctx)
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not authenticated",
			})
			return
		}

		username, exists := auth.GetUsername(ctx)
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not authenticated",
			})
			return
		}

		email, exists := auth.GetEmail(ctx)
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not authenticated",
			})
			return
		}

		// Get JWT manager
		config := app.GetConfig()
		jwtManager := auth.NewJWTManager(config.JWTSecret, time.Duration(config.JWTExpiration)*time.Hour)

		// Generate new token
		token, err := jwtManager.GenerateToken(userID, username, email)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to refresh authentication token",
				"details": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Token refreshed successfully",
			"data": gin.H{
				"token":      token,
				"expires_in": int64(config.JWTExpiration * 3600), // Convert hours to seconds
				"token_type": "Bearer",
			},
		})
	}
}
