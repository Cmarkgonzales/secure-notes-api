package controllers

import (
	"net/http"
	"secure-notes-api/config"
	"secure-notes-api/models"

	"github.com/gin-gonic/gin"
)

// CreateNote handles POST /notes
func CreateNote(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var input struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note := models.Note{
		UserID:  userID,
		Title:   input.Title,
		Content: input.Content, // encryption can be applied here later
	}

	result := config.DB.Create(&note)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"note": note})
}

// GetNotes handles GET /notes
func GetNotes(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var notes []models.Note
	result := config.DB.Where("user_id = ?", userID).Find(&notes)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notes": notes})
}
