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
	taskObserver := NewTaskObserver(db)
	execObserver := NewExecutionObserver(db)
	var t *tasks.Task
	for _, def := range definitions {
		t = LoadTask(def, db, &taskObserver)
		t.AddExecutionObserver(&execObserver)
		t.Schedule(c)
	}

	c.Start()

	time.Sleep(15 * time.Minute)
	c.Stop()
	t.Stopped()
}

func LoadTask(t tasks.TaskDef, db *db.ZenythDatabase, o observer.Observer[tasks.Task]) *tasks.Task {
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

type TaskObserver struct {
	db *db.ZenythDatabase
}

type ExecutionObserver struct {
	db *db.ZenythDatabase
}

func NewTaskObserver(db *db.ZenythDatabase) TaskObserver {
	return TaskObserver{
		db: db,
	}
}

func NewExecutionObserver(db *db.ZenythDatabase) ExecutionObserver {
	return ExecutionObserver{
		db: db,
	}
}

func (o *ExecutionObserver) Notify(event observer.Event, data tasks.Execution) {
	switch event {
	case observer.Terminated:
		o.db.LogExectution(data)
		break
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
