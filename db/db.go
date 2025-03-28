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
	zDB.Db.Create(&repo.TaskEntity{
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
