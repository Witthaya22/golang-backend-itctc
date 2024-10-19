package entities

import (
	"time"
)

type Admin struct {
	AdminID    string    `gorm:"primaryKey;unique;not null"`
	AdminPass  string    `gorm:"type:varchar(256);not null"`
	UserID     string    `gorm:"type:varchar(64);unique;not null"`
	User       User      `gorm:"foreignKey:UserID;references:UserID"`
	IsArchived bool      `gorm:"not null;default:false"`
	CreatedAt  time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt  time.Time `gorm:"not null;autoUpdateTime"`
}
