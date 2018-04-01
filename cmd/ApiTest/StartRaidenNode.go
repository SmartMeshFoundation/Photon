package main

import (
	"github.com/larspensjo/config"
	"log"
	"os/exec"
	"time"
)

func Exec_shell(cmdstr string, param []string) bool {
	cmd := exec.Command(cmdstr, param...)
	//stdout, _ := cmd.StdoutPipe()
	//stderr, _ := cmd.StderrPipe()
	err := cmd.Start()

	if err != nil {
		log.Println(err)
		return false
	}
	//
	//reader := bufio.NewReader(stdout)
	//readererr := bufio.NewReader(stderr)
	//
	////本地注释：实时循环读取输出流中的一行内容
	////A real-time loop reads a line in the output stream.
	//go func() {
	//	for {
	//		line, err := reader.ReadString('\n')
	//		if err != nil || io.EOF == err {
	//			break
	//		}
	//		log.Println(line)
	//	}
	//}()
	//
	//for {
	//	line, err := readererr.ReadString('\n')
	//	if err != nil || io.EOF == err {
	//		break
	//	}
	//	log.Println(line)
	//}
	//
	//err = cmd.Wait()
	//if err != nil {
	//	log.Println(err)
	//	return false
	//}

	return true
}

func Startraiden(RegistryAddress string) {
	var pstr []string
	//本地注释：公共参数
	//public parameter
	var pstr2 []string
	//本地注释：杀死旧进程
	pstr2 = append(pstr2, "goraiden")
	Exec_shell("/usr/bin/killall", pstr2)
	//本地注释：杀死旧进程后等待释放端口
	time.Sleep(30 * time.Second)

	param := new(RaidenParam)
	c, err := config.ReadDefault("./ApiTest.INI")

	if err != nil {
		log.Println("Read error:", err)
		return
	}

	s, err := c.String("common", "datadir")
	if err != nil {
		log.Println("Read error:", err)
		return
	}
	param.datadir = s
	s, err = c.String("common", "keystore_path")
	if err != nil {
		log.Println("Read error:", err)
		return
	}
	param.keystore_path = s
	s, err = c.String("common", "discovery_contract_address")
	if err != nil {
		log.Println("Read error:", err)
		return
	}
	param.discovery_contract_address = s
	if RegistryAddress == "" {
		s, err = c.String("common", "registry_contract_address")
		if err != nil {
			log.Println("Read error:", err)
			return
		}
		param.registry_contract_address = s
	} else {
		param.registry_contract_address = RegistryAddress
	}
	s, err = c.String("common", "password_file")
	if err != nil {
		log.Println("Read error:", err)
		return
	}
	param.password_file = s
	s, err = c.String("common", "nat")
	if err != nil {
		log.Println("Read error:", err)
		return
	}
	param.nat = s
	s, err = c.String("common", "eth_rpc_endpoint")
	if err != nil {
		log.Println("Read error:", err)
		return
	}
	param.eth_rpc_endpoint = s
	s, err = c.String("common", "conditionquit")
	if err != nil {
		log.Println("Read error:", err)
		return
	}
	param.conditionquit = s
	b, err := c.Bool("common", "debug")
	if err != nil {
		log.Println("Read error:", err)
		return
	}
	param.debug = b
	//本地注释：节点1
	//NODE 1
	s, err = c.String("NODE1", "api_address")
	if err != nil {
		log.Println("Read error:", err)
		return
	}
	param.api_address = s
	s, err = c.String("NODE1", "listen_address")
	if err != nil {
		log.Println("Read error:", err)
		return
	}
	param.listen_address = s
	s, err = c.String("NODE1", "address")
	if err != nil {
		log.Println("Read error:", err)
		return
	}
	param.address = s
	pstr = param.getParam()
	//log.Println(pstr)
	s, err = c.String("NODE1", "raidenpath")
	if err != nil {
		log.Println("Read error:", err)
		return
	}
	Exec_shell(s, pstr)
	//本地注释：节点2
	//NODE 2
	s, err = c.String("NODE2", "api_address")
	if err != nil {
		log.Println("Read error:", err)
		return
	}
	param.api_address = s
	s, err = c.String("NODE2", "listen_address")
	if err != nil {
		log.Println("Read error:", err)
		return
	}
	param.listen_address = s
	s, err = c.String("NODE2", "address")
	if err != nil {
		log.Println("Read error:", err)
		return
	}
	param.address = s
	pstr = param.getParam()
	//log.Println(pstr)
	s, err = c.String("NODE2", "raidenpath")
	if err != nil {
		log.Println("Read error:", err)
		return
	}
	Exec_shell(s, pstr)
	//本地注释：节点3
	//NODE 3
	s, err = c.String("NODE3", "api_address")
	if err != nil {
		log.Println("Read error:", err)
		return
	}
	param.api_address = s
	s, err = c.String("NODE3", "listen_address")
	if err != nil {
		log.Println("Read error:", err)
		return
	}
	param.listen_address = s
	s, err = c.String("NODE3", "address")
	if err != nil {
		log.Println("Read error:", err)
		return
	}
	param.address = s
	pstr = param.getParam()
	//log.Println(pstr)
	s, err = c.String("NODE3", "raidenpath")
	if err != nil {
		log.Println("Read error:", err)
		return
	}
	Exec_shell(s, pstr)
	//本地注释：节点4
	//NODE 4
	s, err = c.String("NODE4", "api_address")
	if err != nil {
		log.Println("Read error:", err)
		return
	}
	param.api_address = s
	s, err = c.String("NODE4", "listen_address")
	if err != nil {
		log.Println("Read error:", err)
		return
	}
	param.listen_address = s
	s, err = c.String("NODE4", "address")
	if err != nil {
		log.Println("Read error:", err)
		return
	}
	param.address = s
	pstr = param.getParam()
	//log.Println(pstr)
	s, err = c.String("NODE4", "raidenpath")
	if err != nil {
		log.Println("Read error:", err)
		return
	}
	Exec_shell(s, pstr)
	//本地注释：节点5
	//NODE 5
	s, err = c.String("NODE5", "api_address")
	if err != nil {
		log.Println("Read error:", err)
		return
	}
	param.api_address = s
	s, err = c.String("NODE5", "listen_address")
	if err != nil {
		log.Println("Read error:", err)
		return
	}
	param.listen_address = s
	s, err = c.String("NODE5", "address")
	if err != nil {
		log.Println("Read error:", err)
		return
	}
	param.address = s
	pstr = param.getParam()
	//log.Println(pstr)
	s, err = c.String("NODE5", "raidenpath")
	if err != nil {
		log.Println("Read error:", err)
		return
	}
	Exec_shell(s, pstr)
	//本地注释：节点6
	//NODE 6
	s, err = c.String("NODE6", "api_address")
	if err != nil {
		log.Println("Read error:", err)
		return
	}
	param.api_address = s
	s, err = c.String("NODE6", "listen_address")
	if err != nil {
		log.Println("Read error:", err)
		return
	}
	param.listen_address = s
	s, err = c.String("NODE6", "address")
	if err != nil {
		log.Println("Read error:", err)
		return
	}
	param.address = s
	pstr = param.getParam()
	//log.Println(pstr)
	s, err = c.String("NODE6", "raidenpath")
	if err != nil {
		log.Println("Read error:", err)
		return
	}
	Exec_shell(s, pstr)
}
