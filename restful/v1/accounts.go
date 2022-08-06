package v1

import (
	"fmt"

	"github.com/SmartMeshFoundation/Photon/dto"
	"github.com/SmartMeshFoundation/Photon/params"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ethereum/go-ethereum/common"
)

/*
Address is api of /api/1/address
*/
func Address(w rest.ResponseWriter, r *rest.Request) {
	writejson(w, dto.NewSuccessAPIResponse(params.Cfg.MyAddress.String()))
}

/*
GetBalanceByTokenAddress : get account's balance and locked account on token
*/
func GetBalanceByTokenAddress(w rest.ResponseWriter, r *rest.Request) {
	var resp *dto.APIResponse
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> GetBalanceByTokenAddress ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()
	tokenAddressStr := r.PathParam("tokenaddress")
	var tokenAddress common.Address
	if tokenAddressStr == "" {
		tokenAddress = utils.EmptyAddress
	} else {
		tokenAddress = common.HexToAddress(tokenAddressStr)
	}
	result, err := API.GetBalanceByTokenAddress(tokenAddress)
	resp = dto.NewAPIResponse(err, result)
}
