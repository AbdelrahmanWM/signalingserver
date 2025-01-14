package signalingserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"

	"github.com/AbdelrahmanWM/signalingserver/signalingserver/message"
	"github.com/AbdelrahmanWM/signalingserver/utils"
	"github.com/gorilla/websocket"
)

type SignalingServer struct {
	peers map[string]*websocket.Conn

	// Length of each peer ID
	idLength int

	webSocketUpgrader websocket.Upgrader

	// This flag controls whether the server should identify a message being sent to the same peer
	// (i.e., the message is being sent to the "self" peer) by appending a label like (self) or (To self) to the message.
	identifyMessageSender bool

	// This flag controls whether the server should include the requesting peer ID
	// to the 'GetAllPeerIDs' message.
	addSelfToGetPeerIDs bool
}

func NewSignalingServer(id_length int, identifyMessageSender, addSelfToGetAllPeerIDs bool) *SignalingServer {
	peers := make(map[string]*websocket.Conn)
	var webSocketUpgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // all origins for now
		},
	}
	return &SignalingServer{peers, id_length, webSocketUpgrader, identifyMessageSender, addSelfToGetAllPeerIDs}
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
		var msg message.Message = message.Message{}
		var responseMsg message.Message = message.Message{
			Kind:    message.TextMessage,
			Reach:   message.Self,
			Sender:  connID,
			PeerID:  connID,
			Content: json.RawMessage{},
		}
		if !s.identifyMessageSender {
			responseMsg.Sender = ""
		}
		err = json.Unmarshal(p, &msg)
		if err != nil {
			log.Printf("Error unmarshaling message %v\n", err)
			responseMsg.Content, err = json.Marshal(message.TextMessageContent{"error", "Invalid message structure"})
			if err != nil {
				log.Println("Error marshalling message")
			}
			conn.WriteJSON(responseMsg)
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

			responseMsg.Content, err = json.Marshal(message.GetAllPeerIDsContent{peerIDs})
			if err != nil {
				log.Printf("Error marshalling peer IDs: %v", err)
				responseMsg.Content, err = json.Marshal(message.TextMessageContent{"error", "Failed to fetch peer IDs"})
				if err != nil {
					log.Println("Error marshaling error message")
					continue
				}
				break
			}
			responseMsg.Kind = message.GetAllPeerIDs

		case message.TextMessage, message.Offer, message.Answer, message.ICECandidate:
			responseMsg.Kind = msg.Kind
			responseMsg.Content = msg.Content
			responseMsg.PeerID = msg.PeerID
		case message.Disconnect:
			log.Printf("Disconnect message received from %s", connID)
			delete(s.peers, connID)
			continue
		case message.IdentifySelf:
			responseMsg.Kind = msg.Kind
			responseMsg.Reach = message.Self
			msgContent, err := json.Marshal(message.IdentifySelfContent{connID})
			if err != nil {
				log.Printf("Error marshalling msg content: %v", err)
			}
			responseMsg.Content = msgContent
		default:
			log.Printf("unexpected Message type: %v", msg.Kind)
			conn.WriteMessage(websocket.TextMessage, []byte("Unexpected message type"))
			continue
		}

		if msg.Reach == message.Self || msg.PeerID == connID {
			responseMsg.Sender = "server"
		}

		switch msg.Reach {
		case message.OnePeer:
			peerConn, exist := s.peers[msg.PeerID]
			if !exist {
				logMsg := fmt.Sprintf("Peer ID %s does not exist", msg.PeerID)
				log.Println(logMsg)
				responseMsg.Kind = message.TextMessage
				responseMsg.Content, err = json.Marshal(message.TextMessageContent{"error", logMsg})
				responseMsg.Reach = message.Self
				responseMsg.Sender = "server"
				if err != nil {
					log.Println("Error marshaling error message")
					continue
				}
				conn.WriteJSON(responseMsg)
				break
			}
			err = peerConn.WriteJSON(responseMsg)

			if err != nil {
				log.Printf("Failed to send message to peer %s: %v\n", msg.PeerID, err)
			}
		case message.AllPeers:
			for peerID, peerConn := range s.peers {
				if peerID != connID {
					err = peerConn.WriteJSON(responseMsg)
					if err != nil {
						log.Printf("Failed to send message to peer %s: %v\n", peerID, err)
					}
				}
			}
		case message.Self:
			err = conn.WriteJSON(responseMsg)
			if err != nil {
				log.Printf("Failed to send message to self %s: %v\n", connID, err)
			}
		case message.None:
			continue
		default:
			log.Printf("unexpected message reach type: %v", msg.Reach)
			conn.WriteMessage(websocket.TextMessage, []byte("Unexpected message reach type"))
			continue
		}
	}
}
