package main

import (
	"net/http"

	"Golangapi/config"
	"Golangapi/mdl"
	"Golangapi/src"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/healthz", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	v1 := router.Group("/v1")
	src.AddUserRouter(v1)

	go func() {
		config.SetupDatabaseConnection()
		config.DB.AutoMigrate(&mdl.Dbtable{})
	}()

	router.Run(":3000")
	// successful

}
