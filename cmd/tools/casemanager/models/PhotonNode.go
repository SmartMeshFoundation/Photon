package models

import (
	"log"
	"time"

	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"

	"math/big"

	"path/filepath"

	"strings"

	"github.com/SmartMeshFoundation/Photon/params"
	"github.com/ethereum/go-ethereum/common"
)

// PhotonNodeRuntime case运行过程中存储临时数据的地方
type PhotonNodeRuntime struct {
	MainChainBalance *big.Int // 主链货币余额
}

// PhotonNode a photon node
type PhotonNode struct {
	Host          string
	Address       string
	Name          string
	APIAddress    string
	ListenAddress string
	ConditionQuit *params.ConditionQuit
	DebugCrash    bool
	Running       bool
	Runtime       PhotonNodeRuntime
}

// Start start a photon node
func (node *PhotonNode) startInternal(env *TestEnv, otherflag ...bool) {
	logfile := fmt.Sprintf("./log/%s.log", env.CaseName+"-"+node.Name)
	dopprof := false
	nonetwork := false //是否完全不和其他节点通信
	if len(otherflag) > 0 {
		dopprof = otherflag[0]
	}
	if len(otherflag) > 1 {
		nonetwork = otherflag[1]
	}
	go ExecShell(env.Main, node.getParamStr(env, dopprof, nonetwork), logfile, true)

	count := 0
	t := time.Now()
	for !node.IsRunning() {
		Logger.Printf("waiting for %s to start, sleep 100ms...\n", node.Name)
		time.Sleep(time.Millisecond * 100)
		count++
		if count > 400 {
			if node.ConditionQuit != nil {
				Logger.Printf("NODE %s %s start with %s TIMEOUT\n", node.Address, node.Host, node.ConditionQuit.QuitEvent)
			} else {
				Logger.Printf("NODE %s %s start TIMEOUT\n", node.Address, node.Host)
			}
			if !env.UseMatrix {
				panic("Start photon node TIMEOUT")
			}
		}
	}
	used := time.Since(t)
	if node.DebugCrash {
		Logger.Printf("NODE %s %s start with %s in %fs", node.Address, node.Host, node.ConditionQuit.QuitEvent, used.Seconds())
	} else {
		Logger.Printf("NODE %s %s start in %fs", node.Address, node.Host, used.Seconds())
	}
	node.Running = true
	if nonetwork { //启动以后不更新其他节点信息,这样就无法和别的节点进行通信了
		return
	}
	if !env.UseMatrix && env.XMPPServer == "" {
		for _, n := range env.Nodes {
			if n.Running {
				//n.UpdateMeshNetworkNodes(env.Nodes...)
			}
		}
	}
}

// Start start a photon node
func (node *PhotonNode) Start(env *TestEnv, otherarg ...bool) {
	node.startInternal(env, otherarg...)
	if env.UseMatrix {
		time.Sleep(time.Second * 5)
	}
}

func removeParam(params []string, remove string) []string {
	//从参数中找到remove,然后删除
	i := 0
	for i = 0; i < len(params); i++ {
		if params[i] == remove {
			break
		}
	}
	if i >= 0 && i < len(params) {
		params = append(params[:i], params[i+1:]...)
	}
	return params
}

// StartWithFeeAndPFS :
func (node *PhotonNode) StartWithFeeAndPFS(env *TestEnv) {
	logfile := fmt.Sprintf("./log/%s.log", env.CaseName+"-"+node.Name)
	params := node.getParamStr(env, false, false)
	//从参数中找到diable-fee,然后删除
	params = removeParam(params, "--disable-fee")
	params = removeParam(params, "--debug-udp-only")
	// 添加casemanager自带的pfs
	params = append(params, "--pfs=http://127.0.0.1:17000")
	go ExecShell(env.Main, params, logfile, true)

	count := 0
	t := time.Now()
	for !node.IsRunning() {
		Logger.Printf("waiting for %s to StartWithFee, sleep 100ms...\n", node.Name)
		time.Sleep(time.Millisecond * 100)
		count++
		if count > 400 {
			if node.ConditionQuit != nil {
				Logger.Printf("NODE %s %s StartWithFee with %s TIMEOUT\n", node.Address, node.Host, node.ConditionQuit.QuitEvent)
			} else {
				Logger.Printf("NODE %s %s StartWithFee TIMEOUT\n", node.Address, node.Host)
			}
			panic("Start photon node TIMEOUT")
		}
	}
	used := time.Since(t)
	if node.DebugCrash {
		Logger.Printf("NODE %s %s StartWithFee with %s in %fs", node.Address, node.Host, node.ConditionQuit.QuitEvent, used.Seconds())
	} else {
		Logger.Printf("NODE %s %s StartWithFee in %fs", node.Address, node.Host, used.Seconds())
	}
	time.Sleep(10 * time.Second)
	node.Running = true
}

// StartWithoutUpdateMeshNetworkNodes : Start start a photon node
func (node *PhotonNode) StartWithoutUpdateMeshNetworkNodes(env *TestEnv, pprof ...bool) {
	logfile := fmt.Sprintf("./log/%s.log", env.CaseName+"-"+node.Name)
	dopprof := false
	if len(pprof) > 0 {
		dopprof = pprof[0]
	}
	go ExecShell(env.Main, node.getParamStr(env, dopprof, false), logfile, true)

	count := 0
	t := time.Now()
	for !node.IsRunning() {
		Logger.Printf("waiting for %s to start, sleep 100ms...\n", node.Name)
		time.Sleep(time.Millisecond * 100)
		count++
		if count > 400 {
			if node.ConditionQuit != nil {
				Logger.Printf("NODE %s %s start with %s TIMEOUT\n", node.Address, node.Host, node.ConditionQuit.QuitEvent)
			} else {
				Logger.Printf("NODE %s %s start TIMEOUT\n", node.Address, node.Host)
			}
			panic("Start photon node TIMEOUT")
		}
	}
	used := time.Since(t)
	if node.DebugCrash {
		Logger.Printf("NODE %s %s start with %s in %fs", node.Address, node.Host, node.ConditionQuit.QuitEvent, used.Seconds())
	} else {
		Logger.Printf("NODE %s %s start in %fs", node.Address, node.Host, used.Seconds())
	}
	node.Running = true
}

// ReStartWithoutConditionQuitAndUpdateMethNetworkNodes : Restart start a photon node
func (node *PhotonNode) ReStartWithoutConditionQuitAndUpdateMethNetworkNodes(env *TestEnv) {
	node.DebugCrash = false
	node.ConditionQuit = nil
	node.Name = node.Name + "Restart"
	node.StartWithoutUpdateMeshNetworkNodes(env)
}

// ReStartWithoutConditionquit : Restart start a photon node
func (node *PhotonNode) ReStartWithoutConditionquit(env *TestEnv) {
	node.DebugCrash = false
	node.ConditionQuit = nil
	node.Name = node.Name + "Restart"
	node.Start(env, false)
}

// ReStartWithoutConditionquitAndNetwork : Restart start a photon node
func (node *PhotonNode) ReStartWithoutConditionquitAndNetwork(env *TestEnv) {
	node.DebugCrash = false
	node.ConditionQuit = nil
	node.Name = node.Name + "Restart"
	node.Start(env, false, true)
}

func (node *PhotonNode) getParamStr(env *TestEnv, pprof bool, nonetwork bool) []string {
	var param []string
	param = append(param, "--datadir="+env.DataDir)
	param = append(param, "--api-address="+node.APIAddress)
	param = append(param, "--listen-address="+node.ListenAddress)
	param = append(param, "--address="+node.Address)
	param = append(param, "--keystore-path="+env.KeystorePath)
	param = append(param, "--registry-contract-address="+env.TokenNetworkAddress)
	param = append(param, "--password-file="+env.PasswordFile)
	param = append(param, "--disable-fee")
	param = append(param, "--debug-mdns-interval=50ms")
	param = append(param, "--debug-mdns-keepalive=1s")
	param = append(param, "--debug-mdns-servicetag="+env.MDNSServiceTag)
	if pprof {
		param = append(param, "--pprof")
	}
	if nonetwork {
		param = append(param, "--debug-nonetwork")
	}

	if !env.UseMatrix {
		param = append(param, "--debug-udp-only")
		if env.XMPPServer != "" {
			param = append(param, "--xmpp-server="+env.XMPPServer)
		}
	} else {
		param = append(param, "--matrix")
		if time.Now().Nanosecond()%2 == 0 {
			param = append(param, "--matrix-server=transport01.smartmesh.cn")
		} else {
			param = append(param, "--matrix-server=transport13.smartmesh.cn")
		}
	}
	param = append(param, "--eth-rpc-endpoint="+env.EthRPCEndpoint)
	param = append(param, fmt.Sprintf("--verbosity=%d", env.Verbosity))
	param = append(param, "--debug")
	if node.DebugCrash == true {
		buf, err := json.Marshal(node.ConditionQuit)
		if err != nil {
			panic(err)
		}
		param = append(param, "--debugcrash")
		param = append(param, "--conditionquit="+string(buf))
	}
	return param
}

// StartWithConditionQuit start a photon node whit condition quit
func (node *PhotonNode) StartWithConditionQuit(env *TestEnv, c *params.ConditionQuit) {
	node.ConditionQuit = c
	node.DebugCrash = true
	node.Start(env)
}

// GetAddress :
func (node *PhotonNode) GetAddress() common.Address {
	return common.HexToAddress(node.Address)
}

// ClearHistoryData :
func (node *PhotonNode) ClearHistoryData(dataDir string) {
	if dataDir == "" {
		return
	}
	userDbPath := strings.ToLower(node.Address[2:10])
	err := filepath.Walk(dataDir, func(path string, fi os.FileInfo, err error) error {
		if nil == fi {
			return err
		}
		if !fi.IsDir() {
			return nil
		}
		name := fi.Name()

		if name == userDbPath {
			err := os.RemoveAll(path)
			if err != nil {
				fmt.Println("delete dir error:", err)
			}
			Logger.Printf("Clear history data of node %s %s SUCCESS", node.Name, node.Address)
		}
		return nil
	})
	err = filepath.Walk(".", func(path string, fi os.FileInfo, err error) error {
		if nil == fi {
			return err
		}
		if fi.IsDir() {
			return nil
		}
		name := fi.Name()
		if name == userDbPath {
			err := os.RemoveAll(path)
			if err != nil {
				fmt.Println("delete dir error:", err)
			}
			Logger.Println("Clear pfs history data SUCCESS ")
		}
		return nil
	})
	if err != nil {
		Logger.Println("No history data ")
	}
}

// ExecShell : run shell commands
func ExecShell(cmdstr string, param []string, logfile string, canquit bool) bool {
	var err error
	/* #nosec */
	cmd := exec.Command(cmdstr, param...)

	stdout, err := cmd.StdoutPipe()
	stderr, err := cmd.StderrPipe()

	err = cmd.Start()
	if err != nil {
		log.Println(err)
		return false
	}

	reader := bufio.NewReader(stdout)
	readererr := bufio.NewReader(stderr)

	logFile, err := os.Create(logfile)
	defer logFile.Close()
	if err != nil {
		log.Fatalln("Create log file error !", logfile)
	}

	debugLog := log.New(logFile, "", 0)
	//A real-time loop reads a line in the output stream.
	go func() {
		for {
			line, readErr := reader.ReadString('\n')
			if readErr != nil || io.EOF == readErr {
				break
			}
			//log.Println(line)
			debugLog.Println(line)
		}
	}()

	//go func() {
	for {
		line, readErr := readererr.ReadString('\n')
		if readErr != nil || io.EOF == readErr {
			break
		}
		//log.Println(line)
		debugLog.Println(line)
	}
	//}()

	err = cmd.Wait()

	if !canquit {
		log.Println("cmd ", cmdstr, " exited:", param)
	}

	if err != nil {
		//log.Println(err)
		debugLog.Println(err)
		if !canquit {
			os.Exit(-1)
		}
		return false
	}
	if !canquit {
		os.Exit(-1)
	}
	return true
}
