package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/czaorz/agscheduler"
	"io/ioutil"
	"net/http"
	"time"
)

var scheduler agscheduler.AGScheduler

type HttpTask struct {
	Url    string
	Method string
}

func (ht HttpTask) Run(ctx context.Context) {
	request, err := http.NewRequestWithContext(context.Background(), ht.Method, ht.Url, new(bytes.Buffer))
	if err != nil {
		panic(err)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[%s][%d][%v][%s]\n", ht.Method, response.StatusCode, time.Now(), ht.Url)
}

type HttpHandler struct{}

func (h HttpHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	switch request.RequestURI {
	case "/":
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
	case "/add":
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		requestHttpTask := struct {
			Name string
			HttpTask
		}{}
		err = json.Unmarshal(body, &requestHttpTask)
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		job := agscheduler.Job{
			Name: requestHttpTask.Name,
			Trigger: agscheduler.IntervalTrigger{
				Interval: time.Second * 30,
			},
			Task: HttpTask{
				Url:    requestHttpTask.Url,
				Method: requestHttpTask.Method,
			},
		}
		err = scheduler.AddJob(&job)
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
	case "/delete":
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		job := agscheduler.Job{}
		err = json.Unmarshal(body, &job)
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		err = scheduler.DelJob(&job)
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		_, _ = writer.Write([]byte("delete ok"))
	case "/update":
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		job := agscheduler.Job{}
		err = json.Unmarshal(body, &job)
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		err = scheduler.UpdateJob(&job)
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		_, _ = writer.Write([]byte("update ok"))
	default:
		_, _ = writer.Write([]byte("hello default."))
	}
}

func Serve() {
	err := http.ListenAndServe(":8080", HttpHandler{})
	if err != nil {
		panic(err)
	}
}

func main() {
	agscheduler.RegisterAllTasks(&HttpTask{})
	scheduler = agscheduler.AGScheduler{
		Store: &agscheduler.PostgresStore{},
	}
	go scheduler.Start()
	Serve()
}
