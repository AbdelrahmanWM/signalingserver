package message

import (
	"encoding/json"
	"fmt"
)

type ReachType int

const (
	Self ReachType = iota
	OnePeer
	AllPeers
	None
)

func (r ReachType) MarshalJSON() ([]byte, error) {
	switch r {
	case Self:
		return json.Marshal("Self")
	case OnePeer:
		return json.Marshal("OnePeer")
	case AllPeers:
		return json.Marshal("AllPeers")
	case None:
		return json.Marshal("None")
	default:
		return nil, fmt.Errorf("unknown ReachType %d", r)
	}
}
func (r *ReachType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch s {
	case "Self":
		*r = Self
	case "OnePeer":
		*r = OnePeer
	case "AllPeers":
		*r = AllPeers
	case "None":
		*r = None
	default:
		return fmt.Errorf("unknown ReachType string %s", s)
	}
	return nil
}
