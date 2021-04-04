package models

import (
	"github.com/golangbb/golangbb/v2/internal/database"
	"gorm.io/gorm"
	"log"
)

type Discussion struct {
	gorm.Model
	Title    string `gorm:"not null" gorm:"size:128"`
	Author   User   `gorm:"foreignKey:AuthorID"`
	AuthorID uint   `gorm:"not null"`
	Topic    Topic  `gorm:"foreignKey:TopicID"`
	TopicID  uint   `gorm:"not null"`
	Posts    []Post
}

func CreateDiscussion(discussion *Discussion) error {
	if discussion.Title == "" {
		return ErrEmptyTitle
	}

	if discussion.AuthorID == 0 {
		return ErrEmptyUserID
	}

	if discussion.TopicID == 0 {
		return ErrEmptyTopicID
	}

	if len(discussion.Posts) != 1 {
		return ErrDiscussionWithoutSinglePost
	}

	err := database.DBConnection.Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("Author", "Topic", "Posts").Create(discussion).Error; err != nil {
			log.Println("[CREATE_DISCUSSION]::DB_INSERT_DISCUSSION_ERROR 💥")
			return err
		}

		for i, _ := range discussion.Posts {
			discussion.Posts[i].DiscussionID = discussion.ID
			discussion.Posts[i].AuthorID = discussion.AuthorID
		}

		if err := tx.Omit("Author", "Discussion").CreateInBatches(&discussion.Posts, 10).Error; err != nil {
			log.Println("[CREATE_DISCUSSION]::DB_INSERT_POST_ERROR 💥")
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
