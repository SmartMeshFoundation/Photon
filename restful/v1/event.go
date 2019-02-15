package v1

import (
	"fmt"

	"github.com/SmartMeshFoundation/Photon/rerr"

	"net/url"
	"strconv"

	"github.com/SmartMeshFoundation/Photon/dto"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ethereum/go-ethereum/common"
)

/*
EventNetwork returns all events related to Photon network
*/
func EventNetwork(w rest.ResponseWriter, r *rest.Request) {
	var resp *dto.APIResponse
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> EventNetwork ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()
	fromBlock, toBlock := getFromTo(r)
	events, err := API.GetNetworkEvents(fromBlock, toBlock)
	resp = dto.NewAPIResponse(err, events)
}

/*
EventTokens returns all events about the token specified
*/
func EventTokens(w rest.ResponseWriter, r *rest.Request) {
	var resp *dto.APIResponse
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> EventTokens ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()
	fromBlock, toBlock := getFromTo(r)
	var token common.Address
	tokenstr := r.PathParam("token")
	fmt.Println("tokenstr ", tokenstr)
	if len(tokenstr) != len(token.String()) {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError)
		return
	}
	token, err := utils.HexToAddress(tokenstr)
	if err != nil {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError.AppendError(err))
		return
	}
	events, err := API.GetTokenNetworkEvents(token, fromBlock, toBlock)
	resp = dto.NewAPIResponse(err, events)
}

/*
EventChannels returns all events about the channel specified
*/
func EventChannels(w rest.ResponseWriter, r *rest.Request) {
	var resp *dto.APIResponse
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> EventChannels ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()
	fromBlock, toBlock := getFromTo(r)
	log.Trace(fmt.Sprintf("from=%d,toblock=%d", fromBlock, toBlock))
	var channel common.Hash
	channelstr := r.PathParam("channel")
	log.Trace(fmt.Sprintf("channels %s", channelstr))
	if len(channelstr) != len(channel.String()) {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError)
		return
	}
	channel = common.HexToHash(channelstr)
	events, err := API.GetChannelEvents(channel, fromBlock, toBlock)
	resp = dto.NewAPIResponse(err, events)
}

func getFromTo(r *rest.Request) (fromBlock, toBlock int64) {
	fromBlock = -1
	toBlock = -1
	var fromBlockStr = ""
	var toBlockStr = ""
	m, err := url.ParseQuery(r.Request.URL.RawQuery)
	if err != nil {
		log.Error(fmt.Sprintf("ParseQuery err %s", err))
		return
	}
	if len(m["from_block"]) > 0 {
		fromBlockStr = m["from_block"][0]
	}
	if len(m["to_block"]) > 0 {
		toBlockStr = m["to_block"][0]
	}
	if _, err := strconv.Atoi(fromBlockStr); err == nil {
		fromBlock, err = strconv.ParseInt(fromBlockStr, 10, 64)
		if err != nil {
			log.Error(fmt.Sprintf("fromBlock %s parse err %s", fromBlockStr, err))
		}
	}
	if _, err := strconv.Atoi(toBlockStr); err == nil {
		toBlock, err = strconv.ParseInt(toBlockStr, 10, 64)
		if err != nil {
			log.Error(fmt.Sprintf("toBlock %s parse err %s", toBlockStr, err))
		}
	}
	return fromBlock, toBlock
}
