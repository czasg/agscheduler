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

func main() {
	scheduler := schedulers.NewScheduler(AGScheduler.WorksMap{}, stores.NewMemoryStore())

	dateTrigger, _ := triggers.NewDateTrigger(time.Now().Add(time.Hour))
	task := tasks.NewTask("task1", func(args []interface{}) {}, []interface{}{}, dateTrigger)
	_ = scheduler.AddTask(task)

	go func() {
		time.Sleep(time.Second * 5)
		dateTrigger, _ := triggers.NewDateTrigger(time.Now().Add(time.Second * 5))
		task := tasks.NewTask("task2", func(args []interface{}) {
			fmt.Println(args)
		}, []interface{}{"this", "is", "task2"}, dateTrigger)
		_ = scheduler.AddTask(task)

		time.Sleep(time.Second * 10)
		os.Exit(0)
	}()

	scheduler.Start()
}
