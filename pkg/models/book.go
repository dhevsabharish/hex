package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title           string         `gorm:"not null"`
	Author          string         `gorm:"not null"`
	PublicationDate datatypes.Date `gorm:"type:date;not null"`
	Genre           string
	Availability    uint `gorm:"default:1"`
}

func (b *Book) BeforeCreate(tx *gorm.DB) error {
	if b.PublicationDate == (datatypes.Date{}) {
		b.PublicationDate = datatypes.Date(time.Now())
	}
	return nil
}

func (b *Book) BeforeUpdate(tx *gorm.DB) error {
	if b.PublicationDate == (datatypes.Date{}) {
		b.PublicationDate = datatypes.Date(time.Now())
	}
	return nil
}
