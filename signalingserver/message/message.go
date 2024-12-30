package message

import (
	"encoding/json"
	"fmt"
)

type Message struct {
	Kind    MessageType `json:"kind"`
	PeerID  string      `json:"peerID"`
	Content string      `json:"content"`
}

type MessageType int

const (
	GetAllPeerIDs MessageType = iota
	SendToOnePeer
	SendToAllPeers
	Disconnect
	End
)

func (m MessageType) MarshalJSON() ([]byte, error) {
	switch m {
	case GetAllPeerIDs:
		return json.Marshal("GetAllPeerIDs")
	case SendToOnePeer:
		return json.Marshal("SendToOnePeer")
	case SendToAllPeers:
		return json.Marshal("SendToAllPeers")
	case Disconnect:
		return json.Marshal("Disconnect")
	default:
		return nil, fmt.Errorf("unknown MessageType: %d", m)
	}
}

func (m *MessageType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch s {
	case "GetAllPeerIDs":
		*m = GetAllPeerIDs
	case "SendToOnePeer":
		*m = SendToOnePeer
	case "SendToAllPeers":
		*m = SendToAllPeers
	case "Disconnect":
		*m = Disconnect
	default:
		return fmt.Errorf("unknown MessageType string: %s", s)
	}
	return nil
}
