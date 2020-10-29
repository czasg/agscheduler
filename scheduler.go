package AGScheduler

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var EmptyDateTime time.Time
var MaxDateTime = time.Now().Add(time.Duration(math.MaxInt64))

type Scheduler struct {
	StoresMap   map[string]IStore
	Logger      *logrus.Entry
	Controller  *Controller
	CloseCancel context.CancelFunc
}

func NewScheduler(store IStore) *Scheduler {
	return &Scheduler{
		StoresMap: map[string]IStore{
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
	s.CloseCancel = cancel
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch,
			os.Interrupt,
			syscall.SIGINT,
			syscall.SIGTERM,
		)
		exitSignal := <-ch
		logger.Warnf("receive a signal of %v", exitSignal)
		s.Close()
	}()

	for {
		select {
		case <-closeContext.Done():
			goto Exit
		default:
			now := time.Now()
			nextCallTime := time.Time{}
			for _, store := range s.StoresMap {
				{
					dueTasks := store.GetDueTasks(now) // Gets the tasks that should be scheduled
					for _, dueTask := range dueTasks {
						dueTaskRunTime := dueTask.GetNextFireTime(now)
						dueTask.Go(dueTaskRunTime)
						dueTaskNextRunTime := dueTask.GetNextFireTime(now)
						if dueTaskNextRunTime.Equal(EmptyDateTime) {
							err := store.DelTask(dueTask)
							if err != nil {
								logger.WithFields(logrus.Fields{
									"TaskName": dueTask.Name,
								}).WithError(err).Errorln("del task failure")
							} else {
								logger.Info("del task success: " + dueTask.Name)
							}
							continue
						}
						err := store.UpdateTask(dueTask)
						if err != nil {
							logger.WithFields(logrus.Fields{
								"TaskName": dueTask.Name,
							}).WithError(err).Errorln("update task failure")
						}
					}
				}
				{
					nextRunTime := store.GetNextRunTime()
					if nextRunTime.Equal(EmptyDateTime) {
						continue
					}
					if nextCallTime.Equal(EmptyDateTime) {
						nextCallTime = nextRunTime
					}
					if nextCallTime.After(nextRunTime) {
						nextCallTime = nextRunTime
					}
				}
			}
			{
				if nextCallTime.Equal(EmptyDateTime) {
					logger.Info("wait task")
					nextCallTime = MaxDateTime // block until new task to wake
				}
				s.Controller.Reset(nextCallTime)
				<-s.Controller.Deadline.Done()
			}
		}
	}
Exit:
	logger.Warning("AGScheduler server closed")
}

func (s *Scheduler) Close() {
	s.CloseCancel()
	time.Sleep(time.Second)
	s.Wake()
}

func (s *Scheduler) Wake() {
	s.Controller.Cancel()
}

func (s *Scheduler) AddTask(task *Task) error {
	logger := s.Logger.WithFields(logrus.Fields{
		"Func": "AddTask",
	})
	task.Scheduler = s
	taskName := task.Name
	_, ok := WorksMap[taskName]
	if ok {
		return errors.New(taskName + " is conflict with TasksMap")
	}
	err := s.StoresMap["default"].AddTask(task)
	if err != nil {
		return err
	}
	logger.Info("add task success: " + taskName)
	s.Wake()
	return nil
}

func (s *Scheduler) AddTaskFromTasksMap(name, taskMapKey string, trigger ITrigger, args ...interface{}) error {
	logger := s.Logger.WithFields(logrus.Fields{
		"Func": "AddTaskFromTasksMap",
	})
	_, ok := WorksMap[name]
	if ok {
		return errors.New(name + " is conflict with TasksMap")
	}
	detail, ok := WorksMap[taskMapKey]
	if !ok {
		return errors.New(name + " is not define in TasksMap")
	}
	if len(args) == 0 {
		args = detail.Args
	}
	task := NewTask(name, trigger, detail.Func, args...)
	err := s.AddTask(task)
	if err != nil {
		return err
	}
	logger.Info("add task success: " + name)
	return nil
}

func (s *Scheduler) GetTaskByName(name string) (*Task, error) {
	for _, store := range s.StoresMap {
		task, err := store.GetTaskByName(name)
		if err != nil {
			continue
		}
		return task, nil
	}
	return nil, errors.New("not found task")
}

func (s *Scheduler) GetAllTasks() []*Task {
	return s.StoresMap["default"].GetAllTasks()
}

func (s *Scheduler) UpdateTask(task *Task) error {
	logger := s.Logger.WithFields(logrus.Fields{
		"Func": "UpdateTask",
	})
	err := s.StoresMap["default"].UpdateTask(task)
	if err != nil {
		return err
	}
	logger.Info("update task success: " + task.Name)
	s.Wake()
	return nil
}

func (s *Scheduler) DelTask(task *Task) error {
	logger := s.Logger.WithFields(logrus.Fields{
		"Func": "DelTask",
	})
	err := s.StoresMap["default"].DelTask(task)
	if err != nil {
		return err
	}
	logger.Info("del task success: " + task.Name)
	return nil
}
