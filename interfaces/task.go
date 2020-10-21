package interfaces

import "time"

type ITask interface {
	Go(now time.Time)
	GetName() string
	GetNextRunTime(now time.Time) time.Time
}
