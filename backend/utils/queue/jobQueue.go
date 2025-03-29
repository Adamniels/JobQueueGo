package queue

import (
	"fmt"
	"sync"
)

type JobQueue struct {
	lock  sync.Mutex
	queue []Job
}

// Vad vill jag ha för funktioner

// create/new
func NewJobQueue() *JobQueue {
	return &JobQueue{
		queue: make([]Job, 0),
	}
}

func (q *JobQueue) Enqueue(job Job) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.queue = append(q.queue, job)
}

func (q *JobQueue) Dequeue() (Job, bool) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if len(q.queue) == 0 {
		return Job{}, false
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

func (q *JobQueue) GetAll() []Job {
	q.lock.Lock()
	defer q.lock.Unlock()

	jobsCopy := make([]Job, len(q.queue))
	copy(jobsCopy, q.queue)
	return jobsCopy
}

func (q *JobQueue) PrintJobQueue() {
	for _, job := range q.queue {
		fmt.Println(job)
	}
}
