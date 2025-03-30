package mobile

import (
	"encoding/hex"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/common"

	"errors"

	"github.com/SmartMeshFoundation/Photon/rerr"
	"github.com/SmartMeshFoundation/Photon/utils"

	"fmt"

	"runtime/debug"

	"github.com/SmartMeshFoundation/Photon/cmd/photon/mainimpl"
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
func StartUp(privateKeyBinHex, ethRPCEndPoint, dataDir, apiAddr, listenAddr, logFile, registryAddress string, otherArgs *Strings) (api *API, err error) {
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
	os.Args = append(os.Args, "--mobile")
	os.Args = append(os.Args, fmt.Sprintf("--mobile-private-key-hex=%s", privateKeyBinHex)) //直接传递私玥,避免耗费资源加载keystore
	//os.Args = append(os.Args, fmt.Sprintf("--address=%s", address))
	//os.Args = append(os.Args, fmt.Sprintf("--keystore-path=%s", keystorePath))
	os.Args = append(os.Args, fmt.Sprintf("--eth-rpc-endpoint=%s", ethRPCEndPoint))
	os.Args = append(os.Args, fmt.Sprintf("--datadir=%s", dataDir))
	//os.Args = append(os.Args, fmt.Sprintf("--password-file=%s", passwordfile))
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

//ecrecover is a wrapper for crypto.Ecrecover ,处理v=27,28的情况
func ecrecover(hash common.Hash, signature []byte) (addr common.Address, err error) {
	if len(signature) != 65 {
		err = fmt.Errorf("signature errr, len=%d,signature=%s", len(signature), hex.EncodeToString(signature))
		return
	}
	pubkey, err := crypto.Ecrecover(hash[:], signature)
	if err != nil {
		signature[len(signature)-1] -= 27 //why?
		pubkey, err = crypto.Ecrecover(hash[:], signature)
		if err != nil {
			signature[len(signature)-1] += 27
			return
		}
	}
	addr = utils.PubkeyToAddress(pubkey)
	signature[len(signature)-1] += 27
	return
}

/*
GetMinerFromSignature : 根据全节点提供的签名,以及我的钱包校验全节点地址是否正确
walletAddr: 抵押的钱包地址
sig: 矿工对WalletAddr的签名
如果不出错,应该返回矿工地址
*/
func GetMinerFromSignature(walletAddr, sig string) (minerAddr string, err error) {
	wa := common.HexToAddress(walletAddr)
	if strings.Index(sig, "0x") == 0 {
		sig = sig[2:]
	}
	signature, err := hex.DecodeString(sig)
	if err != nil {
		err = fmt.Errorf("signature format errr %s,sig=%s", err, sig)
		return
	}
	hash := crypto.Keccak256(wa.Bytes())
	ma, err := ecrecover(common.BytesToHash(hash), signature)
	if err != nil {
		err = fmt.Errorf("Ecrecover err %s", err)
		return
	}
	minerAddr = ma.String()
	return
}
