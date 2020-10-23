package AGScheduler

import (
	_ "github.com/CzaOrz/AGScheduler/interfaces"
	_ "github.com/CzaOrz/AGScheduler/schedulers"
	_ "github.com/CzaOrz/AGScheduler/stores"
	_ "github.com/CzaOrz/AGScheduler/tasks"
	_ "github.com/CzaOrz/AGScheduler/triggers"
	"math"
	"time"
)

var EmptyDateTime time.Time
var MaxDateTime = time.Now().Add(time.Duration(math.MaxInt64))

type WorksMap map[string]WorkDetail

type WorkDetail struct {
	Func func([]interface{})
	Args []interface{}
}
