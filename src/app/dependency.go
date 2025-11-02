package app

import (
	"database/sql"
	"sync"

	"github.com/alex-1900/wishlist/src/repository"
	"github.com/gin-gonic/gin"
)

var (
	appInstance *App
	appOnce     sync.Once
)

// GetInstance returns the singleton App instance, initializing it if necessary
func GetInstance() *App {
	appOnce.Do(func() {
		appInstance = buildApp()
	})
	return appInstance
}

// GetGinEngine returns the Gin engine from the App instance
func GetGinEngine() *gin.Engine {
	return GetInstance().GinEngine
}

// GetConfig returns the configuration from the App instance
func GetConfig() AppConfig {
	return GetInstance().Config
}

// GetDB returns the database connection from the App instance
func GetDB() *sql.DB {
	return GetInstance().DB
}

// GetRepository returns the repository manager from the App instance
func GetRepository() *repository.RepositoryManager {
	return GetInstance().Repository
}

// ResetApp resets the singleton instance (mainly for testing)
func ResetApp() {
	appOnce = sync.Once{}
	appInstance = nil
}
