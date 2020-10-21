package tasks

import (
	"github.com/CzaOrz/AGScheduler/interfaces"
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
}

func NewTask(name string, method func(args []interface{}), args []interface{}, trigger interfaces.ITrigger) *Task {
	return &Task{
		Name:    name,
		Func:    method,
		Args:    args,
		Trigger: trigger,
	}
}

func (t *Task) Go(runTime time.Time) {
	t.PreviousRunTime = runTime
	go t.Func(t.Args)
}

func (t *Task) GetName() string {
	return t.Name
}

func (t *Task) GetNextRunTime(now time.Time) time.Time {
	return t.Trigger.NextFireTime(t.PreviousRunTime, now)
}
