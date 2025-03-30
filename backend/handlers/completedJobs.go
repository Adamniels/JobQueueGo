package handlers

import (
	//"encoding/json"
	//"fmt"
	"net/http"
)

func CompletedJobs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}
	// TODO:
}
