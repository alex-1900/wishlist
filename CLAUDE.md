# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a wishlist SNS (Social Networking Service) built with Go and the Gin web framework. The project follows clean architecture principles with clear separation of concerns across domain models, repositories, and HTTP layers. The application implements user authentication, profile management, and email verification with a modular structure.

## Architecture

The application follows a dependency injection pattern with a singleton manager and clean architecture principles:

### Core Layers

- **src/main.go**: Entry point that obtains dependencies through the dependency manager
- **src/app/**: Core application with dependency injection system
  - `types.go`: **Type definitions** for `App` struct and `AppConfig` (central type definitions)
  - `config.go`: Application configuration including database settings
  - `providers.go`: Dependency creation functions (`buildApp`, `buildGinEngine`, database connection)
  - `dependency.go`: **Central dependency manager** providing singleton access:
    - `GetInstance()`: Returns singleton App instance
    - `GetGinEngine()`: Direct access to Gin engine
    - `GetConfig()`: Direct access to configuration
    - `GetDB()`: Direct access to database connection
    - `GetRepository()`: Direct access to repository manager
    - `ResetApp()`: Reset singleton (for testing)
- **src/model/**: Domain models and business logic
  - `user.go`: User domain model with validation, request/response types
  - `types.go`: Package exports and type aliases
- **src/repository/**: Data access layer implementing repository pattern
  - `user_repository.go`: User repository with full CRUD operations
  - `repository.go`: Repository manager and interfaces
- **src/database/**: Database schema and migrations
  - `migrations.go`: Database table creation and connection verification
- **src/module/**: HTTP layer with modular routing
  - `routes.go`: Main route definition that delegates to modules
  - `account/`: Account module handling user authentication and profile management
    - `module.go`: Account module route registration
    - `action/`: Account-related handler functions (user, auth, db operations)

### Dependency Flow
1. `main.go` → `app.GetInstance()` → `buildApp()` (in providers.go)
2. `buildApp()` creates App instance with database connection and repositories
3. Repository manager is initialized with database connection
4. Dependency manager maintains singleton instance throughout application lifecycle
5. All type definitions are centralized in `src/app/types.go`

### Database Integration
- PostgreSQL database connection managed through dependency injection
- Users table with fields: id, username, email, gender, password_hash, created_at, updated_at
- Database migrations run automatically on application startup
- Repository pattern provides clean data access abstraction

### Authentication System
- JWT-based authentication managed through `src/auth/jwt.go`
- Auth middleware (`auth.AuthMiddleware`) protects routes requiring authentication
- Token management includes generation, validation, and refresh capabilities
- User context available in protected routes via `auth.GetUserID()`, `auth.GetUsername()`, `auth.GetEmail()`

## Common Commands

### Build and Run
```bash
# Build the application
go build -o wishlist ./src/main.go

# Run the application directly
go run ./src/main.go

# Run with timeout for testing
timeout 5 ./wishlist

# The application runs on :8080 by default
```

### Development
```bash
# Format code
go fmt ./...

# Vet code for potential issues
go vet ./...

# Get dependencies
go mod tidy

# Download dependencies
go mod download

# Build specific packages
go build ./src/model
go build ./src/repository
```

### Testing
```bash
# Run tests (when test files exist)
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific packages
go test ./src/model
go test ./src/repository
```

## Key Dependencies

- **github.com/gin-gonic/gin**: Web framework for HTTP routing and middleware
- **github.com/golang-jwt/jwt/v5**: JWT token implementation for authentication
- **github.com/lib/pq**: PostgreSQL driver for database connectivity
- **golang.org/x/crypto**: Cryptographic functions including bcrypt for password hashing

## Development Notes

### Architecture Principles
- **Dependency Injection**: All dependencies should be obtained through `app/dependency.go` functions
- **Singleton Pattern**: App instance is managed as a singleton to ensure consistent dependency access
- **Repository Pattern**: Data access is abstracted through repository interfaces
- **Domain Models**: Business logic and validation are encapsulated in domain models
- **Clean Architecture**: Clear separation between HTTP, domain, and data layers

### Database Usage
- Database connection is automatically established on application startup
- Schema migrations run automatically (users table creation)
- Use `app.GetRepository().User()` to access user repository operations
- Repository provides: Create, GetByID, GetByUsername, GetByEmail, Update, Delete, List, ExistsByUsername, ExistsByEmail, UpdatePassword operations

### Module Structure and Routing
- HTTP routes are organized by business domain in separate modules under `src/module/`
- Each module has its own `module.go` with `RegisterRoutes()` function
- Main `src/module/routes.go` delegates to individual modules
- Routing follows semantic naming with kebab-case (e.g., `/user-register`, `/update-user-profile`)
- Only GET and POST methods are used per project requirements

### Testing Guidelines
- Use `app.ResetApp()` to reset singleton state between tests
- Mock repositories can be injected for unit testing through the repository interface
- Database integration tests should use a separate test database

### Type Organization
- All core types (`App`, `AppConfig`) are centralized in `src/app/types.go`
- Domain models are in `src/model/` with validation and business logic
- Repository interfaces are defined in `src/model/` and implemented in `src/repository/`

### Code Style
- Package names follow Go conventions (lowercase, single words)
- Error handling uses Go's error wrapping with context
- Logging is included for important operations and errors
- Database operations use prepared statements and proper error handling

### User Model and Validation
- User model supports gender (male, female, unknown) with comprehensive validation
- Password validation requires minimum 8 characters with complexity requirements
- Email validation uses regex patterns and length restrictions
- Username validation allows alphanumeric characters with underscores/hyphens
- All validation logic is centralized in the domain model with detailed error messages

## Current API Endpoints

### Public Endpoints
- `GET /ping`: Health check endpoint returning `{"message": "pong"}`
- `GET /db-test`: Database connectivity test endpoint (returns connection status)
- `POST /user-register`: User registration with email, username, gender, and password
- `POST /user-login`: User authentication with email and password
- `POST /send-verification-code`: Send email verification code (placeholder implementation)
- `POST /confirm-verification-code`: Confirm email verification code (placeholder implementation)

### Protected Endpoints (require JWT authentication)
- `GET /user-profile`: Get authenticated user's profile information
- `POST /update-user-profile`: Update user profile (username, email, gender, password)
- `POST /user-logout`: User logout (placeholder for token blacklisting)
- `POST /refresh-auth-token`: Refresh JWT authentication token

### Testing Endpoints
- `POST /create-test-user`: Create test user with random credentials for development
- `GET /list-users`: List all users (for testing purposes)

### Repository Access Pattern
```go
// Get user repository
userRepo := app.GetRepository().User()

// Create a user
user := &model.User{
    Username: "testuser",
    Email: "test@example.com",
    Gender: model.GenderUnknown, // or GenderMale/GenderFemale
    PasswordHash: "hashed_password", // Use auth.HashPassword()
}
user.BeforeCreate() // Set timestamps
err := userRepo.Create(user)

// Find user by username/email
user, err := userRepo.GetByUsername("testuser")
user, err := userRepo.GetByEmail("test@example.com")

// Check existence
exists, err := userRepo.ExistsByUsername("testuser")
exists, err := userRepo.ExistsByEmail("test@example.com")
```

### Authentication Pattern
```go
// Login flow
config := app.GetConfig()
jwtManager := auth.NewJWTManager(config.JWTSecret, time.Duration(config.JWTExpiration)*time.Hour)

// Generate token for authenticated user
token, err := jwtManager.GenerateToken(user.ID, user.Username, user.Email)

// In protected routes, get user context
userID, exists := auth.GetUserID(ctx)
username, exists := auth.GetUsername(ctx)
email, exists := auth.GetEmail(ctx)
```