package server

import (
	"log"

	"github.com/Meplos/zenyth/db"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

type ZenythServer struct {
	app  *fiber.App
	port string
	db   *db.ZenythDatabase
}

func Init() *ZenythServer {
	engine := html.New("./layouts", ".html")
	zServer := &ZenythServer{
		app: fiber.New(fiber.Config{
			Views: engine,
		}),

		port: ":9823",
		db:   db.Connect(),
	}
	zServer.app.Get("/", func(c *fiber.Ctx) error {
		tasks := zServer.db.ListTask()
		return c.Render("index", fiber.Map{
			"Tasks": tasks,
		}, "main")
	})

	zServer.app.Static("/public", "./public/")
	return zServer
}

func (s *ZenythServer) Run() {
	s.app.Listen(s.port)
}

func (s *ZenythServer) Stop() {
	log.Printf("Closing...")
}
