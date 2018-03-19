package main

import (
	"github.com/aaa59891/account_backend/src/configs"
	"github.com/aaa59891/account_backend/src/controllers"
	"github.com/aaa59891/account_backend/src/db"
	"github.com/aaa59891/account_backend/src/inits"
	"github.com/aaa59891/account_backend/src/middlewares"
	"github.com/aaa59891/account_backend/src/routers"
	"github.com/gin-gonic/gin"
)

func init() {
	inits.CreateTable()
	inits.RegisterStruct()
}

func main() {
	defer db.DB.Close()
	config := configs.GetConfig()
	r := gin.Default()
	r.Use(middlewares.Cors)

	r.NoRoute(controllers.ErrorHandler404)
	routers.SetRoutes(r)

	r.Run(config.Server.Port)
}
