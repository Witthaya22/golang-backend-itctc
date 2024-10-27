package userrepository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Witthaya22/golang-backend-itctc/entities"
	"gorm.io/gorm"
)

type IUserRepository interface {
	RegisterUser(user *entities.User) error
	FindByUserID(userID string) (*entities.User, error)
	FindDepartmentByID(departmentID string) (*entities.Department, error)
	FindOneUserByUserID(userID string) (*entities.UserCredentialCheck, error)
	InsertOauthUser(oauth *entities.Oauth) error
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

func (r *userRepository) FindOneUserByUserID(userID string) (*entities.UserCredentialCheck, error) {
	user := new(entities.UserCredentialCheck)
	err := r.db.Table("users").Where("user_id = ?", userID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found with ID: %s", userID)
		}
		return nil, fmt.Errorf("error finding user: %v", err)
	}
	return user, nil
}

func (r *userRepository) InsertOauthUser(oauth *entities.Oauth) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := r.db.WithContext(ctx).Create(oauth).Error; err != nil {
		return fmt.Errorf("insert oauth failed: %v", err)
	}

	return nil
}
