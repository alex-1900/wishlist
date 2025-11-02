package account

import (
	"time"

	"github.com/alex-1900/wishlist/src/app"
	"github.com/alex-1900/wishlist/src/auth"
	"github.com/alex-1900/wishlist/src/module/account/action"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all account-related routes following the new routing principles
func RegisterRoutes(router *gin.Engine) {
	// Health check endpoints (keep them for now)
	router.GET("/ping", action.ActionPing())
	router.GET("/db-test", action.ActionDBTest())

	// User registration endpoint
	router.POST("/user-register", action.ActionCreateUser())

	// Email verification endpoints (placeholder implementation)
	router.POST("/send-verification-code", action.ActionSendVerificationCode())
	router.POST("/confirm-verification-code", action.ActionConfirmVerificationCode())

	// Authentication endpoint - email and password login
	router.POST("/user-login", action.ActionLogin())

	// Create auth middleware for protected routes
	config := app.GetConfig()
	jwtManager := auth.NewJWTManager(config.JWTSecret, time.Duration(config.JWTExpiration)*time.Hour)
	authMiddleware := auth.AuthMiddleware(jwtManager)

	// Protected routes (require authentication)
	protected := router.Group("/")
	protected.Use(authMiddleware)
	{
		// Profile management - get user profile
		protected.GET("/user-profile", action.ActionGetProfile())

		// Profile management - update user profile (username, email, gender, password)
		protected.POST("/update-user-profile", action.ActionUpdateProfile())

		// Authentication management
		protected.POST("/user-logout", action.ActionLogout())
		protected.POST("/refresh-auth-token", action.ActionRefreshToken())
	}

	// Testing endpoints (keep for development)
	router.POST("/create-test-user", action.ActionCreateTestUser())
	router.GET("/list-users", action.ActionListUsers())
}