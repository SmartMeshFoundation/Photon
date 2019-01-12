package photon

import (
	"math/big"

	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

/*
request from user
todo  we need a seperate rpc server .
*/
//key for map, no pointer
type swapKey struct {
	LockSecretHash common.Hash
	FromToken      common.Address
	FromAmount     string //string of  big int
}

//TokenSwap for tokenswap api
type TokenSwap struct {
	LockSecretHash  common.Hash
	Secret          common.Hash // maker will use
	FromToken       common.Address
	FromAmount      *big.Int
	FromNodeAddress common.Address //the node address of the owner of the `from_token`
	ToToken         common.Address
	ToAmount        *big.Int
	ToNodeAddress   common.Address //the node address of the owner of the `to_token`
}

const transferReqName = "transfer"
const newChannelReqName = "newchannel"
const closeChannelReqName = "closechannel"
const settleChannelReqName = "settlechannel"
const cooperativeSettleChannelReqName = "cooperativeSettlechannel"
const prepareForCooperativeSettleReqName = "mark channel cooperative settle"
const cancelPrepareForCooperativeSettleReqName = "cancel mark cooperative settle"
const withdrawReqName = "withdraw"
const prepareWithdrawReqName = "mark withdraw"
const cancelPrepareWithdrawReqName = "cancel mark withdraw"
const tokenSwapMakerReqName = "tokenswapmaker"
const tokenSwapTakerReqName = "tokenswaptaker"
const cancelTransfer = "canceltransfer"
const allowRevealSecretReqName = "AllowRevealSecret"
const registerSecretReqName = "RegisterSecret"
const getUnfinishedReceviedTransferReqName = "GetUnfinishedReceivedTransfer"
const forceUnlockReqName = "ForceUnlock"

/*
transfer api
*/
type transferReq struct {
	TokenAddress     common.Address
	Amount           *big.Int
	Target           common.Address
	Fee              *big.Int
	Secret           common.Hash
	IsDirectTransfer bool
	Data             string
}

/*
new channel api
*/
type newChannelReq struct {
	tokenAddress   common.Address
	partnerAddress common.Address
	settleTimeout  int
	amount         *big.Int
	isNewChannel   bool
}

/*
close channel api
settle channel api
*/
type closeSettleChannelReq struct {
	addr common.Hash //channel address
}

type withdrawReq struct {
	addr   common.Hash //channel address
	amount *big.Int
}

/*
depsoit  to channel api
*/
type depositChannelReq struct {
	addr   common.Hash
	amount *big.Int
}

/*
maker's token swap
*/
type tokenSwapMakerReq struct {
	tokenSwap *TokenSwap
}

/*
taker's token swap api
*/
type tokenSwapTakerReq struct {
	tokenSwap *TokenSwap
}

/*
cancel transfer api
*/
type cancelTransferReq struct {
	LockSecretHash common.Hash
	TokenAddress   common.Address
}

/*
general req's wraper
*/
type apiReq struct {
	ReqID  string
	Name   string      //operation name
	Req    interface{} //operatoin
	result chan *utils.AsyncResult
}

/*
Transfer `amount` between this node and `target`.

       This method will start an asyncronous transfer, the transfer might fail
       or succeed depending on a couple of factors:

           - Existence of a path that can be used, through the usage of direct
             or intermediary channels.
           - Network speed, making the transfer sufficiently fast so it doesn't
             expire.
*/
func (rs *Service) transferAsyncClient(tokenAddress common.Address, amount *big.Int, fee *big.Int, target common.Address, secret common.Hash, isDirectTransfer bool, data string) *utils.AsyncResult {
	req := &apiReq{
		ReqID: utils.RandomString(10),
		Name:  transferReqName,
		Req: &transferReq{
			TokenAddress:     tokenAddress,
			Amount:           amount,
			Target:           target,
			Secret:           secret,
			Fee:              fee,
			IsDirectTransfer: isDirectTransfer,
			Data:             data,
		},
	}
	return rs.sendReqClient(req)
	//return rs.startMediatedTransfer(tokenAddress, target, amount, identifier)
}
func (rs *Service) sendReqClient(req *apiReq) *utils.AsyncResult {
	req.result = make(chan *utils.AsyncResult, 1)
	rs.UserReqChan <- req
	ar := <-req.result
	return ar
}
func (rs *Service) depositAndOpenChannelClient(token, partner common.Address, settleTimeout int, deposit *big.Int, isNewChannel bool) *utils.AsyncResult {
	req := &apiReq{
		ReqID: utils.RandomString(10),
		Name:  newChannelReqName,
		Req: &newChannelReq{
			tokenAddress:   token,
			partnerAddress: partner,
			settleTimeout:  settleTimeout,
			amount:         deposit,
			isNewChannel:   isNewChannel,
		},
	}
	return rs.sendReqClient(req)
}
func (rs *Service) closeChannelClient(channelIdentifier common.Hash) *utils.AsyncResult {
	req := &apiReq{
		ReqID: utils.RandomString(10),
		Name:  closeChannelReqName,
		Req: &closeSettleChannelReq{
			addr: channelIdentifier,
		},
	}
	return rs.sendReqClient(req)
}
func (rs *Service) settleChannelClient(channelIdentifier common.Hash) *utils.AsyncResult {
	req := &apiReq{
		ReqID: utils.RandomString(10),
		Name:  settleChannelReqName,
		Req: &closeSettleChannelReq{
			addr: channelIdentifier,
		},
	}
	return rs.sendReqClient(req)
}
func (rs *Service) cooperativeSettleChannelClient(channelIdentifier common.Hash) *utils.AsyncResult {
	req := &apiReq{
		ReqID: utils.RandomString(10),
		Name:  cooperativeSettleChannelReqName,
		Req: &closeSettleChannelReq{
			addr: channelIdentifier,
		},
	}
	return rs.sendReqClient(req)
}
func (rs *Service) markChannelForCooperativeSettleClient(channelIdentifier common.Hash) *utils.AsyncResult {
	req := &apiReq{
		ReqID: utils.RandomString(10),
		Name:  prepareForCooperativeSettleReqName,
		Req: &closeSettleChannelReq{
			addr: channelIdentifier,
		},
	}
	return rs.sendReqClient(req)
}
func (rs *Service) cancelMarkChannelForCooperativeSettleClient(channelIdentifier common.Hash) *utils.AsyncResult {
	req := &apiReq{
		ReqID: utils.RandomString(10),
		Name:  cancelPrepareForCooperativeSettleReqName,
		Req: &closeSettleChannelReq{
			addr: channelIdentifier,
		},
	}
	return rs.sendReqClient(req)
}
func (rs *Service) withdrawClient(channelIdentifier common.Hash, amount *big.Int) *utils.AsyncResult {
	req := &apiReq{
		ReqID: utils.RandomString(10),
		Name:  withdrawReqName,
		Req: &withdrawReq{
			addr:   channelIdentifier,
			amount: amount,
		},
	}
	return rs.sendReqClient(req)
}
func (rs *Service) markWithdrawClient(channelIdentifier common.Hash) *utils.AsyncResult {
	req := &apiReq{
		ReqID: utils.RandomString(10),
		Name:  prepareWithdrawReqName,
		Req: &closeSettleChannelReq{
			addr: channelIdentifier,
		},
	}
	return rs.sendReqClient(req)
}
func (rs *Service) cancelMarkWithdrawClient(channelIdentifier common.Hash) *utils.AsyncResult {
	req := &apiReq{
		ReqID: utils.RandomString(10),
		Name:  cancelPrepareWithdrawReqName,
		Req: &closeSettleChannelReq{
			addr: channelIdentifier,
		},
	}
	return rs.sendReqClient(req)
}
func (rs *Service) tokenSwapMakerClient(tokenswap *TokenSwap) *utils.AsyncResult {
	req := &apiReq{
		ReqID: utils.RandomString(10),
		Name:  tokenSwapMakerReqName,
		Req:   &tokenSwapMakerReq{tokenswap},
	}
	return rs.sendReqClient(req)
}
func (rs *Service) tokenSwapTakerClient(tokenswap *TokenSwap) *utils.AsyncResult {
	req := &apiReq{
		ReqID: utils.RandomString(10),
		Name:  tokenSwapTakerReqName,
		Req:   &tokenSwapTakerReq{tokenswap},
	}
	return rs.sendReqClient(req)
}
func (rs *Service) cancelTransferClient(lockSecretHash common.Hash, tokenAddress common.Address) *utils.AsyncResult {
	req := &apiReq{
		ReqID: utils.RandomString(10),
		Name:  cancelTransfer,
		Req: &cancelTransferReq{
			LockSecretHash: lockSecretHash,
			TokenAddress:   tokenAddress,
		},
	}
	return rs.sendReqClient(req)
}

type allowRevealSecretReq struct {
	LockSecretHash common.Hash
	TokenAddress   common.Address
}

func (rs *Service) allowRevealSecretClient(lockSecretHash common.Hash, tokenAddress common.Address) *utils.AsyncResult {
	req := &apiReq{
		ReqID: utils.RandomString(10),
		Name:  allowRevealSecretReqName,
		Req: &allowRevealSecretReq{
			LockSecretHash: lockSecretHash,
			TokenAddress:   tokenAddress,
		},
	}
	return rs.sendReqClient(req)
}

type registerSecretReq struct {
	Secret       common.Hash
	TokenAddress common.Address
}

func (rs *Service) registerSecretClient(secret common.Hash, tokenAddress common.Address) *utils.AsyncResult {
	req := &apiReq{
		ReqID: utils.RandomString(10),
		Name:  registerSecretReqName,
		Req: &registerSecretReq{
			Secret:       secret,
			TokenAddress: tokenAddress,
		},
	}
	return rs.sendReqClient(req)
}

type getUnfinishedReceivedTransferReq struct {
	LockSecretHash common.Hash
	TokenAddress   common.Address
}

func (rs *Service) getUnfinishedReceivedTransferClient(lockSecretHash common.Hash, tokenAddress common.Address) *utils.AsyncResult {
	req := &apiReq{
		ReqID: utils.RandomString(10),
		Name:  getUnfinishedReceviedTransferReqName,
		Req: &getUnfinishedReceivedTransferReq{
			LockSecretHash: lockSecretHash,
			TokenAddress:   tokenAddress,
		},
	}
	return rs.sendReqClient(req)
}

type forceUnlockReq struct {
	LockSecretHash    common.Hash
	Secret            common.Hash
	ChannelIdentifier common.Hash
}

func (rs *Service) forceUnlockClient(lockSecretHash, secret, channelIdentifier common.Hash) *utils.AsyncResult {
	req := &apiReq{
		ReqID: utils.RandomString(10),
		Name:  forceUnlockReqName,
		Req: &forceUnlockReq{
			LockSecretHash:    lockSecretHash,
			ChannelIdentifier: channelIdentifier,
			Secret:            secret,
		},
	}
	return rs.sendReqClient(req)
}
