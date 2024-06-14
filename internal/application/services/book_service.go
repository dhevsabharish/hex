package services

import (
	"hex/internal/adapters/persistence"
	"hex/pkg/models"
	"strconv"
	"time"

	"gorm.io/datatypes"
)

type BookService struct {
	repo persistence.BookRepository
}

func NewBookService(repo persistence.BookRepository) *BookService {
	return &BookService{repo: repo}
}

func (s *BookService) CreateBook(book *models.Book, publicationDateStr string) error {
	layout := "2006-01-02"
	parsedDate, err := time.Parse(layout, publicationDateStr)
	if err != nil {
		return err
	}
	book.PublicationDate = datatypes.Date(parsedDate)
	return s.repo.Create(book)
}

func (s *BookService) ViewAllBooks() ([]models.Book, error) {
	return s.repo.GetAll()
}

func (s *BookService) UpdateBook(book *models.Book, publicationDateStr string) error {
	if publicationDateStr != "" {
		layout := "2006-01-02"
		parsedDate, err := time.Parse(layout, publicationDateStr)
		if err != nil {
			return err
		}
		book.PublicationDate = datatypes.Date(parsedDate)
	}

	return s.repo.Update(book)
}

func (s *BookService) DeleteBook(id string) error {
	return s.repo.Delete(id)
}

func (s *BookService) GetBookByID(id string) (*models.Book, error) {
	bookID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, err
	}

	return s.repo.GetByID(uint(bookID))
}
