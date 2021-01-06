package AGScheduler

import (
	"container/list"
	"errors"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/sirupsen/logrus"
	"time"
)

/*
*	MemoryStore
**/
type MemoryStore struct {
	Tasks    *list.List
	TasksMap map[string]*list.Element
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		Tasks:    list.New(),
		TasksMap: map[string]*list.Element{},
	}
}

func (m *MemoryStore) GetDueTasks(now time.Time) []*Task {
	var dueTasks []*Task
	for el := m.Tasks.Front(); el != nil; el = el.Next() {
		task := el.Value.(*Task)
		startTime := task.GetNextFireTime(now)
		if startTime.Before(now) {
			dueTasks = append(dueTasks, task)
			continue
		}
		break
	}
	return dueTasks
}

func (m *MemoryStore) GetTaskByName(name string) (*Task, error) {
	el, ok := m.TasksMap[name]
	if ok {
		task := el.Value.(*Task)
		return task, nil
	}
	return nil, errors.New("not found task")
}

func (m *MemoryStore) GetAllTasks() []*Task {
	var allTasks []*Task
	for el := m.Tasks.Front(); el != nil; el = el.Next() {
		task := el.Value.(*Task)
		allTasks = append(allTasks, task)
	}
	return allTasks
}

func (m *MemoryStore) AddTask(task *Task) error {
	now := time.Now()
	startTime := task.GetNextFireTime(now)
	for el := m.Tasks.Front(); el != nil; el = el.Next() {
		iTask := el.Value.(*Task)
		iTaskRunTime := iTask.GetNextFireTime(now)
		if startTime.After(iTaskRunTime) {
			continue
		}
		m.TasksMap[task.Name] = m.Tasks.InsertBefore(task, el)
		return nil
	}
	ele := m.Tasks.PushBack(task)
	m.TasksMap[task.Name] = ele
	return nil
}

func (m *MemoryStore) DelTask(task *Task) error {
	element, ok := m.TasksMap[task.Name]
	if !ok {
		return errors.New("not found task in TasksMap")
	}
	delete(m.TasksMap, task.Name)
	m.Tasks.Remove(element)
	return nil
}

func (m *MemoryStore) UpdateTask(task *Task) error {
	element, ok := m.TasksMap[task.Name]
	if !ok {
		return errors.New("not found task in TasksMap")
	}
	nextStartTime := task.NextRunTime
	for el := m.Tasks.Front(); el != nil; el = el.Next() {
		iTask := el.Value.(*Task)
		if iTask.Name == task.Name {
			continue
		}
		endTime := iTask.NextRunTime
		if nextStartTime.After(endTime) {
			continue
		}
		m.Tasks.MoveBefore(element, el)
		return nil
	}
	m.Tasks.MoveToBack(element)
	return nil
}

func (m *MemoryStore) GetNextRunTime() time.Time {
	if m.Tasks.Len() == 0 {
		return EmptyDateTime
	}
	task := m.Tasks.Front().Value.(*Task)
	return task.NextRunTime
}

/*
*	PostGreSQL
**/
type PgStore struct {
	Pg     *pg.DB
	Logger *logrus.Entry
}

func FromPgTask(pgTask *Task, trigger ITrigger, workDetail WorkDetail) *Task {
	var task *Task
	if len(pgTask.Args) == 0 {
		task = NewTask(pgTask.Name, trigger, workDetail.Func, workDetail.Args...)
	} else {
		task = NewTask(pgTask.Name, trigger, workDetail.Func, pgTask.Args...)
	}
	task.TriggerState = pgTask.TriggerState
	task.WorkKey = pgTask.WorkKey
	task.Running = pgTask.Running
	task.Coalesce = pgTask.Coalesce
	task.Count = pgTask.Count
	return task
}

func NewPgStore(pg *pg.DB) (*PgStore, error) {
	if len(WorksMap) == 0 {
		return nil, errors.New("PG instance need define WorksMap")
	}
	err := pg.Model((*Task)(nil)).CreateTable(&orm.CreateTableOptions{
		IfNotExists: true,
	})
	if err != nil {
		return nil, err
	}
	return &PgStore{
		Pg: pg,
		Logger: logrus.WithFields(logrus.Fields{
			"Module": "stores.PgStore",
		}),
	}, nil
}

func (p *PgStore) GetDueTasks(now time.Time) []*Task {
	var dueTasks []*Task
	err := p.Pg.Model(&dueTasks).Where("next_run_time <= ?", now).Select()
	if err != nil {
		p.Logger.WithFields(logrus.Fields{
			"Func": "GetDueTasks",
		}).Errorln(err)
		return dueTasks
	}
	for index, task := range dueTasks {
		workDetail, ok := WorksMap[task.WorkKey]
		if !ok {
			continue
		}
		trigger, err := FromTriggerState(task.TriggerState)
		if err != nil {
			continue
		}
		dueTasks[index] = FromPgTask(task, trigger, workDetail)
	}
	return dueTasks
}

func (p *PgStore) GetTaskByName(name string) (*Task, error) {
	task := Task{}
	err := p.Pg.Model(&task).Where("name = ?", name).Returning("*").Select()
	if err != nil {
		return nil, err
	}
	workDetail, ok := WorksMap[task.WorkKey]
	if !ok {
		return nil, errors.New("WorksKey not existed")
	}
	trigger, err := FromTriggerState(task.TriggerState)
	if err != nil {
		return nil, err
	}
	return FromPgTask(&task, trigger, workDetail), nil
}

func (p *PgStore) GetAllTasks() []*Task {
	var allTasks []*Task
	err := p.Pg.Model(&allTasks).Select()
	if err != nil {
		p.Logger.WithFields(logrus.Fields{
			"Func": "GetAllTasks",
		}).Errorln(err)
		return allTasks
	}
	for index, task := range allTasks {
		workDetail, ok := WorksMap[task.WorkKey]
		if !ok {
			continue
		}
		trigger, err := FromTriggerState(task.TriggerState)
		if err != nil {
			continue
		}
		allTasks[index] = FromPgTask(task, trigger, workDetail)
	}
	return allTasks
}

func (p *PgStore) AddTask(task *Task) error {
	if task.WorkKey == "" {
		return errors.New("task not define WorkKey")
	}
	_, ok := WorksMap[task.WorkKey]
	if !ok {
		return errors.New("WorkKey not existed")
	}
	exist, err := p.Pg.Model(&Task{}).Where("name = ?", task.Name).Exists()
	if err != nil {
		return err
	}
	if exist {
		return errors.New("task has existed")
	}
	task.TriggerState = task.Trigger.GetTriggerState()
	_, err = p.Pg.Model(task).Insert()
	return err
}

func (p *PgStore) DelTask(task *Task) error {
	_, err := p.Pg.Model(task).Where("name = ?name").Delete()
	if err != nil {
		return err
	}
	return nil
}

func (p *PgStore) UpdateTask(task *Task) error {
	task.TriggerState = task.Trigger.GetTriggerState()
	_, err := p.Pg.Model(task).Where("name = ?name").Update()
	if err != nil {
		return err
	}
	return nil
}

func (p *PgStore) GetNextRunTime() time.Time {
	task := Task{}
	err := p.Pg.Model(&task).Order("next_run_time ASC").Returning("next_run_time").Limit(1).Select()
	if err != nil {
		return EmptyDateTime
	}
	return task.NextRunTime
}
