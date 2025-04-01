package main

import (
	"flag"
	"fmt"

	"workers/workers"
)

func main() {
	host := flag.String("host", "localhost", "Serverns IP eller hostname")
	port := flag.String("port", "8080", "Serverns port")
	flag.Parse()

	wsURL := fmt.Sprintf("ws://%s:%s/ws/worker", *host, *port)
	fmt.Println(wsURL)

	workers.Start(wsURL)
}
