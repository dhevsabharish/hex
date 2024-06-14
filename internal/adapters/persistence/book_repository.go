package persistence

import (
	"hex/pkg/models"

	"gorm.io/gorm"
)

type BookRepository struct {
	DB *gorm.DB
}

func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{DB: db}
}

func (r *BookRepository) Create(book *models.Book) error {
	return r.DB.Create(book).Error
}

func (r *BookRepository) GetAll() ([]models.Book, error) {
	var books []models.Book
	err := r.DB.Find(&books).Error
	return books, err
}

func (r *BookRepository) Update(book *models.Book) error {
	return r.DB.Save(book).Error
}

func (r *BookRepository) Delete(id string) error {
	return r.DB.Delete(&models.Book{}, id).Error
}

func (r *BookRepository) GetByID(id uint) (*models.Book, error) {
	var book models.Book
	err := r.DB.First(&book, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &book, nil
}
