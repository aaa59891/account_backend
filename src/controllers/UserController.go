package controllers

import (
	"net/http"

	"github.com/aaa59891/account_backend/src/models"
	"github.com/gin-gonic/gin"
)

func SingUp(c *gin.Context) {
	user := models.User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		GoToErrorResponse(http.StatusBadRequest, c, err)
		return
	}
	if err := models.Transactional(user.Insert); err != nil {
		status := http.StatusInternalServerError
		if err == models.ErrEmailExist {
			status = http.StatusBadRequest
		}
		GoToErrorResponse(status, c, err)
		return
	}
	c.JSON(http.StatusOK, nil)
}

func SignIn(c *gin.Context) {
	user := models.User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		GoToErrorResponse(http.StatusBadRequest, c, err)
		return
	}
	if err := user.CheckPassword(); err != nil {
		status := http.StatusInternalServerError
		if err == models.ErrWrongPassword {
			status = http.StatusBadRequest
		}
		GoToErrorResponse(status, c, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}
