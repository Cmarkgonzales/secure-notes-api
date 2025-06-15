package models

import (
	"gorm.io/gorm"
)

type Note struct {
	gorm.Model
	UserID  uint   `json:"userId"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
