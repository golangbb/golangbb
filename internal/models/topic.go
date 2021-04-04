package models

import (
	"github.com/golangbb/golangbb/v2/internal/database"
	"gorm.io/gorm"
	"log"
)

type Topic struct {
	gorm.Model
	Title    string `gorm:"uniqueIndex" gorm:"not null" gorm:"size:96"`
	ParentID *uint  `gorm:"TYPE:integer REFERENCES topics"`
	Parent   *Topic
	Author   User `gorm:"foreignKey:AuthorID"`
	AuthorID uint
}

func CreateTopic(topic *Topic) error {
	if topic.Title == "" {
		return ErrEmptyTitle
	}

	if topic.AuthorID == 0 {
		return ErrEmptyUserID
	}

	err := database.DBConnection.Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("Parent", "Author").Create(topic).Error; err != nil {
			log.Println("[CREATE_TOPIC]::DB_INSERT_TOPIC_ERROR 💥")
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
