package signalingserver

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/AbdelrahmanWM/signalingserver/signalingserver/message"
	"github.com/AbdelrahmanWM/signalingserver/utils"
	"github.com/gorilla/websocket"
)

type SignalingServer struct {
	peers map[string]*websocket.Conn

	// Length of each peer ID
	idLength int

	webSocketUpgrader websocket.Upgrader

	// This flag determines whether the signaling server should prepend the sender's ID to the message's content.
	// If set to true, the server would include the ID of the peer sending the message (i.e., the sender's connID)
	// before the actual message content.
	prependMessageSenderIDToMessage bool

	// This flag controls whether the server should identify a message being sent to the same peer
	// (i.e., the message is being sent to the "self" peer) by appending a label like (self) or (To self) to the message.
	identifySelfInMessages bool

	// This flag controls whether the server should include the requesting peer ID
	// to the 'GetAllPeerIDs' message.
	addSelfToGetPeerIDs bool
}

func NewSignalingServer(id_length int, prependMessageSenderIDToMessage, identifySelfInMessages, addSelfToGetAllPeerIDs bool) *SignalingServer {
	peers := make(map[string]*websocket.Conn)
	var webSocketUpgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // all origins for now
		},
	}
	return &SignalingServer{peers, id_length, webSocketUpgrader, prependMessageSenderIDToMessage, identifySelfInMessages, addSelfToGetAllPeerIDs}
}

func (s *SignalingServer) generateRandomID() string {
	return utils.GenerateRandomID(s.idLength)
}
func (s *SignalingServer) upgradeToWebSocketConn(responseWriter http.ResponseWriter, request *http.Request, responseHeader http.Header) (*websocket.Conn, error) {
	return s.webSocketUpgrader.Upgrade(responseWriter, request, responseHeader)
}
func (s *SignalingServer) GetAllPeerIDs() []string {
	var keys []string
	for key, _ := range s.peers {
		keys = append(keys, key)
	}
	return keys
}
func (s *SignalingServer) HandleWebSocketConn(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgradeToWebSocketConn(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade to webSocket connection", err)
	}
	connID := s.generateRandomID()
	s.peers[connID] = conn
	log.Println("New socket connection: ", connID)
	defer conn.Close()
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			if closeErr, ok := err.(*websocket.CloseError); ok {
				switch closeErr.Code {
				case websocket.CloseNormalClosure, websocket.CloseGoingAway:
					log.Println("Client disconnected gracefully")
				default:
					log.Printf("Web socket closed with unexpected code %d: %v\n", closeErr.Code, err)
				}
				delete(s.peers, connID)
			} else {
				log.Printf("Error reading message: %v\n", err)
			}
			return
		}
		var msg *message.Message = &message.Message{}
		err = json.Unmarshal(p, msg)
		if err != nil {
			log.Printf("Error unmarshaling message %v\n", err)
			conn.WriteMessage(websocket.TextMessage, []byte("Invalid message format"))
			continue
		}
		switch msg.Kind {
		case message.GetAllPeerIDs:
			peerIDs := s.GetAllPeerIDs()
			if !s.addSelfToGetPeerIDs {
				selfIndex := slices.Index(peerIDs, connID)
				if selfIndex != -1 {
					peerIDs = slices.Delete(peerIDs, selfIndex, selfIndex+1)
				} else {
					log.Printf("connID %v not found in peerIDs", connID)
				}
			}
			payload, err := json.Marshal(peerIDs)
			if err != nil {
				log.Printf("Error marshalling peer IDs: %v", err)
				conn.WriteMessage(websocket.TextMessage, []byte("Failed to fetch peer IDs"))
				continue
			}

			if s.addSelfToGetPeerIDs && s.identifySelfInMessages {
				index := strings.Index(string(payload), connID)
				payload = []byte(string(payload[0:index]) + "(self) " + string(payload[index:]))
			}
			err = conn.WriteMessage(websocket.TextMessage, payload)
			if err != nil {
				log.Println("Failed to send Peers IDs")
			}
		case message.SendToAllPeers:
			if s.prependMessageSenderIDToMessage {
				msg.Content = "Message from " + connID + ": " + msg.Content // identifying the sender
			}
			for peerID, peerConn := range s.peers {
				if peerID != connID {
					err = peerConn.WriteMessage(websocket.TextMessage, []byte(msg.Content))
					if err != nil {
						log.Printf("Failed to send message to peer %s: %v\n", peerID, err)
					}
				}
			}
		case message.SendToOnePeer:
			peerConn, exist := s.peers[msg.PeerID]
			if !exist {
				log.Printf("Peer ID %s does not exist\n", msg.PeerID)
				conn.WriteMessage(websocket.TextMessage, []byte("Peer ID doesn't exist"))
				continue
			}
			if s.identifySelfInMessages && msg.PeerID == connID { // if the sender send to himself
				msg.Content = "(To self) " + msg.Content
			}
			if s.prependMessageSenderIDToMessage {
				msg.Content = "From " + connID + ": " + msg.Content // identifying the sender
			}
			err = peerConn.WriteMessage(websocket.TextMessage, []byte(msg.Content))

			if err != nil {
				log.Printf("Failed to send message to peer %s: %v\n", msg.PeerID, err)
				conn.WriteMessage(websocket.TextMessage, []byte("Failed to send message to peer"))
			}
		case message.Disconnect:
			log.Printf("Disconnect message received from %s", connID)
			delete(s.peers, connID)
		default:
			log.Printf("unexpected message.MessageType: %v", msg.Kind)
			conn.WriteMessage(websocket.TextMessage, []byte("Unexpected message type"))
			continue
		}
	}
}
