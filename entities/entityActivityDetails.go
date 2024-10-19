package entities

import (
	"time"
)

type ActivityDetails struct {
	ID          string    `gorm:"primaryKey;unique;not null"`
	UserID      string    `gorm:"type:varchar(64);not null"`
	ActivityID  string    `gorm:"type:varchar(64);not null"`
	Reservation bool      `gorm:"not null"`
	JoinSucc    bool      `gorm:"not null"`
	LeaveSucc   bool      `gorm:"not null"`
	Status      string    `gorm:"type:varchar(20);not null"`
	User        User      `gorm:"foreignKey:UserID"`
	Activity    Activity  `gorm:"foreignKey:ActivityID"`
	IsArchived  bool      `gorm:"not null;default:false"`
	CreatedAt   time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"not null;autoUpdateTime"`
}
