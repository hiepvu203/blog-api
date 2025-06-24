package config

import (
	"blog-api/internal/entities"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func LoadEnv(){
	err := godotenv.Load()
	if err != nil{
		log.Fatal("Error when load .env file")
	}
}

func ConnectDB(){
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	var err error

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Cannot connect database: ", err)
	}

	fmt.Println("Connected")
}

func InitDB() {
	err := DB.AutoMigrate(
		&entities.User{},
		&entities.Category{},
		&entities.Post{},
		&entities.Comment{},
	)

	if err != nil{
		log.Fatal("AutoMigrate failed: ", err)
	}
}