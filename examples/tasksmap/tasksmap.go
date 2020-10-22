package main

import (
	"fmt"
	"github.com/CzaOrz/AGScheduler"
	"github.com/CzaOrz/AGScheduler/schedulers"
	"github.com/CzaOrz/AGScheduler/stores"
	"github.com/CzaOrz/AGScheduler/triggers"
	"os"
	"time"
)

var TasksMap = AGScheduler.WorksMap{
	"task1": AGScheduler.WorkDetail{
		Func: func(args []interface{}) {
			fmt.Println(args, time.Now())
		},
		Args: []interface{}{"this", "is", "task1"},
	},
	"task2": AGScheduler.WorkDetail{
		Func: func(args []interface{}) {
			fmt.Println(args, time.Now())
		},
		Args: []interface{}{"this", "is", "task2"},
	},
	"task3": AGScheduler.WorkDetail{
		Func: func(args []interface{}) {
			fmt.Println(args, time.Now())
		},
		Args: []interface{}{"this", "is", "task3"},
	},
	"task4": AGScheduler.WorkDetail{
		Func: func(args []interface{}) {
			fmt.Println(args, time.Now())
		},
		Args: []interface{}{"this", "is", "task4"},
	},
	"task5": AGScheduler.WorkDetail{
		Func: func(args []interface{}) {
			fmt.Println(args, time.Now())
		},
		Args: []interface{}{"this", "is", "task5"},
	},
}

func main() {
	store := stores.NewMemoryStore()
	scheduler := schedulers.NewScheduler(TasksMap, store)

	now := time.Now()
	fmt.Println(now)
	dateTrigger1, _ := triggers.NewDateTrigger(now.Add(time.Second * 1))
	dateTrigger2, _ := triggers.NewDateTrigger(now.Add(time.Second * 2))
	dateTrigger3, _ := triggers.NewDateTrigger(now.Add(time.Second * 3))
	dateTrigger4, _ := triggers.NewDateTrigger(now.Add(time.Second * 4))
	dateTrigger5, _ := triggers.NewDateTrigger(now.Add(time.Second * 5))

	_ = scheduler.AddTaskFromTasksMap("t-task1", "task1", []interface{}{}, dateTrigger1)
	_ = scheduler.AddTaskFromTasksMap("t-task2", "task2", []interface{}{}, dateTrigger2)
	_ = scheduler.AddTaskFromTasksMap("t-task3", "task3", []interface{}{}, dateTrigger3)
	_ = scheduler.AddTaskFromTasksMap("t-task4", "task4", []interface{}{"new", "args", "task4"}, dateTrigger4)
	_ = scheduler.AddTaskFromTasksMap("t-task5", "task5", []interface{}{"new", "args", "task5"}, dateTrigger5)

	go func() {
		time.Sleep(time.Second * 10)
		os.Exit(0)
	}()

	scheduler.Start()
}
