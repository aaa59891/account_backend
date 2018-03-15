package models

import (
	"errors"
	"strings"
	"time"

	"github.com/aaa59891/account_backend/src/db"

	"github.com/aaa59891/account_backend/src/constants"

	"golang.org/x/crypto/bcrypt"

	"github.com/aaa59891/account_backend/src/configs"

	"github.com/jinzhu/gorm"
)

var (
	ErrNoEmail       = errors.New("this email does not exist")
	ErrWrongPassword = errors.New("wrong password")
	ErrEmailExist    = errors.New("this email already existed")
)

type User struct {
	Id        uint   `gorm:"primary_key"`
	Email     string `gorm:"type:varchar(100);unique_index;not null" binding:"required"`
	Password  string `gorm:"not null;" binding:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (u *User) Insert(tx *gorm.DB) error {
	u.Email = strings.ToLower(u.Email)
	if err := u.encryptPassword(); err != nil {
		return err
	}
	if err := tx.Create(u).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			return ErrEmailExist
		}
		return err
	}
	return nil
}

func (u *User) Update(tx *gorm.DB) error {
	return constants.ErrNoUsedFunction
}

func (u *User) Delete(tx *gorm.DB) error {
	return constants.ErrNoUsedFunction
}

func (u *User) UpdatePassword(tx *gorm.DB) error {
	if err := u.CheckPassword(); err != nil {
		return err
	}
	if err := u.encryptPassword(); err != nil {
		return err
	}
	return tx.Model(u).Where("email = ?", strings.ToLower(u.Email)).Update("password", u.Password).Error
}

func (u *User) CheckPassword() error {
	dbUser := User{}
	if err := db.DB.Find(&dbUser, "email = ?", strings.ToLower(u.Email)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrNoEmail
		}
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(u.Password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return ErrWrongPassword
		}
		return err
	}
	return nil
}

func (u *User) encryptPassword() error {
	bcryptCost := configs.GetConfig().Security.BcryptCost
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcryptCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}
