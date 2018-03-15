package controllers

import (
	"net/http"

	"github.com/aaa59891/account_backend/src/models"
	"github.com/gin-gonic/gin"
)

func CreateCategory(c *gin.Context) {
	category := models.Category{}
	if err := c.ShouldBindJSON(&category); err != nil {
		GoToErrorResponse(http.StatusBadRequest, c, err)
		return
	}

	if err := models.Transactional(category.Insert); err != nil {
		GoToErrorResponse(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}
