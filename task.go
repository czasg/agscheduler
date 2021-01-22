package agscheduler

import (
	"context"
	"reflect"
	"sync"
)

var (
	allITasks = map[string]ITask{}
	taskLock  = sync.Mutex{}
)

type ITask interface {
	Run(ctx context.Context)
}

func RegisterAllTasks(tasks ...ITask) {
	for _, task := range tasks {
		reflect.TypeOf(task)
	}
}
