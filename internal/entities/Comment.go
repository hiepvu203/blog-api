package entities

import (
	"time"
	"gorm.io/gorm"
)

type Comment struct {
	ID        uint      `gorm:"primaryKey"`
	PostID    uint
	UserID    uint
	Content   string    `gorm:"type:text;not null"`
	CreatedAt time.Time

	Post Post
	User User	

	DeletedAt   gorm.DeletedAt `gorm:"index"`
}