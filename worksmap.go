package AGScheduler

import "errors"

var WorksMap = map[string]WorkDetail{}

type WorkDetail struct {
	Func func(args ...interface{})
	Args []interface{}
}

func RegisterWorksMap(worksMap map[string]WorkDetail) error {
	if len(WorksMap) == 0 {
		WorksMap = worksMap
		return nil
	}
	for name, _ := range worksMap {
		_, ok := WorksMap[name]
		if ok {
			return errors.New(name + " has existed")
		}
	}
	for name, value := range worksMap {
		WorksMap[name] = value
	}
	return nil
}
