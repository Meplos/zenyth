package main

import (
	// "github.com/Meplos/zenyth/app"
	"github.com/Meplos/zenyth/server"
)

const TASK_FILENAME = "zenyth.tasks.json"

func main() {
	// app.Init().Run()
	server.Init().Run()
}
