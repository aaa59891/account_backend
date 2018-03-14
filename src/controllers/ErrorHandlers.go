package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler404(c *gin.Context) {
	GoToErrorResponse(http.StatusNotFound, c, errors.New("Page not found"))
}
