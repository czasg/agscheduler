package agscheduler

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	AGSContext = context.Background()
)

type AGScheduler struct {
	Store  IStore
	Logger *logrus.Entry
	Status STATUS
	// context.
	Context    context.Context
	WaitCancel context.CancelFunc
}

func (ags *AGScheduler) FillByDefault() {
	if ags.Store == nil {
		ags.Store = &MemoryStore{}
	}
	if ags.Logger == nil {
		ags.Logger = Log.WithFields(GenASGModule("scheduler"))
	}
	if ags.Status == "" {
		ags.Status.SetPaused()
	}
	if ags.Context == nil {
		ags.Context = AGSContext
	}
}

func (ags *AGScheduler) listenSignal() {
	ags.FillByDefault()
	ags.Status.SetRunning()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGTERM, syscall.SIGQUIT)
	ags.Logger.Warningln(fmt.Sprintf("receive signal[%s], exiting.", (<-ch).String()))
	_ = ags.Close()
}

func (ags *AGScheduler) Start() {
	ags.FillByDefault()
	go ags.listenSignal()
	for {
		now := time.Now()
		jobs, err := ags.Store.GetSchedulingJobs(now)
		if err != nil {
			ags.Logger.WithError(err).Errorln("there exist an error when get jobs from store.")
		}
		for _, job := range jobs {
			job.FillByDefault()
			runTimes := job.GetRunTimes(now)
			if len(runTimes) > 0 {
				job.Run(runTimes)
				job.NextRunTime = job.Trigger.GetNextRunTime(runTimes[len(runTimes)-1], now)
				if job.NextRunTime.Equal(MinDateTime) {
					if err = ags.Store.DelJob(job); err != nil {
						ags.Logger.WithError(err).Errorln("there exist an error when delete jobs.")
					}
					continue
				}
				if err = ags.Store.UpdateJob(job); err != nil {
					ags.Logger.WithError(err).Errorln("there exist an error when update jobs.")
				}
			}
		}
		nextRunTime, err := ags.Store.GetNextRunTime()
		if err != nil {
			ags.Logger.WithError(err).Errorln("there exist an error when get next run time.")
		}
		ags.WaitWithTime(nextRunTime)
		if ags.Status.IsPaused() {
			ags.Logger.Warningln("scheduler paused.")
			ags.WaitWithTime(MaxDateTime)
		}
		if ags.Status.IsStopped() {
			break
		}
	}
	ags.Logger.WithFields(logrus.Fields{
		"Status": ags.Status,
	}).Info("scheduler over.")
}

func (ags *AGScheduler) WaitWithTime(waitTime time.Time) {
	ctx, cancel := context.WithDeadline(ags.Context, waitTime)
	ags.WaitCancel = cancel
	<-ctx.Done()
}

func (ags *AGScheduler) Pause() {
	ags.FillByDefault()
	ags.Logger.Warningln("scheduler is pausing.")
	ags.Status.SetPaused()
	ags.Wake()
}

func (ags *AGScheduler) Wake() {
	ags.FillByDefault()
	if ags.WaitCancel == nil {
		return
	}
	ags.WaitCancel()
}

func (ags *AGScheduler) Close() error {
	ags.FillByDefault()
	ags.Status.SetStopped()
	defer ags.Wake()
	return nil
}

func (ags *AGScheduler) AddJob(jobs ...*Job) (err error) {
	ags.FillByDefault()
	done := []*Job{}
	for _, job := range jobs {
		job.FillByDefault()
		err = ags.Store.AddJob(job)
		if err != nil {
			_ = ags.DelJob(done...)
			return err
		}
		done = append(done, job)
	}
	defer ags.Wake()
	return
}

func (ags *AGScheduler) DelJob(jobs ...*Job) (err error) {
	ags.FillByDefault()
	for _, job := range jobs {
		job.FillByDefault()
		err = ags.Store.DelJob(job)
		if err != nil {
			return err
		}
	}
	defer ags.Wake()
	return
}

func (ags *AGScheduler) UpdateJob(jobs ...*Job) (err error) {
	ags.FillByDefault()
	for _, job := range jobs {
		job.FillByDefault()
		err = ags.Store.UpdateJob(job)
		if err != nil {
			return err
		}
	}
	defer ags.Wake()
	return
}

func (ags *AGScheduler) GetAllJobs() (jobs []*Job, err error) {
	ags.FillByDefault()
	jobs, err = ags.Store.GetAllJobs()
	if err != nil {
		return
	}
	for _, job := range jobs {
		job.FillByDefault()
	}
	return
}

func (ags *AGScheduler) GetJobByJobName(jobName string) (job *Job, err error) {
	ags.FillByDefault()
	job, err = ags.Store.GetJobByName(jobName)
	if err != nil {
		return
	}
	job.FillByDefault()
	return
}
