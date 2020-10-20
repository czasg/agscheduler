package AGScheduler

type WorksMap map[string]WorkDetail

type WorkDetail struct {
	Func func([]interface{})
	Args []interface{}
}
