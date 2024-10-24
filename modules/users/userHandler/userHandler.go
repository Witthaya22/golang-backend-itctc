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
	signUpUserErrCode userHandlerErrCode = "user-001"
)

type IUserHandler interface {
	SignUpUser(c *fiber.Ctx) error
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
