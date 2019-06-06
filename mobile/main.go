package mobile

import (
	"os"
	"time"

	"github.com/go-errors/errors"

	"github.com/SmartMeshFoundation/Photon/utils"

	"github.com/SmartMeshFoundation/Photon/rerr"

	"fmt"

	"runtime/debug"

	"github.com/SmartMeshFoundation/Photon/cmd/photon/mainimpl"
	"github.com/SmartMeshFoundation/Photon/params"
)

var apiMonitor = make(map[*API]struct{})

func init() {
	debug.SetTraceback("crash")
}

/*
StartUp is entry point for mobile photon.
address is the Node address,such as 0x1a9ec3b0b807464e6d3398a59d6b0a369bf422fa.
keystorePath is the address of the private key,  geth keystore directory . eg ~/.geth/keystore.
ethRpcEndPoint is the URL connected to geth ,such as:ws://10.0.0.2:8546.
dataDir is the working directory of a node, such as ~/.photon .
passwordfile is the file to storage password eg ~/.geth/pass.txt .
apiAddr is  127.0.0.1:5001 for product,0.0.0.1:5001 for test .
listenAddr is the listenning address for incomming message from peers.
registryAddress is the contract address working on.
otherArgs is an array of other arguments.
todo 启动参数需要重构
1. 缺省的不用传递参数默认都不要传了,如果确实有需要可以走otherArgs
	包括(apiAddr,listenAddr,registryAddress,logFile)
2.默认启用的参数--verbosity和--debug应该去掉,尤其是--debug会自动上传日志
3. DefaultRevealTimeout 需要修改,不能在默认用3了,这个纯粹是为了测试
*/
func StartUp(address, keystorePath, ethRPCEndPoint, dataDir, passwordfile, apiAddr, listenAddr, logFile, registryAddress string, otherArgs *Strings) (api *API, err error) {
	if len(apiMonitor) > 0 {
		var s = ""
		for a := range apiMonitor {
			s += fmt.Sprintf("%s\n", a.startTime.String())
		}
		err = errors.New(utils.Marshal(rerr.ErrPhotonAlreadyRunning.WithData(s)))
		return
	}
	os.Args = make([]string, 0, 20)
	os.Args = append(os.Args, "photonmobile")
	os.Args = append(os.Args, fmt.Sprintf("--address=%s", address))
	os.Args = append(os.Args, fmt.Sprintf("--keystore-path=%s", keystorePath))
	os.Args = append(os.Args, fmt.Sprintf("--eth-rpc-endpoint=%s", ethRPCEndPoint))
	os.Args = append(os.Args, fmt.Sprintf("--datadir=%s", dataDir))
	os.Args = append(os.Args, fmt.Sprintf("--password-file=%s", passwordfile))
	os.Args = append(os.Args, fmt.Sprintf("--api-address=%s", apiAddr))
	os.Args = append(os.Args, fmt.Sprintf("--listen-address=%s", listenAddr))
	os.Args = append(os.Args, fmt.Sprintf("--ignore-mediatednode-request"))
	os.Args = append(os.Args, fmt.Sprintf("--verbosity=5")) //需要移除
	os.Args = append(os.Args, fmt.Sprintf("--registry-contract-address=%s", registryAddress))

	if len(logFile) > 0 {
		os.Args = append(os.Args, fmt.Sprintf("--logfile=%s", logFile))
	}
	if otherArgs != nil {
		os.Args = append(os.Args, otherArgs.strs...)
	}
	//panicOnNullValue()
	params.MobileMode = true
	params.DefaultRevealTimeout = 3 //todo 需要移除
	rapi, err := mainimpl.StartMain()
	if err != nil {
		err = errors.New(utils.Marshal(err))
		return
	}
	api = &API{
		startTime: time.Now(),
		api:       rapi,
	}
	apiMonitor[api] = struct{}{}
	return
}
