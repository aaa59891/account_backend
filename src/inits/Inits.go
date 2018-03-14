package inits

import (
	"github.com/aaa59891/account_backend/src/db"
)

func CreateTable() {
	modelArr := make([]interface{}, 0)

	for _, model := range modelArr {
		db.DB.Set("gorm:table_options", "CHARACTER SET = utf8").AutoMigrate(model)
	}
}

func RegisterStruct() {
}
