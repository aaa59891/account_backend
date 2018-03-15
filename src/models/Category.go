package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Category struct {
	Id        uint   `gorm:"primary_key"`
	Email     string `gorm:"type:varchar(100);not null"`
	Name      string `gorm:"not null;"`
	CreatedAt time.Time
	UpdatedAt time.Time
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
