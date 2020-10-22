package main

import (
	"fmt"
	"github.com/CzaOrz/AGScheduler"
	"github.com/CzaOrz/AGScheduler/schedulers"
	"github.com/CzaOrz/AGScheduler/stores"
	"github.com/CzaOrz/AGScheduler/tasks"
	"github.com/CzaOrz/AGScheduler/triggers"
	"os"
	"time"
)

func taskFunc(args []interface{}) {
	fmt.Println(args, time.Now())
}

func main() {
	store := stores.NewMemoryStore()
	scheduler := schedulers.NewScheduler(AGScheduler.WorksMap{}, store)

	now := time.Now()
	fmt.Println(now)

	dateTrigger1, _ := triggers.NewDateTrigger(now.Add(time.Second * 1))
	dateTrigger2, _ := triggers.NewDateTrigger(now.Add(time.Second * 2))
	dateTrigger3, _ := triggers.NewDateTrigger(now.Add(time.Second * 3))
	dateTrigger4, _ := triggers.NewDateTrigger(now.Add(time.Second * 4))
	dateTrigger5, _ := triggers.NewDateTrigger(now.Add(time.Second * 5))

	task1 := tasks.NewTask("task1", taskFunc, []interface{}{"this", "is", "task1"}, dateTrigger1)
	task2 := tasks.NewTask("task2", taskFunc, []interface{}{"this", "is", "task2"}, dateTrigger2)
	task3 := tasks.NewTask("task3", taskFunc, []interface{}{"this", "is", "task3"}, dateTrigger3)
	task4 := tasks.NewTask("task4", taskFunc, []interface{}{"this", "is", "task4"}, dateTrigger4)
	task5 := tasks.NewTask("task5", taskFunc, []interface{}{"this", "is", "task5"}, dateTrigger5)

	_ = scheduler.AddTask(task1)
	_ = scheduler.AddTask(task2)
	_ = scheduler.AddTask(task3)
	_ = scheduler.AddTask(task4)
	_ = scheduler.AddTask(task5)

	go func() {
		time.Sleep(time.Second * 10)
		now := time.Now()
		dateTrigger7, _ := triggers.NewDateTrigger(now.Add(time.Second * 1))
		dateTrigger8, _ := triggers.NewDateTrigger(now.Add(time.Second * 2))
		task6 := tasks.NewTask("task6", taskFunc, []interface{}{"this", "is", "task6"}, dateTrigger7)
		task7 := tasks.NewTask("task7", taskFunc, []interface{}{"this", "is", "task7"}, dateTrigger8)
		_ = scheduler.AddTask(task6)
		_ = scheduler.AddTask(task7)

		time.Sleep(time.Second * 5)
		os.Exit(0)
	}()

	scheduler.Start()
}
