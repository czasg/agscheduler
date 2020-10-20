package interfaces

import "time"

type ITask interface {
	Go()
	Pause() error
	Resume() error
	GetName() string
	GetRunTimes(now time.Time) []time.Time
	SetScheduler(scheduler IScheduler)
}
