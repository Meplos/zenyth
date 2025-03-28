package db

import (
	"log"

	repo "github.com/Meplos/zenyth/db/repository"
	"github.com/Meplos/zenyth/tasks"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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
}

func (zDB *ZenythDatabase) CreateTask(t tasks.Task) {
	newTask := repo.TaskEntity{
		Name:    t.Name,
		Exec:    t.Exec,
		LogFile: t.LogFile,
		State:   string(t.State),
		Hash:    t.Hash,

		Cron:       t.Cron,
		Second:     t.Second,
		Minute:     t.Minute,
		Hour:       t.Hour,
		DayInMonth: t.DayInMonth,
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
	log.Printf("Task found %v", t)
	log.Printf("Task found ID:%v", t.ID)
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
		Hash:    t.Hash,

		Cron:       t.Cron,
		Second:     t.Second,
		Minute:     t.Minute,
		Hour:       t.Hour,
		DayInMonth: t.DayInMonth,
		DayInWeek:  t.DayInWeek,
	})
}
