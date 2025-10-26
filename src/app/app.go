package app

import (
	"fmt"

	"github.com/alex-1900/wishlist/src/http/route"
	"github.com/gin-gonic/gin"
)

type App struct {
	Config    AppConfig
	GinEngine *gin.Engine
}

func NewApp() *App {
	app := new(App)
	app.Config = config
	app.GinEngine = buildGinEngine()
	return app
}

func (a App) PrepareRoutes() {
	route.RouteDefination(a.GinEngine)
}

func (a App) Listen() {
	if a.GinEngine.Run() != nil {
		fmt.Println("gin error.")
	}
}
