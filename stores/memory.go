package stores

import (
	"container/list"
	"errors"
	"github.com/CzaOrz/AGScheduler"
	"github.com/CzaOrz/AGScheduler/interfaces"
	"time"
)

type MemoryStore struct {
	Tasks    list.List
	TasksMap map[string]*list.Element
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		TasksMap: map[string]*list.Element{},
	}
}

func (m *MemoryStore) GetDueTasks(now time.Time) []interfaces.ITask {
	var dueTasks []interfaces.ITask
	for el := m.Tasks.Front(); el != nil; el = el.Next() {
		task := el.Value.(interfaces.ITask)
		startTime := task.GetNextRunTime(now)
		if startTime.Before(now) {
			dueTasks = append(dueTasks, task)
			continue
		}
		break
	}
	return dueTasks
}

func (m *MemoryStore) GetAllTasks() []interfaces.ITask {
	var allTasks []interfaces.ITask
	for el := m.Tasks.Front(); el != nil; el = el.Next() {
		task := el.Value.(interfaces.ITask)
		allTasks = append(allTasks, task)
	}
	return allTasks
}

func (m *MemoryStore) GetTask(name string) (interfaces.ITask, error) {
	el, ok := m.TasksMap[name]
	if ok {
		task := el.Value.(interfaces.ITask)
		return task, nil
	}
	return nil, errors.New("not found task")
}

func (m *MemoryStore) AddTask(task interfaces.ITask) error {
	now := time.Now()
	startTime := task.GetNextRunTime(now)
	for el := m.Tasks.Front(); el != nil; el = el.Next() {
		iTask := el.Value.(interfaces.ITask)
		iTaskRunTime := iTask.GetNextRunTime(now)
		if startTime.After(iTaskRunTime) {
			continue
		}
		m.TasksMap[task.GetName()] = m.Tasks.InsertBefore(task, el)
		return nil
	}
	ele := m.Tasks.PushBack(task)
	m.TasksMap[task.GetName()] = ele
	return nil
}

func (m *MemoryStore) DelTask(task interfaces.ITask) error {
	element, ok := m.TasksMap[task.GetName()]
	if !ok {
		return errors.New("not found task in TasksMap")
	}
	delete(m.TasksMap, task.GetName())
	m.Tasks.Remove(element)
	return nil
}

func (m *MemoryStore) UpdateTask(task interfaces.ITask, now time.Time) error {
	element, ok := m.TasksMap[task.GetName()]
	if !ok {
		return errors.New("not found task in TasksMap")
	}
	nextStartTime := task.GetNextRunTime(now)
	for el := m.Tasks.Front(); el != nil; el = el.Next() {
		iTask := el.Value.(interfaces.ITask)
		if iTask.GetName() == task.GetName() {
			continue
		}
		endTime := iTask.GetNextRunTime(now)
		if nextStartTime.After(endTime) {
			continue
		}
		m.Tasks.MoveBefore(element, el)
		return nil
	}
	m.Tasks.MoveToBack(element)
	return nil
}

func (m *MemoryStore) GetNextRunTime(now time.Time) time.Time {
	if m.Tasks.Len() == 0 {
		return AGScheduler.EmptyDateTime
	}
	task := m.Tasks.Front().Value.(interfaces.ITask)
	return task.GetNextRunTime(now)
}
