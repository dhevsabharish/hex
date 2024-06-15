package services

import (
	"fmt"
	"hex/internal/adapters/persistence"
	"hex/pkg/models"
	"time"

	"hex/internal/adapters/logging"
)

type BorrowingService interface {
	BorrowBook(bookID uint, memberID uint) error
	ReturnBook(borrowingRecordID uint, memberID uint) error
	GetMyBorrowings(memberID uint) ([]models.BorrowingRecord, error)
	GetAllBorrowingRecords() ([]models.BorrowingRecord, error)
}

type borrowingService struct {
	bookRepo      persistence.BookRepository
	borrowingRepo persistence.BorrowingRepository
	logger        *logging.MongoDBLogger
}

func NewBorrowingService(bookRepo persistence.BookRepository, borrowingRepo persistence.BorrowingRepository, logger *logging.MongoDBLogger) BorrowingService {
	return &borrowingService{
		bookRepo:      bookRepo,
		borrowingRepo: borrowingRepo,
		logger:        logger,
	}
}

func (s *borrowingService) BorrowBook(bookID uint, memberID uint) error {
	// Check if the book is available
	book, err := s.bookRepo.GetByID(bookID)
	if err != nil {
		s.logger.Log("ERROR", "Failed to get book by ID: "+err.Error())
		return err
	}
	if book == nil {
		err := fmt.Errorf("book not found")
		s.logger.Log("ERROR", err.Error())
		return err
	}
	if book.Availability == 0 {
		err := fmt.Errorf("book is not available")
		s.logger.Log("ERROR", err.Error())
		return err
	}

	// Create a new borrowing record
	borrowingRecord := models.BorrowingRecord{
		BookID:     bookID,
		MemberID:   memberID,
		BorrowDate: time.Now(),
	}
	if err := s.borrowingRepo.Create(&borrowingRecord); err != nil {
		s.logger.Log("ERROR", "Failed to create borrowing record: "+err.Error())
		return err
	}

	// Update the book availability
	book.Availability--
	if err := s.bookRepo.Update(book); err != nil {
		s.logger.Log("ERROR", "Failed to update book availability: "+err.Error())
		return err
	}

	s.logger.Log("INFO", fmt.Sprintf("Book borrowed: bookID=%d, memberID=%d", bookID, memberID))
	return nil
}

func (s *borrowingService) ReturnBook(borrowingRecordID uint, memberID uint) error {
	borrowingRecord, err := s.borrowingRepo.GetByID(borrowingRecordID)
	if err != nil {
		s.logger.Log("ERROR", "Failed to get borrowing record by ID: "+err.Error())
		return err
	}
	if borrowingRecord == nil {
		err := fmt.Errorf("borrowing record not found")
		s.logger.Log("ERROR", err.Error())
		return err
	}

	// Check if the book belongs to the user
	if borrowingRecord.MemberID != memberID {
		err := fmt.Errorf("unauthorized: you can only return books you borrowed")
		s.logger.Log("ERROR", err.Error())
		return err
	}

	// Check if the book is already returned
	if !borrowingRecord.ReturnDate.IsZero() {
		err := fmt.Errorf("book is already returned")
		s.logger.Log("ERROR", err.Error())
		return err
	}

	book, err := s.bookRepo.GetByID(borrowingRecord.BookID)
	if err != nil {
		s.logger.Log("ERROR", "Failed to get book by ID: "+err.Error())
		return err
	}
	if book == nil {
		err := fmt.Errorf("book not found")
		s.logger.Log("ERROR", err.Error())
		return err
	}

	borrowingRecord.ReturnDate = time.Now()
	if err := s.borrowingRepo.Update(borrowingRecord); err != nil {
		s.logger.Log("ERROR", "Failed to update borrowing record: "+err.Error())
		return err
	}

	book.Availability++
	if err := s.bookRepo.Update(book); err != nil {
		s.logger.Log("ERROR", "Failed to update book availability: "+err.Error())
		return err
	}

	s.logger.Log("INFO", fmt.Sprintf("Book returned: bookID=%d, memberID=%d", borrowingRecord.BookID, memberID))
	return nil
}

func (s *borrowingService) GetMyBorrowings(memberID uint) ([]models.BorrowingRecord, error) {
	borrowingRecords, err := s.borrowingRepo.GetByMemberID(memberID)
	if err != nil {
		s.logger.Log("ERROR", "Failed to get borrowing records by member ID: "+err.Error())
		return nil, err
	}

	s.logger.Log("INFO", fmt.Sprintf("Retrieved borrowing records for memberID=%d", memberID))
	return borrowingRecords, nil
}

func (s *borrowingService) GetAllBorrowingRecords() ([]models.BorrowingRecord, error) {
	borrowingRecords, err := s.borrowingRepo.GetAll()
	if err != nil {
		s.logger.Log("ERROR", "Failed to get all borrowing records: "+err.Error())
		return nil, err
	}

	s.logger.Log("INFO", "Retrieved all borrowing records")
	return borrowingRecords, nil
}
