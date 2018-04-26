package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ethereum/go-ethereum/common"
)

func getFromTo(r *rest.Request) (fromBlock, toBlock int64) {
	fromBlockStr := r.PathParam("from_block")
	toBlockStr := r.PathParam("to_block")
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
func EventNetwork(w rest.ResponseWriter, r *rest.Request) {
	fromBlock, toBlock := getFromTo(r)
	events, err := RaidenApi.GetNetworkEvents(fromBlock, toBlock)
	if err != nil {
		log.Error(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(events)
}
func EventTokens(w rest.ResponseWriter, r *rest.Request) {
	fromBlock, toBlock := getFromTo(r)
	var token common.Address
	tokenstr := r.PathParam("token")
	fmt.Println("tokenstr ", tokenstr)
	if len(tokenstr) != len(token.String()) {
		rest.Error(w, "address error", http.StatusBadRequest)
		return
	}
	token = common.HexToAddress(tokenstr)
	events, err := RaidenApi.GetTokenNetworkEvents(token, fromBlock, toBlock)
	if err != nil {
		log.Error(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(events)
}
func EventChannels(w rest.ResponseWriter, r *rest.Request) {
	fromBlock, toBlock := getFromTo(r)
	log.Trace(fmt.Sprintf("from=%d,toblock=%d", fromBlock, toBlock))
	var channel common.Address
	channelstr := r.PathParam("channel")
	log.Trace(fmt.Sprintf("channels %s", channelstr))
	if len(channelstr) != len(channel.String()) {
		rest.Error(w, "adderss error", http.StatusBadRequest)
		return
	}
	channel = common.HexToAddress(channelstr)
	events, err := RaidenApi.GetChannelEvents(channel, fromBlock, toBlock)
	if err != nil {
		log.Error(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(events)
}
