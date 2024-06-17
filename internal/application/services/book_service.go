package services

import (
	"strconv"

	"hex/internal/adapters/persistence"
	"hex/pkg/models"

	"hex/internal/adapters/logging"
)

type BookService struct {
	repo   persistence.BookRepository
	logger *logging.MongoDBLogger
}

func NewBookService(repo persistence.BookRepository, logger *logging.MongoDBLogger) *BookService {
	return &BookService{repo: repo, logger: logger}
}

func (s *BookService) CreateBook(book *models.Book) error {
	if err := s.repo.Create(book); err != nil {
		s.logger.Log("ERROR", "Failed to create book: "+err.Error())
		return err
	}
	s.logger.Log("INFO", "Book created: "+book.Title)
	return nil
}

func (s *BookService) ViewAllBooks() ([]models.Book, error) {
	books, err := s.repo.GetAll()
	if err != nil {
		s.logger.Log("ERROR", "Failed to retrieve books: "+err.Error())
		return nil, err
	}
	s.logger.Log("INFO", "Retrieved all books")
	return books, nil
}

func (s *BookService) UpdateBook(book *models.Book) error {
	if err := s.repo.Update(book); err != nil {
		s.logger.Log("ERROR", "Failed to update book: "+err.Error())
		return err
	}
	s.logger.Log("INFO", "Book updated: "+book.Title)
	return nil
}

func (s *BookService) DeleteBook(id string) error {
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
