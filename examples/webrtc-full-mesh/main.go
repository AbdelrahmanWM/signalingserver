package main

import (
	"log"
	"net/http"

	"github.com/AbdelrahmanWM/signalingserver/signalingserver"
)

func main() {
	signalingServer := signalingserver.NewSignalingServer(10, true, false)
	http.HandleFunc("/signalingserver", signalingServer.HandleWebSocketConn)
	log.Println("Signaling server available at localhost:8090")
	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		log.Fatal("Error serving the signaling server.")
	}

}
