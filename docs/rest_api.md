# Photon REST API Reference  

Hey guys, welcome to Photon REST API Reference page. This is an API Spec for Photon version 1.0, which adds a lot more new features, such as, support multi-token functions, channel charging,etc. Please note that this reference is still updating. If any problem, feel free to submit at our Issue.

##  Channel Structure  
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
        "reveal_timeout": 5
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

##  Query node address

 `GET /api/1/address`

 Return the address of photon node.

**Example Request :** 

`GET  http://{{ip1}}/api/1/address`

**Example Response :**  
```json
{
    "our_address": "0x69C5621db8093ee9a26cc2e253f929316E6E5b92"
}
```
**Status Codes:**  
- `200 OK` 


##  Query the registered token
 ` GET /api/1/tokens` 
  
  Return the token address which can be used in the offchain transfer.

**Example Request :**  

`GET  http://{{ip1}}/api/1/tokens`

**Example Response :**  
```json
[
     "0xC07D1D6e8F20F2a90B205762a0BAC0B611c490DC",
    "0x2a7Af974B7bB88703180d6AFF9a656BB4Dbba809",
    "0x8B916406c1ecCC5B15865b7BC7aF5fA90c01Fc59",
    "0x489CEE6beAA894898d0890f4c6d750cA3D8176A4"
]
```
**Status Codes:**  
- `200 OK` 
If the node has not registered the token, then respond the message " NULL".

## Get all the channel partners of this token
 `GET /api/1/tokens/*(token_address)*/partners`

   Return all channels of this node under this token.

**Example Request :**  

`GET http://{{ip2}}/api/1/tokens/0x2a7Af974B7bB88703180d6AFF9a656BB4Dbba809/partners`

**Example Response: **
```json
[
    {
        "partner_address": "0xd5dC7504e0b448b1c62D86306AE8e4a5836Fc1A1",
        "channel": "api/1/channles/0x019ed640b5c6f8a714a77a754e793cd162df164f7e96f88a2beefbd1c576980d"
    },
    {
        "partner_address": "0xC445a8C326A8fD5a3e250C7dc0EFc566eDcB263B",
        "channel": "api/1/channles/0x081f7a9771994de9f06edb52cb60a0fe3b9bbebd4c1240c267967c7e3fa433f5"
    },
    {
        "partner_address": "0x97Cd7291f93F9582Ddb8E9885bF7E77e3f34Be40",
        "channel": "api/1/channles/0xf25edb59e35544e060ecfcef6e6a0ba619ff905a132295957e11ffdc2206fc24"
    }
]
```
**Status Codes:**  
- `200 OK` 

 If the node has not created the channel with the token, then respond the message " NULL".

 
## Query all the channels of the node
   `GET /api/1/channels`  

Return all the unsettled channels of the node.

**Example Request :**  

 `GET  http://{{ip1}}/api/1/channels`

**Example Response :**  
```json
[
     {
        "channel_identifier": "0x8b48df693d6ceeb40c6285b9820171e204d2218088f506b6d8dd415ef690edd7",
        "open_block_number": 14495292,
        "partner_address": "0x3bC7726c489E617571792aC0Cd8b70dF8A5D0e22",
        "balance": 100000000000000000000,
        "partner_balance": 100000000000000000000,
        "locked_amount": 0,
        "partner_locked_amount": 0,
        "token_address": "0x8B916406c1ecCC5B15865b7BC7aF5fA90c01Fc59",
        "state": 1,
        "state_string": "opened",
        "settle_timeout": 100,
        "reveal_timeout": 30
    },
    {
        "channel_identifier": "0xc704dad871ed767dd2ed3d40c7ed2db6c047e82749613bf1458d7ab0a65ba4f1",
        "open_block_number": 14495300,
        "partner_address": "0x3bC7726c489E617571792aC0Cd8b70dF8A5D0e22",
        "balance": 0,
        "partner_balance": 100000000000000000000,
        "locked_amount": 0,
        "partner_locked_amount": 0,
        "token_address": "0x489CEE6beAA894898d0890f4c6d750cA3D8176A4",
        "state": 1,
        "state_string": "opened",
        "settle_timeout": 100,
        "reveal_timeout": 30
    },
    {
        "channel_identifier": "0xf25edb59e35544e060ecfcef6e6a0ba619ff905a132295957e11ffdc2206fc24",
        "open_block_number": 14660842,
        "partner_address": "0x3bC7726c489E617571792aC0Cd8b70dF8A5D0e22",
        "balance": 7987,
        "partner_balance": 12013,
        "locked_amount": 0,
        "partner_locked_amount": 0,
        "token_address": "0x2a7Af974B7bB88703180d6AFF9a656BB4Dbba809",
        "state": 1,
        "state_string": "opened",
        "settle_timeout": 100,
        "reveal_timeout": 30
    }
]
```
**Status Codes:**  
- `200 OK` 

If the node has not created the channel with other nodes, then respond the message " NULL".

## Query specific channel of the node
  `GET /api/1/channels/*(channel_identifier)* `

Query the specific channel and return all the information about the channel.

**Example Request :**  

`GET http://{{ip1}}/api/1/channels/0xf25edb59e35544e060ecfcef6e6a0ba619ff905a132295957e11ffdc2206fc24`

**Example Response :**  
```json
{
    "channel_identifier": "0xc943251676c4e53b2669fbbf17ebcbb850da9cb0a907200c40f1342a37629489",
    "open_block_number": 2899911,
    "partner_address": "0x31DdaC67e610c22d19E887fB1937BEE3079B56Cd",
    "balance": 80000000000000000000,
    "patner_balance": 120000000000000000000,
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
        "TransferAmount": 20000000000000000000,
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
**Status Codes:**  
- `200 OK` 
- `404 Not Found` - not found

## Deposit to the channel
 `  PUT /api/1/deposit `

   Deposit to the channel (if there is no channel, the interface can be reused to create the channel and deposit).

Parameter |type |SON format|description
--|--|--|--
partnerAddress|string|partner_address|The address of the partner
tokenAddress|string|token_address| which token to deposit
settleTimeout|string|settle_timeout|The settlement window period 
balanceStr|big.Int|balance|The deposited amount which must be greater than 0.
newChannel|bool|new_channel|Judge whether the channel exists or not.If the channel doesnot exist, `deposit` will create a new channel and deposit, else only deposit.  

deposit interfaces contain two behaviors：

1. Create channel and deposit
    - `new_channel`sets`true`, which means open a new channel and deposit to the channel;if there is no channel between the participants, `false`is no meaning for `new_channel`，which will return the error message "There is no channel".
    - `settle_timeout`represent the settlement window for new channel, for example,settle_timeout：100; if the `settle_timeout` set to 0,the default window period is used which is 600 block.

2. only deposit:
   - `new_channel`must set to `false`，which means the channel has been existed；If the channel has been existed ,there is Meaningless to set the `new_channel`statue as `ture`，which will response the error message "The channel has already existed". 
   - `settle_timeout`must set to 0,because the channel has already existed.

 **Example Request :**  

  `PUT http://{{ip1}}/api/1/deposit`

**PAYLOAD:** 
```json 
{
    "partner_address": "0x7d289f1cBd70d5c3c6F56c09f812F6407f6458B7",
    "token_address": "0xadE88bC1519867e7091f83D763cf61918d50244a",
    "balance": 10000000000000000000000,
    "settle_timeout": 100,
    "new_channel": true
}
```

**Example Response :**  
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




## Withdraw from the channel  

` PUT /api/1/withdraw/*(channel_identifier)* `

CooperateWithdraw available when both channel participants online.
When you’re ready to withdraw, you can switch the channel state to `"preparewithdraw"` by setting the `"op"`:`"preparewithdraw"` and refuse to accept the transaction.When no new block is received from the connection point for more than one minute, an error message will be given when calling `withdraw`."call smc SyncProgress err, client is closed”,which means that the connection point need to synchronize new blocks.

 **Example Request :**  

`PUT http://{{ip2}}/api/1/withdraw/0x081f7a9771994de9f06edb52cb60a0fe3b9bbebd4c1240c267967c7e3fa433f5`

**PAYLOAD:**  
```json
{
		"op":"preparewithdraw"
}
```
**Example Response :**  

```json
{
    "channel_identifier": "0x623c5bf569977f6da37ff39da9a917eb500089ba7ae95ee894b9349db4320b16",
    "open_block_number": 4135231,
    "partner_address": "0x31DdaC67e610c22d19E887fB1937BEE3079B56Cd",
    "balance": 100000000000000000000,
    "partner_balance": 200000000000000000000,
    "locked_amount": 0,
    "partner_locked_amount": 0,
    "token_address": "0xc0dfdD7821c762eF38F86225BD45ff4e912fFA20",
    "state": 9,
    "StateString": "prepareForWithdraw",
    "settle_timeout": 150,
    "reveal_timeout": 30
}
```
When you want to cancel the state of the `preparewithdraw`, you can switch the channel state to the`opened` through the parameter`"op":"cancelprepare"`.

**PAYLOAD:**   
```json
{
		"op":"cancelprepare"
}
```
**Example Response: **
```json
{
    "channel_identifier": "0x623c5bf569977f6da37ff39da9a917eb500089ba7ae95ee894b9349db4320b16",
    "open_block_number": 4135231,
    "partner_address": "0x31DdaC67e610c22d19E887fB1937BEE3079B56Cd",
    "balance": 100000000000000000000,
    "partner_balance": 200000000000000000000,
    "locked_amount": 0,
    "partner_locked_amount": 0,
    "token_address": "0xc0dfdD7821c762eF38F86225BD45ff4e912fFA20",
    "state": 9,
    "StateString": "opened",
    "settle_timeout": 150,
    "reveal_timeout": 30
}
```
Of course, as long as both channels are online and there is no lock, then you can directly `withdraw`, `op` parameters are not necessary.
When `amount`is greater than 0, the `op` parameter is meaningless.

**PAYLOAD:**     
```json
{
	"amount":50000000000000000000,

}
```
**Example Response :**  
```json
{
    "channel_identifier": "0x47235d9d81eb6c19dea2b695b3d6ba1cf76c169d329dc60d188390ba5549d025",
    "open_block_number": 3613578,
    "partner_address": "0x31DdaC67e610c22d19E887fB1937BEE3079B56Cd",
    "balance": 190000000000000000000,
    "partner_balance": 100000000000000000000,
    "locked_amount": 0,
    "partner_locked_amount": 0,
    "token_address": "0xF2747ea1AEE15D23F3a49E37A146d3967e2Ea4E5",
    "state": 7,
    "StateString": "withdrawing",
    "settle_timeout": 150,
    "reveal_timeout": 30
}
```
If the withdrawn amount is larger than the available balance of the channel, an error message will be returned.such as "Error": "invalid withdraw amount, availabe=399999999999999999999,want=1000000000000000000000"”.

##  Close the channel
`PATCH /api/1/channels/*(channel_identifier)* `

Close the channel, which includes the unilateral close the channel and cooperative settle the channel.
set `force` default to `false`, meaning that channel participants cooperate settle the channel.When no new block is received from the connection point for more than one minute, an error message will be given when calling to close the channel."call smc SyncProgress err, client is closed”,which means that the connection point need to synchronize new blocks.

**Example Request :**    

`PATCH /api/1/channels/0x97f73562938f6d538a07780b29847330e97d40bb8d0f23845a798912e76970e1`
   
**PAYLOAD:**      
```json
{"state":"closed"，
  "force":false
}
```
**Example Response :**    
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
    "settle_timeout": 35,
    "reveal_timeout": 5
}
```
Once channel partner is offline or has the locks, the cooperate settle can't be carried out.The participant should alter the`force` to `true`, wait for settle_timeout and unilateral settle the channel.

**PAYLOAD:**     
```json 
{"state":"closed",
  "force":true
}
```
**Example Response :**    
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
##  Settle the Channel
`PATCH /api/1/channels/(channel_identifier)`

The interface of unilaterally settling channel is reused with closing channel, which the parameters are different.
 
 After unilaterally closing the channel, it is necessary to call the settlement channel to settle the closed channel.Once the half of the settle_timeout block has passed,the PMS can submit the balanceproof of the delegate and unlock the registered transaction,and the channel participants can submit the balanceproof to undate the channel and unlock the registered transaction at any time during the settlement window period.
 Tips:When no new block is received from the connection point for more than one minute, an error message will be given when calling to settle the channel."call smc SyncProgress err, client is closed”,which means that the connection point need to synchronize new blocks.

Note: Since settle_timeout does not include the penalty period (in spectrum, which is 257  block, about an hour), the actual settlement time is about 410 block.

**Example Request :**  

`PATCH /api/1/channels/0x97f73562938f6d538a07780b29847330e97d40bb8d0f23845a798912e76970e1`   

**PAYLOAD:**   
```json
{
    "state":"settled"
}
```

**Example Response :**  
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
    "state": 1,
    "StateString": "settled",
    "settle_timeout": 150,
    "reveal_timeout": 30
}
```
**Status Codes :**  
- `200 OK` - close/settle success
- `409 Conflict` - State conflicts, such as, "failed to estimate gas needed: gas required exceeds allowance or always failing transaction",or "channel is still open".

## Initiate the payment
`POST /api/1/transfers/*(token_address)*/*(target_address)*`  

This interface is used to initiate a transfer transaction, which is currently associated with PFS by default.

**Example Request :**   

`POST /api/1/transfers/0x7B874444681F7AEF18D48f330a0Ba093d3d0fDD2/0xf2234A51c827196ea779a440df610F9091ffd570`

**PAYLOAD:**     
```json
{
    "amount":200000000000000000000000,
    "fee":0,
    "is_direct":false,
    "Sync":false,
    "data":"hello word"
}
```
**Parameter implication:** 
- amount: Transfer amount 
- fee: Specify the total cost of the transaction(When using specified transaction costs, the fee calculated by PFS is not used to send the transfer amount)
- is_direct: whether it is a direct transfer. The default is false(MediatedTransfer)
- Sync: whether it is a sync or not. The default is false,that is,  after a transaction is initiated, it immediately returns the `lockSecretHash` of the transaction.
- data: Incidental information of the transaction. The length is not more than 256 byte.

**Example Response :**    
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


## Initiate the transfer with specified secret

The normal transfer secret is automatically generated by photon. If the user wants to precisely control the success or failure of the transaction, he can use the transfer of the specified `secret`. Currently a major application scenario is tokenswap.

**Example Request :**   

`POST: http://{{ip1}}/api/1/transfers/0xF2747ea1AEE15D23F3a49E37A146d3967e2Ea4E5/0xf0f6E53d6bbB9Debf35Da6531eC9f1141cd549d5`   

```json
{
    "amount":20000000000000000000,
    "is_direct":false,
    "secret":"0xad96e0d02aa2f4db096e3acdba0831f95bb09d876a5c6f44bc3f7325a0a45ea1"
}
```
Note: The specified secret is obtained by the interface  `/api/1/secret`.

## Get secret
` GET /api/1/secret `

Through calling the interface,the caller will Get `lock_secret_hash` / `secret` pair,which can be used in Specified secret transaction or tokenswap.

 **Example Request :**   

 `GET  http://{{ip1}}/api/1/secret`

**Example Response: **
```json
{
    "lock_secret_hash": "0x8e90b850fdc5475efb04600615a1619f0194be97a6c394848008f33823a7ee03",
    "secret": "0x40a6994181d0b98efcf80431ff38f9bae6fefda303f483e7cf5b7de7e341502a"
}
```

## Allow disclosure of secret
 `Post /api/1/transfers/allowrevealsecret`

This interface is used in combination with the specified secret transfer interface. When a transfer with the special secret is sent, if the interface is not be called to unlock the Secret, the initiator will not accept the SecretRequest from the recipient. So when sending the transaction with the specified secret, the sender must actively call this interface ,then the transaction can be successfully continued.
  
  **Example Request :**   

`Post  http://{{ip1}}/api/1/transfers/allowrevealsecret`

**PAYLOAD:**   
```json
{
	"lock_secret_hash":"0xd575975dc6fe745b4abee09804b8b97c16dc9842035d39cf474041315374ef02",
	"token_address":"0x37346b78de60f4F5C6f6dF6f0d2b4C0425087a06"
}

```

- lock_secret_hash: Refers to the lock_secret_hash corresponding to Secret in the send Secret transaction.
- token_address: Token of transactions

**Example Response:**  
**200 ok**

## Query the sent successful transfer 
  `GET /api/1/querysenttransfer` 

For the sender of the transfer, the interface can be used to query the history information of all successful transfer which sent from itself, so that the user can accurately master the situation of the transfered funds. If there is too much history, you can use block filtering.Such as:".../querysenttransfer?from_block=3000&to_block=5000"

**Example Request :**   

`GET http://{{ip1}}/api/1/querysenttransfer`

**Example Response: **
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
## Query the received successful transfer 
   `GET /api/1/queryreceivedtransfer`
   
For the receiver of the transfer, the interface can be used to query the history information of all successful transfer which received from other partners, so that the user can accurately master the situation of the received funds. If there is too much history, you can use block filtering.Such as:".../queryreceivedtransfer?from_block=3000&to_block=5000"

**Example Request：**

`GET http://{{ip2}}/api/1/queryreceivedtransfer`

**Example Response : ** 
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
##  Query the transaction that have not yet been received
   ` GET /api/1/getunfinishedreceivedtransfer/*(tokenaddress)*/*(locksecrethash)* `  

 This interface is called by the receiver, also used to specify the secret transaction scenario. The  receiver can find out that  a transaction has been received through the interface, but there is no secret. The receiver can request the sender to call allowrevealsecret to complete the transaction, otherwise the transaction will be returned after expiration.
 
 **Example Request :**    
`GET /api/1/getunfinishedreceivedtransfer/0xD82E6be96a1457d33B35CdED7e9326E1A40c565D/0x2fb55cec26a26d0212cf6bd6022aaa7426410916de09133be3b353ac1a91d843`   

 **Example Response:  **
```json
{
    "initiator_address": "0x201B20123b3C489b47Fde27ce5b451a0fA55FD60",
    "target_address": "0x151E62a787d0d8d9EfFac182Eae06C559d1B68C2",
    "token_address": "0xD82E6be96a1457d33B35CdED7e9326E1A40c565D",
    "amount": 30000000000000000000,
    "secret": "",
    "lock_secret_hash": "0x2fb55cec26a26d0212cf6bd6022aaa7426410916de09133be3b353ac1a91d843",
    "expiration": 131,
    "is_direct": false
}
```

## Query the transaction status
`GET /api/1/transferstatus/*(token_address)*/*(locksecrethash)* ` 

There are two ways for users to send and receive transactions, that is, synchronous and asynchronous. If the asynchronous mode is used (sync is false, that is, the default mode), the interface can be called to query the status information of the current transaction. Among them, locksecrethash is obtained from the message returned by the asynchronous transfer transaction.

**Example Request :**  

`GET /api/1/transferstatus/0xD82E6be96a1457d33B35CdED7e9326E1A40c565D/0xdb0d663a82d04fedf4f558f75d7be801ab6707ea765662919063bad93cd71c82` 

**Example Response :**  
```json
{
    "LockSecretHash": "0xdb0d663a82d04fedf4f558f75d7be801ab6707ea765662919063bad93cd71c82",
    "Status": 0,
    "StatusMessage": "MediatedTransfer is sending target=151e\nMediatedTransfer sending success\n"
}
```

**Response JSON Array of Objects :**
- `Status`  
  - 0 -TransferStatusInit                    init  
  - 1 -TransferStatusCanCancel               transfer can cancel right now  
  - 2 -TransferStatusCanNotCancel            transfer can not cancel     
  - 3 -TransferStatusSuccess                 transfer already success  
  - 4 -TransferStatusCanceled                transfer cancel by user request 
  - 5 -TransferStatusFailed                  transfer already failed

## cancel the transaction
  ` Post /api/1/transfercancel/*(token)*/*(locksecrethash)*`

This interface is for cancellation of a transaction, which the transaction is in a cancelable state.

In the asynchronous transaction transfer process, if the current transaction is in the cancelable state (status code is 1) through the transaction status query, and the waiting time is too long, the interface can be used to cancel the transaction.

**Example Request :**  

`POST /api/1/transfercancel/0xD82E6be96a1457d33B35CdED7e9326E1A40c565D/0xe0f8d65ddb4f70899b97f36795925a97c1b286582f58f56a041f141d345acdca`

**Example Response :**  
**200 OK**

Note: Before using this interface, you need to query the corresponding transaction status through the interface `/api/1/transferstatus`. If it is not in the cancelable state, the interface will return an Error:"can not found transfer".

## Token exchange
  ` PUT /api/1/token_swaps/*(target_address)*/*(lock_secret_hash)*`

Token Swap can be used to atomic exchange within two types of tokens.


 Under the circumstances that valid routing strategies are existed, first invoke `taker` then `maker`,  It should be noted that the preimage of the `lock_secret_hash` which the maker introduced must be equal to the `secret` when taker adopted. Note that both taker and maker request the lock_secret_hash, the secret was given in the  maker's request parameters.

Note: With help of the interface  `/api/1/secret` ,  channel participants can receive a `lock_secret_hash` / `secret` pair.

**Example Request :**    

**the taker:  **
`PUT /api/1/token_swaps/0x31DdaC67e610c22d19E887fB1937BEE3079B56Cd/0x8e90b850fdc5475efb04600615a1619f0194be97a6c394848008f33823a7ee03` 

**PAYLOAD:**    
```json
{
    "role": "taker",
    "sending_amount": 10000000000000000000,
    "sending_token": "0x7B874444681F7AEF18D48f330a0Ba093d3d0fDD2",
    "receiving_amount": 100000000000000000000,
    "receiving_token": "0x9E7c6C6bf3A60751df8AAee9DEB406f037279C2a"
}
``` 
**Example Request :**  

**the maker:**   
`PUT /api/1/token_swaps/0x69C5621db8093ee9a26cc2e253f929316E6E5b92/0x8e90b850fdc5475efb04600615a1619f0194be97a6c394848008f33823a7ee03`  

**PAYLOAD:**     
```json
{
    "role": "maker",
    "sending_amount": 100000000000000000000,
    "sending_token": "0x9E7c6C6bf3A60751df8AAee9DEB406f037279C2a",
    "receiving_amount": 10000000000000000000,
    "receiving_token": "0x7B874444681F7AEF18D48f330a0Ba093d3d0fDD2",
    "secret": "0x40a6994181d0b98efcf80431ff38f9bae6fefda303f483e7cf5b7de7e341502a"
}
```
**Status Codes :**  
- `201 Created` - success 
- `400 Bad Request` - "no route available"(the scene may be happen when the old secret was used or the intermediate node has no corresponding token)


## switch to no network
 `GET /api/1/switch/*(Boolean)*` 

This interface is provided to switch to no-network state,by the way,this is mainly provided for APP, if the interface is called, the indirect transaction is prohibited.

**Example Request :**  

`GET ：http://{{ip2}}/api/1/switch/true`

**Parameters:**
- Boolean  
  - `true` - Switch to no network
  - `false` -With the network connection

   When switching to no network state, only `direct transactions` can be accepted.

  Note: Although the node can use this interface to switch to the no-network state, if the node sends a direct transfer transaction to other nodes with direct channels, it can still succeed. Therefore, the interface is not completely "no network", but only achieves the shielding function of indirect transactions.

##  Update node registered information
 ` POST /api/1/updatenodes` 

 It is necessary to update node information in order to ensure normal transaction at the state of no-network.
when the node information were registered, the transfer will take the UDP mode of communication for inproving TPS.Both nodes need to call the interface to update the other party's information,and if the registered node is restarted,the information need to re-registered since the message is stored in memory.

**Example Request :**  

`POST http://{{ip2}}/api/1/updatenodes`

**Request parameters(Suppose the other party's port is:192.168.14.13：60002):**

```json
[{
   "address":"0x151E62a787d0d8d9EfFac182Eae06C559d1B68C2",
   "ip_port":"192.168.14.13:60002"
}]
```
**Example Response :**   
**200 OK**  


## Set the fee policy
 ` POST /api/1/fee_policy `

The interface is called by the user to provide the PFS with the charging rate of the local node. When sending a transaction, the node requests routing from the PFS, and the PFS calculates the shortest path of the total cost based on the rate conditions submitted by all nodes. The node can use the PFS recommended fee plan and give enough money to transfer the transaction. If the amount is not enough, the transaction fails. If the node does not set a rate, PFS will use the global default value for route charging calculations, and there may be cases where the transaction fails. In addition, even if the node rate is set, if the sender clarifies the cost of the transfer in the transfer fee field (requires a fee plan greater than or equal to the PFS recommendation, or there is no available route), the actual charge is executed at the user-specified fee (last hop in the route collects extra fees).

**Example Request :**  

`POST http://{{ip1}}/api/1/fee_policy`

**PAYLOAD:**   

```json 
{
    "account_fee":{
        "fee_constant":5,
        "fee_percent":10000
    },
    "token_fee_map":{
        "0x83073FCD20b9D31C6c6B3aAE1dEE0a539458d0c5":{
            "fee_constant":5,
            "fee_percent":10000
        }
    },
    "channel_fee_map":{
        "0xa7712241a1a10abdada1c228c6935a71a9db80aa0bf2a13b59940159aa4eb4b5":{
            "fee_constant":5,
            "fee_percent":10000
        }
    }
}
```
- fee_constant: Fixed charge 
- fee_percent: fee rate
  
  Where `fee_constant` is the fixed rate, for example, 5 means that the fixed fee is 5 tokens, and setting it to 0 means no charge. `fee_percent` is the proportional rate, calculated as the transaction amount/`fee_percent`, such as transaction amount 50000000000000000000000, `fee_percent`=10000, then the commission ratio part = 50000000000000000000000/10000=5000000000000000000, set to 0 means no charge.
 Charge rule fee = `fee_constant` + amount/`fee_percent`

 There are three charging modes for a node:
- account_fee    Node charging
- token_fee      Node charging  on Specific token
- channel_fee    Node charging at a certain channel
 The priority of the three charging modes is：`channel_fee`>`token_fee`>`account_fee`

## Query the fee policy
`  GET /api/1/fee_policy `

Query the node charging information, which connect to default PFS server. If the rate has been set, return the fee rate information, otherwise, return the default information.

**Example Request :**  

`GET：http://{{ip1}}//api/1/fee_policy`

**Example Response :**  

**200 OK**   

```json  
{
    "Key": "feePolicy",
    "account_fee": {
        "fee_constant": 0,
        "fee_percent": 10000,
        "signature": null
    },
    "token_fee_map": {},
    "channel_fee_map": {}
}
```
## Query node charge record
 ` GET  /api/1/fee`

 The user can query all feecharge record of the intermediate node in different channels through the interface, so as to verify the revenue situation.

 **Example Request :**  

  `GET： http://{{ip2}}/api/1/fee`

**Example Response :**  
200 OK

```json 
{
    "error_code": "0000",
    "error_msg": "SUCCESS",
    "data": {
        "total_fee": {
            "0x2a7af974b7bb88703180d6aff9a656bb4dbba809": 15
        },
        "details": [
            {
                "key": "0x4e50d0211bc09079583a0d902f6e8e5bc6fa89b4b2d8e8f0ee52316f7f5439eb",
                "lock_secret_hash": "0x6d0a349cb75de3020b90b3b8d05be127e56dfc15d16686488ecfb99ff79e91b5",
                "token_address": "0x2a7af974b7bb88703180d6aff9a656bb4dbba809",
                "transfer_from": "0x97cd7291f93f9582ddb8e9885bf7e77e3f34be40",
                "transfer_to": "0xc445a8c326a8fd5a3e250c7dc0efc566edcb263b",
                "transfer_amount": 100000000000000000000,
                "in_channel": "0xf25edb59e35544e060ecfcef6e6a0ba619ff905a132295957e11ffdc2206fc24",
                "out_channel": "0x081f7a9771994de9f06edb52cb60a0fe3b9bbebd4c1240c267967c7e3fa433f5",
                "fee": 5,
                "timestamp": 1548151954
            },
            {
                "key": "0x6ade0365b8a2c4cdfbcd5bbc40cb46665bdb4e5453a644a6dd49ba7717a6f8f8",
                "lock_secret_hash": "0x7ad052f2dd5c30a5cadc0fc3d34eb4f728e338c714fcc305402760c12099efa1",
                "token_address": "0x2a7af974b7bb88703180d6aff9a656bb4dbba809",
                "transfer_from": "0x97cd7291f93f9582ddb8e9885bf7e77e3f34be40",
                "transfer_to": "0xc445a8c326a8fd5a3e250c7dc0efc566edcb263b",
                "transfer_amount": 100000000000000000000,
                "in_channel": "0xf25edb59e35544e060ecfcef6e6a0ba619ff905a132295957e11ffdc2206fc24",
                "out_channel": "0x081f7a9771994de9f06edb52cb60a0fe3b9bbebd4c1240c267967c7e3fa433f5",
                "fee": 5,
                "timestamp": 1548151970
            },
            {
                "key": "0x81cf65031f17ede87fdc9022943077f0105184b01f07c85660906b80be15fe00",
                "lock_secret_hash": "0xc5b85dcd61f3874caa5363998b452365e0df37e3740cd1728a838cd2e4cd94d2",
                "token_address": "0x2a7af974b7bb88703180d6aff9a656bb4dbba809",
                "transfer_from": "0x97cd7291f93f9582ddb8e9885bf7e77e3f34be40",
                "transfer_to": "0xc445a8c326a8fd5a3e250c7dc0efc566edcb263b",
                "transfer_amount": 100000000000000000000,
                "in_channel": "0xf25edb59e35544e060ecfcef6e6a0ba619ff905a132295957e11ffdc2206fc24",
                "out_channel": "0x081f7a9771994de9f06edb52cb60a0fe3b9bbebd4c1240c267967c7e3fa433f5",
                "fee": 5,
                "timestamp": 1548151963
            }
        ]
    }
}

```
## Query the charging route from the node to the target node
` GET /api/1/path/{target_address}/{token_address}/"amount"`

The user invokes the interface to query whether the target node has available routes and fees. If there are multiple routes with the same cost, they are given together.

**Example Request :**  

`GET：http://{{ip1}}/api/1/path/0xEfB2e46724f675381ce0b3F70Ea66383061924E9/0x5b9d594750bb54f95E372F17a04a70E488284f64/100`
  
**Example Response :**  

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





