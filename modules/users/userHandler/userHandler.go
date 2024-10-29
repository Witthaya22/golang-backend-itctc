package userhandler

import (
	"strings"

	"github.com/Witthaya22/golang-backend-itctc/config"
	"github.com/Witthaya22/golang-backend-itctc/entities"
	customresponse "github.com/Witthaya22/golang-backend-itctc/modules/customResponse"
	userusecase "github.com/Witthaya22/golang-backend-itctc/modules/users/userUsecase"
	"github.com/gofiber/fiber/v2"
)

type userHandlerErrCode string

const (
	signUpUserErrCode      userHandlerErrCode = "user-001"
	signInErrCode          userHandlerErrCode = "user-002"
	refreshPassportErrCode userHandlerErrCode = "user-003"
	signOutErrCode         userHandlerErrCode = "user-004"
)

type IUserHandler interface {
	SignUpUser(c *fiber.Ctx) error
	SignIn(c *fiber.Ctx) error
	RefreshPassport(c *fiber.Ctx) error
	SignOut(c *fiber.Ctx) error
}

type userHandler struct {
	conf        *config.Config
	userUsecase userusecase.IUserUsecase
}

func UserHandler(conf *config.Config, userUsecase userusecase.IUserUsecase) IUserHandler {
	return &userHandler{
		conf:        conf,
		userUsecase: userUsecase,
	}
}

func (h *userHandler) SignUpUser(c *fiber.Ctx) error {
	// Parse request
	req := new(entities.UserRegisterReq)
	if err := c.BodyParser(req); err != nil {
		return customresponse.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpUserErrCode),
			"invalid request format",
		).Res()
	}

	// Validate email format
	if !req.IsEmail() {
		return customresponse.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpUserErrCode),
			"invalid email format",
		).Res()
	}

	// Validate required fields
	if req.DepartmentID == "" || req.UserFirstName == "" || req.UserLastName == "" || req.UserPassword == "" {
		return customresponse.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpUserErrCode),
			"missing required fields",
		).Res()
	}

	// Register user
	result, err := h.userUsecase.RegisterUser(req)
	if err != nil {
		// Determine appropriate error code based on error type
		statusCode := fiber.ErrInternalServerError.Code
		if strings.Contains(err.Error(), "department validation failed") {
			statusCode = fiber.ErrBadRequest.Code
		}

		return customresponse.NewResponse(c).Error(
			statusCode,
			string(signUpUserErrCode),
			err.Error(),
		).Res()
	}

	return customresponse.NewResponse(c).Success(
		fiber.StatusCreated,
		result,
	).Res()
}

func (h *userHandler) SignIn(c *fiber.Ctx) error {
	req := new(entities.UserCredential)
	if err := c.BodyParser(req); err != nil {
		return customresponse.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signInErrCode),
			err.Error(),
		).Res()
	}

	passport, err := h.userUsecase.GetPassport(req)
	if err != nil {
		return customresponse.NewResponse(c).Error(
			fiber.ErrNotFound.Code,
			string(signInErrCode),
			err.Error(),
		).Res()
	}

	return customresponse.NewResponse(c).Success(
		fiber.StatusOK,
		passport,
	).Res()
}

func (h *userHandler) RefreshPassport(c *fiber.Ctx) error {
	req := new(entities.UserRefresnCredential)
	if err := c.BodyParser(req); err != nil {
		return customresponse.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(refreshPassportErrCode),
			err.Error(),
		).Res()
	}

	passport, err := h.userUsecase.RefreshPassport(req)
	if err != nil {
		return customresponse.NewResponse(c).Error(
			fiber.ErrNotFound.Code,
			string(refreshPassportErrCode),
			err.Error(),
		).Res()
	}

	return customresponse.NewResponse(c).Success(
		fiber.StatusOK,
		passport,
	).Res()
}

func (h *userHandler) SignOut(c *fiber.Ctx) error {
	req := new(entities.UserRemoveCredentials)
	if err := c.BodyParser(req); err != nil {
		return customresponse.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signOutErrCode),
			err.Error(),
		).Res()
	}

	if err := h.userUsecase.DeleteOauth(req.Id); err != nil {
		return customresponse.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signOutErrCode),
			err.Error(),
		).Res()
	}

	return customresponse.NewResponse(c).Success(
		fiber.StatusOK,
		nil,
	).Res()
}
