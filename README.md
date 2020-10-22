# AGScheduler

##### Framework
* Scheduler: 调度核心
* Task: 任务模块
* Store： 任务存储模块
    * Memory (Done): 存储在内存中
    * Postgrosql (Todo): 存储在PG中
    * Redis (Todo): 存储在Redis中
* Trigger: 任务触发模块
    * Date (Done): 执行一次后删除任务，参考`time.NewTimer()`
    * Interval (Done): 定期执行，参考`time.NewTicker()`
    * Cron (Done): 根据cron指令周期性执行任务

##### how to use
执行date
```golang
func main() {
	store := stores.NewMemoryStore()
	scheduler := schedulers.NewScheduler(AGScheduler.WorksMap{}, store)

	dateTrigger, _ := triggers.NewDateTrigger(time.Now().Add(time.Second * 1))
	task := tasks.NewTask("task1", func(args []interface{}) {}, []interface{}{"this", "is", "task1"}, dateTrigger)
	_ = scheduler.AddTask(task)

	go func() {
		time.Sleep(time.Second * 10)
		os.Exit(0)
	}()
	scheduler.Start()
}
```

执行interval
```golang
func main() {
	store := stores.NewMemoryStore()
	scheduler := schedulers.NewScheduler(AGScheduler.WorksMap{}, store)

	now := time.Now()
	trigger, _ := triggers.NewIntervalTrigger(now.Add(time.Second*1), now.Add(time.Second*30), time.Second*5)
	task := tasks.NewTask("task1", func(args []interface{}) {}, []interface{}{"this", "is", "task1"}, trigger)

	_ = scheduler.AddTask(task)

	go func() {
		time.Sleep(time.Second * 60)
		os.Exit(0)
	}()
	scheduler.Start()
}
```
执行cron
```golang
func main() {
	store := stores.NewMemoryStore()
	scheduler := schedulers.NewScheduler(AGScheduler.WorksMap{}, store)

	cronTrigger, _ := triggers.NewCronTrigger("*/5 * * * *")

	task := tasks.NewTask("task1", taskFunc, []interface{}{"this", "is", "task1"}, cronTrigger)
	_ = scheduler.AddTask(task)

	go func() {
		time.Sleep(time.Second * 60)
		os.Exit(0)
	}()
	scheduler.Start()
}
```
停止和启动任务
```golang
func main() {
	now := time.Now()
	scheduler := schedulers.NewScheduler(AGScheduler.WorksMap{}, stores.NewMemoryStore())

	trigger1, _ := triggers.NewIntervalTrigger(now.Add(time.Second*1), AGScheduler.EmptyDateTime, time.Second*5)
	task1 := tasks.NewTask("task1", func(args []interface{}) {}, []interface{}{"this", "is", "task1"}, trigger1)
	_ = scheduler.AddTask(task1)

	go func() {
		time.Sleep(time.Second * 10)
		task1.Pause()
		fmt.Println("Pause", time.Now())
		time.Sleep(time.Second * 20)
		fmt.Println("Resume", time.Now())
		task1.Resume()
	}()

	scheduler.Start()
}
```


