package main

import (
	"fmt"
	"log"
	"net/http"

	"JobQueueGo/handlers"
	"JobQueueGo/utils/queue"
)

func main() {
	fmt.Println("hello world")

	jobQueue := queue.NewJobQueue()
	router := setupRouter(jobQueue)

	fmt.Println("Listening on https://localhost:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", router))
}

func setupRouter(queue *queue.JobQueue) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/job", handlers.MakeJobHandler(queue))
	mux.HandleFunc("/status", handlers.MakeStatusHandler(queue))

	return mux
}
