package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"JobQueueGo/utils"
	"JobQueueGo/utils/queue"
)

func MakeJobHandler(jobQueue *queue.JobQueue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
			return
		}

		var job queue.Job
		err := json.NewDecoder(r.Body).Decode(&job)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		job.Id = utils.GenerateJobID()
		jobQueue.Enqueue(job)

		fmt.Printf("Enqueued job: %+v\n", job)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "enqueued",
			"id":     job.Id,
		})
	}
}
