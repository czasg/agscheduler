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
	State      string
	WorksMap   AGScheduler.WorksMap
	StoresMap  map[string]interfaces.IStore
	Logger     *logrus.Entry
	Controller *Controller
}

func NewScheduler(worksMap AGScheduler.WorksMap, store interfaces.IStore) *Scheduler {
	return &Scheduler{
		WorksMap: worksMap,
		StoresMap: map[string]interfaces.IStore{
			"default": store,
		},
		Logger: logrus.WithFields(logrus.Fields{
			"Module": "AGScheduler.Scheduler",
		}),
		Controller: NewController(),
	}
}

func (s *Scheduler) Start() {
	logger := s.Logger.WithFields(logrus.Fields{
		"Func": "Start",
	})
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
			logger.Warning("AGScheduler server closed")
			os.Exit(1)
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
								logger.WithFields(logrus.Fields{
									"TaskName": dueTask.GetName(),
								}).WithError(err).Errorln("del task failure")
							} else {
								logger.Info("del task success: " + dueTask.GetName())
							}
							continue
						}
						err := store.UpdateTask(dueTask, now)
						if err != nil {
							logger.WithFields(logrus.Fields{
								"TaskName": dueTask.GetName(),
							}).WithError(err).Errorln("update task failure")
						}
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
			{
				if nextCallTime.Equal(AGScheduler.EmptyDateTime) {
					logger.Info("wait task")
					nextCallTime = AGScheduler.MaxDateTime // block until new task to wake
				}
				s.Controller.Reset(nextCallTime)
				select {
				case <-s.Controller.Deadline.Done():
					continue
				}
			}
		}
	}
}

func (s *Scheduler) AddWorksMap(worksMap AGScheduler.WorksMap) error {
	logger := s.Logger.WithFields(logrus.Fields{
		"Func": "AddWorksMap",
	})
	if s.WorksMap == nil {
		s.Logger.Info("empty works map")
		for workName, _ := range worksMap {
			logger.Info("add map work: " + workName)
		}
		s.WorksMap = worksMap
		return nil
	}
	for workName, workDetail := range worksMap {
		_, ok := s.WorksMap[workName]
		if ok {
			logger.Warning("ignore conflict map work: " + workName)
			continue
		}
		logger.Info("add map work: " + workName)
		s.WorksMap[workName] = workDetail
	}
	return nil
}

func (s *Scheduler) AddTask(task interfaces.ITask) error {
	logger := s.Logger.WithFields(logrus.Fields{
		"Func": "AddTask",
	})
	taskName := task.GetName()
	_, ok := s.WorksMap[taskName]
	if ok {
		return errors.New(taskName + " is conflict with TasksMap")
	}
	err := s.StoresMap["default"].AddTask(task)
	if err != nil {
		return err
	}
	logger.Info("add task success: " + taskName)
	task.SetScheduler(s)
	s.Wake()
	return nil
}

func (s *Scheduler) AddTaskFromTasksMap(name, taskMapKey string, args []interface{}, trigger interfaces.ITrigger) error {
	logger := s.Logger.WithFields(logrus.Fields{
		"Func": "AddTaskFromTasksMap",
	})
	_, ok := s.WorksMap[name]
	if ok {
		return errors.New(name + " is conflict with TasksMap")
	}
	detail, ok := s.WorksMap[taskMapKey]
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
	logger.Info("add task success: " + name)
	return nil
}

func (s *Scheduler) DelTask(task interfaces.ITask) error {
	logger := s.Logger.WithFields(logrus.Fields{
		"Func": "DelTask",
	})
	for _, store := range s.StoresMap {
		err := store.DelTask(task)
		if err != nil {
			return err
		}
	}
	logger.Info("del task success: " + task.GetName())
	return nil
}

func (s *Scheduler) SetStore(name string, store interfaces.IStore) {
	s.StoresMap[name] = store
}

func (s *Scheduler) Wake() {
	s.Controller.Cancel()
}
