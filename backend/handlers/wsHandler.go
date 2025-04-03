package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"JobQueueGo/utils"
	"JobQueueGo/utils/queue"
	"JobQueueGo/utils/resultstore"
	"JobQueueGo/utils/types"

	"github.com/gorilla/websocket"
)

// create an instance of a upgrader that allows connections from anyone
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	workers     = make(map[*websocket.Conn]*utils.Worker)
	id      int = 0
)

func MakeWorkerWebSocketHandler(jobQueue *queue.JobQueue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

			var result types.Result
			err = json.Unmarshal(message, &result)
			if err != nil {
				// Inte JSON – skriv ut som vanlig text
				fmt.Printf("Textmeddelande från worker: %s\n", string(message))
				continue
			}

			if !result.Success {
				job := queue.Job{
					Id:    result.JobId,
					Type:  result.Type,
					Input: result.Input,
					// TODO: borde bara skicka jobbet fram och tillbaka också istället för att ha med alla delar bara
					// Attempts: ,
				}
				jobQueue.Enqueue(job)
			}

			// det är ett giltigt JSON-resultat skriv ut resultatet
			fmt.Printf("Result från worker %s:\n", worker.ID)
			fmt.Printf("  Job type:  %s\n", result.Type)
			fmt.Printf("  Job ID:    %s\n", result.JobId)
			fmt.Printf("  Result:    %s\n", result.Result)
			fmt.Printf("  Duration:  %d ms\n", result.Duration)
			fmt.Printf("  Success:   %t", result.Success)

			resultstore.SaveResult(types.Result{
				Type:     result.Type,
				JobId:    result.JobId,
				Result:   result.Result,
				Duration: result.Duration,
				Success:  result.Success,
			})

			utils.SetWorkerFree(conn)

			// Här kan jag svara med något
			conn.WriteMessage(websocket.TextMessage, []byte("ack result was recived"))
		}
	}
}
