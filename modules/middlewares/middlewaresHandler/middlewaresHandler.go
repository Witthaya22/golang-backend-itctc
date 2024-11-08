package middlewareshandler

import (
	"strings"

	"github.com/Witthaya22/golang-backend-itctc/config"
	customresponse "github.com/Witthaya22/golang-backend-itctc/modules/customResponse"
	middlewaresusecase "github.com/Witthaya22/golang-backend-itctc/modules/middlewares/middlewaresUsecase"
	"github.com/Witthaya22/golang-backend-itctc/pkg/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type midedlewareHandlerErrCode string

const (
	routerCheckErr midedlewareHandlerErrCode = "middleware-001"
	jwtAuthErr     midedlewareHandlerErrCode = "middleware-002"
)

type IMiddlewaresHandler interface {
	Cors() fiber.Handler
	RouterCheck() fiber.Handler
	Logger() fiber.Handler
	JwtAuth() fiber.Handler
}

type middlewaresHandler struct {
	conf              *config.Config
	middlewareUsecase middlewaresusecase.IMiddlewaresUsecase
}

func MiddlewaresHandler(conf *config.Config, middlewareUsecase middlewaresusecase.IMiddlewaresUsecase) IMiddlewaresHandler {
	return &middlewaresHandler{
		conf:              conf,
		middlewareUsecase: middlewareUsecase,
	}
}

func (h *middlewaresHandler) Cors() fiber.Handler {
	return cors.New(cors.Config{
		Next:             cors.ConfigDefault.Next,
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0,
	})
}

func (h *middlewaresHandler) RouterCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return customresponse.NewResponse(c).Error(
			fiber.ErrNotFound.Code,
			string(routerCheckErr),
			"Router not found",
		).Res()
	}
}

func (h *middlewaresHandler) Logger() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "${time} [${ip}] ${status} - ${method} ${path}\n",
		TimeFormat: "02/01/2006",
		TimeZone:   "bangkok/asia",
	})
}

func (h *middlewaresHandler) JwtAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// ดึง token จาก header
		token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
		if token == "" {
			return customresponse.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(jwtAuthErr),
				"missing token",
			).Res()
		}

		// Parse token
		jwtConf := *h.conf.Jwt
		result, err := auth.ParseToken(token, jwtConf)
		if err != nil {
			return customresponse.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(jwtAuthErr),
				err.Error(),
			).Res()
		}

		// ตรวจสอบ token ในฐานข้อมูล
		if !h.middlewareUsecase.FindAccessToken(result.Claims.UserID, token) {
			return customresponse.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(jwtAuthErr),
				"invalid access token",
			).Res()
		}

		// เก็บข้อมูลใน context
		c.Locals("userId", result.Claims.UserID)
		c.Locals("role", result.Claims.Role)

		return c.Next()
	}
}
