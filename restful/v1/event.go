package v1

import (
	"fmt"
	"net/http"

	"net/url"
	"strconv"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ethereum/go-ethereum/common"
)

/*
EventNetwork returns all events related to raiden network
*/
func EventNetwork(w rest.ResponseWriter, r *rest.Request) {
	fromBlock, toBlock := getFromTo(r)
	events, err := RaidenAPI.GetNetworkEvents(fromBlock, toBlock)
	if err != nil {
		log.Error(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = w.WriteJson(events)
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}

/*
EventTokens returns all events about the token specified
*/
func EventTokens(w rest.ResponseWriter, r *rest.Request) {
	fromBlock, toBlock := getFromTo(r)
	var token common.Address
	tokenstr := r.PathParam("token")
	fmt.Println("tokenstr ", tokenstr)
	if len(tokenstr) != len(token.String()) {
		rest.Error(w, "address error", http.StatusBadRequest)
		return
	}
	token, err := utils.HexToAddress(tokenstr)
	if err != nil {
		log.Error(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	events, err := RaidenAPI.GetTokenNetworkEvents(token, fromBlock, toBlock)
	if err != nil {
		log.Error(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = w.WriteJson(events)
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}

/*
EventChannels returns all events about the channel specified
*/
func EventChannels(w rest.ResponseWriter, r *rest.Request) {
	fromBlock, toBlock := getFromTo(r)
	log.Trace(fmt.Sprintf("from=%d,toblock=%d", fromBlock, toBlock))
	var channel common.Hash
	channelstr := r.PathParam("channel")
	log.Trace(fmt.Sprintf("channels %s", channelstr))
	if len(channelstr) != len(channel.String()) {
		rest.Error(w, "adderss error", http.StatusBadRequest)
		return
	}
	channel = common.HexToHash(channelstr)
	events, err := RaidenAPI.GetChannelEvents(channel, fromBlock, toBlock)
	if err != nil {
		log.Error(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = w.WriteJson(events)
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
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
