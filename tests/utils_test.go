package tests

import (
	"log"
	"os"
	"secure-notes-api/config"
	"secure-notes-api/models"
	"secure-notes-api/router"
	"secure-notes-api/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var TestRouter *gin.Engine

func InitTestSuite() {
	gin.SetMode(gin.TestMode)
	loadEnv()
	initTestDB()
	initRouter()
}

func loadEnv() {
	// Use test-specific secrets
	os.Setenv("JWT_SECRET", "test_secret_123456")
	os.Setenv("ENCRYPTION_KEY", "jcXeHPAikPW3XmhpF++6/B5wacjH+3UHBc29ngFJAHc=")
	log.Println("[TEST] Environment variables loaded")

	utils.LoadEncryptionKey()
	log.Println("[TEST] Encryption key loaded from environment")
}

func initTestDB() {
	var err error

	// Use in-memory sqlite for isolation
	config.DB, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto-migrate schemas
	err = config.DB.AutoMigrate(&models.User{}, &models.Note{})
	if err != nil {
		log.Fatalf("Failed to migrate test database: %v", err)
	}

	log.Println("[TEST] Test database initialized")
}

func initRouter() {
	log.Println("[TEST] Encryption key loaded")
	TestRouter = router.SetupRouter()
	log.Println("[TEST] Router initialized")
}

func ClearDB() {
	// Clean all test data
	config.DB.Exec("DELETE FROM notes")
	config.DB.Exec("DELETE FROM users")
	log.Println("[TEST] Database cleared")
}
