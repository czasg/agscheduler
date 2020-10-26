package AGScheduler

import (
	"container/list"
	"errors"
	"time"
)

/*
*	MemoryStore
**/
type MemoryStore struct {
	Tasks    list.List
	TasksMap map[string]*list.Element
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		TasksMap: map[string]*list.Element{},
	}
}

func (m *MemoryStore) GetDueTasks(now time.Time) []Task {
	var dueTasks []Task
	for el := m.Tasks.Front(); el != nil; el = el.Next() {
		task := el.Value.(Task)
		startTime := task.GetNextFireTime(now)
		if startTime.Before(now) {
			dueTasks = append(dueTasks, task)
			continue
		}
		break
	}
	return dueTasks
}

func (m *MemoryStore) GetAllTasks() []Task {
	var allTasks []Task
	for el := m.Tasks.Front(); el != nil; el = el.Next() {
		task := el.Value.(Task)
		allTasks = append(allTasks, task)
	}
	return allTasks
}

func (m *MemoryStore) GetTaskByName(name string) (Task, error) {
	el, ok := m.TasksMap[name]
	if ok {
		task := el.Value.(Task)
		return task, nil
	}
	return Task{}, errors.New("not found task")
}

func (m *MemoryStore) AddTask(task Task) error {
	now := time.Now()
	startTime := task.GetNextFireTime(now)
	for el := m.Tasks.Front(); el != nil; el = el.Next() {
		iTask := el.Value.(Task)
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

func (m *MemoryStore) DelTask(task Task) error {
	element, ok := m.TasksMap[task.Name]
	if !ok {
		return errors.New("not found task in TasksMap")
	}
	delete(m.TasksMap, task.Name)
	m.Tasks.Remove(element)
	return nil
}

func (m *MemoryStore) UpdateTask(task Task) error {
	element, ok := m.TasksMap[task.Name]
	if !ok {
		return errors.New("not found task in TasksMap")
	}
	nextStartTime := task.NextRunTime
	for el := m.Tasks.Front(); el != nil; el = el.Next() {
		iTask := el.Value.(Task)
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
	task := m.Tasks.Front().Value.(Task)
	return task.NextRunTime
}
