package entities

import (
	"time"
)

type User struct {
	UserID        string    `gorm:"primaryKey;unique;not null"`
	UserFirstName string    `gorm:"type:varchar(64);not null"`
	UserLastName  string    `gorm:"type:varchar(64);not null"`
	UserPassword  string    `gorm:"type:varchar(256);not null"`
	UserImage     string    `gorm:"type:varchar(256)"`
	DepartmentID  string    `gorm:"type:varchar(64);not null"`
	Role          string    `gorm:"type:varchar(20);not null"`
	IsArchived    bool      `gorm:"not null;default:false"`
	CreatedAt     time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt     time.Time `gorm:"not null;autoUpdateTime"`
}
