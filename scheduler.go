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
	if ags.Context == nil {
		ags.Context = context.Background()
	}
}

func (ags *AGScheduler) Start() {
	ags.FillByDefault()
	ags.Status.SetRunning()
	go func() {
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT)
		ags.Logger.Warningln(fmt.Sprintf("receive signal[%s], exiting.", (<-ch).String()))
		_ = ags.Close()
	}()
	for {
		now := time.Now()
		jobs, err := ags.Store.GetSchedulingJobs(now)
		if err != nil {
			ags.Logger.WithError(err).Errorln("there exist an error when get jobs from store.")
		}
		for _, job := range jobs {
			runTimes := job.GetRunTimes(now)
			if len(runTimes) > 1 && job.Coalesce {
				runTimes = runTimes[len(runTimes)-1:]
			}
			job.Run(runTimes)
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
	if ags.Status.IsPaused() {
		if ags.WaitCancel == nil {
			ags.Logger.Warningln("scheduler is paused and there is not WaitCancelFunc to wake it.")
			return
		}
		ags.WaitCancel()
	}
}

func (ags *AGScheduler) Close() error {
	ags.FillByDefault()
	ags.Status.SetPaused()
	ags.Status.SetStopped()
	return nil
}
