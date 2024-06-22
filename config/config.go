package config

import (
	"hex/internal/adapters/logging"
	"hex/pkg/models"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	DB           *gorm.DB
	Port         string
	RailsAPIURL  string
	Logger       *logging.MongoDBLogger
	SeedDatabase bool
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file, using environment variables")
	}

	var db *gorm.DB
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		dsn := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_HOST") + ":3306)/" + os.Getenv("DB_NAME") + "?charset=utf8mb4&parseTime=True&loc=Local"
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(time.Second * 5)
	}
	if err != nil {
		log.Fatalf("Error connecting to database after %d attempts: %v", maxRetries, err)
	}

	db.AutoMigrate(&models.Book{}, &models.BorrowingRecord{})

	logger := logging.NewMongoDBLogger(os.Getenv("MONGODB_URI"), os.Getenv("MONGODB_DB"), os.Getenv("MONGODB_COLLECTION"))

	seedDatabase, err := strconv.ParseBool(os.Getenv("SEED_DATABASE"))
	if err != nil {
		seedDatabase = false
	}

	return &Config{
		DB:           db,
		Port:         os.Getenv("PORT"),
		RailsAPIURL:  os.Getenv("RAILS_API_URL"),
		Logger:       logger,
		SeedDatabase: seedDatabase,
	}
}
