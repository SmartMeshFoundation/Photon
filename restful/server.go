package restful

import (
	photon "github.com/SmartMeshFoundation/Photon"
	"github.com/SmartMeshFoundation/Photon/params"
	"github.com/SmartMeshFoundation/Photon/restful/v1"
)

func init() {

}

/*
Start restful server
PhotonAPI is the interface of photon network
config is the configuration of photon network
*/
func Start(API *photon.API, config *params.Config) {
	v1.API = API
	v1.Config = config
	v1.Start()
}
