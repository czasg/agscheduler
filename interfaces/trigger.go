package interfaces

import "time"

type ITrigger interface {
	NextFireTime(previous, now time.Time) (time.Time, error)
}
