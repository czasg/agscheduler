package AGScheduler

import (
	"github.com/sirupsen/logrus"
	"runtime"
	"time"
)

type Task struct {
	Id              int64                                                   `json:"id" pg:",pk"`
	Name            string                                                  `json:"name" pg:",use_zero"`
	Func            func(args []interface{}, kwargs map[string]interface{}) `json:"func" pg:"-"`
	Args            []interface{}                                           `json:"args" pg:",use_zero"`
	KwArgs          map[string]interface{}                                  `json:"kwargs" pg:",use_zero"`
	Scheduler       *Scheduler                                              `json:"scheduler" pg:"-"`
	Trigger         ITrigger                                                `json:"trigger" pg:"-"`
	PreviousRunTime time.Time                                               `json:"previous_run_time" pg:",use_zero"`
	NextRunTime     time.Time                                               `json:"next_run_time" pg:",use_zero"`
	Logger          *logrus.Entry                                           `json:"logger" pg:"-"`
	Running         bool                                                    `json:"running" pg:",use_zero"`
	Coalesce        bool                                                    `json:"coalesce" pg:",use_zero"`
	Count           int64                                                   `json:"cound" pg:",use_zero"`
}

func NewTask(
	name string,
	method func(args []interface{}, kwargs map[string]interface{}),
	args []interface{},
	kwargs map[string]interface{},
	trigger ITrigger,
) *Task {
	return &Task{
		Name:        name,
		Func:        method,
		Args:        args,
		KwArgs:      kwargs,
		Trigger:     trigger,
		NextRunTime: trigger.NextFireTime(EmptyDateTime, time.Now()),
		Logger: logrus.WithFields(logrus.Fields{
			"Module":   "AGScheduler.Task",
			"TaskName": name,
		}),
		Running:  true,
		Coalesce: true,
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
		t.Func(t.Args, t.KwArgs)
	}()
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

func (t *Task) UpdateTrigger(trigger ITrigger) {
	t.Trigger = trigger
	t.PreviousRunTime = EmptyDateTime
	t.Scheduler.Wake()
	t.Logger.Info("update trigger")
}

func (t *Task) GetNextFireTime(now time.Time) time.Time {
	if !t.Running {
		return MaxDateTime
	}
	t.NextRunTime = t.Trigger.NextFireTime(t.PreviousRunTime, now)
	if t.Coalesce && t.NextRunTime.Before(now) {
		t.NextRunTime = now.Add(-time.Duration(1))
	}
	return t.NextRunTime
}
