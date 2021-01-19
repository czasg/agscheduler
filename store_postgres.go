package agscheduler

import "time"

type PostgresStore struct{}

func (ps *PostgresStore) GetSchedulingJobs(now time.Time) ([]*Job, error) {
	return nil, nil
}

func (ps *PostgresStore) GetJobByName(name string) (*Job, error) {
	return nil, nil
}

func (ps *PostgresStore) GetAllJobs() ([]*Job, error) {
	return nil, nil
}

func (ps *PostgresStore) AddJob(job *Job) error {
	return nil
}

func (ps *PostgresStore) DelJob(job *Job) error {
	return nil
}

func (ps *PostgresStore) UpdateJob(job *Job) error {
	return nil
}

func (ps *PostgresStore) GetNextRunTime() (time.Time, error) {
	return MinDateTime, nil
}
