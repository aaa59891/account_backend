package controllers

import (
	"fmt"
	"math"
	"net/http"
	"os"

	"github.com/aaa59891/account_backend/src/db"
	"github.com/aaa59891/account_backend/src/models/pagination"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func GoToErrorResponse(statusCode int, c *gin.Context, err error) {
	defer func() {
		panic(err)
	}()

	if os.Getenv("GO_ENV") == "test" {
		fmt.Println(err)
	}

	errMsg := err.Error()

	c.JSON(statusCode, gin.H{
		"message": errMsg,
	})
}

func Pagination(numPerPage, currentPage int, prepareDb *gorm.DB, model interface{}) (pagination.Pagination, error) {
	var totalCount int
	pg := pagination.Pagination{}
	if err := prepareDb.Model(model).Count(&totalCount).Error; err != nil {
		return pg, err
	}

	totalPage := int(math.Ceil(float64(totalCount) / float64(numPerPage)))

	if err := prepareDb.Limit(numPerPage).Offset(numPerPage * (currentPage - 1)).Find(model).Error; err != nil {
		return pg, err
	}

	if currentPage == 0 {
		currentPage = 1
	}
	pg.Current = currentPage
	pg.Next = currentPage + 1
	pg.Previous = currentPage - 1
	pg.Total = totalPage
	return pg, nil
}

func GetModelById(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if err := db.DB.Find(model, "id = ?", id).Error; err != nil {
			statusCode := http.StatusInternalServerError
			if err == gorm.ErrRecordNotFound {
				statusCode = http.StatusNotFound
			}
			GoToErrorResponse(statusCode, c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": model,
		})
	}
}
