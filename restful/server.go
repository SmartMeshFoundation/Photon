package restful

import (
	"github.com/SmartMeshFoundation/SmartRaiden"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
	"github.com/SmartMeshFoundation/SmartRaiden/restful/v1"
)

func init() {

}

/*
Start restful server
RaidenAPI is the interface of raiden network
config is the configuration of raiden network
*/
func Start(RaidenAPI *smartraiden.RaidenAPI, config *params.Config) {
	v1.RaidenAPI = RaidenAPI
	v1.Config = config
	v1.Start()
}
