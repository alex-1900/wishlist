package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/alex-1900/wishlist/src/model"
)

// UserRepository implements the model.UserRepository interface
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *sql.DB) model.UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create creates a new user in the database
func (r *UserRepository) Create(user *model.User) error {
	query := `
		INSERT INTO users (username, email, gender, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var id int
	err := r.db.QueryRow(
		query,
		user.Username,
		user.Email,
		user.Gender,
		user.PasswordHash,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&id)

	if err != nil {
		log.Printf("Error creating user: %v", err)
		return fmt.Errorf("failed to create user: %w", err)
	}

	user.ID = id
	log.Printf("User created successfully with ID: %d", id)
	return nil
}

// GetByID retrieves a user by their ID
func (r *UserRepository) GetByID(id int) (*model.User, error) {
	query := `
		SELECT id, username, email, gender, password_hash, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &model.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Gender,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with ID %d not found", id)
		}
		log.Printf("Error getting user by ID %d: %v", id, err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetByUsername retrieves a user by their username
func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	query := `
		SELECT id, username, email, gender, password_hash, created_at, updated_at
		FROM users
		WHERE username = $1
	`

	user := &model.User{}
	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Gender,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with username '%s' not found", username)
		}
		log.Printf("Error getting user by username '%s': %v", username, err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetByEmail retrieves a user by their email
func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	query := `
		SELECT id, username, email, gender, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &model.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Gender,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with email '%s' not found", email)
		}
		log.Printf("Error getting user by email '%s': %v", email, err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// Update updates an existing user in the database
func (r *UserRepository) Update(user *model.User) error {
	query := `
		UPDATE users
		SET username = $2, email = $3, gender = $4, password_hash = $5, updated_at = $6
		WHERE id = $1
	`

	user.BeforeUpdate() // Update the timestamp
	result, err := r.db.Exec(
		query,
		user.ID,
		user.Username,
		user.Email,
		user.Gender,
		user.PasswordHash,
		user.UpdatedAt,
	)

	if err != nil {
		log.Printf("Error updating user with ID %d: %v", user.ID, err)
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected for user update: %v", err)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", user.ID)
	}

	log.Printf("User with ID %d updated successfully", user.ID)
	return nil
}

// Delete deletes a user from the database
func (r *UserRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting user with ID %d: %v", id, err)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected for user deletion: %v", err)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", id)
	}

	log.Printf("User with ID %d deleted successfully", id)
	return nil
}

// List retrieves all users from the database
func (r *UserRepository) List() ([]*model.User, error) {
	query := `
		SELECT id, username, email, gender, password_hash, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("Error listing users: %v", err)
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf("Error closing rows: %v", closeErr)
		}
	}()

	var users []*model.User
	for rows.Next() {
		user := &model.User{}
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Gender,
			&user.PasswordHash,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning user row: %v", err)
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating over user rows: %v", err)
		return nil, fmt.Errorf("error iterating over users: %w", err)
	}

	log.Printf("Retrieved %d users from database", len(users))
	return users, nil
}

// Helper methods for common operations

// ExistsByUsername checks if a user with the given username exists
func (r *UserRepository) ExistsByUsername(username string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE username = $1`

	var count int
	err := r.db.QueryRow(query, username).Scan(&count)
	if err != nil {
		log.Printf("Error checking if username exists: %v", err)
		return false, fmt.Errorf("failed to check username existence: %w", err)
	}

	return count > 0, nil
}

// ExistsByEmail checks if a user with the given email exists
func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE email = $1`

	var count int
	err := r.db.QueryRow(query, email).Scan(&count)
	if err != nil {
		log.Printf("Error checking if email exists: %v", err)
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return count > 0, nil
}

// GetTotalCount returns the total number of users in the database
func (r *UserRepository) GetTotalCount() (int, error) {
	query := `SELECT COUNT(*) FROM users`

	var count int
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		log.Printf("Error getting total user count: %v", err)
		return 0, fmt.Errorf("failed to get user count: %w", err)
	}

	return count, nil
}

// UpdatePassword updates only the password hash for a user
func (r *UserRepository) UpdatePassword(userID int, passwordHash string) error {
	query := `
		UPDATE users
		SET password_hash = $2, updated_at = $3
		WHERE id = $1
	`

	result, err := r.db.Exec(query, userID, passwordHash, time.Now().UTC())
	if err != nil {
		log.Printf("Error updating password for user ID %d: %v", userID, err)
		return fmt.Errorf("failed to update password: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected for password update: %v", err)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", userID)
	}

	log.Printf("Password updated successfully for user ID %d", userID)
	return nil
}
