package models

import (
	"gorm.io/gorm"
)

type Topic struct {
	gorm.Model
	Title string `gorm:"uniqueIndex" gorm:"not null" gorm:"size:96"`
	ParentID *uint `gorm:"TYPE:integer REFERENCES topics"`
	Parent *Topic
}