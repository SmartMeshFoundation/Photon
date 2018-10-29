package models

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/SmartMeshFoundation/Photon/network/rpc/contracts"
)

// PhotonEnvReader : save all data about photon nodes and refresh in time
type PhotonEnvReader struct {
	RegisterContractAddress string        `json:"register_contract_address"`
	PhotonNodes             []*PhotonNode `json:"photon_nodes"` // 节点列表
	Tokens                  []*Token      `json:"tokens"`       // Token列表
}

// NewPhotonEnvReader : construct
func NewPhotonEnvReader(hosts []string) *PhotonEnvReader {
	var env = new(PhotonEnvReader)
	// init hosts
	if hosts == nil || len(hosts) == 0 {
		panic("At least need one photon node")
	}
	for _, host := range hosts {
		env.PhotonNodes = append(env.PhotonNodes, &PhotonNode{
			Host: host,
		})
	}
	env.Refresh()
	return env
}

// Refresh : refresh all data by photon query api
func (env *PhotonEnvReader) Refresh() {
	// 1. refresh node address
	env.RefreshNodes()
	// 2. refresh tokens
	env.RefreshTokens()
	// 3. refresh channels
	env.RefreshChannels()
}

// RefreshNodes :
func (env *PhotonEnvReader) RefreshNodes() {
	for _, node := range env.PhotonNodes {
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
		err = json.Unmarshal(body, &addr)
		if err != nil {
			panic(err)
		}
		node.AccountAddress = addr.OurAddress
	}
	log.Println("PhotonEnvReader refresh nodes done")
}

// RefreshTokens :
func (env *PhotonEnvReader) RefreshTokens() {
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
	err = json.Unmarshal(body, &tokenAddrs)
	if err != nil {
		panic(err)
	}
	if len(tokenAddrs) == 0 {
		panic("no token")
	}
	env.Tokens = []*Token{}
	for _, addr := range tokenAddrs {
		if env.HasToken(addr) {
			continue
		}
		env.Tokens = append(env.Tokens, &Token{
			Address:      addr,
			IsRegistered: true,
		})
	}
	log.Println("PhotonEnvReader refresh tokens done")
}

// RefreshChannels :
func (env *PhotonEnvReader) RefreshChannels() {
	// clear old data
	for _, token := range env.Tokens {
		token.Channels = []Channel{}
	}
	// set new data
	for _, node := range env.PhotonNodes {
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
		err = json.Unmarshal(body, &nodeChannels)
		if err != nil {
			log.Println(err)
		}
		if len(nodeChannels) == 0 {
			continue
		}
		for _, channel := range nodeChannels {
			channel.SelfAddress = node.AccountAddress
			for _, token := range env.Tokens {
				if channel.TokenAddress == token.Address && !token.hasChannel(channel.ChannelIdentifier) {
					token.Channels = append(token.Channels, channel)
					break
				}
			}
		}
	}
	log.Println("PhotonEnvReader refresh channels done")
}

// HasToken ：
func (env *PhotonEnvReader) HasToken(tokenAddress string) bool {
	for _, token := range env.Tokens {
		if token.Address == tokenAddress {
			return true
		}
	}
	return false
}

// SaveToFile : save all data to file
func (env *PhotonEnvReader) SaveToFile(filepath string) {
	dataFile, err := os.Create(filepath)
	defer dataFile.Close()
	if err != nil {
		log.Fatalln("Create " + filepath + " file error !")
	}
	data, err := json.MarshalIndent(env, "", "\t")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = dataFile.Write(data)
	if err != nil {
		panic(err)
	}
	log.Println("Write env data to " + filepath + " done")
}

// RandomNode : get a random photon node
func (env *PhotonEnvReader) RandomNode() *PhotonNode {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	num := len(env.PhotonNodes)
	if num == 0 {
		return nil
	}
	return env.PhotonNodes[r.Intn(num)]
}

// RandomToken : get a random Token
func (env *PhotonEnvReader) RandomToken() *Token {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	num := len(env.Tokens)
	if num == 0 {
		return nil
	}
	return env.Tokens[r.Intn(num)]
}

// GetChannelsOfNode : get all channels of a photon node
func (env *PhotonEnvReader) GetChannelsOfNode(nodeAccountAddress string) (channels []Channel) {
	for _, token := range env.Tokens {
		for _, channel := range token.Channels {
			if channel.SelfAddress == nodeAccountAddress {
				channels = append(channels, channel)
			}
			if channel.PartnerAddress == nodeAccountAddress {
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

// GetChannelsOfNodeByState get all channels of a photon node by channel state
func (env *PhotonEnvReader) GetChannelsOfNodeByState(nodeAccountAddress string, state int) (channels []Channel) {
	all := env.GetChannelsOfNode(nodeAccountAddress)
	for _, channel := range all {
		if channel.State == state {
			channels = append(channels, channel)
		}
	}
	return channels
}

// GetChannelsByState : get all channels by channel state
func (env *PhotonEnvReader) GetChannelsByState(state int) (channels []Channel) {
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
func (env *PhotonEnvReader) GetNodeByAccountAddress(accountAddress string) (node *PhotonNode) {
	for _, n := range env.PhotonNodes {
		if n.AccountAddress == accountAddress {
			node = n
		}
	}
	return node
}

// HasOpenedChannelBetween :
func (env *PhotonEnvReader) HasOpenedChannelBetween(node1 *PhotonNode, node2 *PhotonNode, token *Token) bool {
	for _, channel := range token.Channels {
		if channel.State == contracts.ChannelStateOpened &&
			((channel.SelfAddress == node1.AccountAddress && channel.PartnerAddress == node2.AccountAddress) ||
				(channel.PartnerAddress == node1.AccountAddress && channel.SelfAddress == node2.AccountAddress)) {
			return true
		}
	}
	return false
}
