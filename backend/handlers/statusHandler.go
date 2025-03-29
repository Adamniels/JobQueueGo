package handlers

import (
	"encoding/json"
	"net/http"

	"JobQueueGo/utils/queue"
)

func MakeStatusHandler(jobQueue *queue.JobQueue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
			return
		}

		statusResp := map[string]any{
			"queue_length": jobQueue.Length(),
			"jobs":         jobQueue.GetAll(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(statusResp)
	}
}
