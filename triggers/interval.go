package triggers

import (
	"errors"
	"github.com/CzaOrz/AGScheduler"
	"time"
)

type IntervalTrigger struct {
	Interval     time.Duration
	StartRunTime time.Time
	EndRunTime   time.Time
}

func NewIntervalTrigger(startTime, endTime time.Time, interval time.Duration) (*IntervalTrigger, error) {
	if startTime.After(endTime) && !endTime.Equal(AGScheduler.EmptyDateTime) {
		return nil, errors.New("Invalid Interval time: endTime should be AGScheduler.EmptyDateTime")
	}
	return &IntervalTrigger{
		Interval:     interval,
		StartRunTime: startTime,
		EndRunTime:   endTime,
	}, nil
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
