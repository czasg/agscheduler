package tasks

import (
	"github.com/CzaOrz/AGScheduler"
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
	Running         bool
}

func NewTask(name string, method func(args []interface{}), args []interface{}, trigger interfaces.ITrigger) *Task {
	return &Task{
		Name:    name,
		Func:    method,
		Args:    args,
		Trigger: trigger,
		Logger: logrus.WithFields(logrus.Fields{
			"Module":   "AGScheduler.Task",
			"TaskName": name,
		}),
		Running: true,
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
				t.Logger.WithFields(logrus.Fields{
					"Func": "Go",
				}).Errorf("cron: panic running task: %v\n%s", r, buf)
			}
		}()
		t.Func(t.Args)
	}()
}

func (t *Task) GetName() string {
	return t.Name
}

func (t *Task) Pause() {
	t.Running = false
	t.Scheduler.Wake()
}

func (t *Task) Resume() {
	t.Running = true
	t.PreviousRunTime = time.Now()
	t.Scheduler.Wake()
}

func (t *Task) SetScheduler(scheduler interfaces.IScheduler) {
	t.Scheduler = scheduler
}

func (t *Task) UpdateTrigger(trigger interfaces.ITrigger) {
	t.Trigger = trigger
	t.PreviousRunTime = AGScheduler.EmptyDateTime
	t.Scheduler.Wake()
	t.Logger.Info("update trigger")
}

func (t *Task) GetNextRunTime(now time.Time) time.Time {
	if t.Running {
		nextFireTime := t.Trigger.NextFireTime(t.PreviousRunTime, now)
		if nextFireTime.Equal(AGScheduler.EmptyDateTime) {
			return nextFireTime
		}
		if nextFireTime.Before(now) {
			return now.Add(-time.Duration(1))
		}
		return nextFireTime
	} else {
		return AGScheduler.MaxDateTime
	}
}
