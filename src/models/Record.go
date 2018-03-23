package models

import (
	"errors"
	"time"

	"github.com/aaa59891/account_backend/src/db"

	"github.com/jinzhu/gorm"
)

var (
	ErrNoCategory = errors.New("there is no category")
)

type Record struct {
	Id         uint      `gorm:"primary_key" json:"id"`
	Email      string    `gorm:"type:varchar(100);not null;index;" json:"email"`
	Amount     int       `json:"amount"`
	Category   Category  `gorm:"foreignkey:CategoryId" json:"category"`
	CategoryId uint      `json:"categoryId"`
	Date       time.Time `gorm:"type:date;index" json:"date" time_format:"2006-01-02"`
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

type RecordForm struct {
	Email string    `form:"email"`
	Date  time.Time `form:"date" time_format:"2006-01-02"`
	Mode  string    `form:"mode"`
}

const (
	RECORD_FETCH_MODE_ALL   = "All"
	RECORD_FETCH_MODE_DATE  = "Date"
	RECORD_FETCH_MODE_MONTH = "Month"
)

func (rf RecordForm) GetRecords() (data []Record, err error) {
	d := db.DB.Preload("Category").Where("email = ?", rf.Email)
	if rf.Mode != RECORD_FETCH_MODE_ALL {
		if rf.Date.Year() != 1 {
			switch rf.Mode {
			case RECORD_FETCH_MODE_MONTH:
				d = d.Where("MONTH(`date`) = ?", rf.Date.Month())
			case RECORD_FETCH_MODE_DATE:
				d = d.Where("`date` = ?", rf.Date)
			}
		}
	}
	err = d.Find(&data).Error
	return
}
