# Photon’s Mobile API Documentation
## Installation
Photon mobile SDk compilation must require the gomobile tool to work properly. Please refer to [gomobile](https://godoc.org/golang.org/x/mobile) for gomobile installation 

```bash
Cd mobile
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
Photon relies on gomobile to automate interface encapsulation. Because it is a cross-language call, it is unavoidable that it is a type conversion problem.
In order to avoid such problems, Photon provides interfaces to almost all basic types (int, string, error).

### Starting a photon node
func StartUp(address, keystorePath, ethRPCEndPoint, dataDir, passwordfile, apiAddr, listenAddr, logFile string, registryAddress string, otherArgs *Strings) (api *API, err error)

parameter:
* `address string` – the account address used by the photon node
* `keystorePath string` – account private key save path
* `ethRPCEndPoint string` – public chain node host, http protocol
* `dataDir string` – photon db path
* `passwordfile string` – account password file path
* `apiAddr string` – http api listening port
* `listenAddr string` – udp listening port
* `logFile string` – log file path
* `registryAddress string` – TokenNetworkRegistry contract address
* `otherArgs mobile.Strings` – other parameters, see photon -h

If you need to pass other parameters than the default parameters, you can refer to the following ways:
```go
otherArgs := mobile.NewStrings(2)
Err = otherArgs.Set(0, fmt.Sprintf("--registry-contract-address=%s", registryContractAddress))
If err != nil {
    Return err
}
Err = otherArgs.Set(1, fmt.Sprintf("--help"))
If err != nil {
    Return err
}
```
return:
* `api *API` - startup successfully returns api handle
* `err error` – error message

### Stop a photon node
func (a *API) Stop()

### Switching the photon operating environment
func (a *API) SwitchNetwork(isMesh bool)

Switch network environment, Mesh or Internet
In the Mesh network, the nodes directly communicate using the UDP protocol, and the App needs to notify the Photon other nodes through the UpdateMeshNetworkNodes.

### Notify the photon node network to disconnect
func (a *API) NotifyNetworkDown() error

Proactively inform the photon node that the network is disconnected and let the photon node start trying to reconnect

In the mobile phone network environment, due to network complexity, such as WiFi disconnection, these events Photon can not be directly perceived from the system, the App needs to actively tell Photon to take appropriate processing.

### Subscribe to photon events
func (a *API) Subscribe(handler NotifyHandler) (sub *Subscription, err error)

Subscribe to photon node events, including transaction notifications, error notifications, etc.
```go
// NotifyHandler is a client-side subscription callback to invoke on events and
// subscription failure.
Type NotifyHandler interface {
//some unexpected error
OnError(errCode int, failure string)
//OnStatusChange server connection status change
OnStatusChange(s string)
//OnReceivedTransfer receive a transfer
OnReceivedTransfer(tr string)
//OnSentTransfer a transfer sent success
OnSentTransfer(tr string)
// OnNotify get some important message raiden want to notify upper application
OnNotify(level int, info string)
}
```

#### OnError
 Notify Photon that an unrecoverable error has occurred. Any function of Photon must be restarted before it can be used. Since it is considered that the integration of Photon may be single-process, The unknown error may cause the App to quit, so even if there is an unpredictable error inside Photon, Photon will intercept and report to the App, and the App will decide whether to exit immediately or continue to use it.
- `errCode` is the error code
- `failure` is an error message description
Restart the Photon mode:
```go
api.Stop()
newAPI, err:=Startup(...)
```
#### OnStatusChange
`s` is the json code of the structure below
```go
//ConnectionStatus status of network connection
Type ConnectionStatus struct {
XMPPStatus netshare.Status
EthStatus netshare.Status
LastBlockTime string
}
```
Where `XMPPStatus` and `EthStatus` are defined as follows:
```go
// Status shows actual connection status.
Type Status int

Const (
//Disconnected init status
Disconnected = Status(iota)
//Connected connection status
Connected
//Closed user closed
Closed
//Reconnecting connection error
Reconnecting
)
```
#### OnReceivedTransfer
`tr` is the json code of the following structure
```go
//ReceivedTransfer tokens I have received and where it comes from
Type ReceivedTransfer struct {
Key string `storm:"id"`
BlockNumber int64 `json:"block_number" storm:"index"`
OpenBlockNumber int64
ChannelIdentifier common.Hash `json:"channel_identifier"`
TokenAddress common.Address `json:"token_address"`
FromAddress common.Address `json:"from_address"`
Nonce uint64 `json:"nonce"`
Amount *big.Int `json:"amount"`
}
```
Note: This interface does not contain transactions that participate as intermediate intermediate nodes.
#### OnSentTransfer
`tr` is the json encoding of the following structure
```go
//SentTransfer transfer's I have sent and success.
Type SentTransfer struct {
Key string `storm:"id"`
BlockNumber int64 `json:"block_number" storm:"index"`
OpenBlockNumber int64
ChannelIdentifier common.Hash `json:"channel_identifier"`
ToAddress common.Address `json:"to_address"`
TokenAddress common.Address `json:"token_address"`
Nonce uint64 `json:"nonce"`
Amount *big.Int `json:"amount"`
}
```
Note: This interface does not contain transactions that participate as transit points in the middle.
#### OnNotify
`level` is defined as follows
```go
Type Level int

Const (
// LevelInfo :
LevelInfo = iota
// LevelWarn :
LevelWarn
// LevelError :
LevelError
)
```
Where `info` is the corresponding message, I hope that through this interface, I hope the App can intercept and pop up the corresponding MessageBox.

### Manually registering node information
func (a *API) UpdateMeshNetworkNodes(nodesstr string) (err error)

Manually register a communicable node address to photon
Example data:
```json
[{
   "address": "0x292650fee408320D888e06ed89D938294Ea42f99",
   "ip_port": "127.0.0.1:40001"
},
{
     "address":"0x4B89Bff01009928784eB7e7d10Bf773e6D166066",
    "ip_port":"127.0.0.1:40002"
}
]
```
Tell Photon how to work with 0x292650fee408320D888e06ed89D938294Ea42f99 and 0x4B89Bff01009928784eB7e7d10Bf773e6D166066

## Query interface
### Get the account address of the running photon node
func (a *API) Address() (addr string)

Return example:
`0x7B874444681F7AEF18D48f330a0Ba093d3d0fDD2`
### Get the delegate data provided to a third party
func (a *API) ChannelFor3rdParty(channelIdentifier, thirdPartyAddress string) (r string, err error)

Because Photon works, if a node goes offline for a long time, it will bring its own financial security risk. So if Photon is likely to be offline for a long time (relative to the settlement window time specified by the `OpenChannel`),
Then the relevant proof of income should be delegated to the third-party service (Photon-Monitoring), and the third-party service submits the relevant BalanceProof when needed.
For how to use Photon-Monitoring, please refer to [Photon-Monitoring](https://github.com/SmartMeshFoundation/Photon-Monitoring)

The returned data should be submitted intact to your trusted third party monitoring service.
Return example:
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
### Get the current public link connection status
func (a *API) EthereumStatus() (r string, err error)

Already obsolete, you should use NotifyHandler to get connection state changes.
Return err means no public chain connection, otherwise the connection is normal
### Get all the registered token lists
func (a *API) Tokens() (tokens string)

Return example:
```json
[
    "0x7B874444681F7AEF18D48f330a0Ba093d3d0fDD2"
]
```
### Query all the channels that you participate in under a token.
func (a *API) TokenPartners(tokenAddress string) (channels string, err error)

Return example:
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
### Get channel list
func (a *API) GetChannelList() (channels string, err error)

Return example:
```json
[
    {
        "channel_address": "0xc502076485a3cff65f83c00095dc55e745f790eee4c259ea963969a343fc792a",
        "open_block_number": 5228715,
        "partner_address": "0x4B89Bff01009928784eB7e7d10Bf773e6D166066",
        "balance": 499490,
        "partner_balance": 1500506,
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
### Get information about a channel
func (a *API) GetOneChannel(channelIdentifier string) (channel string, err error)

Return example:
```json
{
    "channel_identifier": "0xc502076485a3cff65f83c00095dc55e745f790eee4c259ea963969a343fc792a",
    "open_block_number": 5228715,
    "partner_address": "0x4B89Bff01009928784eB7e7d10Bf773e6D166066",
    "balance": 499490,
    "patner_balance": 1500506,
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
### Query the list of transactions received
func (a *API) GetReceivedTransfers(from, to int64) (r string, err error)

Convenient App query history transaction, return data is an array of `ReceivedTransfer`
### Query the list of transactions sent
func (a *API) GetSentTransfers(from, to int64) (r string, err error)
Convenient App query history transaction, return data is an array of `SenTransfer`
## Transaction/Channel related interface, asynchronous
### Initiating a transaction
func (a *API) Transfers(tokenAddress, targetAddress string, amountstr string, feestr string, secretStr string, isDirect bool) (transfer string, err error)

Initiate a transaction, asynchronous interface, use the token_address + lockSecretHash returned to call the GetTransferStatus interface to query the transaction status

parameter:
* `tokenAddress string` – transaction token
* `targetAddress string` – payee address
* `amountstr string` – amount
* `feestr string` – the amount of the fee
* `secretStr string` – the transaction password, which can be ""
* `isDirect string` – whether direct channel trading

Return example:
```json
{
    "initiator_address": "0x33Df901ABc22DcB7F33c2a77aD43CC98FbFa0790",
    "target_address": "0x1a9eC3b0b807464e6D3398a59d6b0a369Bf422fA",
    "token_address": "0x4092cE58b448abDFB59fbC84a0E30689F004d02E",
    "amount": 1,
    "lockSecretHash": "0x0676b190e483c6ce425492e45726797c8b538a620a761371d68b1c96e7a8538e",
    "is_direct": true
}
```
### Query the status of the transaction initiated by yourself
func (a *API) GetTransferStatus(tokenAddressStr string, lockSecretHashStr string) (r string, err error)

parameter:
* `tokenAddress string` – transaction token
* `lockSecretHashStr string` – lockSecretHash in the Transfers interface return

Return example:
```json
{
    "Key": "0xf9c7a5491439238ad55c0a8e5a1b97eb205cb14f8137705c898d8d24fcf32465",
    "LockSecretHash": "0x0676b190e483c6ce425492e45726797c8b538a620a761371d68b1c96e7a8538e",
    "TokenAddress": "0x4092ce58b448abdfb59fbc84a0e30689f004d02e",
    "Status": 0,
    "StatusMessage": "DirectTransfer is sending \n"
}
```
The status of the status is as follows:
* 0 - Transfer init
* 1 - transfer can cancel
* 2 - transfer can not cancel
* 3 - transfer already success
* 4 - transfer cancel by user request
* 5 - transfer already failed

### Initiating a token swap transaction
func (a *API) TokenSwap(role string, lockSecretHash string, SendingAmountStr, ReceivingAmountStr string, SendingToken, ReceivingToken, TargetAddress string, SecretStr string) (callID string, err error)

This interface implements the decentralized atomic interchange operation of two Tokens.
This transaction process is generally called by `taker` first, calling the example:
```
    "role": "taker",
    "lockSecretHash": "0x8e90b850fdc5475efb04600615a1619f0194be97a6c394848008f33823a7ee03",
    "TargeAddress": "0x31DdaC67e610c22d19E887fB1937BEE3079B56Cd",
    "SendingAmountStr": 10,
    "SendingToken": "0x7B874444681F7AEF18D48f330a0Ba093d3d0fDD2",
    "ReceivingAmountStr": 100,
    "ReceivingToken": "0x9E7c6C6bf3A60751df8AAee9DEB406f037279C2a"
    "SecretStr":"",
```
Then call it as a maker on another phone, calling the example:
```
    "role": "maker",
    "TargetAddress": "0x69C5621db8093ee9a26cc2e253f929316E6E5b92",
    "SendingAmountStr": 100,
    "SendingToken": "0x9E7c6C6bf3A60751df8AAee9DEB406f037279C2a",
    "ReceivingAmountStr": 10,
    "ReceivingToken": "0x7B874444681F7AEF18D48f330a0Ba093d3d0fDD2",
    "SecretStr": "0x40a6994181d0b98efcf80431ff38f9bae6fefda303f483e7cf5b7de7e341502a",
    "lockSecretHash":""
```
This function will return immediately, and the result of the exchange can be passed to GetCallResult.
Returns a callID used to call the GetCallResult interface to query the result of the call.

### Registering a token with the photon network
func (a *API) RegisterToken(tokenAddress string) (callID string, err error)

Immediately return a callID, used to call the GetCallResult interface to query the call result
### Create a channel
func (a *API) OpenChannel(partnerAddress, tokenAddress string, settleTimeout int, balanceStr string) (callID string, err error)

Note: settleTimeout is the settlement window time, in block units. In actual use, for security reasons, you should set a larger value, such as 600 blocks (9000 seconds, more than two hours).
This means that it takes more than two hours to return to your account without collaborating to close the channel.
Returns a callID used to call the GetCallResult interface to query the result of the call.
### Deposit a corresponding token into a channel
func (a *API) DepositChannel(channelIdentifier string, balanceStr string) (callID string, err error)

Returns a callID used to call the GetCallResult interface to query the result of the call.
### Close a channel
func (a *API) CloseChannel(channelIdentifier string, force bool) (callID string, err error)

`force` is false, it will seek to close the channel with the other party, in the case of consensus, you can immediately  (to wait for one or two blocks) get the token to your account.
`force` is true, it will not negotiate with the other party, which means that the channel will be closed first, then wait for `settleTimeout`   blocks before the `SettleChannel` can be executed, and finally the Token will return to its own account.
Returns a callID used to call the GetCallResult interface to query the result of the call.
### Settle a channel
func (a *API) SettleChannel(channelIdentifier string) (callID string, err error)

Returns a callID used to call the GetCallResult interface to query the result of the call.
### Query result according to callID query
func (a *API) GetCallResult(callID string) (r string, done bool, err error)

return:
* `r string`– interface call returns, example reference [rest_api.md](rest_api.md)
* `err error` – the interface calls an error message, returning a description indicating that the result is being processed yet
