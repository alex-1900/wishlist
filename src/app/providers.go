package app

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/alex-1900/wishlist/src/database"
	"github.com/alex-1900/wishlist/src/repository"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func buildApp() *App {
	app := new(App)
	app.Config = config

	// Build database connection
	db, err := buildDatabaseConnection(app.Config.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	app.DB = db

	// Initialize database schema
	if err := database.InitializeSchema(db); err != nil {
		log.Fatalf("Failed to initialize database schema: %v", err)
	}

	// Initialize repository manager
	app.Repository = repository.NewRepositoryManager(db)

	app.GinEngine = buildGinEngine()
	return app
}

func buildGinEngine() *gin.Engine {
	return gin.Default()
}

func buildDatabaseConnection(dbConfig DatabaseConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.SSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL database")
	return db, nil
}
