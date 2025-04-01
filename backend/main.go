package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"JobQueueGo/handlers"
	"JobQueueGo/utils"
	"JobQueueGo/utils/queue"
)

func main() {
	fmt.Println("hello world")

	jobQueue := queue.NewJobQueue()
	router := setupRouter(jobQueue)

	go matchWorkersWithJobs(jobQueue)

	ipaddr, _ := getLocalIP()
	fmt.Printf("listening on ip addr: %s, port: %s\n", ipaddr, "8080")
	fmt.Println("Listening on https://localhost:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", router))
}

func setupRouter(queue *queue.JobQueue) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/job", handlers.MakeJobHandler(queue))
	mux.HandleFunc("/status", handlers.MakeStatusHandler(queue))
	mux.HandleFunc("/completedJobs", handlers.CompletedJobs)

	mux.HandleFunc("/ws/worker", handlers.WorkerWebSocketHandler)

	return mux
}

func matchWorkersWithJobs(jobQueue *queue.JobQueue) {
	for {
		worker := utils.GetFreeWorker()
		if worker == nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		job, ok := jobQueue.Dequeue()
		if !ok {
			utils.SetWorkerFree(worker.Conn) // markera workern ledig igen
			time.Sleep(100 * time.Millisecond)
			continue
		}

		err := worker.Conn.WriteJSON(job)
		if err != nil {
			fmt.Println("Failed to send job:", err)
			utils.RemoveWorker(worker.Conn)
			continue
		}

		fmt.Printf("Assigned job %s to worker %s\n", job.Id, worker.ID)
	}
}

func getLocalIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() || ip.To4() == nil {
				continue
			}

			return ip.String(), nil
		}
	}

	return "", fmt.Errorf("no IP address found")
}
