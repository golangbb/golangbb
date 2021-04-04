package models

import (
	"github.com/golangbb/golangbb/v2/internal/database"
	"gorm.io/gorm"
	"log"
)

type Post struct {
	gorm.Model
	Content      string     `gorm:"size:4096" gorm:"not null"`
	Author       User       `gorm:"foreignKey:AuthorID"`
	AuthorID     uint       `gorm:"not null"`
	Discussion   Discussion `gorm:"foreignKey:DiscussionID"`
	DiscussionID uint       `gorm:"not null"`
}

func CreatePost(post *Post) error {
	if post.Content == "" {
		return ErrEmptyContent
	}

	if post.AuthorID == 0 {
		return ErrEmptyUserID
	}

	if post.DiscussionID == 0 {
		return ErrEmptyDiscussionID
	}

	err := database.DBConnection.Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("Author", "Discussion").Create(post).Error; err != nil {
			log.Println("[CREATE_POST]::DB_INSERT_POST_ERROR 💥")
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}