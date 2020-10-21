package triggers

import (
	"github.com/CzaOrz/AGScheduler"
	"time"
)

type IntervalTrigger struct {
	Interval     time.Duration
	StartRunTime time.Time
	EndRunTime   time.Time
}

func NewIntervalTrigger(startTime, endTime time.Time, interval time.Duration) *IntervalTrigger {
	if startTime.After(endTime) && !endTime.Equal(AGScheduler.EmptyDateTime) {
		panic("Invalid Interval time: endTime should be AGScheduler.EmptyDateTime")
	}
	return &IntervalTrigger{
		Interval:     interval,
		StartRunTime: startTime,
		EndRunTime:   endTime,
	}
}

func (i IntervalTrigger) NextFireTime(previous, now time.Time) time.Time {
	if !i.EndRunTime.Equal(AGScheduler.EmptyDateTime) && i.EndRunTime.Before(now) {
		return AGScheduler.EmptyDateTime
	}
	if !previous.Equal(AGScheduler.EmptyDateTime) {
		return previous.Add(i.Interval)
	}
	return i.StartRunTime
}
