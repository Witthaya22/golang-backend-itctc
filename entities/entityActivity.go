package entities

import "time"

type Activity struct {
	ActivityID      string    `gorm:"primaryKey;unique;not null"`
	ActivityName    string    `gorm:"type:varchar(128);not null"`
	ActivityDate    time.Time `gorm:"not null"`
	Location        string    `gorm:"type:varchar(128);not null"`
	ActivityMaxOpen int64     `gorm:"not null"`
	ActivityDes     string    `gorm:"type:text"`
	Status          string    `gorm:"type:varchar(20);not null"`
	StartTime       time.Time `gorm:"not null"`
	EndTime         time.Time `gorm:"not null"`
	IsArchived      bool      `gorm:"not null;default:false"`
	CreatedAt       time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt       time.Time `gorm:"not null;autoUpdateTime"`
}
