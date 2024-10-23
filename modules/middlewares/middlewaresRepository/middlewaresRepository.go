package middlewaresrepository

import "gorm.io/gorm"

type IMiddlewaresRepository interface {
}

type middlewaresRepository struct {
	db *gorm.DB
}

func MiddlewaresRepository(db *gorm.DB) IMiddlewaresRepository {
	return &middlewaresRepository{db: db}
}
