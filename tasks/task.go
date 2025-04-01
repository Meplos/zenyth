package tasks

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Meplos/zenyth/db/repository"
	"github.com/Meplos/zenyth/observer"
	"github.com/Meplos/zenyth/runners"
	"github.com/robfig/cron"
)

type TaskDef struct {
	Name   string `json:"name"`
	Exec   string `json:"exec"`
	Cron   string `json:"cron"`
	Runner string `json:"runner"`
}

type (
	TaskState    string
	ProcessState string
)

const (
	SUCCESS ProcessState = "SUCCESS"
	FAILURE ProcessState = "FAILURE"
)

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
	Hash    string
	State   TaskState
	Runner  string

	Cron string

	Second     string
	Minute     string
	Hour       string
	DayInMonth string
	Month      string
	DayInWeek  string

	logger       *log.Logger
	taskObserver []observer.Observer[Task]
	execObserver []observer.Observer[Execution]
}

func NewTask(def TaskDef) *Task {
	cronExpr := strings.Split(def.Cron, " ")
	if len(cronExpr) != 6 {
		log.Fatalf("Invalid CRON expression for %v", def.Name)
	}
	bytes := taskToBytes(def)
	hash := md5.Sum(bytes)
	logFile := fmt.Sprintf("%v.log", def.Name)
	f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0660)
	if err != nil {
		log.Printf("%v", err)
		log.Fatalf("Cannot open logfile %v", logFile)
	}
	logger := log.New(f, def.Name, log.LstdFlags)
	return &Task{
		Name:    def.Name,
		Exec:    def.Exec,
		LogFile: logFile,
		State:   PENDING,
		Runner:  def.Runner,
		Hash:    string(hash[:]),

		Cron:       def.Cron,
		Second:     cronExpr[0],
		Minute:     cronExpr[1],
		Hour:       cronExpr[2],
		DayInMonth: cronExpr[3],
		Month:      cronExpr[4],
		DayInWeek:  cronExpr[5],

		logger:       logger,
		taskObserver: make([]observer.Observer[Task], 0),
		execObserver: make([]observer.Observer[Execution], 0),
	}
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
	start := time.Now()
	t.Running()
	log.Printf("Task starting %v, with command %v [hash=%v, state=%v]", t.Name, t.Exec, string(t.Hash[:]), t.State)
	// Execute commande
	runner := runners.GetRunner(t.Runner)
	output, err := runner.Exec(t.Exec)
	end := time.Now()
	lines := strings.Split(string(output), "\n")
	for _, l := range lines {
		t.logger.Printf(l)
	}
	if err != nil {
		t.Errored()
		t.EndProcess(start, end, FAILURE)
		log.Printf("Task FAILED %v, with command %v [hash=%v, state=%v]", t.Name, t.Exec, string(t.Hash[:]), t.State)
		return
	}

	t.Pending()
	log.Printf("Task termindated %v, with command %v [hash=%v, state=%v]", t.Name, t.Exec, string(t.Hash[:]), t.State)
	t.EndProcess(start, end, SUCCESS)
}

func (t *Task) EndProcess(start, end time.Time, state ProcessState) {
	execution := NewExecution(t.Name, start, end, state)
	t.NotifyExecution(observer.Terminated, execution)
}

func (t *Task) Schedule(c *cron.Cron) {
	c.AddJob(t.Cron, t)
}

func (t *Task) AddTaskObserver(o observer.Observer[Task]) {
	t.taskObserver = append(t.taskObserver, o)
}

func (t *Task) AddExecutionObserver(o observer.Observer[Execution]) {
	t.execObserver = append(t.execObserver, o)
}

func (t *Task) Notify(event observer.Event) {
	for _, o := range t.taskObserver {
		o.Notify(event, *t)
	}
}

func (t *Task) NotifyExecution(event observer.Event, data Execution) {
	for _, o := range t.execObserver {
		o.Notify(event, data)
	}
}

func taskToBytes(task TaskDef) []byte {
	bytes, err := json.Marshal(task)
	if err != nil {
		log.Fatalf("Unable to marshall task : %v", task.Name)
	}
	return bytes
}

func FromEntity(t repository.TaskEntity) Task {
	return Task{
		Name:    t.Name,
		Exec:    t.Exec,
		State:   TaskState(t.State),
		LogFile: t.LogFile,
		Runner:  t.Runner,
		Hash:    t.Hash,

		Cron:       t.Cron,
		Second:     t.Second,
		Minute:     t.Minute,
		Hour:       t.Hour,
		DayInMonth: t.DayInMonth,
		Month:      t.Month,
		DayInWeek:  t.DayInWeek,

		taskObserver: make([]observer.Observer[Task], 0),
	}
}
