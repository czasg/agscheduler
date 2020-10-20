package stores

import (
	"fmt"
	"github.com/CzaOrz/AGScheduler/interfaces"
	"time"
)

type MemoryStore struct {
	Tasks []interfaces.ITask
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (m *MemoryStore) GetDueTasks() []interfaces.ITask {
	fmt.Println(len(m.Tasks), m.Tasks)
	return m.Tasks
}
func (m *MemoryStore) GetAllTasks() []interfaces.ITask {
	return m.Tasks
}
func (m *MemoryStore) AddTask(task interfaces.ITask) error {
	m.Tasks = append(m.Tasks, task)
	return nil
}
func (m *MemoryStore) DelTask(task interfaces.ITask) error {
	return nil
}
func (m *MemoryStore) NextRunTime() (time.Time, error) {
	return time.Time{}, nil
}
