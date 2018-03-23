package db

import (
	"fmt"
	"log"
	"os"

	"github.com/aaa59891/account_backend/src/configs"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB
var err error

func init() {
	databaseConfig := configs.GetConfig().Database

	DB, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", databaseConfig.Account, databaseConfig.Password, databaseConfig.Host, databaseConfig.Port, databaseConfig.DatabaseName))

	if err != nil {
		log.Fatal(err)
	}

	if os.Getenv("GO_ENV") == "dev" || os.Getenv("GO_ENV") == "test" {
		DB.LogMode(true)
	}

	fmt.Println("Connected to database.")
}
