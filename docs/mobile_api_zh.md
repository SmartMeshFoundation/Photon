# SmartRaiden’s Mobile API 文档
<!-- @import "[TOC]" {cmd="toc" depthFrom=1 depthTo=6 orderedList=false} -->

<!-- code_chunk_output -->

* [SmartRaiden’s Mobile API 文档](#smartraidens-mobile-api-文档)
	* [安装](#安装)
		* [android使用](#android使用)
		* [iOS使用](#ios使用)
		* [其他已知问题](#其他已知问题)
	* [节点管理相关接口](#节点管理相关接口)
		* [启动一个雷电节点](#启动一个雷电节点)
		* [停止一个雷电节点](#停止一个雷电节点)
		* [切换雷电运行环境](#切换雷电运行环境)
		* [通知雷电节点网络断开](#通知雷电节点网络断开)
		* [订阅雷电事件](#订阅雷电事件)
			* [OnError](#onerror)
			* [OnStatusChange](#onstatuschange)
			* [OnReceivedTransfer](#onreceivedtransfer)
			* [OnSentTransfer](#onsenttransfer)
			* [OnNotify](#onnotify)
		* [手动注册节点信息](#手动注册节点信息)
	* [查询接口](#查询接口)
		* [获取运行雷电节点的账户地址](#获取运行雷电节点的账户地址)
		* [获取提供给第三方的委托数据](#获取提供给第三方的委托数据)
		* [获取当前公链连接状态](#获取当前公链连接状态)
		* [获取所有已经注册的token列表](#获取所有已经注册的token列表)
		* [查询某个token下自己参与的所有channel](#查询某个token下自己参与的所有channel)
		* [获取通道列表](#获取通道列表)
		* [获取一个通道的信息](#获取一个通道的信息)
		* [查询收到的交易列表](#查询收到的交易列表)
		* [查询发出的交易列表](#查询发出的交易列表)
	* [交易/通道相关接口,异步](#交易通道相关接口异步)
		* [发起一笔交易](#发起一笔交易)
		* [查询自己发起的交易状态](#查询自己发起的交易状态)
		* [发起一笔token swap交易](#发起一笔token-swap交易)
		* [向雷电网络注册一个token](#向雷电网络注册一个token)
		* [创建一个channel](#创建一个channel)
		* [向一个channel里面存入对应token](#向一个channel里面存入对应token)
		* [关闭一个channel](#关闭一个channel)
		* [结算一个channel](#结算一个channel)
		* [根据callID查询调用结果](#根据callid查询调用结果)

<!-- /code_chunk_output -->

## 安装
SmartRaiden mobile SDk编译必须要求gomobile工具可以正常使用. gomobile的安装编译工作请参考[gomobile](https://godoc.org/golang.org/x/mobile)
```bash
cd mobile
#build android
./build_android.sh
#build iOS
./build_iOS.sh
```
### android使用
将mobile.aar 集成到项目即可
### iOS使用
将Mobile.framework集成到项目即可

### 其他已知问题
由于gomobile的工作方式限制,如果项目中同时有两个gomobile编译的sdk(比如你的项目还依赖ethereum的mobile包),程序无法正常运行.

## 节点管理相关接口
SmartRaiden依赖gomobile自动进行接口封装,因为是跨语言调用,无法避免的就是类型转换问题. 
为了规避此类问题,SmartRaiden对外提供接口几乎都是基本类型(int,string,error).

### 启动一个雷电节点
func StartUp(...) (api *API, err error)

参数:
* `address string`– 雷电节点所使用的账户地址
* `keystorePath string` – 账户私钥保存路径
* `ethRPCEndPoint string` – 公链节点host,http协议
* `dataDir string` – smartraiden db路径
* `passwordfile string` – 账户密码文件路径
* `apiAddr string` – http api 监听端口
* `listenAddr string` – udp 监听端口
* `logFile string` – 日志文件路径
* `registryAddress string` – TokenNetworkRegistry合约地址
* `otherArgs mobile.Strings` – 其他参数,参考smartraiden -h   

如果需要传递默认参数以外的其他参数,可以参考如下方式:
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
返回:
* `api *API` – 启动成功返回api句柄
* `err error` – 错误信息

### 停止一个雷电节点
func (a *API) Stop()

### 切换雷电运行环境
func (a *API) SwitchNetwork(isMesh bool)

切换网络环境,Mesh or Internet
在Mesh网络下,节点之间直接使用UDP协议通信,需要App通过UpdateMeshNetworkNodes来告知SmartRaiden其他节点信息.

### 通知雷电节点网络断开
func (a *API) NotifyNetworkDown() error

主动告知雷电节点网络断开,并让雷电节点开始尝试重连

在手机网络环境下,由于网络复杂性,比如WiFi断开等,这些事件SmartRaiden不能直接从系统感知,需要App主动告诉SmartRaiden,让其采取相应的处理.

### 订阅雷电事件
func (a *API) Subscribe(handler NotifyHandler) (sub *Subscription, err error)

订阅雷电节点事件,包含交易通知,错误通知等
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
	// OnNotify get some important message raiden want to notify upper application
	OnNotify(level int, info string)
}
```

#### OnError
 通知SmartRaiden内部发生了不可恢复的错误,SmartRaiden的任何功能必须立即重启才能使用. 由于考虑到SmartRaiden的集成方式可能是单进程方式,我们不希望因为SmartRaiden
 的未知错误导致App闪退,因此即使SmartRaiden内部发生了不可预知的错误,也会有SmartRaiden截获并报告给App,由App来决定是立即退出还是继续使用.
- `errCode`是错误代码
- `failure`是错误信息描述
重启SmartRaiden方式:
```go
api.Stop()
newAPI,err:=Startup(...)
```
#### OnStatusChange
s 是如下结构体的json编码
```go
//ConnectionStatus status of network connection
type ConnectionStatus struct {
	XMPPStatus    netshare.Status
	EthStatus     netshare.Status
	LastBlockTime string
}
```
其中`XMPPStatus`和`EthStatus`定义如下:
```go
// Status shows actual connection status.
type Status int

const (
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
tr是如下结构体的json编码
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
注意: 此接口并不包含作为中间中转节点参与的交易
#### OnSentTransfer
tr是如下结构体的编码
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
注意: 此接口不包含作为中间中转接点参与的交易
#### OnNotify
`level`定义如下
```go
type Level int

const (
	// LevelInfo :
	LevelInfo = iota
	// LevelWarn :
	LevelWarn
	// LevelError :
	LevelError
)
```
其中info为对应的消息,希望通过此接口,希望App能够截获并弹出相应的MessageBox.

### 手动注册节点信息
func (a *API) UpdateMeshNetworkNodes(nodesstr string) (err error)

手动注册一个可通信的节点地址到smartraiden
example data:
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
告诉SmartRaiden如何与0x292650fee408320D888e06ed89D938294Ea42f99和0x4B89Bff01009928784eB7e7d10Bf773e6D166066两个节点进行通信.


## 查询接口
### 获取运行雷电节点的账户地址
func (a *API) Address() (addr string)

返回示例:
``0x7B874444681F7AEF18D48f330a0Ba093d3d0fDD2``
### 获取提供给第三方的委托数据
func (a *API) ChannelFor3rdParty(channelIdentifier, thirdPartyAddress string) (r string, err error)

因为SmartRaiden的工作原理决定了,如果一个节点长时间离线,将会带来自身的资金安全性风险.因此如果SmartRaiden有可能较长时间(相对于创建通道指定的结算窗口时间)离线,
那么应该把相关的收益证明委托给第三方服务(SmartRaiden-Monitoring),由第三方服务在需要的时候提交相关的BalanceProof.
SmartRaiden-Monitoring如何使用,请参考[SmartRaiden-Monitoring](https://github.com/SmartMeshFoundation/SmartRaiden-Monitoring)

返回数据应该原封不动提交给您可信赖的第三方监控服务.
返回示例:
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
### 获取当前公链连接状态
func (a *API) EthereumStatus() (r string, err error)

已经废弃,应该使用NotifyHandler来获取连接状态变化.
返回err表示无公链连接,反之连接正常
### 获取所有已经注册的token列表
func (a *API) Tokens() (tokens string)

返回示例:
```json
[
    "0x7B874444681F7AEF18D48f330a0Ba093d3d0fDD2"
]
```
### 查询某个token下自己参与的所有channel
func (a *API) TokenPartners(tokenAddress string) (channels string, err error)

返回示例:
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
### 获取通道列表
func (a *API) GetChannelList() (channels string, err error)

返回示例:
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
### 获取一个通道的信息
func (a *API) GetOneChannel(channelIdentifier string) (channel string, err error)

返回示例:
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
### 查询收到的交易列表
func (a *API) GetReceivedTransfers(from, to int64) (r string, err error)

方便App查询历史交易,返回数据是`ReceivedTransfer`的数组
### 查询发出的交易列表
func (a *API) GetSentTransfers(from, to int64) (r string, err error)
方便App查询历史交易,返回数据是`SenTransfer`的数组
## 交易/通道相关接口,异步
### 发起一笔交易
func (a *API) Transfers(...) (transfer string, err error)

发起一笔交易,异步接口,使用返回里面的token_address + lockSecretHash调用GetTransferStatus接口查询交易状态

参数:
* `tokenAddress string`– 交易token
* `targetAddress string` – 收款方地址
* `amountstr string` – 金额
* `feestr string` – 手续费金额
* `secretStr string` – 交易密码,可为""
* `isDirect string` – 是否直接通道交易

返回示例:
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
### 查询自己发起的交易状态
func (a *API) GetTransferStatus(tokenAddressStr string, lockSecretHashStr string) (r string, err error)

参数:
* `tokenAddress string`– 交易token
* `lockSecretHashStr string` – Transfers接口返回中的lockSecretHash

返回示例:
```json
{
    "Key": "0xf9c7a5491439238ad55c0a8e5a1b97eb205cb14f8137705c898d8d24fcf32465",
    "LockSecretHash": "0x0676b190e483c6ce425492e45726797c8b538a620a761371d68b1c96e7a8538e",
    "TokenAddress": "0x4092ce58b448abdfb59fbc84a0e30689f004d02e",
    "Status": 0,
    "StatusMessage": "DirectTransfer 正在发送\n"
}
```
其中status取值如下:
* 0 - Transfer init
* 1 - transfer can cancel
* 2 - transfer can not cancel
* 3 - transfer already success
* 4 - transfer cancel by user request
* 5 - transfer already failed

### 发起一笔token swap交易
func (a *API) TokenSwap(role string, lockSecretHash string, SendingAmountStr, ReceivingAmountStr string, SendingToken, ReceivingToken, TargetAddress string, SecretStr string) (callID string, err error)

该接口实现两种Token的去中心化原子互换操作. 
此交易过程一般是现有`taker`调用,调用示例
```
    "role": "taker",
    "lockSecretHash":"0x8e90b850fdc5475efb04600615a1619f0194be97a6c394848008f33823a7ee03",
    "TargeAddress":"0x31DdaC67e610c22d19E887fB1937BEE3079B56Cd",
    "SendingAmountStr": 10,
    "SendingToken": "0x7B874444681F7AEF18D48f330a0Ba093d3d0fDD2",
    "ReceivingAmountStr": 100,
    "ReceivingToken": "0x9E7c6C6bf3A60751df8AAee9DEB406f037279C2a"
    "SecretStr":"",
```
然后在另一台手机上作为maker调用,调用示例:
```
    "role": "maker",
    "TargetAddress":"0x69C5621db8093ee9a26cc2e253f929316E6E5b92",
    "SendingAmountStr": 100,
    "SendingToken": "0x9E7c6C6bf3A60751df8AAee9DEB406f037279C2a",
    "ReceivingAmountStr": 10,
    "ReceivingToken": "0x7B874444681F7AEF18D48f330a0Ba093d3d0fDD2",
    "SecretStr": "0x40a6994181d0b98efcf80431ff38f9bae6fefda303f483e7cf5b7de7e341502a",
    "lockSecretHash":""
```
此函数会立即返回,交换结果可以通过GetCallResult
返回一个callID,用于调用GetCallResult接口查询调用结果

### 向雷电网络注册一个token
func (a *API) RegisterToken(tokenAddress string) (callID string, err error)

立即返回一个callID,用于调用GetCallResult接口查询调用结果
### 创建一个channel
func (a *API) OpenChannel(partnerAddress, tokenAddress string, settleTimeout int, balanceStr string) (callID string, err error)

注意: settleTimeout就是结算窗口时间,以块为单位,在实际使用过程中,为了安全起见,应该设置一个较大的值,比如600块(9000秒两个多小时). 
这意味着在不合作关闭通道的情况下,需要两个多小时Token才能返回自己的账户. 
返回一个callID,用于调用GetCallResult接口查询调用结果
### 向一个channel里面存入对应token
func (a *API) DepositChannel(channelIdentifier string, balanceStr string) (callID string, err error)

返回一个callID,用于调用GetCallResult接口查询调用结果
### 关闭一个channel
func (a *API) CloseChannel(channelIdentifier string, force bool) (callID string, err error)

force 为false,则会寻求和对方协商关闭通道,在协商一致的情况下可以立即(等待一两个块的时间)将Token返回到自己账户
force 为true,则不会与对方协商,意味着会首先关闭通道,然后等待`settleTimeout`这么多块,然后才可以进行SettleChannel,最终Token才会返回自己的账户
返回一个callID,用于调用GetCallResult接口查询调用结果
### 结算一个channel
func (a *API) SettleChannel(channelIdentifier string) (callID string, err error)

返回一个callID,用于调用GetCallResult接口查询调用结果
### 根据callID查询调用结果
func (a *API) GetCallResult(callID string) (r string, done bool, err error)

返回:
* `r string`– 接口调用返回,示例参考http接口文档
* `err error`– 接口调用错误信息,返回dealing说明正在处理尚未收到结果
