package models

import "time"

type BorrowingRecord struct {
	ID         uint `gorm:"primaryKey"`
	BookID     uint
	Book       Book `gorm:"foreignKey:BookID"`
	MemberID   uint
	BorrowDate time.Time
	ReturnDate time.Time `gorm:"default:null"`
}
