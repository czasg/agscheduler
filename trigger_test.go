package agscheduler

import (
	"fmt"
	"testing"
	"time"
)

func TestDeserializeTrigger(t *testing.T) {
	type args struct {
		job *Job
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "date pass",
			args: args{
				job: &Job{
					Name: "test",
					TriggerMeta: TriggerMeta{
						Type:        "agscheduler.DateTrigger",
						NextRunTime: time.Now(),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "interval pass",
			args: args{
				job: &Job{
					Name: "test",
					TriggerMeta: TriggerMeta{
						Type:         "agscheduler.IntervalTrigger",
						NextRunTime:  time.Now(),
						Interval:     time.Minute,
						StartRunTime: MinDateTime,
						EndRunTime:   MinDateTime,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "cron pass",
			args: args{
				job: &Job{
					Name: "test",
					TriggerMeta: TriggerMeta{
						Type:        "agscheduler.CronTrigger",
						NextRunTime: time.Now(),
						CronCmd:     "* * * * *",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeserializeTrigger(tt.args.job); (err != nil) != tt.wantErr {
				t.Errorf("DeserializeTrigger() error = %v, wantErr %v", err, tt.wantErr)
			}
			fmt.Println(tt.args.job.Trigger)
		})
	}
}

func TestSerializeTrigger(t *testing.T) {
	now := time.Now()
	type args struct {
		job *Job
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "date pass",
			args: args{
				job: &Job{
					Name: "test",
					Trigger: &DateTrigger{
						NextRunTime: now,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "interval pass",
			args: args{
				job: &Job{
					Name: "test",
					Trigger: &IntervalTrigger{
						Interval: time.Minute,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "cron pass",
			args: args{
				job: &Job{
					Name: "test",
					Trigger: &CronTrigger{
						CronCmd: "* * * * *",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SerializeTrigger(tt.args.job); (err != nil) != tt.wantErr {
				t.Errorf("SerializeTrigger() error = %v, wantErr %v", err, tt.wantErr)
			}
			fmt.Println(tt.args.job.TriggerMeta)
		})
	}
}
