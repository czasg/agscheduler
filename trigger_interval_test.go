package agscheduler

import (
	"reflect"
	"testing"
	"time"
)

func TestIntervalTrigger_GetNextRunTime(t1 *testing.T) {
	now := time.Now()
	type fields struct {
		Interval     time.Duration
		StartRunTime time.Time
		EndRunTime   time.Time
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
			name:   "test end time",
			fields: fields{},
			args:   args{},
			want:   MinDateTime,
		},
		{
			name: "test start time",
			fields: fields{
				StartRunTime: MaxDateTime,
			},
			args: args{},
			want: MaxDateTime,
		},
		{
			name:   "test equal time",
			fields: fields{},
			args: args{
				now: now,
			},
			want: now,
		},
		{
			name: "test pass",
			fields: fields{
				Interval: time.Second,
			},
			args: args{
				previous: now,
				now:      now,
			},
			want: now.Add(time.Second),
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := IntervalTrigger{
				Interval:     tt.fields.Interval,
				StartRunTime: tt.fields.StartRunTime,
				EndRunTime:   tt.fields.EndRunTime,
			}
			if got := t.GetNextRunTime(tt.args.previous, tt.args.now); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("GetNextRunTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
