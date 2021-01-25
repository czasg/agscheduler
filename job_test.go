package agscheduler

import (
	"github.com/sirupsen/logrus"
	"reflect"
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
		NotCoalesce  bool
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
				NotCoalesce:  tt.fields.NotCoalesce,
				MaxInstances: tt.fields.MaxInstances,
				Scheduler:    tt.fields.Scheduler,
				NextRunTime:  tt.fields.NextRunTime,
				Logger:       tt.fields.Logger,
			}
			j.FillByDefault()
			logStr, _ := j.Logger.String()
			if logStr != "time=\"0001-01-01T00:00:00Z\" level=panic AGSVersion=0.0.1 ASGModule=job JobName=\n" {
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
		NotCoalesce  bool
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
			name: "date empty",
			fields: fields{
				Trigger: &DateTrigger{NextRunTime: now.Add(time.Second)},
			},
			args: args{
				now: now,
			},
			want: []time.Time{},
		},
		{
			name: "date pass",
			fields: fields{
				Trigger: &DateTrigger{NextRunTime: now.Add(time.Second)},
			},
			args: args{
				now: now.Add(time.Minute),
			},
			want: []time.Time{now.Add(time.Second)},
		},
		{
			name: "interval now",
			fields: fields{
				Trigger: &IntervalTrigger{
					Interval: time.Minute,
				},
			},
			args: args{
				now: now,
			},
			want: []time.Time{now.Add(-time.Nanosecond)},
		},
		{
			name: "interval pass with NotCoalesce",
			fields: fields{
				Trigger: &IntervalTrigger{
					Interval:     time.Second,
					StartRunTime: now.Add(-time.Second * 2),
				},
				NotCoalesce: true,
			},
			args: args{
				now: now,
			},
			want: []time.Time{now.Add(-time.Second * 2), now.Add(-time.Second), now},
		},
		{
			name: "interval pass with coalesce",
			fields: fields{
				Trigger: &IntervalTrigger{
					Interval:     time.Second,
					StartRunTime: now.Add(-time.Second * 2),
				},
			},
			args: args{
				now: now,
			},
			want: []time.Time{now},
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
				NotCoalesce:  tt.fields.NotCoalesce,
				MaxInstances: tt.fields.MaxInstances,
				Scheduler:    tt.fields.Scheduler,
				NextRunTime:  tt.fields.NextRunTime,
				Logger:       tt.fields.Logger,
			}
			if got := j.GetRunTimes(tt.args.now); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRunTimes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteInstance(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "pass",
			args: args{
				key: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DeleteInstance(tt.args.key)
		})
	}
}

func TestIncreaseInstance(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "pass",
			args: args{
				key: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IncreaseInstance(tt.args.key); got != tt.want {
				t.Errorf("IncreaseInstance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJob_Run(t *testing.T) {
	now := time.Now()
	type fields struct {
		tableName      struct{}
		Id             int
		Name           string
		Task           ITask
		Trigger        ITrigger
		Status         STATUS
		NotCoalesce    bool
		MaxInstances   int
		DelayGraceTime time.Duration
		Scheduler      AGScheduler
		NextRunTime    time.Time
		Logger         *logrus.Entry
		TriggerMeta    TriggerMeta
		TaskMeta       map[string]interface{}
	}
	type args struct {
		runTimes []time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "pass",
			fields: fields{
				Name: "test",
				Task: TestTask{
					Name: "test-name",
					Age:  1,
				},
			},
			args: args{
				runTimes: []time.Time{now},
			},
		},
		{
			name: "out of instance",
			fields: fields{
				Name: "test",
				Task: TestTask{
					Name: "test-name",
					Age:  1,
				},
				DelayGraceTime: time.Duration(1),
			},
			args: args{
				runTimes: []time.Time{now},
			},
		},
		{
			name: "out of grace",
			fields: fields{
				Name: "test1",
				Task: TestTask{
					Name: "test-name",
					Age:  1,
				},
				DelayGraceTime: time.Duration(1),
			},
			args: args{
				runTimes: []time.Time{now},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Job{
				tableName:      tt.fields.tableName,
				Id:             tt.fields.Id,
				Name:           tt.fields.Name,
				Task:           tt.fields.Task,
				Trigger:        tt.fields.Trigger,
				Status:         tt.fields.Status,
				NotCoalesce:    tt.fields.NotCoalesce,
				MaxInstances:   tt.fields.MaxInstances,
				DelayGraceTime: tt.fields.DelayGraceTime,
				Scheduler:      tt.fields.Scheduler,
				NextRunTime:    tt.fields.NextRunTime,
				Logger:         tt.fields.Logger,
				TriggerMeta:    tt.fields.TriggerMeta,
				TaskMeta:       tt.fields.TaskMeta,
			}
			j.Run(tt.args.runTimes)
		})
	}
}

func TestReduceInstance(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "pass",
			args: args{
				key: "test-only-one",
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReduceInstance(tt.args.key); got != tt.want {
				t.Errorf("ReduceInstance() = %v, want %v", got, tt.want)
			}
		})
	}
}
