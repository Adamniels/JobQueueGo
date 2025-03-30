package resultstore

import (
	"sync"
)

// TODO: göra om detta till att använda en databas istället så att det kan spara över tid

type Result struct {
	Type     string
	JobId    string
	Result   string
	Duration int64
}

var (
	lock          sync.RWMutex
	storageResult = make(map[string]Result)
)

func SaveResult(res Result) {
	lock.Lock()
	defer lock.Unlock()
	storageResult[res.JobId] = res
}

func GetResultId(id string) (Result, bool) {
	lock.RLock()
	defer lock.RUnlock()
	res, ok := storageResult[id]
	return res, ok
}

func GetAll() ([]Result, bool) {
	lock.RLock()
	defer lock.RUnlock()
	if len(storageResult) == 0 {
		return nil, false
	}
	results := make([]Result, 0, len(storageResult))
	for _, res := range storageResult {
		results = append(results, res)
	}
	return results, true
}
