package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"JobQueueGo/utils/resultstore"
)

func CompletedJobs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	storage, _ := resultstore.GetAll()

	statusResp := map[string]any{
		"completedJobs": len(storage),
		"jobs":          storage,
	}

	fmt.Println("all completed jobs was requested")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(statusResp)
}
