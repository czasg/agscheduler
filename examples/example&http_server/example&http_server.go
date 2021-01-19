package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/czaorz/agscheduler"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func serverByNetListen() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	conn, err := listener.Accept()
	if err != nil {
		panic(err)
	}

	buf := make([]byte, 1024)
	_, err = conn.Read(buf)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(buf))
}

var scheduler agscheduler.AGScheduler

type HttpHandler struct{}

func (h HttpHandler) Run(ctx context.Context) {}

func (h HttpHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	switch request.RequestURI {
	case "/":
		_, _ = writer.Write([]byte("hello index"))
	case "/scheduler":
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
	case "/transfer":
		req, _ := http.NewRequestWithContext(context.Background(), request.Method, "http://www.baidu.com", request.Body)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		_, _ = writer.Write(body)
	default:
		_, _ = writer.Write([]byte("hello world"))
	}
}

func serverByHttpListenAndServe() {
	err := scheduler.AddJob(&agscheduler.Job{
		Name: "http-test",
		Trigger: agscheduler.IntervalTrigger{
			Interval: time.Minute,
		},
		Task: &HttpHandler{},
	})
	go scheduler.Start()
	err = http.ListenAndServe(":8080", HttpHandler{})
	if err != nil {
		panic(err)
	}
}

func main() {
	serverByHttpListenAndServe()
}
