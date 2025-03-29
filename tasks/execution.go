package tasks

import "time"

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
