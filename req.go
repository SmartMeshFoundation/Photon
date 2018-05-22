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
type swapKey struct {
	Identifier uint64
	FromToken  common.Address
	FromAmount string //string of  big int
}

//TokenSwap for tokenswap api
type TokenSwap struct {
	Identifier      uint64
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
const depositChannelReqName = "deposit"
const tokenSwapMakerReqName = "tokenswapmaker"
const tokenSwapTakerReqName = "tokenswaptaker"

/*
transfer api
*/
type transferReq struct {
	TokenAddress common.Address
	Amount       *big.Int
	Target       common.Address
	Identifier   uint64
	Fee          *big.Int
}

/*
new channel api
*/
type newChannelReq struct {
	tokenAddress   common.Address
	partnerAddress common.Address
	settleTimeout  int
}

/*
close channel api
settle channel api
*/
type closeSettleChannelReq struct {
	addr common.Address //channel address
}

/*
depsoit  to channel api
*/
type depositChannelReq struct {
	addr   common.Address
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
general req's wraper
*/
type apiReq struct {
	ReqID  string
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
func (rs *RaidenService) mediatedTransferAsyncClient(tokenAddress common.Address, amount *big.Int, fee *big.Int, target common.Address, identifier uint64) *network.AsyncResult {
	req := &apiReq{
		ReqID: utils.RandomString(10),
		Name:  transferReqName,
		Req: &transferReq{
			TokenAddress: tokenAddress,
			Amount:       amount,
			Target:       target,
			Identifier:   identifier,
			Fee:          fee,
		},
	}
	return rs.sendReqClient(req)
	//return rs.startMediatedTransfer(tokenAddress, target, amount, identifier)
}
func (rs *RaidenService) sendReqClient(req *apiReq) *network.AsyncResult {
	req.result = make(chan *network.AsyncResult, 1)
	rs.UserReqChan <- req
	ar := <-req.result
	return ar
}
func (rs *RaidenService) newChannelClient(token, partner common.Address, settleTimeout int) *network.AsyncResult {
	req := &apiReq{
		ReqID: utils.RandomString(10),
		Name:  newChannelReqName,
		Req: &newChannelReq{
			tokenAddress:   token,
			partnerAddress: partner,
			settleTimeout:  settleTimeout,
		},
	}
	return rs.sendReqClient(req)
}
func (rs *RaidenService) depositChannelClient(channelAddres common.Address, amount *big.Int) *network.AsyncResult {
	req := &apiReq{
		ReqID: utils.RandomString(10),
		Name:  depositChannelReqName,
		Req: &depositChannelReq{
			addr:   channelAddres,
			amount: amount,
		},
	}
	return rs.sendReqClient(req)
}
func (rs *RaidenService) closeChannelClient(channelAddress common.Address) *network.AsyncResult {
	req := &apiReq{
		ReqID: utils.RandomString(10),
		Name:  closeChannelReqName,
		Req: &closeSettleChannelReq{
			addr: channelAddress,
		},
	}
	return rs.sendReqClient(req)
}
func (rs *RaidenService) settleChannelClient(channelAddress common.Address) *network.AsyncResult {
	req := &apiReq{
		ReqID: utils.RandomString(10),
		Name:  settleChannelReqName,
		Req: &closeSettleChannelReq{
			addr: channelAddress,
		},
	}
	return rs.sendReqClient(req)
}
func (rs *RaidenService) tokenSwapMakerClient(tokenswap *TokenSwap) *network.AsyncResult {
	req := &apiReq{
		ReqID: utils.RandomString(10),
		Name:  tokenSwapMakerReqName,
		Req:   &tokenSwapMakerReq{tokenswap},
	}
	return rs.sendReqClient(req)
}
func (rs *RaidenService) tokenSwapTakerClient(tokenswap *TokenSwap) *network.AsyncResult {
	req := &apiReq{
		ReqID: utils.RandomString(10),
		Name:  tokenSwapTakerReqName,
		Req:   &tokenSwapTakerReq{tokenswap},
	}
	return rs.sendReqClient(req)
}
