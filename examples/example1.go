package main

import (
	"fmt"
	"github.com/CzaOrz/AGScheduler"
	"github.com/CzaOrz/AGScheduler/schedulers"
	"github.com/CzaOrz/AGScheduler/tasks"
	"github.com/CzaOrz/AGScheduler/triggers"
	"time"
)

var TasksMap = AGScheduler.WorksMap{
	"task1": AGScheduler.WorkDetail{
		Func: func(args []interface{}) {
			for index, value := range args {
				println(index, value)
			}
		},
		Args: []interface{}{"this", "is", "task1"},
	},
}

func taskFunc(args []interface{}) {
	for index, value := range args {
		fmt.Println(index, value)
	}
}

func main() {
	trigger := triggers.NewDateTrigger(time.Now())
	task2 := tasks.NewTask("task2", taskFunc, []interface{}{"this", "is", "task2"}, trigger)
	scheduler := schedulers.NewScheduler(TasksMap)

	_ = scheduler.AddTask(task2)
	scheduler.Start()
}
