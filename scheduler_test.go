package agscheduler

import (
	"context"
	"github.com/sirupsen/logrus"
	"reflect"
	"testing"
	"time"
)

func TestAGScheduler_AddJob(t *testing.T) {
	type fields struct {
		Store      IStore
		Logger     *logrus.Entry
		Status     STATUS
		Context    context.Context
		WaitCancel context.CancelFunc
	}
	type args struct {
		jobs []*Job
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "pass",
			args: args{
				jobs: []*Job{
					{
						Name: "test",
						Trigger: DateTrigger{
							NextRunTime: time.Now(),
						},
						Task: TestTask{},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ags := &AGScheduler{
				Store:      tt.fields.Store,
				Logger:     tt.fields.Logger,
				Status:     tt.fields.Status,
				Context:    tt.fields.Context,
				WaitCancel: tt.fields.WaitCancel,
			}
			if err := ags.AddJob(tt.args.jobs...); (err != nil) != tt.wantErr {
				t.Errorf("AddJob() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAGScheduler_Close(t *testing.T) {
	type fields struct {
		Store      IStore
		Logger     *logrus.Entry
		Status     STATUS
		Context    context.Context
		WaitCancel context.CancelFunc
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "pass",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ags := &AGScheduler{
				Store:      tt.fields.Store,
				Logger:     tt.fields.Logger,
				Status:     tt.fields.Status,
				Context:    tt.fields.Context,
				WaitCancel: tt.fields.WaitCancel,
			}
			if err := ags.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAGScheduler_DelJob(t *testing.T) {
	job := Job{
		Name: "test",
		Trigger: DateTrigger{
			NextRunTime: time.Now(),
		},
		Task: TestTask{},
	}
	type fields struct {
		Store      IStore
		Logger     *logrus.Entry
		Status     STATUS
		Context    context.Context
		WaitCancel context.CancelFunc
	}
	type args struct {
		jobs []*Job
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "pass",
			args: args{
				jobs: []*Job{&job},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ags := &AGScheduler{
				Store:      tt.fields.Store,
				Logger:     tt.fields.Logger,
				Status:     tt.fields.Status,
				Context:    tt.fields.Context,
				WaitCancel: tt.fields.WaitCancel,
			}
			_ = ags.AddJob(&job)
			if err := ags.DelJob(tt.args.jobs...); (err != nil) != tt.wantErr {
				t.Errorf("DelJob() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAGScheduler_FillByDefault(t *testing.T) {
	type fields struct {
		Store      IStore
		Logger     *logrus.Entry
		Status     STATUS
		Context    context.Context
		WaitCancel context.CancelFunc
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "pass",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ags := &AGScheduler{
				Store:      tt.fields.Store,
				Logger:     tt.fields.Logger,
				Status:     tt.fields.Status,
				Context:    tt.fields.Context,
				WaitCancel: tt.fields.WaitCancel,
			}
			ags.FillByDefault()
		})
	}
}

func TestAGScheduler_GetAllJobs(t *testing.T) {
	job := Job{
		Name: "test",
		Trigger: DateTrigger{
			NextRunTime: time.Now(),
		},
		Task: TestTask{},
	}
	type fields struct {
		Store      IStore
		Logger     *logrus.Entry
		Status     STATUS
		Context    context.Context
		WaitCancel context.CancelFunc
	}
	tests := []struct {
		name     string
		fields   fields
		wantJobs []*Job
		wantErr  bool
	}{
		{
			name: "pass",
			wantJobs: []*Job{
				&job,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ags := &AGScheduler{
				Store:      tt.fields.Store,
				Logger:     tt.fields.Logger,
				Status:     tt.fields.Status,
				Context:    tt.fields.Context,
				WaitCancel: tt.fields.WaitCancel,
			}
			_ = ags.AddJob(&job)
			gotJobs, err := ags.GetAllJobs()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllJobs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotJobs, tt.wantJobs) {
				t.Errorf("GetAllJobs() gotJobs = %v, want %v", gotJobs, tt.wantJobs)
			}
		})
	}
}

func TestAGScheduler_GetJobByJobName(t *testing.T) {
	job := Job{
		Name: "test",
		Trigger: DateTrigger{
			NextRunTime: time.Now(),
		},
		Task: TestTask{},
	}
	type fields struct {
		Store      IStore
		Logger     *logrus.Entry
		Status     STATUS
		Context    context.Context
		WaitCancel context.CancelFunc
	}
	type args struct {
		jobName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantJob *Job
		wantErr bool
	}{
		{
			name: "pass",
			args: args{
				jobName: job.Name,
			},
			wantJob: &job,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ags := &AGScheduler{
				Store:      tt.fields.Store,
				Logger:     tt.fields.Logger,
				Status:     tt.fields.Status,
				Context:    tt.fields.Context,
				WaitCancel: tt.fields.WaitCancel,
			}
			_ = ags.AddJob(&job)
			gotJob, err := ags.GetJobByJobName(tt.args.jobName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetJobByJobName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotJob, tt.wantJob) {
				t.Errorf("GetJobByJobName() gotJob = %v, want %v", gotJob, tt.wantJob)
			}
		})
	}
}

func TestAGScheduler_Pause(t *testing.T) {
	type fields struct {
		Store      IStore
		Logger     *logrus.Entry
		Status     STATUS
		Context    context.Context
		WaitCancel context.CancelFunc
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "pass",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ags := &AGScheduler{
				Store:      tt.fields.Store,
				Logger:     tt.fields.Logger,
				Status:     tt.fields.Status,
				Context:    tt.fields.Context,
				WaitCancel: tt.fields.WaitCancel,
			}
			ags.Pause()
		})
	}
}

func TestAGScheduler_Start(t *testing.T) {
	type fields struct {
		Store      IStore
		Logger     *logrus.Entry
		Status     STATUS
		Context    context.Context
		WaitCancel context.CancelFunc
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "pass",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ags := &AGScheduler{
				Store:      tt.fields.Store,
				Logger:     tt.fields.Logger,
				Status:     tt.fields.Status,
				Context:    tt.fields.Context,
				WaitCancel: tt.fields.WaitCancel,
			}
			go func() {
				time.Sleep(time.Second)
				_ = ags.Close()
			}()
			ags.Start()
		})
	}
}

func TestAGScheduler_UpdateJob(t *testing.T) {
	type fields struct {
		Store      IStore
		Logger     *logrus.Entry
		Status     STATUS
		Context    context.Context
		WaitCancel context.CancelFunc
	}
	type args struct {
		jobs []*Job
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ags := &AGScheduler{
				Store:      tt.fields.Store,
				Logger:     tt.fields.Logger,
				Status:     tt.fields.Status,
				Context:    tt.fields.Context,
				WaitCancel: tt.fields.WaitCancel,
			}
			if err := ags.UpdateJob(tt.args.jobs...); (err != nil) != tt.wantErr {
				t.Errorf("UpdateJob() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAGScheduler_WaitWithTime(t *testing.T) {
	type fields struct {
		Store      IStore
		Logger     *logrus.Entry
		Status     STATUS
		Context    context.Context
		WaitCancel context.CancelFunc
	}
	type args struct {
		waitTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "pass",
			args: args{
				waitTime: MaxDateTime,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ags := &AGScheduler{
				Store:      tt.fields.Store,
				Logger:     tt.fields.Logger,
				Status:     tt.fields.Status,
				Context:    tt.fields.Context,
				WaitCancel: tt.fields.WaitCancel,
			}
			go func() {
				time.Sleep(time.Second)
				ags.Wake()
			}()
			ags.WaitWithTime(MaxDateTime)
		})
	}
}

func TestAGScheduler_Wake(t *testing.T) {
	type fields struct {
		Store      IStore
		Logger     *logrus.Entry
		Status     STATUS
		Context    context.Context
		WaitCancel context.CancelFunc
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "pass",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ags := &AGScheduler{
				Store:      tt.fields.Store,
				Logger:     tt.fields.Logger,
				Status:     tt.fields.Status,
				Context:    tt.fields.Context,
				WaitCancel: tt.fields.WaitCancel,
			}
			go func() {
				time.Sleep(time.Second)
				ags.Wake()
			}()
			ags.WaitWithTime(MaxDateTime)
		})
	}
}

func TestAGScheduler_listenSignal(t *testing.T) {
	type fields struct {
		Store      IStore
		Logger     *logrus.Entry
		Status     STATUS
		Context    context.Context
		WaitCancel context.CancelFunc
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "pass",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ags := &AGScheduler{
				Store:      tt.fields.Store,
				Logger:     tt.fields.Logger,
				Status:     tt.fields.Status,
				Context:    tt.fields.Context,
				WaitCancel: tt.fields.WaitCancel,
			}
			go ags.listenSignal()
		})
	}
}
