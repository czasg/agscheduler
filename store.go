package agscheduler

import "time"

type IStore interface {
	GetSchedulingJobs() ([]*Job, error)
	GetJobByName(name string) (*Job, error)
	GetAllJobs() ([]*Job, error)
	AddJob(job *Job) error
	DelJob(job *Job) error
	UpdateJob(job *Job) error
	GetNextRunTime() (time.Time, error)
}
