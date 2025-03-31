package tasks

import (
	"time"

	"github.com/Meplos/zenyth/db/repository"
)

type Execution struct {
	Task     string
	Start    time.Time
	End      time.Time
	Duration int64
	Status   ProcessState
}

func NewExecution(name string, start, end time.Time, status ProcessState) Execution {
	execDuration := end.Sub(start).Milliseconds()
	return Execution{
		Task:     name,
		Start:    start,
		End:      end,
		Duration: execDuration,
		Status:   status,
	}
}

func ExecutionFromEntity(entity repository.ExecutionEntity) Execution {
	return Execution{
		Task:     entity.Task,
		Start:    entity.Start,
		End:      entity.End,
		Duration: entity.Duration,
		Status:   ProcessState(entity.Status),
	}
}
