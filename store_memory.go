package agscheduler

import (
	"container/list"
	"time"
)

type MemoryStore struct {
	Jobs    *list.List
	JobsMap map[string]*list.Element
}

func (ms *MemoryStore) FillByDefault() {
	if ms.Jobs == nil {
		ms.Jobs = list.New()
	}
	if ms.JobsMap == nil {
		ms.JobsMap = map[string]*list.Element{}
	}
}

func (ms *MemoryStore) GetSchedulingJobs() ([]*Job, error) {
	ms.FillByDefault()
	return nil, nil
}

func (ms *MemoryStore) GetJobByName(name string) (*Job, error) {
	ms.FillByDefault()
	return nil, nil
}

func (ms *MemoryStore) GetAllJobs() ([]*Job, error) {
	ms.FillByDefault()
	return nil, nil
}

func (ms *MemoryStore) AddJob(job *Job) error {
	ms.FillByDefault()
	return nil
}

func (ms *MemoryStore) DelJob(job *Job) error {
	ms.FillByDefault()
	return nil
}

func (ms *MemoryStore) UpdateJob(job *Job) error {
	ms.FillByDefault()
	return nil
}

func (ms *MemoryStore) GetNextRunTime() (time.Time, error) {
	ms.FillByDefault()
	return time.Time{}, nil
}
