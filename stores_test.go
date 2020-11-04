package AGScheduler

import (
	"container/list"
	"reflect"
	"testing"
	"time"
)

func TestMemoryStore_AddTask(t *testing.T) {
	now := time.Now()
	interval, _ := NewIntervalTrigger(now, EmptyDateTime, time.Second)

	type fields struct {
		Tasks    *list.List
		TasksMap map[string]*list.Element
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
				Tasks:    list.New(),
				TasksMap: map[string]*list.Element{},
			},
			args: args{
				task: NewTask(
					"task",
					interval,
					func(args ...interface{}) {},
				),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemoryStore{
				Tasks:    tt.fields.Tasks,
				TasksMap: tt.fields.TasksMap,
			}
			if err := m.AddTask(tt.args.task); (err != nil) != tt.wantErr {
				t.Errorf("AddTask() error = %v, wantErr %v", err, tt.wantErr)
			}
			if m.Tasks.Len() != 1 {
				t.Errorf("add fail")
			}
		})
	}
}

func TestMemoryStore_DelTask(t *testing.T) {
	now := time.Now()
	interval, _ := NewIntervalTrigger(now, EmptyDateTime, time.Second)

	task := NewTask(
		"task",
		interval,
		func(args ...interface{}) {},
	)

	taskList := list.New()
	ele := taskList.PushFront(task)

	type fields struct {
		Tasks    *list.List
		TasksMap map[string]*list.Element
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
			name: "del",
			fields: fields{
				Tasks: taskList,
				TasksMap: map[string]*list.Element{
					"task": ele,
				},
			},
			args: args{
				task: task,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemoryStore{
				Tasks:    tt.fields.Tasks,
				TasksMap: tt.fields.TasksMap,
			}
			if err := m.DelTask(tt.args.task); (err != nil) != tt.wantErr {
				t.Errorf("DelTask() error = %v, wantErr %v", err, tt.wantErr)
			}
			if m.Tasks.Len() != 0 {
				t.Error("del fail")
			}
			if len(m.TasksMap) != 0 {
				t.Error("del fail")
			}
		})
	}
}

func TestMemoryStore_GetAllTasks(t *testing.T) {
	cron, _ := NewCronTrigger("* * * * *")

	task1 := NewTask("task1", cron, func(args ...interface{}) {})
	task2 := NewTask("task2", cron, func(args ...interface{}) {})

	taskList := list.New()
	ele1 := taskList.PushBack(task1)
	ele2 := taskList.PushBack(task2)

	type fields struct {
		Tasks    *list.List
		TasksMap map[string]*list.Element
	}
	tests := []struct {
		name   string
		fields fields
		want   []*Task
	}{
		{
			name: "get all",
			fields: fields{
				Tasks: taskList,
				TasksMap: map[string]*list.Element{
					"task1": ele1,
					"task2": ele2,
				},
			},
			want: []*Task{task1, task2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemoryStore{
				Tasks:    tt.fields.Tasks,
				TasksMap: tt.fields.TasksMap,
			}
			if got := m.GetAllTasks(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllTasks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryStore_GetDueTasks(t *testing.T) {
	var empty []*Task
	now := time.Now()
	trigger1, _ := NewIntervalTrigger(now, EmptyDateTime, time.Minute)
	trigger2, _ := NewIntervalTrigger(now, EmptyDateTime, time.Hour)

	task1 := NewTask("task1", trigger1, func(args ...interface{}) {})
	task2 := NewTask("task2", trigger2, func(args ...interface{}) {})

	taskList := list.New()
	ele1 := taskList.PushBack(task1)
	ele2 := taskList.PushBack(task2)

	type fields struct {
		Tasks    *list.List
		TasksMap map[string]*list.Element
	}
	type args struct {
		now time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*Task
	}{
		{
			name: "get due empty",
			fields: fields{
				Tasks: taskList,
				TasksMap: map[string]*list.Element{
					"task1": ele1,
					"task2": ele2,
				},
			},
			args: args{
				now: now,
			},
			want: empty,
		},
		{
			name: "get due all",
			fields: fields{
				Tasks: taskList,
				TasksMap: map[string]*list.Element{
					"task1": ele1,
					"task2": ele2,
				},
			},
			args: args{
				now: now.Add(time.Second),
			},
			want: []*Task{task1, task2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemoryStore{
				Tasks:    tt.fields.Tasks,
				TasksMap: tt.fields.TasksMap,
			}
			if got := m.GetDueTasks(tt.args.now); len(got) != len(tt.want) {
				t.Errorf("GetDueTasks() = %v, want %v", got, tt.want)
			}
			if got := m.GetDueTasks(tt.args.now); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDueTasks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryStore_GetNextRunTime(t *testing.T) {
	now := time.Now()
	trigger1, _ := NewIntervalTrigger(now, EmptyDateTime, time.Minute)
	trigger2, _ := NewIntervalTrigger(now, EmptyDateTime, time.Hour)

	task1 := NewTask("task1", trigger1, func(args ...interface{}) {})
	task2 := NewTask("task2", trigger2, func(args ...interface{}) {})

	taskList := list.New()
	ele1 := taskList.PushBack(task1)
	ele2 := taskList.PushBack(task2)

	type fields struct {
		Tasks    *list.List
		TasksMap map[string]*list.Element
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Time
	}{
		{
			name: "get",
			fields: fields{
				Tasks: taskList,
				TasksMap: map[string]*list.Element{
					"task1": ele1,
					"task2": ele2,
				},
			},
			want: now,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemoryStore{
				Tasks:    tt.fields.Tasks,
				TasksMap: tt.fields.TasksMap,
			}
			if got := m.GetNextRunTime(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNextRunTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryStore_GetTaskByName(t *testing.T) {
	now := time.Now()
	trigger1, _ := NewIntervalTrigger(now, EmptyDateTime, time.Minute)
	trigger2, _ := NewIntervalTrigger(now, EmptyDateTime, time.Hour)

	task1 := NewTask("task1", trigger1, func(args ...interface{}) {})
	task2 := NewTask("task2", trigger2, func(args ...interface{}) {})

	taskList := list.New()
	ele1 := taskList.PushBack(task1)
	ele2 := taskList.PushBack(task2)

	type fields struct {
		Tasks    *list.List
		TasksMap map[string]*list.Element
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
			name: "query",
			fields: fields{
				Tasks: taskList,
				TasksMap: map[string]*list.Element{
					"task1": ele1,
					"task2": ele2,
				},
			},
			args: args{
				name: "task1",
			},
			want:    task1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemoryStore{
				Tasks:    tt.fields.Tasks,
				TasksMap: tt.fields.TasksMap,
			}
			got, err := m.GetTaskByName(tt.args.name)
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

func TestMemoryStore_UpdateTask(t *testing.T) {
	now := time.Now()
	trigger1, _ := NewIntervalTrigger(now, EmptyDateTime, time.Minute)
	trigger2, _ := NewIntervalTrigger(now, EmptyDateTime, time.Hour)

	task1 := NewTask("task1", trigger1, func(args ...interface{}) {})
	task1.PreviousRunTime = now
	task2 := NewTask("task2", trigger2, func(args ...interface{}) {})
	task2.PreviousRunTime = now

	taskList := list.New()
	ele1 := taskList.PushBack(task1)
	ele2 := taskList.PushBack(task2)

	type fields struct {
		Tasks    *list.List
		TasksMap map[string]*list.Element
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
			name: "test",
			fields: fields{
				Tasks: taskList,
				TasksMap: map[string]*list.Element{
					"task1": ele1,
					"task2": ele2,
				},
			},
			args: args{
				task: task1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemoryStore{
				Tasks:    tt.fields.Tasks,
				TasksMap: tt.fields.TasksMap,
			}
			if err := m.UpdateTask(tt.args.task); (err != nil) != tt.wantErr {
				t.Errorf("UpdateTask() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
