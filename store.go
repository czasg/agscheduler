package agscheduler

import "time"

// GetSchedulingJobs: get jobs which should be scheduled now.
// GetJobByName: find job by name.
// GetNextRunTime: get the first schedule time from all jobs.
type IStore interface {
	GetSchedulingJobs(now time.Time) ([]*Job, error)
	GetJobByName(name string) (*Job, error)
	GetAllJobs() ([]*Job, error)
	AddJob(job *Job) error
	DelJob(job *Job) error
	UpdateJob(job *Job) error
	GetNextRunTime() (time.Time, error)
}
