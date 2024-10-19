package databases

import "gorm.io/gorm"

type IDatabase interface {
	ConnectionGetting() *gorm.DB
}
