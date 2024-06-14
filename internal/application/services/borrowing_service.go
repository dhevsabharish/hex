package services

import (
	"hex/internal/adapters/persistence"
	"hex/pkg/models"
)

type BorrowingService struct {
	repo persistence.BorrowingRepository
}

func NewBorrowingService(repo persistence.BorrowingRepository) *BorrowingService {
	return &BorrowingService{repo: repo}
}

func (s *BorrowingService) BorrowBook(record *models.BorrowingRecord) error {
	return s.repo.Borrow(record)
}

func (s *BorrowingService) ReturnBook(recordID string) error {
	return s.repo.Return(recordID)
}

func (s *BorrowingService) ShowAllRecords() ([]models.BorrowingRecord, error) {
	return s.repo.GetAllRecords()
}

func (s *BorrowingService) ShowRecordsByUserID(userID string) ([]models.BorrowingRecord, error) {
	return s.repo.GetRecordsByUserID(userID)
}
