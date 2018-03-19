package models

import (
	"time"

	"github.com/aaa59891/account_backend/src/db"

	"github.com/jinzhu/gorm"
)

type Category struct {
	Id        uint      `gorm:"primary_key" json:"id"`
	Email     string    `gorm:"type:varchar(100);not null" json:"email"`
	Name      string    `gorm:"not null;" json:"name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (category *Category) Insert(tx *gorm.DB) error {
	if len(category.Email) == 0 {
		return ErrNoEmail
	}
	return tx.Create(category).Error
}

func (category *Category) Update(tx *gorm.DB) error {
	return tx.Model(category).Update("name", category.Name).Error
}

func (category *Category) Delete(tx *gorm.DB) error {
	return tx.Delete(category).Error
}

func GetCategoriesByEmail(email string) (categories []Category, err error) {
	err = db.DB.Find(&categories, "email = ?", email).Error
	return
}
