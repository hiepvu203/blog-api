package entities

import (
	"time"
	"gorm.io/gorm"
)

type Post struct {
	ID          uint      `gorm:"primaryKey"`
	Title       string    `gorm:"type:varchar(200);not null"`
	Slug        string    `gorm:"type:varchar(200);unique;not null"`
	Content     string    `gorm:"type:text;not null"`
	Thumbnail   string    `gorm:"type:text"`
	CategoryID  uint
	AuthorID    uint
	Status      string    `gorm:"type:post_status;default:'draft'"` // ENUM
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// Relationships
	Author   User
	Category Category
	Comments []Comment `gorm:"foreignKey:PostID"`

	DeletedAt   gorm.DeletedAt `gorm:"index"`
}