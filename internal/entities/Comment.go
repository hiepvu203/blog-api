package entities

import "time"

type Comment struct {
	ID        uint      `gorm:"primaryKey"`
	PostID    uint
	UserID    uint
	Content   string    `gorm:"type:text;not null"`
	CreatedAt time.Time

	Post Post
	User User	
}