package netshare

// Status shows actual connection status.
type Status int

const (
	//Disconnected init status
	Disconnected = Status(iota)
	//Connected connection status
	Connected
	//Closed user closed
	Closed
	//Reconnecting connection error
	Reconnecting
)
