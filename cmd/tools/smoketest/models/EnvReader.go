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
	// 1. refresh node address
	env.RefreshNodes()
	// 2. refresh tokens
	env.RefreshTokens()
	// 3. refresh channels
	env.RefreshChannels()
}

// RefreshNodes :
func (env *RaidenEnvReader) RefreshNodes() {
	for _, node := range env.RaidenNodes {
		req := &Req{
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
	log.Println("RaidenEnvReader refresh nodes done")
}

// RefreshTokens :
func (env *RaidenEnvReader) RefreshTokens() {
	req := &Req{
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
		if env.HasToken(addr) {
			continue
		}
		env.Tokens = append(env.Tokens, &Token{
			Address:      strings.ToUpper(addr),
			IsRegistered: true,
		})
	}
	log.Println("RaidenEnvReader refresh tokens done")
}

// RefreshChannels :
func (env *RaidenEnvReader) RefreshChannels() {
	// clear old data
	for _, token := range env.Tokens {
		token.Channels = []Channel{}
	}
	// set new data
	for _, node := range env.RaidenNodes {
		req := &Req{
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
		var nodeChannels []Channel
		json.Unmarshal(body, &nodeChannels)
		if len(nodeChannels) == 0 {
			continue
		}
		for _, channel := range nodeChannels {
			channel.ChannelAddress = strings.ToUpper(channel.ChannelAddress)
			channel.SelfAddress = node.AccountAddress
			channel.TokenAddress = strings.ToUpper(channel.TokenAddress)
			channel.PartnerAddress = strings.ToUpper(channel.PartnerAddress)
			for _, token := range env.Tokens {
				if channel.TokenAddress == token.Address && !token.hasChannel(channel.ChannelAddress) {
					token.Channels = append(token.Channels, channel)
					break
				}
			}
		}
	}
	log.Println("RaidenEnvReader refresh channels done")
}

// HasToken ：
func (env *RaidenEnvReader) HasToken(tokenAddress string) bool {
	for _, token := range env.Tokens {
		if token.Address == strings.ToUpper(tokenAddress) {
			return true
		}
	}
	return false
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
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	num := len(env.RaidenNodes)
	if num == 0 {
		return nil
	}
	return env.RaidenNodes[r.Intn(num)]
}

// RandomToken : get a random TokenNetworkAddres
func (env *RaidenEnvReader) RandomToken() *Token {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	num := len(env.Tokens)
	if num == 0 {
		return nil
	}
	return env.Tokens[r.Intn(num)]
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

// GetChannelsOfNodeByState get all channels of a smartraiden node by channel state
func (env *RaidenEnvReader) GetChannelsOfNodeByState(nodeAccountAddress string, state string) (channels []Channel) {
	all := env.GetChannelsOfNode(nodeAccountAddress)
	for _, channel := range all {
		if channel.State == state {
			channels = append(channels, channel)
		}
	}
	return channels
}

// GetChannelsByState : get all channels by channel state
func (env *RaidenEnvReader) GetChannelsByState(state string) (channels []Channel) {
	for _, token := range env.Tokens {
		for _, channel := range token.Channels {
			if channel.State == state {
				channels = append(channels, channel)
			}
		}
	}
	return channels
}

// GetNodeByAccountAddress :
func (env *RaidenEnvReader) GetNodeByAccountAddress(accountAddress string) (node *RaidenNode) {
	for _, n := range env.RaidenNodes {
		if n.AccountAddress == accountAddress {
			node = n
		}
	}
	return node
}

// HasOpenedChannelBetween :
func (env *RaidenEnvReader) HasOpenedChannelBetween(node1 *RaidenNode, node2 *RaidenNode, token *Token) bool {
	for _, channel := range token.Channels {
		if channel.State == "opened" &&
			((channel.SelfAddress == node1.AccountAddress && channel.PartnerAddress == node2.AccountAddress) ||
				(channel.PartnerAddress == node1.AccountAddress && channel.SelfAddress == node2.AccountAddress)) {
			return true
		}
	}
	return false
}
