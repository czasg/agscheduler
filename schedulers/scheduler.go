package schedulers

import (
	"context"
	"errors"
	"github.com/CzaOrz/AGScheduler"
	"github.com/CzaOrz/AGScheduler/interfaces"
	"github.com/CzaOrz/AGScheduler/tasks"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Scheduler struct {
	State     string
	TasksMap  AGScheduler.WorksMap
	StoresMap map[string]interfaces.IStore
	Logger    *logrus.Entry
}

func NewScheduler(worksMap AGScheduler.WorksMap, store interfaces.IStore) *Scheduler {
	return &Scheduler{
		TasksMap: worksMap,
		StoresMap: map[string]interfaces.IStore{
			"memory": store,
		},
		Logger: logrus.WithFields(logrus.Fields{
			"Module": "AGScheduler.Scheduler",
		}),
	}
}

func (s *Scheduler) Start() {
	ticker := time.NewTicker(time.Second)
	closeContext, cancel := context.WithCancel(context.Background())
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch,
			os.Interrupt,
			syscall.SIGINT,
			syscall.SIGTERM,
		)
		<-ch
		cancel()
	}()

	for {
		select {
		case <-closeContext.Done():
			break
		default:
			now := time.Now()
			nextCallTime := time.Time{}
			for _, store := range s.StoresMap {
				{
					dueTasks := store.GetDueTasks(now) // Gets the tasks that should be scheduled
					for _, dueTask := range dueTasks {
						dueTaskRunTime := dueTask.GetNextRunTime(now)
						dueTask.Go(dueTaskRunTime)
						dueTaskNextRunTime := dueTask.GetNextRunTime(now)
						if dueTaskNextRunTime.Equal(AGScheduler.EmptyDateTime) {
							err := store.DelTask(dueTask)
							if err != nil {
								s.Logger.Errorln("del task failure: " + err.Error())
							} else {
								s.Logger.Info("del task success: " + dueTask.GetName())
							}
							continue
						}
						_ = store.UpdateTask(dueTask, now)
					}
				}
				{
					nextRunTime := store.GetNextRunTime(now)
					if nextRunTime.Equal(AGScheduler.EmptyDateTime) {
						continue
					}
					if nextCallTime.Equal(AGScheduler.EmptyDateTime) {
						nextCallTime = nextRunTime
					}
					if nextCallTime.After(nextRunTime) {
						nextCallTime = nextRunTime
					}
				}
			}
			if nextCallTime.Equal(AGScheduler.EmptyDateTime) {
				<-ticker.C
				s.Logger.Debug("wait task")
				continue
			}
			time.Sleep(time.Duration(nextCallTime.Unix()-now.Unix()) * time.Second)
		}
	}
}

func (s *Scheduler) AddTask(task interfaces.ITask) error {
	taskName := task.GetName()
	_, ok := s.TasksMap[taskName]
	if ok {
		return errors.New(taskName + " is conflict with TasksMap")
	}
	err := s.StoresMap["memory"].AddTask(task)
	if err != nil {
		return err
	}
	s.Logger.Info("add task success: " + taskName)
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
	s.Logger.Info("add task success: " + name)
	return nil
}

func (s *Scheduler) DelTask(task interfaces.ITask) error {
	for _, store := range s.StoresMap {
		err := store.DelTask(task)
		if err != nil {
			return err
		}
	}
	s.Logger.Info("del task success: " + task.GetName())
	return nil
}

func (s *Scheduler) SetStore(name string, store interfaces.IStore) {
	s.StoresMap[name] = store
}
