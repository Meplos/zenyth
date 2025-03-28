package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	config "github.com/Meplos/zenyth/config"
	"github.com/robfig/cron"
)

const TASK_FILENAME = "zenyth.tasks.json"

func main() {
	c := cron.New()
	bytes, err := os.ReadFile(TASK_FILENAME)
	if err != nil {
		panic(err)
	}
	tasks := config.ParseTaskDef(string(bytes))

	for _, t := range tasks {
		fmt.Printf("Stating tasks %v\n", t.Name)
		c.AddFunc(t.Cron, func() {
			exec.Command(t.Exec, ">>", "/home/aerard/Documents/output.log")
		})

	}

	c.Start()

	time.Sleep(5 * time.Second)
	c.Stop()
}
