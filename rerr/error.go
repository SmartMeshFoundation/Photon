package rerr

import (
	"errors"

	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func newrerr(msg1, msg2 string) error {
	return fmt.Errorf("%s: %s", msg1, msg2)
}

/*
ErrHashLengthNot32 Raised if the length of the provided element is not 32 bytes in length,
    a keccak hash is required to include the element in the merkle tree.
*/
var ErrHashLengthNot32 = errors.New("HashLengthNot32")

/*
ChannelNotFound Raised when a provided channel via the REST api is not found in the
    internal data structures
*/
func ChannelNotFound(msg string) error {
	return fmt.Errorf("ChannelNotFound %s", msg)
}

/*
ErrInsufficientFunds Raised when provided account doesn't have token funds to complete the
    requested deposit.

    Used when a *user* tries to deposit a given amount of token in a channel,
    but his account doesn't have enough funds to pay for the deposit
*/
var ErrInsufficientFunds = errors.New("InsufficientFunds")

//InvalidAddress Raised when the user provided value is not a valid address.
func InvalidAddress(msg string) error {
	return newrerr("InvalidAddress", msg)
}

/*
ErrInvalidAmount Raised when the user provided value is not a integer and cannot be used
    to defined a transfer value
*/
var ErrInvalidAmount = errors.New("InvalidAmount")

/*
ErrInvalidSettleTimeout Raised when the user provided timeout value is less than the minimum
    settle timeout
*/
var ErrInvalidSettleTimeout = errors.New("ErrInvalidSettleTimeout")

/*
ErrNoPathError Raised when there is no path to the requested target address in the
    payment network.

    This exception is raised if there is not a single path in the network to
    reach the target, it's not used if there is a path but the transfre failed
    because of the lack of capacity or network problems.
*/
var ErrNoPathError = errors.New("NoPathError")

/*
ErrSamePeerAddress Raised when a user tries to create a channel where the address of both
    peers is the same.
*/
var ErrSamePeerAddress = errors.New("SamePeerAddress")

/*
InvalidState Raised when the user requested action cannot be done due to the current
    state of the channel.
*/
func InvalidState(msg string) error {
	return newrerr("InvalidState", msg)
}

//TransferWhenClosed Raised when a user tries to request a transfer is a closed channel.
func TransferWhenClosed(msg string) error {
	return newrerr("TransferWhenClosed", msg)
}

/*
UnknownAddress Raised when the user provided address is valid but is not from a known
    node.
*/
func UnknownAddress(msg string) error {
	return fmt.Errorf("UnknownAddress: %s", msg)
}

/*
ErrInsufficientBalance Raised when the netting channel doesn't enough available capacity to
    pay for the transfer.

    Used for the validation of an *incoming* messages.
*/
var ErrInsufficientBalance = errors.New("InsufficientBalance")

/*
InvalidLocksRoot Raised when the received message has an invalid locksroot.

    Used to reject a message when a pending lock is missing from the locksroot,
    otherwise if the message is accepted there is a pontential loss of token.
*/
func InvalidLocksRoot(expectedLocksroot, gotLocksroot common.Hash) error {
	return fmt.Errorf("Locksroot mismatch. Expected %s but got %s", utils.HPex(expectedLocksroot), utils.HPex(gotLocksroot))
}

/*
InvalidNonce Raised when the received messages has an invalid value for the nonce.

    The nonce field must change incrementally
*/
func InvalidNonce(msg string) error {
	return newrerr("InvalidNonce", msg)
}

/*
ErrTransferUnwanted Raised when the node is not receiving new transfers.
*/
var ErrTransferUnwanted = errors.New("TransferUnwanted")

//UnknownTokenAddress token address is unkown
func UnknownTokenAddress(msg string) error {
	return fmt.Errorf("UnknownTokenAddress: %s", msg)
}

//ErrSTUNUnavailableException cannot reach stun server
var ErrSTUNUnavailableException = errors.New("STUNUnavailableException")

//ErrEthNodeCommunicationError eth communication error
var ErrEthNodeCommunicationError = errors.New("EthNodeCommunicationError")

/*
ErrAddressWithoutCode Raised on attempt to execute contract on address without a code.
*/
var ErrAddressWithoutCode = errors.New("AddressWithoutCode")

/*
ErrNoTokenManager Manager for a given token does not exist.
*/
var ErrNoTokenManager = errors.New("NoTokenManager")

/*
ErrDuplicatedChannelError Raised if someone tries to create a channel that already exists
*/
var ErrDuplicatedChannelError = errors.New("DuplicatedChannelError")

/*
TransactionThrew Raised when, after waiting for a transaction to be mined,
    the receipt has a 0x0 status field
*/
func TransactionThrew(txName string, receipt *types.Receipt) error {
	return fmt.Errorf("%s transaction threw. Receipt=%s", txName, receipt)
}

//ErrTransferTimeout  timeout error
var ErrTransferTimeout = errors.New("TransferTimeout")
