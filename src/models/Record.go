package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

var (
	ErrNoCategory = errors.New("there is no category")
)

type Record struct {
	Id         uint   `gorm:"primary_key"`
	Email      string `gorm:"type:varchar(100);not null;index;"`
	Amount     int
	CategoryId uint
	Date       time.Time `gorm:"type:date();index"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (record *Record) Insert(tx *gorm.DB) error {
	if len(record.Email) == 0 {
		return ErrNoEmail
	}
	if record.CategoryId == 0 {
		return ErrNoCategory
	}
	return tx.Create(record).Error
}

func (record *Record) Update(tx *gorm.DB) error {
	return tx.Model(record).Update(Record{CategoryId: record.CategoryId, Amount: record.Amount}).Error
}

func (record *Record) Delete(tx *gorm.DB) error {
	return tx.Delete(record).Error
}
