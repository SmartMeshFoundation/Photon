package v1

import (
	"net/http"

	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
	"github.com/ant0ine/go-json-rest/rest"
)

/*
RaidenAPI is the interface of raiden network
should be set before start restful server
*/
var RaidenAPI *smartraiden.RaidenAPI

/*
Config is the configuration of raiden network
should be set before start restful server
*/
var Config *params.Config

/*
Start the restful server
*/
func Start() {

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/api/1/address", Address),
		rest.Get("/api/1/tokens", Tokens),
		rest.Get("/api/1/tokens/:token/partners", TokenPartners),
		rest.Put("/api/1/tokens/:token", RegisterToken),
		/*
			transfer
		*/
		rest.Put("/api/1/token_swaps/:target/:id", TokenSwap),
		rest.Post("/api/1/transfers/:token/:target", Transfers),
		rest.Get("/api/1/querysenttransfer", GetSentTransfers),
		rest.Get("/api/1/queryreceivedtransfer", GetReceivedTransfers),
		/*
			test
		*/
		rest.Get("/api/1/stop", Stop),
		rest.Get("/api/1/switch/:mesh", SwitchNetwork),
		rest.Post("/api/1/updatenodes", UpdateMeshNetworkNodes),
		/*
			channels
		*/
		rest.Get("/api/1/channels/:channel", SpecifiedChannel),
		rest.Get("/api/1/channels", GetChannelList),
		rest.Put("/api/1/channels", OpenChannel),
		rest.Patch("/api/1/channels/:channel", CloseSettleDepositChannel),
		rest.Get("/api/1/thirdparty/:channel/:3rd", ChannelFor3rdParty),
		/*
			1. withdraw
			{ "amount":3333,}
			2. prepare for withdraw:
			{"op":"preparewithdraw",}
			3. cancel prepare:
			{"op": "cancelprepare"}
		*/
		rest.Put("/api/1/withdraw/:channel", withdraw),
		/*
			1. prepare for withdraw:
			{"op":"preparesettle",}
			3. cancel prepare:
			{"op": "cancelprepare"}
		*/
		rest.Put("/api/1/settle/:channel", nil),
		/*
			connections
		*/
		rest.Get("/api/1/connections", GetConnections),
		rest.Put("/api/1/connections/:token", ConnectTokenNetwork),
		rest.Delete("/api/1/connections/:token", LeaveTokenNetwork),

		/*
			events
		*/
		rest.Get("/api/1/events/network", EventNetwork),
		rest.Get("/api/1/events/tokens/:token", EventTokens),
		rest.Get("/api/1/events/channels/:channel", EventChannels),
		/*
			for debug only
		*/
		rest.Get("/api/1/debug/balance/:token/:addr", Balance),
		rest.Get("/api/1/debug/transfer/:token/:addr/:value", TransferToken),
		rest.Get("/api/1/debug/ethbalance/:addr", EthBalance),
		rest.Get("/api/1/debug/ethstatus", EthereumStatus),
	)
	if err != nil {
		log.Crit(fmt.Sprintf("maker router :%s", err))
	}
	api.SetApp(router)
	listen := fmt.Sprintf("%s:%d", Config.APIHost, Config.APIPort)
	log.Crit(fmt.Sprintf("http listen and serve :%s", http.ListenAndServe(listen, api.MakeHandler())))
}
