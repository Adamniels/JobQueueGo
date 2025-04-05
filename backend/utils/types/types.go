package types

/*
type Result struct {
	RespType string `json:"restType"` // "result"
	Type     string `json:"type"`     // "type of job"
	JobId    string `json:"jobId"`    // koppla till rätt jobb
	Result   string `json:"result"`   // resultatet ifall jobbet returnerar det
	Input    string `json:"input"`    // så att det kan läggas till i kön igen om det misslyckas
	Duration int64  `json:"duration"` // hur lång tid jobbet tog i ms
	Success  bool   `json:"success"`  // om det lyckades eller inte
}
*/

type MsgFromWorker struct {
	RespType string
	Res      Result
}
type Result struct {
	Job      Job    `json:"job"`      // jobbet
	Result   string `json:"result"`   // valfritt: kan vara text, hash etc.
	Duration int64  `json:"duration"` // hur lång tid jobbet tog i ms
	Success  bool   `json:"success"`
}
type Job struct {
	Id       string `json:"id"`
	Type     string `json:"type"`
	Duration int    `json:"duration,omitempty"`
	Input    string `json:"input,omitempty"`
	Attempts int    `json:"attempts"`
}
