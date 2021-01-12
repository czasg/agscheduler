package agscheduler

type ITask interface {
	Run()
}

type Job struct {
	Id        int64       `json:"id" pg:",pk"`
	Name      string      `json:"name" pg:",use_zero"`
	Task      ITask       `json:"-" pg:"-"`
	Scheduler AGScheduler `json:"-" pg:"-"`
}
