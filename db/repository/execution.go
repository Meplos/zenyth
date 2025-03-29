package repository

import (
	"time"
)

type ExecutionEntity struct {
	Task     string
	Start    time.Time
	End      time.Time
	Duration int64
	Status   string
}
