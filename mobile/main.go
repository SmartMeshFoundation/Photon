package mobile

import (
	"os"

	"fmt"

	"runtime/debug"

	"github.com/SmartMeshFoundation/SmartRaiden/cmd/smartraiden/mainimpl"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
)

func init() {
	debug.SetTraceback("crash")
}

/*
StartUp is entry point for mobile raiden
address :Node address,such as 0x1a9ec3b0b807464e6d3398a59d6b0a369bf422fa
keystorePath:The address of the private key,  geth keystore directory . eg ~/.geth/keystore
ethRpcEndPoint:URL connected to geth ,such as:ws://10.0.0.2:8546
dataDir:The working directory of a node, such as ~/.smartraiden
passwordfile: file to storage password eg ~/.geth/pass.txt
apiAddr: 127.0.0.1:5001 for product,0.0.0.1:5001 for test
*/
func StartUp(address, keystorePath, ethRPCEndPoint, dataDir, passwordfile, listenAddr, logFile string, otherArgs ...string) (api *API, err error) {
	os.Args = make([]string, 0, 20)
	os.Args = append(os.Args, "smartraidenmobile")
	os.Args = append(os.Args, fmt.Sprintf("--address=%s", address))
	os.Args = append(os.Args, fmt.Sprintf("--keystore-path=%s", keystorePath))
	os.Args = append(os.Args, fmt.Sprintf("--eth-rpc-endpoint=%s", ethRPCEndPoint))
	os.Args = append(os.Args, fmt.Sprintf("--datadir=%s", dataDir))
	os.Args = append(os.Args, fmt.Sprintf("--password-file=%s", passwordfile))
	os.Args = append(os.Args, fmt.Sprintf("--listen-address=%s", listenAddr))
	os.Args = append(os.Args, fmt.Sprintf("--ignore-mediatednode-request"))
	os.Args = append(os.Args, fmt.Sprintf("--verbosity=5"))
	os.Args = append(os.Args, fmt.Sprintf("--debug"))
	//os.Args = append(os.Args, fmt.Sprintf("--enable-health-check"))
	if len(logFile) > 0 {
		os.Args = append(os.Args, fmt.Sprintf("--logfile=%s", logFile))
	}
	os.Args = append(os.Args, otherArgs...)
	//panicOnNullValue()
	params.MobileMode = true
	rapi, err := mainimpl.StartMain()
	if err != nil {
		return
	}
	api = &API{rapi}
	return
}
