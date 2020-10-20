package interfaces

type IScheduler interface {
	Start()
	AddTask(task ITask) error
	AddTaskFromTasksMap(name, taskMapKey string, args []interface{}, trigger ITrigger) error
	DelTask(task ITask) error
}
