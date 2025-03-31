package db

import (
	"log"

	repo "github.com/Meplos/zenyth/db/repository"
	"github.com/Meplos/zenyth/tasks"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ZenythDatabase struct {
	Db *gorm.DB
}

func Connect() *ZenythDatabase {
	db, err := gorm.Open(sqlite.Open("zenyth.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Canno't open database")
	}
	return &ZenythDatabase{
		Db: db,
	}
}

func (zDB *ZenythDatabase) Init() {
	zDB.Db.AutoMigrate(&repo.TaskEntity{})
	zDB.Db.AutoMigrate(&repo.ExecutionEntity{})
}

func (zDB *ZenythDatabase) ListTask() []tasks.Task {
	var dbTasks []repo.TaskEntity
	result := zDB.Db.Find(&dbTasks)
	if result.Error != nil {
		log.Fatalf("Query all tasks fail")
	}

	ts := make([]tasks.Task, 0)
	for _, t := range dbTasks {
		ts = append(ts, tasks.FromEntity(t))
	}
	return ts
}

func (zDB *ZenythDatabase) ListExecution(name string) []tasks.Execution {
	var dbExec []repo.ExecutionEntity
	result := zDB.Db.Where("task = ?", name).Order(clause.OrderByColumn{
		Column: clause.Column{Name: "Start"},
		Desc:   true,
	}).Find(&dbExec)
	if result.Error != nil {
		log.Fatalf("Query all tasks fail")
	}

	execs := make([]tasks.Execution, 0)
	for _, e := range dbExec {
		execs = append(execs, tasks.ExecutionFromEntity(e))
	}
	return execs
}

func (zDB *ZenythDatabase) CreateTask(t tasks.Task) {
	newTask := repo.TaskEntity{
		Name:    t.Name,
		Exec:    t.Exec,
		LogFile: t.LogFile,
		State:   string(t.State),
		Runner:  t.Runner,
		Hash:    t.Hash,

		Cron:       t.Cron,
		Minute:     t.Minute,
		Hour:       t.Hour,
		DayInMonth: t.DayInMonth,
		Month:      t.Month,
		DayInWeek:  t.DayInWeek,
	}
	zDB.Db.Create(&newTask)
	log.Printf("Task Created with ID=%v", newTask.ID)
}

func (zDB *ZenythDatabase) FindTask(name string) *tasks.Task {
	var t repo.TaskEntity
	result := zDB.Db.Where("name = ?", name).First(&t)
	if result.Error != nil {
		return nil
	}
	task := tasks.FromEntity(t)
	return &task
}

func (zDB *ZenythDatabase) UpdateTask(t tasks.Task) {
	var model repo.TaskEntity
	zDB.Db.First(&model, "name=?", t.Name)
	zDB.Db.Model(&model).Updates(repo.TaskEntity{
		Name:    t.Name,
		Exec:    t.Exec,
		LogFile: t.LogFile,
		State:   string(t.State),
		Runner:  t.Runner,
		Hash:    t.Hash,

		Cron:       t.Cron,
		Minute:     t.Minute,
		Hour:       t.Hour,
		DayInMonth: t.DayInMonth,
		Month:      t.Month,
		DayInWeek:  t.DayInWeek,
	})
}

func (zDB *ZenythDatabase) LogExectution(e tasks.Execution) {
	newTask := repo.ExecutionEntity{
		Task:     e.Task,
		Start:    e.Start,
		End:      e.End,
		Duration: e.Duration,
		Status:   string(e.Status),
	}
	zDB.Db.Create(&newTask)
}
