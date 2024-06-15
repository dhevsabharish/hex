package services

import (
	"fmt"
	"hex/internal/adapters/persistence"
	"hex/pkg/models"
	"time"
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
}

func NewBorrowingService(bookRepo persistence.BookRepository, borrowingRepo persistence.BorrowingRepository) BorrowingService {
	return &borrowingService{
		bookRepo:      bookRepo,
		borrowingRepo: borrowingRepo,
	}
}

func (s *borrowingService) BorrowBook(bookID uint, memberID uint) error {
	// Check if the book is available
	book, err := s.bookRepo.GetByID(bookID)
	if err != nil {
		return err
	}
	if book == nil {
		return fmt.Errorf("book not found")
	}
	if book.Availability == 0 {
		return fmt.Errorf("book is not available")
	}

	// Create a new borrowing record
	borrowingRecord := models.BorrowingRecord{
		BookID:     bookID,
		MemberID:   memberID,
		BorrowDate: time.Now(),
	}
	if err := s.borrowingRepo.Create(&borrowingRecord); err != nil {
		return err
	}

	// Update the book availability
	book.Availability--
	if err := s.bookRepo.Update(book); err != nil {
		return err
	}

	return nil
}

func (s *borrowingService) ReturnBook(borrowingRecordID uint, memberID uint) error {
	borrowingRecord, err := s.borrowingRepo.GetByID(borrowingRecordID)
	if err != nil {
		return err
	}
	if borrowingRecord == nil {
		return fmt.Errorf("borrowing record not found")
	}

	// Check if the book belongs to the user
	if borrowingRecord.MemberID != memberID {
		return fmt.Errorf("unauthorized: you can only return books you borrowed")
	}

	// Check if the book is already returned
	if !borrowingRecord.ReturnDate.IsZero() {
		return fmt.Errorf("book is already returned")
	}

	book, err := s.bookRepo.GetByID(borrowingRecord.BookID)
	if err != nil {
		return err
	}
	if book == nil {
		return fmt.Errorf("book not found")
	}

	borrowingRecord.ReturnDate = time.Now()
	if err := s.borrowingRepo.Update(borrowingRecord); err != nil {
		return err
	}

	book.Availability++
	if err := s.bookRepo.Update(book); err != nil {
		return err
	}

	return nil
}

func (s *borrowingService) GetMyBorrowings(memberID uint) ([]models.BorrowingRecord, error) {
	borrowingRecords, err := s.borrowingRepo.GetByMemberID(memberID)
	if err != nil {
		return nil, err
	}

	return borrowingRecords, nil
}

func (s *borrowingService) GetAllBorrowingRecords() ([]models.BorrowingRecord, error) {
	borrowingRecords, err := s.borrowingRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return borrowingRecords, nil
}