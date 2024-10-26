package servers

import (
	checkAppHandler "github.com/Witthaya22/golang-backend-itctc/modules/checkApp/handler"
	middlewareshandler "github.com/Witthaya22/golang-backend-itctc/modules/middlewares/middlewaresHandler"
	middlewaresrepository "github.com/Witthaya22/golang-backend-itctc/modules/middlewares/middlewaresRepository"
	middlewaresusecase "github.com/Witthaya22/golang-backend-itctc/modules/middlewares/middlewaresUsecase"
	userhandler "github.com/Witthaya22/golang-backend-itctc/modules/users/userHandler"
	userrepository "github.com/Witthaya22/golang-backend-itctc/modules/users/userRepository"
	userusecase "github.com/Witthaya22/golang-backend-itctc/modules/users/userUsecase"
	"github.com/gofiber/fiber/v2"
)

type IModuleFatory interface {
	CheckAppModule()
	UserModule()
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

func (m *moduleFactory) UserModule() {
	repository := userrepository.UserRepository(m.s.db)
	usecase := userusecase.UserUsecase(m.s.conf, repository)
	handler := userhandler.UserHandler(m.s.conf, usecase)

	router := m.r.Group("/user")

	router.Post("/signup", handler.SignUpUser)
	router.Post("/login", handler.SignIn)
}
