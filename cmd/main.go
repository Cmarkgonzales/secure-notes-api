package main

import (
	"secure-notes-api/config"
	"secure-notes-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	config.ConnectDatabase()

	r := gin.Default()
	routes.RegisterRoutes(r)

	port := config.GetEnv("PORT", "8086")
	r.Run(":" + port)
}
