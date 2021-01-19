package agscheduler

import "time"

type RedisStore struct{}

func (ps *RedisStore) GetSchedulingJobs(now time.Time) ([]*Job, error) {
	return nil, nil
}

func (ps *RedisStore) GetJobByName(name string) (*Job, error) {
	return nil, nil
}

func (ps *RedisStore) GetAllJobs() ([]*Job, error) {
	return nil, nil
}

func (ps *RedisStore) AddJob(job *Job) error {
	return nil
}

func (ps *RedisStore) DelJob(job *Job) error {
	return nil
}

func (ps *RedisStore) UpdateJob(job *Job) error {
	return nil
}

func (ps *RedisStore) GetNextRunTime() (time.Time, error) {
	return MinDateTime, nil
}
