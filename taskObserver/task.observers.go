package taskobserver

import (
	"github.com/Meplos/zenyth/db"
	"github.com/Meplos/zenyth/observer"
	"github.com/Meplos/zenyth/tasks"
)

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
	switch event {
	case observer.Create:
		o.db.CreateTask(data)
		break
	case observer.Update:
		o.db.UpdateTask(data)
		break
	}
}
