package models

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	Name   string `gorm:"not null" gorm:"size:64"`
	Users  []User `gorm:"many2many:users_groups;"`
	Author User   `gorm:"foreignKey:AuthorID"`
	AuthorID uint
}
