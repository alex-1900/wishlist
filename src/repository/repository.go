package repository

import (
	"database/sql"

	"github.com/alex-1900/wishlist/src/model"
)

// RepositoryManager manages all repository instances
type RepositoryManager struct {
	UserRepo model.UserRepository
}

// NewRepositoryManager creates a new repository manager with all repositories
func NewRepositoryManager(db *sql.DB) *RepositoryManager {
	return &RepositoryManager{
		UserRepo: NewUserRepository(db),
	}
}

// Repository interface for easier testing and dependency injection
type Repository interface {
	User() model.UserRepository
}

// Ensure RepositoryManager implements the Repository interface
var _ Repository = (*RepositoryManager)(nil)

// User returns the user repository
func (rm *RepositoryManager) User() model.UserRepository {
	return rm.UserRepo
}
