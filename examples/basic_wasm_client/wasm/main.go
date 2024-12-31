//go:build js && wasm

// +build: js,wasm
package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/AbdelrahmanWM/signalingserver/signalingserver/message"
)

var sockets []js.Value
var socket js.Value

func connectToWebSocket(this js.Value, p []js.Value) interface{} {
	for _, skt := range sockets {
		if skt.Get("readyState").Int() == 1 {
			skt.Call("close")
		}
	}
	url := "ws://localhost:8090/signalingserver"   // Websocket url
	socket = js.Global().Get("WebSocket").New(url) // Create websocket instance

	// Define Websocket onopen handler
	socket.Set("onopen", js.FuncOf(func(this js.Value, p []js.Value) any {
		log("Websocket connected!")
		return nil
	}))
	// define webSocket on message event handler
	socket.Set("onmessage", js.FuncOf(func(this js.Value, p []js.Value) any {
		event := p[0]
		messageData := event.Get("data").String() // Get message from event
		var message message.Message
		err := json.Unmarshal([]byte(messageData), &message)
		if err != nil {
			log("Error on unmarshaling message: " + err.Error())
		}
		log("(" + message.Sender + ") " + string(message.Content))
		return nil
	}))

	// Define Websocket onclose event handler

	socket.Set("onclose", js.FuncOf(func(this js.Value, p []js.Value) any {
		log("Websocket connection closed.")
		return nil
	}))

	//Define webSocket on error event handler
	socket.Set("onerror", js.FuncOf(func(this js.Value, p []js.Value) any {
		log("Error with Websocket connection.")
		return nil
	}))
	sockets = append(sockets, socket)
	return nil
}
func disconnectFromWebSocket(v js.Value, p []js.Value) any {
	if socket.Get("readyState").Int() == 1 {
		socket.Call("close")
	}
	return nil
}
func getAllPeerIDs(v js.Value, p []js.Value) any {
	if socket.IsUndefined() {
		log("Socket connection not found.")
		return nil
	}
	getAllPeerIDsMsg := message.Message{
		Kind:    message.GetAllPeerIDs,
		PeerID:  "",
		Content: nil,
		Reach:   message.Self,
		Sender:  "",
	}
	msgJSON, err := json.Marshal(getAllPeerIDsMsg)
	if err != nil {
		log("Error marshalling message:" + err.Error())
		return nil
	}
	log("Sending message: " + string(msgJSON))
	socket.Call("send", string(msgJSON))
	return nil
}
func log(message string) {
	el := getElementByID("logArea")
	el.Set("innerHTML", el.Get("innerHTML").String()+"* "+message+"<br>")
}
func getElementByID(id string) js.Value {
	el := js.Global().Get("document").Call("getElementById", id)
	if el.IsNull() {
		log(fmt.Sprintf("Element with id '%s' not found", id))
	}
	return el
}
func sendToPeer(v js.Value, p []js.Value) any {
	if socket.IsUndefined() {
		log("Socket connection not found.")
		return nil
	}
	peerID := getElementByID("peerID").Get("value").String()
	msg := getElementByID("message").Get("value").String()
	fmt.Println(peerID)
	fmt.Println(msg)
	msgContent, err := json.Marshal(message.TextMessageContent{msg})
	if err != nil {
		log("Error marshalling message")
	}
	messageToPeer := message.Message{
		Kind:    message.TextMessage,
		PeerID:  peerID,
		Reach:   message.OnePeer,
		Content: msgContent,
	}
	msgJSON, err := json.Marshal(messageToPeer)
	if err != nil {
		log("Error marshalling message:" + err.Error())
		return nil
	}
	log("Sending message: " + string(msgJSON))
	socket.Call("send", string(msgJSON))
	return nil

}
func sendToAll(v js.Value, p []js.Value) any {
	if socket.IsUndefined() {
		log("Socket connection not found.")
		return nil
	}
	msg := getElementByID("message").Get("value").String()
	fmt.Println(msg)
	msgContent, err := json.Marshal(message.TextMessageContent{msg})
	if err != nil {
		log("Error marshalling message")
	}
	messageToPeer := message.Message{
		Kind:    message.TextMessage,
		Reach:   message.AllPeers,
		PeerID:  "",
		Content: msgContent,
	}
	msgJSON, err := json.Marshal(messageToPeer)
	if err != nil {
		log("Error marshalling message:" + err.Error())
		return nil
	}
	log("Sending message: " + string(msgJSON))
	socket.Call("send", string(msgJSON))
	return nil

}
func main() {
	js.Global().Set("connectToWebSocket", js.FuncOf(connectToWebSocket))
	js.Global().Set("disconnectFromWebSocket", js.FuncOf(disconnectFromWebSocket))
	js.Global().Set("getAllPeerIDs", js.FuncOf(getAllPeerIDs))
	js.Global().Set("sendToAll", js.FuncOf(sendToAll))
	js.Global().Set("sendToPeer", js.FuncOf(sendToPeer))
	fmt.Println("Go Web Assembly Websocket client")
	select {} // keep the go connection running

}
