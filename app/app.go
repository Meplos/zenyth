package app

import (
	"log"
	"os"
	"os/signal"

	"github.com/Meplos/zenyth/config"
	"github.com/Meplos/zenyth/db"
	"github.com/Meplos/zenyth/observer"
	taskobserver "github.com/Meplos/zenyth/taskObserver"
	"github.com/Meplos/zenyth/tasks"
	"github.com/robfig/cron"
)

type App struct {
	db        *db.ZenythDatabase
	scheduler *cron.Cron
	taskFile  string
}

func Init() *App {
	zdb := db.Connect()
	zdb.Init()
	scheduler := cron.New()
	return &App{
		db:        zdb,
		scheduler: scheduler,
		taskFile:  "zenyth.tasks.json",
	}
}

func (a *App) Run() {
	bytes, err := os.ReadFile(a.taskFile)
	if err != nil {
		panic(err)
	}
	definitions := config.ParseTaskDef(string(bytes))
	to := taskobserver.NewTaskObserver(a.db)
	eo := taskobserver.NewExecutionObserver(a.db)
	var t *tasks.Task
	for _, def := range definitions {
		t = a.loadTask(def, a.db, &to)
		t.AddExecutionObserver(&eo)
		t.Schedule(a.scheduler)
	}

	a.scheduler.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			a.Stop()
		}
	}()

	for true {
	}
}

func (a *App) Stop() {
	a.scheduler.Stop()
	os.Exit(0)
}

func (a *App) loadTask(t tasks.TaskDef, db *db.ZenythDatabase, o observer.Observer[tasks.Task]) *tasks.Task {
	log.Printf("Load %v\n", t.Name)
	saveTask := db.FindTask(t.Name)
	newTask := tasks.NewTask(t)
	newTask.AddTaskObserver(o)
	if saveTask == nil {
		log.Printf("New task detected : %v", newTask.Name)
		o.Notify(observer.Create, *newTask)
		return newTask
	}

	if saveTask.Hash == newTask.Hash {
		log.Printf("Task modification detected: %v", newTask.Name)
		o.Notify(observer.Update, *newTask)
		return newTask
	}

	log.Printf("No change detected for : %v", newTask.Name)

	return newTask
}
