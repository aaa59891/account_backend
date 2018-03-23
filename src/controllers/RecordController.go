package controllers

import (
	"net/http"

	"github.com/aaa59891/account_backend/src/models"
	"github.com/gin-gonic/gin"
)

func CreateRecord(c *gin.Context) {
	record := models.Record{}
	if err := c.ShouldBindJSON(&record); err != nil {
		GoToErrorResponse(http.StatusBadRequest, c, err)
		return
	}

	if err := models.Transactional(record.Insert); err != nil {
		status := http.StatusInternalServerError
		if err == models.ErrNoCategory || err == models.ErrNoEmail {
			status = http.StatusBadRequest
		}
		GoToErrorResponse(status, c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": record,
	})
}

func FetchRecords(c *gin.Context) {
	rf := models.RecordForm{}
	if err := c.ShouldBindQuery(&rf); err != nil {
		GoToErrorResponse(http.StatusBadRequest, c, err)
		return
	}

	data, err := rf.GetRecords()
	if err != nil {
		GoToErrorResponse(http.StatusInternalServerError, c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}
