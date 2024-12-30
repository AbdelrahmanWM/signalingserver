package main

import (
	"log"
	"net/http"

	"github.com/AbdelrahmanWM/signalingserver/signalingserver"
)

func main() {
	signalingServer := signalingserver.NewSignalingServer(20, false, false)
	http.HandleFunc("/signalingserver", signalingServer.HandleWebSocketConn)
	log.Println("Server listening on :8090")
	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		log.Fatal("Server failed:", err)
	}

}
