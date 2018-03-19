package routers

import (
	"net/http"

	"github.com/aaa59891/account_backend/src/controllers"

	"github.com/gin-gonic/gin"
)

func SetRoutes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})
	r.POST("/user", controllers.SingUp)
	r.POST("/signin", controllers.SignIn)

	r.POST("/category", controllers.CreateCategory)
	r.PUT("/category", controllers.UpdateCategory)
	r.DELETE("/category/:id", controllers.DeleteCategory)
	r.GET("/category/:email", controllers.FetchCategories)

	r.POST("/record", controllers.CreateRecord)
}
