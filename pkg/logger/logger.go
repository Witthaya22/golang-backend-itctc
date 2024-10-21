package logger

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Witthaya22/golang-backend-itctc/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type ILogger interface {
	Print() ILogger
	Save()
	SetQuery(c *fiber.Ctx)
	SetBody(c *fiber.Ctx)
	SetResponse(res any)
}

type logger struct {
	Time       string `json:"time"`
	Ip         string `json:"ip"`
	Method     string `json:"method"`
	StatusCode int    `json:"status_code"`
	Path       string `json:"path"`
	Query      any    `json:"query"`
	Body       any    `json:"body"`
	Response   any    `json:"response"`
}

func InitLogger(c *fiber.Ctx, res any) ILogger {
	log := &logger{
		Time:       time.Now().Format("2006-01-02 15:04:05"),
		Ip:         c.IP(),
		Method:     c.Method(),
		Path:       c.Path(),
		StatusCode: c.Response().StatusCode(),
	}
	log.SetQuery(c)
	log.SetBody(c)
	log.SetResponse(res)
	return log
}

func (l *logger) Print() ILogger {
	utils.Debug(l)
	return l
}
func (l *logger) Save() {
	data := utils.Output(l)

	filename := fmt.Sprintf("./assets/log/backendlogger_%v.txt", time.Now().Format("2006-01-02"))
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()
	file.WriteString(string(data) + "\n")
}
func (l *logger) SetQuery(c *fiber.Ctx) {
	if len(c.Request().URI().QueryString()) == 0 {
		l.Query = nil
		return
	}

	var body map[string]interface{}
	if err := c.QueryParser(&body); err != nil {
		log.Printf("query parser error: %v", err)
		l.Query = nil
		return
	}
	l.Query = body
}
func (l *logger) SetBody(c *fiber.Ctx) {
	// Check if the request has a body
	if len(c.Body()) == 0 {
		l.Body = nil
		return
	}

	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		log.Printf("body parser error: %v", err)
		l.Body = nil
		return
	}

	switch l.Path {
	case "v1/users/signup":
		l.Body = "never gonna give you up"
	default:
		l.Body = body
	}
}
func (l *logger) SetResponse(res any) {
	l.Response = res
}
