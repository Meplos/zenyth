package main

import (
	"log"
	"os"
	"time"

	config "github.com/Meplos/zenyth/config"
	"github.com/Meplos/zenyth/db"
	"github.com/Meplos/zenyth/observer"
	"github.com/Meplos/zenyth/tasks"
	"github.com/robfig/cron"
)

const TASK_FILENAME = "zenyth.tasks.json"

func main() {
	Init()
}

func Init() {
	db := db.Connect()
	db.Init()
	c := cron.New()
	bytes, err := os.ReadFile(TASK_FILENAME)
	if err != nil {
		panic(err)
	}
	definitions := config.ParseTaskDef(string(bytes))
	// runnableTasks := make([]tasks.Task, len(definitions))
	taskObserver := New(db)
	var t tasks.Task
	for _, def := range definitions {
		t = tasks.NewTask(def, &taskObserver)
		t.Schedule(c)
	}

	c.Start()

	time.Sleep(5 * time.Second)
	c.Stop()
	t.Stopped()
}

type TaskObserver struct {
	db *db.ZenythDatabase
}

func New(db *db.ZenythDatabase) TaskObserver {
	return TaskObserver{
		db: db,
	}
}

func (o *TaskObserver) Notify(event observer.Event, data tasks.Task) {
	log.Printf("Event=%v, Tasks=%v, State=%v", event, data.Name, data.State)
	switch event {
	case observer.Create:
		o.db.CreateTask(data)
		break
	case observer.Update:
		o.db.UpdateTask(data)
		break
	}
}
