package models

import (
	"gorm.io/gorm"
	"time"
)

type Email struct {
	Email     string `gorm:"primaryKey" gorm:"size:128"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	User      User           `gorm:"foreignkey:UserID"`
	UserID    *uint
}
