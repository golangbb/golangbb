package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Content string `gorm:"size:4096"`
	Author User   `gorm:"foreignKey:AuthorID"`
	AuthorID *uint
}
