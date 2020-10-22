package interfaces

import "time"

type ITask interface {
	Go(now time.Time)
	Pause()
	Resume()
	GetName() string
	SetScheduler(scheduler IScheduler)
	GetNextRunTime(now time.Time) time.Time
}
