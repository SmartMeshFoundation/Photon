# Photon’s Mobile API Documentation
## Installation
Photon mobile SDk compilation must require the gomobile tool to work properly. Please refer to [gomobile](https://godoc.org/golang.org/x/mobile) for gomobile installation. 
```bash
cd mobile
#build android
./build_android.sh
#build iOS
./build_iOS.sh
```
### android use
Integrate mobile.aar into your project
### iOS use
Integrate Mobile.framework into your project.

### Other known issues
Due to the working restrictions of gomobile, if there are two gomobile compiled sdk in the project (for example, your project also depends on the mobile package of ethereum), the program will not run normally.

## Node Management Related Interface
Photon relies on gomobile to automate interface encapsulation. Because it is a cross-language call, it is unavoidable that there is a type conversion problem.
In order to avoid such problems, Photon provides interfaces to almost all basic types (int, string, error).

### Starting a photon node
func StartUp(address, keystorePath, ethRPCEndPoint, dataDir, passwordfile, apiAddr, listenAddr, logFile string, registryAddress string, otherArgs *Strings) (api *API, err error)

parameter:
* `address string`– the account address used by the photon node
* `keystorePath string` – the path of account private key 
* `ethRPCEndPoint string` – public chain node host, http protocol
* `dataDir string` – photon db path
* `passwordfile string` –  password file path
* `apiAddr string` – http api listening port
* `listenAddr string` – udp listening port
* `logFile string` – log file path
* `registryAddress string` – TokenNetworkRegistry contract address
* `otherArgs mobile.Strings` – other parameters, see photon -h   

If you need to pass other parameters than the default parameters, you can refer to the following ways:
```go
otherArgs := mobile.NewStrings(2)
err = otherArgs.Set(0, fmt.Sprintf("--registry-contract-address=%s", registryContractAddress))
if err != nil {
    return err
}
err = otherArgs.Set(1, fmt.Sprintf("--help"))
if err != nil {
    return err
}
```
return:
* `api *API` – startup successfully which will return the api handle
* `err error` – error message

### Stop a photon node
func (a *API) Stop()

### Switching the photon operating environment
func (a *API) SwitchNetwork(isMesh bool)

Switch network environment,Mesh or Internet
In the Mesh network, the nodes directly communicate using the UDP protocol, and the App needs to notify the Photon other nodes information through the UpdateMeshNetworkNodes.

### Notify the photon node network to disconnect
func (a *API) NotifyNetworkDown() error

Proactively inform the photon node that the network is disconnected and let the photon node start trying to reconnect

In the mobile phone network environment, due to network complexity, such as WiFi disconnection, these events Photon can not be directly perceived from the system, the App needs to actively tell Photon to take appropriate processing.

### Subscribe to photon events
func (a *API) Subscribe(handler NotifyHandler) (sub *Subscription, err error)

Subscribe to photon node events, including transaction notifications, error notifications, etc.When the photon changes internally, the notification is actively pushed to the App to avoid App polling and to improve efficiency.
```go
// NotifyHandler is a client-side subscription callback to invoke on events and
// subscription failure.
type NotifyHandler interface {
	//some unexpected error
	OnError(errCode int, failure string)
	//OnStatusChange server connection status change
	OnStatusChange(s string)
	//OnReceivedTransfer  receive a transfer
	OnReceivedTransfer(tr string)
	//OnSentTransfer a transfer sent success
	OnSentTransfer(tr string)
	// OnNotify get some important message photon want to notify upper application
	OnNotify(level int, info string)
}
```

#### OnError
  Notify Photon that an unrecoverable error has occurred. Any function of Photon must be restarted before it can be used. Since it is considered that the integration of Photon may be single-process, The unknown error may cause the App to quit, so even if there is an unpredictable error inside Photon, Photon will intercept and report to the App, and the App will decide whether to exit immediately or continue to use it.
- `errCode`is the error code
- `failure`is an error message description
Restart the Photon mode:
```go
api.Stop()
newAPI,err:=Startup(...)
```
#### OnStatusChange
This interface is used to notify the App that the status from Subscribe to the current public link and XMPP link has changed. If there is no change, it will not be notified.
`s` is the json code of the structure below.
```go
//ConnectionStatus status of network connection
{
 "xmpp_status":1,
 "eth_status":1,
 "last_block_time":"2019-01-23"
}
```
Where `XMPPStatus` and `EthStatus` are defined as follows:
```go
// Status shows actual connection status.
const (
	//Disconnected init status
	Disconnected = 0
	//Connected connection status
	Connected =1
	//Closed user closed
	Closed =2
	//Reconnecting connection error
	Reconnecting =3
)
```
#### OnReceivedTransfer
This interface is used to notify the app that a new transaction has been received.
`tr`is the json code of the following structure
```go
//ReceivedTransfer tokens I have received and where it comes from
type ReceivedTransfer struct {
	Key               string `storm:"id"`
	BlockNumber       int64  `json:"block_number" storm:"index"`
	OpenBlockNumber   int64
	ChannelIdentifier common.Hash    `json:"channel_identifier"`
	TokenAddress      common.Address `json:"token_address"`
	FromAddress       common.Address `json:"from_address"`
	Nonce             uint64         `json:"nonce"`
	Amount            *big.Int       `json:"amount"`
}
```
Note: This interface does not contain transactions that participate as intermediate nodes.
#### OnSentTransfer
This interface is used to notify the app that a transaction just sent out has succeeded.
`tr` is the json encoding of the following structure
```go
//SentTransfer transfer's I have sent and success.
type SentTransfer struct {
	Key               string `storm:"id"`
	BlockNumber       int64  `json:"block_number" storm:"index"`
	OpenBlockNumber   int64
	ChannelIdentifier common.Hash    `json:"channel_identifier"`
	ToAddress         common.Address `json:"to_address"`
	TokenAddress      common.Address `json:"token_address"`
	Nonce             uint64         `json:"nonce"`
	Amount            *big.Int       `json:"amount"`
}
```
Note: This interface does not contain transactions that participate as intermediate nodes.
#### OnNotify
The interface has two functions, mainly based on the type in the second parameter.
If the type is 0, it represents string information that photon hopes to  push to the user and let the user know  there is a change inside photon. At this time, the first parameter indicates the importance of the information.
If the type is 1, it means that the status of a transaction initiated by the user has changed. For details, please refer to[Query the status of the transaction initiated by yourself](#Query the status of the transaction initiated by yourself) 
`level` is defined as follows
```go
type Level int

const (
	// LevelInfo : 0
	LevelInfo = 0
	// LevelWarn : 1
	LevelWarn =1
	// LevelError : 2
	LevelError =2
)
```
Where `info` is 
```go
 type InfoStruct struct {
		Type    int
		Message interface{}
}
```
 
### Manually registering node information
func (a *API) UpdateMeshNetworkNodes(nodesstr string) (err error)

Manually register a communicable node address to photon. This interface is mainly used for information registration of nodes under no-network conditions. After registering the node information, the UDP mode communication will be prioritized, which can improve tps. In addition, in order to enable nodes to communicate (mutually recognize each other) in the no-network state, both nodes need to call the interface to update the other party's information,and if the registered node is restarted,the information need to re-registered since the message is stored in memory.

Example data:

```json
[{
   "address":"0x292650fee408320D888e06ed89D938294Ea42f99",
   "ip_port":"127.0.0.1:40001"
},
{
     "address":"0x4B89Bff01009928784eB7e7d10Bf773e6D166066",
    "ip_port":"127.0.0.1:40002"
}
]
```
Tell Photon how to communicate with 0x292650fee408320D888e06ed89D938294Ea42f99 and 0x4B89Bff01009928784eB7e7d10Bf773e6D166066

## Force network reconnection
func (a *API) OnResume() (err error) 

Called when the phone switches from the background to the foreground, allowing the photon node to start trying to reconnect.

## Query system status
func (a *API) GetSystemStatus() (r string, err error)

 Example Response：
```json
{
    "error_code": "0000",
    "error_msg": "SUCCESS",
    "data": {
        "eth_rpc_endpoint": "ws://192.168.124.13:5555",
        "eth_rpc_status": "connected",
        "node_address": "0x97Cd7291f93F9582Ddb8E9885bF7E77e3f34Be40",
        "registry_address": "0xf1d87c419a586Bd480Ce33067180F8e710B9931F",
        "token_to_token_network": {
            "0x2a7af974b7bb88703180d6aff9a656bb4dbba809": "0x0000000000000000000000000000000000000000",
            "0x489cee6beaa894898d0890f4c6d750ca3d8176a4": "0x0000000000000000000000000000000000000000",
            "0x8b916406c1eccc5b15865b7bc7af5fa90c01fc59": "0x0000000000000000000000000000000000000000",
            "0xc07d1d6e8f20f2a90b205762a0bac0b611c490dc": "0x0000000000000000000000000000000000000000"
        },
        "block_number": 15555306,
        "last_block_number_time": "2019-01-28T11:32:15.4144738+08:00",
        "is_mobile_mode": false,
        "network_type": "xmpp-udp",
        "fee_policy": {
            "Key": "",
            "account_fee": {
                "fee_constant": 0,
                "fee_percent": 10000,
                "signature": null
            },
            "token_fee_map": {},
            "channel_fee_map": {}
        },
        "channel_num": 3,
        "transfers": {
            "send_num": 0,
            "receive_num": 0,
            "dealing_num": 0
        }
    }
}
```

## Channel Structure   
```json
    {
        "channel_identifier": "0x47235d9d81eb6c19dea2b695b3d6ba1cf76c169d329dc60d188390ba5549d025",
        "open_block_number": 3158573,
        "partner_address": "0x31DdaC67e610c22d19E887fB1937BEE3079B56Cd",
        "balance": 100000000000000000000,
        "partner_balance": 100000000000000000000,
        "locked_amount": 0,
        "partner_locked_amount": 0,
        "token_address": "0xF2747ea1AEE15D23F3a49E37A146d3967e2Ea4E5",
        "state": 1,
        "StateString": "opened",
        "settle_timeout": 150,
        "reveal_timeout": 30
    }
```
Channel structure description ： 
- `channel_identifier`:  Address for a channel
- `open_block_number` : Block height when a channel opens
- `partner_address`: The address of the other participant of the channel
- `balance`: Available Balance of the channel participant
- `partner_balance` : Available Balance of the other participant of the channel 
- `locked_amount`: The locked amount of the participant
- `partner_locked_amount`: The locked amount of the other participant 
- `token_address`: Address for tokens in this channel
- `state` :The digits denoting for the channel states
- `StateString` :The string literal for the Channel States
-  `settle_timeout`: Some amount of block denoting time period for transaction settlement,which must greater than `reveal_timeout`.
-  `reveal_timeout`: The block height at which nodes registering `secret`,the default value is 30, and if modified, it can be setting at node startup with `-- reveal-timeout` 

State|StateString|Description
---|---|---
0 |inValid|Channels do not exist
1|opened|Channel open status,which can carry out normal offchain transactions
2|closed|The channel is closed, no more transactions can be initiated, but Ongoing transactions can be accepted.
3|settled|The channel is settled which  the token will return to the respective accounts on the blockchain, and the channel will be invalid.
4|closing| The participant initiated a request to close the channel, the transactions which is being processed  can continue to finish, but the participant  cannot initiate new transactions.
5|settling|The participant initiated a settlement request and is processing. Normally, there should be no uncompleted transaction and no new transactions can be initiated. The settling transaction are being submitted to the chain and have not yet been successfully packaged.
6|withdrawing|When the participant receives or sends a `withdraw` request,  just at this moment,he receive the transaction request of the other node,the ongoing transaction can only be abandoned immediately.
7|cooperativeSettling| Once the participant receives or sends the `cooperative settle`requests, just at this moment,he receive the transaction request of the other node, the ongoing transactions will  be abandoned immediately.
8|prepareForCooperativeSettle| The participant received  ` CooperativeSettle ` request,but there is ongoing transactions and the channel cannot be cooperatively settled. At this time , if the participant still want to cooperative settle the channel, he can wait until the transaction is completed. In order to prevent new transactions from occurring during the waiting period,  the 'prepareForCooperativeSettle' can be set as the mark to stop accepting new transactions and wait for the current transaction to be completed. Then he can call the CooperativeSettle to settle the channel. 
 9|prepareForWithdraw|The participant receives the request to initiate `withdraw`,but the participant or the partner still hold the locks,he cannot withdraw tokens from the channel. At this time , if the participant still want to withdraw tokens from the channel, he need to wait for the locked transaction to be unlocked. In order to prevent new transactions from occurring during the waiting period,  the 'prepareForWithdraw' can be set as the mark to stop accepting new transactions and wait for the current transaction to be unlocked. Then he can call the `withdraw` to withdraw the token from the channel. 
10|unkown|StateError

## Query interface
### Query the account address of the running photon node
func (a *API) Address() (addr string)

Example Response：
```json
[
"0x7B874444681F7AEF18D48f330a0Ba093d3d0fDD2"
]
```
### Query the list of all the registered token
func (a *API) Tokens() (tokens string)

Example Response：
```json
[
    "0x7B874444681F7AEF18D48f330a0Ba093d3d0fDD2"
]
```
### Query all the channels that you participate in under some token.
func (a *API) TokenPartners(tokenAddress string) (channels string, err error)

Example Response：
```json
[
    {
        "partner_address": "0x151E62a787d0d8d9EfFac182Eae06C559d1B68C2",
        "channel": "api/1/channles/0x79b789e88c3d2173af4048498f8c1ce66f019f33a6b8b06bedef51dde72bbbc1"
    },
    {
        "partner_address": "0x201B20123b3C489b47Fde27ce5b451a0fA55FD60",
        "channel": "api/1/channles/0xd971f803c7ea39ee050bf00ec9919269cf63ee5d0e968d5fe33a1a0f0004f73d"
    }
]
```
### Query all the channels of the node
func (a *API) GetChannelList() (channels string, err error)

Example Response：
```json
[
    {
        "channel_address": "0xc502076485a3cff65f83c00095dc55e745f790eee4c259ea963969a343fc792a",
        "open_block_number": 5228715,
        "partner_address": "0x4B89Bff01009928784eB7e7d10Bf773e6D166066",
        "balance": 499490000000000000000000,
        "partner_balance": 1500506000000000000000000,
        "locked_amount": 0,
        "partner_locked_amount": 0,
        "token_address": "0x663495a1b8e9Be17083b37924cFE39e17858F9e8",
        "state": 1,
        "StateString": "opened",
        "settle_timeout": 100000,
        "reveal_timeout": 5000
    }
]
```
### Query information about special channel
func (a *API) GetOneChannel(channelIdentifier string) (channel string, err error)

Example Response：
```json
{
    "channel_identifier": "0xc502076485a3cff65f83c00095dc55e745f790eee4c259ea963969a343fc792a",
    "open_block_number": 5228715,
    "partner_address": "0x4B89Bff01009928784eB7e7d10Bf773e6D166066",
    "balance": 499490000000000000000000,
    "patner_balance": 1500506000000000000000000,
    "locked_amount": 0,
    "partner_locked_amount": 0,
    "token_address": "0x663495a1b8e9Be17083b37924cFE39e17858F9e8",
    "state": 1,
    "StateString": "opened",
    "settle_timeout": 100000,
    "reveal_timeout": 0,
    "ClosedBlock": 0,
    "SettledBlock": 0,
    "OurUnkownSecretLocks": {},
    "OurKnownSecretLocks": {},
    "PartnerUnkownSecretLocks": {},
    "PartnerKnownSecretLocks": {},
    "OurLeaves": null,
    "PartnerLeaves": null,
    "OurBalanceProof": {
        "Nonce": 0,
        "TransferAmount": 0,
        "LocksRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "ChannelIdentifier": {
            "ChannelIdentifier": "0x0000000000000000000000000000000000000000000000000000000000000000",
            "OpenBlockNumber": 0
        },
        "MessageHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "Signature": null,
        "ContractTransferAmount": 0,
        "ContractNonce": 0,
        "ContractLocksRoot": "0x0000000000000000000000000000000000000000000000000000000000000000"
    },
    "PartnerBalanceProof": {
        "Nonce": 0,
        "TransferAmount": 0,
        "LocksRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "ChannelIdentifier": {
            "ChannelIdentifier": "0x0000000000000000000000000000000000000000000000000000000000000000",
            "OpenBlockNumber": 0
        },
        "MessageHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "Signature": null,
        "ContractTransferAmount": 0,
        "ContractNonce": 0,
        "ContractLocksRoot": "0x0000000000000000000000000000000000000000000000000000000000000000"
    },
    "Signature": null
}
```
##  Function interface of channel
### Deposit to the channel

func (a *API) Deposit(partnerAddress, tokenAddress string, settleTimeout int, balanceStr string, newChannel bool) (callID string, err error)

Deposit to the channel (if there is no channel, the interface can be reused to create the channel and deposit).

Parameter |type |JSON format|description
--|--|--|--
partnerAddress|string|partner_address|The address of the partner
tokenAddress|string|token_address| which token to deposit
settleTimeout|string|settle_timeout|The settlement window period 
balanceStr|big.Int|balance|The deposited amount which must be greater than 0.
newChannel|bool|new_channel|Judge whether the channel exists or not.If the channel doesnot exist, `deposit` will create a new channel and deposit, else only deposit.  

Example Request：
```json 
{
    "partnerAddress": "0x7d289f1cBd70d5c3c6F56c09f812F6407f6458B7",
    "tokenAddress": "0xadE88bC1519867e7091f83D763cf61918d50244a",
     "settleTimeout": 100,
     "balanceStr": 10000000000000000000000,
    "newChannel": true
}
```

deposit interfaces contain two behaviors：

1. Create channel and deposit
    - `new_channel`sets`true`, which means open a new channel and deposit to the channel;if there is no channel between the participants, `false`is no meaning for `new_channel`，which will return the error message "There is no channel".
    - `settle_timeout`represent the settlement window for new channel, for example,settle_timeout：100; if the `settle_timeout` set to 0,the default window period is used which is 600 block.

2. only deposit:
   - `new_channel`must set to `false`，which means the channel has been existed；If the channel has been existed ,there is Meaningless to set the `new_channel`statue as `ture`，which will response the error message "The channel has already existed". 
   - `settle_timeout`must set to 0,because the channel has already existed.


Example Response:
```json
{
    "channel_identifier": "0x16305a3a4e1b8f1ee167be895c60a9a77551ea1db40077a3a897cb1a75dadab1",
    "open_block_number": 1607480,
    "partner_address": "0x7d289f1cBd70d5c3c6F56c09f812F6407f6458B7",
    "balance": 10000000000000000000000,
    "partner_balance": 0,
    "locked_amount": 0,
    "partner_locked_amount": 0,
    "token_address": "0xadE88bC1519867e7091f83D763cf61918d50244a",
    "state": 1,
    "state_string": "opened",
    "settle_timeout": 100,
    "reveal_timeout": 30
}
```
**Status Codes:**  
- `200 OK`  
- `409 Conflict` 

Possible conflict situations:

If the channel exists and the parameter `new_channel` is set to true, an "Error" is prompted: "channel already exist";

If the channel does not exist and the parameter `new_channel` is set to false, an "Error" is prompted: "channel does not exist"

Setting `settle_timeout` to be non-zero when the channel already exists will prompt "settleTimeout must be zero when newChannel is false"

### Withdraw from the channel  
func (a *API) Withdraw(channelIdentifierHashStr, amountstr, op string) (callID string, err error) 

CooperateWithdraw available when both channel participants online.

parameter:
- channelIdentifierHashStr        Channel address
- amountstr                      Amount of taken which will be withdrawn
- op                             Option
  -  When you’re ready to withdraw, you can switch the channel state to `"preparewithdraw"` by setting the `"op"`:`"preparewithdraw"` and refuse to accept the transaction.
  -  When you want to cancel the state of the `preparewithdraw`, you can switch the channel state to the`opened` through the parameter`"op":"cancelprepare"`.

Parameter:    
```json
{
    "channelIdentifierHashStr":"0xa7712241a1a10abdada1c228c6935a71a9db80aa0bf2a13b59940159aa4eb4b5",
    "amountstr":0,
    "op":"preparewithdraw"
}
```
Example Response:  
```json
{
    "channel_identifier": "0xa7712241a1a10abdada1c228c6935a71a9db80aa0bf2a13b59940159aa4eb4b5",
    "open_block_number": 8100682,
    "partner_address": "0x10b256b3C83904D524210958FA4E7F9cAFFB76c6",
    "balance": 1000000000000000000000,
    "partner_balance": 1000000000000000000000000,
    "locked_amount": 0,
    "partner_locked_amount": 0,
    "token_address": "0x83073FCD20b9D31C6c6B3aAE1dEE0a539458d0c5",
    "state": 9,
    "state_string": "prepareForWithdraw",
    "settle_timeout": 150,
    "reveal_timeout": 10
}
```

When you want to cancel the state of the `preparewithdraw`, you can switch the channel state to the`opened` through the parameter`"op":"cancelprepare"`.

Parameter:  
```json
{
    "channelIdentifierHashStr":"0xa7712241a1a10abdada1c228c6935a71a9db80aa0bf2a13b59940159aa4eb4b5",
    "amountstr":0,
    "op":"cancelprepare"
}
```
Example Response:     
```json
{
    "channel_identifier": "0xa7712241a1a10abdada1c228c6935a71a9db80aa0bf2a13b59940159aa4eb4b5",
    "open_block_number": 8100682,
    "partner_address": "0x10b256b3C83904D524210958FA4E7F9cAFFB76c6",
    "balance": 1000000000000000000000,
    "partner_balance": 1000000000000000000000000,
    "locked_amount": 0,
    "partner_locked_amount": 0,
    "token_address": "0x83073FCD20b9D31C6c6B3aAE1dEE0a539458d0c5",
    "state": 9,
    "state_string": "opened",
    "settle_timeout": 150,
    "reveal_timeout": 10
}
```

Of course, as long as both channels are online and there is no lock, then you can directly `withdraw`, `op` parameters are not necessary.
When `amount`is greater than 0, the `op` parameter is meaningless.

Parameter:  
```json
{
    "channelIdentifierHashStr":"0xa7712241a1a10abdada1c228c6935a71a9db80aa0bf2a13b59940159aa4eb4b5",
    "amountstr":1000000000000000000000,
    
}
```
Example Response:   
```json
{
    "channel_identifier": "0xa7712241a1a10abdada1c228c6935a71a9db80aa0bf2a13b59940159aa4eb4b5",
    "open_block_number": 8100682,
    "partner_address": "0x201B20123b3C489b47Fde27ce5b451a0fA55FD60",
    "balance": 1000000000000000000000000,
    "partner_balance": 1000000000000000000000,
    "locked_amount": 0,
    "partner_locked_amount": 0,
    "token_address": "0x83073FCD20b9D31C6c6B3aAE1dEE0a539458d0c5",
    "state": 6,
    "state_string": "withdrawing",
    "settle_timeout": 150,
    "reveal_timeout": 10
}
```
If the withdrawn amount is larger than the available balance of the channel, an error message will be returned.such as "Error": "invalid withdraw amount, availabe=399999999999999999999,want=1000000000000000000000"”.

###  Close the channel

func (a *API) CloseChannel(channelIdentifier string, force bool) (callID string, err error)

Close the channel, which includes the unilateral close the channel and cooperative settle the channel.Set `force` default to `false`, meaning that channel participants cooperate settle the channel.In the case of consensus, the token can be returned to both accounts immediately (waiting for a few blocks);
Set `force`  to `true`, it will not negotiate with the other party, which means that the channel will be forcibly closed, waiting for the `settleTimeout` passed, then the settle channel can be performed,  finally the token will return to the accounts of both parties.

The return parameter is a callID, which is used to call the GetCallResult interface to query the call result.

Example Request:  
```json
{"channelIdentifier":"closed"，
  "force":false
}
```
Example Response:     
```json
{
    "channel_identifier": "0xf1fa19fa6a54912e32d6e6e1aa0baa14d530385c60266886ef7c18838f6e9bdc",
    "open_block_number": 2498052,
    "partner_address": "0x6B9E4D89EE3828e7a477eA9AA7B62810260e27E9",
    "balance": 0,
    "partner_balance": 0,
    "locked_amount": 0,
    "partner_locked_amount": 0,
    "token_address": "0x7B874444681F7AEF18D48f330a0Ba093d3d0fDD2",
    "state": 8,
    "StateString": "cooperativeSettling",
    "settle_timeout": 150,
    "reveal_timeout": 30
}
```
Once channel partner is offline or has the locks, the cooperate settle can't be carried out.The participant should alter the`force` to `true`, wait for settle_timeout and unilateral settle the channel.

Parameter:  

```json 
{"state":"closed",
  "force":true
}
```
Example Response:  
```json
{
    "channel_identifier": "0xc943251676c4e53b2669fbbf17ebcbb850da9cb0a907200c40f1342a37629489",
    "open_block_number": 2560169,
    "partner_address": "0x69C5621db8093ee9a26cc2e253f929316E6E5b92",
    "balance": 100000000000000000000,
    "partner_balance": 100000000000000000000,
    "locked_amount": 0,
    "partner_locked_amount": 0,
    "token_address": "0x7B874444681F7AEF18D48f330a0Ba093d3d0fDD2",
    "state": 2,
    "StateString": "closed",
    "settle_timeout": 150,
    "reveal_timeout": 30
}
```


### Settle the Channel
func (a *API) SettleChannel(channelIdentifier string) (callID string, err error)

In the case that the channel has been closed, after the settlement window period, the user can settle the channel.

The return parameter is a callID, which is used to call the GetCallResult interface to query the call result.

 Tips:When no new block is received from the connection point for more than one minute, an error message will be given when calling to settle the channel."call smc SyncProgress err, client is closed”,which means that the connection point need to synchronize new blocks.

Note: Since settle_timeout does not include the penalty period (in spectrum, which is 257  block, about an hour), the actual settlement time is about 410 block.

Parameter:   
 
```json
{
    "state":"settled"
}
```

Example Response:   
```json

{
    "channel_identifier": "0xc943251676c4e53b2669fbbf17ebcbb850da9cb0a907200c40f1342a37629489",
    "open_block_number": 2575160,
    "partner_address": "0x69C5621db8093ee9a26cc2e253f929316E6E5b92",
    "balance": 100000000000000000000,
    "partner_balance": 50000000000000000000,
    "locked_amount": 0,
    "partner_locked_amount": 0,
    "token_address": "0x7B874444681F7AEF18D48f330a0Ba093d3d0fDD2",
    "state": 3,
    "StateString": "settled",
    "settle_timeout": 150,
    "reveal_timeout": 30
}
```
**Status Codes :**  
- `200 OK`   
- `409 Conflict` ,such as, "failed to estimate gas needed: gas required exceeds allowance or always failing transaction",or "channel is still open".

## Transaction related interface (asynchronous)
### Initiate a transaction
func (a *API) Transfers(tokenAddress, targetAddress string, amountstr string, feestr string, secretStr string, isDirect bool, data string) (transfer string, err error)

This interface is used to initiate a transfer transaction, which is currently associated with PFS by default.

Parameters:
* `tokenAddress string`– Transaction token
* `targetAddress string` – Payee address
* `amountstr string` – Transfer amount
* `feestr string` – Specify the total cost of the transaction(When using specified transaction costs, the fee calculated by PFS is not used to send the transfer amount)
* `secretStr string` – transaction secret,which may be " ",if it is designated, the transaction should use the special secret. 
* `isDirect string` – whether it is a direct transfer. The default is false(MediatedTransfer)
* `data` -  Incidental information of the transaction. The length is not more than 256 byte.

Example Request:  
```json
{
    "amountstr":200000000000000000000000,
    "feestr":0,
    "isDirect":false,
    "data":"hello word"
}
```

Example Response: 

```json
{
    "initiator_address": "0x151E62a787d0d8d9EfFac182Eae06C559d1B68C2",
    "target_address": "0x10b256b3C83904D524210958FA4E7F9cAFFB76c6",
    "token_address": "0x3e9f443405072BA0147F06708E9c0b4663D1D645",
    "amount": 200000000000000000000000,
    "lockSecretHash": "0x98c04dd2a7e479f72b54af90728742f59f40ff89339c18ebe19846969009c883",
    "data": "hello word"
}
```
Note: In general, the parameter "fee" uses the default value of 0, that is, the total cost of the transfer is not specified. The sender will refer to the PFS recommended fee plan for transfer; in the case where "fee" is not 0, the amount is theoretically greater than or equal to the cost value recommended by PFS, otherwise the transfer may fail, prompting “no available route”.
In addition, the payment can also use the specified 'secret' method (not commonly used), the user needs to refer to the relevant interface in the Http REST API to generate a pair of `lock_secret_hash` / `secret`, which we will introduce in the later upgrade.

### Query the transaction status
func (a *API) GetTransferStatus(tokenAddressStr string, lockSecretHashStr string) (r string, err error)

There are two ways for users to send and receive transactions, that is, synchronous and asynchronous. If the asynchronous mode is used (sync is false), the GetTransferStatus interface can be called to query the status information of the current transaction. Among them, locksecrethash and token_address can be  obtained from the message returned by the asynchronous transfer transaction.

Parameters:
* `tokenAddress string`– Transaction token
* `lockSecretHashStr string` – The lockSecretHash returned from the Transfers interface

Example Response:
```json
{
    "Key": "0xf9c7a5491439238ad55c0a8e5a1b97eb205cb14f8137705c898d8d24fcf32465",
    "LockSecretHash": "0x0676b190e483c6ce425492e45726797c8b538a620a761371d68b1c96e7a8538e",
    "TokenAddress": "0x4092ce58b448abdfb59fbc84a0e30689f004d02e",
    "Status": 0,
    "StatusMessage": "DirectTransfer is sending\n"
}
```
Response JSON Array of Objects :
* 0 - Transfer init
* 1 - transfer can cancel
* 2 - transfer can not cancel
* 3 - transfer already success
* 4 - transfer cancel by user request
* 5 - transfer already failed
### Query the received successful transfer 
func (a *API) GetReceivedTransfers(from, to int64) (r string, err error)

The interface can be used to query the history information of all successful transfer which received from other partners, return data is an array of `ReceivedTransfer`

Example Response:  
```json
[
    {
        "Key": "0x79b789e88c3d2173af4048498f8c1ce66f019f33a6b8b06bedef51dde72bbbc1-2",
        "block_number": 4492809,
        "OpenBlockNumber": 0,
        "channel_identifier": "0x79b789e88c3d2173af4048498f8c1ce66f019f33a6b8b06bedef51dde72bbbc1",
        "token_address": "0xd82e6be96a1457d33b35cded7e9326e1a40c565d",
        "from_address": "0x201b20123b3c489b47fde27ce5b451a0fa55fd60",
        "nonce": 2,
        "amount": 10000000000000000000
    },
    {
        "Key": "0x79b789e88c3d2173af4048498f8c1ce66f019f33a6b8b06bedef51dde72bbbc1-6",
        "block_number": 4493353,
        "OpenBlockNumber": 0,
        "channel_identifier": "0x79b789e88c3d2173af4048498f8c1ce66f019f33a6b8b06bedef51dde72bbbc1",
        "token_address": "0xd82e6be96a1457d33b35cded7e9326e1a40c565d",
        "from_address": "0x201b20123b3c489b47fde27ce5b451a0fa55fd60",
        "nonce": 6,
        "amount": 20000000000000000000
    }
]
```

### Query the sent successful transfer 
func (a *API) GetSentTransfers(from, to int64) (r string, err error)
For the sender of the transfer, the interface can be used to query the history information of all successful transfer which sent from itself, return data is an array of `SenTransfer`

Example Response: 
```json
[
    {
        "Key": "0xd971f803c7ea39ee050bf00ec9919269cf63ee5d0e968d5fe33a1a0f0004f73d-3",
        "block_number": 4490372,
        "OpenBlockNumber": 0,
        "channel_identifier": "0xd971f803c7ea39ee050bf00ec9919269cf63ee5d0e968d5fe33a1a0f0004f73d",
        "to_address": "0x151e62a787d0d8d9effac182eae06c559d1b68c2",
        "token_address": "0xd82e6be96a1457d33b35cded7e9326e1a40c565d",
        "nonce": 3,
        "amount": 10000000000000000000
    },
    {
        "Key": "0xd971f803c7ea39ee050bf00ec9919269cf63ee5d0e968d5fe33a1a0f0004f73d-5",
        "block_number": 4490580,
        "OpenBlockNumber": 0,
        "channel_identifier": "0xd971f803c7ea39ee050bf00ec9919269cf63ee5d0e968d5fe33a1a0f0004f73d",
        "to_address": "0x151e62a787d0d8d9effac182eae06c559d1b68c2",
        "token_address": "0xd82e6be96a1457d33b35cded7e9326e1a40c565d",
        "nonce": 5,
        "amount": 10000000000000000000
    }
]
```
### Initiate a token swap transaction
func (a *API) TokenSwap(role string, lockSecretHash string, SendingAmountStr, ReceivingAmountStr string, SendingToken, ReceivingToken, TargetAddress string, SecretStr string) (callID string, err error)

This interface implements the decentralized atomic interchange operation of two Tokens.

When taker is called, the introduced `lockSecretHash` should be equal to the hash value of the `SecretStr` passed to the maker,that is ,the preimage of the `lock_secret_hash` in the taker request should  equal to the  `SecretStr` in the maker request.

The `taker` should  called  first during the token swap transaction.

The taker(called on one phone): 
```
    "role": "taker",
    "lockSecretHash":"0x8e90b850fdc5475efb04600615a1619f0194be97a6c394848008f33823a7ee03",
    "TargeAddress":"0x31DdaC67e610c22d19E887fB1937BEE3079B56Cd",
    "SendingAmountStr": 10000000000000000000,
    "SendingToken": "0x7B874444681F7AEF18D48f330a0Ba093d3d0fDD2",
    "ReceivingAmountStr": 100000000000000000000,
    "ReceivingToken": "0x9E7c6C6bf3A60751df8AAee9DEB406f037279C2a"
    "SecretStr":"",
```
Then the interface is called as a maker on another phone.

The maker:
```
    "role": "maker",
    "TargetAddress":"0x69C5621db8093ee9a26cc2e253f929316E6E5b92",
    "SendingAmountStr": 100000000000000000000,
    "SendingToken": "0x9E7c6C6bf3A60751df8AAee9DEB406f037279C2a",
    "ReceivingAmountStr": 10000000000000000000,
    "ReceivingToken": "0x7B874444681F7AEF18D48f330a0Ba093d3d0fDD2",
    "SecretStr": "0x40a6994181d0b98efcf80431ff38f9bae6fefda303f483e7cf5b7de7e341502a",
    "lockSecretHash":""
```

This function will return immediately, and the result of the exchange can be obtained by a callID through the GetCallResult.

Note: At present, the token swap transaction in the mobile phone API is not commonly used. Therefore, we do not generate the `lock_secret_hash` / `secret` interface. The user needs to use the corresponding interface of the http REST API to generate the `lock_secret_hash` / `secret`.

### Query result according to callID

func (a *API) GetCallResult(callID string) (r string)

Query the processing result of the asynchronous call, and return the result as follows:
```
{
"status":2,
"message":"error happening"
}
```
There are three status:
- 0: indicates that there is no result yet, and the message is empty at this time.
- 1: indicates that the processing is successful, the corresponding json result is included in the message, and the sample can refer to  http interface [rest_api.md](https://github.com/SmartMeshFoundation/Photon/blob/master/docs/rest_api.md)
- 2: indicates that the processing is failed. The message contains the corresponding Error information.

## Third party service
### Get the delegated data for Photon monitoring service  
func (a *API) ChannelFor3rdParty(channelIdentifier, thirdPartyAddress string) (r string, err error)

 Photon's principle determines that if a node is offline for a long time, it will bring its own financial security risk. Therefore, if Photon is likely to be offline for a long time (relative to the settletimeout specified by the channel creation),Then the relevant balanceproof should be delegated to the third-party service (Photon-Monitoring), and the Photon-Monitoring will submits the relevant BalanceProof when needed.

How to use the Photon-Monitoring,please refer to the [Photon-Monitoring](https://github.com/SmartMeshFoundation/Photon-Monitoring)

The returned data should be submitted directly to the trusted Photon-monitoring.

Example Response:
```json
{
    "channel_identifier": "0x029a853513e98050e670eb6d5f36217998a2c689ef2f1c65b5954051490d5965",
    "open_block_number": 2644876,
    "token_network_address": "0xa3b6481d1c6aa8ba538e8fa9d4d8b1dbadfd379c",
    "partner_address": "0x64d11d0cbb3f4f9bb3ee09709d4254f0899a6381",
    "update_transfer": {
        "nonce": 0,
        "transfer_amount": null,
        "locksroot": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "extra_hash": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "closing_signature": null,
        "non_closing_signature": null
    },
    "unlocks": null,
    "punishes": [
        {
            "lock_hash": "0xd4ec833949fa91e5f30b4e5e8b2e88cca10e8192a68e51bdb24d18220b3f519d",
            "additional_hash": "0xe800ff8e78b8e367fb165b76f6e0cd1f31d46e7fda640e02134eed4f5e983d53",
            "signature": "i24Lz6KVvDnlqsxhQzDu+IIx6jJKC4gdVyWg6NpkrfsEejzGV8F0CPB0oUUJjDZ2wmChKG6XjZQx24QkDmhsKhs="
        }
    ]
}
```
### Query the charging route from the node to the target node
func (a *API) FindPath(targetStr, tokenStr, amountStr string) (r string, err error) 

The user invokes the interface to query whether the target node has available routes and fees. If there are multiple routes with the same cost, they are given together.

Parameters:
- targetStr  The address of the target node 
- tokenStr   Token address
- amountstr  Transfer amount
 
Example Response:

200 OK

```json 
{
        "path_id": 0,
        "path_hop": 2,
        "fee": 10,
        "result": [
            "0x3bc7726c489e617571792ac0cd8b70df8a5d0e22",
            "0x8a32108d269c11f8db859ca7fac8199ca87a2722",
            "0xefb2e46724f675381ce0b3f70ea66383061924e9"
        ]
    } 
```
