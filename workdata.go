package AGScheduler

import (
	"math"
	"time"
)

var EmptyDateTime time.Time
var MaxDateTime = time.Now().Add(time.Duration(math.MaxInt64))

type WorksMap map[string]WorkDetail

type WorkDetail struct {
	Func   func(args []interface{}, kwargs map[string]interface{})
	Args   []interface{}
	KwArgs map[string]interface{}
}

type CronTriggerState struct {
	Name string
}
