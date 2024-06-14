package persistence

import (
	"hex/pkg/models"

	"gorm.io/gorm"
)

type BorrowingRepository struct {
	db *gorm.DB
}

func NewBorrowingRepository(db *gorm.DB) *BorrowingRepository {
	return &BorrowingRepository{db: db}
}

func (r *BorrowingRepository) Borrow(record *models.BorrowingRecord) error {
	return r.db.Create(record).Error
}
func (r *BorrowingRepository) Return(recordID string) error {
	return r.db.Model(&models.BorrowingRecord{}).Where("id = ?", recordID).Update("returned", true).Error
}

func (r *BorrowingRepository) GetAllRecords() ([]models.BorrowingRecord, error) {
	var records []models.BorrowingRecord
	err := r.db.Find(&records).Error
	return records, err
}

func (r *BorrowingRepository) GetRecordsByUserID(userID string) ([]models.BorrowingRecord, error) {
	var records []models.BorrowingRecord
	err := r.db.Where("user_id = ?", userID).Find(&records).Error
	return records, err
}
