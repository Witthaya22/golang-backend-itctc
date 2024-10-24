package servers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/Witthaya22/golang-backend-itctc/config"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type IServer interface {
	Start()
}

type server struct {
	app  *fiber.App
	db   *gorm.DB
	conf *config.Config
}

func NewServer(conf *config.Config, db *gorm.DB) IServer {
	return &server{
		conf: conf,
		db:   db,
		app: fiber.New(fiber.Config{
			AppName:      conf.Server.Name,
			BodyLimit:    int(conf.Server.BodyLimit),
			WriteTimeout: conf.Server.Timeout,
			ReadTimeout:  conf.Server.ReadTimeout,
			JSONEncoder:  json.Marshal,
			JSONDecoder:  json.Unmarshal,
		}),
	}
}

func (s *server) Start() {
	// middleware
	middleware := InitMiddleware(s)
	s.app.Use(middleware.Logger())
	s.app.Use(middleware.Cors())

	// modules
	v1 := s.app.Group("/v1")
	modules := InitModule(v1, s, middleware)

	modules.CheckAppModule()
	modules.UserModule()

	s.app.Use(middleware.RouterCheck())

	// Graceful shutdown ค่อยๆปิดระบบทุกอย่าง
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Println("server shutting down...")
		s.app.Shutdown()
	}()

	log.Printf("server starting on port %d", s.conf.Server.Port)
	if err := s.app.Listen(fmt.Sprintf(":%d", s.conf.Server.Port)); err != nil {
		log.Panicf("Failed to start server: %v", err)
	}
}
