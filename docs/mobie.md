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
### Android use
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
If the type is 1, it means that the status of a transaction initiated by the user has changed. For details, please refer to `Query the transaction status`. 

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
 ##### Type Description in InfoStruct
Level|name|value|description
---|---|----|----
Info|InfoTypeString|0 |Simple string notification,the format is not fixed, **the type which has been deprecated**
Info|InfoTypeSentTransferDetail|1|The status of the transaction initiated by the initiator has changed, and the format is fixed.
Info|InfoTypeChannelCallID|2|The operation on the channel has a result, the format is not fixed, and the caller decides what operation is based on the CallID.
Info|InfoTypeChannelStatus|3|Channel status has changed, which including the balance,the patner_balance,the locked_amount,the partner_locked_amount,the state,and so on
Info|InfoTypeContractCallTXInfo|4|User initiated Tx execution result notification, the format is relatively fixed
Info|InfoTypeInconsistentDatabase|5|During the transaction, it was found that the database of the receiving party was inconsistent. **Note that this message can only be used as a reference, and the other party may be maliciously falsified**
Error|InfoTypeBalanceNotEnoughError|6|The SMT in the account  is insufficient, and the bottom line of the on-chain transaction cannot be guaranteed. It must be recharged as soon as possible.
Error|InfoTypeCooperateSettleRefused|7|The  Cooperative settlement background execution was failed to close the channel, the other party refuses the request.
Error|InfoTypeCooperateSettleFailed|8|The  Cooperative settlement background execution was failed to close the channel, the TX is failure.
Error|InfoTypeWithdrawRefused|9|The  withdraw background execution was failed , the other party refuses the request.
Error|InfoTypeWithdrawFailed|10|The  withdraw background execution was failed ,  the TX is failure.
Info|InfoTypeReceivedMediatedTransfer|11|If the receiver receives MediatedTransfer, it does not mean that the transaction is successful, but only on behalf of receiving the message. If the transaction is successfully received, please use `OnReceivedTransfer`

**Info corresponding to 0,Warn corresponding to  1,Error corresponding to  2**
###### InfoTypeInconsistentDatabase
Message:
```go
	type inconsistentDatabase struct {
		ChannelIdentifier common.Hash    `json:"channel_identifier"`
		Target            common.Address `json:"target"`
    }
```
######   InfoTypeBalanceNotEnoughError
```go
	type notEnough struct {
		Need *big.Int `json:"need"`
		Have *big.Int `json:"have"`
    }
```
 ###### InfoTypeCooperateSettleRefused
 以及InfoTypeCooperateSettleFailed,   InfoTypeWithdrawRefused,InfoTypeWithdrawFailed
```go
type failedCooperate struct {
	Channel   common.Hash `json:"channel"`
	ErrorCode int         `json:"error_code"`
	ErrorMsg  string      `json:"error_message"`
}
```
##### InfoTypeReceivedMediatedTransfer
```go
type receivedTranser struct {
	Token      common.Address `json:"token"`
	From       common.Address `json:"from"`
	Amount     *big.Int       `json:"amount"`
	ID         string         `json:"id"`
	Expiration int64          `json:"expiration"`
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
    "error_code": 0,
    "error_message": "SUCCESS",
    "data": {
        "channel_identifier": "0xfe738aa39610416e4100036130af7ae00930021d5a51be60b55b96c12b1f4af5",
        "open_block_number": 1872482,
        "partner_address": "0x3bC7726c489E617571792aC0Cd8b70dF8A5D0e22",
        "balance": 1e+22,
        "partner_balance": 1e+22,
        "locked_amount": 0,
        "partner_locked_amount": 0,
        "token_address": "0xB31567308AD3c42D864FB41684bB40d3A2c57E1b",
        "state": 1,
        "state_string": "opened",
        "settle_timeout": 100,
        "reveal_timeout": 30,
        "closed_block": 0,
        "settled_block": 0,
        "our_balance_proof": {
            "nonce": 0,
            "transfer_amount": 0,
            "locks_root": "0x0000000000000000000000000000000000000000000000000000000000000000",
            "channel_identifier": {
                "channel_identifier": "0x0000000000000000000000000000000000000000000000000000000000000000",
                "open_block_number": 0
            },
            "message_hash": "0x0000000000000000000000000000000000000000000000000000000000000000",
            "contract_transfer_amount": 0,
            "contract_locksroot": "0x0000000000000000000000000000000000000000000000000000000000000000"
        },
        "partner_balance_proof": {
            "nonce": 0,
            "transfer_amount": 0,
            "locks_root": "0x0000000000000000000000000000000000000000000000000000000000000000",
            "channel_identifier": {
                "channel_identifier": "0x0000000000000000000000000000000000000000000000000000000000000000",
                "open_block_number": 0
            },
            "message_hash": "0x0000000000000000000000000000000000000000000000000000000000000000",
            "contract_transfer_amount": 0,
            "contract_locksroot": "0x0000000000000000000000000000000000000000000000000000000000000000"
        }
    }
}

```
Channel structure description ： 

- `error_code`:  error code

- `error_message`: Error code description

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

-  `closed_block`: Block height when the channel is closed

-  `settled_block`: Block height at the time of channel settlement

-  `our_balance_proof`: Our balanceproof data

-  `partner_balance_proof`: Counterparty balanceproof  of the Channel 

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

Among them, App can see only 1-9 states, other states can not be directly observed, which is internal use.  **prepareForWithdraw and prepareForCooperativeSettle will not appear on the mobile phone** , only appear when the meshbox  be used as intermediate node of the transaction.

DelegateState indicates whether the channel-related balanceProof is delegated to the PMS

delegateState|name|Description
---|---|---
0|ChannelDelegateStateNoNeed| Appears only when it is explicitly specified that no PMS is required.
1|ChannelDelegateStateWaiting|Waiting for delegate to PMS
2|ChannelDelegateStateSuccess|Successful delegation
3|ChannelDelegateStateFail|Delegate failure
4|ChannelDelegateStateFailAndNoEffectiveChain| The delegation is failed and there is no effective public chain

Currently, the interface results are changed from polling to synchronization. All  returns of the interfaces contain error codes and error messages. ErrCode 0 indicates success, and others indicate errors. It is meaningful to parse data fields when ErrCode is 0. Below are some error codes and message descriptions.

errorcode|errormessage|Description
---|---|---
0 |SUCCESS|Successful call
-1|unknown error|Unknown error
1|ArgumentError|Parameter error
2|PhotonAlreadyRunning|Start multiple photon instances
1000|HashLengthNot32|Parameter error
1001|Not found|Not found
1002|InsufficientBalance|There is not enough banalce in the channel to pay for transfers.
1003|InvalidAmount|The values supplied by the User are not integers and cannot be used to define a transfer value.
1005|NoPathError|No route to the requested destination address, excluding the case of inadequate channel capacity.
1006|SamePeerAddress|When a user attempts to create a channel, the addresses of the nodes on both sides are the same.
1007|InvalidState|The user's request for behavior is inconsistent with the current channel state.
1008|TransferWhenClosed|When the channel is closed, the user attempts to initiate a request for transfer.
1009|UnknownAddress|The addresses provided by users are valid, but not from known nodes.
1010|Locksroot mismatch|The received message contains an invalid locksroot, which is rejected when a pending lock is lost from the locksroot.
1011|InvalidNonce|The messages received from the partner contain invalid nonce values, which must be incremented in turn.
1012|TransferUnwanted|Nodes did not receive new transfers
1013|new transactions are not allowed|Stop creating new transfers and reject new transactions
1014|no mediated transfer on mesh only network| Indirect transfer is not allowed on Mesh network.
1015|secret and token cannot duplicate| Same token and same secret transactions are not allowed.
1016|NodeOffline|When sending a message, the other party is not online.
1017|TranasferCannotCancel| Failure to attempt to cancel a transfer that the secret has leaked. 
1018|DBError| Uncategorized database errors
1019|duplicate key| Duplicate key
1020|ErrTransferTimeout|Transaction timeout ,which do not mean that the transaction will succeed or fail, but the transaction is not succeeded in a given time.
1021|ErrUpdateButHaveTransfer|Trying to upgrade and discovering that there are still transactions going on.
1022|ErrNotChargeFee|Operations related to charges are performed, but charges are not enabled.
2000|insufficient balance to pay for gas|Not enough balance to pay gas
2001|closeChannel|An error occurred while closing the channel on the chain.
2002|RegisterSecret|An error occurred while registering a secret on the chain.
2003|Unlock|An error occurred while unlock the locks on the chain.
2004|UpdateBalanceProof|An error occurred while submitting balance proof on the chain.
2005|punish|An error occurred while executing punish on the chain.
2006|settle|An error occurred while performing settle on the chain.
2007|deposit|An error occurred while executing deposit on the chain.
2008|ErrSpectrumNotConnected|Not connected to the public chain（spectrum).
2009|ErrTxWaitMined|Wait for returning error of mining.
2010|ErrTxReceiptStatus|The transfer was packaged, but it failed.
2011|ErrSecretAlreadyRegistered|Attempt to connect to the public chain to register the secret, but the secret has been registered.
2012|ErrSpectrumSyncError|Photon has connected to the public chain, but did not create the block for a long time or was synchronized.
2013|ErrSpectrumBlockError|The number of locally processed blocks is not consistent with the number which public chain reporting blocks.
2999|unkown spectrum rpc error|Other Ethereum RPC errors
3001|TokenNotFound|No corresponding token was found
3002|ChannelNotFound|No corresponding channel was found
3003|NoAvailabeRoute|No available routes
3004|TransferNotFound|No corresponding transfer was found.
3005|ChannelAlreadExist|Channels already exist.
5000|CannotWithdarw|Channels are not cooperatively withdraw now, such as transactions in progress.
5001|ErrChannelState|The channel state in which the corresponding operation cannot be performed, one attempt to execute certain transactions, such as initiating transactions on closed channels.
5002|Channel only can settle after timeout|Attempt the settle the channel before the timeout
5003|NotParticipant|The given address is not one of the participants of the channel.
5004|ChannelNoSuchLock|There is no corresponding lock in the channel.
5005|ErrChannelEndStateNoSuchLock|The corresponding lock cannot be found in the current participant of the channel
5006|ErrChannelLockAlreadyExpired|The lock in the channel has expired.
5007|ErrChannelBalanceDecrease|There has been a reduction in channel balance(which means the balance in the contract).
5008|ErrChannelTransferAmountMismatch|Transferamount was mismatched in received transactions.
5009|ErrChannelBalanceProofAlreadyRegisteredOnChain| Attempts to modify local balance proof after submitting balanceproof
5010|ErrChannelDuplicateLock|A lock for this secret already exists in the channel.
5011|ErrChannelTransferAmountDecrease|The transaction is received, but transferamount became smaller.
5012|ErrRemoveNotExpiredLock|Attempt to remove an unexpired lock.
5013|ErrUpdateBalanceProofAfterClosed|Trying to update balance proof of the  participant or the partner after the channel closed
5014|ErrChannelIdentifierMismatch|Channel ID mismatch
5015|ErrChannelInvalidSender|Receiving transactions from unknown participants
5016|ErrChannelBalanceNotMatch|Cooperating to close the channel, the amount check of withdraw was mismatched.
5017|ErrChannelLockMisMatch|The specified locks in the received transaction do not match local locks.
5018|ErrChannelWithdrawAmount|Excessive amount to withdraw
5019|ErrChannelLockExpirationTooLarge|Receiving a transaction, the specified expiration time is too long.
5020|ErrChannelRevealTimeout|The specified reveal timeout is illegal. 
5021|ErrChannelBalanceProofNil|The balanceproof is null.
5022|ErrChannelCloseClosedChannel|Attempts to close closed channel.
5023|ErrChannelBackgroundTx|BackgroundError in transaction execution. 
5024|ErrChannelWithdrawButHasLocks|Withdraw requests cannot be sent in the existence of locks.
5025|ErrChannelCooperativeSettleButHasLocks| CooperativeSettle requests cannot be sent in the existence of locks.
5026|ErrInvalidSettleTimeout|The timeout value submitted by the user is less than the minimum settle timeout value.
6000|transport type error|Unknown transport layer errors.
6001|ErrSubScribeNeighbor|Subscriber online information error
9999|ErrUnknown|Unknown error, the code should be incomplete, the error classification is not detailed enough
## Query interface
### Query the account address of the running photon node
func (a *API) Address() (addr string)

Example Response：
```json
{
    "error_code": 0,
    "error_message": "SUCCESS",
    "data": "0x97Cd7291f93F9582Ddb8E9885bF7E77e3f34Be40"
}
```
### Query the list of all the registered token
func (a *API) Tokens() (tokens string)

Example Response：
```json
{
    "error_code": 0,
    "error_message": "SUCCESS",
    "data": [
        "0xF0123C3267Af5CbBFAB985d39171f5F5758C0900",
        "0xB31567308AD3c42D864FB41684bB40d3A2c57E1b",
        "0x270831A3C8dB8e515ba4ee2c6b3087E58e8DD1C7",
        "0x481Df7AC195d000546592e7D39488134FdCd042A",
        "0xB5F80e9013d62A891B062595C3E864B3D4612a78"
    ]
}
```
### Query all the channels that you participate in under some token.
func (a *API) TokenPartners(tokenAddress string) (channels string, err error)

Example Response：
```json
{
    "error_code": 0,
    "error_message": "SUCCESS",
    "data": [
        {
            "partner_address": "0xC445a8C326A8fD5a3e250C7dc0EFc566eDcB263B",
            "channel": "api/1/channles/0xe4c61eac5f3f45ea62c7f021cc0aa6a774feb14fed3eaa28af16b512f7fec966"
        },
        {
            "partner_address": "0x97Cd7291f93F9582Ddb8E9885bF7E77e3f34Be40",
            "channel": "api/1/channles/0xfe738aa39610416e4100036130af7ae00930021d5a51be60b55b96c12b1f4af5"
        }
    ]
}
```
### Query all the channels of the node
func (a *API) GetChannelList() (channels string, err error)

Example Response：
```json
{
  "error_code": 0,
  "error_message": "SUCCESS",
  "data": [
    {
      "channel_identifier": "0xfe738aa39610416e4100036130af7ae00930021d5a51be60b55b96c12b1f4af5",
      "open_block_number": 1872482,
      "partner_address": "0x3bC7726c489E617571792aC0Cd8b70dF8A5D0e22",
      "balance": 10000000000000000000000,
      "partner_balance": 10000000000000000000000,
      "locked_amount": 0,
      "partner_locked_amount": 0,
      "token_address": "0xB31567308AD3c42D864FB41684bB40d3A2c57E1b",
      "state": 1,
      "state_string": "opened",
      "settle_timeout": 100,
      "reveal_timeout": 30
    }
  ]
}
```
### Query information about special channel
func (a *API) GetOneChannel(channelIdentifier string) (channel string, err error)

Example Response：
```json
{
    "error_code": 0,
    "error_message": "SUCCESS",
    "data": {
        "channel_identifier": "0xfe738aa39610416e4100036130af7ae00930021d5a51be60b55b96c12b1f4af5",
        "open_block_number": 1872482,
        "partner_address": "0x3bC7726c489E617571792aC0Cd8b70dF8A5D0e22",
        "balance": 1e+22,
        "partner_balance": 1e+22,
        "locked_amount": 0,
        "partner_locked_amount": 0,
        "token_address": "0xB31567308AD3c42D864FB41684bB40d3A2c57E1b",
        "state": 1,
        "state_string": "opened",
        "settle_timeout": 100,
        "reveal_timeout": 30,
        "closed_block": 0,
        "settled_block": 0,
        "our_balance_proof": {
            "nonce": 0,
            "transfer_amount": 0,
            "locks_root": "0x0000000000000000000000000000000000000000000000000000000000000000",
            "channel_identifier": {
                "channel_identifier": "0x0000000000000000000000000000000000000000000000000000000000000000",
                "open_block_number": 0
            },
            "message_hash": "0x0000000000000000000000000000000000000000000000000000000000000000",
            "contract_transfer_amount": 0,
            "contract_locksroot": "0x0000000000000000000000000000000000000000000000000000000000000000"
        },
        "partner_balance_proof": {
            "nonce": 0,
            "transfer_amount": 0,
            "locks_root": "0x0000000000000000000000000000000000000000000000000000000000000000",
            "channel_identifier": {
                "channel_identifier": "0x0000000000000000000000000000000000000000000000000000000000000000",
                "open_block_number": 0
            },
            "message_hash": "0x0000000000000000000000000000000000000000000000000000000000000000",
            "contract_transfer_amount": 0,
            "contract_locksroot": "0x0000000000000000000000000000000000000000000000000000000000000000"
        }
    }
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

1.Create channel and deposit:
 
   - `new_channel` sets `true`，which means open a new channel and deposit to the channel; if there is no channel between the participants. `false`is no meaning for `new_channel`，which will return the error message "There is no channel".
   
   - `settle_timeout`represent the settlement window for new channel,for example,settle_timeout：100; if the `settle_timeout` set to 0,the default window period is used which is 600 block. 
    
2.Only deposit:

   - `new_channel` must set to `false`，which means the channel has been existed；If the channel has been existed ,there is Meaningless to set the `new_channel`statue as `ture`，which will response the error message "The channel has already existed". 
   
   - `settle_timeout`must set to 0,because the channel has already existed.


Example Response:
```json
{
    "error_code": 0,
    "error_message": "SUCCESS",
    "data": {
        "channel_identifier": "0xfe738aa39610416e4100036130af7ae00930021d5a51be60b55b96c12b1f4af5",
        "open_block_number": 1872482,
        "partner_address": "0x3bC7726c489E617571792aC0Cd8b70dF8A5D0e22",
        "balance": 10000000000000000000000,
        "partner_balance": 0,
        "locked_amount": 0,
        "partner_locked_amount": 0,
        "token_address": "0xB31567308AD3c42D864FB41684bB40d3A2c57E1b",
        "state": 1,
        "state_string": "opened",
        "settle_timeout": 100,
        "reveal_timeout": 30
    }
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
  "error_code": 0,
  "error_message": "SUCCESS",
  "data": {
    "channel_identifier": "0xfe738aa39610416e4100036130af7ae00930021d5a51be60b55b96c12b1f4af5",
    "open_block_number": 1872482,
    "partner_address": "0x3bC7726c489E617571792aC0Cd8b70dF8A5D0e22",
    "balance": 20000000000000000000000,
    "partner_balance": 10000000000000000000000,
    "locked_amount": 0,
    "partner_locked_amount": 0,
    "token_address": "0xB31567308AD3c42D864FB41684bB40d3A2c57E1b",
    "state": 9,
    "state_string": "prepareForWithdraw",
    "settle_timeout": 100,
    "reveal_timeout": 30
  }
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
  "error_code": 0,
  "error_message": "SUCCESS",
  "data": {
    "channel_identifier": "0xfe738aa39610416e4100036130af7ae00930021d5a51be60b55b96c12b1f4af5",
    "open_block_number": 1872482,
    "partner_address": "0x3bC7726c489E617571792aC0Cd8b70dF8A5D0e22",
    "balance": 20000000000000000000000,
    "partner_balance": 10000000000000000000000,
    "locked_amount": 0,
    "partner_locked_amount": 0,
    "token_address": "0xB31567308AD3c42D864FB41684bB40d3A2c57E1b",
    "state": 1,
    "state_string": "opened",
    "settle_timeout": 100,
    "reveal_timeout": 30
  }
}
```

Of course, as long as both channels are online and there is no lock, then you can directly `withdraw`, `op` parameters are not necessary.
When `amount`is greater than 0, the `op` parameter is meaningless.

Parameter:  
```json
{
    "channelIdentifierHashStr":"0xfe738aa39610416e4100036130af7ae00930021d5a51be60b55b96c12b1f4af5",
    "amountstr":100000000000000000000,
    
}
```
Example Response:   
```json
{
  "error_code": 0,
  "error_message": "SUCCESS",
  "data": {
    "channel_identifier": "0xfe738aa39610416e4100036130af7ae00930021d5a51be60b55b96c12b1f4af5",
    "open_block_number": 1872482,
    "partner_address": "0x3bC7726c489E617571792aC0Cd8b70dF8A5D0e22",
    "balance": 20000000000000000000000,
    "partner_balance": 10000000000000000000000,
    "locked_amount": 0,
    "partner_locked_amount": 0,
    "token_address": "0xB31567308AD3c42D864FB41684bB40d3A2c57E1b",
    "state": 6,
    "state_string": "withdrawing",
    "settle_timeout": 100,
    "reveal_timeout": 30
  }
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
  "error_code": 0,
  "error_message": "SUCCESS",
  "data": {
    "channel_identifier": "0xe4c61eac5f3f45ea62c7f021cc0aa6a774feb14fed3eaa28af16b512f7fec966",
    "open_block_number": 1694460,
    "partner_address": "0xC445a8C326A8fD5a3e250C7dc0EFc566eDcB263B",
    "balance": 10000000000000000000000,
    "partner_balance": 10000000000000000000000,
    "locked_amount": 0,
    "partner_locked_amount": 0,
    "token_address": "0xB31567308AD3c42D864FB41684bB40d3A2c57E1b",
    "state": 7,
    "state_string": "cooperativeSettling",
    "settle_timeout": 100,
    "reveal_timeout": 30
  }
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
  "error_code": 0,
  "error_message": "SUCCESS",
  "data": {
    "channel_identifier": "0xe4c61eac5f3f45ea62c7f021cc0aa6a774feb14fed3eaa28af16b512f7fec966",
    "open_block_number": 1890493,
    "partner_address": "0xC445a8C326A8fD5a3e250C7dc0EFc566eDcB263B",
    "balance": 10000000000000000000000,
    "partner_balance": 10000000000000000000000,
    "locked_amount": 0,
    "partner_locked_amount": 0,
    "token_address": "0xB31567308AD3c42D864FB41684bB40d3A2c57E1b",
    "state": 4,
    "state_string": "closing",
    "settle_timeout": 100,
    "reveal_timeout": 30
  }
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
  "error_code": 0,
  "error_message": "SUCCESS",
  "data": {
    "channel_identifier": "0xe4c61eac5f3f45ea62c7f021cc0aa6a774feb14fed3eaa28af16b512f7fec966",
    "open_block_number": 1890493,
    "partner_address": "0xC445a8C326A8fD5a3e250C7dc0EFc566eDcB263B",
    "balance": 10000000000000000000000,
    "partner_balance": 10000000000000000000000,
    "locked_amount": 0,
    "partner_locked_amount": 0,
    "token_address": "0xB31567308AD3c42D864FB41684bB40d3A2c57E1b",
    "state": 5,
    "state_string": "settling",
    "settle_timeout": 100,
    "reveal_timeout": 30
  }
}
```
**Status Codes :**  
- `200 OK`   
- `409 Conflict` ,such as, "failed to estimate gas needed: gas required exceeds allowance or always failing transaction",or "channel is still open".

### Query contract call

func (a *API)  ContractCallTXQuery(channelIdentifierStr string, openBlockNumber int, tokenAddressStr, txTypeStr, txStatusStr string) (result string)

Query the result of the channel function interface contract call. Mainly include:deposit, close, settle, withdraw, cooperating settle, etc.

Example Request: 

```json
{
	"token_address":"",
	"tx_type":"ChannelSettle"
}
```

Example Response:

```json
{
    "error_code": 0,
    "error_message": "SUCCESS",
    "data": [
        {
            "tx_hash": "0x980c44a5a75224ed140549cf94cd37fa481c1d2e5bee12e507acea695af0f30a",
            "channel_identifier": "0xb943bb364667a2c9c06526be5a9c03e99544ded5b2380be035608e016bdfa8ac",
            "open_block_number": 1251429,
            "token_address": "0x2158c8c27ab31602f462084bdc47ab5c9d339b26",
            "type": "ChannelSettle",
            "is_self_call": true,
            "tx_params": "{\"token_address\":\"0x2158c8c27ab31602f462084bdc47ab5c9d339b26\",\"p1_address\":\"0x97251ddfe70ea44be0e5156c4e3aadd30328c6a5\",\"p1_transfer_amount\":0,\"p1_locks_root\":\"0x0000000000000000000000000000000000000000000000000000000000000000\",\"p2_address\":\"0x3de45febbd988b6e417e4ebd2c69e42630fefbf0\",\"p2_transfer_amount\":0,\"p2_locks_root\":\"0x0278c2c9445b7930a2d95cce4d3de0fc6845782a721d7bf92ddaf8189b00a936\"}",
            "tx_status": "success",
            "events": null,
            "pack_block_number": 1252608,
            "call_time": 1551753432,
            "pack_time": 1551753433,
            "gas_price": 20000000000,
            "gas_used": 44568
        },
    ]
        }

```
The transaction parameters are structured as follows:

DepositTXParams :

Deposit parameter

```json
{
  "token_address":"0x2158c8c27ab31602f462084bdc47ab5c9d339b26",
	"participant_address":"0x3de45febbd988b6e417e4ebd2c69e42630fefbf0",
	"partner_address":"0x97251ddfe70ea44be0e5156c4e3aadd30328c6a5",
	"amount":0,
	"locks_root":"0x0000000000000000000000000000000000000000000000000000000000000000",
	"settle_timeout":500,
}
```
Parameter Description:

	TokenAddress      token address

	ParticipantAddress  Own address

	PartnerAddress     Address of the other party

	Amount            Deposit amount

	SettleTimeout      When equal to 0, it means Deposit, and when it is greater than 0, it means OpenAndDeposit.

ChannelCloseOrChannelUpdateBalanceProofTXParams：

Close the channel or UpdateBalanceProof parameters, the two operations are multiplexed, according to the Type in the upper layer TXInfo

```json
{
  "token_address":"0x2158c8c27ab31602f462084bdc47ab5c9d339b26",
	"participant_address":"0x3de45febbd988b6e417e4ebd2c69e42630fefbf0",
	"partner_address":"0x97251ddfe70ea44be0e5156c4e3aadd30328c6a5",
	"transfer_amount":0,
	"locks_root":"0x0000000000000000000000000000000000000000000000000000000000000000",
	"nonce":0,
	"extra_hash":"0x0000000000000000000000000000000000000000000000000000000000000000",
	"signature":null
}
```
 Parameter Description:

	TokenAddress          token address

	ParticipantAddress    Own address

	PartnerAddress      Address of the other party

	TransferAmount      The transferAmount of the other party

	LocksRoot           The locksroot  of the other party

	Nonce               The nonce of the other party

	ExtraHash             Metadata

	Signature          The other party's BalanceProof signature


ChannelSettleTXParams：

 Channel settlement parameters, p1 for yourself, p2 for the other

```json
{
 "token_address":"0x2158c8c27ab31602f462084bdc47ab5c9d339b26",
	"p1_address":"0x3de45febbd988b6e417e4ebd2c69e42630fefbf0",
	"p1_transfer_amount":1,
	"p1_locks_root":"0x0000000000000000000000000000000000000000000000000000000000000000",
	"p2_address":"0x97251ddfe70ea44be0e5156c4e3aadd30328c6a5",
	"p2_transfer_amount":1,
	"p2_locks_root":"0x0000000000000000000000000000000000000000000000000000000000000000",
}
```

Parameter Description:

	TokenAddress     token  address 

	P1Address         Own address

	P1TransferAmount  Own Transferamount

	P1LocksRoot       Own locksroot

	P2Address         Address of the other party 

	P2TransferAmount  Transferamount  of the other party

	P2LocksRoot       Locksroot  of the other party



 ChannelWithDrawTXParams:
 
Channel withdraw parameters, p1 for yourself, p2 for the other

  ```json
{
    "token_address":"0x2158c8c27ab31602f462084bdc47ab5c9d339b26",
	"p1_address":"0x3de45febbd988b6e417e4ebd2c69e42630fefbf0",
	"p2_address":"0x97251ddfe70ea44be0e5156c4e3aadd30328c6a5",
	"p1_balance":10,
	"p1_withdraw":5,
	"p1_signature":"9tu3sP8vYLl8OTvs3TPmsftvLjRb+HiUKFmp7mYvbANlBzbslHa/y90D35yC/bDUygXtpgPyqvJZvdyespdklhs=",
	"p2_signature":"9tu3sP8vYLl8OTvs3TPmsftvLjRb+HiUKFmp7mYvbANlBzbslHa/y90D35yC/bDUygXtpgPyqvJZvdyespdklhs=",
}
```

 Parameter Description:

	TokenAddress   token address

	P1Address      Own address

	P2Address      Address of the other party 

	P1Balance     Own balance

	P1Withdraw    Own withdraw amount

	P1Signature   Own signature

	P2Signature   The signature of the other party

 ChannelCooperativeSettleTXParams：
 
   The parameters of the channel cooperation close, p1 is itself, p2 is the other party

 ```json
{
   "token_address":"0x2158c8c27ab31602f462084bdc47ab5c9d339b26",
	"p1_address":"0x3de45febbd988b6e417e4ebd2c69e42630fefbf0",
	"p1_balance":10,
	"p2_address":"0x97251ddfe70ea44be0e5156c4e3aadd30328c6a5",
	"p2_balance":5,
	"p1_signature":"9tu3sP8vYLl8OTvs3TPmsftvLjRb+HiUKFmp7mYvbANlBzbslHa/y90D35yC/bDUygXtpgPyqvJZvdyespdklhs=",
	"p2_signature":"9tu3sP8vYLl8OTvs3TPmsftvLjRb+HiUKFmp7mYvbANlBzbslHa/y90D35yC/bDUygXtpgPyqvJZvdyespdklhs=",
}
``` 
Parameter Description:

	TokenAddress   token address

	P1Address      Own address

	P1Balance      Own balance

	P2Address       The address of the other party

	P2Balance      The balance of the other party

	P1Signature    Own signature

	P2Signature    The signature of the other party
## Transaction related interface (asynchronous)
### Initiate a transaction
func (a *API) Transfers(tokenAddress, targetAddress string, amountstr string, secretStr string, isDirect bool, data string, routeInfoStr string) (result string)

This interface is used to initiate a transfer transaction, which is currently associated with PFS by default.

Parameters:

* `tokenAddress string`– Transaction token

* `targetAddress string` – Payee address

* `amountstr string` – Transfer amount

* `secretStr string` – transaction secret,which may be " ",if it is designated, the transaction should use the special secret. 

* `isDirect string` – whether it is a direct transfer. The default is false(MediatedTransfer)

* `data` -  Incidental information of the transaction. The length is not more than 256 byte.
* `routeInfoStr string` – Specify the route and total cost of the transaction.

Example Request:  
```json
{
    "amountstr":1000000000000000000000,
      "isDirect":false,
    "data":"hello word"
     "route_info":[
    {
        "path_id": 0,
        "path_hop": 2,
        "fee": 23611121,
        "result": [
            "0x3bc7726c489e617571792ac0cd8b70df8a5d0e22",
            "0xc445a8c326a8fd5a3e250c7dc0efc566edcb263b",
            "0xd5dc7504e0b448b1c62d86306ae8e4a5836fc1a1"
        ]
}
```

Example Response: 

```json
{
  "error_code": 0,
  "error_message": "SUCCESS",
  "data": {
    "initiator_address": "0x3bC7726c489E617571792aC0Cd8b70dF8A5D0e22",
    "target_address": "0x97Cd7291f93F9582Ddb8E9885bF7E77e3f34Be40",
    "token_address": "0xB31567308AD3c42D864FB41684bB40d3A2c57E1b",
    "amount": 1000000000000000000000,
    "lockSecretHash": "0xa27591f7a7eb6922d6dac202fe08352cc2af79ce43b7692d04fe9e72524940b3",
    "data": ""
  }
}
```
Note: In this version, the transfer route and total cost of the transfer is specified. The sender will refer to the PFS recommended fee plan for transfer;  the transfer amount is theoretically greater than or equal to the cost value recommended by PFS, otherwise the transfer may fail, prompting “no available route”.
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
{
  "error_code": 0,
  "error_message": "SUCCESS",
  "data": [
    {
      "Key": "0xfe738aa39610416e4100036130af7ae00930021d5a51be60b55b96c12b1f4af5-1890429-2",
      "block_number": 1890583,
      "OpenBlockNumber": 1890429,
      "channel_identifier": "0xfe738aa39610416e4100036130af7ae00930021d5a51be60b55b96c12b1f4af5",
      "token_address": "0xb31567308ad3c42d864fb41684bb40d3a2c57e1b",
      "initiator_address": "0x3bc7726c489e617571792ac0cd8b70df8a5d0e22",
      "nonce": 2,
      "amount": 1000000000000000000000,
      "data": "",
      "time_stamp": "2019-02-18T15:22:10+08:00"
    },
    {
      "Key": "0xfe738aa39610416e4100036130af7ae00930021d5a51be60b55b96c12b1f4af5-1890429-4",
      "block_number": 1890656,
      "OpenBlockNumber": 1890429,
      "channel_identifier": "0xfe738aa39610416e4100036130af7ae00930021d5a51be60b55b96c12b1f4af5",
      "token_address": "0xb31567308ad3c42d864fb41684bb40d3a2c57e1b",
      "initiator_address": "0x3bc7726c489e617571792ac0cd8b70df8a5d0e22",
      "nonce": 4,
      "amount": 1000000000000000000000,
      "data": "",
      "time_stamp": "2019-02-18T15:40:47+08:00"
    }
  ]
}
```

### Query the sent successful transfer 
func (a *API) GetSentTransfers(from, to int64) (r string, err error)
For the sender of the transfer, the interface can be used to query the history information of all successful transfer which sent from itself, return data is an array of `SenTransfer`

Example Response: 
```json
{
  "error_code": 0,
  "error_message": "SUCCESS",
  "data": [
    {
      "Key": "0xfe738aa39610416e4100036130af7ae00930021d5a51be60b55b96c12b1f4af5-1890429-3",
      "block_number": 1890583,
      "open_block_number": 1890429,
      "channel_identifier": "0xfe738aa39610416e4100036130af7ae00930021d5a51be60b55b96c12b1f4af5",
      "target_address": "0x97cd7291f93f9582ddb8e9885bf7e77e3f34be40",
      "token_address": "0xb31567308ad3c42d864fb41684bb40d3a2c57e1b",
      "nonce": 3,
      "amount": 1000000000000000000000,
      "data": "",
      "time_stamp": "2019-02-18T15:22:10+08:00"
    },
    {
      "Key": "0xfe738aa39610416e4100036130af7ae00930021d5a51be60b55b96c12b1f4af5-1890429-5",
      "block_number": 1890656,
      "open_block_number": 1890429,
      "channel_identifier": "0xfe738aa39610416e4100036130af7ae00930021d5a51be60b55b96c12b1f4af5",
      "target_address": "0x97cd7291f93f9582ddb8e9885bf7e77e3f34be40",
      "token_address": "0xb31567308ad3c42d864fb41684bb40d3a2c57e1b",
      "nonce": 5,
      "amount": 1000000000000000000000,
      "data": "",
      "time_stamp": "2019-02-18T15:40:47+08:00"
    }
  ]
}
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
    "error_code": 0,
    "error_message": "SUCCESS",
    "data": [
        {
            "path_id": 0,
            "path_hop": 1,
            "fee": 1000000000000000005,
            "result": [
                "0x3bc7726c489e617571792ac0cd8b70df8a5d0e22",
                "0xc445a8c326a8fd5a3e250c7dc0efc566edcb263b"
            ]
        }
    ]
}
```
### Version query
func (a *Api) Version() string

 Version Information Query Interface
 
 **Example Response :**
 ```json
   {
            "error_code": 0,
            "error_message": "SUCCESS",
            "data": {
                "go_version": "goversiongo1.11linux/amd64",
                "git_commit": "17b4d194449e2da643c7b0309063720b602a0b2d",
                "build_date": "2019-03-19-17:01:58",
                "version": "1.1.0--17b4"
            }
        }
```
##  Interface for temporary access to photon information 

func NewSimpleAPI(datadir, address string) (api *SimpleAPI, err error)

This temporary interface is added because the photon query channel information is not started in the app.

Parameter 1: The datadir parameter of the Startup function, which is the fourth parameter.
Parameter 2: Startup function address parameter, which is the first parameter

### Example
```go
a, err := NewSimpleApi("/Users/bai/sm/Photon/cmd/photon/.photon", "0x292650fee408320D888e06ed89D938294Ea42f99")
r = a.BalanceAvailabelOnPhoton("0x6601F810eaF2fa749EEa10533Fd4CC23B8C791dc")
a.Stop()
```
### BalanceAvailabelOnPhoton 

The parameter is the token you want to query.

Return example


```json
{
    "error_code":0,
    "error_message":"SUCCESS",
    "data":20
}
```
**Note**

 NewSimpleApi returns, **Make sure to call Stop** after it is used, otherwise it will cause subsequent photon to fail to start.
