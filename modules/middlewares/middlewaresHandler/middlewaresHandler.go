package middlewareshandler

import (
	"github.com/Witthaya22/golang-backend-itctc/config"
	customresponse "github.com/Witthaya22/golang-backend-itctc/modules/customResponse"
	middlewaresusecase "github.com/Witthaya22/golang-backend-itctc/modules/middlewares/middlewaresUsecase"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type midedlewareHandlerErrCode string

const (
	routerCheckErr midedlewareHandlerErrCode = "middleware-001"
)

type IMiddlewaresHandler interface {
	Cors() fiber.Handler
	RouterCheck() fiber.Handler
	Logger() fiber.Handler
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
