package SignalInterface

import "github.com/ethereum/go-ethereum/common"

// SignalProxy represents connection to signal server.
type SignalProxy interface {
	// TryReach test if addr is online or not
	TryReach(addr common.Address) error
	// ExchangeSdp exchange Sdp information with another raiden node
	ExchangeSdp(addr common.Address, sdp string) (partnerSdp string, err error)
	// Connected allows to check that client connected to Centrifugo at moment.
	Connected() bool
	// Close closes connection to Centrifugo and clears all subscriptions.
	Close()
}

//SdpHandler handels a ICE connection request from a raiden node
type SdpHandler func(from common.Address, sdp string) (mysdp string, err error)
