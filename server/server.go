package server

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

type ZenythServer struct {
	app  *fiber.App
	port string
}

func Init() *ZenythServer {
	server := fiber.New()
	server.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Helloworld")
	})
	return &ZenythServer{
		app:  server,
		port: ":9823",
	}
}

func (s *ZenythServer) Run() {
	s.app.Listen(s.port)
}

func (s *ZenythServer) Stop() {
	log.Printf("Closing...")
}
