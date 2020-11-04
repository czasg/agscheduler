package AGScheduler

import (
	"github.com/sirupsen/logrus"
	"reflect"
	"testing"
	"time"
)

func TestNewTask(t *testing.T) {
	cron, err := NewCronTrigger("*/5 * * * *")
	if err != nil {
		panic(err)
	}

	type args struct {
		name    string
		method  func(args ...interface{})
		args    []interface{}
		trigger ITrigger
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "ensure args",
			args: args{
				name:    "",
				method:  func(args ...interface{}) {},
				args:    []interface{}{},
				trigger: cron,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewTask(tt.args.name, tt.args.trigger, tt.args.method, tt.args.args...)
		})
	}
}

func TestTask_GetNextFireTime(t1 *testing.T) {
	now := time.Now()
	interval, _ := NewIntervalTrigger(now, EmptyDateTime, time.Second)
	date, _ := NewDateTrigger(now)

	type fields struct {
		Id              int64
		Name            string
		Func            func(args ...interface{})
		Args            []interface{}
		Scheduler       *Scheduler
		Trigger         ITrigger
		PreviousRunTime time.Time
		NextRunTime     time.Time
		Logger          *logrus.Entry
		Running         bool
		Coalesce        bool
		Count           int64
		ErrorCount      int64
	}
	type args struct {
		now time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   time.Time
	}{
		{
			name: "test-interval",
			fields: fields{
				Id:              1,
				Name:            "task",
				Func:            func(args ...interface{}) {},
				Args:            []interface{}{},
				Scheduler:       nil,
				Trigger:         interval,
				PreviousRunTime: EmptyDateTime,
				NextRunTime:     EmptyDateTime,
				Logger:          nil,
				Running:         true,
				Coalesce:        true,
				Count:           0,
				ErrorCount:      0,
			},
			args: args{
				now: now,
			},
			want: now,
		},
		{
			name: "test-date-succ",
			fields: fields{
				Id:              1,
				Name:            "task",
				Func:            func(args ...interface{}) {},
				Args:            []interface{}{},
				Scheduler:       nil,
				Trigger:         date,
				PreviousRunTime: EmptyDateTime,
				NextRunTime:     EmptyDateTime,
				Logger:          nil,
				Running:         true,
				Coalesce:        true,
				Count:           0,
				ErrorCount:      0,
			},
			args: args{
				now: now,
			},
			want: now,
		},
		{
			name: "test-date-empty",
			fields: fields{
				Id:              1,
				Name:            "task",
				Func:            func(args ...interface{}) {},
				Args:            []interface{}{},
				Scheduler:       nil,
				Trigger:         date,
				PreviousRunTime: now,
				NextRunTime:     EmptyDateTime,
				Logger:          nil,
				Running:         true,
				Coalesce:        true,
				Count:           0,
				ErrorCount:      0,
			},
			args: args{
				now: now,
			},
			want: EmptyDateTime,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Task{
				Id:              tt.fields.Id,
				Name:            tt.fields.Name,
				Func:            tt.fields.Func,
				Args:            tt.fields.Args,
				Scheduler:       tt.fields.Scheduler,
				Trigger:         tt.fields.Trigger,
				PreviousRunTime: tt.fields.PreviousRunTime,
				NextRunTime:     tt.fields.NextRunTime,
				Logger:          tt.fields.Logger,
				Running:         tt.fields.Running,
				Coalesce:        tt.fields.Coalesce,
				Count:           tt.fields.Count,
			}
			if got := t.GetNextFireTime(tt.args.now); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("GetNextFireTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTask_Go(t1 *testing.T) {
	now := time.Now()
	interval, err := NewIntervalTrigger(now, EmptyDateTime, time.Second)
	if err != nil {
		panic(err)
	}
	intChan := make(chan int)

	type fields struct {
		Id              int64
		Name            string
		Func            func(args ...interface{})
		Args            []interface{}
		Scheduler       *Scheduler
		Trigger         ITrigger
		PreviousRunTime time.Time
		NextRunTime     time.Time
		Logger          *logrus.Entry
		Running         bool
		Coalesce        bool
		Count           int64
		ErrorCount      int64
	}
	type args struct {
		runTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "test",
			fields: fields{
				Id:   1,
				Name: "task",
				Func: func(args ...interface{}) {
					iChan := args[0].(chan int)
					iChan <- 0
				},
				Args:            []interface{}{intChan},
				Scheduler:       nil,
				Trigger:         interval,
				PreviousRunTime: EmptyDateTime,
				NextRunTime:     EmptyDateTime,
				Logger:          nil,
				Running:         true,
				Coalesce:        true,
				Count:           0,
				ErrorCount:      0,
			},
			args: args{
				runTime: now,
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Task{
				Id:              tt.fields.Id,
				Name:            tt.fields.Name,
				Func:            tt.fields.Func,
				Args:            tt.fields.Args,
				Scheduler:       tt.fields.Scheduler,
				Trigger:         tt.fields.Trigger,
				PreviousRunTime: tt.fields.PreviousRunTime,
				NextRunTime:     tt.fields.NextRunTime,
				Logger:          tt.fields.Logger,
				Running:         tt.fields.Running,
				Coalesce:        tt.fields.Coalesce,
				Count:           tt.fields.Count,
			}
			t.Go(now)
			<-intChan
			if t.Count != 1 {
				t1.Errorf("func not work")
			}
		})
	}
}
