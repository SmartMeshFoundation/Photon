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

//NewError create an error and check duplicates
func NewError(errCode int, errMsg string) StandardError {
	_, ok := m[errCode]
	if ok {
		panic(fmt.Sprintf("errCode %d already exist ", errCode))
	}
	m[errCode] = struct{}{}
	return StandardError{errCode, errMsg}
}

var (
	//ErrSuccess 成功
	ErrSuccess = NewError(0, "success")
	//ErrUnrecognized 未知错误,
	ErrUnrecognized = NewError(-1, "unknown error")
	//ErrArgumentError 参数错误
	ErrArgumentError = NewError(1, "ArgumentError")
	//ErrPhotonAlreadyRunning 启动了多个photon实例
	ErrPhotonAlreadyRunning = NewError(2, "PhotonAlreadyRunning")
	//ErrHashLengthNot32  参数错误
	ErrHashLengthNot32 = NewError(1000, "HashLengthNot32")
	/*ErrNotFound Raised when something not found
	 */
	ErrNotFound = NewError(1001, "Not found")
	/*ErrInsufficientBalance Raised when the netting channel doesn't enough available capacity to
	  pay for the transfer.

	  Used for the validation of an *incoming* messages.
	*/
	ErrInsufficientBalance = NewError(1002, "InsufficientBalance")
	/*ErrInvalidAmount Raised when the user provided value is not a integer and cannot be used
	  to defined a transfer value
	*/
	ErrInvalidAmount = NewError(1003, "InvalidAmount")

	/*ErrNoPathError Raised when there is no path to the requested target address in the
	  payment network.

	  This exception is raised if there is not a single path in the network to
	  reach the target, it's not used if there is a path but the transfre failed
	  because of the lack of capacity or network problems.
	*/
	ErrNoPathError = NewError(1005, "NoPathError")
	/*ErrSamePeerAddress Raised when a user tries to create a channel where the address of both
	  peers is the same.
	*/
	ErrSamePeerAddress = NewError(1006, "SamePeerAddress")
	/*ErrInvalidState Raised when the user requested action cannot be done due to the current
	  state of the channel.
	*/
	ErrInvalidState = NewError(1007, "InvalidState")
	//ErrTransferWhenClosed Raised when a user tries to request a transfer is a closed channel.
	ErrTransferWhenClosed = NewError(1008, "TransferWhenClosed")
	/*ErrUnknownAddress Raised when the user provided address is valid but is not from a known
	  node.
	*/
	ErrUnknownAddress = NewError(1009, "UnknownAddress")
	/*ErrInvalidLocksRoot Raised when the received message has an invalid locksroot.

	  Used to reject a message when a pending lock is missing from the locksroot,
	  otherwise if the message is accepted there is a pontential loss of token.
	*/
	ErrInvalidLocksRoot = NewError(1010, "Locksroot mismatch")
	/*ErrInvalidNonce Raised when the received messages has an invalid value for the nonce.

	  The nonce field must change incrementally
	*/
	ErrInvalidNonce = NewError(1011, "InvalidNonce")

	/*ErrTransferUnwanted Raised when the node is not receiving new transfers.
	 */
	ErrTransferUnwanted = NewError(1012, "TransferUnwanted")
	// ErrStopCreateNewTransfer reject new transactions
	ErrStopCreateNewTransfer = NewError(1013, "new transactions are not allowed")

	//ErrNotAllowMediatedTransfer not allow mediated transfer when mesh
	ErrNotAllowMediatedTransfer = NewError(1014, "can not send mediated transfer when photon works without effective chain")
	//ErrDuplicateTransfer token和secret相同的交易
	ErrDuplicateTransfer = NewError(1015, "secret and token cannot duplicate")
	//ErrNodeNotOnline 发送消息时,对方不在线
	ErrNodeNotOnline = NewError(1016, "NodeOffline")
	//ErrTransferCannotCancel 试图取消已经泄露秘密的交易
	ErrTransferCannotCancel = NewError(1017, "TranasferCannotCancel")
	/*
		DB error
	*/

	//ErrGeneralDBError 未归类数据库错误,需要进一步细化
	ErrGeneralDBError = NewError(1018, "DBError")
	//ErrDBDuplicateKey 重复的key
	ErrDBDuplicateKey = NewError(1019, "duplicate key")
	//ErrTransferTimeout 交易超时,不代表交易肯定会成功或者失败,只是在给定时间内交易没有成功而已
	ErrTransferTimeout = NewError(1020, "ErrTransferTimeout")
	//ErrUpdateButHaveTransfer 试图升级,发现还有交易在进行
	ErrUpdateButHaveTransfer = NewError(1021, "ErrUpdateButHaveTransfer")
	//ErrNotChargeFee 进行与收费相关的操作,但是没有启用收费
	ErrNotChargeFee = NewError(1022, "ErrNotChargeFee")
	//ErrNotAllowDirectTransfer not allow mediated transfer when mesh
	ErrNotAllowDirectTransfer = NewError(1023, "can not send direct transfer after photon worked without effective chain for a long time")
	/*
		以太坊报公链节点报的错误


	*/

	//ErrInsufficientBalanceForGas gas problem
	ErrInsufficientBalanceForGas = NewError(2000, "insufficient balance to pay for gas")

	/*
		Tx 相关
		链上操作相关
	*/

	//ErrCloseChannel 链上执行关闭通道时发生了错误
	ErrCloseChannel = NewError(2001, "closeChannel")
	//ErrRegisterSecret 链上注册密码的时候发生了错误
	ErrRegisterSecret = NewError(2002, "RegisterSecret")
	//ErrUnlock 链上unlock的时候发生了错误
	ErrUnlock = NewError(2003, "Unlock")
	//ErrUpdateBalanceProof 链上提交balance proof发生错误
	ErrUpdateBalanceProof = NewError(2004, "UpdateBalanceProof")
	//ErrPunish 链上执行punish的时候发生错误
	ErrPunish = NewError(2005, "punish")
	//ErrSettle 链上执行settle操作的时候发生错误
	ErrSettle = NewError(2006, "settle")
	//ErrDeposit 链上执行deposit发生错误
	ErrDeposit = NewError(2007, "deposit")
	//ErrSpectrumNotConnected 没有连接到公链.
	ErrSpectrumNotConnected = NewError(2008, "ErrSpectrumNotConnected")
	//ErrTxWaitMined waitMined return error
	ErrTxWaitMined = NewError(2009, "ErrTxWaitMined")
	//ErrTxReceiptStatus tx 被打包了,但是结果失败
	ErrTxReceiptStatus = NewError(2010, "ErrTxReceiptStatus")
	//ErrSecretAlreadyRegistered 尝试连上注册密码,但是密码已经注册了
	ErrSecretAlreadyRegistered = NewError(2011, "ErrSecretAlreadyRegistered")
	//ErrSpectrumSyncError 连接到的公链长时间不出块或者正在同步
	ErrSpectrumSyncError = NewError(2012, "ErrSpectrumSyncError")
	//ErrSpectrumBlockError 本地已处理的块数和公链汇报块数不一致,比如我本地已经处理到了50000块,但是公链节点报告现在只有3000块
	ErrSpectrumBlockError = NewError(2013, "ErrSpectrumBlockError")
	//ErrUnkownSpectrumRPCError 其他以太坊rpc错误
	ErrUnkownSpectrumRPCError = NewError(2999, "unkown spectrum rpc error")
	/*ErrTokenNotFound Raised when token not found
	 */
	ErrTokenNotFound = NewError(3001, "TokenNotFound")
	/*ErrChannelNotFound Raised when token not found
	 */
	ErrChannelNotFound = NewError(3002, "ChannelNotFound")
	//ErrNoAvailabeRoute no availabe route
	ErrNoAvailabeRoute = NewError(3003, "NoAvailabeRoute")
	//ErrTransferNotFound not found transfer
	ErrTransferNotFound = NewError(3004, "TransferNotFound")
	//ErrChannelAlreadExist 通道已存在
	ErrChannelAlreadExist = NewError(3005, "ChannelAlreadExist")
	//ErrRejectTransferBecauseChannelHoldingTooMuchLock 通道中已经存在过多的锁,暂时拒绝交易
	ErrRejectTransferBecauseChannelHoldingTooMuchLock = NewError(3006, "channel too busy, reject mediated transfer for a while")
	//ErrRejectTransferBecausePayerChannelClosed 上家通道状态已经关闭,拒绝交易
	ErrRejectTransferBecausePayerChannelClosed = NewError(3007, "payer's channel already closed ,reject mediated transfer")
	// ErrChannelNoEnoughBalance 通道余额不足
	ErrChannelNoEnoughBalance = NewError(3008, "no enough balance")
	/*ErrPFS PFS Error
	向PFS发起请求错误
	*/
	ErrPFS = NewError(4000, "ErrorPFS")

	/*
		Channel Error
	*/

	//ErrChannelNotAllowWithdraw 通道现在不能合作取现,比如通道已经关闭或者正在withdraw等
	ErrChannelNotAllowWithdraw = NewError(5000, "CannotWithdarw")
	//ErrChannelState 在不能执行相应操作的通道状态,试图执行某些交易,比如在关闭的通道上发起交易
	ErrChannelState = NewError(5001, "ErrChannelState")
	//ErrChannelSettleTimeout 没到settle时间尝试去settle
	ErrChannelSettleTimeout = NewError(5002, "Channel only can settle after timeout")
	//ErrChannelNotParticipant 给定地址不是通道的任何参与一方
	ErrChannelNotParticipant = NewError(5003, "NotParticipant")
	//ErrChannelLockSecretHashNotFound 通道中没有相应的锁
	ErrChannelLockSecretHashNotFound = NewError(5004, "ChannelNoSuchLock")
	//ErrChannelEndStateNoSuchLock 通道当前参与方中找不到相应的锁
	ErrChannelEndStateNoSuchLock = NewError(5005, "ErrChannelEndStateNoSuchLock")
	//ErrChannelLockAlreadyExpired 通道中锁已过期
	ErrChannelLockAlreadyExpired = NewError(5006, "ErrChannelLockAlreadyExpired")
	//ErrChannelBalanceDecrease 发生了降低通道balance(指的是合约中的balance)的行为
	ErrChannelBalanceDecrease = NewError(5007, "ErrChannelBalanceDecrease")
	//ErrChannelTransferAmountMismatch 收到的交易中transferamount不匹配
	ErrChannelTransferAmountMismatch = NewError(5008, "ErrChannelTransferAmountMismatch")
	//ErrChannelBalanceProofAlreadyRegisteredOnChain  已经提交过balanceproof以后试图修改本地balance proof
	ErrChannelBalanceProofAlreadyRegisteredOnChain = NewError(5009, "ErrChannelBalanceProofAlreadyRegisteredOnChain")
	//ErrChannelDuplicateLock 通道中已存在该密码的锁
	ErrChannelDuplicateLock = NewError(5010, "ErrChannelDuplicateLock")
	//ErrChannelTransferAmountDecrease 收到交易,TransferAmount变小了
	ErrChannelTransferAmountDecrease = NewError(5011, "ErrChannelTransferAmountDecrease")
	//ErrRemoveNotExpiredLock 试图移除没有过期的锁
	ErrRemoveNotExpiredLock = NewError(5012, "ErrRemoveNotExpiredLock")
	//ErrUpdateBalanceProofAfterClosed 试图在通道关闭以后还更新对方或者我的balance proof,基本意思和ErrChannelBalanceProofAlreadyRegisteredOnChain一样
	ErrUpdateBalanceProofAfterClosed = NewError(5013, "ErrUpdateBalanceProofAfterClosed")
	//ErrChannelIdentifierMismatch 通道id不匹配
	ErrChannelIdentifierMismatch = NewError(5014, "ErrChannelIdentifierMismatch")
	//ErrChannelInvalidSender 收到来自未知参与方的交易
	ErrChannelInvalidSender = NewError(5015, "ErrChannelInvalidSender")
	//ErrChannelBalanceNotMatch  合作关闭通道,取现时金额检查不匹配,
	ErrChannelBalanceNotMatch = NewError(5016, "ErrChannelBalanceNotMatch")
	//ErrChannelLockMisMatch 收到交易中指定的锁和本地不匹配
	ErrChannelLockMisMatch = NewError(5017, "ErrChannelLockMisMatch")
	//ErrChannelWithdrawAmount  合作取现的金额过大
	ErrChannelWithdrawAmount = NewError(5018, "ErrChannelWithdrawAmount")
	//ErrChannelLockExpirationTooLarge 收到交易,指定的过期时间太长了,这可能是一种攻击
	ErrChannelLockExpirationTooLarge = NewError(5019, "ErrChannelLockExpirationTooLarge")
	//ErrChannelRevealTimeout 指定的reveal timeout 非法
	ErrChannelRevealTimeout = NewError(5020, "ErrChannelRevealTimeout")
	//ErrChannelBalanceProofNil balanceproof为空
	ErrChannelBalanceProofNil = NewError(5021, "ErrChannelBalanceProofNil")
	//ErrChannelCloseClosedChannel 试图关闭已经关闭的通道
	ErrChannelCloseClosedChannel = NewError(5022, "ErrChannelCloseClosedChannel")
	//ErrChannelBackgroundTx 后台执行Tx发生错误
	ErrChannelBackgroundTx = NewError(5023, "ErrChannelBackgroundTx")

	/*ErrChannelWithdrawButHasLocks : we can't send a request for withdraw when there are locks.
	 */
	ErrChannelWithdrawButHasLocks = NewError(5024, "ErrChannelWithdrawButHasLocks")
	/*ErrChannelCooperativeSettleButHasLocks : we can't send a request for settle when there are locks.
	 */
	ErrChannelCooperativeSettleButHasLocks = NewError(5025, "ErrChannelCooperativeSettleButHasLocks")
	/*ErrChannelInvalidSettleTimeout Raised when the user provided timeout value is less than the minimum
	  settle timeout
	*/
	ErrChannelInvalidSettleTimeout = NewError(5026, "ErrInvalidSettleTimeout")
	/*ErrOpenChannelWithSelf 不能自己与自己创建通道
	 */
	ErrOpenChannelWithSelf = NewError(5027, "ErrOpenChannelWithSelf")
	/*
		Transport error
	*/

	//ErrTransportTypeUnknown  未知的transport层错误,
	ErrTransportTypeUnknown = NewError(6000, "transport type error")
	//ErrSubScribeNeighbor 订阅节点在线信息错误
	ErrSubScribeNeighbor = NewError(6001, "ErrSubScribeNeighbor")
	//ErrContractQueryError 合约查询发生错误
	ErrContractQueryError = NewError(6002, "ErrContractQueryError")

	// ErrUnknown 未知错误
	ErrUnknown = NewError(9999, "unknown error")
)
