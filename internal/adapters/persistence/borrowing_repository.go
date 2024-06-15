package persistence

import (
	"hex/pkg/models"

	"gorm.io/gorm"
)

type BorrowingRepository struct {
	DB *gorm.DB
}

func NewBorrowingRepository(db *gorm.DB) *BorrowingRepository {
	return &BorrowingRepository{DB: db}
}

func (r *BorrowingRepository) Create(borrowingRecord *models.BorrowingRecord) error {
	return r.DB.Create(borrowingRecord).Error
}

func (r *BorrowingRepository) GetByID(id uint) (*models.BorrowingRecord, error) {
	var borrowingRecord models.BorrowingRecord
	err := r.DB.Preload("Book").First(&borrowingRecord, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &borrowingRecord, nil
}

func (r *BorrowingRepository) GetByMemberID(memberID uint) ([]models.BorrowingRecord, error) {
	var borrowingRecords []models.BorrowingRecord
	err := r.DB.Where("member_id = ?", memberID).Preload("Book").Find(&borrowingRecords).Error
	return borrowingRecords, err
}

func (r *BorrowingRepository) GetAll() ([]models.BorrowingRecord, error) {
	var borrowingRecords []models.BorrowingRecord
	err := r.DB.Preload("Book").Find(&borrowingRecords).Error
	return borrowingRecords, err
}

func (r *BorrowingRepository) Update(borrowingRecord *models.BorrowingRecord) error {
	return r.DB.Save(borrowingRecord).Error
}
