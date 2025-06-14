package routes

import (
	"secure-notes-api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	authorized := r.Group("/")
	authorized.Use(controllers.AuthMiddleware())
	{
		authorized.POST("/notes", controllers.CreateNote)
		authorized.GET("/notes", controllers.GetNotes)
	}
}
