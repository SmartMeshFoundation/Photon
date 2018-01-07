package v1

import (
	"fmt"

	"github.com/SmartMeshFoundation/raiden-network"
	"github.com/SmartMeshFoundation/raiden-network/params"
	"github.com/astaxie/beego"
	"github.com/fatedier/frp/src/utils/log"
)

var RaidenApi *raiden_network.RaidenApi
var Config *params.Config

func Start() {
	beego.BConfig.RunMode = beego.DEV
	beego.BConfig.CopyRequestBody = true
	log.Info(fmt.Sprintf("api server running on %s:%d", Config.ApiHost, Config.ApiPort))
	beego.Run(fmt.Sprintf("%s:%d", Config.ApiHost, Config.ApiPort))
}
