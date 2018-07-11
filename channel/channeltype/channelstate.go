package channeltype

//State is all possible state of channel
type State int

const (
	//StateInValid channel never exist
	StateInValid = iota
	//StateOpened channel is ready for transfer
	StateOpened
	//StateClosed 不能再发起交易了,还可以接受交易.
	StateClosed
	//StateSettled 通道已经彻底结算,和 invalid 状态意义相同
	StateSettled
	/*
		StateClosing 用户发起了关闭通道的请求,正在处理
		正在进行交易,可以继续,不再新开交易
	*/
	StateClosing
	/*
		StateSettling 用户发起了 结算请求,正在处理
		正常情况下此时不应该还有未完成交易,
		不能新开交易,正在进行的交易也没必要继续了.因为已经提交到链上了.
	*/
	StateSettling

	//StateWithdraw 用户收到或者发出了 withdraw 请求,这时候正在进行的交易只能立即放弃,因为没有任何意义了
	StateWithdraw
	//StateCooprativeSettle 用户收到或者发出了 cooperative settle 请求,这时候正在进行的交易只能立即放弃,因为没有任何意义了
	StateCooprativeSettle
	/*
		StatePrepareForSettling 收到了用户 cooperative 请求,但是有正在处理的交易,这时候不再接受新的交易了,可以等待一段时间,然后settle
		已开始交易,可以继续
	*/
	StatePrepareForSettling
	/*
		StatePrepareForWithdraw收到用户请求,要发起 withdraw, 但是目前还持有锁,不再发起或者接受任何交易,可以等待一段时间进行 withdraw
		已开始交易,可以继续
	*/
	StatePrepareForWithdraw
)

var TransferCannotBeContinuedMap map[State]bool
var CanTransferMap map[State]bool

func init() {
	TransferCannotBeContinuedMap = make(map[State]bool)
	CanTransferMap = make(map[State]bool)
	CanTransferMap[StateOpened] = true

	TransferCannotBeContinuedMap[StateSettling] = true
	TransferCannotBeContinuedMap[StateWithdraw] = true
	TransferCannotBeContinuedMap[StateCooprativeSettle] = true
}

func (s State) String() string {
	switch s {
	case StateInValid:
		return "inValid"
	case StateOpened:
		return "opened"
	case StateClosed:
		return "closed"
	case StateSettled:
		return "settled"
	case StateClosing:
		return "closing"
	case StateSettling:
		return "settling"
	case StateWithdraw:
		return "withdrawing"
	case StateCooprativeSettle:
		return "cooperativeSettling"
	default:
		return "unkown"
	}
}
