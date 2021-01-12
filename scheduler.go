package agscheduler

import (
	"github.com/sirupsen/logrus"
)

//var (
//	MinDateTime = time.Time{}
//	MaxDateTime = time.Now().Add(time.Duration(math.MaxInt64))
//)

type AGScheduler struct {
	Store  IStore
	Logger *logrus.Entry
	Status STATUS
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
	ags.Status.SetRunning()
}

func (ags *AGScheduler) Close() error {
	ags.FillByDefault()
	return nil
}

func (ags *AGScheduler) Wake() {
	ags.FillByDefault()
}
