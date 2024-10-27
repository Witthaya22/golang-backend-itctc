package userusecase

import (
	"fmt"

	"github.com/Witthaya22/golang-backend-itctc/config"
	"github.com/Witthaya22/golang-backend-itctc/entities"
	userrepository "github.com/Witthaya22/golang-backend-itctc/modules/users/userRepository"
	"github.com/Witthaya22/golang-backend-itctc/pkg/auth"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecase interface {
	RegisterUser(req *entities.UserRegisterReq) (*entities.User, error)
	GetPassport(req *entities.UserCredential) (*entities.UserPassport, error)
}

type userUsecase struct {
	conf           *config.Config
	userRepository userrepository.IUserRepository
}

func UserUsecase(conf *config.Config, userRepository userrepository.IUserRepository) IUserUsecase {
	return &userUsecase{
		conf:           conf,
		userRepository: userRepository,
	}
}

func (u *userUsecase) RegisterUser(req *entities.UserRegisterReq) (*entities.User, error) {
	// Validate department existence
	department, err := u.userRepository.FindDepartmentByID(req.DepartmentID)
	if err != nil {
		return nil, fmt.Errorf("department validation failed: %v", err)
	}

	// Hash password
	if err := req.BcryptHashing(); err != nil {
		return nil, fmt.Errorf("password hashing failed: %v", err)
	}

	// Create user entity
	user := &entities.User{
		UserID:        req.UserID,
		UserFirstName: req.UserFirstName,
		UserLastName:  req.UserLastName,
		UserPassword:  req.UserPassword,
		DepartmentID:  department.DepartmentID, // Updated to use correct field name
		Role:          "user",
	}

	// Register user
	if err = u.userRepository.RegisterUser(user); err != nil {
		return nil, fmt.Errorf("user registration failed: %v", err)
	}

	return user, nil
}

func (u *userUsecase) GetPassport(req *entities.UserCredential) (*entities.UserPassport, error) {
	user, err := u.userRepository.FindOneUserByUserID(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(req.UserPassword)); err != nil {
		return nil, fmt.Errorf("invalid password: %v", err)
	}

	accessToken, err := auth.NewAuth(auth.Access, *u.conf.Jwt, &entities.UserClaims{
		UserID: user.UserID,
		Role:   user.Role,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating access token: %v", err)
	}

	refreshToken, err := auth.NewAuth(auth.Refresh, *u.conf.Jwt, &entities.UserClaims{
		UserID: user.UserID,
		Role:   user.Role,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating refresh token: %v", err)
	}

	tokenID := uuid.New().String()

	passport := &entities.UserPassport{
		User: &entities.UserResponse{
			UserID:        user.UserID,
			UserFirstName: user.UserFirstName,
			UserLastName:  user.UserLastName,
			Role:          user.Role,
		},
		Token: &entities.UserToken{
			ID:           tokenID,
			AccessToken:  accessToken.SingToken(),
			RefreshToken: refreshToken.SingToken(),
		},
	}

	if err := u.userRepository.InsertOauthUser(&entities.Oauth{
		ID:           tokenID,
		UserID:       user.UserID,
		AccessToken:  accessToken.SingToken(),
		RefreshToken: refreshToken.SingToken(),
	}); err != nil {
		return nil, fmt.Errorf("failed to insert oauth: %v", err)
	}

	return passport, nil
}
