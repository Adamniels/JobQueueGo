package utils

import (
	"fmt"
	"sync"
)

var (
	jobCounter     int64
	jobCounterLock sync.Mutex
)

func GenerateJobID() string {
	jobCounterLock.Lock()
	defer jobCounterLock.Unlock()
	jobCounter++
	return fmt.Sprintf("job-%d", jobCounter)
}
