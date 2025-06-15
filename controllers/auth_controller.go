package controllers

import (
	"net/http"
	"secure-notes-api/config"
	"secure-notes-api/models"
	"secure-notes-api/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func generateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JWT_SECRET))
}

func Register(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	hashedPassword, err := hashPassword(input.Password)
	if err != nil {
		utils.InternalError(c, err)
		return
	}

	user := models.User{Username: input.Username, Password: hashedPassword}
	if err := config.DB.Create(&user).Error; err != nil {
		utils.InternalError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	var user models.User
	if err := config.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		utils.RespondError(c, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	if !checkPasswordHash(input.Password, user.Password) {
		utils.RespondError(c, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	token, err := generateToken(user.ID)
	if err != nil {
		utils.InternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
