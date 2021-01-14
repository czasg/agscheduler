package agscheduler

import (
	"context"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

var (
	MaxTolerance = 1024
	InstanceLock = sync.Mutex{}
	InstanceMap  = map[string]int{}
)

func IncreaseInstance(key string) int {
	InstanceLock.Lock()
	defer InstanceLock.Unlock()
	instance, ok := InstanceMap[key]
	if !ok {
		instance = 0
	}
	InstanceMap[key] = instance + 1
	return instance
}

func ReduceInstance(key string) int {
	InstanceLock.Lock()
	defer InstanceLock.Unlock()
	instance, ok := InstanceMap[key]
	if !ok {
		return 0
	}
	InstanceMap[key] = instance - 1
	return instance
}

func DeleteInstance(key string) {
	InstanceLock.Lock()
	defer InstanceLock.Unlock()
	delete(InstanceMap, key)
}

type ITask interface {
	Run(ctx context.Context)
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
		j.Logger = Log.WithFields(GenASGModule("job")).WithField("JobName", j.Name)
	}
}

func (j *Job) Run(runTimes []time.Time) {
	j.FillByDefault()
	for _, runTime := range runTimes {
		logger := j.Logger.WithFields(logrus.Fields{"RunTime": runTime})
		instance := IncreaseInstance(j.Name)
		if instance > j.MaxInstances {
			logger.WithFields(logrus.Fields{"MaxInstances": j.MaxInstances + 1}).
				Warningln("Out Of MaxInstances, ignore this scheduler")
			break
		}
		go func() {
			logger.Debugln("task run")
			defer func() {
				ReduceInstance(j.Name)
				if r := recover(); r != nil {
					logger.WithFields(logrus.Fields{"TaskRecover": r}).
						Errorln("Panic! please ensure task right.")
				}
			}()
			j.Task.Run(j.Scheduler.Context)
		}()
	}
}

func (j *Job) GetRunTimes(now time.Time) []time.Time {
	j.FillByDefault()
	var (
		runTimes    = []time.Time{}
		nextRunTime = j.NextRunTime
		tolerance   = 0
	)
	for {
		nextRunTime = j.Trigger.GetNextRunTime(nextRunTime, now)
		if nextRunTime.After(now) {
			break
		}
		if nextRunTime.Equal(MinDateTime) {
			break
		}
		tolerance++
		if tolerance >= MaxTolerance {
			j.Logger.Warningln("Abnormal RunTimes, Please Ensure Trigger Right.")
			break
		}
		runTimes = append(runTimes, nextRunTime)
	}
	return runTimes
}
