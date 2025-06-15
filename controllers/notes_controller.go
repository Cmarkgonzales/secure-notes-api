package controllers

import (
	"net/http"
	"secure-notes-api/config"
	"secure-notes-api/models"
	"secure-notes-api/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func extractUserId(c *gin.Context) (uint, bool) {
	userIdInterface, exists := c.Get("userId")
	if !exists {
		return 0, false
	}
	userId, ok := userIdInterface.(uint)
	return userId, ok
}

func CreateNote(c *gin.Context) {
	var note models.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	userId, ok := extractUserId(c)
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "User ID missing")
		return
	}
	note.UserID = userId

	encryptedTitle, err := utils.Encrypt(note.Title)
	if err != nil {
		utils.InternalError(c, err)
		return
	}
	encryptedContent, err := utils.Encrypt(note.Content)
	if err != nil {
		utils.InternalError(c, err)
		return
	}

	note.Title = encryptedTitle
	note.Content = encryptedContent

	if err := config.DB.Create(&note).Error; err != nil {
		utils.InternalError(c, err)
		return
	}

	c.JSON(http.StatusCreated, note)
}

func GetNotes(c *gin.Context) {
	userId, ok := extractUserId(c)
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "User ID missing")
		return
	}

	search := c.Query("search")
	var notes []models.Note
	if err := config.DB.Where("user_id = ?", userId).Find(&notes).Error; err != nil {
		utils.InternalError(c, err)
		return
	}

	var decryptedNotes []models.Note
	for _, note := range notes {
		decryptedTitle, err1 := utils.Decrypt(note.Title)
		decryptedContent, err2 := utils.Decrypt(note.Content)
		if err1 != nil || err2 != nil {
			continue
		}
		note.Title = decryptedTitle
		note.Content = decryptedContent

		if search == "" || utils.ContainsIgnoreCase(note.Title, search) || utils.ContainsIgnoreCase(note.Content, search) {
			decryptedNotes = append(decryptedNotes, note)
		}
	}

	c.JSON(http.StatusOK, decryptedNotes)
}

func UpdateNote(c *gin.Context) {
	userId, ok := extractUserId(c)
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "User ID missing")
		return
	}

	noteId := c.Param("id")
	id, err := strconv.Atoi(noteId)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid note ID")
		return
	}

	var existingNote models.Note
	if err := config.DB.Where("id = ? AND user_id = ?", id, userId).First(&existingNote).Error; err != nil {
		utils.RespondError(c, http.StatusNotFound, "Note not found")
		return
	}

	var updatedNote models.Note
	if err := c.ShouldBindJSON(&updatedNote); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	encryptedTitle, err := utils.Encrypt(updatedNote.Title)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Encryption failed")
		return
	}
	encryptedContent, err := utils.Encrypt(updatedNote.Content)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Encryption failed")
		return
	}

	existingNote.Title = encryptedTitle
	existingNote.Content = encryptedContent

	if err := config.DB.Save(&existingNote).Error; err != nil {
		utils.InternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, existingNote)
}

func DeleteNote(c *gin.Context) {
	userId, ok := extractUserId(c)
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "User ID missing")
		return
	}

	noteId := c.Param("id")
	id, err := strconv.Atoi(noteId)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid note ID")
		return
	}

	var note models.Note
	if err := config.DB.Where("id = ? AND user_id = ?", id, userId).First(&note).Error; err != nil {
		utils.RespondError(c, http.StatusNotFound, "Note not found")
		return
	}

	if err := config.DB.Delete(&note).Error; err != nil {
		utils.InternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
}
