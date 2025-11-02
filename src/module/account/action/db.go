package action

import (
	"net/http"

	"github.com/alex-1900/wishlist/src/app"
	"github.com/gin-gonic/gin"
)

func ActionDBTest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db := app.GetDB()

		// Test database connection with a simple query
		var result int
		if err := db.QueryRow("SELECT 1").Scan(&result); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Database connection failed",
				"error":   err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Database connection successful",
			"result":  result,
		})
	}
}
