package ice

import "errors"

var (
	errCheckRetry         = errors.New("retry again")
	errTriedTooManyTimes  = errors.New("have tried too many times")
	errInvalidStunMessage = errors.New("Invalid STUN message")
	errStunInvalidLength  = errors.New("Invalid STUN message length")
	errStunUnknownType    = errors.New("Invalid or unexpected STUN message type")
	errStunTimeout        = errors.New("STUN transaction has timed out")

	errStunTooManyAttributes = errors.New("Too many STUN attributes")
	errStunAttributeLength   = errors.New("Invalid STUN attribute length")
)
