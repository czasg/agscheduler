package triggers

import (
	"github.com/CzaOrz/AGScheduler"
	"github.com/robfig/cron"
	"time"
)

type CronTrigger struct {
	StartTime time.Time
	Schedule  cron.Schedule
}

func NewCronTrigger(cronCmd string) (*CronTrigger, error) {
	scheduler, err := cron.Parse(cronCmd)
	if err != nil {
		return nil, err
	}
	return &CronTrigger{Schedule: scheduler}, nil
}

func (c *CronTrigger) NextFireTime(previous, now time.Time) time.Time {
	if previous.Equal(AGScheduler.EmptyDateTime) {
		if c.StartTime.Equal(AGScheduler.EmptyDateTime) {
			c.StartTime = c.Schedule.Next(now)
		}
		return c.StartTime
	}
	return c.Schedule.Next(previous)
}
