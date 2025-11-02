package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// InitializeSchema creates the users table for the application
func InitializeSchema(db *sql.DB) error {
	// Create users table
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		gender VARCHAR(10) DEFAULT 'unknown' NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	)`

	if _, err := db.Exec(usersTable); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	// Add gender column to existing table if it doesn't exist (for backward compatibility)
	addGenderColumn := `
	DO $$
	BEGIN
		IF NOT EXISTS (
			SELECT 1 FROM information_schema.columns
			WHERE table_name = 'users' AND column_name = 'gender'
		) THEN
			ALTER TABLE users ADD COLUMN gender VARCHAR(10) DEFAULT 'unknown' NOT NULL;
		END IF;
	END
	$$;`

	if _, err := db.Exec(addGenderColumn); err != nil {
		return fmt.Errorf("failed to add gender column: %w", err)
	}

	// Add gender constraint after ensuring column exists
	genderConstraint := `
	DO $$
	BEGIN
		IF NOT EXISTS (
			SELECT 1 FROM information_schema.table_constraints
			WHERE table_name = 'users' AND constraint_name = 'check_gender'
		) THEN
			ALTER TABLE users ADD CONSTRAINT check_gender
			CHECK (gender IN ('male', 'female', 'unknown'));
		END IF;
	END
	$$;`

	if _, err := db.Exec(genderConstraint); err != nil {
		return fmt.Errorf("failed to add gender constraint: %w", err)
	}

	log.Println("Users table created successfully with gender field")
	return nil
}

// CheckConnection verifies the database connection is working
func CheckConnection(db *sql.DB) error {
	if err := db.Ping(); err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}

	// Check if we can query the database
	var result int
	if err := db.QueryRow("SELECT 1").Scan(&result); err != nil {
		return fmt.Errorf("database query failed: %w", err)
	}

	if result != 1 {
		return fmt.Errorf("unexpected database query result: %d", result)
	}

	log.Println("Database connection verified successfully")
	return nil
}
