//go:build js && wasm

// +build: js,wasm
package signalingserverconn

import (
	"encoding/json"
	"fmt"
	"os"
	"syscall/js"
	. "webrtc-full-mesh/utils"

	"github.com/AbdelrahmanWM/signalingserver/signalingserver/message"
)

const signalingServerURL = "ws://localhost:8090/signalingserver"

type Connection interface {
	SetLocalDescription(input json.RawMessage) error
	SetRemoteDescription(input json.RawMessage) error
	AddICECandidate(input json.RawMessage) error
	SendPendingICECandidates() error
}
type SignalingServerConn struct {
	socket    js.Value
	peerID    string
	peerConns map[string]Connection
}

func NewSignalingServerConn(peerConns map[string]Connection) *SignalingServerConn {
	return &SignalingServerConn{peerConns: peerConns}
}
func (conn *SignalingServerConn) SindIdentifySelfMessage() {
	identifySelfMsgContent := message.IdentifySelfContent{ID: ""}
	identifySelfMsgContentJson, err := json.Marshal(identifySelfMsgContent)
	if err != nil {
		Log(fmt.Sprintf("Error parsing message content %v", err))
		os.Exit(1)
	}
	identifySelfMsg := message.Message{
		Kind:    message.IdentifySelf,
		Reach:   message.Self,
		PeerID:  "",
		Content: identifySelfMsgContentJson,
	}
	identifySelfMsgJson, err := json.Marshal(identifySelfMsg)
	if err != nil {
		Log(fmt.Sprintf("Error marshalling message %v", err))
		os.Exit(1)
	}
	conn.Send(string(identifySelfMsgJson))
}
func (conn *SignalingServerConn) Connect(v js.Value, p []js.Value) any {
	socket := js.Global().Get("WebSocket").New(signalingServerURL)
	socket.Set("onopen", js.FuncOf(conn.handleSocketOnOpen))
	socket.Set("onmessage", js.FuncOf(conn.handleSocketOnMessage))
	socket.Set("onclose", js.FuncOf(conn.handleSocketOnClose))
	socket.Set("onError", js.FuncOf(conn.handleSocketOnError))
	conn.socket = socket
	Log("WebSocket connection attempted")
	return nil
}
func (conn *SignalingServerConn) Disconnect(v js.Value, p []js.Value) any {
	if conn.socket.Get("readyState").Int() == 1 {
		conn.socket.Call("close")
		Log("WebSocket connection closed.")
	} else {
		Log("WebSocket is not open, cannot disconnect.")
	}
	return nil
}
func (conn *SignalingServerConn) Socket() js.Value {
	return conn.socket
}
func (conn *SignalingServerConn) Send(message string) error {
	if conn.socket.Get("readyState").Int() != 1 {
		return fmt.Errorf("Websocket is not open")
	}
	conn.socket.Call("send", message)
	return nil
}
func (conn *SignalingServerConn) handleSocketOnOpen(v js.Value, p []js.Value) any {
	Log("Websocket connected!")
	return nil
}
func (conn *SignalingServerConn) handleSocketOnMessage(v js.Value, p []js.Value) any {
	event := p[0]
	messageData := event.Get("data").String() // Get message from event
	var msg message.Message
	err := json.Unmarshal([]byte(messageData), &msg)
	if err != nil {
		Log("Error on unmarshaling message: " + err.Error())
		return nil
	}
	Log("(" + msg.Sender + ") " + string(msg.Content))
	switch msg.Kind {
	case message.IdentifySelf:
		var identifyMsgContent message.IdentifySelfContent
		err := json.Unmarshal(msg.Content, &identifyMsgContent)
		if err != nil {
			Log("Error unmarshaling message content " + err.Error())
		}
		conn.peerID = identifyMsgContent.ID

	case message.Offer:
		targetPeer := msg.Sender
		targetPeerConn, ok := conn.peerConns[targetPeer]
		if !ok {
			Log("Error setting remote description")
			Log(fmt.Sprintf("%s->%#v", targetPeer, conn.peerConns))
			break
		}
		err := targetPeerConn.SetRemoteDescription(msg.Content)
		if err != nil {
			Log(fmt.Sprintf("Error setting remote description: %v", err))
		}
		err = targetPeerConn.SendPendingICECandidates()
		if err != nil {
			Log("Error sending ICE candidates: " + err.Error())
		}
	case message.Answer:
		targetPeer := msg.Sender
		targetPeerConn, ok := conn.peerConns[targetPeer]
		if !ok {
			Log("Error setting remote description")
			break
		}
		err := targetPeerConn.SetRemoteDescription(msg.Content)
		if err != nil {
			Log(fmt.Sprintf("Error setting remote description: %v", err))
		}
		err = targetPeerConn.SendPendingICECandidates()
		if err != nil {
			Log("Error sending ICE candidates: " + err.Error())
		}
	case message.ICECandidate:
		targetPeer := msg.Sender
		targetPeerConn, ok := conn.peerConns[targetPeer]
		if !ok {
			Log("Error setting remote description")
			break
		}
		if candidateErr := targetPeerConn.AddICECandidate(msg.Content); candidateErr != nil {
			Log("Error adding ICE candidate: " + candidateErr.Error())
		}
	}

	return nil
}
func (conn *SignalingServerConn) handleSocketOnClose(v js.Value, p []js.Value) any {
	Log("Connection closed.")
	return nil
}
func (conn *SignalingServerConn) handleSocketOnError(v js.Value, p []js.Value) any {
	Log("Error with websocket connection")
	return nil
}
func (conn *SignalingServerConn) AddNewPeerConnection(id string, peerConn Connection) {
	conn.peerConns[id] = peerConn
}
func (conn *SignalingServerConn) RemovePeerConnection(id string) {
	delete(conn.peerConns, id)
}
func (conn *SignalingServerConn) PeerID() string {
	return conn.peerID
}
