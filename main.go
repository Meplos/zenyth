package main

import (
	"bufio"
	"log"
	"os"

	"github.com/Meplos/zenyth/app"
	"github.com/Meplos/zenyth/manager"
	"github.com/Meplos/zenyth/server"
)

const TASK_FILENAME = "zenyth.tasks.json"

func main() {
	InitLogger()

	manager := manager.New()
	app.Init(manager).Run()
	server.Init(manager).Run()
}

func InitLogger() {
	err := os.Mkdir("/var/log/.zenyth", 0750)
	if err != nil && !os.IsExist(err) {
		log.Fatal("Could not create log dir")
	}

	logFile, err := os.OpenFile("/var/log/.zenyth/zth.log", os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Fatal("Could not create log file")
	}

	log.SetOutput(logFile)
}
