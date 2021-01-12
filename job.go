package agscheduler

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

type ITask interface {
	Run()
}

type Job struct {
	Id           int      `json:"id" pg:",pk"`
	Name         string   `json:"name" pg:",use_zero"`
	Task         ITask    `json:"-" pg:"-"`
	Trigger      ITrigger `json:"-" pg:"-"`
	Status       STATUS   `json:"status" pg:",use_zero"`
	Coalesce     bool     `json:"coalesce" pg:",use_zero"`
	MaxInstances int      `json:"max_instances" pg:",use_zero"`
	/* should be init by AGS. */
	Scheduler   AGScheduler   `json:"-" pg:"-"`
	NextRunTime time.Time     `json:"next_run_time" pg:",use_zero"`
	Logger      *logrus.Entry `json:"-" pg:"-"`
}

func (j *Job) FillByDefault() {
	if j.Logger == nil {
		j.Logger = Log.WithFields(GenASGModule("job"))
	}
}

func (j *Job) GetRunTimes(now time.Time) []time.Time {
	runTimes := []time.Time{}
	nextRunTime := j.NextRunTime
	for {
		nextRunTime := j.Trigger.GetNextRunTime(nextRunTime, now)
		fmt.Println(nextRunTime)
		if nextRunTime.Equal(MinDateTime) {
			break
		}
		runTimes = append(runTimes, nextRunTime)
	}
	return runTimes
}
