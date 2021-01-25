package main

import (
	"context"
	"encoding/json"
	"github.com/czaorz/agscheduler"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

var scheduler agscheduler.AGScheduler

type HttpHandler struct {
	Name string
}

func (h HttpHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	switch request.RequestURI {
	case "/": // curl http://localhost:8080
		jobs, err := scheduler.GetAllJobs()
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		body, err := json.Marshal(jobs)
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		_, _ = writer.Write(body)
	case "/total": // curl http://localhost:8080/total
		jobs, err := scheduler.GetAllJobs()
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		_, _ = writer.Write([]byte(string(len(jobs))))
	case "/job": // curl -d "{\"name\":\"intervalJob2\"}" http://localhost:8080/job
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		req := struct {
			Name string `json:"name"`
		}{}
		err = json.Unmarshal(body, &req)
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		job, err := scheduler.GetJobByJobName(req.Name)
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		body, err = json.Marshal(job)
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		_, _ = writer.Write(body)
	}
}

func (h HttpHandler) Run(ctx context.Context) {
	logrus.WithField("Name", h.Name).Info("run")
}

func main() {
	err := http.ListenAndServe("0.0.0.0:8080", HttpHandler{})
	if err != nil {
		panic(err)
	}
}

func init() {
	now := time.Now()
	dateJob1 := &agscheduler.Job{
		Name: "dateJob1",
		Trigger: &agscheduler.DateTrigger{
			NextRunTime: now.Add(time.Hour),
		},
		Task: HttpHandler{Name: "dateJob1"},
	}
	dateJob2 := &agscheduler.Job{
		Name: "dateJob2",
		Trigger: &agscheduler.DateTrigger{
			NextRunTime: now.Add(time.Minute * 5),
		},
		Task: HttpHandler{Name: "dateJob2"},
	}
	intervalJob1 := &agscheduler.Job{
		Name: "intervalJob1",
		Trigger: &agscheduler.IntervalTrigger{
			Interval: time.Minute,
		},
		Task: HttpHandler{Name: "intervalJob1"},
	}
	intervalJob2 := &agscheduler.Job{
		Name: "intervalJob2",
		Trigger: &agscheduler.IntervalTrigger{
			Interval: time.Minute * 30,
		},
		Task: HttpHandler{Name: "intervalJob2"},
	}
	cronJob1 := &agscheduler.Job{
		Name: "cronJob1",
		Trigger: &agscheduler.CronTrigger{
			CronCmd: "* 1 * * *",
		},
		Task: HttpHandler{Name: "cronJob1"},
	}
	cronJob2 := &agscheduler.Job{
		Name: "cronJob2",
		Trigger: &agscheduler.CronTrigger{
			CronCmd: "* * 1 * *",
		},
		Task: HttpHandler{Name: "cronJob2"},
	}
	scheduler = agscheduler.AGScheduler{}
	err := scheduler.AddJob(dateJob1, dateJob2, intervalJob1, intervalJob2, cronJob1, cronJob2)
	if err != nil {
		panic(err)
	}
	go scheduler.Start()
}
