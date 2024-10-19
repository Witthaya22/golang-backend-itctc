package entities

import (
	"time"
)

type ActivityResults struct {
	ID           string    `gorm:"primaryKey;unique;not null"`
	DepartmentID string    `gorm:"type:varchar(64);not null"`
	UserID       string    `gorm:"type:varchar(64);not null"`
	ActivityID   string    `gorm:"type:text;not null"`
	Reservation  bool      `gorm:"not null"`
	Status       string    `gorm:"type:varchar(20);not null"`
	IsArchived   bool      `gorm:"not null;default:false"`
	CreatedAt    time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"not null;autoUpdateTime"`
}
