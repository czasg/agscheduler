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
