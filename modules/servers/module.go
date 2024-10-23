package servers

import (
	checkAppHandler "github.com/Witthaya22/golang-backend-itctc/modules/checkApp/handler"
	middlewareshandler "github.com/Witthaya22/golang-backend-itctc/modules/middlewares/middlewaresHandler"
	middlewaresrepository "github.com/Witthaya22/golang-backend-itctc/modules/middlewares/middlewaresRepository"
	middlewaresusecase "github.com/Witthaya22/golang-backend-itctc/modules/middlewares/middlewaresUsecase"
	"github.com/gofiber/fiber/v2"
)

type IModuleFatory interface {
	CheckAppModule()
}

type moduleFactory struct {
	r   fiber.Router
	s   *server
	mid middlewareshandler.IMiddlewaresHandler
}

func InitModule(r fiber.Router, s *server, mid middlewareshandler.IMiddlewaresHandler) IModuleFatory {
	return &moduleFactory{
		r:   r,
		s:   s,
		mid: mid,
	}
}

func InitMiddleware(s *server) middlewareshandler.IMiddlewaresHandler {
	repository := middlewaresrepository.MiddlewaresRepository(s.db)
	usecase := middlewaresusecase.MiddlewaresUsecase(repository)
	return middlewareshandler.MiddlewaresHandler(s.conf, usecase)
}
func (m *moduleFactory) CheckAppModule() {
	handler := checkAppHandler.NewCheckAppHandler(m.s.conf)

	m.r.Get("/", handler.HealthCheck)
}
