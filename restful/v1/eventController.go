package v1

import (
	"strconv"

	"net/http"

	"fmt"

	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/fatedier/frp/src/utils/log"
)

type EventsController struct {
	BaseController
}

func (this *EventsController) getFromTo() (fromBlock, toBlock int64) {
	fromBlockStr := this.GetString("from_block")
	toBlockStr := this.GetString("to_block")
	fromBlock = -1
	toBlock = -1
	if _, err := strconv.Atoi(fromBlockStr); err == nil {
		fromBlock, _ = strconv.ParseInt(fromBlockStr, 10, 64)
	}
	if _, err := strconv.Atoi(toBlockStr); err == nil {
		toBlock, _ = strconv.ParseInt(toBlockStr, 10, 64)
	}
	return fromBlock, toBlock
}
func (this *EventsController) Network() {
	fromBlock, toBlock := this.getFromTo()
	events, err := RaidenApi.GetNetworkEvents(fromBlock, toBlock)
	if err != nil {
		if err != nil {
			log.Error(err.Error())
		}
		this.Abort(http.StatusInternalServerError)
		return
	}
	this.Data["json"] = events
	this.ServeJSON()
}
func (this *EventsController) Tokens() {
	fromBlock, toBlock := this.getFromTo()
	fmt.Println(fromBlock, toBlock)
	url := this.Ctx.Input.URL()
	var token common.Address
	sep := "events/tokens/"
	i := strings.Index(url, sep)
	if i > 0 {
		tokenstr := url[i+len(sep):]
		fmt.Println("tokenstr ", tokenstr)
		if len(tokenstr) != len(token.String()) {
			this.Abort(http.StatusBadRequest)
			return
		}
		token = common.HexToAddress(tokenstr)
	}
	events, err := RaidenApi.GetTokenNetworkEvents(token, fromBlock, toBlock)
	if err != nil {
		if err != nil {
			log.Error(err.Error())
		}
		this.Abort(http.StatusInternalServerError)
		return
	}
	this.Data["json"] = events
	this.ServeJSON()
}
func (this *EventsController) Channels() {
	fromBlock, toBlock := this.getFromTo()
	fmt.Println(fromBlock, toBlock)
	url := this.Ctx.Input.URL()
	var channel common.Address
	sep := "events/channels/"
	i := strings.Index(url, sep)
	if i > 0 {
		tokenstr := url[i+len(sep):]
		fmt.Println("channels ", tokenstr)
		if len(tokenstr) != len(channel.String()) {
			this.Abort(http.StatusBadRequest)
			return
		}
		channel = common.HexToAddress(tokenstr)
	}
	events, err := RaidenApi.GetChannelEvents(channel, fromBlock, toBlock)
	if err != nil {
		if err != nil {
			log.Error(err.Error())
		}
		this.Abort(http.StatusInternalServerError)
		return
	}
	this.Data["json"] = events
	this.ServeJSON()
}
