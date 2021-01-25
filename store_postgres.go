package agscheduler

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/sirupsen/logrus"
	"time"
)

var PG *pg.DB

func NewPG() *pg.DB {
	pG := pg.Connect(&pg.Options{
		Addr:         Config.PG.Addr,
		User:         Config.PG.User,
		Password:     Config.PG.Password,
		Database:     Config.PG.Database,
		PoolSize:     Config.PG.PoolSize,
		MaxRetries:   3,
		MinIdleConns: 2,
	})
	if err := pG.Model((*Job)(nil)).CreateTable(&orm.CreateTableOptions{
		IfNotExists: true,
	}); err != nil {
		AGSLog.WithError(err).Panic("create table error.")
	}
	return pG
}

type PostgresStore struct {
	Logger *logrus.Entry
	PG     *pg.DB
}

func (ps *PostgresStore) FillByDefault() {
	if ps.Logger == nil {
		ps.Logger = AGSLog.WithFields(GenASGModule("storePostgres"))
	}
	if ps.PG == nil {
		if PG == nil {
			PG = NewPG()
		}
		ps.PG = PG.WithContext(AGSContext)
	}
}

func (ps *PostgresStore) GetSchedulingJobs(now time.Time) ([]*Job, error) {
	ps.FillByDefault()
	jobs := []*Job{}
	err := ps.PG.Model(&jobs).Where("next_run_time <= ?", now).Select()
	if err != nil {
		return nil, err
	}
	for index, job := range jobs {
		err = DeserializeTask(job)
		if err != nil {
			return nil, err
		}
		err = DeserializeTrigger(job)
		if err != nil {
			return nil, err
		}
		jobs[index] = job
	}
	return jobs, nil
}

func (ps *PostgresStore) GetJobByName(name string) (*Job, error) {
	ps.FillByDefault()
	job := &Job{}
	err := ps.PG.Model(job).Where("name = ?", name).Returning("*").Select()
	if err != nil {
		return nil, err
	}
	err = DeserializeTrigger(job)
	if err != nil {
		return nil, err
	}
	err = DeserializeTask(job)
	if err != nil {
		return nil, err
	}
	job.FillByDefault()
	return job, nil
}

func (ps *PostgresStore) GetAllJobs() ([]*Job, error) {
	ps.FillByDefault()
	jobs := []*Job{}
	err := ps.PG.Model(&jobs).Select()
	fmt.Println(jobs)
	if err != nil {
		return nil, err
	}
	for index, job := range jobs {
		err = DeserializeTrigger(job)
		if err != nil {
			return nil, err
		}
		err = DeserializeTask(job)
		if err != nil {
			return nil, err
		}
		jobs[index] = job
	}
	return jobs, nil
}

func (ps *PostgresStore) AddJob(job *Job) error {
	ps.FillByDefault()
	job.FillByDefault()
	err := SerializeTask(job)
	if err != nil {
		return err
	}
	err = SerializeTrigger(job)
	if err != nil {
		return err
	}
	exist, err := ps.PG.Model(&Job{}).Where("name = ?", job.Name).Exists()
	if err != nil {
		return err
	}
	if exist {
		return fmt.Errorf("job[%s] has existed.", job.Name)
	}
	_, err = ps.PG.Model(job).Insert()
	if err != nil {
		return err
	}
	return nil
}

func (ps *PostgresStore) DelJob(job *Job) error {
	ps.FillByDefault()
	_, err := ps.PG.Model(job).Where("name = ?name").Delete()
	if err != nil {
		return err
	}
	return nil
}

func (ps *PostgresStore) UpdateJob(job *Job) error {
	ps.FillByDefault()
	job.FillByDefault()
	err := SerializeTask(job)
	if err != nil {
		return err
	}
	err = SerializeTrigger(job)
	if err != nil {
		return err
	}
	_, err = ps.PG.Model(job).Where("name = ?name").Update()
	if err != nil {
		return err
	}
	return nil
}

func (ps *PostgresStore) GetNextRunTime() (time.Time, error) {
	ps.FillByDefault()
	job := Job{}
	err := ps.PG.Model(&job).Order("next_run_time ASC").Returning("next_run_time").Limit(1).Select()
	if err != nil {
		return MinDateTime, err
	}
	return job.NextRunTime, nil
}
