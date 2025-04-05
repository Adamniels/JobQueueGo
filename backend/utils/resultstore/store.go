package resultstore

import (
	"sync"

	"JobQueueGo/utils/types"
)

// TODO: göra om detta till att använda en databas istället så att det kan spara över tid

var (
	lock          sync.RWMutex
	storageResult = make(map[string]types.Result)
)

func SaveResult(res types.Result) {
	lock.Lock()
	defer lock.Unlock()
	storageResult[res.Job.Id] = res
}

func GetResultId(id string) (types.Result, bool) {
	lock.RLock()
	defer lock.RUnlock()
	res, ok := storageResult[id]
	return res, ok
}

func GetAll() ([]types.Result, bool) {
	lock.RLock()
	defer lock.RUnlock()
	if len(storageResult) == 0 {
		return nil, false
	}
	results := make([]types.Result, 0, len(storageResult))
	for _, res := range storageResult {
		results = append(results, res)
	}
	return results, true
}
