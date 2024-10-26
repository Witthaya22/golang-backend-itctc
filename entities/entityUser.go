package entities

import (
	"fmt"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserID        string    `gorm:"primaryKey;unique;not null" json:"userID"`
	UserFirstName string    `gorm:"type:varchar(64);not null" json:"userFirstName"`
	UserLastName  string    `gorm:"type:varchar(64);not null" json:"userLastName"`
	UserPassword  string    `gorm:"type:varchar(256);not null" json:"userPassword"`
	UserImage     string    `gorm:"type:varchar(256)" json:"userImage"`
	DepartmentID  string    `gorm:"type:varchar(64);not null" json:"departmentID"`
	Role          string    `gorm:"type:varchar(20);not null" json:"role"`
	AccessToken   string    `gorm:"type:varchar(256);" json:"accessToken"`
	RefreshToken  string    `gorm:"type:varchar(256);" json:"refreshToken"`
	IsArchived    bool      `gorm:"not null;default:false" json:"isArchived"`
	CreatedAt     time.Time `gorm:"not null;autoCreateTime" json:"createdAt"`
	UpdatedAt     time.Time `gorm:"not null;autoUpdateTime" json:"updatedAt"`
}

type UserRegisterReq struct {
	UserID        string `json:"userID" form:"userID"`
	UserPassword  string `json:"userPassword" form:"userPassword"`
	DepartmentID  string `json:"departmentID" form:"departmentID"`
	UserFirstName string `json:"userFirstName" form:"userFirstName"`
	UserLastName  string `json:"userLastName" form:"userLastName" `
}

func (obj *UserRegisterReq) BcryptHashing() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(obj.UserPassword), 10)
	if err != nil {
		return fmt.Errorf("error hashing password: %v", err)
	}
	obj.UserPassword = string(hashedPassword)
	return nil
}

func (obj *UserRegisterReq) IsEmail() bool {
	match, err := regexp.MatchString(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`, obj.UserID)
	if err != nil {
		return false
	}

	return match
}

type UserResponse struct {
	UserID        string `json:"userID"`
	UserFirstName string `json:"userFirstName"`
	UserLastName  string `json:"userLastName"`
	Role          string `json:"role"`
}

type UserCredential struct {
	UserID       string `json:"userID" form:"userID"`
	UserPassword string `json:"userPassword" form:"userPassword"`
}

type UserCredentialCheck struct {
	UserID        string `json:"userID" form:"userID"`
	UserPassword  string `json:"userPassword" form:"userPassword"`
	UserFirstName string `json:"userFirstName" form:"userFirstName"`
	UserLastName  string `json:"userLastName" form:"userLastName" `
	Role          string `json:"role" form:"role"`
}

type UserPassport struct {
	User  *UserResponse `json:"user"`
	Token *UserToken    `json:"token"`
}

type UserToken struct {
	UserID       string `json:"userID"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
