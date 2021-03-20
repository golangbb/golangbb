package models

import (
	"errors"
	"github.com/golangbb/golangbb/v2/internal/database"
	"gorm.io/gorm"
	"log"
)

type User struct {
	gorm.Model
	UserName    string `gorm:"uniqueIndex" gorm:"size:32"`
	DisplayName string `gorm:"not null" gorm:"size:32"`
	Password    string `gorm:"not null" gorm:"size:64"`
	Emails      []Email
	Groups      []Group `gorm:"many2many:users_groups;"`
}

var ErrEmptyUserName = errors.New("empty UserName not allowed")
var ErrEmptyPassword = errors.New("empty Password not allowed")

func CreateUser(user *User) error {
	if user.UserName == "" {
		log.Println("[CREATE_USER]::EMPTY_USER_NAME_ERROR 💥")
		return ErrEmptyUserName
	}

	if user.Password == "" {
		log.Println("[CREATE_USER]::EMPTY_PASSWORD_ERROR 💥")
		return ErrEmptyPassword
	}

	if user.DisplayName == "" {
		log.Println("[CREATE_USER]::DEFAULTING_DISPLAY_NAME_WARNING ⚠️")
		user.DisplayName = user.UserName
	}

	err := database.DBConnection.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			log.Println("[CREATE_USER]::DB_INSERT_ERROR 💥")
			log.Println(err)
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
