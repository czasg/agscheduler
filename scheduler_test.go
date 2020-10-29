package AGScheduler

import (
	"context"
	"github.com/sirupsen/logrus"
	"reflect"
	"testing"
	"time"
)

func TestNewScheduler(t *testing.T) {
	type args struct {
		store IStore
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "new",
			args: args{
				store: NewMemoryStore(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scheduler := NewScheduler(tt.args.store)
			scheduler.Wake()
		})
	}
}

func TestScheduler_AddTask(t *testing.T) {
	now := time.Now()
	interval, _ := NewIntervalTrigger(now, EmptyDateTime, time.Second)
	task := NewTask("task", interval, func(args ...interface{}) {})

	type fields struct {
		StoresMap   map[string]IStore
		Logger      *logrus.Entry
		Controller  *Controller
		CloseCancel context.CancelFunc
	}
	type args struct {
		task *Task
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "add",
			fields: fields{
				StoresMap: map[string]IStore{
					"default": NewMemoryStore(),
				},
				Logger: logrus.New().WithFields(logrus.Fields{
					"Module": "GoTest",
				}),
				Controller: NewController(),
			},
			args: args{
				task: task,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scheduler{
				StoresMap:   tt.fields.StoresMap,
				Logger:      tt.fields.Logger,
				Controller:  tt.fields.Controller,
				CloseCancel: tt.fields.CloseCancel,
			}
			if err := s.AddTask(tt.args.task); (err != nil) != tt.wantErr {
				t.Errorf("AddTask() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestScheduler_AddTaskFromTasksMap(t *testing.T) {
	err := RegisterWorksMap(map[string]WorkDetail{
		"test": WorkDetail{
			Func: func(args ...interface{}) {},
			Args: []interface{}{},
		},
	})
	if err != nil {
		panic(err)
	}
	now := time.Now()
	interval, _ := NewIntervalTrigger(now, EmptyDateTime, time.Second)

	type fields struct {
		StoresMap   map[string]IStore
		Logger      *logrus.Entry
		Controller  *Controller
		CloseCancel context.CancelFunc
	}
	type args struct {
		name       string
		taskMapKey string
		args       []interface{}
		trigger    ITrigger
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "add-from-map",
			fields: fields{
				StoresMap: map[string]IStore{
					"default": NewMemoryStore(),
				},
				Logger: logrus.New().WithFields(logrus.Fields{
					"Module": "GoTest",
				}),
				Controller: NewController(),
			},
			args: args{
				name:       "task",
				taskMapKey: "test",
				args:       []interface{}{},
				trigger:    interval,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scheduler{
				StoresMap:   tt.fields.StoresMap,
				Logger:      tt.fields.Logger,
				Controller:  tt.fields.Controller,
				CloseCancel: tt.fields.CloseCancel,
			}
			if err := s.AddTaskFromTasksMap(tt.args.name, tt.args.taskMapKey, tt.args.trigger, tt.args.args...); (err != nil) != tt.wantErr {
				t.Errorf("AddTaskFromTasksMap() error = %v, wantErr %v", err, tt.wantErr)
			}
			task, err := s.GetTaskByName("task")
			if err != nil {
				t.Error("add fail")
				return
			}
			if task.Name != "task" {
				t.Error("add fail")
			}
		})
	}
}

func TestScheduler_Close(t *testing.T) {
	type fields struct {
		StoresMap   map[string]IStore
		Logger      *logrus.Entry
		Controller  *Controller
		CloseCancel context.CancelFunc
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "close",
			fields: fields{
				StoresMap: map[string]IStore{
					"default": NewMemoryStore(),
				},
				Logger:     logrus.New().WithFields(logrus.Fields{}),
				Controller: NewController(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scheduler{
				StoresMap:   tt.fields.StoresMap,
				Logger:      tt.fields.Logger,
				Controller:  tt.fields.Controller,
				CloseCancel: tt.fields.CloseCancel,
			}
			go s.Close()
			s.Start()
		})
	}
}

func TestScheduler_DelTask(t *testing.T) {
	now := time.Now()
	interval, _ := NewIntervalTrigger(now, EmptyDateTime, time.Second)
	task := NewTask("task", interval, func(args ...interface{}) {})

	memory1 := NewMemoryStore()
	memory2 := NewMemoryStore()
	err := memory2.AddTask(task)
	if err != nil {
		panic(err)
	}

	type fields struct {
		StoresMap   map[string]IStore
		Logger      *logrus.Entry
		Controller  *Controller
		CloseCancel context.CancelFunc
	}
	type args struct {
		task *Task
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "del-fail",
			fields: fields{
				StoresMap: map[string]IStore{
					"default": memory1,
				},
				Logger: logrus.New().WithFields(logrus.Fields{
					"Module": "GoTest",
				}),
				Controller: NewController(),
			},
			args: args{
				task: task,
			},
			wantErr: true,
		},
		{
			name: "del-success",
			fields: fields{
				StoresMap: map[string]IStore{
					"default": memory2,
				},
				Logger: logrus.New().WithFields(logrus.Fields{
					"Module": "GoTest",
				}),
				Controller: NewController(),
			},
			args: args{
				task: task,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scheduler{
				StoresMap:   tt.fields.StoresMap,
				Logger:      tt.fields.Logger,
				Controller:  tt.fields.Controller,
				CloseCancel: tt.fields.CloseCancel,
			}
			if err := s.DelTask(tt.args.task); (err != nil) != tt.wantErr {
				t.Errorf("DelTask() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestScheduler_GetAllTasks(t *testing.T) {
	now := time.Now()
	interval, _ := NewIntervalTrigger(now, EmptyDateTime, time.Second)
	task1 := NewTask("task1", interval, func(args ...interface{}) {})
	task2 := NewTask("task2", interval, func(args ...interface{}) {})

	memory := NewMemoryStore()
	_ = memory.AddTask(task1)
	_ = memory.AddTask(task2)

	type fields struct {
		StoresMap   map[string]IStore
		Logger      *logrus.Entry
		Controller  *Controller
		CloseCancel context.CancelFunc
	}
	tests := []struct {
		name   string
		fields fields
		want   []*Task
	}{
		{
			name: "del-fail",
			fields: fields{
				StoresMap: map[string]IStore{
					"default": memory,
				},
				Logger: logrus.New().WithFields(logrus.Fields{
					"Module": "GoTest",
				}),
				Controller: NewController(),
			},
			want: []*Task{task2, task1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scheduler{
				StoresMap:   tt.fields.StoresMap,
				Logger:      tt.fields.Logger,
				Controller:  tt.fields.Controller,
				CloseCancel: tt.fields.CloseCancel,
			}
			if got := s.GetAllTasks(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllTasks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScheduler_GetTaskByName(t *testing.T) {
	now := time.Now()
	interval, _ := NewIntervalTrigger(now, EmptyDateTime, time.Second)
	task := NewTask("task", interval, func(args ...interface{}) {})
	memory := NewMemoryStore()
	_ = memory.AddTask(task)

	type fields struct {
		StoresMap   map[string]IStore
		Logger      *logrus.Entry
		Controller  *Controller
		CloseCancel context.CancelFunc
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Task
		wantErr bool
	}{
		{
			name: "succ",
			fields: fields{
				StoresMap: map[string]IStore{
					"default": memory,
				},
				Logger: logrus.New().WithFields(logrus.Fields{
					"Module": "GoTest",
				}),
				Controller: NewController(),
			},
			args: args{
				name: "task",
			},
			want:    task,
			wantErr: false,
		},
		{
			name: "fail",
			fields: fields{
				StoresMap: map[string]IStore{
					"default": NewMemoryStore(),
				},
				Logger: logrus.New().WithFields(logrus.Fields{
					"Module": "GoTest",
				}),
				Controller: NewController(),
			},
			args: args{
				name: "task",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scheduler{
				StoresMap:   tt.fields.StoresMap,
				Logger:      tt.fields.Logger,
				Controller:  tt.fields.Controller,
				CloseCancel: tt.fields.CloseCancel,
			}
			got, err := s.GetTaskByName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTaskByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTaskByName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScheduler_Start(t *testing.T) {
	type fields struct {
		StoresMap   map[string]IStore
		Logger      *logrus.Entry
		Controller  *Controller
		CloseCancel context.CancelFunc
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "start",
			fields: fields{
				StoresMap: map[string]IStore{
					"default": NewMemoryStore(),
				},
				Logger:     logrus.New().WithFields(logrus.Fields{}),
				Controller: NewController(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scheduler{
				StoresMap:   tt.fields.StoresMap,
				Logger:      tt.fields.Logger,
				Controller:  tt.fields.Controller,
				CloseCancel: tt.fields.CloseCancel,
			}
			go s.Close()
			s.Start()
		})
	}
}

func TestScheduler_UpdateTask(t *testing.T) {
	now := time.Now()
	interval, _ := NewIntervalTrigger(now, EmptyDateTime, time.Second)
	task := NewTask("task", interval, func(args ...interface{}) {})
	memory := NewMemoryStore()
	_ = memory.AddTask(task)

	type fields struct {
		StoresMap   map[string]IStore
		Logger      *logrus.Entry
		Controller  *Controller
		CloseCancel context.CancelFunc
	}
	type args struct {
		task *Task
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "update",
			fields: fields{
				StoresMap: map[string]IStore{
					"default": memory,
				},
				Logger: logrus.New().WithFields(logrus.Fields{
					"Module": "GoTest",
				}),
				Controller: NewController(),
			},
			args: args{
				task: task,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scheduler{
				StoresMap:   tt.fields.StoresMap,
				Logger:      tt.fields.Logger,
				Controller:  tt.fields.Controller,
				CloseCancel: tt.fields.CloseCancel,
			}
			if err := s.UpdateTask(tt.args.task); (err != nil) != tt.wantErr {
				t.Errorf("UpdateTask() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestScheduler_Wake(t *testing.T) {
	type fields struct {
		StoresMap   map[string]IStore
		Logger      *logrus.Entry
		Controller  *Controller
		CloseCancel context.CancelFunc
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "wake",
			fields: fields{
				StoresMap: map[string]IStore{
					"default": NewMemoryStore(),
				},
				Logger: logrus.New().WithFields(logrus.Fields{
					"Module": "GoTest",
				}),
				Controller: NewController(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scheduler{
				StoresMap:   tt.fields.StoresMap,
				Logger:      tt.fields.Logger,
				Controller:  tt.fields.Controller,
				CloseCancel: tt.fields.CloseCancel,
			}
			s.Wake()
		})
	}
}
