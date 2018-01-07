package v1

import "github.com/astaxie/beego"

func init() {
	ns :=
		beego.NewNamespace("/api/1",
			beego.NSRouter("/tokens/:token/partners", &Controller{}, "get:TokenPartners"),
			beego.NSRouter("/tokens/:token", &Controller{}, "put:RegisterToken"),
			//beego.NSRouter("/channels/:channel", &Controller{}, "get:SpecifiedChannel;patch:CloseSettleDepositChannel"),
			//beego.NSRouter("/channels", &Controller{}, "put:OpenChannel"),
			beego.NSRouter(" /token_swaps/:target/:id", &Controller{}, "put:TokenSwap"),
			beego.NSRouter("/transfers/:token/:target", &Controller{}, "post:Transfers"),
			beego.NSRouter("/connections/?:token", &ConnectionsController{}),
			beego.NSRouter("/channels/?:channel", &ChannelsController{}),
			beego.NSAutoRouter(&EventsController{}),
			beego.NSAutoRouter(&Controller{}),
		)
	//注册 namespace
	beego.AddNamespace(ns)
}
