package v1

import (
	"net/http"

	"encoding/json"

	"github.com/SmartMeshFoundation/raiden-network/params"
	"github.com/ethereum/go-ethereum/common"
)

type ConnectionsController struct {
	BaseController
}

func (this *ConnectionsController) Get() {
	this.Data["json"] = RaidenApi.GetConnectionManagersInfo()
	this.ServeJSON()
}

//Connecting to a token network
func (this *ConnectionsController) Put() {
	type Req struct {
		Funds int64 `json:"funds"`
	}
	tokenstr := this.Ctx.Input.Param(":token")
	var token common.Address
	if len(tokenstr) != len(token.String()) {
		this.Abort(http.StatusInternalServerError)
		return
	}
	token = common.HexToAddress(tokenstr)
	req := &Req{}
	err := json.Unmarshal(this.Ctx.Input.RequestBody, req)
	if err != nil {
		this.Abort(http.StatusInternalServerError)
		return
	}
	if req.Funds <= 0 {
		this.Abort(http.StatusInternalServerError)
		return
	}
	err = RaidenApi.ConnectTokenNetwork(token, req.Funds, params.DEFAULT_INITIAL_CHANNEL_TARGET, params.DEFAULT_JOINABLE_FUNDS_TARGET)
	if err != nil {
		this.Abort(http.StatusInternalServerError)
		return
	}
	this.Abort(http.StatusNoContent)
}

//leave a token network
func (this *ConnectionsController) Delete() {
	type Req struct {
		OnlyReceivingChannels bool `json:"only_receiving_channels"`
	}
	tokenstr := this.Ctx.Input.Param(":token")
	var token common.Address
	if len(tokenstr) != len(token.String()) {
		this.Abort(http.StatusInternalServerError)
		return
	}
	token = common.HexToAddress(tokenstr)
	req := &Req{}
	err := json.Unmarshal(this.Ctx.Input.RequestBody, req)
	if err != nil {
		this.Abort(http.StatusInternalServerError)
		return
	}
	chs, err := RaidenApi.LeaveTokenNetwork(token, req.OnlyReceivingChannels)
	if err != nil {
		this.Abort(http.StatusInternalServerError)
		return
	} else if len(chs) > 0 {
		var addrs []string
		for _, c := range chs {
			addrs = append(addrs, c.MyAddress.String())
		}
		this.Data["json"] = addrs
		this.ServeJSON()
		return
	} else {
		this.Abort(http.StatusNoContent)
	}

}
