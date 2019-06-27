package v1

import (
	"context"
	"fmt"
	"math/big"

	"github.com/SmartMeshFoundation/Photon/utils"

	"github.com/SmartMeshFoundation/Photon/dto"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/rerr"
	"github.com/ant0ine/go-json-rest/rest"
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// TransferSMT 主币转账
func TransferSMT(w rest.ResponseWriter, r *rest.Request) {
	var resp *dto.APIResponse
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> TransferSMT ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()
	// 1. 参数解析
	targetAddr := common.HexToAddress(r.PathParam("addr"))
	if targetAddr == utils.EmptyAddress {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError)
		return
	}
	valueStr := r.PathParam("value")
	value, b := new(big.Int).SetString(valueStr, 0)
	if !b {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError)
		return
	}
	if value == nil || value.Cmp(big.NewInt(0)) <= 0 {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError)
		return
	}
	// 2. 构造交易
	conn := API.Photon.Chain.Client
	ctx := context.Background()
	auth := bind.NewKeyedTransactor(API.Photon.Chain.PrivKey)
	nonceUsed := false
	nonce, err := bind.GetValidNonceAndLock(conn, auth.From, ctx)
	if err != nil {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError.Append(err.Error()))
		return
	}
	defer func() {
		bind.ConfirmNonceAndUnlock(auth.From, nonce, nonceUsed)
	}()
	msg := ethereum.CallMsg{From: auth.From, To: &targetAddr, Value: value, Data: nil}
	gasLimit, err := conn.EstimateGas(ctx, msg)
	if err != nil {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError.Append(fmt.Sprintf("failed to estimate gas needed: %v", err)))
		return
	}
	gasPrice, err := conn.SuggestGasPrice(ctx)
	if err != nil {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError.Append(fmt.Sprintf("failed to suggest gas price: %v", err)))
		return
	}
	networkID, err := conn.NetworkID(ctx)
	if err != nil {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError.Append(fmt.Sprintf("failed to get networkID : %v", err)))
		return
	}
	log.Info(fmt.Sprintf("gasLimit=%d,gasPrice=%s", gasLimit, gasPrice.String()))
	rawTx := types.NewTransaction(nonce, targetAddr, value, gasLimit, gasPrice, nil)
	// Create the transaction, sign it and schedule it for execution
	signedTx, err := auth.Signer(types.NewEIP155Signer(networkID), auth.From, rawTx)
	if err != nil {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError.Append(fmt.Sprintf("failed to sign rawTX : %v", err)))
		return
	}
	if err = conn.SendTransaction(ctx, signedTx); err != nil {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError.Append(fmt.Sprintf("failed to send rawTX : %v", err)))
		return
	}
	nonceUsed = true
	receipt, err := bind.WaitMined(ctx, conn, signedTx)
	if err != nil {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError.Append(fmt.Sprintf("tx fail : %v", err)))
		return
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError.Append("tx fail"))
		return
	}
	fmt.Printf("transfer from %s to %s amount=%s\n", auth.From.String(), targetAddr.String(), value)
	resp = dto.NewSuccessAPIResponse(nil)
	return
}
