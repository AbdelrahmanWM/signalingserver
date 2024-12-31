package message

type ReachType int

const (
	Self ReachType = iota
	OnePeer
	AllPeers
	None
)

