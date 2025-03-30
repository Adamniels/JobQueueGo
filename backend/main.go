package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"JobQueueGo/handlers"
	"JobQueueGo/utils"
	"JobQueueGo/utils/queue"
)

func main() {
	fmt.Println("hello world")

	jobQueue := queue.NewJobQueue()
	router := setupRouter(jobQueue)

	go matchWorkersWithJobs(jobQueue)

	fmt.Println("Listening on https://localhost:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", router))
}

func setupRouter(queue *queue.JobQueue) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/job", handlers.MakeJobHandler(queue))
	mux.HandleFunc("/status", handlers.MakeStatusHandler(queue))
	mux.HandleFunc("completedJobs", handlers.CompletedJobs)

	mux.HandleFunc("/ws/worker", handlers.WorkerWebSocketHandler)

	return mux
}

func matchWorkersWithJobs(jobQueue *queue.JobQueue) {
	for {
		worker := utils.GetFreeWorker()
		if worker == nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		job, ok := jobQueue.Dequeue()
		if !ok {
			utils.SetWorkerFree(worker.Conn) // markera workern ledig igen
			time.Sleep(100 * time.Millisecond)
			continue
		}

		err := worker.Conn.WriteJSON(job)
		if err != nil {
			fmt.Println("Failed to send job:", err)
			utils.RemoveWorker(worker.Conn)
			continue
		}

		fmt.Printf("Assigned job %s to worker %s\n", job.Id, worker.ID)
	}
}
