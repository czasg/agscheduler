package agscheduler

import (
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
		return []*Job{}, err
	}
	return jobs, nil
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
