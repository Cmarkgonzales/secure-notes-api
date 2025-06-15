package router

import (
	"net/http"
	"secure-notes-api/controllers"
	"secure-notes-api/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all origins for now
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	api := r.Group("/api/v1")

	api.POST("/register", controllers.Register)
	api.POST("/login", controllers.Login)

	api.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "version": "1.0.0"})
	})

	notes := api.Group("/notes", middleware.AuthMiddleware())
	notes.POST("", controllers.CreateNote)
	notes.GET("", controllers.GetNotes)
	notes.PUT("/:id", controllers.UpdateNote)
	notes.DELETE("/:id", controllers.DeleteNote)

	return r
}
