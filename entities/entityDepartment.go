package entities

import "time"

type Department struct {
	DepartmentID   string    `gorm:"primaryKey;unique;not null"`
	DepartmentName string    `gorm:"type:varchar(64);not null"`
	Users          []User    `gorm:"foreignKey:DepartmentID"`
	IsArchived     bool      `gorm:"not null;default:false"`
	CreatedAt      time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt      time.Time `gorm:"not null;autoUpdateTime"`
}
