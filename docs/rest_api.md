# SmartRaiden REST API Reference
欢迎使用SmartRaiden REST API Reference,这是一份`v1`版本的粗略api参考文档，供开发者尝鲜使用（后续会不断更新完善)

文档主要介绍几个大类：

- Token 
  
- Channel
  
- Transfer


## Token  
### `GET /api/1/tokens`
查询已经注册的token ,没有请求参数  

**Example Response:**  
```json
[
    "0x7B874444681F7AEF18D48f330a0Ba093d3d0fDD2"
]
```
### `PUT /api/1/tokens/<token_address>`  
注册新的token  

**Example Request:**

`PUT /api/1/tokens/0x9E7c6C6bf3A60751df8AAee9DEB406f037279C2a`



**Example Response:**  
```json
{
    "channel_manager_address": "0xBb1e95363b0181De7bBf394f18eaC7D4230e391A"
}
```
   
### `GET /api/1/address`     
查询你的smartraiden地址

**Example Response:**  
```json
{
    "our_address": "0xf0f6E53d6bbB9Debf35Da6531eC9f1141cd549d5"
}
```

## Channel  
### Channel Object
```json
{
    "channel_address": "0x47235d9d81eb6c19dea2b695b3d6ba1cf76c169d329dc60d188390ba5549d025",
    "open_block_number": 2996350,
    "partner_address": "0x31DdaC67e610c22d19E887fB1937BEE3079B56Cd",
    "balance": 500,
    "partner_balance": 100,
    "locked_amount": 0,
    "partner_locked_amount": 0,
    "token_address": "0xF2747ea1AEE15D23F3a49E37A146d3967e2Ea4E5",
    "state": 1,
    "StateString": "opened",
    "settle_timeout": 150,
    "reveal_timeout": 5
}
```

State|StateString|Description
--|--|--
0|InValid|无效的通道
1|Opened|可以正常交易
2|Closed|不能再发起交易了,还可以接受交易
3|BalanceProofUpdated|已经提交过证据,未完成的交易不再继续,不能接收 unlock 消息
4|Settled|通道已经彻底结算,和 invalid 状态意义相同
5|Closing|StateClosing 用户发起了关闭通道的请求,正在处理正在进行交易,可以继续,不再新开交易
6|Settling|StateSettling 用户发起了 结算请求,正在处理正常情况下此时不应该还有未完成交易，不能新开交易,正在进行的交易也没必要继续了.因为已经提交到链上了
7|Withdraw|StateWithdraw 用户收到或者发出了 withdraw 请求,这时候正在进行的交易只能立即放弃,因为没有任何意义了
8|CooprativeSettle|StateCooprativeSettle 用户收到或者发出了 cooperative settle 请求,这时候正在进行的交易只能立即放弃,因为没有任何意义了
9|PrepareForCooperativeSettle|StatePrepareForCooperativeSettle 收到了用户 cooperative 请求,但是有正在处理的交易,这时候不再接受新的交易了,可以等待一段时间,然后settle已开始交易,可以继续
10|PrepareForWithdraw|StatePrepareForWithdraw 收到用户请求,要发起 withdraw, 但是目前还持有锁,不再发起或者接受任何交易,可以等待一段时间进行 withdraw已开始交易,可以继续
11|Error|StateError 比如收到了明显错误的消息,又是对方签名的,如何处理?比如自己未发送 withdrawRequest,但是收到了 withdrawResponse。todo 这种情况应该的实现是关闭通道.这样真的合理吗?


### `GET /api/1/channels` 
查询节点所有未结算的通道   
 
**Example Response:**  
```json
[
    {
        "channel_address": "0xc943251676c4e53b2669fbbf17ebcbb850da9cb0a907200c40f1342a37629489",
        "open_block_number": 2560169,
        "partner_address": "0x31DdaC67e610c22d19E887fB1937BEE3079B56Cd",
        "balance": 100,
        "partner_balance": 100,
        "locked_amount": 0,
        "partner_locked_amount": 0,
        "token_address": "0x7B874444681F7AEF18D48f330a0Ba093d3d0fDD2",
        "state": 1,
        "StateString": "opened",
        "settle_timeout": 150,
        "reveal_timeout": 5
    }
]
```

### `POST /api/1/channels`  
开启一个通道  

**PAYLOAD:**
```json
{
    "partner_address": "0xf0f6E53d6bbB9Debf35Da6531eC9f1141cd549d5",
    "token_address": "0x7B874444681F7AEF18D48f330a0Ba093d3d0fDD2",
    "balance": 50,
    "settle_timeout": 150
}
```

**Example Response:**  

```json 
{
    "channel_address": "0x97f73562938f6d538a07780b29847330e97d40bb8d0f23845a798912e76970e1",
    "open_block_number": 2560271,
    "partner_address": "0xf0f6E53d6bbB9Debf35Da6531eC9f1141cd549d5",
    "balance": 50,
    "partner_balance": 0,
    "locked_amount": 0,
    "partner_locked_amount": 0,
    "token_address": "0x7B874444681F7AEF18D48f330a0Ba093d3d0fDD2",
    "state": 1,
    "StateString": "opened",
    "settle_timeout": 150,
    "reveal_timeout": 0
}
```
### `GET/api/1/channels/<channel_address>`  
查询特定的通道 ,可以看到通道的详细信息

**Example Request:**

`GET /api/1/channels/0xc943251676c4e53b2669fbbf17ebcbb850da9cb0a907200c40f1342a37629489`  

**Example Response:**
```json
{
    "channel_address": "0xc943251676c4e53b2669fbbf17ebcbb850da9cb0a907200c40f1342a37629489",
    "open_block_number": 2899911,
    "partner_address": "0x31DdaC67e610c22d19E887fB1937BEE3079B56Cd",
    "balance": 80,
    "patner_balance": 120,
    "locked_amount": 0,
    "partner_locked_amount": 0,
    "token_address": "0x7B874444681F7AEF18D48f330a0Ba093d3d0fDD2",
    "state": 1,
    "StateString": "opened",
    "settle_timeout": 150,
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
        "Nonce": 2,
        "TransferAmount": 20,
        "LocksRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "ChannelIdentifier": {
            "ChannelIdentifier": "0xc943251676c4e53b2669fbbf17ebcbb850da9cb0a907200c40f1342a37629489",
            "OpenBlockNumber": 2899911
        },
        "MessageHash": "0x93a656c5b673759c76083439790a9f7b91c7656b41ef8884e098517e15461427",
        "Signature": "BCspERU5NQvgm3zB55mK/YWRBErqhgcPiGZMVgIfgz1bzO7iplEOQ/An6F8cLIXMt06RjQmsfOc4yjWRDFSzYBw=",
        "ContractTransferAmount": 0,
        "ContractNonce": 2,
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

### `PATCH /api/1/channels/<channel_address>` 

`PATCH /api/1/channels/0xc943251676c4e53b2669fbbf17ebcbb850da9cb0a907200c40f1342a37629489`

向一个通道里面存钱  

**PAYLOAD:**   
```json
{
    "balance": 100
}
```
**Example Response:**  
```json
{
    "channel_address": "0xc943251676c4e53b2669fbbf17ebcbb850da9cb0a907200c40f1342a37629489",
    "open_block_number": 2560169,
    "partner_address": "0x69C5621db8093ee9a26cc2e253f929316E6E5b92",
    "balance": 100,
    "partner_balance": 100,
    "locked_amount": 0,
    "partner_locked_amount": 0,
    "token_address": "0x7B874444681F7AEF18D48f330a0Ba093d3d0fDD2",
    "state": 1,
    "StateString": "opened",
    "settle_timeout": 150,
    "reveal_timeout": 5
}
```
`PATCH /api/1/channels/0x97f73562938f6d538a07780b29847330e97d40bb8d0f23845a798912e76970e1` 

关闭一个通道,参数`force`默认为`false`，表示合作结算通道。  

**PAYLOAD:**  
```json
{"state":"closed"，
  "force":false
	
}
```
**Example Response:**   
```json
{
    "channel_address": "0xf1fa19fa6a54912e32d6e6e1aa0baa14d530385c60266886ef7c18838f6e9bdc",
    "open_block_number": 2498052,
    "partner_address": "0x6B9E4D89EE3828e7a477eA9AA7B62810260e27E9",
    "balance": 0,
    "partner_balance": 0,
    "locked_amount": 0,
    "partner_locked_amount": 0,
    "token_address": "0x7B874444681F7AEF18D48f330a0Ba093d3d0fDD2",
    "state": 8,
    "StateString": "cooperativeSettling",
    "settle_timeout": 35,
    "reveal_timeout": 5
}
```
当通道对方不在线时，或者不想合作结算通道，可将`force`设置为`true`,等待`settle_timeout`后再结算 

**PAYLOAD：**  
```json
{"state":"closed",
  "force":true
}
```
**Example Response:**  
```json 
{
    "channel_address": "0xc943251676c4e53b2669fbbf17ebcbb850da9cb0a907200c40f1342a37629489",
    "open_block_number": 2560169,
    "partner_address": "0x69C5621db8093ee9a26cc2e253f929316E6E5b92",
    "balance": 100,
    "partner_balance": 100,
    "locked_amount": 0,
    "partner_locked_amount": 0,
    "token_address": "0x7B874444681F7AEF18D48f330a0Ba093d3d0fDD2",
    "state": 2,
    "StateString": "closed",
    "settle_timeout": 150,
    "reveal_timeout": 5
}
```

`PATCH /api/1/channels/0x97f73562938f6d538a07780b29847330e97d40bb8d0f23845a798912e76970e1`   
结算通道，当通道已经关闭且`settle_timeout`已过，可结算通道  

**PAYLOAD:**  
```json
{
    "state":"settled"
}
```
**Example Response:**  
```json

{
    "channel_address": "0xc943251676c4e53b2669fbbf17ebcbb850da9cb0a907200c40f1342a37629489",
    "open_block_number": 2575160,
    "partner_address": "0x69C5621db8093ee9a26cc2e253f929316E6E5b92",
    "balance": 100,
    "partner_balance": 50,
    "locked_amount": 0,
    "partner_locked_amount": 0,
    "token_address": "0x7B874444681F7AEF18D48f330a0Ba093d3d0fDD2",
    "state": 1,
    "StateString": "settled",
    "settle_timeout": 150,
    "reveal_timeout": 5
}
```

## Transfer
### `transfer/<token_address>/<target_address>`
`POST /api/1/transfers/0x7B874444681F7AEF18D48f330a0Ba093d3d0fDD2/0xf2234A51c827196ea779a440df610F9091ffd570`  

当通道是`open`状态且资金充足的情况下，可以进行转账

**PAYLOAD**
```json
{
    "amount":20,
    "fee":0, //收费金额
    "is_direct":false //是否直接转账
   
}
```
**Example Response:** 
```json
{
    "initiator_address": "0x69C5621db8093ee9a26cc2e253f929316E6E5b92",
    "target_address": "0xf2234A51c827196ea779a440df610F9091ffd570",
    "token_address": "0x7B874444681F7AEF18D48f330a0Ba093d3d0fDD2",
    "amount": 20,
    "secret": "",
    "fee": 0,
    "is_direct": false
}
```
## Token Swap
### `/api/1/token_swaps/<target_address>/<lock_secret_hash>` 
Token Swap 可以用来进行两种token的交换，在保证有效路由的情况下，先调用`taker`再调用`maker`，可通过接口`/api/1/secret/`获取一对`lock_secret_hash`和`secret` 

`PUT /api/1/token_swaps/0x31DdaC67e610c22d19E887fB1937BEE3079B56Cd/0x8e90b850fdc5475efb04600615a1619f0194be97a6c394848008f33823a7ee03`  

**PAYLOAD**
```json
{
    "role": "taker",
    "sending_amount": 10,
    "sending_token": "0x7B874444681F7AEF18D48f330a0Ba093d3d0fDD2",
    "receiving_amount": 100,
    "receiving_token": "0x9E7c6C6bf3A60751df8AAee9DEB406f037279C2a"
}
```

**Example Response:** 

`201 Created`

`PUT /api/1/token_swaps/0x69C5621db8093ee9a26cc2e253f929316E6E5b92/0x8e90b850fdc5475efb04600615a1619f0194be97a6c394848008f33823a7ee03`  

**PAYLOAD** 
```json
{
    "role": "maker",
    "sending_amount": 100,
    "sending_token": "0x9E7c6C6bf3A60751df8AAee9DEB406f037279C2a",
    "receiving_amount": 10,
    "receiving_token": "0x7B874444681F7AEF18D48f330a0Ba093d3d0fDD2",
    "secret": "0x40a6994181d0b98efcf80431ff38f9bae6fefda303f483e7cf5b7de7e341502a"
}
```

**Example Response:** 

`201 Created`

查看对应通道的token余额，会发现token swap 成功

获取一组`lock_secret_hash`和`secret`  

`GET /api/1/secret`


**Example Response:** 
```json
{
    "lock_secret_hash": "0x8e90b850fdc5475efb04600615a1619f0194be97a6c394848008f33823a7ee03",
    "secret": "0x40a6994181d0b98efcf80431ff38f9bae6fefda303f483e7cf5b7de7e341502a"
}
```