package triggers

import (
	"github.com/CzaOrz/AGScheduler"
	"time"
)

type DateTrigger struct {
	RunDateTime time.Time
}

func NewDateTrigger(runDateTime time.Time) *DateTrigger {
	return &DateTrigger{
		RunDateTime: runDateTime,
	}
}

func (d DateTrigger) NextFireTime(previous, now time.Time) time.Time {
	if !previous.Equal(AGScheduler.EmptyDateTime) {
		return AGScheduler.EmptyDateTime
	}
	return d.RunDateTime
}
