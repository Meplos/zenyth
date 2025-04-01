package main

import (
	"github.com/Meplos/zenyth/app"
	"github.com/Meplos/zenyth/manager"
	"github.com/Meplos/zenyth/server"
)

const TASK_FILENAME = "zenyth.tasks.json"

func main() {
	manager := manager.New()
	go server.Init(manager).Run()
	app.Init(manager).Run()
}
