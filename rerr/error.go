package rerr

import (
	"github.com/SmartMeshFoundation/Photon/channel/channeltype"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

/*
InvalidState Raised when the user requested action cannot be done due to the current
    state of the channel.
*/
func InvalidState(msg string) error {
	return ErrInvalidState.Append(msg)
}

//TransferWhenClosed Raised when a user tries to request a transfer is a closed channel.
func TransferWhenClosed(msg string) error {
	return ErrTransferWhenClosed.Append(msg)
}

/*
UnknownAddress Raised when the user provided address is valid but is not from a known
    node.
*/
func UnknownAddress(msg string) error {
	return ErrUnknownAddress.Append(msg)
}

/*
InvalidLocksRoot Raised when the received message has an invalid locksroot.

    Used to reject a message when a pending lock is missing from the locksroot,
    otherwise if the message is accepted there is a pontential loss of token.
*/
func InvalidLocksRoot(expectedLocksroot, gotLocksroot common.Hash) error {
	return ErrInvalidLocksRoot.Printf("Expected %s but got %s", utils.HPex(expectedLocksroot), utils.HPex(gotLocksroot))
}

/*
InvalidNonce Raised when the received messages has an invalid value for the nonce.

    The nonce field must change incrementally
*/
func InvalidNonce(msg string) StandardError {
	return ErrInvalidNonce.Append(msg)
}

//ChannelStateError  在不能执行相应操作的通道状态,试图执行某些交易,比如在关闭的通道上发起交易
func ChannelStateError(state channeltype.State) StandardError {
	return ErrChannelState.Printf("state=%s", state)
}

//ChannelNotFound 找不到通道错误
func ChannelNotFound(info string) StandardError {
	return ErrChannelNotFound.Append(info)
}
