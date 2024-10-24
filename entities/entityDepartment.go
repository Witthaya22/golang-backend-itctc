package entities

import "time"

// type Department struct {
// 	DepartmentID   string    `gorm:"primaryKey;unique;not null" json:"id"`
// 	DepartmentName string    `gorm:"type:varchar(64);not null" json:"name"`
// 	Users          []User    `gorm:"foreignKey:DepartmentID"`
// 	IsArchived     bool      `gorm:"not null;default:false"`
// 	CreatedAt      time.Time `gorm:"not null;autoCreateTime" json:"createdAt`
// 	UpdatedAt      time.Time `gorm:"not null;autoUpdateTime" json:"updatedAt"`
// }

type Department struct {
	DepartmentID   string    `gorm:"primaryKey;column:department_id" json:"departmentID"`
	DepartmentName string    `gorm:"type:varchar(64);not null" json:"departmentName"`
	IsArchived     bool      `gorm:"not null;default:false" json:"isArchived"`
	CreatedAt      time.Time `gorm:"not null;autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"not null;autoUpdateTime" json:"updatedAt"`
}
