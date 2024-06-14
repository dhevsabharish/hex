package services

import (
	"fmt"
	"strconv"
	"time"

	"hex/internal/adapters/persistence"
	"hex/internal/application/auth"
	"hex/pkg/models"

	"gorm.io/datatypes"
)

type BookService struct {
	repo        persistence.BookRepository
	authService auth.AuthService
}

func NewBookService(repo persistence.BookRepository, authService auth.AuthService) *BookService {
	return &BookService{repo: repo, authService: authService}
}

func (s *BookService) CreateBook(book *models.Book, publicationDateStr string, token string) error {
	_, role, err := s.authService.Authenticate(token)
	if err != nil {
		return err
	}

	if role != "admin" && role != "librarian" {
		return fmt.Errorf("unauthorized: only admins and librarians can create books")
	}
	layout := "2006-01-02"
	parsedDate, err := time.Parse(layout, publicationDateStr)
	if err != nil {
		return err
	}
	book.PublicationDate = datatypes.Date(parsedDate)
	return s.repo.Create(book)
}

func (s *BookService) ViewAllBooks(token string) ([]models.Book, error) {
	_, _, err := s.authService.Authenticate(token)
	if err != nil {
		return nil, err
	}
	return s.repo.GetAll()
}

func (s *BookService) UpdateBook(book *models.Book, publicationDateStr string, token string) error {
	_, role, err := s.authService.Authenticate(token)
	if err != nil {
		return err
	}

	if role != "admin" && role != "librarian" {
		return fmt.Errorf("unauthorized: only admins and librarians can update books")
	}
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

func (s *BookService) DeleteBook(id string, token string) error {
	_, role, err := s.authService.Authenticate(token)
	if err != nil {
		return err
	}

	if role != "admin" && role != "librarian" {
		return fmt.Errorf("unauthorized: only admins and librarians can delete books")
	}
	return s.repo.Delete(id)
}

func (s *BookService) GetBookByID(id string) (*models.Book, error) {
	bookID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, err
	}

	return s.repo.GetByID(uint(bookID))
}
