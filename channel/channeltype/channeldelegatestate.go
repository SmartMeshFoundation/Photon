package channeltype

//ChannelDelegateState is all possible state of channel delegate status
type ChannelDelegateState int

const (
	//ChannelDelegateStateNoNeed only used when photon works without pms
	ChannelDelegateStateNoNeed = iota
	//ChannelDelegateStateWaiting 正在等待委托到pms
	ChannelDelegateStateWaiting
	// ChannelDelegateStateSuccess 委托成功
	ChannelDelegateStateSuccess
	// ChannelDelegateStateFail 委托失败
	ChannelDelegateStateFail
	// ChannelDelegateStateFailAndNoEffectiveChain 委托失败且无有效公链
	ChannelDelegateStateFailAndNoEffectiveChain
)

// String  for read
func (s ChannelDelegateState) String() string {
	switch s {
	case ChannelDelegateStateNoNeed:
		return "no need"
	case ChannelDelegateStateWaiting:
		return "waiting to delegate to pms"
	case ChannelDelegateStateSuccess:
		return "delegate success"
	case ChannelDelegateStateFail:
		return "delegate fail"
	case ChannelDelegateStateFailAndNoEffectiveChain:
		return "delegate fail since no effective chain,dangerous"
	default:
		return "unknown"
	}
}
