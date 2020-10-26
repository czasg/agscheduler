package AGScheduler

import (
	"github.com/robfig/cron"
	"reflect"
	"testing"
	"time"
)

func TestCronTrigger_NextFireTime(t *testing.T) {
	now := time.Now()
	cronIns, err := cron.Parse("*/5 * * * *")
	if err != nil {
		panic(err)
	}

	type fields struct {
		StartTime time.Time
		Schedule  cron.Schedule
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
			name: "StartTime empty",
			fields: fields{
				StartTime: EmptyDateTime,
				Schedule:  cronIns,
			},
			args: args{
				previous: EmptyDateTime,
				now:      now,
			},
			want: cronIns.Next(now),
		},
		{
			name: "StartTime is now",
			fields: fields{
				StartTime: now,
				Schedule:  cronIns,
			},
			args: args{
				previous: EmptyDateTime,
				now:      now,
			},
			want: now,
		},
		{
			name: "previous is now",
			fields: fields{
				StartTime: now,
				Schedule:  cronIns,
			},
			args: args{
				previous: now,
				now:      now,
			},
			want: cronIns.Next(now),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CronTrigger{
				StartTime: tt.fields.StartTime,
				Schedule:  tt.fields.Schedule,
			}
			if got := c.NextFireTime(tt.args.previous, tt.args.now); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NextFireTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDateTrigger_NextFireTime(t *testing.T) {
	now := time.Now()

	type fields struct {
		RunDateTime time.Time
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
			name: "previous is empty",
			fields: fields{
				RunDateTime: now,
			},
			args: args{
				previous: EmptyDateTime,
				now:      now,
			},
			want: now,
		},
		{
			name: "previous is not empty",
			fields: fields{
				RunDateTime: now,
			},
			args: args{
				previous: now,
				now:      now,
			},
			want: EmptyDateTime,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DateTrigger{
				RunDateTime: tt.fields.RunDateTime,
			}
			if got := d.NextFireTime(tt.args.previous, tt.args.now); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NextFireTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntervalTrigger_NextFireTime(t *testing.T) {
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
			name: "EndRunTime&previous is empty",
			fields: fields{
				Interval:     time.Duration(1),
				StartRunTime: now,
				EndRunTime:   EmptyDateTime,
			},
			args: args{
				previous: EmptyDateTime,
				now:      now,
			},
			want: now,
		},
		{
			name: "EndRunTime before StartRunTime is empty",
			fields: fields{
				Interval:     time.Second,
				StartRunTime: now,
				EndRunTime:   now.Add(-time.Second),
			},
			args: args{
				previous: EmptyDateTime,
				now:      now,
			},
			want: EmptyDateTime,
		},
		{
			name: "normal",
			fields: fields{
				Interval:     time.Second,
				StartRunTime: now,
				EndRunTime:   now.Add(time.Second * 2),
			},
			args: args{
				previous: now,
				now:      now,
			},
			want: now.Add(time.Second),
		},
		{
			name: "stop",
			fields: fields{
				Interval:     time.Second,
				StartRunTime: now,
				EndRunTime:   now.Add(time.Second * 2),
			},
			args: args{
				previous: now,
				now:      now.Add(time.Second * 3),
			},
			want: EmptyDateTime,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := IntervalTrigger{
				Interval:     tt.fields.Interval,
				StartRunTime: tt.fields.StartRunTime,
				EndRunTime:   tt.fields.EndRunTime,
			}
			if got := i.NextFireTime(tt.args.previous, tt.args.now); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NextFireTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCronTrigger(t *testing.T) {
	cronIns, err := cron.Parse("*/5 * * * *")
	if err != nil {
		panic(err)
	}

	type args struct {
		cronCmd string
	}
	tests := []struct {
		name    string
		args    args
		want    *CronTrigger
		wantErr bool
	}{
		{
			"want succ",
			args{"*/5 * * * *"},
			&CronTrigger{Schedule: cronIns},
			false,
		},
		{
			"want err",
			args{"*/5 * * * * * * * *"},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCronTrigger(tt.args.cronCmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCronTrigger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCronTrigger() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDateTrigger(t *testing.T) {
	now := time.Now()

	type args struct {
		runDateTime time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    *DateTrigger
		wantErr bool
	}{
		{
			"want succ",
			args{
				runDateTime: now,
			},
			&DateTrigger{RunDateTime: now},
			false,
		},
		{
			"want err",
			args{
				runDateTime: EmptyDateTime.Add(-time.Duration(1)),
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDateTrigger(tt.args.runDateTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDateTrigger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDateTrigger() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewIntervalTrigger(t *testing.T) {
	now := time.Now()

	type args struct {
		startTime time.Time
		endTime   time.Time
		interval  time.Duration
	}
	tests := []struct {
		name    string
		args    args
		want    *IntervalTrigger
		wantErr bool
	}{
		{
			"want succ",
			args{
				startTime: now,
				endTime:   now,
				interval:  time.Second,
			},
			&IntervalTrigger{time.Second, now, now},
			false,
		},
		{
			"want err",
			args{
				startTime: now,
				endTime:   now.Add(-time.Second),
				interval:  time.Second,
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewIntervalTrigger(tt.args.startTime, tt.args.endTime, tt.args.interval)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewIntervalTrigger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIntervalTrigger() got = %v, want %v", got, tt.want)
			}
		})
	}
}
