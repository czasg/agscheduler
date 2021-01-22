package agscheduler

import (
	"math"
	"time"
)

var (
	MinDateTime = time.Time{}
	MaxDateTime = time.Now().Add(time.Duration(math.MaxInt64))
)

// GetNextRunTime: if result is MinDateTime, it mean this Job is over, should be delete.
type ITrigger interface {
	GetNextRunTime(previous, now time.Time) time.Time
}

type TriggerMeta struct {
	Type         string        `json:"type"`
	NextRunTime  time.Time     `json:"next_run_time"`
	Interval     time.Duration `json:"interval"`
	StartRunTime time.Time     `json:"start_run_time"`
	EndRunTime   time.Time     `json:"end_run_time"`
	CronCmd      string        `json:"cron_cmd"`
}
