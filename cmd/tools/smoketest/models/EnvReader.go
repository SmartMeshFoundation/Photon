package models

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

// RaidenEnvReader : save all data about raiden nodes and refresh in time
type RaidenEnvReader struct {
	RegisterContractAddress string        `json:"register_contract_address"`
	RaidenNodes             []*RaidenNode `json:"raiden_nodes"` // 节点列表
	Tokens                  []*Token      `json:"tokens"`       // Token列表
}

// NewRaidenEnvReader : construct
func NewRaidenEnvReader(hosts []string) *RaidenEnvReader {
	var env = new(RaidenEnvReader)
	// init hosts
	if hosts == nil || len(hosts) == 0 {
		panic("At least need one raiden node")
	}
	for _, host := range hosts {
		env.RaidenNodes = append(env.RaidenNodes, &RaidenNode{
			Host: host,
		})
	}
	env.Refresh()
	return env
}

// Refresh : refresh all data by raiden query api
func (env *RaidenEnvReader) Refresh() {
	var req *Req
	// 1. refresh node address
	for _, node := range env.RaidenNodes {
		req = &Req{
			APIName: "QueryNodeAddress",
			FullURL: node.Host + "/api/1/address",
			Method:  http.MethodGet,
			Payload: "",
			Timeout: time.Second * 30,
		}
		_, body, err := req.Invoke()
		if err != nil {
			panic(err)
		}
		var addr struct {
			OurAddress string `json:"our_address"`
		}
		json.Unmarshal(body, &addr)
		node.AccountAddress = strings.ToUpper(addr.OurAddress)
	}
	// 2. refresh tokens
	req = &Req{
		APIName: "QueryRegisteredTokens",
		FullURL: env.RandomNode().Host + "/api/1/tokens",
		Method:  http.MethodGet,
		Payload: "",
		Timeout: time.Second * 30,
	}
	_, body, err := req.Invoke()
	if err != nil {
		panic(err)
	}
	var tokenAddrs []string
	json.Unmarshal(body, &tokenAddrs)
	env.Tokens = []*Token{}
	for _, addr := range tokenAddrs {
		env.Tokens = append(env.Tokens, &Token{
			Address:      strings.ToUpper(addr),
			IsRegistered: true,
		})
	}
	// 3. refresh channels
	for _, node := range env.RaidenNodes {
		req = &Req{
			APIName: "QueryNodeAllChannels",
			FullURL: node.Host + "/api/1/channels",
			Method:  http.MethodGet,
			Payload: "",
			Timeout: time.Second * 30,
		}
		_, body, err := req.Invoke()
		if err != nil {
			panic(err)
		}
		var channels []Channel
		json.Unmarshal(body, &channels)
		if len(channels) == 0 {
			continue
		}
		// clear old data
		for _, token := range env.Tokens {
			token.Channels = []Channel{}
		}
		// set new data
		for _, channel := range channels {
			channel.SelfAddress = node.AccountAddress
			channel.TokenAddress = strings.ToUpper(channel.TokenAddress)
			for _, token := range env.Tokens {
				if channel.TokenAddress == token.Address && !token.hasChannel(channel.ChannelAddress) {
					token.Channels = append(token.Channels, channel)
					break
				}
			}
		}
	}
	log.Println("RaidenEnvReader refresh done")
}

// SaveToFile : save all data to file
func (env *RaidenEnvReader) SaveToFile(filepath string) {
	dataFile, err := os.Create(filepath)
	defer dataFile.Close()
	if err != nil {
		log.Fatalln("Create " + filepath + " file error !")
	}
	data, err := json.MarshalIndent(env, "", "\t")
	if err != nil {
		log.Fatalln(err)
	}
	dataFile.Write(data)
	log.Println("Write env data to " + filepath + " done")
}

// RandomNode : get a random raiden node
func (env *RaidenEnvReader) RandomNode() *RaidenNode {
	return env.RaidenNodes[rand.Intn(len(env.RaidenNodes))]
}

// RandomToken : get a random Token
func (env *RaidenEnvReader) RandomToken() *Token {
	return env.Tokens[rand.Intn(len(env.Tokens))]
}

// GetChannelsOfNode : get all channels of a smartraiden node
func (env *RaidenEnvReader) GetChannelsOfNode(nodeAccountAddress string) (channels []Channel) {
	for _, token := range env.Tokens {
		for _, channel := range token.Channels {
			if channel.SelfAddress == strings.ToUpper(nodeAccountAddress) {
				channels = append(channels, channel)
			}
			if strings.ToUpper(channel.PartnerAddress) == strings.ToUpper(nodeAccountAddress) {
				// deep copy
				t := channel
				new := &t
				// revert
				new.SelfAddress, new.PartnerAddress = channel.PartnerAddress, channel.SelfAddress
				new.Balance, new.PartnerBalance = channel.PartnerBalance, channel.Balance
				new.LockedAmount, new.PartnerLockedAmount = channel.PartnerLockedAmount, channel.LockedAmount
				channels = append(channels, *new)
			}
		}
	}
	return channels
}

// GetChannelsOfNodeByState get all open channels of a smartraiden node
func (env *RaidenEnvReader) GetChannelsOfNodeByState(nodeAccountAddress string, state string) (channels []Channel) {
	all := env.GetChannelsOfNode(nodeAccountAddress)
	for _, channel := range all {
		if channel.State == state {
			channels = append(channels, channel)
		}
	}
	return channels
}
