package agscheduler

import (
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"time"
)

type CronTrigger struct {
	CronCmd      string        `json:"cron_cmd"`
	StartRunTime time.Time     `json:"start_run_time"`
	EndRunTime   time.Time     `json:"end_run_time"`
	CronIns      cron.Schedule `json:"-"`
}

func (t *CronTrigger) GetNextRunTime(previous, now time.Time) time.Time {
	if t.EndRunTime.After(now) {
		return MinDateTime
	}
	if t.StartRunTime.After(now) {
		return t.StartRunTime
	}
	if t.CronIns == nil {
		cronIns, err := cron.Parse(t.CronCmd)
		if err != nil {
			AGSLog.WithError(err).
				WithFields(GenASGModule("cron-trigger")).
				WithFields(logrus.Fields{
					"Func":    "GetNextRunTime",
					"CronCmd": t.CronCmd,
				}).Errorln("cron init err, invalid cmd.")
			return MinDateTime
		}
		t.CronIns = cronIns
	}
	if previous.Equal(MinDateTime) {
		if t.StartRunTime.Equal(MinDateTime) {
			t.StartRunTime = t.CronIns.Next(now)
		}
		return t.StartRunTime
	}
	return t.CronIns.Next(previous)
}
