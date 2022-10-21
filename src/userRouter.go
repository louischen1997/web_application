package src

import (
	"Golangapi/call"

	"github.com/gin-gonic/gin"
)

func AddUserRouter(r *gin.RouterGroup) {
	user := r.Group("/account")

	user.GET("/getall", call.GetAllUsers)
	user.GET("/:id", call.HandleBA, call.GetUsers)

	user.PUT("/:id", call.HandleBA, call.PutUser)

	user.POST("", call.PostUsers)

	user.DELETE("/:id", call.DeleteUser)
	user.DELETE("/delete", call.DeleteAllUsers)

}
