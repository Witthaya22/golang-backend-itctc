package entities

import "time"

type Activity struct {
	ActivityID      string `gorm:"primaryKey;unique;not null;autoIncrement"`
	Title           string `gorm:"type:varchar(128);not null"`
	ActivityDate    time.Time
	Location        string `gorm:"type:varchar(128);not null"`
	ActivityMaxOpen int64
	Description     string  `gorm:"type:text"`
	Score           float32 `gorm:"type:float"`
	Status          string  `gorm:"type:varchar(20);"`
	Images          string  `gorm:"type:varchar(512);"`
	StartTime       time.Time
	EndTime         time.Time
	IsArchived      bool      `gorm:"not null;default:false"`
	CreatedAt       time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt       time.Time `gorm:"not null;autoUpdateTime"`
}
