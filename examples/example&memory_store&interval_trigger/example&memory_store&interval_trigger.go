package main

import (
	"context"
	"fmt"
	"github.com/czaorz/agscheduler"
	"github.com/sirupsen/logrus"
	"time"
)

type MSIntervalTask struct {
	Url    string
	Method string
}

func (m MSIntervalTask) Run(ctx context.Context) {
	fmt.Printf("[%v]%s/%s\n", time.Now(), m.Method, m.Url)
}

func main() {
	agscheduler.AGSLog.SetLevel(logrus.DebugLevel)
	job := agscheduler.Job{
		Name: "http-task",
		Trigger: agscheduler.IntervalTrigger{
			Interval: time.Second * 20,
		},
		Task: &MSIntervalTask{
			Url:    "/api",
			Method: "/GET",
		},
		DelayGraceTime: time.Second,
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
