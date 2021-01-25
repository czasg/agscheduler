package agscheduler

import (
	"context"
	"fmt"
	"reflect"
	"testing"
)

type TestTask struct {
	Name string
	Age  float64
}

func (t TestTask) Run(ctx context.Context) {
	fmt.Println(t.Name, t.Age)
}

func TestDeserializeTask(t *testing.T) {
	job := Job{
		Name: "test",
		Task: &TestTask{
			Name: "test-task",
			Age:  18,
		},
	}
	RegisterAllTasks(&TestTask{})
	err := SerializeTask(&job)
	if err != nil {
		panic(err)
	}
	job.Task = nil
	type args struct {
		job *Job
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "pass",
			args: args{
				job: &job,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeserializeTask(tt.args.job); (err != nil) != tt.wantErr {
				t.Errorf("DeserializeTask() error = %v, wantErr %v", err, tt.wantErr)
			}
			tt.args.job.Task.Run(context.Background())
		})
	}
}

func TestRegisterAllTasks(t *testing.T) {
	type args struct {
		tasks []ITask
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "pass",
			args: args{
				[]ITask{&TestTask{}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RegisterAllTasks(tt.args.tasks...)
			if len(allITasks) == 0 {
				t.Errorf("tasks should not be null.")
			}
		})
	}
}

func TestSerializeTask(t *testing.T) {
	RegisterAllTasks(&TestTask{})
	type args struct {
		job *Job
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "pass",
			args: args{
				job: &Job{
					Name: "test",
					Task: &TestTask{
						Name: "test-task",
						Age:  18,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SerializeTask(tt.args.job); (err != nil) != tt.wantErr {
				t.Errorf("SerializeTask() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.args.job.TaskMeta == nil {
				t.Errorf("task meta should not be nil.")
			}
			fmt.Println(tt.args.job.TaskMeta)
		})
	}
}

func Test_addTask(t *testing.T) {
	type args struct {
		taskName string
		taskType reflect.Type
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "pass",
			args: args{
				taskName: "test",
				taskType: reflect.TypeOf("test"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addTask(tt.args.taskName, tt.args.taskType)
			if allITasks["test"] == nil {
				t.Errorf("tasks should not't be nil.")
			}
		})
	}
}
