package agscheduler

import (
	"container/list"
	"github.com/sirupsen/logrus"
	"reflect"
	"testing"
	"time"
)

func TestMemoryStore_AddJob(t *testing.T) {
	job := &Job{
		Name: "test",
	}
	testList := list.New()
	testList.PushBack(job)
	errList := list.New()
	errList.PushBack("test")
	type fields struct {
		Jobs    *list.List
		JobsMap map[string]*list.Element
		Logger  *logrus.Entry
	}
	type args struct {
		job *Job
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "empty add",
			args: args{
				job: job,
			},
			wantErr: false,
		},
		{
			name: "pass",
			fields: fields{
				Jobs: testList,
			},
			args: args{
				job: job,
			},
			wantErr: false,
		},
		{
			name: "not job",
			fields: fields{
				Jobs: errList,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MemoryStore{
				Jobs:    tt.fields.Jobs,
				JobsMap: tt.fields.JobsMap,
				Logger:  tt.fields.Logger,
			}
			if err := ms.AddJob(tt.args.job); (err != nil) != tt.wantErr {
				t.Errorf("AddJob() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemoryStore_DelJob(t *testing.T) {
	job := &Job{Name: "test"}
	testList := list.New()
	type fields struct {
		Jobs    *list.List
		JobsMap map[string]*list.Element
		Logger  *logrus.Entry
	}
	type args struct {
		job *Job
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "not found",
			args: args{
				job: &Job{Name: "test"},
			},
			wantErr: true,
		},
		{
			name: "pass",
			fields: fields{
				Jobs: testList,
				JobsMap: map[string]*list.Element{
					"test": testList.PushFront(job),
				},
			},
			args: args{
				job: &Job{Name: "test"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MemoryStore{
				Jobs:    tt.fields.Jobs,
				JobsMap: tt.fields.JobsMap,
				Logger:  tt.fields.Logger,
			}
			if err := ms.DelJob(tt.args.job); (err != nil) != tt.wantErr {
				t.Errorf("DelJob() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemoryStore_FillByDefault(t *testing.T) {
	type fields struct {
		Jobs    *list.List
		JobsMap map[string]*list.Element
		Logger  *logrus.Entry
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "pass",
			fields: fields{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MemoryStore{
				Jobs:    tt.fields.Jobs,
				JobsMap: tt.fields.JobsMap,
				Logger:  tt.fields.Logger,
			}
			ms.FillByDefault()
			if ms.Jobs == nil || ms.JobsMap == nil || ms.Logger == nil {
				t.Error("FillByDefault(), should not be nil")
			}
		})
	}
}

func TestMemoryStore_GetAllJobs(t *testing.T) {
	Job1 := &Job{Name: "test1"}
	Job2 := &Job{Name: "test2"}
	testList := list.New()
	testList.PushBack(Job1)
	testList.PushBack(Job2)
	testListB := list.New()
	testListB.PushBack("test")
	type fields struct {
		Jobs    *list.List
		JobsMap map[string]*list.Element
		Logger  *logrus.Entry
	}
	tests := []struct {
		name    string
		fields  fields
		want    []*Job
		wantErr bool
	}{
		{
			name:    "empty",
			want:    []*Job{},
			wantErr: false,
		},
		{
			name: "not job",
			fields: fields{
				Jobs: testListB,
			},
			want:    []*Job{},
			wantErr: true,
		},
		{
			name: "pass",
			fields: fields{
				Jobs: testList,
			},
			want:    []*Job{Job1, Job2},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MemoryStore{
				Jobs:    tt.fields.Jobs,
				JobsMap: tt.fields.JobsMap,
				Logger:  tt.fields.Logger,
			}
			got, err := ms.GetAllJobs()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllJobs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllJobs() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryStore_GetJobByName(t *testing.T) {
	job := &Job{
		Name: "test",
	}
	type fields struct {
		Jobs    *list.List
		JobsMap map[string]*list.Element
		Logger  *logrus.Entry
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Job
		wantErr bool
	}{
		{
			name:    "not found",
			fields:  fields{},
			args:    args{},
			want:    nil,
			wantErr: true,
		},
		{
			name: "not job",
			fields: fields{
				JobsMap: map[string]*list.Element{
					"test": list.New().PushBack("test"),
				},
			},
			args: args{
				name: "test",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "not job",
			fields: fields{
				JobsMap: map[string]*list.Element{
					"test": list.New().PushBack("test"),
				},
			},
			args: args{
				name: "test",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "pass",
			fields: fields{
				JobsMap: map[string]*list.Element{
					"test": list.New().PushBack(job),
				},
			},
			args: args{
				name: "test",
			},
			want:    job,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MemoryStore{
				Jobs:    tt.fields.Jobs,
				JobsMap: tt.fields.JobsMap,
				Logger:  tt.fields.Logger,
			}
			got, err := ms.GetJobByName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetJobByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetJobByName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryStore_GetNextRunTime(t *testing.T) {
	job := &Job{Name: "test"}
	testList := list.New()
	testList.PushBack(job)
	type fields struct {
		Jobs    *list.List
		JobsMap map[string]*list.Element
		Logger  *logrus.Entry
	}
	tests := []struct {
		name    string
		fields  fields
		want    time.Time
		wantErr bool
	}{
		{
			name:    "empty",
			want:    MaxDateTime,
			wantErr: false,
		},
		{
			name: "pass",
			fields: fields{
				Jobs: testList,
			},
			want:    MinDateTime,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MemoryStore{
				Jobs:    tt.fields.Jobs,
				JobsMap: tt.fields.JobsMap,
				Logger:  tt.fields.Logger,
			}
			got, err := ms.GetNextRunTime()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNextRunTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNextRunTime() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryStore_GetSchedulingJobs(t *testing.T) {
	type fields struct {
		Jobs    *list.List
		JobsMap map[string]*list.Element
		Logger  *logrus.Entry
	}
	type args struct {
		now time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*Job
		wantErr bool
	}{
		{
			name:    "empty",
			fields:  fields{},
			args:    args{},
			want:    []*Job{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MemoryStore{
				Jobs:    tt.fields.Jobs,
				JobsMap: tt.fields.JobsMap,
				Logger:  tt.fields.Logger,
			}
			got, err := ms.GetSchedulingJobs(tt.args.now)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSchedulingJobs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSchedulingJobs() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryStore_UpdateJob(t *testing.T) {
	job := &Job{Name: "test"}
	testList := list.New()
	type fields struct {
		Jobs    *list.List
		JobsMap map[string]*list.Element
		Logger  *logrus.Entry
	}
	type args struct {
		job *Job
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "not found",
			fields: fields{},
			args: args{
				job: job,
			},
			wantErr: true,
		},
		{
			name: "pass",
			fields: fields{
				Jobs: testList,
				JobsMap: map[string]*list.Element{
					"test": testList.PushBack(job),
				},
			},
			args: args{
				job: job,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MemoryStore{
				Jobs:    tt.fields.Jobs,
				JobsMap: tt.fields.JobsMap,
				Logger:  tt.fields.Logger,
			}
			if err := ms.UpdateJob(tt.args.job); (err != nil) != tt.wantErr {
				t.Errorf("UpdateJob() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
