package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Cors(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
	c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, access-control-allow-origin, access-control-allow-headers, api_token")
	if c.Request.Method == http.MethodOptions {
		c.AbortWithStatus(200)
		return
	}
	c.Next()
}
