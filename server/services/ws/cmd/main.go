package main

import (
	"log"
	"net/http"

	"github.com/priyanshubiswas1/barista/ws"
)

func main() {
	manager := ws.NewManager()
	http.HandleFunc("/ws", manager.ServeWS)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
