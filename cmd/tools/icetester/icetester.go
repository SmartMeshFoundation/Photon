package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"flag"
	"fmt"
	"os"

	"github.com/SmartMeshFoundation/SmartRaiden/network"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
)

const key1 = "9be657694d9b460fad0b00d798ecd4b9a1b2a034bc2e6567262d1f426e8d31c2"
const key2 = "9fbc4420d0df20e8830b86052cf8775d2bb19d0f89f7698a6df5aefa313c7d91"

var use2 bool

func init() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlTrace, utils.MyStreamHandler(os.Stderr)))
	flag.BoolVar(&use2, "2", false, "use second address")
	flag.Parse()
}

type testreceiver struct{}

func (testreceiver) Receive(data []byte, host string, port int) {
	log.Error(fmt.Sprintf("from %s:%d, recevied:%s", host, port, string(data)))
}

func main() {
	var err error
	var mykey *ecdsa.PrivateKey
	var myaddr, partneraddr common.Address
	key1bin, _ := hex.DecodeString(key1)
	key2bin, _ := hex.DecodeString(key2)
	privkey1, _ := crypto.ToECDSA(key1bin)
	privkey2, _ := crypto.ToECDSA(key2bin)
	addr1 := crypto.PubkeyToAddress(privkey1.PublicKey)
	addr2 := crypto.PubkeyToAddress(privkey2.PublicKey)
	log.Info(fmt.Sprintf("addr1=%s\naddr2=%s\n", addr1.String(), addr2.String()))
	network.InitIceTransporter("182.254.155.208:3478", "bai", "bai", "139.199.6.114:5222")
	log.Info(fmt.Sprintf("use2=%v", use2))
	if use2 {
		mykey = privkey2
		myaddr = addr2
		partneraddr = addr1
	} else {
		mykey = privkey1
		myaddr = addr1
		partneraddr = addr2
	}
	log.Info(fmt.Sprintf("myaddr=%s,partneraddr=%s\n", myaddr.String(), partneraddr.String()))
	it := network.NewIceTransporter(mykey, "client1")
	it.Start()
	it.Register(new(testreceiver))
	for {
		var cmd string
		fmt.Printf("input s to start, q to quit")
		fmt.Scanf("%s", &cmd)
		if cmd == "s" {
			err = it.Send(partneraddr, "", 0, []byte("aaaaaaa"))
			if err != nil {
				log.Error(fmt.Sprintf("send message error: %s", err))
				return
			}
		} else if cmd == "q" {
			it.Stop()
			return
		} else {
			log.Info("command error")
		}
	}
}
