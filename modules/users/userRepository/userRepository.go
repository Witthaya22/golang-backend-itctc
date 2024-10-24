package userrepository

import (
	"errors"
	"fmt"

	"github.com/Witthaya22/golang-backend-itctc/entities"
	"gorm.io/gorm"
)

type IUserRepository interface {
	RegisterUser(user *entities.User) error
	FindByUserID(userID string) (*entities.User, error)
	FindDepartmentByID(departmentID string) (*entities.Department, error)
}

type userRepository struct {
	db *gorm.DB
}

func UserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) RegisterUser(user *entities.User) error {
	existingUser := new(entities.User)
	if err := r.db.Where("user_id = ?", user.UserID).First(&existingUser).Error; err == nil {
		return errors.New("user already exists")
	}

	return r.db.Create(user).Error
}

func (r *userRepository) FindByUserID(userID string) (*entities.User, error) {
	user := new(entities.User)
	err := r.db.Where("user_id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) FindDepartmentByID(departmentID string) (*entities.Department, error) {
	department := new(entities.Department)
	err := r.db.Table("departments").Where("department_id = ?", departmentID).First(&department).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("department not found with ID: %s", departmentID)
		}
		return nil, fmt.Errorf("error finding department: %v", err)
	}
	return department, nil
}
