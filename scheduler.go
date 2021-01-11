package agscheduler

import (
	"github.com/sirupsen/logrus"
	"math"
	"time"
)

type STATUS int64

var (
	AGSLog = logrus.New()
	Log    = AGSLog.WithFields(logrus.Fields{
		"AGSVersion": Version,
	})
	STATUS_RUNNING = STATUS(0)
	STATUS_PAUSED  = STATUS(1)
	STATUS_STOPPED = STATUS(2)
)

var (
	MinDateTime = time.Time{}
	MaxDateTime = time.Now().Add(time.Duration(math.MaxInt64))
)

type AGScheduler struct {
	Store  IStore
	Logger *logrus.Entry
}

func (ags *AGScheduler) FillByDefault() {
	if ags.Store == nil {
		ags.Store = &MemoryStore{}
	}
	if ags.Logger == nil {
		ags.Logger = Log.WithFields(logrus.Fields{
			"ASGModule": "Scheduler",
		})
	}
}

func (ags *AGScheduler) Start() {
	ags.FillByDefault()
}

func (ags *AGScheduler) Close() error {
	ags.FillByDefault()
	return nil
}

func (ags *AGScheduler) Wake() {
	ags.FillByDefault()
}
