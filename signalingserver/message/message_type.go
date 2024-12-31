package message

import (
	"encoding/json"
	"fmt"
)
type MessageType int

const (
	GetAllPeerIDs MessageType = iota
	TextMessage
	Disconnect
	Offer  // webrtc specific
	Answer // webrtc specific
	ICECandidate // webrtc specific
	End
)
func (m MessageType) MarshalJSON() ([]byte, error) {
	switch m {
	case GetAllPeerIDs:
		return json.Marshal("GetAllPeerIDs")
	case TextMessage:
		return json.Marshal("TextMessage")
	case Disconnect:
		return json.Marshal("Disconnect")
	case Offer:
		return json.Marshal("Offer")
	case Answer:
		return json.Marshal("Answer")
	case ICECandidate:
		return json.Marshal("ICECandidate")
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
	case "TextMessage":
		*m = TextMessage
	case "Disconnect":
		*m = Disconnect
	case "Offer":
		*m = Offer
	case "Answer":
		*m = Answer
	case "ICECandidate":
		*m = ICECandidate
	default:
		return fmt.Errorf("unknown MessageType string: %s", s)
	}
	return nil
}

type GetAllPeerIDsContent struct {
	PeersIDs []string `json:"peersIDs"`
}
type TextMessageContent struct {
	Message string `json:"message"`
}
type DisconnectContent struct{

}
type OfferContent struct {
	Type string  `json:"type"`
	SDP string `json:"sdp"`
}
type AnswerContent struct {
    Type string `json:"type"`
    SDP  string `json:"sdp"`
}
type ICECandidateContent struct {
	Candidate string `json:"candidate"`
	SdpMid string `json:"sdpMid"`
	SdpMLineIndex int `json:"sdpMLineIndex"`
}