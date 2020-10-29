package main

//
//import (
//	"fmt"
//	"github.com/CzaOrz/AGScheduler"
//	"time"
//)
//
//type Count struct {
//	Date     int64
//	Interval int64
//	Cron     int64
//}
//
//func AddCount(args ...interface{}) {
//	index := args[0].(int)
//	count := args[1].(*Count)
//	switch index {
//	case 0:
//		count.Date += 1
//	case 1:
//		count.Interval += 1
//	case 2:
//		count.Cron += 1
//	}
//	fmt.Println(count)
//}
//
//func main() {
//	count := Count{0, 0, 0}
//	now := time.Now()
//	date, _ := AGScheduler.NewDateTrigger(now.Add(time.Minute))
//	interval, _ := AGScheduler.NewIntervalTrigger(now, AGScheduler.EmptyDateTime, time.Second*5)
//	cron, _ := AGScheduler.NewCronTrigger("*/10 * * * *")
//
//	dateTask := AGScheduler.NewTask("date", date, AddCount, 0, &count)
//	intervalTask := AGScheduler.NewTask("interval", interval, AddCount, 1, &count)
//	cronTask := AGScheduler.NewTask("cron", cron, AddCount, 2, &count)
//
//	scheduler := AGScheduler.NewScheduler(AGScheduler.WorksMap{}, AGScheduler.NewMemoryStore())
//	_ = scheduler.AddTask(dateTask)
//	_ = scheduler.AddTask(intervalTask)
//	_ = scheduler.AddTask(cronTask)
//
//	scheduler.Start()
//}
