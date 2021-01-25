package agscheduler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sync"
)

var (
	allITasks = map[string]reflect.Type{}
	taskLock  = sync.Mutex{}
)

type ITask interface {
	Run(ctx context.Context)
}

func addTask(taskName string, taskType reflect.Type) {
	defer taskLock.Unlock()
	taskLock.Lock()
	_, ok := allITasks[taskName]
	if !ok {
		allITasks[taskName] = taskType
	}
}

func SerializeTask(job *Job) error {
	if job.Task == nil {
		return errors.New("task is nil.")
	}
	if job.TaskMeta == nil {
		job.TaskMeta = map[string]interface{}{}
	}
	body, err := json.Marshal(job.Task)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &job.TaskMeta)
	if err != nil {
		return err
	}
	taskT := reflect.TypeOf(job.Task)
	if taskT.Kind() == reflect.Ptr {
		taskT = taskT.Elem()
	}
	job.TaskMeta[""] = taskT.String()
	return nil
}

func DeserializeTask(job *Job) error {
	if job.TaskMeta == nil {
		return errors.New("task meta is nil.")
	}
	taskNameI, ok := job.TaskMeta[""]
	if !ok {
		return errors.New("task mate not define task name.")
	}
	taskName, ok := taskNameI.(string)
	if !ok {
		return errors.New("task name is not string.")
	}
	taskT, ok := allITasks[taskName]
	if !ok {
		return errors.New("task is not type of reflect.Type.")
	}
	newTaskT := reflect.New(taskT)
	if newTaskT.Kind() == reflect.Ptr {
		newTaskT = newTaskT.Elem()
	}
	for i := 0; i < taskT.NumField(); i++ {
		fieldName := taskT.Field(i).Name
		fieldValue, ok := job.TaskMeta[fieldName]
		if !ok {
			return errors.New("task's field value is not define.")
		}
		newTaskT.Field(i).Set(reflect.ValueOf(fieldValue))
	}
	newTask, ok := newTaskT.Interface().(ITask)
	if !ok {
		return fmt.Errorf("task[%v] is not type of ITask.", newTaskT.String())
	}
	job.Task = newTask
	return nil
}

func RegisterAllTasks(tasks ...ITask) {
	for _, task := range tasks {
		taskT := reflect.TypeOf(task)
		if taskT.Kind() == reflect.Ptr {
			taskT = taskT.Elem()
		}
		addTask(taskT.String(), taskT)
	}
}
