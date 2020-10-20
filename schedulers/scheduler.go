package schedulers

import (
	"errors"
	"github.com/CzaOrz/AGScheduler"
	"github.com/CzaOrz/AGScheduler/interfaces"
	"github.com/CzaOrz/AGScheduler/stores"
	"github.com/CzaOrz/AGScheduler/tasks"
	"time"
)

type Scheduler struct {
	State     string
	TasksMap  AGScheduler.WorksMap
	StoresMap map[string]interfaces.IStore
}

func NewScheduler(worksMap AGScheduler.WorksMap) *Scheduler {
	store := stores.NewMemoryStore()

	return &Scheduler{
		TasksMap: worksMap,
		StoresMap: map[string]interfaces.IStore{
			"memory": store,
		},
	}
}

func (s *Scheduler) Start() {
	for {
		nextCallTime := 10
		for _, store := range s.StoresMap {
			dueTasks := store.GetDueTasks()
			for _, dueTask := range dueTasks {
				dueTask.Go()
			}
		}
		time.Sleep(time.Duration(nextCallTime) * time.Second)
		break
	}
}

func (s *Scheduler) AddTask(task interfaces.ITask) error {
	taskName := task.GetName()
	_, ok := s.TasksMap[taskName]
	if ok {
		return errors.New(taskName + " is conflict with TasksMap")
	}
	task.SetScheduler(s)
	_ = s.StoresMap["memory"].AddTask(task)
	return nil
}

func (s *Scheduler) AddTaskFromTasksMap(name, taskMapKey string, args []interface{}, trigger interfaces.ITrigger) error {
	_, ok := s.TasksMap[name]
	if ok {
		return errors.New(name + " is conflict with TasksMap")
	}
	detail, ok := s.TasksMap[taskMapKey]
	if !ok {
		return errors.New(name + " is not define in TasksMap")
	}
	if len(args) == 0 {
		args = detail.Args
	}
	task := tasks.NewTask(name, detail.Func, args, trigger)
	err := s.AddTask(task)
	if err != nil {
		return err
	}
	return nil
}

func (s *Scheduler) DelTask(task interfaces.ITask) error {
	return nil
}

func (s *Scheduler) SetStore(name string, store interfaces.IStore) {
	s.StoresMap[name] = store
}
