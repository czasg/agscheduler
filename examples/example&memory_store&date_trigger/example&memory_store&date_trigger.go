package main

import (
	"context"
	"fmt"
	"github.com/czaorz/agscheduler"
	"github.com/sirupsen/logrus"
	"time"
)

type MSDateTask struct {
	Url    string
	Method string
}

func (m MSDateTask) Run(ctx context.Context) {
	fmt.Printf("[%v]%s/%s\n", time.Now(), m.Method, m.Url)
}

func main() {
	agscheduler.AGSLog.SetLevel(logrus.DebugLevel)
	now := time.Now()
	fmt.Println(now)
	job := agscheduler.Job{
		Name: "http-task",
		Trigger: agscheduler.DateTrigger{
			NextRunTime: now.Add(time.Second * 30),
		},
		Task: &MSDateTask{
			Url:    "/api",
			Method: "/GET",
		},
	}
	scheduler := agscheduler.AGScheduler{}
	err := scheduler.AddJob(&job)
	if err != nil {
		panic(err)
	}
	jobObj, err := scheduler.GetJobByJobName("http-task")
	if err != nil {
		panic(err)
	}
	fmt.Println(jobObj.Name, jobObj.Trigger, jobObj.Task)
	scheduler.Start()
}
