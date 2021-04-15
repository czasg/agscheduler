package agscheduler

import "time"

type IntervalTrigger struct {
	Interval     time.Duration `json:"interval"`
	StartRunTime time.Time     `json:"start_run_time"`
	EndRunTime   time.Time     `json:"end_run_time"`
}

func (t *IntervalTrigger) GetNextRunTime(previous, now time.Time) time.Time {
	if t.EndRunTime.After(now) {
		return MinDateTime
	}
	if t.StartRunTime.After(now) {
		return t.StartRunTime
	}
	if previous.Equal(MinDateTime) {
		if t.StartRunTime.Equal(MinDateTime) {
			t.StartRunTime = now.Add(-time.Nanosecond)
		}
		return t.StartRunTime
	}
	return previous.Add(t.Interval)
}
