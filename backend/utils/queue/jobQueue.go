package queue

import (
	"JobQueueGo/utils/types"
	"fmt"
	"sync"
)

type JobQueue struct {
	lock  sync.Mutex
	queue []types.Job
}

// Vad vill jag ha för funktioner

// create/new
func NewJobQueue() *JobQueue {
	return &JobQueue{
		queue: make([]types.Job, 0),
	}
}

func (q *JobQueue) Enqueue(job types.Job) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.queue = append(q.queue, job)
}

func (q *JobQueue) Dequeue() (types.Job, bool) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if len(q.queue) == 0 {
		return types.Job{}, false
	}

	firstJob := q.queue[0]
	q.queue = q.queue[1:]

	return firstJob, true
}

func (q *JobQueue) Length() int {
	q.lock.Lock()
	defer q.lock.Unlock()
	return len(q.queue)
}

func (q *JobQueue) GetAll() []types.Job {
	q.lock.Lock()
	defer q.lock.Unlock()

	jobsCopy := make([]types.Job, len(q.queue))
	copy(jobsCopy, q.queue)
	return jobsCopy
}

func (q *JobQueue) PrintJobQueue() {
	for _, job := range q.queue {
		fmt.Println(job)
	}
}
