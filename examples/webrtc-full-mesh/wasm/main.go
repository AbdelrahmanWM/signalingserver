//go:build js && wasm

// +build: js,wasm
package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"
	"webrtc-full-mesh/signalingserverconn"
	"webrtc-full-mesh/utils"
	. "webrtc-full-mesh/utils"
	"webrtc-full-mesh/webrtcpeerconn"

	"github.com/AbdelrahmanWM/signalingserver/signalingserver/message"
	"github.com/pion/webrtc/v4"
)

type Peer struct {
	signalingServerConn *signalingserverconn.SignalingServerConn
	peerConnections     map[string]*webrtcpeerconn.PeerConnection
}

func NewPeer() *Peer {
	peerConnections := make(map[string]*webrtcpeerconn.PeerConnection)
	connectionMap := make(map[string]signalingserverconn.Connection)
	for key, pc := range peerConnections {
		connectionMap[key] = pc
	}
	signalingServerConn := signalingserverconn.NewSignalingServerConn(connectionMap)

	return &Peer{signalingServerConn, peerConnections}
}
func (p *Peer) ConnectToSignalingServer(v js.Value, pp []js.Value) any {
	return p.signalingServerConn.Connect(v, pp)
}
func (p *Peer) DisconnectFromSignalingServer(v js.Value, pp []js.Value) any {
	return p.signalingServerConn.Disconnect(v, pp)
}

func (p *Peer) NewPeerConnection(peerConnectionID string) error {
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}
	conn, err := webrtcpeerconn.NewPeerConnection(&config, p.signalingServerConn, peerConnectionID)
	if err != nil {
		return err
	}
	p.AddNewPeerConnection(peerConnectionID, conn)
	return nil
}
func (peer *Peer) SendOffer(peerID string) error {
	return peer.peerConnections[peerID].SendOffer()
}
func (peer *Peer) SendAnswer(peerID string) error {
	return peer.peerConnections[peerID].SendAnswer()
}

func (peer *Peer) NewPeerConnectionJS(v js.Value, p []js.Value) any {
	Log("Attempting to establish new peer connection...")
	peerID := utils.GetElementByID("peerIDInput").Get("value").String()
	if peerID == "" {
		return nil
	}
	err := peer.NewPeerConnection(peerID)
	if err != nil {
		Log("Error establishing new peer connection: " + err.Error())
	}
	Log(fmt.Sprintf("Successfully established a new peer connection: '%s'<->'%s'", peer.signalingServerConn.PeerID(), peerID))
	peer.RenderAllPeerConnections() ///////
	return nil
}
func (peer *Peer) GetAllPeerIDs(v js.Value, p []js.Value) any {
	if peer.signalingServerConn.Socket().IsUndefined() {
		Log("Socket connection not found.")
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
		Log("Error marshalling message:" + err.Error())
		return nil
	}
	Log("Sending message: " + string(msgJSON))
	peer.signalingServerConn.Send(string(msgJSON))
	return nil
}
func (p *Peer) RenderAllPeerConnections() {
	document := js.Global().Get("document")
	parentDiv := utils.GetElementByID("peerConnections")
	parentDiv.Set("innerHTML", "")

	for id, peerConnection := range p.peerConnections {

		div := document.Call("createElement", "div")
		// input:=document.Call("createElement","input")
		peerID := document.Call("createElement", "p")
		peerID.Set("innerHTML", id)
		div.Call("appendChild", peerID)

		offerButton := document.Call("createElement", "button")
		offerButton.Call("addEventListener", "click", js.FuncOf(func(this js.Value, p []js.Value) any {
			go func() { peerConnection.SendOffer() }()

			return nil
		}))
		offerButton.Set("innerHTML", "Send offer message")
		div.Call("appendChild", offerButton)

		answerButton := document.Call("createElement", "button")
		answerButton.Call("addEventListener", "click", js.FuncOf(func(this js.Value, p []js.Value) any {
			go func() { peerConnection.SendAnswer() }()
			return nil
		}))
		answerButton.Set("innerHTML", "Send answer message")
		div.Call("appendChild", answerButton)

		sendMessage := document.Call("createElement", "button")
		sendMessage.Call("addEventListener", "click", js.FuncOf(peerConnection.SendMessageJS))
		sendMessage.Set("innerHTML", "Send message via webrtc")
		div.Call("appendChild", sendMessage)

		parentDiv.Call("appendChild", div)
	}
}
func (p *Peer) SendToAll(v js.Value, pp []js.Value) any {
	message := utils.GetElementByID("message").Get("value").String()
	for _, pc := range p.peerConnections {
		pc.SendMessage([]byte(message))
	}
	return nil
}
func (p *Peer)SindIdentifySelfMessageJS(v js.Value, pp[]js.Value)any{
	p.signalingServerConn.SindIdentifySelfMessage()
	return nil
}
func main() {
	peer := NewPeer()
	Log("New peer!")
	js.Global().Set("connectToSignalingServer", js.FuncOf(peer.ConnectToSignalingServer))
	js.Global().Set("disconnectFromSignalingServer", js.FuncOf(peer.DisconnectFromSignalingServer))
	js.Global().Set("getAllPeerIDs", js.FuncOf(peer.GetAllPeerIDs))
	js.Global().Set("newPeerConnection", js.FuncOf(peer.NewPeerConnectionJS))
	js.Global().Set("clearLog", js.FuncOf(utils.ClearLog))
	js.Global().Set("sendToAll", js.FuncOf(peer.SendToAll))
	js.Global().Set("sindIdentifySelfMessage",js.FuncOf(peer.SindIdentifySelfMessageJS))
	select {}
}
func (p *Peer) AddNewPeerConnection(id string, pc *webrtcpeerconn.PeerConnection) {
	p.peerConnections[id] = pc
	p.signalingServerConn.AddNewPeerConnection(id, pc)
}
func (p *Peer) RemovePeerConnection(id string) {
	delete(p.peerConnections, id)
	p.signalingServerConn.RemovePeerConnection(id)
}
