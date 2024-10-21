package checkAppHandler

import (
	"github.com/Witthaya22/golang-backend-itctc/config"
	checkapp "github.com/Witthaya22/golang-backend-itctc/modules/checkApp"
	customresponse "github.com/Witthaya22/golang-backend-itctc/modules/customResponse"
	"github.com/gofiber/fiber/v2"
)

type ICheckAppHandler interface {
	HealthCheck(c *fiber.Ctx) error
}

type checkAppHandler struct {
	conf *config.Config
}

func NewCheckAppHandler(conf *config.Config) ICheckAppHandler {
	return &checkAppHandler{
		conf: conf,
	}
}

func (h *checkAppHandler) HealthCheck(c *fiber.Ctx) error {
	res := &checkapp.CheckApp{
		Name:    h.conf.Server.Name,
		Version: h.conf.Server.Version,
	}

	return customresponse.NewResponse(c).Success(fiber.StatusOK, res).Res()
}
