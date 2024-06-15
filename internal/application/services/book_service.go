package services

import (
	"fmt"
	"strconv"
	"time"

	"hex/internal/adapters/persistence"
	"hex/internal/application/auth"
	"hex/pkg/models"

	"hex/internal/adapters/logging"

	"gorm.io/datatypes"
)

type BookService struct {
	repo        persistence.BookRepository
	authService auth.AuthService
	logger      *logging.MongoDBLogger
}

func NewBookService(repo persistence.BookRepository, authService auth.AuthService, logger *logging.MongoDBLogger) *BookService {
	return &BookService{repo: repo, authService: authService, logger: logger}
}

func (s *BookService) CreateBook(book *models.Book, publicationDateStr string, token string) error {
	_, role, err := s.authService.Authenticate(token)
	if err != nil {
		s.logger.Log("ERROR", "Authentication failed: "+err.Error())
		return err
	}

	if role != "admin" && role != "librarian" {
		err := fmt.Errorf("unauthorized: only admins and librarians can create books")
		s.logger.Log("ERROR", err.Error())
		return err
	}
	layout := "2006-01-02"
	parsedDate, err := time.Parse(layout, publicationDateStr)
	if err != nil {
		s.logger.Log("ERROR", "Failed to parse publication date: "+err.Error())
		return err
	}
	book.PublicationDate = datatypes.Date(parsedDate)
	if err := s.repo.Create(book); err != nil {
		s.logger.Log("ERROR", "Failed to create book: "+err.Error())
		return err
	}
	s.logger.Log("INFO", "Book created: "+book.Title)
	return nil
}

func (s *BookService) ViewAllBooks(token string) ([]models.Book, error) {
	_, _, err := s.authService.Authenticate(token)
	if err != nil {
		s.logger.Log("ERROR", "Authentication failed: "+err.Error())
		return nil, err
	}
	books, err := s.repo.GetAll()
	if err != nil {
		s.logger.Log("ERROR", "Failed to retrieve books: "+err.Error())
		return nil, err
	}
	s.logger.Log("INFO", "Retrieved all books")
	return books, nil
}

func (s *BookService) UpdateBook(book *models.Book, publicationDateStr string, token string) error {
	_, role, err := s.authService.Authenticate(token)
	if err != nil {
		s.logger.Log("ERROR", "Authentication failed: "+err.Error())
		return err
	}

	if role != "admin" && role != "librarian" {
		err := fmt.Errorf("unauthorized: only admins and librarians can update books")
		s.logger.Log("ERROR", err.Error())
		return err
	}
	if publicationDateStr != "" {
		layout := "2006-01-02"
		parsedDate, err := time.Parse(layout, publicationDateStr)
		if err != nil {
			s.logger.Log("ERROR", "Failed to parse publication date: "+err.Error())
			return err
		}
		book.PublicationDate = datatypes.Date(parsedDate)
	}

	if err := s.repo.Update(book); err != nil {
		s.logger.Log("ERROR", "Failed to update book: "+err.Error())
		return err
	}
	s.logger.Log("INFO", "Book updated: "+book.Title)
	return nil
}

func (s *BookService) DeleteBook(id string, token string) error {
	_, role, err := s.authService.Authenticate(token)
	if err != nil {
		s.logger.Log("ERROR", "Authentication failed: "+err.Error())
		return err
	}

	if role != "admin" && role != "librarian" {
		err := fmt.Errorf("unauthorized: only admins and librarians can delete books")
		s.logger.Log("ERROR", err.Error())
		return err
	}
	if err := s.repo.Delete(id); err != nil {
		s.logger.Log("ERROR", "Failed to delete book: "+err.Error())
		return err
	}
	s.logger.Log("INFO", "Book deleted: ID "+id)
	return nil
}

func (s *BookService) GetBookByID(id string) (*models.Book, error) {
	bookID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		s.logger.Log("ERROR", "Invalid book ID: "+err.Error())
		return nil, err
	}

	book, err := s.repo.GetByID(uint(bookID))
	if err != nil {
		s.logger.Log("ERROR", "Failed to get book by ID: "+err.Error())
		return nil, err
	}
	s.logger.Log("INFO", "Retrieved book by ID: "+id)
	return book, nil
}
