package mediatedtransfer

import (
	"encoding/gob"

	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// EventSendMediatedTransfer A mediated transfer that must be sent to `node_address`.
type EventSendMediatedTransfer struct {
	Token          common.Address
	Amount         *big.Int
	LockSecretHash common.Hash
	Initiator      common.Address
	Target         common.Address
	Expiration     int64
	Receiver       common.Address
	Fee            *big.Int // target should get amount-fee.
	/*
		which channel received a mediated transfer and then I have to send another mediated transfer,
		因为哪个 channel 收到了 MediatedTransfer, 导致我需要发送新的 Transfer.
		如果是我主动发起的 MediatedTransfer, 那么 FromChannel 应该为空
	*/
	// no matter which channel received a mediated transfer, I have to send another mediated transfer,
	// because which channel receives MediatedTransfer and leads me to send a new Transfer
	// If I am the transfer initiator, then FromChannel should be null.
	FromChannel common.Hash
	Path        []common.Address //2019-03 消息升级后,带全路径path
}

//NewEventSendMediatedTransfer create EventSendMediatedTransfer
func NewEventSendMediatedTransfer(transfer *LockedTransferState, receiver common.Address, path []common.Address) *EventSendMediatedTransfer {
	return &EventSendMediatedTransfer{
		Token:          transfer.Token,
		Amount:         new(big.Int).Set(transfer.Amount),
		LockSecretHash: transfer.LockSecretHash,
		Initiator:      transfer.Initiator,
		Target:         transfer.Target,
		Expiration:     transfer.Expiration,
		Receiver:       receiver,
		Fee:            transfer.Fee,
		Path:           path,
	}
}

/*
EventSendRevealSecret Sends a RevealSecret to another node.

    This event is used once the secret is known locally and an action must be
    performed on the receiver:

        - For receivers in the payee role, it informs the node that the lock has
            been released and the token can be withdrawn, either on-chain or
            off-chain.
        - For receivers in the payer role, it tells the payer that the payee
            knows the secret and wants to withdraw the lock off-chain, so the payer
            may unlock the lock and send an up-to-date balance proof to the payee,
            avoiding on-chain payments which would require the channel to be
            closed.

    For any mediated transfer:

        - The initiator will only perform the payer role.
        - The target will only perform the payee role.
        - The mediators will have `n` channels at the payee role and `n` at the
          payer role, where `n` is equal to `1 + number_of_refunds`.

    Note:
        The payee must only update its local balance once the payer sends an
        up-to-date balance-proof message. This is a requirement for keeping the
        nodes synchronized. The reveal secret message flows from the receiver
        to the sender, so when the secret is learned it is not yet time to
        update the balance.
*/
type EventSendRevealSecret struct {
	LockSecretHash common.Hash
	Secret         common.Hash
	Token          common.Address
	Receiver       common.Address
	Sender         common.Address
	Data           string
}

/*
EventSendBalanceProof send a balance-proof to the counter-party, used after a lock
    is unlocked locally allowing the counter-party to withdraw.

    Used by payers: The initiator and mediator nodes.

    Note:
        This event has a dual role, it serves as a synchronization and as
        balance-proof for the netting channel smart contract.

        Nodes need to keep the last known merkle root synchronized. This is
        required by the receiving end of a transfer in order to properly
        validate. The rule is "only the party that owns the current payment
        channel may change it" (remember that a netting channel is composed of
        two uni-directional channels), as a consequence the merkle root is only
        updated by the receiver once a balance proof message is received.
*/
type EventSendBalanceProof struct {
	LockSecretHash    common.Hash
	ChannelIdentifier common.Hash
	Token             common.Address
	Receiver          common.Address
}

/*
EventSendSecretRequest used by a target node to request the secret from the initiator
    (`receiver`).
*/
type EventSendSecretRequest struct {
	ChannelIdentifier common.Hash
	LockSecretHash    common.Hash
	Amount            *big.Int
	Receiver          common.Address
}

/*
EventSendAnnounceDisposed used to cleanly backtrack the current node in the route.

    This message will pay back the same amount of token from the receiver to
    the sender, allowing the sender to try a different route without the risk
    of losing token.
*/
type EventSendAnnounceDisposed struct {
	Amount         *big.Int
	LockSecretHash common.Hash
	Expiration     int64
	Token          common.Address
	Receiver       common.Address
}

/*
EventSendAnnounceDisposedResponse 收到对方AnnounceDisposed,需要给以应答
这时候我可能会一次发出两条消息,
一条是 Reponse, 另一条是 MediatedTransfer.
我极可能是中间节点,也可能是交易发起人,但是不会是接收方.
*/
/*
 *	EventSendAnnounceDisposedResponse : after received AnnounceDisposed message from his partner,
 * 	he needs to respond and there are two messages :
 *	1. Response
 *	2. MediatedTransfer
 *	This participant has strong possibility to be a mediated node, or transfer initiator, but never be recipient.
 */
type EventSendAnnounceDisposedResponse struct {
	LockSecretHash common.Hash
	Token          common.Address
	Receiver       common.Address
}

/*
EventContractSendRegisterSecret emitted to register a secret on chain

    This event is used when a node needs to prepare the channel to withdraw
    on-chain.
*/
type EventContractSendRegisterSecret struct {
	Secret common.Hash
}

/*
EventContractSendUnlock emitted when the lock must be withdrawn on-chain.
上家不知什么原因要关闭channel，我一旦知道密码，应该立即到链上提现。
channel 自己会关注是否要提现，但是如果是在关闭以后才获取到密码的呢？
目前完全无用,如果 unlock 放在 settle 之后,还有可能有用.
*/
/*
EventContractSendUnlock : emmited when the lock must be unlock on-chain.
 当一个密码在链上注册成功以后,这时候自己(接收方)必须判断通道是否关闭,因为如果已经关闭,
说明自己已经提交过updateBalanceProof,并且在合约上调用过相关的unlock,当然肯定unlock失败了.

比如A-B-C交易,A可能在交易未完成的时候就关闭通道,当B拿到了C的密码以后,B实际上已经在链上提交过BalanceProof,并且调用过相关的unlock.
这时候B的unlock肯定是失败的,因为还没有注册密码.
*/
type EventContractSendUnlock struct {
	LockSecretHash    common.Hash
	ChannelIdentifier common.Hash
}

//EventUnlockSuccess emitted when a lock unlock succeded ,emit this event after receive a revealsecret message
type EventUnlockSuccess struct {
	LockSecretHash common.Hash
}

/*
下家没有在expiration之内收到balanceproof，也没有选择在链上兑现。
能想到的应对就是移除失效的lock
*/

//EventUnlockFailed emitted when a lock unlock failed.
type EventUnlockFailed struct {
	LockSecretHash    common.Hash
	ChannelIdentifier common.Hash
	Reason            string
}

//EventWithdrawSuccess emitted when a lock withdraw succeded.
type EventWithdrawSuccess struct {
	LockSecretHash common.Hash
}

/*
EventWithdrawFailed : 上家没有在expiration之内给我balanceproof，我也没有在链上兑现（因为没有密码）。
必须等待上家的 RemoveExpiredHashlockTransfer, 然后移除.
*/

// EventWithdrawFailed : cases that previous node does not transfer BalanceProof to him within expiration block, and
// 	this participant also does not register secret on-chain (because he does not have secret).
// 	So this participant must wait for RemoveExpiredHashlock from his previous node, then remove this lock.
type EventWithdrawFailed struct {
	LockSecretHash    common.Hash
	ChannelIdentifier common.Hash
	Reason            string
}

// EventSaveFeeChargeRecord :
// 记录本次中转收取手续费的流水
type EventSaveFeeChargeRecord struct {
	LockSecretHash common.Hash    `json:"lock_secret_hash"`
	TokenAddress   common.Address `json:"token_address"`
	TransferFrom   common.Address `json:"transfer_from"`
	TransferTo     common.Address `json:"transfer_to"`
	TransferAmount *big.Int       `json:"transfer_amount"`
	InChannel      common.Hash    `json:"in_channel"`  // 我收款的channelID
	OutChannel     common.Hash    `json:"out_channel"` // 我付款的channelID
	Fee            *big.Int       `json:"fee"`
	Timestamp      int64          `json:"timestamp"` // 时间戳,time.Unix()
}

func init() {
	gob.Register(&EventSendMediatedTransfer{})
	gob.Register(&EventSendRevealSecret{})
	gob.Register(&EventSendBalanceProof{})
	gob.Register(&EventSendSecretRequest{})
	gob.Register(&EventSendAnnounceDisposed{})
	gob.Register(&EventContractSendRegisterSecret{})
	gob.Register(&EventContractSendUnlock{})
	gob.Register(&EventUnlockSuccess{})
	gob.Register(&EventUnlockFailed{})
	gob.Register(&EventWithdrawSuccess{})
	gob.Register(&EventWithdrawFailed{})
}
