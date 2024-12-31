package message

import (
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