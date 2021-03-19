package models

import "gorm.io/gorm"

type Discussion struct {
	gorm.Model
	Author User   `gorm:"foreignKey:AuthorID"`
	AuthorID *uint
	Topic Topic   `gorm:"foreignKey:TopicID"`
	TopicID *uint
}