package AGScheduler

import (
	"errors"
	"github.com/robfig/cron"
	"time"
)

/*
*	Cron
**/
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
	if previous.Equal(EmptyDateTime) {
		if c.StartTime.Equal(EmptyDateTime) {
			c.StartTime = c.Schedule.Next(now)
		}
		return c.StartTime
	}
	return c.Schedule.Next(previous)
}

func (c *CronTrigger) GetTriggerState(previous, now time.Time) map[string]interface{} {
	return map[string]interface{}{
		"name": "cron",
		"cmd":  "",
	}
}

/*
*	Date
**/
type DateTrigger struct {
	RunDateTime time.Time
}

func NewDateTrigger(runDateTime time.Time) (*DateTrigger, error) {
	if runDateTime.Before(EmptyDateTime) {
		return nil, errors.New("Invalid Run Date Time")
	}
	return &DateTrigger{runDateTime}, nil
}

func (d DateTrigger) NextFireTime(previous, now time.Time) time.Time {
	if !previous.Equal(EmptyDateTime) {
		return EmptyDateTime
	}
	return d.RunDateTime
}

/*
*	Interval
**/
type IntervalTrigger struct {
	Interval     time.Duration
	StartRunTime time.Time
	EndRunTime   time.Time
}

func NewIntervalTrigger(startTime, endTime time.Time, interval time.Duration) (*IntervalTrigger, error) {
	if startTime.After(endTime) && !endTime.Equal(EmptyDateTime) {
		return nil, errors.New("Invalid Interval time: endTime should be AGScheduler.EmptyDateTime")
	}
	return &IntervalTrigger{
		Interval:     interval,
		StartRunTime: startTime,
		EndRunTime:   endTime,
	}, nil
}

func (i IntervalTrigger) NextFireTime(previous, now time.Time) time.Time {
	if !i.EndRunTime.Equal(EmptyDateTime) && i.EndRunTime.Before(now) {
		return EmptyDateTime
	}
	if previous.Equal(EmptyDateTime) {
		return i.StartRunTime
	}
	return previous.Add(i.Interval)
}
