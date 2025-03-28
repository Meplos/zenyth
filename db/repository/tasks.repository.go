package repository

import (
	"gorm.io/gorm"
)

type TaskEntity struct {
	gorm.Model
	Name    string
	Exec    string
	LogFile string
	Hash    [16]byte
	State   string

	Cron       string
	Second     string
	Minute     string
	Hour       string
	DayInMonth string
	DayInWeek  string
}
