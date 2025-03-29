package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// Matchar din Job-struktur på serversidan
type Job struct {
	Id       string `json:"id"`
	Type     string `json:"type"`
	Duration int    `json:"duration,omitempty"` // används för sleep
	Input    string `json:"input,omitempty"`    // används för hash
}

type Result struct {
	Type     string `json:"type"`     // "result"
	JobId    string `json:"jobId"`    // koppla till rätt jobb
	Result   string `json:"result"`   // valfritt: kan vara text, hash etc.
	Duration int64  `json:"duration"` // hur lång tid jobbet tog i ms
}

func main() {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws/worker", nil)
	if err != nil {
		log.Fatal("WebSocket connection failed:", err)
	}
	defer conn.Close()

	log.Println("Connected to server!")

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
			Type:     "result",
			JobId:    job.Id,
			Result:   resultText,
			Duration: elapsed,
		}

		msg, _ := json.Marshal(result)
		conn.WriteMessage(websocket.TextMessage, msg)
		log.Println("Sent result:", string(msg))
	}
}
