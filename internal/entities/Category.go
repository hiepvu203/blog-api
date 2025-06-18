package entities

type Category struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"type:varchar(100);not null"`
	Slug string `gorm:"type:varchar(100);unique;not null"`
	
	Posts []Post `gorm:"foreignKey:CategoryID"`
}