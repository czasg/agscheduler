package interfaces

import "time"

type IStore interface {
	GetDueTasks(now time.Time) []ITask
	GetAllTasks() []ITask
	GetTask(name string) (ITask, error)
	AddTask(task ITask) error
	DelTask(task ITask) error
	UpdateTask(task ITask, now time.Time) error
	GetNextRunTime(now time.Time) time.Time
}
