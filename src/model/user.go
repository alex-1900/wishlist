package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Gender represents the gender type for users
type Gender string

// Gender constants
const (
	GenderUnknown Gender = "unknown"
	GenderMale    Gender = "male"
	GenderFemale  Gender = "female"
)

// ParseGender parses a string into a Gender type
func ParseGender(gender string) Gender {
	switch strings.ToLower(gender) {
	case "male":
		return GenderMale
	case "female":
		return GenderFemale
	default:
		return GenderUnknown
	}
}

// IsValid checks if the gender value is valid
func (g Gender) IsValid() bool {
	return g == GenderMale || g == GenderFemale || g == GenderUnknown
}

// String returns the string representation of Gender
func (g Gender) String() string {
	return string(g)
}

// User represents the user domain model
type User struct {
	ID           int       `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	Gender       Gender    `json:"gender" db:"gender"`
	PasswordHash string    `json:"-" db:"password_hash"` // Hidden from JSON output
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(user *User) error
	GetByID(id int) (*User, error)
	GetByUsername(username string) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id int) error
	List() ([]*User, error)
	ExistsByUsername(username string) (bool, error)
	ExistsByEmail(email string) (bool, error)
	GetTotalCount() (int, error)
	UpdatePassword(userID int, passwordHash string) error
}

// UserCreateRequest represents the request structure for creating a user
type UserCreateRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Gender   string `json:"gender" binding:"omitempty,oneof=male female unknown"`
	Password string `json:"password" binding:"required,min=8"`
}

// UserUpdateRequest represents the request structure for updating a user
type UserUpdateRequest struct {
	Username *string `json:"username,omitempty" binding:"omitempty,min=3,max=50"`
	Email    *string `json:"email,omitempty" binding:"omitempty,email"`
	Gender   *string `json:"gender,omitempty" binding:"omitempty,oneof=male female unknown"`
	Password *string `json:"password,omitempty" binding:"omitempty,min=8"`
}

// UserResponse represents the safe response structure for user data
type UserResponse struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Gender    Gender    `json:"gender"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validation constants
const (
	UsernameMinLength = 3
	UsernameMaxLength = 50
	EmailMaxLength    = 100
	PasswordMinLength = 8
)

// Validation regex patterns
var (
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
)

// Validate validates the UserCreateRequest fields
func (ucr *UserCreateRequest) Validate() error {
	if err := validateUsername(ucr.Username); err != nil {
		return fmt.Errorf("username validation failed: %w", err)
	}

	if err := validateEmail(ucr.Email); err != nil {
		return fmt.Errorf("email validation failed: %w", err)
	}

	if err := validateGender(ucr.Gender); err != nil {
		return fmt.Errorf("gender validation failed: %w", err)
	}

	if err := validatePassword(ucr.Password); err != nil {
		return fmt.Errorf("password validation failed: %w", err)
	}

	return nil
}

// Validate validates the UserUpdateRequest fields
func (uur *UserUpdateRequest) Validate() error {
	if uur.Username != nil {
		if err := validateUsername(*uur.Username); err != nil {
			return fmt.Errorf("username validation failed: %w", err)
		}
	}

	if uur.Email != nil {
		if err := validateEmail(*uur.Email); err != nil {
			return fmt.Errorf("email validation failed: %w", err)
		}
	}

	if uur.Gender != nil {
		if err := validateGender(*uur.Gender); err != nil {
			return fmt.Errorf("gender validation failed: %w", err)
		}
	}

	if uur.Password != nil {
		if err := validatePassword(*uur.Password); err != nil {
			return fmt.Errorf("password validation failed: %w", err)
		}
	}

	return nil
}

// validateUsername validates the username field
func validateUsername(username string) error {
	if len(username) < UsernameMinLength {
		return errors.New("username is too short")
	}

	if len(username) > UsernameMaxLength {
		return errors.New("username is too long")
	}

	if !usernameRegex.MatchString(username) {
		return errors.New("username can only contain alphanumeric characters, underscores, and hyphens")
	}

	return nil
}

// validateEmail validates the email field
func validateEmail(email string) error {
	if len(email) > EmailMaxLength {
		return errors.New("email is too long")
	}

	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}

	return nil
}

// validateGender validates the gender field
func validateGender(gender string) error {
	if gender == "" {
		return nil // Optional field, default to unknown
	}

	parsedGender := ParseGender(gender)
	if !parsedGender.IsValid() {
		return errors.New("gender must be one of: male, female, unknown")
	}

	return nil
}

// validatePassword validates the password field
func validatePassword(password string) error {
	if len(password) < PasswordMinLength {
		return errors.New("password is too short")
	}

	// You can add more password complexity requirements here
	// For example: check for uppercase, lowercase, numbers, special characters
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	if !hasUpper || !hasLower || !hasNumber {
		return errors.New("password must contain at least one uppercase letter, one lowercase letter, and one number")
	}

	return nil
}

// ToResponse converts a User to a UserResponse (safe for API responses)
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Gender:    u.Gender,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// BeforeCreate sets the CreatedAt and UpdatedAt fields before creating a new user
func (u *User) BeforeCreate() {
	now := time.Now().UTC()
	u.CreatedAt = now
	u.UpdatedAt = now
}

// BeforeUpdate updates the UpdatedAt field before updating an existing user
func (u *User) BeforeUpdate() {
	u.UpdatedAt = time.Now().UTC()
}

// Value implements the driver.Valuer interface for database operations
func (u User) Value() (driver.Value, error) {
	return json.Marshal(u)
}

// Scan implements the sql.Scanner interface for database operations
func (u *User) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("cannot scan non-byte value into User")
	}

	return json.Unmarshal(bytes, u)
}
