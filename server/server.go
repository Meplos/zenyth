package server

import (
	"log"
	"time"

	"github.com/Meplos/zenyth/db"
	"github.com/Meplos/zenyth/tasks"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

type FormattedExecution struct {
	Status   tasks.ProcessState
	Start    string
	End      string
	Duration int64
}

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

	zServer.app.Get("/exec", func(c *fiber.Ctx) error {
		name := c.Query("task")
		exec := zServer.db.ListExecution(name)
		formatted := make([]FormattedExecution, 0)
		for _, e := range exec {
			formatted = append(formatted, formatExec(e))
		}

		return c.Render("exec", fiber.Map{
			"Title": name,
			"Execs": formatted,
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

func formatExec(exec tasks.Execution) FormattedExecution {
	return FormattedExecution{
		Status:   exec.Status,
		Start:    exec.Start.UTC().Format(time.DateTime),
		End:      exec.End.UTC().Format(time.DateTime),
		Duration: exec.Duration,
	}
}
