package config

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB         *gorm.DB
	JWT_SECRET string
)

func LoadEnv() {
	JWT_SECRET = os.Getenv("JWT_SECRET")
}

func ConnectDatabase() {
	var err error
	DB, err = gorm.Open(sqlite.Open("secure-notes.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
}
