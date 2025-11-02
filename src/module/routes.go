package module

import (
	"github.com/alex-1900/wishlist/src/module/account"
	"github.com/gin-gonic/gin"
)

// RouteDefinition registers all application routes
func RouteDefinition(router *gin.Engine) {
	// Register account module routes
	account.RegisterRoutes(router)
}
