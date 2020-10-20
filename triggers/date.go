package triggers

import (
	"errors"
	"time"
)

var DateTriggerError = errors.New("Date Trigger Done")

type DateTrigger struct {
	Time time.Time
	Done bool
}

func NewDateTrigger(time time.Time) *DateTrigger {
	return &DateTrigger{
		Time: time,
		Done: false,
	}
}

func (t *DateTrigger) NextFireTime(previous, now time.Time) (time.Time, error) {
	if t.Done {
		return time.Time{}, DateTriggerError
	}
	t.Done = true
	return t.Time, nil
}
