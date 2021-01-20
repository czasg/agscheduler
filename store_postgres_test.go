package agscheduler

import (
	"github.com/go-pg/pg/v10"
	"github.com/sirupsen/logrus"
	"reflect"
	"testing"
	"time"
)

func TestPostgresStore_AddJob(t *testing.T) {
	type fields struct {
		Logger *logrus.Entry
		PG     *pg.DB
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &PostgresStore{
				Logger: tt.fields.Logger,
				PG:     tt.fields.PG,
			}
			if err := ps.AddJob(tt.args.job); (err != nil) != tt.wantErr {
				t.Errorf("AddJob() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPostgresStore_DelJob(t *testing.T) {
	type fields struct {
		Logger *logrus.Entry
		PG     *pg.DB
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &PostgresStore{
				Logger: tt.fields.Logger,
				PG:     tt.fields.PG,
			}
			if err := ps.DelJob(tt.args.job); (err != nil) != tt.wantErr {
				t.Errorf("DelJob() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPostgresStore_FillByDefault(t *testing.T) {
	type fields struct {
		Logger *logrus.Entry
		PG     *pg.DB
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "empty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &PostgresStore{
				Logger: tt.fields.Logger,
				PG:     tt.fields.PG,
			}
			ps.FillByDefault()
			if ps.Logger == nil || ps.PG == nil {
				t.Errorf("FillByDefault() error")
			}
		})
	}
}

func TestPostgresStore_GetAllJobs(t *testing.T) {
	type fields struct {
		Logger *logrus.Entry
		PG     *pg.DB
	}
	tests := []struct {
		name    string
		fields  fields
		want    []*Job
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &PostgresStore{
				Logger: tt.fields.Logger,
				PG:     tt.fields.PG,
			}
			got, err := ps.GetAllJobs()
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

func TestPostgresStore_GetJobByName(t *testing.T) {
	type fields struct {
		Logger *logrus.Entry
		PG     *pg.DB
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &PostgresStore{
				Logger: tt.fields.Logger,
				PG:     tt.fields.PG,
			}
			got, err := ps.GetJobByName(tt.args.name)
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

func TestPostgresStore_GetNextRunTime(t *testing.T) {
	type fields struct {
		Logger *logrus.Entry
		PG     *pg.DB
	}
	tests := []struct {
		name    string
		fields  fields
		want    time.Time
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &PostgresStore{
				Logger: tt.fields.Logger,
				PG:     tt.fields.PG,
			}
			got, err := ps.GetNextRunTime()
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

func TestPostgresStore_GetSchedulingJobs(t *testing.T) {
	type fields struct {
		Logger *logrus.Entry
		PG     *pg.DB
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &PostgresStore{
				Logger: tt.fields.Logger,
				PG:     tt.fields.PG,
			}
			got, err := ps.GetSchedulingJobs(tt.args.now)
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

func TestPostgresStore_UpdateJob(t *testing.T) {
	type fields struct {
		Logger *logrus.Entry
		PG     *pg.DB
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &PostgresStore{
				Logger: tt.fields.Logger,
				PG:     tt.fields.PG,
			}
			if err := ps.UpdateJob(tt.args.job); (err != nil) != tt.wantErr {
				t.Errorf("UpdateJob() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
