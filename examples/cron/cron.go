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

	cron1, _ := triggers.NewCronTrigger("*/5 * * * *")

	task1 := tasks.NewTask("task1", taskFunc, []interface{}{"this", "is", "task1"}, cron1)
	_ = scheduler.AddTask(task1)

	go func() {
		time.Sleep(time.Second * 60)
		os.Exit(0)
	}()

	scheduler.Start()
}
