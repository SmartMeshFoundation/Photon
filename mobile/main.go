package mobile

import (
	"os"

	"fmt"

	"runtime/debug"

	"github.com/SmartMeshFoundation/Photon/cmd/photon/mainimpl"
	"github.com/SmartMeshFoundation/Photon/params"
)

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
*/
func StartUp(address, keystorePath, ethRPCEndPoint, dataDir, passwordfile, apiAddr, listenAddr, logFile, registryAddress string, otherArgs *Strings) (api *API, err error) {
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
	os.Args = append(os.Args, fmt.Sprintf("--verbosity=5"))
	os.Args = append(os.Args, fmt.Sprintf("--debug"))
	os.Args = append(os.Args, fmt.Sprintf("--registry-contract-address=%s", registryAddress))
	os.Args = append(os.Args, fmt.Sprintf("--disable-fee"))
	//os.Args = append(os.Args, fmt.Sprintf("--enable-health-check"))
	if len(logFile) > 0 {
		os.Args = append(os.Args, fmt.Sprintf("--logfile=%s", logFile))
	}
	if otherArgs != nil {
		os.Args = append(os.Args, otherArgs.strs...)
	}
	//panicOnNullValue()
	params.MobileMode = true
	params.DefaultRevealTimeout = 3
	rapi, err := mainimpl.StartMain()
	if err != nil {
		return
	}
	api = &API{rapi}
	return
}
