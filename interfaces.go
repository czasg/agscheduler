package AGScheduler

import "time"

type IStore interface {
	GetDueTasks(now time.Time) []*Task
	GetTaskByName(name string) (*Task, error)
	GetAllTasks() []*Task
	AddTask(task *Task) error
	DelTask(task *Task) error
	UpdateTask(task *Task) error
	GetNextRunTime() time.Time
}

type ITrigger interface {
	NextFireTime(previous, now time.Time) time.Time
	GetTriggerState() TriggerState
}
