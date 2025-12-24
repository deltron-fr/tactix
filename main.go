package main

import (
	"log"
	"net/http"
)

func main() {
	startRepl()

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func SetupAPI() {
	manager := NewManager()

	http.HandleFunc("/ws", manager.serveWS)
}
