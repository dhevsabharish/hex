package services

import (
	"fmt"
	"hex/internal/adapters/persistence"
	"hex/internal/application/auth"
	"hex/pkg/models"
	"time"
)

type BorrowingService interface {
	BorrowBook(bookID uint, memberID uint, token string) error
	ReturnBook(borrowingRecordID uint, memberID uint, token string) error
	GetMyBorrowings(memberID uint, token string) ([]models.BorrowingRecord, error)
	GetAllBorrowingRecords(token string) ([]models.BorrowingRecord, error)
}

type borrowingService struct {
	bookRepo      persistence.BookRepository
	borrowingRepo persistence.BorrowingRepository
	authService   auth.AuthService
}

func NewBorrowingService(bookRepo persistence.BookRepository, borrowingRepo persistence.BorrowingRepository, authService auth.AuthService) BorrowingService {
	return &borrowingService{
		bookRepo:      bookRepo,
		borrowingRepo: borrowingRepo,
		authService:   authService,
	}
}

func (s *borrowingService) BorrowBook(bookID uint, memberID uint, token string) error {
	// Check if the user is a member
	_, role, err := s.authService.Authenticate(token)
	if err != nil {
		return err
	}
	if role != "member" {
		return fmt.Errorf("unauthorized: only members can borrow books")
	}

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

func (s *borrowingService) ReturnBook(borrowingRecordID uint, memberID uint, token string) error {
	_, role, err := s.authService.Authenticate(token)
	if err != nil {
		return err
	}
	if role != "member" {
		return fmt.Errorf("unauthorized: only members can return books")
	}

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

func (s *borrowingService) GetMyBorrowings(memberID uint, token string) ([]models.BorrowingRecord, error) {
	_, role, err := s.authService.Authenticate(token)
	if err != nil {
		return nil, err
	}
	if role != "member" {
		return nil, fmt.Errorf("unauthorized: only members can view their borrowings")
	}

	borrowingRecords, err := s.borrowingRepo.GetByMemberID(memberID)
	if err != nil {
		return nil, err
	}

	return borrowingRecords, nil
}

func (s *borrowingService) GetAllBorrowingRecords(token string) ([]models.BorrowingRecord, error) {
	_, role, err := s.authService.Authenticate(token)
	if err != nil {
		return nil, err
	}
	if role != "admin" && role != "librarian" {
		return nil, fmt.Errorf("unauthorized: only admins and librarians can view all borrowing records")
	}

	borrowingRecords, err := s.borrowingRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return borrowingRecords, nil
}
