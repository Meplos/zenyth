package main

import (
	"os"
	"time"

	config "github.com/Meplos/zenyth/config"
	"github.com/Meplos/zenyth/tasks"
	"github.com/robfig/cron"
)

const TASK_FILENAME = "zenyth.tasks.json"

func main() {
	go Init()
}

func Init() {
	c := cron.New()
	bytes, err := os.ReadFile(TASK_FILENAME)
	if err != nil {
		panic(err)
	}
	definitions := config.ParseTaskDef(string(bytes))
	// runnableTasks := make([]tasks.Task, len(definitions))
	for _, def := range definitions {
		t := tasks.NewTask(def)
		t.Schedule(c)
	}

	c.Start()

	time.Sleep(5 * time.Second)
	c.Stop()
}
