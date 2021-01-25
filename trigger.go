package agscheduler

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"time"
)

var (
	MinDateTime = time.Time{}
	MaxDateTime = time.Now().Add(time.Duration(math.MaxInt64))
)

// GetNextRunTime: if result is MinDateTime, it mean this Job is over, should be delete.
type ITrigger interface {
	GetNextRunTime(previous, now time.Time) time.Time
}

type TriggerMeta struct {
	Type         string        `json:"type"`
	NextRunTime  time.Time     `json:"next_run_time"`
	Interval     time.Duration `json:"interval"`
	StartRunTime time.Time     `json:"start_run_time"`
	EndRunTime   time.Time     `json:"end_run_time"`
	CronCmd      string        `json:"cron_cmd"`
}

func SerializeTrigger(job *Job) error {
	if job.Trigger == nil {
		return errors.New("trigger is nil.")
	}
	body, err := json.Marshal(job.Trigger)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &job.TriggerMeta)
	if err != nil {
		return err
	}
	taskT := reflect.TypeOf(job.Trigger)
	if taskT.Kind() == reflect.Ptr {
		taskT = taskT.Elem()
	}
	job.TriggerMeta.Type = taskT.String()
	return nil
}

func DeserializeTrigger(job *Job) error {
	switch job.TriggerMeta.Type {
	case "agscheduler.DateTrigger":
		job.Trigger = &DateTrigger{
			NextRunTime: job.TriggerMeta.NextRunTime,
		}
	case "agscheduler.IntervalTrigger":
		job.Trigger = &IntervalTrigger{
			Interval:     time.Minute,
			StartRunTime: job.TriggerMeta.StartRunTime,
			EndRunTime:   job.TriggerMeta.EndRunTime,
		}
	case "agscheduler.CronTrigger":
		job.Trigger = &CronTrigger{
			CronCmd:      job.TriggerMeta.CronCmd,
			StartRunTime: job.TriggerMeta.StartRunTime,
			EndRunTime:   job.TriggerMeta.EndRunTime,
		}
	default:
		return fmt.Errorf("trigger[%s] is not define.", job.TriggerMeta.Type)
	}
	return nil
}
