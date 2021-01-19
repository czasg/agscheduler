package agscheduler

import (
	"github.com/robfig/cron"
	"reflect"
	"testing"
	"time"
)

func TestCronTrigger_GetNextRunTime(t1 *testing.T) {
	now := time.Now()
	cronCmd := "* * * * *"
	cronIns, _ := cron.Parse(cronCmd)
	type fields struct {
		CronCmd      string
		StartRunTime time.Time
		EndRunTime   time.Time
		CronIns      cron.Schedule
	}
	type args struct {
		previous time.Time
		now      time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   time.Time
	}{
		{
			name:   "out of end run time",
			fields: fields{},
			args: args{
				now: MinDateTime,
			},
			want: MinDateTime,
		},
		{
			name: "out of start run time",
			fields: fields{
				StartRunTime: now,
			},
			args: args{
				now: MinDateTime,
			},
			want: now,
		},
		{
			name: "cron cmd err",
			fields: fields{
				CronCmd: "test",
			},
			args: args{
				now: now,
			},
			want: MinDateTime,
		},
		{
			name: "previous empty",
			fields: fields{
				CronCmd:      "* * * * *",
				StartRunTime: now,
			},
			args: args{
				now: now,
			},
			want: now,
		},
		{
			name: "previous & startRunTime empty",
			fields: fields{
				CronCmd: cronCmd,
			},
			args: args{
				now: now,
			},
			want: cronIns.Next(now),
		},
		{
			name: "pass",
			fields: fields{
				CronCmd: cronCmd,
			},
			args: args{
				previous: now,
			},
			want: cronIns.Next(now),
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &CronTrigger{
				CronCmd:      tt.fields.CronCmd,
				StartRunTime: tt.fields.StartRunTime,
				EndRunTime:   tt.fields.EndRunTime,
				CronIns:      tt.fields.CronIns,
			}
			if got := t.GetNextRunTime(tt.args.previous, tt.args.now); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("GetNextRunTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
