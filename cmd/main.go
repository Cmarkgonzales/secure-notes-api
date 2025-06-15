package main

import (
	"secure-notes-api/config"
	"secure-notes-api/models"
	"secure-notes-api/router"
	"secure-notes-api/utils"
)

func main() {
	config.LoadEnv()
	utils.LoadEncryptionKey()
	config.ConnectDatabase()
	config.DB.AutoMigrate(&models.User{}, &models.Note{})

	r := router.SetupRouter()
	r.Run(":8086")
}
