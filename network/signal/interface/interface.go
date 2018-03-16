package SignalInterface

import "github.com/ethereum/go-ethereum/common"

// SignalProxy represents connection to signal server.
type SignalProxy interface {
	// Subscribe allows to subscribe on channel and react on various subscription events.
	TryReach(addr common.Address) error
	// ClientID returns client ID that Centrifugo gave to connection. Or empty string if no client ID issued yet.
	ExchangeSdp(addr common.Address, sdp string) (partnerSdp string, err error)
	// Connected allows to check that client connected to Centrifugo at moment.
	Connected() bool
	// Close closes connection to Centrifugo and clears all subscriptions.
	Close()
}

type SdpHandler func(from common.Address, sdp string) (mysdp string, err error)
