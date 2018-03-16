package inits

import (
	"github.com/aaa59891/account_backend/src/db"
	"github.com/aaa59891/account_backend/src/models"
)

func CreateTable() {
	modelArr := make([]interface{}, 0)
	modelArr = append(modelArr, models.User{})
	modelArr = append(modelArr, models.Category{})
	modelArr = append(modelArr, models.Record{})

	for _, model := range modelArr {
		db.DB.Set("gorm:table_options", "CHARACTER SET = utf8").AutoMigrate(model)
	}
}

func RegisterStruct() {
}
