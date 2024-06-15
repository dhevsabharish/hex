package config

import (
	"log"
	"os"

	"hex/pkg/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	DB           *gorm.DB
	Port         string
	RailsAPIURL  string
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dsn := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(127.0.0.1:3306)/" + os.Getenv("DB_NAME") + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	db.AutoMigrate(&models.Book{}, &models.BorrowingRecord{})

	return &Config{
		DB:           db,
		Port:         os.Getenv("PORT"),
		RailsAPIURL:  os.Getenv("RAILS_API_URL"),
	}
}