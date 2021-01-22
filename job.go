package agscheduler

import (
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

type Job struct {
	tableName      struct{}               `json:"-" pg:"ags_jobs"`
	Id             int                    `json:"id" pg:",pk"`
	Name           string                 `json:"name" pg:",use_zero,unique"`
	Task           ITask                  `json:"-" pg:"-"`
	Trigger        ITrigger               `json:"-" pg:"-"`
	Status         STATUS                 `json:"status" pg:",use_zero"`
	NotCoalesce    bool                   `json:"not_coalesce" pg:",use_zero"`
	MaxInstances   int                    `json:"max_instances" pg:",use_zero"`
	DelayGraceTime time.Duration          `json:"delay_grace_time" pg:",use_zero"`
	Scheduler      AGScheduler            `json:"-" pg:"-"`
	NextRunTime    time.Time              `json:"next_run_time" pg:",use_zero"`
	Logger         *logrus.Entry          `json:"-" pg:"-"`
	TriggerMeta    TriggerMeta            `json:"trigger_meta" pg:",use_zero"`
	TaskMeta       map[string]interface{} `json:"task_meta" pg:",use_zero"`
}

func (j *Job) FillByDefault() {
	if j.Trigger == nil {
		j.Trigger = DateTrigger{NextRunTime: time.Now()}
	}
	if j.NextRunTime.Equal(MinDateTime) {
		j.NextRunTime = j.Trigger.GetNextRunTime(MinDateTime, time.Now())
	}
	if j.Logger == nil {
		j.Logger = Log.WithFields(GenASGModule("job")).WithField("JobName", j.Name)
	}
}

func (j *Job) Run(runTimes []time.Time) {
	j.FillByDefault()
	for _, runTime := range runTimes {
		if j.DelayGraceTime > 0 && runTime.Add(j.DelayGraceTime).Before(time.Now()) {
			j.Logger.WithFields(logrus.Fields{
				"DelayGraceTime": j.DelayGraceTime,
				"RunTime":        runTime,
			}).Errorln("job's run time is out of grace time, shouldn't be scheduled.")
			continue
		}
		instance := IncreaseInstance(j.Name)
		if instance > j.MaxInstances {
			j.Logger.WithFields(logrus.Fields{"MaxInstances": j.MaxInstances + 1}).
				Warningln("Out Of MaxInstances, ignore this scheduler")
			break
		}
		go func() {
			j.Logger.Debugln("task run")
			defer func() {
				ReduceInstance(j.Name)
				if r := recover(); r != nil {
					j.Logger.WithFields(logrus.Fields{"TaskRecover": r}).
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
		nextRunTime = j.Trigger.GetNextRunTime(nextRunTime, now)
	}
	if len(runTimes) > 1 && !j.NotCoalesce {
		runTimes = runTimes[len(runTimes)-1:]
	}
	return runTimes
}
