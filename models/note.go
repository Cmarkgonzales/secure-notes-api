package models

import (
	"gorm.io/gorm"
)

type Note struct {
	gorm.Model
	UserID  uint   `gorm:"not null"`
	Title   string `gorm:"not null"`
	Content string `gorm:"not null"` // (we'll encrypt/decrypt this later)
}
