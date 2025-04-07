package utils

import (
	"fmt"
	"sync"
)

var (
	jobCounter     int64
	jobCounterLock sync.Mutex
)

func SetInitialJobCounter(n int64) {
	jobCounterLock.Lock()
	defer jobCounterLock.Unlock()
	jobCounter = n
}

func GenerateJobID() string {
	jobCounterLock.Lock()
	defer jobCounterLock.Unlock()
	jobCounter++
	return fmt.Sprintf("job-%d", jobCounter)
}

var (
	workerCounter     int64
	workerCounterLock sync.Mutex
)

func GenerateWorkerID() string {
	workerCounterLock.Lock()
	defer workerCounterLock.Unlock()
	workerCounter++
	return fmt.Sprintf("worker-%d", workerCounter)
}
