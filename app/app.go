package app

import (
	"log"
	"os"

	"github.com/Meplos/zenyth/config"
	"github.com/Meplos/zenyth/db"
	"github.com/Meplos/zenyth/manager"
	"github.com/Meplos/zenyth/observer"
	taskobserver "github.com/Meplos/zenyth/taskObserver"
	"github.com/Meplos/zenyth/tasks"
	"github.com/robfig/cron"
)

type App struct {
	db        *db.ZenythDatabase
	scheduler *cron.Cron
	manager   *manager.CronManager
	taskFile  string
}

func Init(manager *manager.CronManager) *App {
	zdb := db.Connect()
	zdb.Init()
	scheduler := cron.New()
	return &App{
		db:        zdb,
		scheduler: scheduler,
		manager:   manager,
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
		t = a.LoadTask(def, a.db, &to)
		t.AddExecutionObserver(&eo)
		a.manager.ScheduleTasks(t)
	}

	a.manager.StartAll()
}

func (a *App) Stop() {
	a.manager.StopAll()
	os.Exit(0)
}

func (a *App) LoadTask(t tasks.TaskDef, db *db.ZenythDatabase, o observer.Observer[tasks.Task]) *tasks.Task {
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
