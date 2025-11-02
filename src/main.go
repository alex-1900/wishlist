package main

import (
	"fmt"

	"github.com/alex-1900/wishlist/src/app"
	"github.com/alex-1900/wishlist/src/module"
)

func main() {
	// Get app instance from dependency manager
	app := app.GetInstance()

	// Register routes
	module.RouteDefinition(app.GinEngine)

	// Start server
	if app.GinEngine.Run() != nil {
		fmt.Println("gin error.")
	}
}
