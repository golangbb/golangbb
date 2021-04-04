package models

import (
	"github.com/golangbb/golangbb/v2/internal/database"
	"gorm.io/gorm"
	"log"
	"time"
)

type Email struct {
	Email     string `gorm:"primaryKey" gorm:"size:128"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	User      User           `gorm:"foreignkey:UserID"`
	UserID    uint
}

func CreateEmail(email *Email) error {
	if email.UserID == 0 {
		return ErrEmptyUserID
	}

	err := database.DBConnection.Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("User").Create(email).Error; err != nil {
			log.Println("[CREATE_EMAIL]::DB_INSERT_EMAIL_ERROR 💥")
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
