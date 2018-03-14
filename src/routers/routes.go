package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetRoutes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})
}
