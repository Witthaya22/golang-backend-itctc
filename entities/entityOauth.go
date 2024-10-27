package entities

import "time"

type Oauth struct {
	ID           string    `gorm:"primaryKey;"`
	UserID       string    `gorm:"type:varchar(64);not null"`
	AccessToken  string    `gorm:"type:varchar(512);" json:"accessToken"`
	RefreshToken string    `gorm:"type:varchar(512);" json:"refreshToken"`
	IsArchived   bool      `gorm:"not null;default:false"`
	CreatedAt    time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"not null;autoUpdateTime"`
	// User         User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user"`
}
