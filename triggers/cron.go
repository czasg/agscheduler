package triggers

import (
	"github.com/CzaOrz/AGScheduler"
	"time"
)

type CronTrigger struct {
	StartRunTime time.Time
}

func NewCronTrigger() *CronTrigger {
	return &CronTrigger{}
}

func (c CronTrigger) NextFireTime(previous, now time.Time) time.Time {
	return AGScheduler.EmptyDateTime
}
