package servers

import (
	checkAppHandler "github.com/Witthaya22/golang-backend-itctc/modules/checkApp/handler"
	"github.com/gofiber/fiber/v2"
)

type IModuleFatory interface {
	CheckAppModule()
}

type moduleFactory struct {
	r fiber.Router
	s *server
}

func InitModule(r fiber.Router, s *server) IModuleFatory {
	return &moduleFactory{
		r: r,
		s: s,
	}
}

func (m *moduleFactory) CheckAppModule() {
	handler := checkAppHandler.NewCheckAppHandler(m.s.conf)

	m.r.Get("/", handler.HealthCheck)
}
