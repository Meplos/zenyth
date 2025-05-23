package repository

import (
	"gorm.io/gorm"
)

type TaskEntity struct {
	gorm.Model
	Name    string `gorm:"primarykey"`
	Exec    string
	LogFile string
	Hash    string
	Runner  string
	State   string

	Cron       string
	Second     string
	Minute     string
	Hour       string
	DayInMonth string
	Month      string
	DayInWeek  string
}
