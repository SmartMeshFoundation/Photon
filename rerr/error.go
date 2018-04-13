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
 Raised if the length of the provided element is not 32 bytes in length,
    a keccak hash is required to include the element in the merkle tree.
*/
var HashLengthNot32 = errors.New("ErrHashLengthNot32")

/*
 Raised when a provided channel via the REST api is not found in the
    internal data structures
*/
func ChannelNotFound(msg string) error {
	return fmt.Errorf("ChannelNotFound %s", msg)
}

/*
Raised when provided account doesn't have token funds to complete the
    requested deposit.

    Used when a *user* tries to deposit a given amount of token in a channel,
    but his account doesn't have enough funds to pay for the deposit
*/
var InsufficientFunds = errors.New("InsufficientFunds")

//Raised when the user provided value is not a valid address.
func InvalidAddress(msg string) error {
	return newrerr("InvalidAddress", msg)
}

/*
Raised when the user provided value is not a integer and cannot be used
    to defined a transfer value
*/
var InvalidAmount = errors.New("InvalidAmount")

/*
Raised when the user provided timeout value is less than the minimum
    settle timeout
*/
var InvalidSettleTimeout = errors.New("InvalidSettleTimeout")

/*
Raised when there is no path to the requested target address in the
    payment network.

    This exception is raised if there is not a single path in the network to
    reach the target, it's not used if there is a path but the transfre failed
    because of the lack of capacity or network problems.
*/
var NoPathError = errors.New("NoPathError")

/*
Raised when a user tries to create a channel where the address of both
    peers is the same.
*/
var SamePeerAddress = errors.New("SamePeerAddress")

/*
Raised when the user requested action cannot be done due to the current
    state of the channel.
*/
func InvalidState(msg string) error {
	return newrerr("InvalidState", msg)
}

//Raised when a user tries to request a transfer is a closed channel.
func TransferWhenClosed(msg string) error {
	return newrerr("TransferWhenClosed", msg)
}

/*
Raised when the user provided address is valid but is not from a known
    node.
*/
//var ErrUnknownAddress = errors.New("UnknownAddress")
func UnknownAddress(msg string) error {
	return fmt.Errorf("UnknownAddress: ", msg)
}

/*
Raised when the netting channel doesn't enough available capacity to
    pay for the transfer.

    Used for the validation of an *incoming* messages.
*/
var InsufficientBalance = errors.New("InsufficientBalance")

/*
Raised when the received message has an invalid locksroot.

    Used to reject a message when a pending lock is missing from the locksroot,
    otherwise if the message is accepted there is a pontential loss of token.
*/
func InvalidLocksRoot(expectedLocksroot, gotLocksroot common.Hash) error {
	return fmt.Errorf("Locksroot mismatch. Expected %s but got %s", utils.HPex(expectedLocksroot), utils.HPex(gotLocksroot))
}

/*
Raised when the received messages has an invalid value for the nonce.

    The nonce field must change incrementally
*/
func InvalidNonce(msg string) error {
	return newrerr("InvalidNonce", msg)
}

/*
Raised when the node is not receiving new transfers.
*/
var TransferUnwanted = errors.New("TransferUnwanted")

func UnknownTokenAddress(msg string) error {
	return fmt.Errorf("UnknownTokenAddress:", msg)
}

var STUNUnavailableException = errors.New("STUNUnavailableException")
var EthNodeCommunicationError = errors.New("EthNodeCommunicationError")

/*
Raised on attempt to execute contract on address without a code.
*/
var AddressWithoutCode = errors.New("AddressWithoutCode")

/*
Manager for a given token does not exist.
*/
var NoTokenManager = errors.New("NoTokenManager")

/*
Raised if someone tries to create a channel that already exists
*/
var DuplicatedChannelError = errors.New("DuplicatedChannelError")

/*
Raised when, after waiting for a transaction to be mined,
    the receipt has a 0x0 status field
*/
func TransactionThrew(txName string, receipt *types.Receipt) error {
	return fmt.Errorf("%s transaction threw. Receipt=%s", txName, receipt)
}

var TransferTimeout = errors.New("TransferTimeout")

func Timeout(msg string) error {
	return newrerr("Timeout", msg)
}

func GoChannelClosed(msg string) error {
	return newrerr("GoChannelClosed", msg)
}
