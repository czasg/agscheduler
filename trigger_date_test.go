package agscheduler

import (
	"reflect"
	"testing"
	"time"
)

func TestDateTrigger_GetNextRunTime(t1 *testing.T) {
	type fields struct {
		NextRunTime time.Time
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
			name:   "empty",
			fields: fields{},
			args:   args{},
			want:   MinDateTime,
		},
		{
			name: "pass",
			fields: fields{
				NextRunTime: MaxDateTime,
			},
			args: args{},
			want: MaxDateTime,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := DateTrigger{
				NextRunTime: tt.fields.NextRunTime,
			}
			if got := t.GetNextRunTime(tt.args.previous, tt.args.now); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("GetNextRunTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
