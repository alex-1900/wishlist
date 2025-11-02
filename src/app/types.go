package app

import (
	"database/sql"

	"github.com/alex-1900/wishlist/src/repository"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type AppConfig struct {
	AppName       string
	Database      DatabaseConfig
	JWTSecret     string
	JWTExpiration int // in hours
}

type App struct {
	Config     AppConfig
	GinEngine  *gin.Engine
	DB         *sql.DB
	Repository *repository.RepositoryManager
}
