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

	trigger1, _ := triggers.NewIntervalTrigger(now.Add(time.Second*1), now.Add(time.Second*10), time.Second*5)
	trigger2, _ := triggers.NewIntervalTrigger(now.Add(time.Second*1), now.Add(time.Second*20), time.Second*5)
	trigger3, _ := triggers.NewIntervalTrigger(now.Add(time.Second*1), now.Add(time.Second*30), time.Second*5)
	trigger4, _ := triggers.NewIntervalTrigger(now.Add(time.Second*1), now.Add(time.Second*40), time.Second*5)
	trigger5, _ := triggers.NewIntervalTrigger(now.Add(time.Second*1), now.Add(time.Second*50), time.Second*5)

	task1 := tasks.NewTask("task1", taskFunc, []interface{}{"this", "is", "task1"}, trigger1)
	task2 := tasks.NewTask("task2", taskFunc, []interface{}{"this", "is", "task2"}, trigger2)
	task3 := tasks.NewTask("task3", taskFunc, []interface{}{"this", "is", "task3"}, trigger3)
	task4 := tasks.NewTask("task4", taskFunc, []interface{}{"this", "is", "task4"}, trigger4)
	task5 := tasks.NewTask("task5", taskFunc, []interface{}{"this", "is", "task5"}, trigger5)

	_ = scheduler.AddTask(task1)
	_ = scheduler.AddTask(task2)
	_ = scheduler.AddTask(task3)
	_ = scheduler.AddTask(task4)
	_ = scheduler.AddTask(task5)

	go func() {
		time.Sleep(time.Second * 60)
		now := time.Now()
		trigger6, _ := triggers.NewIntervalTrigger(now.Add(time.Second*1), now.Add(time.Second*10), time.Second*5)
		trigger7, _ := triggers.NewIntervalTrigger(now.Add(time.Second*1), now.Add(time.Second*20), time.Second*5)
		task6 := tasks.NewTask("task6", taskFunc, []interface{}{"this", "is", "task6"}, trigger6)
		task7 := tasks.NewTask("task7", taskFunc, []interface{}{"this", "is", "task7"}, trigger7)
		_ = scheduler.AddTask(task6)
		_ = scheduler.AddTask(task7)

		time.Sleep(time.Second * 30)
		os.Exit(0)
	}()

	scheduler.Start()
}
