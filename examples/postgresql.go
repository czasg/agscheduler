package main

import (
	"fmt"
	"github.com/CzaOrz/AGScheduler"
	"github.com/go-pg/pg/v10"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

var pgdb *pg.DB
var once sync.Once

func InitPostGreSql() {
	once.Do(func() {
		log := logrus.WithFields(logrus.Fields{
			"FUNC": "remote.GetPGInstance",
		})
		pgdb = pg.Connect(&pg.Options{
			Addr:         "localhost:5432",
			User:         "postgres",
			Password:     "postgres",
			Database:     "monitor_edn",
			PoolSize:     3,
			MaxRetries:   3,
			MinIdleConns: 2,
		})
		_, err := pgdb.Exec("SELECT 1")
		if err != nil {
			log.Fatal(err)
		}
		log.Info("init postgresql successful.")
	})
}

type Count struct {
	Date     int64
	Interval int64
	Cron     int64
}

func AddCount(args ...interface{}) {
	index := args[0].(int)
	count := args[1].(*Count)
	switch index {
	case 0:
		count.Date += 1
	case 1:
		count.Interval += 1
	case 2:
		count.Cron += 1
	}
	fmt.Println(count)
}

func main() {
	InitPostGreSql()

	err := AGScheduler.RegisterWorksMap(map[string]AGScheduler.WorkDetail{
		"add": {
			Func: AddCount,
			Args: []interface{}{},
		},
	})
	if err != nil {
		panic(err)
	}

	count := Count{0, 0, 0}
	now := time.Now()
	date, _ := AGScheduler.NewDateTrigger(now.Add(time.Minute))
	interval, _ := AGScheduler.NewIntervalTrigger(now, AGScheduler.EmptyDateTime, time.Second*5)
	cron, _ := AGScheduler.NewCronTrigger("*/10 * * * *")

	dateTask := AGScheduler.NewTask("date", date, AddCount, 0, &count)
	dateTask.WorkKey = "add"
	intervalTask := AGScheduler.NewTask("interval", interval, AddCount, 1, &count)
	intervalTask.WorkKey = "add"
	cronTask := AGScheduler.NewTask("cron", cron, AddCount, 2, &count)
	cronTask.WorkKey = "add"

	pgStore, err := AGScheduler.NewPgStore(pgdb)
	if err != nil {
		panic(err)
	}
	scheduler := AGScheduler.NewScheduler(pgStore)
	_ = scheduler.AddTask(dateTask)
	_ = scheduler.AddTask(intervalTask)
	_ = scheduler.AddTask(cronTask)

	scheduler.Start()
}
