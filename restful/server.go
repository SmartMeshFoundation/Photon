package restful

import (
	"github.com/SmartMeshFoundation/raiden-network"
	"github.com/SmartMeshFoundation/raiden-network/params"
	"github.com/SmartMeshFoundation/raiden-network/restful/v1"
)

func init() {

}

func Start(RaidenApi *raiden_network.RaidenApi, config *params.Config) {
	v1.RaidenApi = RaidenApi
	v1.Config = config
	v1.Start()
}
