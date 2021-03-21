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
	Emails      []*Email
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
		if err := tx.Omit("Emails", "Groups").Create(user).Error; err != nil {
			log.Println("[CREATE_USER]::DB_INSERT_USER_ERROR 💥")
			return err
		}

		if len(user.Emails) == 0 {
			return nil
		}

		for _, email := range user.Emails {
			email.UserID = user.ID
		}

		if err := tx.Omit("User").CreateInBatches(&user.Emails, 10).Error; err != nil {
			log.Println("[CREATE_USER]::DB_INSERT_EMAIL_ERROR 💥")
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
