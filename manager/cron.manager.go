package manager

import (
	"github.com/Meplos/zenyth/observer"
	"github.com/Meplos/zenyth/tasks"
	"github.com/robfig/cron"
)

type CronManager struct {
	schedulers map[string]*CronTask
}

type CronTask struct {
	Task          *tasks.Task
	CronScheduler *cron.Cron
}

func New() *CronManager {
	return &CronManager{
		schedulers: make(map[string]*CronTask),
	}
}

type Manager interface {
	StartAll()
	StopAll()
	StartOne(key string)
	StopOne(key string)
	ScheduleTask(tasks.Task)
}

func (m *CronManager) ScheduleTasks(task *tasks.Task) {
	scheduler := cron.New()
	task.AddTaskObserver(m)
	scheduler.AddJob(task.Cron, task)
	m.schedulers[task.Hash] = &CronTask{
		Task:          task,
		CronScheduler: scheduler,
	}
}

func (m *CronManager) StartAll() {
	for _, s := range m.schedulers {
		s.CronScheduler.Start()
		s.Task.Pending()
	}
}

func (m *CronManager) StopAll() {
	for _, s := range m.schedulers {
		s.CronScheduler.Stop()
		s.Task.Stopped()

	}
}

func (m *CronManager) StartOne(key string) {
	s := m.schedulers[key]
	if s == nil {
		return
	}
	s.CronScheduler.Start()
	s.Task.Pending()
}

func (m *CronManager) StopOne(key string) {
	s := m.schedulers[key]
	if s == nil {
		return
	}
	s.CronScheduler.Stop()
	s.Task.Stopped()
}

func (m *CronManager) Notify(event observer.Event, data tasks.Task) {
	switch event {
	case observer.Errored:
		{
			m.StopOne(data.Hash)
		}
	}
}
