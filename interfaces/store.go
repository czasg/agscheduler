package interfaces

import "time"

type IStore interface {
	GetDueTasks() []ITask
	GetAllTasks() []ITask
	AddTask(task ITask) error
	DelTask(task ITask) error
	NextRunTime() (time.Time, error)
}
