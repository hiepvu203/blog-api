package entities

import "time"

type User struct {
	ID 			int			`gorm:"primaryKey"`		
	Username	string		`gorm:"type:varchar(50);unique;not null"`
	Email 		string		`gorm:"type:varchar(100);unique;not null"`
	Password	string		`gorm:"type:varchar(255);not null"`
	Role		string 		`gorm:"type:user_role;default:'client';not null"` // ENUM
	CreatedAt	time.Time

	// relationships
	Posts 		[]Post 		`gorm:"foreignKey:AuthorID"`	
	Comments	[]Comment 	`gorm:"foreignKey:UserID"`
}