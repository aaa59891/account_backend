package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

var (
	ErrNoId = errors.New("ID is not specified")
)

type Dao interface {
	Insert(tx *gorm.DB) error
	Update(tx *gorm.DB) error
	Delete(tx *gorm.DB) error
}

/* Special -- not use dao interface */
func DeleteById(model, id interface{}, user, from string) Transaction {
	return func(tx *gorm.DB) error {
		if tx.NewRecord(model) {
			return ErrNoId
		}

		if err := tx.Delete(model, "id = ?", id).Error; err != nil {
			return err
		}
		return nil
	}
}
