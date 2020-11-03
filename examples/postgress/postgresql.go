package main

import (
	"encoding/json"
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

func RequestIns(args ...interface{}) {
	if len(args) != 1 {
		return
	}
	body, err := json.Marshal(args[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	var req Request
	err = json.Unmarshal(body, &req)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(req, time.Now())
}

type Request struct {
	Url     string
	Method  string
	Headers map[string]string
	Count   int64
}

func main() {
	InitPostGreSql()

	err := AGScheduler.RegisterWorksMap(map[string]AGScheduler.WorkDetail{
		"add": {
			Func: RequestIns,
			Args: []interface{}{},
		},
	})
	if err != nil {
		panic(err)
	}

	now := time.Now()
	date, _ := AGScheduler.NewDateTrigger(now.Add(time.Minute))
	interval, _ := AGScheduler.NewIntervalTrigger(now, AGScheduler.EmptyDateTime, time.Second*5)
	cron, _ := AGScheduler.NewCronTrigger("*/10 * * * *")

	dateTask := AGScheduler.NewTask("date", date, RequestIns, Request{Url: "www.date.com", Method: "GET", Headers: map[string]string{"Content-Type": "Application/json"}, Count: 116})
	dateTask.WorkKey = "add"
	intervalTask := AGScheduler.NewTask("interval", interval, RequestIns, Request{Url: "www.interval.com", Method: "GET", Headers: map[string]string{"Content-Type": "Application/json"}, Count: 116})
	intervalTask.WorkKey = "add"
	cronTask := AGScheduler.NewTask("cron", cron, RequestIns, Request{Url: "www.cron.com", Method: "GET", Headers: map[string]string{"Content-Type": "Application/json"}, Count: 116})
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
