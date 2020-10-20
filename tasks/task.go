package tasks

import (
	"github.com/CzaOrz/AGScheduler/interfaces"
	"time"
)

type Task struct {
	Name      string
	Func      func(args []interface{})
	Args      []interface{}
	Scheduler interfaces.IScheduler
	Trigger   interfaces.ITrigger
	Store     interfaces.IStore
}

func NewTask(name string, method func(args []interface{}), args []interface{}, trigger interfaces.ITrigger) *Task {
	return &Task{
		Name:    name,
		Func:    method,
		Args:    args,
		Trigger: trigger,
	}
}
func (t *Task) Go() {
	go func() {
		t.Func(t.Args)
	}()
}
func (t *Task) Pause() error {
	return nil
}
func (t *Task) Resume() error {
	return nil
}
func (t *Task) GetName() string {
	return t.Name
}
func (t *Task) GetRunTimes(now time.Time) []time.Time {
	return []time.Time{}
}
func (t *Task) SetScheduler(scheduler interfaces.IScheduler) {
	return
}
