//Communal func and date structure
package main

import (
	"log"
	"math/big"
	"time"
	//"os"
	"os"
)

const MaxTry int = 20

var HalfLife = true

//For QueryingNodeAddress API  Response http body
type NodeAddressPayload struct {
	OurAddress string `json:"our_address"`
}

//For QueryingNodeAllChannels  and QueryingNodeSpecificChannel  API  Response http body
type NodeChannel struct {
	ChannelAddress string `json:"channel_address"`
	PartnerAddress string `json:"partner_address"`
	Balance        int32  `json:"balance"`
	PatnerBalance  int32  `json:"patner_balance"`
	TokenAddress   string `json:"token_address"`
	State          string `json:"state"`
	SettleTimeout  int32  `json:"settle_timeout"`
	RevealTimeout  int32  `json:"reveal_timeout"`
}

//New QueryingNodeAllChannels  and  QueryingNodeSpecificChannel API  Response http body
type newNodeChannel struct {
	ChannelAddress string `json:"channel_address"`
	PartnerAddress string `json:"partner_address"`
	TokenAddress   string `json:"token_address"`
	Balance        int32  `json:"balance"`
	State          string `json:"state"`
	SettleTimeout  int32  `json:"settle_timeout"`
}

//For QueryingAllPartnersForaToken  API  Response http body
type TokenPartnerPayload struct {
	PartnerAddress string `json:"partner_address"`
	Channel        string `json:"channel"`
}

//For OpenChannel API  http body
type OpenChannelPayload struct {
	PartnerAddress string `json:"partner_address"`
	TokenAddress   string `json:"token_address"`
	Balance        int32  `json:"balance"`
	SettleTimeout  int32  `json:"settle_timeout"`
}

//For CloseChannel API  http body
type CloseChannelPayload struct {
	State string `json:"state"`
}

//For SettleChannel API  http body
type SettleChannelPayload struct {
	State string `json:"state"`
}

//For  InitiatingTransfer API  http body
type TransferRequest struct {
	Amount     int32 `json:"amount"`
	Identifier int64 `json:"identifier"`
}

//For  InitiatingTransfer API  Response http body
type TransferResponse struct {
	InitiatorAddress string `json:"initiator_address"`
	TargetAddress    string `json:"target_address"`
	TokenAddress     string `json:"token_address"`
	Amount           int32  `json:"amount"`
	Identifier       int64  `json:"identifier"`
}

//For TokenSwaps API  htpp body
type TokenSwapsPayload struct {
	Role            string `json:"role"`
	SendingAmount   int32  `json:"sending_amount"`
	SendingToken    string `json:"sending_token"`
	ReceivingAmount int32  `json:"receiving_amount"`
	ReceivingToken  string `json:"receiving_token"`
}

//Desposit2Channel API htpp body
type Desposit2ChannelPayload struct {
	Balance int32 `json:"balance"`
}

//Connecting2TokenNetwork API htpp body
type Connecting2TokenNetworkPayload struct {
	Funds int32 `json:"funds"`
}

//LeavingTokenNetwork API htpp body
type LeavingTokenNetworkPayload struct {
	OnlyReceivingChannels bool `json:"only_receiving_channels"`
}

//QueryingConnectionsDetails API  Response http body
type ConnectionsDetails struct {
	Funds       *big.Int `json:"funds"`
	SumDeposits *big.Int `json:"sum_deposits"`
	Channels    int      `json:"channels"`
}

//Querying general network events API  Response http body
type GeneralNetworkEvents struct {
	EventType             string `json:"event_type"`
	TokenAddress          string `json:"token_address"`
	ChannelManagerAddress string `json:"channel_manager_address"`
}

//Querying token network events API  Response http body
type TokenNetworkEvents struct {
	EventType             string `json:"event_type"`
	SettleTimeout         int32  `json:"settle_timeout`
	TokenAddress          string `json:"token_address"`
	ChannelManagerAddress string `json:"channel_manager_address"`
}

//Querying Channel Events API  Response http body

type ChannelNewBalance struct {
	EventType   string `json:"event_type"`
	participant string `json:"participant`
	Balance     int32  `json:"balance"`
	BlockNumber int64  `json:"block_number"`
}
type TransferUpdated struct {
	EventType             string `json:"event_type"`
	TokenAddress          string `json:"token_address"`
	ChannelManagerAddress string `json:"channel_manager_address"`
}
type EventTransferSentSuccess struct {
	EventType   string `json:"event_type"`
	Identifier  int64  `json:"identifier"`
	BlockNumber int64  `json:"block_number"`
	Amount      int32  `json:"amount"`
	Target      string `json:"target"`
}

//Get time string
func GetTime() string {
	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)
	return tm.Format("2006-01-02 03:04:05 PM")
}

//Print time
func ShowTime() {
	log.Println(GetTime())
}

//Show Error Msg
func ShowError(err error) {
	if err != nil {
		//log.SetFlags(log.Lshortfile | log.LstdFlags)
		//log.Println(err)
		log.Output(3, err.Error())
		os.Exit(-1)
	}
}
