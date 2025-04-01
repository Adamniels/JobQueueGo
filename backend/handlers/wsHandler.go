package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"JobQueueGo/utils"
	"JobQueueGo/utils/resultstore"
	"github.com/gorilla/websocket"
)

// create an instance of a upgrader that allows connections from anyone
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Result struct {
	RespType string `json:"restType"` // "result"
	Type     string `json:"type"`     // "type of job"
	JobId    string `json:"jobId"`    // koppla till rätt jobb
	Result   string `json:"result"`   // valfritt: kan vara text, hash etc.
	Duration int64  `json:"duration"` // hur lång tid jobbet tog i ms
}

var (
	workers     = make(map[*websocket.Conn]*utils.Worker)
	id      int = 0
)

func WorkerWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade failed:", err)
		return
	}
	worker := &utils.Worker{Conn: conn, Busy: false, ID: utils.GenerateWorkerID()}
	workers[conn] = worker
	utils.AddWorker(worker)

	defer func() {
		// tar bort worker från worker listan och stänger connection
		utils.RemoveWorker(conn)
		conn.Close()
	}()

	conn.WriteMessage(websocket.TextMessage, []byte("connected to server"))
	fmt.Println("Worker connected via WebSocket")

	// enter the loop that keeps the websocket open
	for {
		// Läs meddelande från worker
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Worker disconnected:", err)
			break
		}

		var result Result
		err = json.Unmarshal(message, &result)
		if err != nil {
			// Inte JSON – skriv ut som vanlig text
			fmt.Printf("Textmeddelande från worker: %s\n", string(message))
			continue
		}

		// TODO: kanske göra switch case ifall worker kan returnera med än bara result

		// det är ett giltigt JSON-resultat skriv ut resultatet
		fmt.Printf("Result från worker %s:\n", worker.ID)
		fmt.Printf("  Job type:  %s\n", result.Type)
		fmt.Printf("  Job ID:    %s\n", result.JobId)
		fmt.Printf("  Result:    %s\n", result.Result)
		fmt.Printf("  Duration:  %d ms\n", result.Duration)

		// TODO: samma typ, borde kunna återanvända den?
		// TODO: Borde även ha med om det lyckades eller inte
		resultstore.SaveResult(resultstore.Result{
			Type:     result.Type,
			JobId:    result.JobId,
			Result:   result.Result,
			Duration: result.Duration,
		})

		utils.SetWorkerFree(conn)

		// Här kan jag svara med något
		conn.WriteMessage(websocket.TextMessage, []byte("ack result was recived"))
	}
}
