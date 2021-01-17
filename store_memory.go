package agscheduler

import (
	"container/list"
	"fmt"
	"github.com/sirupsen/logrus"
	"reflect"
	"sync"
	"time"
)

var JobsMapLock = sync.Mutex{}

type MemoryStore struct {
	Jobs    *list.List
	JobsMap map[string]*list.Element
	Logger  *logrus.Entry
}

func (ms *MemoryStore) FillByDefault() {
	if ms.Jobs == nil {
		ms.Jobs = list.New()
	}
	if ms.JobsMap == nil {
		ms.JobsMap = map[string]*list.Element{}
	}
	if ms.Logger == nil {
		ms.Logger = AGSLog.WithFields(GenASGModule("storeMemory"))
	}
}

func (ms *MemoryStore) GetSchedulingJobs(now time.Time) ([]*Job, error) {
	ms.FillByDefault()
	jobs := []*Job{}
	for el := ms.Jobs.Front(); el != nil; el = el.Next() {
		elJob, ok := el.Value.(*Job)
		if !ok {
			return []*Job{}, fmt.Errorf("invalid store type[%s]", reflect.TypeOf(el.Value).String())
		}
		if elJob.NextRunTime.Before(now) {
			jobs = append(jobs, elJob)
			continue
		}
		break
	}
	ms.Logger.WithFields(logrus.Fields{
		"Func":    "GetSchedulingJobs",
		"JobsNum": len(jobs),
	}).Debugln("OK")
	return jobs, nil
}

func (ms *MemoryStore) GetJobByName(name string) (*Job, error) {
	ms.FillByDefault()
	el, ok := ms.JobsMap[name]
	if !ok {
		return nil, fmt.Errorf("job[%s] not exist", name)
	}
	job, ok := el.Value.(*Job)
	if !ok {
		return nil, fmt.Errorf("invalid store type[%s]", reflect.TypeOf(el.Value).String())
	}
	return job, nil
}

func (ms *MemoryStore) GetAllJobs() ([]*Job, error) {
	ms.FillByDefault()
	var allJobs []*Job
	for el := ms.Jobs.Front(); el != nil; el = el.Next() {
		task, ok := el.Value.(*Job)
		if !ok {
			return []*Job{}, fmt.Errorf("invalid store type[%s]", reflect.TypeOf(el.Value).String())
		}
		allJobs = append(allJobs, task)
	}
	ms.Logger.WithFields(logrus.Fields{
		"Func":    "GetAllJobs",
		"JobsNum": len(allJobs),
	}).Debugln("OK")
	return allJobs, nil
}

func (ms *MemoryStore) AddJob(job *Job) error {
	ms.FillByDefault()
	JobsMapLock.Lock()
	defer JobsMapLock.Unlock()
	for el := ms.Jobs.Front(); el != nil; el = el.Next() {
		elJob, ok := el.Value.(*Job)
		if !ok {
			return fmt.Errorf("invalid store type[%s]", reflect.TypeOf(el.Value).String())
		}
		if job.NextRunTime.After(elJob.NextRunTime) {
			continue
		}
		ms.JobsMap[job.Name] = ms.Jobs.InsertBefore(job, el)
		return nil
	}
	ms.JobsMap[job.Name] = ms.Jobs.PushBack(job)
	ms.Logger.WithFields(logrus.Fields{
		"Func":    "AddJob",
		"JobName": job.Name,
	}).Debugln("OK")
	return nil
}

func (ms *MemoryStore) DelJob(job *Job) error {
	ms.FillByDefault()
	el, ok := ms.JobsMap[job.Name]
	if !ok {
		return fmt.Errorf("job[%s] not exist", job.Name)
	}
	JobsMapLock.Lock()
	defer JobsMapLock.Unlock()
	delete(ms.JobsMap, job.Name)
	ms.Jobs.Remove(el)
	ms.Logger.WithFields(logrus.Fields{
		"Func":    "DelJob",
		"JobName": job.Name,
	}).Debugln("OK")
	return nil
}

func (ms *MemoryStore) UpdateJob(job *Job) error {
	ms.FillByDefault()
	element, ok := ms.JobsMap[job.Name]
	if !ok {
		return fmt.Errorf("job[%s] not exist", job.Name)
	}
	JobsMapLock.Lock()
	defer JobsMapLock.Unlock()
	for el := ms.Jobs.Front(); el != nil; el = el.Next() {
		elJob, ok := el.Value.(*Job)
		if !ok {
			return fmt.Errorf("invalid store type[%s]", reflect.TypeOf(el.Value).String())
		}
		if elJob.Name == job.Name {
			continue
		}
		if job.NextRunTime.After(elJob.NextRunTime) {
			continue
		}
		ms.Jobs.MoveBefore(element, el)
		return nil
	}
	ms.Jobs.MoveToBack(element)
	return nil
}

func (ms *MemoryStore) GetNextRunTime() (time.Time, error) {
	ms.FillByDefault()
	if ms.Jobs.Len() == 0 {
		return MaxDateTime, nil
	}
	job, ok := ms.Jobs.Front().Value.(*Job)
	if !ok {
		return MaxDateTime, fmt.Errorf("invalid store type[%s]", reflect.TypeOf(ms.Jobs.Front().Value).String())
	}
	return job.NextRunTime, nil
}
