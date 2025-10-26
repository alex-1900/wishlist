package route

import (
	"github.com/alex-1900/wishlist/src/http/module/action"
	"github.com/gin-gonic/gin"
)

func RouteDefination(router *gin.Engine) {
	router.GET("/ping", action.ActionPing())
}
