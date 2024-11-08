package middlewaresrepository

import (
	"github.com/Witthaya22/golang-backend-itctc/entities"
	"gorm.io/gorm"
)

type IMiddlewaresRepository interface {
	FindAccessToken(userId, accessToken string) bool
}

type middlewaresRepository struct {
	db *gorm.DB
}

func MiddlewaresRepository(db *gorm.DB) IMiddlewaresRepository {
	return &middlewaresRepository{db: db}
}

func (r *middlewaresRepository) FindAccessToken(userId, accessToken string) bool {
	var oauth entities.Oauth
	err := r.db.Where("user_id = ? AND access_token = ?", userId, accessToken).First(&oauth).Error
	return err == nil
}
