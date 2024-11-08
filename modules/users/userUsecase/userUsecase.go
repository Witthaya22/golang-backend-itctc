package userusecase

import (
	"errors"
	"fmt"
	"log"
	"time"

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
	RefreshPassport(req *entities.UserRefresnCredential) (*entities.UserPassport, error)
	DeleteOauth(oauthId string) error
	AddAdminRole(req *entities.User) error
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

func (u *userUsecase) RefreshPassport(req *entities.UserRefresnCredential) (*entities.UserPassport, error) {
	// 1. เพิ่มการตรวจสอบ refresh token
	if req.RefreshToken == "" {
		return nil, fmt.Errorf("refresh token is required")
	}
	log.Printf("Processing refresh token request with token: %s", req.RefreshToken)

	// 2. ตรวจสอบและ log claims
	claims, err := auth.ParseToken(req.RefreshToken, *u.conf.Jwt)
	if err != nil {
		log.Printf("Failed to parse token: %v", err)
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}
	log.Printf("Successfully parsed claims: %+v", claims)

	// 3. ตรวจสอบ token expiration
	if time.Now().Unix() > claims.ExpiresAt.Unix() {
		return nil, fmt.Errorf("refresh token has expired")
	}

	// 4. ปรับปรุงการค้นหา oauth
	oauth, err := u.userRepository.FindOneOauth(req.RefreshToken)
	if err != nil {
		log.Printf("Failed to find oauth: %v", err)
		return nil, fmt.Errorf("failed to find oauth: %v", err)
	}
	log.Printf("Found OAuth record: %+v", oauth)

	// 5. เพิ่มการตรวจสอบข้อมูล oauth
	if oauth.ID == "" || oauth.UserID == "" {
		return nil, fmt.Errorf("invalid oauth record")
	}

	// 6. ปรับปรุงการดึงข้อมูล profile
	profile, err := u.userRepository.GetProfile(oauth.UserID)
	if err != nil {
		log.Printf("Failed to find profile for user ID %s: %v", oauth.UserID, err)
		return nil, fmt.Errorf("failed to find profile: %v", err)
	}
	log.Printf("Found user profile: %+v", profile)

	// 7. สร้าง claims ใหม่
	newClaims := &entities.UserClaims{
		UserID: profile.UserID,
		Role:   profile.Role,
	}

	// 8. สร้าง tokens ใหม่พร้อมตรวจสอบ error
	accessToken, err := auth.NewAuth(auth.Access, *u.conf.Jwt, newClaims)
	if err != nil {
		log.Printf("Failed to create new access token: %v", err)
		return nil, fmt.Errorf("error creating access token: %v", err)
	}

	// 9. สร้าง refresh token ใหม่
	newRefreshToken := auth.RepeatToken(*u.conf.Jwt, newClaims, claims.ExpiresAt.Unix())
	if newRefreshToken == "" {
		return nil, fmt.Errorf("failed to create new refresh token")
	}

	// 10. สร้าง passport
	passport := &entities.UserPassport{
		User: &entities.UserResponse{
			UserID:        profile.UserID,
			UserFirstName: profile.UserFirstName,
			UserLastName:  profile.UserLastName,
			Role:          profile.Role,
		},
		Token: &entities.UserToken{
			ID:           oauth.ID,
			AccessToken:  accessToken.SingToken(),
			RefreshToken: newRefreshToken,
		},
	}

	// 11. อัพเดท oauth พร้อม logging
	log.Printf("Updating oauth record for user: %s", profile.UserID)
	if err := u.userRepository.UpdateOauth(passport.Token); err != nil {
		log.Printf("Failed to update oauth: %v", err)
		return nil, fmt.Errorf("failed to update oauth: %v", err)
	}
	log.Printf("Successfully updated oauth")

	return passport, nil
}

func (u *userUsecase) DeleteOauth(oauthId string) error {
	if err := u.userRepository.DeleteOauth(oauthId); err != nil {
		log.Printf("Failed to delete oauth: %v", err)
		return fmt.Errorf("failed to delete oauth: %v", err)
	}
	return nil
}

func (u *userUsecase) AddAdminRole(req *entities.User) error {
	// Validate request
	if req.UserID == "" {
		return errors.New("user ID is required")
	}

	// Add admin role to user
	if err := u.userRepository.AddAdminRole(req.UserID); err != nil {
		return fmt.Errorf("failed to add admin role: %v", err)
	}

	return nil
}
