package workers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/gorilla/websocket"
)

// Matchar din Job-struktur på serversidan
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

func Start(wsURL string) {
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatal("WebSocket connection failed:", err)
	}
	defer conn.Close()

	log.Println("Connected to server!")
	success := true

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading:", err)
			break
		}

		var job Job
		err = json.Unmarshal(message, &job)
		if err != nil {
			log.Printf("Text message from server: %s\n", string(message))
			continue
		}

		log.Printf("Received job: %+v\n", job)

		start := time.Now()
		var resultText string

		// Jobbtyper
		switch job.Type {
		case "execurl":
			resultText, err = handleExecURLJob(job)
			if err != nil {
				success = false
			}
		case "sleep":
			time.Sleep(time.Duration(job.Duration) * time.Second)
			resultText = fmt.Sprintf("Slept for %d seconds", job.Duration)

		case "hash":
			// Simulera en hash (du kan byta till riktig SHA256)
			resultText = fmt.Sprintf("fakehash(%s)", job.Input)

		default:
			resultText = "Unknown job type"
		}

		elapsed := time.Since(start).Milliseconds()

		// Skicka tillbaka resultatet
		result := Result{
			Job: job,
			Result:   resultText,
			Duration: elapsed,
			Success: success,
		}

		msgFromWorker := MsgFromWorker{
			RespType: "result",
			Res: result,
		}

		msg, _ := json.Marshal(msgFromWorker)
		conn.WriteMessage(websocket.TextMessage, msg)
		log.Println("Sent result:", string(msg))
	}
}

// TODO: kolla på och förstå
func handleExecURLJob(job Job) (string, error) {
	filename := fmt.Sprintf("job_%s_exec.sh", job.Id)

	// Ladda ner från URL
	resp, err := http.Get(job.Input)
	if err != nil {
		return fmt.Sprintf("Download error: %s", err), err
	}
	defer resp.Body.Close()

	out, err := os.Create(filename)
	if err != nil {
		return fmt.Sprintf("File creation error: %s", err), err
	}
	io.Copy(out, resp.Body)
	out.Close()

	os.Chmod(filename, 0o755)

	// Kör filen
	cmd := exec.Command("./" + filename)
	output, err := cmd.CombinedOutput()
	os.Remove(filename)

	if err != nil {
		return fmt.Sprintf("Exec error: %s\nOutput: %s", err, output), err
	}

	return string(output), nil
}
