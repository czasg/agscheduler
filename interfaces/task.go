package interfaces

import "time"

type ITask interface {
	Go(now time.Time)
	Pause()
	Resume()
	GetName() string
	UpdateTrigger(trigger ITrigger)
	SetScheduler(scheduler IScheduler)
	GetNextRunTime(now time.Time) time.Time
}
