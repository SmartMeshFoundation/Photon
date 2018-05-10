package smartraiden

import (
	"math/big"

	"github.com/SmartMeshFoundation/SmartRaiden/network"
	"github.com/ethereum/go-ethereum/common"
	"github.com/kataras/iris/utils"
)

/*
request from user
todo  we need a seperate rpc server .
*/
//key for map, no pointer
type SwapKey struct {
	Identifier uint64
	FromToken  common.Address
	FromAmount string //string of  big int
}
type TokenSwap struct {
	Identifier      uint64
	FromToken       common.Address
	FromAmount      *big.Int
	FromNodeAddress common.Address //the node address of the owner of the `from_token`
	ToToken         common.Address
	ToAmount        *big.Int
	ToNodeAddress   common.Address //the node address of the owner of the `to_token`
}

const TransferReqName = "transfer"
const NewChannelReqName = "newchannel"
const CloseChannelReqName = "closechannel"
const SettleChannelReqName = "settlechannel"
const DepositChannelReqName = "deposit"
const TokenSwapMakerReqName = "tokenswapmaker"
const TokenSwapTakerReqName = "tokenswaptaker"

/*
transfer api
*/
type TransferReq struct {
	TokenAddress common.Address
	Amount       *big.Int
	Target       common.Address
	Identifier   uint64
	Fee          *big.Int
}

/*
new channel api
*/
type NewChannelReq struct {
	tokenAddress   common.Address
	partnerAddress common.Address
	settleTimeout  int
}

/*
close channel api
settle channel api
*/
type CloseSettleChannelReq struct {
	addr common.Address //channel address
}

/*
depsoit  to channel api
*/
type DepositChannelReq struct {
	addr   common.Address
	amount *big.Int
}

/*
maker's token swap
*/
type TokenSwapMakerReq struct {
	tokenSwap *TokenSwap
}

/*
taker's token swap api
*/
type TokenSwapTakerReq struct {
	tokenSwap *TokenSwap
}

/*
general req's wraper
*/
type ApiReq struct {
	ReqId  string
	Name   string      //operation name
	Req    interface{} //operatoin
	result chan *network.AsyncResult
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
func (this *RaidenService) MediatedTransferAsyncClient(tokenAddress common.Address, amount *big.Int, fee *big.Int, target common.Address, identifier uint64) *network.AsyncResult {
	req := &ApiReq{
		ReqId: utils.RandomString(10),
		Name:  TransferReqName,
		Req: &TransferReq{
			TokenAddress: tokenAddress,
			Amount:       amount,
			Target:       target,
			Identifier:   identifier,
			Fee:          fee,
		},
	}
	return this.sendReqClient(req)
	//return this.StartMediatedTransfer(tokenAddress, target, amount, identifier)
}
func (this *RaidenService) sendReqClient(req *ApiReq) *network.AsyncResult {
	req.result = make(chan *network.AsyncResult, 1)
	this.UserReqChan <- req
	ar := <-req.result
	return ar
}
func (this *RaidenService) NewChannelClient(token, partner common.Address, settleTimeout int) *network.AsyncResult {
	req := &ApiReq{
		ReqId: utils.RandomString(10),
		Name:  NewChannelReqName,
		Req: &NewChannelReq{
			tokenAddress:   token,
			partnerAddress: partner,
			settleTimeout:  settleTimeout,
		},
	}
	return this.sendReqClient(req)
}
func (this *RaidenService) DepositChannelClient(channelAddres common.Address, amount *big.Int) *network.AsyncResult {
	req := &ApiReq{
		ReqId: utils.RandomString(10),
		Name:  DepositChannelReqName,
		Req: &DepositChannelReq{
			addr:   channelAddres,
			amount: amount,
		},
	}
	return this.sendReqClient(req)
}
func (this *RaidenService) CloseChannelClient(channelAddress common.Address) *network.AsyncResult {
	req := &ApiReq{
		ReqId: utils.RandomString(10),
		Name:  CloseChannelReqName,
		Req: &CloseSettleChannelReq{
			addr: channelAddress,
		},
	}
	return this.sendReqClient(req)
}
func (this *RaidenService) SettleChannelClient(channelAddress common.Address) *network.AsyncResult {
	req := &ApiReq{
		ReqId: utils.RandomString(10),
		Name:  SettleChannelReqName,
		Req: &CloseSettleChannelReq{
			addr: channelAddress,
		},
	}
	return this.sendReqClient(req)
}
func (this *RaidenService) TokenSwapMakerClient(tokenswap *TokenSwap) *network.AsyncResult {
	req := &ApiReq{
		ReqId: utils.RandomString(10),
		Name:  TokenSwapMakerReqName,
		Req:   &TokenSwapMakerReq{tokenswap},
	}
	return this.sendReqClient(req)
}
func (this *RaidenService) TokenSwapTakerClient(tokenswap *TokenSwap) *network.AsyncResult {
	req := &ApiReq{
		ReqId: utils.RandomString(10),
		Name:  TokenSwapTakerReqName,
		Req:   &TokenSwapTakerReq{tokenswap},
	}
	return this.sendReqClient(req)
}
