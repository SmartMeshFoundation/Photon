package transfer

import (
	"encoding/gob"

	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

/*
EventTransferSentSuccess emitted by the initiator when a transfer is considered sucessful.

    A transfer is considered sucessful when the initiator's payee hop sends the
    reveal secret message, assuming that each hop in the mediator chain has
    also learned the secret and unlock/withdraw its token.

    This definition of sucessful is used to avoid the following corner case:

    - The reveal secret message is sent, since the network is unreliable and we
      assume byzantine behavior the message is considered delivered without an
      acknowledgement.
    - The transfer is considered sucessful because of the above.
    - The reveal secret message was not delivered because of actual network
      problems.
    - The lock expires and an EventUnlockFailed follows, contradicting the
      EventTransferSentSuccess.

    Note:
        Mediators cannot use this event, since an unlock may be locally
        sucessful but there is no knowledge about the global transfer.
*/
type EventTransferSentSuccess struct {
	Identifier uint64
	Amount     *big.Int
	Target     common.Address
}

/*
EventTransferSentFailed emitted by the payer when a transfer has failed.

    Note:
        Mediators cannot use this event since they don't know when a transfer
        has failed, they may infer about lock successes and failures.
*/
type EventTransferSentFailed struct {
	Identifier uint64
	Reason     string
	Target     common.Address //transfer's target, may be not the same as receipient
}

/*
EventTransferReceivedSuccess emitted when a payee has received a payment.

    Note:
        A payee knows if a lock withdraw has failed, but this is not sufficient
        information to deduce when a transfer has failed, because the initiator may
        try again at a different time and/or with different routes, for this reason
        there is no correspoding `EventTransferReceivedFailed`.
*/
type EventTransferReceivedSuccess struct {
	Identifier uint64
	Amount     *big.Int
	Initiator  common.Address
}

func init() {
	gob.Register(&EventTransferSentSuccess{})
	gob.Register(&EventTransferSentFailed{})
	gob.Register(&EventTransferReceivedSuccess{})
}
