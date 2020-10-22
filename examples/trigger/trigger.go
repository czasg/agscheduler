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

	cronTrigger, _ := triggers.NewCronTrigger("*/5 * * * *")
	intervalTrigger, _ := triggers.NewIntervalTrigger(now.Add(time.Second), AGScheduler.EmptyDateTime, time.Second*6)
	dateTrigger, _ := triggers.NewDateTrigger(now.Add(time.Second * 10))

	task1 := tasks.NewTask("task1", taskFunc, []interface{}{"this", "is", "task1"}, cronTrigger)
	task2 := tasks.NewTask("task2", taskFunc, []interface{}{"this", "is", "task2"}, intervalTrigger)
	task3 := tasks.NewTask("task3", taskFunc, []interface{}{"this", "is", "task3"}, dateTrigger)
	_ = scheduler.AddTask(task1)
	_ = scheduler.AddTask(task2)
	_ = scheduler.AddTask(task3)

	go func() {
		time.Sleep(time.Second * 60)
		os.Exit(0)
	}()

	scheduler.Start()
}
