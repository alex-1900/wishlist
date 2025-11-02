package model

// This file serves as the main entry point for the model package
// It exports all domain models and related types

// Domain Models
// User - defined in user.go

// Request/Response Types
// UserCreateRequest - defined in user.go
// UserUpdateRequest - defined in user.go
// UserResponse - defined in user.go

// Repository Interfaces
// UserRepository - defined in user.go

// Validation Constants and Functions
// All validation logic is defined in user.go

// Export commonly used types for easier importing
type (
	UserModel         = User
	UserCreateReq     = UserCreateRequest
	UserUpdateReq     = UserUpdateRequest
	UserResp          = UserResponse
	UserRepoInterface = UserRepository
)
