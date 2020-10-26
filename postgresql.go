package AGScheduler

import (
	"github.com/go-pg/pg/v10"
	"time"
)

type PgStore struct {
	Pg *pg.DB
}

func NewPgStore(pg *pg.DB) *PgStore {
	return &PgStore{}
}

func (p *PgStore) GetAllTasks() []Task {
	return []Task{}
}

func (p *PgStore) GetTask(name string) Task {
	return Task{}
}

func (p *PgStore) AddTask(task Task) error {
	return nil
}

func (p *PgStore) DelTask(task Task) error {
	return nil
}

func (p *PgStore) UpdateTask(task Task, now time.Time) error {
	return nil
}

func (p *PgStore) GetNextRunTime(now time.Time) error {
	return nil
}
