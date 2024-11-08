package userrepository

import (
	"context"
	"errors"
	"fmt"
	"log"
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
	FindOneOauth(refreshToken string) (*entities.FindOneOauth, error)
	UpdateOauth(req *entities.UserToken) error
	// GetProfile(userID string) (*entities.User, error)
	GetProfile(userID string) (*entities.UserProfile, error)
	DeleteOauth(oauthId string) error
	AddAdminRole(userID string) error
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

func (r *userRepository) FindOneOauth(refreshToken string) (*entities.FindOneOauth, error) {
	var oauth entities.FindOneOauth

	err := r.db.Debug().
		Table("oauths").
		Select("id, user_id, access_token, refresh_token").
		Where("refresh_token = ?", refreshToken).
		First(&oauth).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("oauth not found with refresh token")
		}
		return nil, fmt.Errorf("error finding oauth: %v", err)
	}

	// เพิ่ม logging เพื่อตรวจสอบข้อมูลที่ได้
	log.Printf("Found OAuth: %+v", oauth)

	return &oauth, nil
}

// func (r *userRepository) UpdateOauth(req *entities.UserToken) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	if err := r.db.WithContext(ctx).Model(&entities.User{}).Where("user_id = ?", req.ID).Updates(map[string]interface{}{
// 		"access_token":  req.AccessToken,
// 		"refresh_token": req.RefreshToken,
// 	}).Error; err != nil {
// 		return fmt.Errorf("update user failed: %v", err)
// 	}

// 	return nil
// }

func (r *userRepository) UpdateOauth(token *entities.UserToken) error {
	// แก้จาก
	// err := r.db.Model(&entities.User{}).Where("user_id = ?", token.ID).Updates(map[string]interface{}{
	//     "access_token":  token.AccessToken,
	//     "refresh_token": token.RefreshToken,
	// }).Error

	// เป็น
	err := r.db.Model(&entities.Oauth{}).Where("id = ?", token.ID).Updates(map[string]interface{}{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	}).Error

	if err != nil {
		return fmt.Errorf("update oauth failed: %v", err)
	}
	return nil
}

// func (r *userRepository) GetProfile(userID string) (*entities.User, error) {

// 	// get from user table userid, userfirstname, userlastname, role

// 	user := new(entities.User)
// 	err := r.db.Table("users").Where("user_id = ?", userID).First(&map[string]interface{}{
// 		"user_id":         userID,
// 		"user_first_name": user.UserFirstName,
// 		"user_last_name":  user.UserLastName,
// 		"role":            user.Role,
// 	}).Error
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, fmt.Errorf("user not found with ID: %s", userID)
// 		}
// 		return nil, fmt.Errorf("error finding user: %v", err)
// 	}

// 	return user, nil

// }

func (r *userRepository) GetProfile(userID string) (*entities.UserProfile, error) {
	var user entities.UserProfile
	result := r.db.Debug().Where("user_id = ?", userID).First(&user)
	if result.Error != nil {
		log.Printf("Error querying user: %v", result.Error)
		return nil, fmt.Errorf("error finding user: %v", result.Error)
	}
	log.Printf("Found user: %+v", user)
	return &user, nil
}

func (r *userRepository) DeleteOauth(oauthId string) error {

	if err := r.db.Debug().Where("id = ?", oauthId).Delete(&entities.Oauth{}).Error; err != nil {
		return fmt.Errorf("delete oauth failed: %v", err)
	}

	return nil
}

func (r *userRepository) AddAdminRole(userID string) error {
	// Update user role to admin
	return r.db.Model(&entities.User{}).Where("user_id = ?", userID).Update("role", "admin").Error
}
