package message

import (
	"encoding/json"
	"fmt"
	"log"
)

type Message struct {
	Kind    MessageType     `json:"kind"`
	Reach   ReachType       `json:"reach"`
	Sender  string          `json:"sender"`
	PeerID  string          `json:"peerID"`
	Content json.RawMessage `json:"content"`
}

func (m *Message) UnmarshalContent() (interface{}, error) {
	switch m.Kind {
	case GetAllPeerIDs:
		var getPeerIDs GetAllPeerIDsContent
		err := json.Unmarshal(m.Content, &getPeerIDs)
		return getPeerIDs, err
	case TextMessage:
		var textMessage TextMessageContent
		err := json.Unmarshal(m.Content, &textMessage)
		return textMessage, err
	case Disconnect:
		var disconnect DisconnectContent
		err := json.Unmarshal(m.Content, &disconnect)
		return disconnect, err
	case Offer:
		var offer OfferContent
		err := json.Unmarshal(m.Content, &offer)
		return offer, err
	case Answer:
		var answer AnswerContent
		err := json.Unmarshal(m.Content, &answer)
		return answer, err
	case ICECandidate:
		var iceCandidate ICECandidateContent
		err := json.Unmarshal(m.Content, &iceCandidate)
		return iceCandidate, err
	default:
		log.Printf("Invalid message kind %d\n", m.Kind)
		return nil, fmt.Errorf("invalid message kind %d", m.Kind)
	}
}
