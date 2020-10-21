package AGScheduler

import "time"

var EmptyDateTime time.Time

type WorksMap map[string]WorkDetail

type WorkDetail struct {
	Func func([]interface{})
	Args []interface{}
}
