package queue

type Job struct {
	Id       string `json:"id"`
	Type     string `json:"type"`
	Duration int    `json:"duration,omitempty"` // används för sleep
	Input    string `json:"input,omitempty"`    // används för hash
}
