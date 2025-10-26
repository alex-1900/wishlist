package app

import "github.com/gin-gonic/gin"

func buildGinEngine() *gin.Engine {
	return gin.Default()
}
