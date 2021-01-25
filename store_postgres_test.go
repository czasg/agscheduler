package agscheduler

//import (
//	"context"
//	"github.com/go-pg/pg/v10"
//	"github.com/sirupsen/logrus"
//	"reflect"
//	"testing"
//	"time"
//)
//
//func TestPostgresStore_AddJob(t *testing.T) {
//	type fields struct {
//		Logger *logrus.Entry
//		PG     *pg.DB
//	}
//	type args struct {
//		job *Job
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		{
//			name: "pass",
//			args: args{
//				job: &Job{
//					Name: "test1",
//					Trigger: DateTrigger{
//						NextRunTime: MaxDateTime.Add(-time.Hour),
//					},
//					Task: &TestTask{
//						Name: "test-task",
//						Age:  1,
//					},
//				},
//			},
//			wantErr: false,
//		},
//		{
//			name: "pass",
//			args: args{
//				job: &Job{
//					Name: "test2",
//					Trigger: DateTrigger{
//						NextRunTime: MaxDateTime.Add(-time.Minute),
//					},
//					Task: &TestTask{
//						Name: "test-task",
//						Age:  11,
//					},
//				},
//			},
//			wantErr: false,
//		},
//		{
//			name: "pass",
//			args: args{
//				job: &Job{
//					Name: "test3",
//					Trigger: DateTrigger{
//						NextRunTime: MaxDateTime,
//					},
//					Task: &TestTask{
//						Name: "test-task",
//						Age:  111,
//					},
//				},
//			},
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			ps := &PostgresStore{
//				Logger: tt.fields.Logger,
//				PG:     tt.fields.PG,
//			}
//			if err := ps.AddJob(tt.args.job); (err != nil) != tt.wantErr {
//				t.Errorf("AddJob() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestPostgresStore_DelJob(t *testing.T) {
//	type fields struct {
//		Logger *logrus.Entry
//		PG     *pg.DB
//	}
//	type args struct {
//		job *Job
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		{
//			name: "pass",
//			args: args{
//				job: &Job{
//					Name: "test1",
//				},
//			},
//			wantErr: false,
//		},
//		{
//			name: "pass",
//			args: args{
//				job: &Job{
//					Name: "test2",
//				},
//			},
//			wantErr: false,
//		},
//		{
//			name: "pass",
//			args: args{
//				job: &Job{
//					Name: "test3",
//				},
//			},
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			ps := &PostgresStore{
//				Logger: tt.fields.Logger,
//				PG:     tt.fields.PG,
//			}
//			if err := ps.DelJob(tt.args.job); (err != nil) != tt.wantErr {
//				t.Errorf("DelJob() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestPostgresStore_FillByDefault(t *testing.T) {
//	type fields struct {
//		Logger *logrus.Entry
//		PG     *pg.DB
//	}
//	tests := []struct {
//		name   string
//		fields fields
//	}{
//		{
//			name: "empty",
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			ps := &PostgresStore{
//				Logger: tt.fields.Logger,
//				PG:     tt.fields.PG,
//			}
//			ps.FillByDefault()
//			if ps.Logger == nil || ps.PG == nil {
//				t.Errorf("FillByDefault() error")
//			}
//		})
//	}
//}
//
//func TestPostgresStore_GetAllJobs(t *testing.T) {
//	RegisterAllTasks(&TestTask{})
//	type fields struct {
//		Logger *logrus.Entry
//		PG     *pg.DB
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		want    []*Job
//		wantErr bool
//	}{
//		{
//			name:    "pass",
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			ps := &PostgresStore{
//				Logger: tt.fields.Logger,
//				PG:     tt.fields.PG,
//			}
//			_, err := ps.GetAllJobs()
//			if (err != nil) != tt.wantErr {
//				t.Errorf("GetAllJobs() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//		})
//	}
//}
//
//func TestPostgresStore_GetJobByName(t *testing.T) {
//	RegisterAllTasks(&TestTask{})
//	type fields struct {
//		Logger *logrus.Entry
//		PG     *pg.DB
//	}
//	type args struct {
//		name string
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    *Job
//		wantErr bool
//	}{
//		{
//			name: "pass",
//			args: args{
//				name: "test1",
//			},
//			want: &Job{
//				Name: "test1",
//				Trigger: DateTrigger{
//					NextRunTime: MaxDateTime.Add(-time.Hour),
//				},
//				Task: TestTask{
//					Name: "test-task",
//					Age:  1,
//				},
//			},
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			ps := &PostgresStore{
//				Logger: tt.fields.Logger,
//				PG:     tt.fields.PG,
//			}
//			got, err := ps.GetJobByName(tt.args.name)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("GetJobByName() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if got == nil {
//				return
//			}
//			got.Task.Run(context.Background())
//		})
//	}
//}
//
//func TestPostgresStore_GetNextRunTime(t *testing.T) {
//	type fields struct {
//		Logger *logrus.Entry
//		PG     *pg.DB
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		want    time.Time
//		wantErr bool
//	}{
//		{
//			name:    "pass",
//			want:    MaxDateTime.Add(-time.Hour),
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			ps := &PostgresStore{
//				Logger: tt.fields.Logger,
//				PG:     tt.fields.PG,
//			}
//			got, err := ps.GetNextRunTime()
//			if (err != nil) != tt.wantErr {
//				t.Errorf("GetNextRunTime() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("GetNextRunTime() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestPostgresStore_GetSchedulingJobs(t *testing.T) {
//	RegisterAllTasks(&TestTask{})
//	type fields struct {
//		Logger *logrus.Entry
//		PG     *pg.DB
//	}
//	type args struct {
//		now time.Time
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    []*Job
//		wantErr bool
//	}{
//		{
//			name: "pass",
//			args: args{
//				now: time.Now(),
//			},
//			want:    []*Job{},
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			ps := &PostgresStore{
//				Logger: tt.fields.Logger,
//				PG:     tt.fields.PG,
//			}
//			got, err := ps.GetSchedulingJobs(tt.args.now)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("GetSchedulingJobs() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("GetSchedulingJobs() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestPostgresStore_UpdateJob(t *testing.T) {
//	type fields struct {
//		Logger *logrus.Entry
//		PG     *pg.DB
//	}
//	type args struct {
//		job *Job
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		{
//			name: "pass",
//			args: args{
//				job: &Job{
//					Name: "test1",
//					Trigger: DateTrigger{
//						NextRunTime: MaxDateTime,
//					},
//					Task: &TestTask{
//						Name: "test-task",
//						Age:  666,
//					},
//				},
//			},
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			ps := &PostgresStore{
//				Logger: tt.fields.Logger,
//				PG:     tt.fields.PG,
//			}
//			if err := ps.UpdateJob(tt.args.job); (err != nil) != tt.wantErr {
//				t.Errorf("UpdateJob() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
