# SmartRaiden’s Mobile API 文档

## 节点管理相关接口
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
* `otherArgs string` – 其他参数,参考smartraiden -h

返回:
* `api *API` – 启动成功返回api句柄
* `err error` – 错误信息

### 停止一个雷电节点
func (a *API) Stop()

### 切换雷电运行环境
func (a *API) SwitchNetwork(isMesh bool)

切换网络环境,Mesh or Internet

### 通知雷电节点网络断开
func (a *API) NotifyNetworkDown() error

主动告知雷电节点网络断开,并让雷电节点开始尝试重连

### 订阅雷电事件
func (a *API) Subscribe(handler NotifyHandler) (sub *Subscription, err error)

订阅雷电节点事件,包含交易通知,错误通知等

### 手动注册节点信息
func (a *API) UpdateMeshNetworkNodes(nodesstr string) (err error)

手动注册一个可通信的节点地址到smartraiden

## 查询接口
### 获取运行雷电节点的账户地址
func (a *API) Address() (addr string)


返回示例:
``0x7B874444681F7AEF18D48f330a0Ba093d3d0fDD2``
### 获取提供给第三方的委托数据
func (a *API) ChannelFor3rdParty(channelIdentifier, thirdPartyAddress string) (r string, err error)

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

返回err表示无公链连接,反之连接正常
### 获取所有token列表
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

### 查询发出的交易列表
func (a *API) GetSentTransfers(from, to int64) (r string, err error)

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
func (a *API) TokenSwap(role string, lockSecretHash string, ...) (callID string, err error)

返回一个callID,用于调用GetCallResult接口查询调用结果

### 向雷电网络注册一个token
func (a *API) RegisterToken(tokenAddress string) (callID string, err error)

返回一个callID,用于调用GetCallResult接口查询调用结果
### 创建一个channel
func (a *API) OpenChannel(partnerAddress, tokenAddress string, settleTimeout int, balanceStr string) (callID string, err error)

返回一个callID,用于调用GetCallResult接口查询调用结果
### 向一个channel里面存入对应token
func (a *API) DepositChannel(channelIdentifier string, balanceStr string) (callID string, err error)

返回一个callID,用于调用GetCallResult接口查询调用结果
### 关闭一个channel
func (a *API) CloseChannel(channelIdentifier string, force bool) (callID string, err error)

返回一个callID,用于调用GetCallResult接口查询调用结果
### 结算一个channel
func (a *API) SettleChannel(channelIdentifier string) (callID string, err error)

返回一个callID,用于调用GetCallResult接口查询调用结果
### 根据callID查询调用结果
func (a *API) GetCallResult(callID string) (r string, done bool, err error)

返回:
* `done bool`– 调用是否完成,当done=true时,r和err才有意义
* `r string`– 接口调用返回,示例参考http接口文档
* `err error`– 接口调用错误信息

