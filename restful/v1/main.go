package v1

import (
	"context"
	"net/http"
	"os"

	"fmt"

	photon "github.com/SmartMeshFoundation/Photon"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/params"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ant0ine/go-json-rest/rest"
)

/*
API is the interface of Photon network
should be set before start restful server
*/
var API *photon.API

//QuitChan stop http server
var QuitChan chan struct{}

/*
Start the restful server
*/
func Start() {
	QuitChan = make(chan struct{})
	api := rest.NewApi()
	if params.Cfg.Debug {
		api.Use(rest.DefaultDevStack...)
	} else {
		api.Use(rest.DefaultProdStack...)
	}
	api.Use(rest.DefaultDevStack...)
	if params.Cfg.HTTPUsername != "" && params.Cfg.HTTPPassword != "" {
		api.Use(&rest.AuthBasicMiddleware{
			Realm: "please input username and password",
			Authenticator: func(userId string, password string) bool {
				return userId == params.Cfg.HTTPUsername && password == params.Cfg.HTTPPassword
			},
		})
	}
	router, err := rest.MakeRouter(

		/*
			prepare update
		*/
		rest.Post("/api/1/prepare-update", PrepareUpdate),
		/*
			transfers
		*/
		rest.Get("/api/1/querysenttransfer", GetSentTransferDetails),
		rest.Get("/api/1/queryreceivedtransfer", GetReceivedTransfers),
		rest.Post("/api/1/transfers/:token/:target", Transfers),
		rest.Get("/api/1/transferstatus/:token/:locksecrethash", GetSentTransferDetail),
		rest.Post("/api/1/transfercancel/:token/:locksecrethash", CancelTransfer),
		/*
			transfer with specified secret
		*/
		rest.Post("/api/1/transfers/allowrevealsecret", AllowRevealSecret),
		rest.Get("/api/1/getunfinishedreceivedtransfer/:tokenaddress/:locksecrethash", GetUnfinishedReceivedTransfer),
		rest.Post("/api/1/registersecret", RegisterSecret),
		/*
			token swap
		*/
		rest.Put("/api/1/token_swaps/:target/:locksecrethash", TokenSwap),
		/*
			accounts
		*/
		rest.Get("/api/1/address", Address),
		rest.Get("/api/1/balance", GetBalanceByTokenAddress),
		rest.Get("/api/1/balance/", GetBalanceByTokenAddress),
		rest.Get("/api/1/balance/:tokenaddress", GetBalanceByTokenAddress),
		/*
			channels
		*/
		rest.Get("/api/1/channels/:channel", SpecifiedChannel),
		rest.Get("/api/1/channels", GetChannelList),
		rest.Patch("/api/1/channels/:channel", CloseSettleChannel),
		rest.Put("/api/1/channels/preparecooperatesettle/:channel", prepareCooperateSettle),
		rest.Put("/api/1/channels/cancelcooperatesettle/:channel", cancelCooperateSettle),
		rest.Get("/api/1/thirdparty/:channel/:3rd", ChannelFor3rdParty),
		rest.Get("/api/1/channel-settle-block/:channel", GetChannelSettleBlock),

		/*
			Deposit
		*/
		rest.Put("/api/1/deposit", Deposit),
		/*
			tokens
		*/
		rest.Get("/api/1/tokens", Tokens),
		rest.Get("/api/1/tokens/:token/partners", TokenPartners),
		/*
			contract call tx
		*/
		rest.Post("/api/1/tx/query", ContractCallTXQuery),
		/*
			utils
		*/
		rest.Get("/api/1/path/:target_address/:token/:amount", FindPath),
		rest.Get("/api/1/secret", GetRandomSecret), // api to provide random secret and lockSecretHash pair
		rest.Get("/api/1/version", GetBuildInfo),

		/*
			fee policy
		*/
		rest.Get("/api/1/fee_policy", GetFeePolicy),
		rest.Post("/api/1/fee_policy", SetFeePolicy),
		rest.Get("/api/1/fee", GetAllFeeChargeRecord),

		/*
			income
		*/
		rest.Post("/api/1/income/details", GetIncomeDetails),
		rest.Post("/api/1/income/days", GetDaysIncome),

		/*
			assets
		*/
		rest.Post("/api/1/assets", GetAssetsOnToken),
		/*
			test
		*/
		rest.Get("/api/1/stop", Stop),
		rest.Get("/api/1/switch/:mesh", SwitchNetwork),
		rest.Post("/api/1/updatenodes", UpdateMeshNetworkNodes),

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
			events
		*/
		//rest.Get("/api/1/events/network", EventNetwork),
		//rest.Get("/api/1/events/tokens/:token", EventTokens),
		//rest.Get("/api/1/events/channels/:channel", EventChannels),
		/*
			blockchain proxy
		*/
		rest.Post("/api/1/transfer-smt/:addr/:value", TransferSMT),
		/*
			for debug only
		*/
		rest.Get("/api/1/system-status", GetSystemStatus),
		rest.Get("/api/1/debug/balance/:token/:addr", Balance),
		rest.Get("/api/1/debug/transfer/:token/:addr/:value", TransferToken),
		rest.Get("/api/1/debug/ethbalance/:addr", EthBalance),
		rest.Get("/api/1/debug/ethstatus", EthereumStatus),
		rest.Get("/api/1/debug/force-unlock/:channel/:secret", ForceUnlock),
		rest.Get("/api/1/debug/register-secret-onchain/:secret", RegisterSecretOnChain),
		rest.Get("/api/1/debug/pfs/:channel", BalanceUpdateForPFS),
		rest.Post("/api/1/debug/notify_network_down", NotifyNetworkDown), // notify photon network down
		rest.Get("/api/1/debug/shutdown", func(writer rest.ResponseWriter, request *rest.Request) {
			API.Photon.Stop()
			utils.SystemExit(0)
		}),
		rest.Get("/api/1/debug/change-eth-rpc-endpoint-port/:port", ChangeEthRPCEndpointPort),
		rest.Get("/api/1/debug/upload-log-file", UploadLogFile),
	)
	if err != nil {
		panic(fmt.Sprintf("maker router :%s", err))
	}
	api.SetApp(router)
	listen := fmt.Sprintf("%s:%d", params.Cfg.RestAPIHost, params.Cfg.RestAPIPort)
	server := &http.Server{Addr: listen, Handler: api.MakeHandler()}
	go func() {
		err2 := server.ListenAndServe()
		if err2 != nil {
			log.Error(fmt.Sprintf("ListenAndServe err %s", err2))
		}
	}()
	<-QuitChan
	err = server.Shutdown(context.Background())
	if err != nil {
		panic(fmt.Sprintf("server shutdown err %s", err))
	}
}

/*
Stop for app user, call this api before quit.
*/
func Stop(w rest.ResponseWriter, r *rest.Request) {
	defer close(QuitChan)
	defer os.Exit(0)
	//test only
	API.Stop()
	w.Header().Set("Content-Type", "text/plain")
	_, err := w.(http.ResponseWriter).Write([]byte("ok"))
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}
