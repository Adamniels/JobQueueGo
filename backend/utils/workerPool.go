package utils

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Worker representerar en aktiv WebSocket-anslutning
type Worker struct {
	Conn *websocket.Conn
	Busy bool
	ID   string // kan vara UUID, IP-adress eller något annat
}

// pool hanterar alla aktiva workers
var (
	mu      sync.Mutex
	workers []*Worker
)

// Add lägger till en ny worker i poolen
func AddWorker(worker *Worker) {
	mu.Lock()
	defer mu.Unlock()
	workers = append(workers, worker)
}

// Remove tar bort en worker (t.ex. vid nedkoppling)
func RemoveWorker(conn *websocket.Conn) {
	mu.Lock()
	defer mu.Unlock()
	for i, w := range workers {
		if w.Conn == conn {
			workers = append(workers[:i], workers[i+1:]...)
			break
		}
	}
}

// GetFreeWorker returnerar första lediga workern
func GetFreeWorker() *Worker {
	mu.Lock()
	defer mu.Unlock()
	for _, w := range workers {
		if !w.Busy {
			w.Busy = true // markera som upptagen
			return w
		}
	}
	return nil // ingen ledig
}

// SetWorkerFree sätter Busy=false igen
func SetWorkerFree(conn *websocket.Conn) {
	mu.Lock()
	defer mu.Unlock()
	for _, w := range workers {
		if w.Conn == conn {
			w.Busy = false
			break
		}
	}
}
