package router

import (
	"secure-notes-api/controllers"
	"secure-notes-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	auth := r.Group("/", middleware.AuthMiddleware())
	auth.POST("/notes", controllers.CreateNote)
	auth.GET("/notes", controllers.GetNotes)
	auth.PUT("/notes/:id", controllers.UpdateNote)
	auth.DELETE("/notes/:id", controllers.DeleteNote)

	return r
}
