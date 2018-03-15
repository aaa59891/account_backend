package controllers

import (
	"net/http"
	"strconv"

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
		status := http.StatusInternalServerError
		if err == models.ErrNoEmail {
			status = http.StatusBadRequest
		}
		GoToErrorResponse(status, c, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func UpdateCategory(c *gin.Context) {
	category := models.Category{}
	if err := c.ShouldBindJSON(&category); err != nil {
		GoToErrorResponse(http.StatusBadRequest, c, err)
		return
	}
	if err := models.Transactional(category.Update); err != nil {
		GoToErrorResponse(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		GoToErrorResponse(http.StatusBadRequest, c, err)
		return
	}
	category := models.Category{Id: uint(idInt)}
	if err := models.Transactional(category.Delete); err != nil {
		GoToErrorResponse(http.StatusInternalServerError, c, err)
		return
	}
	c.JSON(http.StatusOK, nil)
}
