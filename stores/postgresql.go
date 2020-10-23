package stores

import (
	"github.com/CzaOrz/AGScheduler/interfaces"
	"github.com/go-pg/pg/v10"
	"time"
)

type PgStore struct {
	Pg *pg.DB
}

func NewPgStore(pg *pg.DB) *PgStore {
	return &PgStore{}
}

func (p *PgStore) GetAllTasks() []interfaces.ITask {
	return []interfaces.ITask{}
}

func (p *PgStore) GetTask(name string) interfaces.ITask {
	return nil
}

func (p *PgStore) AddTask(task interfaces.ITask) error {
	return nil
}

func (p *PgStore) DelTask(task interfaces.ITask) error {
	return nil
}

func (p *PgStore) UpdateTask(task interfaces.ITask, now time.Time) error {
	return nil
}

func (p *PgStore) GetNextRunTime(now time.Time) error {
	return nil
}
