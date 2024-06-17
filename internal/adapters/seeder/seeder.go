package seeder

import (
	"hex/pkg/models"
	"math/rand"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	// Delete existing books and borrowing records
	db.Migrator().DropTable(&models.BorrowingRecord{}, &models.Book{})
	db.AutoMigrate(&models.Book{}, &models.BorrowingRecord{})

	// Create 10 random book records
	books := make([]models.Book, 10)
	for i := 0; i < 10; i++ {
		books[i] = models.Book{
			Title:           generateRandomTitle(),
			Author:          generateRandomAuthor(),
			PublicationDate: datatypes.Date(generateRandomPublicationDate()),
			Genre:           generateRandomGenre(),
			Availability:    1,
		}
	}
	if err := db.Create(&books).Error; err != nil {
		return err
	}

	return nil
}

func generateRandomTitle() string {
	titles := []string{"Book 1", "Book 2", "Book 3", "Book 4", "Book 5"}
	return titles[rand.Intn(len(titles))]
}

func generateRandomAuthor() string {
	authors := []string{"Author 1", "Author 2", "Author 3", "Author 4", "Author 5"}
	return authors[rand.Intn(len(authors))]
}

func generateRandomPublicationDate() time.Time {
	startDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)
	days := int(endDate.Sub(startDate).Hours() / 24)
	randomDays := rand.Intn(days + 1)
	return startDate.AddDate(0, 0, randomDays)
}

func generateRandomGenre() string {
	genres := []string{"Fiction", "Mystery", "Romance", "Thriller", "Science Fiction"}
	return genres[rand.Intn(len(genres))]
}
