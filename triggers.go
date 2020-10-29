package AGScheduler

import (
	"errors"
	"github.com/robfig/cron"
	"time"
)

type TriggerState struct {
	Name     string        `json:"name"`
	Cron     CronState     `json:"cron"`
	Date     DateState     `json:"date"`
	Interval IntervalState `json:"interval"`
}

type CronState struct {
	CronCmd string `json:"cron_cmd"`
}

type DateState struct {
	RunDateTime time.Time `json:"run_date_time"`
}

type IntervalState struct {
	StartRunTime time.Time     `json:"start_run_time"`
	EndRunTime   time.Time     `json:"end_run_time"`
	Interval     time.Duration `json:"interval"`
}

func FromTriggerState(state TriggerState) ITrigger {
	switch state.Name {
	case "date":
		date, err := NewDateTrigger(state.Date.RunDateTime)
		if err != nil {
			return nil
		}
		return date
	case "interval":
		date, err := NewIntervalTrigger(state.Interval.StartRunTime, state.Interval.EndRunTime, state.Interval.Interval)
		if err != nil {
			return nil
		}
		return date
	case "cron":
		date, err := NewCronTrigger(state.Cron.CronCmd)
		if err != nil {
			return nil
		}
		return date
	default:
		return nil
	}
}

/*
*	Cron
**/
type CronTrigger struct {
	CronCmd   string
	StartTime time.Time
	Schedule  cron.Schedule
}

func NewCronTrigger(cronCmd string) (*CronTrigger, error) {
	scheduler, err := cron.Parse(cronCmd)
	if err != nil {
		return nil, err
	}
	return &CronTrigger{
		CronCmd:   cronCmd,
		StartTime: EmptyDateTime,
		Schedule:  scheduler,
	}, nil
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

func (c *CronTrigger) GetTriggerState() TriggerState {
	return TriggerState{
		Name: "cron",
		Cron: CronState{
			CronCmd: c.CronCmd,
		},
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

func (d *DateTrigger) NextFireTime(previous, now time.Time) time.Time {
	if !previous.Equal(EmptyDateTime) {
		return EmptyDateTime
	}
	return d.RunDateTime
}

func (d *DateTrigger) GetTriggerState() TriggerState {
	return TriggerState{
		Name: "date",
		Date: DateState{
			RunDateTime: d.RunDateTime,
		},
	}
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

func (i *IntervalTrigger) NextFireTime(previous, now time.Time) time.Time {
	if !i.EndRunTime.Equal(EmptyDateTime) && i.EndRunTime.Before(now) {
		return EmptyDateTime
	}
	if previous.Equal(EmptyDateTime) {
		return i.StartRunTime
	}
	return previous.Add(i.Interval)
}

func (i *IntervalTrigger) GetTriggerState() TriggerState {
	return TriggerState{
		Name: "interval",
		Interval: IntervalState{
			StartRunTime: i.StartRunTime,
			EndRunTime:   i.EndRunTime,
			Interval:     i.Interval,
		},
	}
}
