package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"strconv"

	"github.com/huamou/config"
	"github.com/kataras/iris/utils"
)

func Exec_shell(cmdstr string, param []string, logfile string, canquit bool) bool {
	if !utils.Exists(cmdstr) {
		log.Println(cmdstr+" is not exist", param)
		os.Exit(-1)
		return false
	}
	cmd := exec.Command(cmdstr, param...)

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	err := cmd.Start()

	if err != nil {
		log.Println(err)
		return false
	}

	reader := bufio.NewReader(stdout)
	readererr := bufio.NewReader(stderr)

	logPath := filepath.Dir(logfile)
	if !utils.Exists(logPath) {
		os.Mkdir(logPath, 0777)
	}

	logFile, err := os.Create(logfile)
	defer logFile.Close()
	if err != nil {
		log.Fatalln("Create log file error !", logfile)
	}

	debugLog := log.New(logFile, "[Debug]", log.Llongfile)
	//本地注释：实时循环读取输出流中的一行内容
	//A real-time loop reads a line in the output stream.
	go func() {
		for {
			line, err := reader.ReadString('\n')
			if err != nil || io.EOF == err {
				break
			}
			//log.Println(line)
			debugLog.Println(line)
		}
	}()

	//go func() {
	for {
		line, err := readererr.ReadString('\n')
		if err != nil || io.EOF == err {
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

func Startraiden(RegistryAddress string) {
	var pstr []string
	//本地注释：公共参数
	//public parameter
	var pstr2 []string
	//本地注释：杀死旧进程
	pstr2 = append(pstr2, "goraiden")
	Exec_shell("/usr/bin/killall", pstr2, "./../../testdata/log/ka.log", true)
	//本地注释：杀死旧进程后等待释放端口
	time.Sleep(10 * time.Second)

	param := new(RaidenParam)
	c, err := config.ReadDefault("./ApiTest.INI")

	if err != nil {
		log.Println("Read error:", err)
		return
	}

	param.datadir = c.RdString("common", "datadir", "/smtwork/share/.goraiden")
	param.keystore_path = c.RdString("common", "keystore_path", "/smtwork/privnet3/data/keystore")
	param.discovery_contract_address = c.RdString("common", "discovery_contract_address", "0x5f014DA6ea514405f641341e42aC0e61B8190653")
	if RegistryAddress == "" {
		param.registry_contract_address = c.RdString("common", "registry_contract_address", "0x069E5c8954b14c7638e8E6479402FDa6F9971036")

	} else {
		param.registry_contract_address = RegistryAddress
	}

	param.password_file = c.RdString("common", "password_file", "")
	param.nat = c.RdString("common", "nat", "none")
	param.eth_rpc_endpoint = c.RdString("common", "eth_rpc_endpoint", "ws://127.0.0.1:8546")
	param.conditionquit = c.RdString("common", "conditionquit", "{\"QuitEvent\":\"RefundTransferRecevieAckxx}")
	param.debug = c.RdBool("common", "debug", true)

	//本地注释：节点1
	//NODE 1
	var NODE string
	for i := 0; i < 6; i++ {
		NODE = "NODE" + strconv.Itoa(i+1)
		param.api_address = c.RdString(NODE, "api_address", "")
		param.listen_address = c.RdString(NODE, "listen_address", "")
		param.address = c.RdString(NODE, "address", "")
		pstr = param.getParam()
		//log.Println(pstr)
		logfile := c.RdString(NODE, "log", "")
		exepath := c.RdString(NODE, "raidenpath", "/project/bin/goraiden.exe")
		go Exec_shell(exepath, pstr, logfile, false)
	}

}
