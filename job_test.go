package agscheduler

import (
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestJob_FillByDefault(t *testing.T) {
	type fields struct {
		Id           int
		Name         string
		Task         ITask
		Trigger      ITrigger
		Status       STATUS
		Coalesce     bool
		MaxInstances int
		Scheduler    AGScheduler
		NextRunTime  time.Time
		Logger       *logrus.Entry
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "empty",
			fields: fields{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Job{
				Id:           tt.fields.Id,
				Name:         tt.fields.Name,
				Task:         tt.fields.Task,
				Trigger:      tt.fields.Trigger,
				Status:       tt.fields.Status,
				Coalesce:     tt.fields.Coalesce,
				MaxInstances: tt.fields.MaxInstances,
				Scheduler:    tt.fields.Scheduler,
				NextRunTime:  tt.fields.NextRunTime,
				Logger:       tt.fields.Logger,
			}
			j.FillByDefault()
			logStr, _ := j.Logger.String()
			if logStr != "time=\"0001-01-01T00:00:00Z\" level=panic AGSVersion=0.0.1 ASGModule=job\n" {
				t.Errorf("GetRunTimes() = %v, want %v", logStr, "time=\"0001-01-01T00:00:00Z\" level=panic AGSVersion=0.0.1 ASGModule=job\n")
			}
		})
	}
}

func TestJob_GetRunTimes(t *testing.T) {
	now := time.Now()
	type fields struct {
		Id           int
		Name         string
		Task         ITask
		Trigger      ITrigger
		Status       STATUS
		Coalesce     bool
		MaxInstances int
		Scheduler    AGScheduler
		NextRunTime  time.Time
		Logger       *logrus.Entry
	}
	type args struct {
		now time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []time.Time
	}{
		{
			name: "pass",
			fields: fields{
				Trigger: &DateTrigger{NextRunTime: now.Add(time.Second)},
			},
			args: args{
				now: now,
			},
			want: []time.Time{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Job{
				Trigger: tt.fields.Trigger,
			}
			j.FillByDefault()
			j.GetRunTimes(now) // todo, System Panic
			//j := &Job{
			//	Id:           tt.fields.Id,
			//	Name:         tt.fields.Name,
			//	Task:         tt.fields.Task,
			//	Trigger:      tt.fields.Trigger,
			//	Status:       tt.fields.Status,
			//	Coalesce:     tt.fields.Coalesce,
			//	MaxInstances: tt.fields.MaxInstances,
			//	Scheduler:    tt.fields.Scheduler,
			//	NextRunTime:  tt.fields.NextRunTime,
			//	Logger:       tt.fields.Logger,
			//}
			//fmt.Println(j.GetRunTimes(now))
			//if got := j.GetRunTimes(tt.args.now); !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("GetRunTimes() = %v, want %v", got, tt.want)
			//}
		})
	}
}
