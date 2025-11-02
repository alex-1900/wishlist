package action

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/alex-1900/wishlist/src/app"
	"github.com/alex-1900/wishlist/src/auth"
	"github.com/alex-1900/wishlist/src/model"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// ActionPing returns a simple health check response
func ActionPing() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	}
}

// ActionCreateUser creates a new user for testing purposes
func ActionCreateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req model.UserCreateRequest

		// Bind JSON request to struct
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request format",
				"details": err.Error(),
			})
			return
		}

		// Validate the request
		if err := req.Validate(); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": err.Error(),
			})
			return
		}

		// Get repository
		userRepo := app.GetRepository().User()

		// Check if username already exists
		if exists, err := userRepo.ExistsByUsername(req.Username); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to check username availability",
				"details": err.Error(),
			})
			return
		} else if exists {
			ctx.JSON(http.StatusConflict, gin.H{
				"error": "Username already exists",
			})
			return
		}

		// Check if email already exists
		if exists, err := userRepo.ExistsByEmail(req.Email); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to check email availability",
				"details": err.Error(),
			})
			return
		} else if exists {
			ctx.JSON(http.StatusConflict, gin.H{
				"error": "Email already exists",
			})
			return
		}

		// Hash the password using auth package
		passwordHash, err := auth.HashPassword(req.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to hash password",
				"details": err.Error(),
			})
			return
		}

		// Create user model with gender
		user := &model.User{
			Username:     req.Username,
			Email:        req.Email,
			Gender:       model.ParseGender(req.Gender),
			PasswordHash: passwordHash,
		}
		user.BeforeCreate()

		// Save user to database
		if err := userRepo.Create(user); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to create user",
				"details": err.Error(),
			})
			return
		}

		// Return user response (without password hash)
		ctx.JSON(http.StatusCreated, gin.H{
			"message": "User created successfully",
			"user":    user.ToResponse(),
		})
	}
}

// ActionCreateTestUser creates a test user with random data for testing purposes
func ActionCreateTestUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Generate random test data
		username := generateRandomString(8, "testuser")
		email := fmt.Sprintf("%s@test.com", generateRandomString(6, "test"))
		password := "TestPassword123!"

		// Create user request
		req := model.UserCreateRequest{
			Username: username,
			Email:    email,
			Password: password,
		}

		// Get repository
		userRepo := app.GetRepository().User()

		// Hash the password
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to hash password",
				"details": err.Error(),
			})
			return
		}

		// Create user model
		user := &model.User{
			Username:     req.Username,
			Email:        req.Email,
			PasswordHash: string(passwordHash),
		}
		user.BeforeCreate()

		// Save user to database
		if err := userRepo.Create(user); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to create test user",
				"details": err.Error(),
			})
			return
		}

		// Return success response with created user info (excluding password)
		ctx.JSON(http.StatusCreated, gin.H{
			"message": "Test user created successfully",
			"user":    user.ToResponse(),
			"test_credentials": gin.H{
				"username": username,
				"password": password,
			},
		})
	}
}

// ActionListUsers returns a list of all users (for testing)
func ActionListUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRepo := app.GetRepository().User()

		users, err := userRepo.List()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to retrieve users",
				"details": err.Error(),
			})
			return
		}

		// Convert to response format (without password hashes)
		userResponses := make([]*model.UserResponse, len(users))
		for i, user := range users {
			userResponses[i] = user.ToResponse()
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Retrieved %d users", len(users)),
			"users":   userResponses,
		})
	}
}

// generateRandomString generates a random string with optional prefix
func generateRandomString(length int, prefix string) string {
	if prefix != "" {
		length = length - len(prefix)
		if length <= 0 {
			return prefix
		}
	}

	bytes := make([]byte, length)
	rand.Read(bytes)
	randomPart := base64.URLEncoding.EncodeToString(bytes)[:length]

	if prefix != "" {
		return prefix + randomPart
	}
	return randomPart
}

// ActionGetProfile retrieves the authenticated user's profile
func ActionGetProfile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get user ID from context (set by auth middleware)
		userID, exists := auth.GetUserID(ctx)
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not authenticated",
			})
			return
		}

		// Get repository
		userRepo := app.GetRepository().User()

		// Find user by ID
		user, err := userRepo.GetByID(userID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "User not found",
				"details": err.Error(),
			})
			return
		}

		// Return user profile
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Profile retrieved successfully",
			"data":    user.ToResponse(),
		})
	}
}

// ActionUpdateProfile updates the authenticated user's profile
func ActionUpdateProfile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get user ID from context (set by auth middleware)
		userID, exists := auth.GetUserID(ctx)
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not authenticated",
			})
			return
		}

		var req model.UserUpdateRequest

		// Bind JSON request to struct
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request format",
				"details": err.Error(),
			})
			return
		}

		// Validate the request
		if err := req.Validate(); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": err.Error(),
			})
			return
		}

		// Get repository
		userRepo := app.GetRepository().User()

		// Get existing user
		user, err := userRepo.GetByID(userID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "User not found",
				"details": err.Error(),
			})
			return
		}

		// Check if new username already exists (if being updated)
		if req.Username != nil && *req.Username != user.Username {
			if exists, err := userRepo.ExistsByUsername(*req.Username); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to check username availability",
					"details": err.Error(),
				})
				return
			} else if exists {
				ctx.JSON(http.StatusConflict, gin.H{
					"error": "Username already exists",
				})
				return
			}
			user.Username = *req.Username
		}

		// Check if new email already exists (if being updated)
		if req.Email != nil && *req.Email != user.Email {
			if exists, err := userRepo.ExistsByEmail(*req.Email); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to check email availability",
					"details": err.Error(),
				})
				return
			} else if exists {
				ctx.JSON(http.StatusConflict, gin.H{
					"error": "Email already exists",
				})
				return
			}
			user.Email = *req.Email
		}

		// Update gender (if provided)
		if req.Gender != nil {
			user.Gender = model.ParseGender(*req.Gender)
		}

		// Update password (if provided)
		if req.Password != nil {
			passwordHash, err := auth.HashPassword(*req.Password)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to hash password",
					"details": err.Error(),
				})
				return
			}
			user.PasswordHash = passwordHash
		}

		// Update timestamp
		user.BeforeUpdate()

		// Save user to database
		if err := userRepo.Update(user); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to update user",
				"details": err.Error(),
			})
			return
		}

		// Return updated user profile
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Profile updated successfully",
			"data":    user.ToResponse(),
		})
	}
}

// EmailVerificationRequest represents the request structure for email verification
type EmailVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ActionSendVerificationCode sends a verification code to the user's email (placeholder)
func ActionSendVerificationCode() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req EmailVerificationRequest

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

		// Check if user exists with this email
		_, err := userRepo.GetByEmail(req.Email)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "User with this email not found",
			})
			return
		}

		// In a real implementation, you would:
		// 1. Generate a verification code
		// 2. Store it with expiration
		// 3. Send it via email service
		// 4. Add rate limiting

		// For now, return a success response with placeholder data
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Verification code sent successfully",
			"data": gin.H{
				"email":              req.Email,
				"code":               "123456", // Placeholder code for testing
				"expires_in_minutes": 10,
			},
		})
	}
}

// EmailVerificationConfirmRequest represents the request structure for confirming email verification
type EmailVerificationConfirmRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
}

// ActionConfirmVerificationCode verifies the email verification code (placeholder)
func ActionConfirmVerificationCode() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req EmailVerificationConfirmRequest

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

		// Check if user exists with this email
		_, err := userRepo.GetByEmail(req.Email)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "User with this email not found",
			})
			return
		}

		// In a real implementation, you would:
		// 1. Verify the code against stored value
		// 2. Check if code has expired
		// 3. Mark email as verified in database

		// For now, accept the placeholder code "123456"
		if req.Code != "123456" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid verification code",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Email verified successfully",
			"data": gin.H{
				"email":    req.Email,
				"verified": true,
			},
		})
	}
}
