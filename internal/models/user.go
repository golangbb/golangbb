package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName    string `gorm:"uniqueIndex" gorm:"not null" gorm:"size:32"`
	DisplayName string
	Password    string  `gorm:"not null" gorm:"size:64"`
	Emails      []Email
	Groups      []Group `gorm:"many2many:users_groups;"`
}
