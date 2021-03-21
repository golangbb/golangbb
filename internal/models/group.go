package models

import (
	"github.com/golangbb/golangbb/v2/internal/database"
	"gorm.io/gorm"
	"log"
)

type Group struct {
	gorm.Model
	Name   string `gorm:"not null" gorm:"size:64"`
	Users  []User `gorm:"many2many:users_groups;"`
	Author User   `gorm:"foreignKey:AuthorID"`
	AuthorID uint
}

func CreateGroup(group *Group) error {
	if group.Name == "" {
		return ErrEmptyName
	}

	if group.AuthorID == 0 {
		return ErrEmptyUserID
	}

	err := database.DBConnection.Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("Users", "Author").Create(group).Error; err != nil {
			log.Println("[CREATE_GROUP]::DB_INSERT_GROUP_ERROR 💥")
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}