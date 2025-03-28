package tasks

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/Meplos/zenyth/observer"
	"github.com/robfig/cron"
)

type TaskDef struct {
	Name string `json:"name"`
	Exec string `json:"exec"`
	Cron string `json:"cron"`
}

type TaskState string

const (
	PENDING TaskState = "PENDING"
	RUNING  TaskState = "RUNING"
	ERRORED TaskState = "ERRORED"
	STOPPED TaskState = "STOPED"
)

type Task struct {
	Name    string
	Exec    string
	LogFile string
	Hash    [16]byte
	State   TaskState

	Cron       string
	Second     string
	Minute     string
	Hour       string
	DayInMonth string
	DayInWeek  string

	observer []observer.Observer[Task]
}

func NewTask(def TaskDef, o observer.Observer[Task]) Task {
	cronExpr := strings.Split(def.Cron, " ")
	if len(cronExpr) != 5 {
		log.Fatalf("Invalid CRON expression for %v", def.Name)
	}
	bytes := taskToBytes(def)

	hash := md5.Sum(bytes)

	t := Task{
		Name:    def.Name,
		Exec:    def.Exec,
		LogFile: fmt.Sprintf("/tmp/.zenyth/%v.log", def.Name),
		State:   PENDING,
		Hash:    hash,

		Cron:       def.Cron,
		Second:     cronExpr[0],
		Minute:     cronExpr[1],
		Hour:       cronExpr[2],
		DayInMonth: cronExpr[3],
		DayInWeek:  cronExpr[4],

		observer: make([]observer.Observer[Task], 0),
	}
	t.AddObserver(o)
	t.Notify(observer.Create)
	return t
}

func (t *Task) Running() {
	t.State = RUNING
	t.Notify(observer.Update)
}

func (t *Task) Pending() {
	t.State = PENDING
	t.Notify(observer.Update)
}

func (t *Task) Errored() {
	t.State = ERRORED
	t.Notify(observer.Update)
}

func (t *Task) Stopped() {
	t.State = STOPPED
	t.Notify(observer.Update)
}

func (t *Task) Run() {
	t.Running()
	log.Printf("Task starting %v, with command %v [hash=%v, state=%v]", t.Name, t.Exec, string(t.Hash[:]), t.State)
	// Execute commande
	output, err := exec.Command(t.Exec).Output()
	if err != nil {
		t.Errored()
		return
	}
	lines := strings.Split(string(output), "\n")
	for _, l := range lines {
		log.Printf(l)
	}

	t.Pending()
	log.Printf("Task termindated %v, with command %v [hash=%v, state=%v]", t.Name, t.Exec, string(t.Hash[:]), t.State)
}

func (t *Task) Schedule(c *cron.Cron) {
	c.AddJob(t.Cron, t)
}

func (t *Task) AddObserver(o observer.Observer[Task]) {
	t.observer = append(t.observer, o)
}

func (t *Task) Notify(event observer.Event) {
	for _, o := range t.observer {
		o.Notify(event, *t)
	}
}

func taskToBytes(task TaskDef) []byte {
	bytes, err := json.Marshal(task)
	if err != nil {
		log.Fatalf("Unable to marshall task : %v", task.Name)
	}
	return bytes
}
