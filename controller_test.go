package AGScheduler

import (
	"context"
	"testing"
	"time"
)

func TestController_Reset(t *testing.T) {
	ctx := context.Background()
	deadline, cancel := context.WithDeadline(ctx, EmptyDateTime)

	type fields struct {
		Ctx      context.Context
		Deadline context.Context
		Cancel   context.CancelFunc
	}
	type args struct {
		deadlineTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "reset",
			fields: fields{
				Ctx:      ctx,
				Deadline: deadline,
				Cancel:   cancel,
			},
			args: args{
				deadlineTime: MaxDateTime,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Controller{
				Ctx:      tt.fields.Ctx,
				Deadline: tt.fields.Deadline,
				Cancel:   tt.fields.Cancel,
			}
			c.Reset(tt.args.deadlineTime)
			go c.Cancel()
			<-c.Deadline.Done()
		})
	}
}

func TestNewController(t *testing.T) {
	tests := []struct {
		name string
		want *Controller
	}{
		{
			name: "new",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewController()
			<-got.Deadline.Done()
		})
	}
}
