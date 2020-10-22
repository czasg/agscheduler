package tasks

import (
	"github.com/CzaOrz/AGScheduler/interfaces"
	"github.com/sirupsen/logrus"
	"runtime"
	"time"
)

type Task struct {
	Name            string
	Func            func(args []interface{})
	Args            []interface{}
	Scheduler       interfaces.IScheduler
	Trigger         interfaces.ITrigger
	Store           interfaces.IStore
	PreviousRunTime time.Time
	Logger          *logrus.Entry
}

func NewTask(name string, method func(args []interface{}), args []interface{}, trigger interfaces.ITrigger) *Task {
	logrus.WithFields(logrus.Fields{})
	return &Task{
		Name:    name,
		Func:    method,
		Args:    args,
		Trigger: trigger,
		Logger: logrus.WithFields(logrus.Fields{
			"Module":   "AGScheduler.Task",
			"TaskName": name,
		}),
	}
}

func (t *Task) Go(runTime time.Time) {
	t.PreviousRunTime = runTime
	go func() {
		defer func() {
			if r := recover(); r != nil {
				const size = 64 << 10
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]
				t.Logger.Errorf("cron: panic running task: %v\n%s", r, buf)
			}
		}()
		t.Func(t.Args)
	}()
}

func (t *Task) GetName() string {
	return t.Name
}

func (t *Task) GetNextRunTime(now time.Time) time.Time {
	return t.Trigger.NextFireTime(t.PreviousRunTime, now)
}
