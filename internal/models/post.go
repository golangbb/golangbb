package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Content      string     `gorm:"size:4096" gorm:"not null"`
	Author       User       `gorm:"foreignKey:AuthorID"`
	AuthorID     uint       `gorm:"not null"`
	Discussion   Discussion `gorm:"foreignKey:DiscussionID"`
	DiscussionID uint       `gorm:"not null"`
}
