package rerr

import (
	"encoding/json"
	"fmt"
)

//StandardError 标准错误，包含错误码和错误信息
type StandardError struct {
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_message"`
}

// Error 实现了 Error接口
func (e StandardError) Error() string {
	return fmt.Sprintf("errorCode: %d, errorMsg %s", e.ErrorCode, e.ErrorMsg)
}

//Printf 附加错误描述信息
func (e StandardError) Printf(format string, a ...interface{}) StandardError {
	s := fmt.Sprintf(format, a...)
	err2 := e
	err2.ErrorMsg = fmt.Sprintf("%s:%s", e.ErrorMsg, s)
	return err2
}

//Errorf alias of printf
func (e StandardError) Errorf(format string, a ...interface{}) StandardError {
	return e.Printf(format, a...)
}

//Append 附加错误描述信息
func (e StandardError) Append(info string) StandardError {
	err2 := e
	err2.ErrorMsg = fmt.Sprintf("%s:%s", e.ErrorMsg, info)
	return err2
}

//AppendError 附加错误信息
func (e StandardError) AppendError(err error) StandardError {
	if err != nil {
		err2 := e
		err2.ErrorMsg = fmt.Sprintf("%s:%s", e.ErrorMsg, err.Error())
		return err2
	}
	return e
}

//StandardDataError 用于有结构化错误描述的场景
type StandardDataError struct {
	StandardError
	Data json.RawMessage `json:"data"`
}

//WithData 附加结构化错误信息
func (e StandardError) WithData(data interface{}) StandardDataError {
	err2 := StandardDataError{
		StandardError: e,
	}
	d, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	err2.Data = json.RawMessage(d)
	return err2
}

//ContractCallError 将合约调用错误转换为standarderror
func ContractCallError(err error) StandardError {
	//金额不够需要特殊处理
	if err.Error() == "insufficient balance to pay for gas" || err.Error() == "insufficient funds for gas * price + value" {
		return ErrInsufficientBalanceForGas

	}
	return ErrUnkownSpectrumRPCError.AppendError(err)
}

var m = make(map[int]struct{})

func newError(errCode int, errMsg string) StandardError {
	_, ok := m[errCode]
	if ok {
		panic(fmt.Sprintf("errCode %d already exist ", errCode))
	}
	m[errCode] = struct{}{}
	return StandardError{errCode, errMsg}
}

var (
	//ErrSuccess 成功
	ErrSuccess = newError(0, "success")
	//ErrUnrecognized 未知错误,
	ErrUnrecognized = newError(-1, "unknown error")
	//ErrArgumentError 参数错误
	ErrArgumentError = newError(1, "ArgumentError")
	//ErrPhotonAlreadyRunning 启动了多个photon实例
	ErrPhotonAlreadyRunning = newError(2, "PhotonAlreadyRunning")
	//ErrHashLengthNot32  参数错误
	ErrHashLengthNot32 = newError(1000, "HashLengthNot32")
	/*ErrNotFound Raised when something not found
	 */
	ErrNotFound = newError(1001, "Not found")
	/*ErrInsufficientBalance Raised when the netting channel doesn't enough available capacity to
	  pay for the transfer.

	  Used for the validation of an *incoming* messages.
	*/
	ErrInsufficientBalance = newError(1002, "InsufficientBalance")
	/*ErrInvalidAmount Raised when the user provided value is not a integer and cannot be used
	  to defined a transfer value
	*/
	ErrInvalidAmount = newError(1003, "InvalidAmount")

	/*ErrNoPathError Raised when there is no path to the requested target address in the
	  payment network.

	  This exception is raised if there is not a single path in the network to
	  reach the target, it's not used if there is a path but the transfre failed
	  because of the lack of capacity or network problems.
	*/
	ErrNoPathError = newError(1005, "NoPathError")
	/*ErrSamePeerAddress Raised when a user tries to create a channel where the address of both
	  peers is the same.
	*/
	ErrSamePeerAddress = newError(1006, "SamePeerAddress")
	/*ErrInvalidState Raised when the user requested action cannot be done due to the current
	  state of the channel.
	*/
	ErrInvalidState = newError(1007, "InvalidState")
	//ErrTransferWhenClosed Raised when a user tries to request a transfer is a closed channel.
	ErrTransferWhenClosed = newError(1008, "TransferWhenClosed")
	/*ErrUnknownAddress Raised when the user provided address is valid but is not from a known
	  node.
	*/
	ErrUnknownAddress = newError(1009, "UnknownAddress")
	/*ErrInvalidLocksRoot Raised when the received message has an invalid locksroot.

	  Used to reject a message when a pending lock is missing from the locksroot,
	  otherwise if the message is accepted there is a pontential loss of token.
	*/
	ErrInvalidLocksRoot = newError(1010, "Locksroot mismatch")
	/*ErrInvalidNonce Raised when the received messages has an invalid value for the nonce.

	  The nonce field must change incrementally
	*/
	ErrInvalidNonce = newError(1011, "InvalidNonce")

	/*ErrTransferUnwanted Raised when the node is not receiving new transfers.
	 */
	ErrTransferUnwanted = newError(1012, "TransferUnwanted")
	// ErrStopCreateNewTransfer reject new transactions
	ErrStopCreateNewTransfer = newError(1013, "new transactions are not allowed")

	//ErrNotAllowMediatedTransfer not allow mediated transfer when mesh
	ErrNotAllowMediatedTransfer = newError(1014, "can not send mediated transfer when photon works without effective chain")
	//ErrDuplicateTransfer token和secret相同的交易
	ErrDuplicateTransfer = newError(1015, "secret and token cannot duplicate")
	//ErrNodeNotOnline 发送消息时,对方不在线
	ErrNodeNotOnline = newError(1016, "NodeOffline")
	//ErrTransferCannotCancel 试图取消已经泄露秘密的交易
	ErrTransferCannotCancel = newError(1017, "TranasferCannotCancel")
	/*
		DB error
	*/

	//ErrGeneralDBError 未归类数据库错误,需要进一步细化
	ErrGeneralDBError = newError(1018, "DBError")
	//ErrDBDuplicateKey 重复的key
	ErrDBDuplicateKey = newError(1019, "duplicate key")
	//ErrTransferTimeout 交易超时,不代表交易肯定会成功或者失败,只是在给定时间内交易没有成功而已
	ErrTransferTimeout = newError(1020, "ErrTransferTimeout")
	//ErrUpdateButHaveTransfer 试图升级,发现还有交易在进行
	ErrUpdateButHaveTransfer = newError(1021, "ErrUpdateButHaveTransfer")
	//ErrNotChargeFee 进行与收费相关的操作,但是没有启用收费
	ErrNotChargeFee = newError(1022, "ErrNotChargeFee")
	/*
		以太坊报公链节点报的错误


	*/

	//ErrInsufficientBalanceForGas gas problem
	ErrInsufficientBalanceForGas = newError(2000, "insufficient balance to pay for gas")

	/*
		Tx 相关
		链上操作相关
	*/

	//ErrCloseChannel 链上执行关闭通道时发生了错误
	ErrCloseChannel = newError(2001, "closeChannel")
	//ErrRegisterSecret 链上注册密码的时候发生了错误
	ErrRegisterSecret = newError(2002, "RegisterSecret")
	//ErrUnlock 链上unlock的时候发生了错误
	ErrUnlock = newError(2003, "Unlock")
	//ErrUpdateBalanceProof 链上提交balance proof发生错误
	ErrUpdateBalanceProof = newError(2004, "UpdateBalanceProof")
	//ErrPunish 链上执行punish的时候发生错误
	ErrPunish = newError(2005, "punish")
	//ErrSettle 链上执行settle操作的时候发生错误
	ErrSettle = newError(2006, "settle")
	//ErrDeposit 链上执行deposit发生错误
	ErrDeposit = newError(2007, "deposit")
	//ErrSpectrumNotConnected 没有连接到公链.
	ErrSpectrumNotConnected = newError(2008, "ErrSpectrumNotConnected")
	//ErrTxWaitMined waitMined return error
	ErrTxWaitMined = newError(2009, "ErrTxWaitMined")
	//ErrTxReceiptStatus tx 被打包了,但是结果失败
	ErrTxReceiptStatus = newError(2010, "ErrTxReceiptStatus")
	//ErrSecretAlreadyRegistered 尝试连上注册密码,但是密码已经注册了
	ErrSecretAlreadyRegistered = newError(2011, "ErrSecretAlreadyRegistered")
	//ErrSpectrumSyncError 连接到的公链长时间不出块或者正在同步
	ErrSpectrumSyncError = newError(2012, "ErrSpectrumSyncError")
	//ErrSpectrumBlockError 本地已处理的块数和公链汇报块数不一致,比如我本地已经处理到了50000块,但是公链节点报告现在只有3000块
	ErrSpectrumBlockError = newError(2013, "ErrSpectrumBlockError")
	//ErrUnkownSpectrumRPCError 其他以太坊rpc错误
	ErrUnkownSpectrumRPCError = newError(2999, "unkown spectrum rpc error")
	/*ErrTokenNotFound Raised when token not found
	 */
	ErrTokenNotFound = newError(3001, "TokenNotFound")
	/*ErrChannelNotFound Raised when token not found
	 */
	ErrChannelNotFound = newError(3002, "ChannelNotFound")
	//ErrNoAvailabeRoute no availabe route
	ErrNoAvailabeRoute = newError(3003, "NoAvailabeRoute")
	//ErrTransferNotFound not found transfer
	ErrTransferNotFound = newError(3004, "TransferNotFound")
	//ErrChannelAlreadExist 通道已存在
	ErrChannelAlreadExist = newError(3005, "ChannelAlreadExist")
	//ErrRejectTransferBecauseChannelHoldingTooMuchLock 通道中已经存在过多的锁,暂时拒绝交易
	ErrRejectTransferBecauseChannelHoldingTooMuchLock = newError(3006, "channel too busy, reject mediated transfer for a while")
	//ErrRejectTransferBecausePayerChannelClosed 上家通道状态已经关闭,拒绝交易
	ErrRejectTransferBecausePayerChannelClosed = newError(3007, "payer's channel already closed ,reject mediated transfer")
	// ErrChannelNoEnoughBalance 通道余额不足
	ErrChannelNoEnoughBalance = newError(3008, "no enough balance")
	/*ErrPFS PFS Error
	向PFS发起请求错误
	*/
	ErrPFS = newError(4000, "ErrorPFS")

	/*
		Channel Error
	*/

	//ErrChannelNotAllowWithdraw 通道现在不能合作取现,比如有交易在进行
	ErrChannelNotAllowWithdraw = newError(5000, "CannotWithdarw")
	//ErrChannelState 在不能执行相应操作的通道状态,试图执行某些交易,比如在关闭的通道上发起交易
	ErrChannelState = newError(5001, "ErrChannelState")
	//ErrChannelSettleTimeout 没到settle时间尝试去settle
	ErrChannelSettleTimeout = newError(5002, "Channel only can settle after timeout")
	//ErrChannelNotParticipant 给定地址不是通道的任何参与一方
	ErrChannelNotParticipant = newError(5003, "NotParticipant")
	//ErrChannelLockSecretHashNotFound 通道中没有相应的锁
	ErrChannelLockSecretHashNotFound = newError(5004, "ChannelNoSuchLock")
	//ErrChannelEndStateNoSuchLock 通道当前参与方中找不到相应的锁
	ErrChannelEndStateNoSuchLock = newError(5005, "ErrChannelEndStateNoSuchLock")
	//ErrChannelLockAlreadyExpired 通道中锁已过期
	ErrChannelLockAlreadyExpired = newError(5006, "ErrChannelLockAlreadyExpired")
	//ErrChannelBalanceDecrease 发生了降低通道balance(指的是合约中的balance)的行为
	ErrChannelBalanceDecrease = newError(5007, "ErrChannelBalanceDecrease")
	//ErrChannelTransferAmountMismatch 收到的交易中transferamount不匹配
	ErrChannelTransferAmountMismatch = newError(5008, "ErrChannelTransferAmountMismatch")
	//ErrChannelBalanceProofAlreadyRegisteredOnChain  已经提交过balanceproof以后试图修改本地balance proof
	ErrChannelBalanceProofAlreadyRegisteredOnChain = newError(5009, "ErrChannelBalanceProofAlreadyRegisteredOnChain")
	//ErrChannelDuplicateLock 通道中已存在该密码的锁
	ErrChannelDuplicateLock = newError(5010, "ErrChannelDuplicateLock")
	//ErrChannelTransferAmountDecrease 收到交易,TransferAmount变小了
	ErrChannelTransferAmountDecrease = newError(5011, "ErrChannelTransferAmountDecrease")
	//ErrRemoveNotExpiredLock 试图移除没有过期的锁
	ErrRemoveNotExpiredLock = newError(5012, "ErrRemoveNotExpiredLock")
	//ErrUpdateBalanceProofAfterClosed 试图在通道关闭以后还更新对方或者我的balance proof,基本意思和ErrChannelBalanceProofAlreadyRegisteredOnChain一样
	ErrUpdateBalanceProofAfterClosed = newError(5013, "ErrUpdateBalanceProofAfterClosed")
	//ErrChannelIdentifierMismatch 通道id不匹配
	ErrChannelIdentifierMismatch = newError(5014, "ErrChannelIdentifierMismatch")
	//ErrChannelInvalidSender 收到来自未知参与方的交易
	ErrChannelInvalidSender = newError(5015, "ErrChannelInvalidSender")
	//ErrChannelBalanceNotMatch  合作关闭通道,取现时金额检查不匹配,
	ErrChannelBalanceNotMatch = newError(5016, "ErrChannelBalanceNotMatch")
	//ErrChannelLockMisMatch 收到交易中指定的锁和本地不匹配
	ErrChannelLockMisMatch = newError(5017, "ErrChannelLockMisMatch")
	//ErrChannelWithdrawAmount  合作取现的金额过大
	ErrChannelWithdrawAmount = newError(5018, "ErrChannelWithdrawAmount")
	//ErrChannelLockExpirationTooLarge 收到交易,指定的过期时间太长了,这可能是一种攻击
	ErrChannelLockExpirationTooLarge = newError(5019, "ErrChannelLockExpirationTooLarge")
	//ErrChannelRevealTimeout 指定的reveal timeout 非法
	ErrChannelRevealTimeout = newError(5020, "ErrChannelRevealTimeout")
	//ErrChannelBalanceProofNil balanceproof为空
	ErrChannelBalanceProofNil = newError(5021, "ErrChannelBalanceProofNil")
	//ErrChannelCloseClosedChannel 试图关闭已经关闭的通道
	ErrChannelCloseClosedChannel = newError(5022, "ErrChannelCloseClosedChannel")
	//ErrChannelBackgroundTx 后台执行Tx发生错误
	ErrChannelBackgroundTx = newError(5023, "ErrChannelBackgroundTx")

	/*ErrChannelWithdrawButHasLocks : we can't send a request for withdraw when there are locks.
	 */
	ErrChannelWithdrawButHasLocks = newError(5024, "ErrChannelWithdrawButHasLocks")
	/*ErrChannelCooperativeSettleButHasLocks : we can't send a request for settle when there are locks.
	 */
	ErrChannelCooperativeSettleButHasLocks = newError(5025, "ErrChannelCooperativeSettleButHasLocks")
	/*ErrChannelInvalidSttleTimeout Raised when the user provided timeout value is less than the minimum
	  settle timeout
	*/
	ErrChannelInvalidSttleTimeout = newError(5026, "ErrInvalidSettleTimeout")
	/*
		Transport error
	*/

	//ErrTransportTypeUnknown  未知的transport层错误,
	ErrTransportTypeUnknown = newError(6000, "transport type error")
	//ErrSubScribeNeighbor 订阅节点在线信息错误
	ErrSubScribeNeighbor = newError(6001, "ErrSubScribeNeighbor")

	// ErrUnknown 未知错误
	ErrUnknown = newError(9999, "unknown error")
)
