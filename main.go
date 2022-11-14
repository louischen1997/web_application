package main

import (
	"net/http"

	"Golangapi/call"
	"Golangapi/config"
	"Golangapi/mdl"
	"Golangapi/mdldb"
	"Golangapi/src"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	call.MyLogger()

	err := godotenv.Load()
	if err != nil {
		call.Logger.Print("Error To Load env File")
	}

	router := gin.Default()
	router.GET("/healthz", func(c *gin.Context) {
		call.Logger.Print("Get Request, Healthz Successful, 200 OK")
		call.HandleMetricCounter("Healthz")
		c.Status(http.StatusOK)
	})
	v1 := router.Group("/v1")
	src.AddUserRouter(v1)

	go func() {
		config.SetupDatabaseConnection()
		config.DB.AutoMigrate(&mdl.User{})
		config.DB.AutoMigrate(&mdldb.Doc{})

	}()
	router.Run(":3000")
	call.Logger.Print("Start The Application, Run On Route 3000")

}
