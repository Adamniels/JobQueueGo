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

var workers = make(map[*websocket.Conn]*utils.Worker)

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

			var msgFromWorker types.MsgFromWorker
			err = json.Unmarshal(message, &msgFromWorker)
			if err != nil {
				// Inte JSON – skriv ut som vanlig text
				fmt.Printf("Textmeddelande från worker: %s\n", string(message))
				continue
			}

			msgFromWorker.Res.Job.Attempts++

			if !msgFromWorker.Res.Success && msgFromWorker.Res.Job.Attempts < 3 {
				jobQueue.Enqueue(msgFromWorker.Res.Job)
			}

			// det är ett giltigt JSON-resultat skriv ut resultatet
			fmt.Printf("Result från worker %s:\n", worker.ID)
			fmt.Printf("  Job type:  %s\n", msgFromWorker.RespType)
			fmt.Printf("  Job ID:    %s\n", msgFromWorker.Res.Job.Id)
			fmt.Printf("  Result:    %s\n", msgFromWorker.Res.Result)
			fmt.Printf("  Duration:  %d ms\n", msgFromWorker.Res.Duration)
			fmt.Printf("  Success:   %t\n", msgFromWorker.Res.Success)
			fmt.Printf("  Attempts:  %d\n", msgFromWorker.Res.Job.Attempts)

			resultstore.SaveResult(types.Result{
				Job:      msgFromWorker.Res.Job,
				Result:   msgFromWorker.Res.Result,
				Duration: msgFromWorker.Res.Duration,
			})

			utils.SetWorkerFree(conn)

			// Här kan jag svara med något
			conn.WriteMessage(websocket.TextMessage, []byte("ack result was recived"))
		}
	}
}
