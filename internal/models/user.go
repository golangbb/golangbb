package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name        string
	DisplayName string
	Email       string `gorm:"uniqueIndex" gorm:"not null"`
	Password    string `gorm:"not null" gorm:"size:64"`
}
