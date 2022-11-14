package src

import (
	"Golangapi/call"

	"github.com/gin-gonic/gin"
)

func AddUserRouter(r *gin.RouterGroup) {
	user := r.Group("/account")
	doc := r.Group("/documents")

	user.GET("/getall", call.GetAllUsers)
	user.GET("/:id", call.HandleBA, call.GetUsers)

	user.PUT("/:id", call.HandleBA, call.PutUser)

	user.POST("", call.PostUsers)

	user.DELETE("/:id", call.DeleteUser)
	user.DELETE("/delete", call.DeleteAllUsers)

	////////////////////////////////
	doc.GET("", call.HandleBA_doc, call.GetAlldocs)
	doc.GET("/:id", call.HandleBA_doc, call.Getdoc)

	doc.POST("", call.HandleBA_doc, call.Postdocs)
	//doc.MaxMultipartMemory = 8 << 20 // 8 MiB

	doc.DELETE("/:id", call.HandleBA_doc, call.DeleteDoc)
	doc.DELETE("/delete", call.HandleBA_doc, call.DeleteAllDoc)

}
