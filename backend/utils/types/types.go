package types

type Result struct {
	RespType string `json:"restType"` // "result"
	Type     string `json:"type"`     // "type of job"
	JobId    string `json:"jobId"`    // koppla till rätt jobb
	Result   string `json:"result"`   // resultatet ifall jobbet returnerar det
	Input    string `json:"input"`    // så att det kan läggas till i kön igen om det misslyckas
	Duration int64  `json:"duration"` // hur lång tid jobbet tog i ms
	Success  bool   `json:"success"`  // om det lyckades eller inte
}
