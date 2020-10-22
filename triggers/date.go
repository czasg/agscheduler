package triggers

import (
	"errors"
	"github.com/CzaOrz/AGScheduler"
	"time"
)

type DateTrigger struct {
	RunDateTime time.Time
}

func NewDateTrigger(runDateTime time.Time) (*DateTrigger, error) {
	if runDateTime.Before(AGScheduler.EmptyDateTime) {
		return nil, errors.New("Invalid Run Date Time")
	}
	return &DateTrigger{runDateTime}, nil
}

func (d DateTrigger) NextFireTime(previous, now time.Time) time.Time {
	if !previous.Equal(AGScheduler.EmptyDateTime) {
		return AGScheduler.EmptyDateTime
	}
	return d.RunDateTime
}
