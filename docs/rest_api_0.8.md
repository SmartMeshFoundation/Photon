# SmartRaiden REST API Reference
Hey guys, welcome to SmartRaiden REST API Reference page. This is an API Spec for SmartRaide version 1.0, which adds a lot more new features, as CooperateWithdraw, CooperateCloseChannel, send specific `secret`, etc. Please note that this reference is still updating. If any problem, feel free to submit at our [Issue](https://github.com/SmartMeshFoundation/SmartRaiden/issues). 
  
## Channel Structure
```json
    {
        "channel_address": "0x47235d9d81eb6c19dea2b695b3d6ba1cf76c169d329dc60d188390ba5549d025",
        "open_block_number": 3158573,
        "partner_address": "0x31DdaC67e610c22d19E887fB1937BEE3079B56Cd",
        "balance": 100,
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

channel variables explanation :  
* `channel_address` : address for a channel  
* `open_block_number` : block height when a channel opens  
* `partner_address` : address for your channel partner  
* `balance` : your token balance in this channel  
* `partner_balance` : token balance for your channel partner  
* `locked_amount` : the amount of token you locked in this channel  
* `partner_locked_amount` : the amount of token your partner locked in this channel  
* `token_address` : address for tokens in this channel  
* `state` : digits denoting transaction states  
* `StateString` : String literal for Channel States  
* `settle_timeout` : some amount of block denoting time period for transaction settlement  
* `reveal_timeout` : block height at which nodes registering `secret`      


State|StateString|Description
---|---|---
0|Invalid|Invalid Channel
1|Opened|Channel opened with normal transfer ongoing
2|Closed|Stop sending transfer but able to receive
3|BalanceProofUpdated|stop unfinished transfer and not receive `unlock`
4|Settled|Channel Settlement completes
5|Closing|StateClosing users request for channel closing, ongoing transfers continue but no more newly-opened transfer.
6|Settling|StateSettling users start a settle request. Transfers cannot be processed, ongoing transfers stops, Deny any newly-opened transfer.
7|Withdraw|StateWithdraw users send/receive withdraw request, and ongoing transfers stop immediately.
8|CooperativeSettle|StateCooperativeSettle users send/receive cooperative settle request, stop any ongoing transfer.
9|PrepareForCooperativeSettle|StatePrepareForCooperativeSettle cooperative request received with ongoing transfers, but no more newly-opened transfer. Channels need to wait for Channel Settle.
10|PrepareForWithdraw|StatePrepareForWithdraw has received user request, and prepares to process withdraw, but there are tokens locked in, and channel participants cannot open/receive any transfer. Can wait for certain block number to process withdraw.
11|Error|StateError 

## GET /api/1/address
Check Node's data, which returns the address of SmartRaiden node.    
**Example Response:**       
```json
{
    "our_address": "0x69C5621db8093ee9a26cc2e253f929316E6E5b92"
}
```
**Status Codes:**  
- `200 OK` - Check Success  
- `404 Not Found` - Check Failure      

## GET /api/1/tokens  
Check registered token    
**Example Response:**    
```json
[
    "0x7B874444681F7AEF18D48f330a0Ba093d3d0fDD2"
]
```
**Status Codes:**  
- `200 OK` - Check Success  
- `404 Not Found` - Check Failure    
## GET /api/1/tokens/*(token_address)*/partners  
Get all the channel partners of this token.    
**Example Request:**        
`GET /api/1/tokens/0xD82E6be96a1457d33B35CdED7e9326E1A40c565D/partners`     
 **Example Response :**      
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
## PUT /api/1/tokens/*(token_address)*
Register another token type     

**Example Request:**      
`PUT /api/1/tokens/0x9E7c6C6bf3A60751df8AAee9DEB406f037279C2a`    

**Example Response:**    
```json
{
    "channel_manager_address": "0xBb1e95363b0181De7bBf394f18eaC7D4230e391A"
}
```
**Status Codes:**  
- `200 OK` - Register Success  
- `400 Bad Request` - Invalid Token Address  
- `409 Conflict` - Token has been registered     


## GET /api/1/channels  
Check all unsettled channels of a node.    
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
**Status Codes:**  
- `200 OK` - Check Success  
- `404 Not Found` - Check Failure    

## POST /api/1/channels
Open a new Channel    
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
**Status Codes:**  
- `200 OK` - Open Channel Success  
- `400 Bad Request` - Invalid Parameter  
- `409 Conflict` - Channel Already Opened    

## GET /api/1/channels/*(channel_address)* 
Check specific channel, can get in-depth channel information    
**Example Request**    
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
**Status Codes:**  
- `200 OK` - Check Success  
- `404 Not Found` - Check Failure    
## PUT /api/1/withdraw/*(channel_address)*  
CooperateWithdraw available when both channel participants online  
**PAYLOAD:**    
```json
{
	"amount":0,
	"op":"preparewithdraw"
}
```
**Request JSON Object:**  
- `op` - Alter Channel States(Optional)  
  - `preparewithdraw` - Alter Channel State to `prepareForWithdraw`, detail in Channel State Chart  
  - `cancelprepare` - cancel prepare/alter channel state to `open`    
 
**Example Response:**     
```json
{
    "channel_address": "0x47235d9d81eb6c19dea2b695b3d6ba1cf76c169d329dc60d188390ba5549d025",
    "open_block_number": 3613578,
    "partner_address": "0x31DdaC67e610c22d19E887fB1937BEE3079B56Cd",
    "balance": 190,
    "partner_balance": 100,
    "locked_amount": 0,
    "partner_locked_amount": 0,
    "token_address": "0xF2747ea1AEE15D23F3a49E37A146d3967e2Ea4E5",
    "state": 7,
    "StateString": "withdrawing",
    "settle_timeout": 150,
    "reveal_timeout": 5
}
```
**Status Codes:**  
- `200 OK ` - Withdraw Success  
- `400 Bad Request` - Invalid Parameter/Low Token Balance    

## PATCH /api/1/channels/*(channel_address)*
Deposit in a channel    
**Example  Request:**      
`PATCH /api/1/channels/0x97f73562938f6d538a07780b29847330e97d40bb8d0f23845a798912e76970e1`    
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
**Status Codes:**  
- `200 OK` - Deposit Success  
- `400 Bad Request` - Invalid Requst Parameter      

Close a channel, set `force` default to `false`, meaning that channel participants cooperate settle channel.    

**Example  Request:**       
`PATCH /api/1/channels/0x97f73562938f6d538a07780b29847330e97d40bb8d0f23845a798912e76970e1`           
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
Once channel partner is offline or do not wish to cooperate settle, then alter `force` to `true`, wait for `settle_timeout` then do channel settle procedure.    
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
Settle Channel. Once channel is closed and `settle_timeout` block has passed, channels can be settled.    
**Example  Request:**        
`PATCH /api/1/channels/0x97f73562938f6d538a07780b29847330e97d40bb8d0f23845a798912e76970e1`       
   
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
**Status Codes :**  
- `200 OK` - Close/Settle Success  
- `400 Bad Request` - Invalid Parameter  
- `409 Conflict` - State Conflicts    


## POST /api/1/transfer/*(token_address)*/*(target_address)*  
When channel state is `open` with sufficient funds, participants can make transfers in it.     
**Example Request :**      
`POST /api/1/transfers/0x7B874444681F7AEF18D48f330a0Ba093d3d0fDD2/0xf2234A51c827196ea779a440df610F9091ffd570`      
**PAYLOAD :**    
```json
{
    "amount":20,  
    "fee":0, // fee for transfer routing   
    "is_direct":false // whether it is a direct transfer  
   
}
```
**Example Response :**   
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
Send transfers with specified `secret`.    

**Example Request :**      
`http://{{ip1}}/api/1/transfers/0xF2747ea1AEE15D23F3a49E37A146d3967e2Ea4E5/0xf0f6E53d6bbB9Debf35Da6531eC9f1141cd549d5`    
**PAYLOAD :**    
```json
{
    "amount":20,
    "fee":0,
    "is_direct":false,
    "secret":"0xad96e0d02aa2f4db096e3acdba0831f95bb09d876a5c6f44bc3f7325a0a45ea1"
}
```
## GET /api/1/querysenttransfer
Query the transaction record that is sent successfully and return all successful transactions list.    
**Example Response :**     
```json
[
    {
        "Key": "0xd971f803c7ea39ee050bf00ec9919269cf63ee5d0e968d5fe33a1a0f0004f73d-3",
        "block_number": 4490372,
        "OpenBlockNumber": 0,
        "channel_address": "0xd971f803c7ea39ee050bf00ec9919269cf63ee5d0e968d5fe33a1a0f0004f73d",
        "to_address": "0x151e62a787d0d8d9effac182eae06c559d1b68c2",
        "token_address": "0xd82e6be96a1457d33b35cded7e9326e1a40c565d",
        "nonce": 3,
        "amount": 10
    },
    {
        "Key": "0xd971f803c7ea39ee050bf00ec9919269cf63ee5d0e968d5fe33a1a0f0004f73d-5",
        "block_number": 4490580,
        "OpenBlockNumber": 0,
        "channel_address": "0xd971f803c7ea39ee050bf00ec9919269cf63ee5d0e968d5fe33a1a0f0004f73d",
        "to_address": "0x151e62a787d0d8d9effac182eae06c559d1b68c2",
        "token_address": "0xd82e6be96a1457d33b35cded7e9326e1a40c565d",
        "nonce": 5,
        "amount": 10
    }
]
```
## GET /api/1/queryreceivedtransfer
Query successfully received transaction record of Unlock message.     
**Example Response :**   
```json
[
    {
        "Key": "0x79b789e88c3d2173af4048498f8c1ce66f019f33a6b8b06bedef51dde72bbbc1-2",
        "block_number": 4492809,
        "OpenBlockNumber": 0,
        "channel_address": "0x79b789e88c3d2173af4048498f8c1ce66f019f33a6b8b06bedef51dde72bbbc1",
        "token_address": "0xd82e6be96a1457d33b35cded7e9326e1a40c565d",
        "from_address": "0x201b20123b3c489b47fde27ce5b451a0fa55fd60",
        "nonce": 2,
        "amount": 10
    },
    {
        "Key": "0x79b789e88c3d2173af4048498f8c1ce66f019f33a6b8b06bedef51dde72bbbc1-6",
        "block_number": 4493353,
        "OpenBlockNumber": 0,
        "channel_address": "0x79b789e88c3d2173af4048498f8c1ce66f019f33a6b8b06bedef51dde72bbbc1",
        "token_address": "0xd82e6be96a1457d33b35cded7e9326e1a40c565d",
        "from_address": "0x201b20123b3c489b47fde27ce5b451a0fa55fd60",
        "nonce": 6,
        "amount": 20
    }
]
```
## GET /api/1/getunfinishedreceivedtransfer/*(tokenaddress)*/*(locksecrethash)*
The receiver of the transaction inquires the transaction that has not yet been received.    
**Example Request :**    
`GET /api/1/getunfinishedreceivedtransfer/0xD82E6be96a1457d33B35CdED7e9326E1A40c565D/0x2fb55cec26a26d0212cf6bd6022aaa7426410916de09133be3b353ac1a91d843`    
**Example Response :**     
```json
{
    "initiator_address": "0x201B20123b3C489b47Fde27ce5b451a0fA55FD60",
    "target_address": "0x151E62a787d0d8d9EfFac182Eae06C559d1B68C2",
    "token_address": "0xD82E6be96a1457d33B35CdED7e9326E1A40c565D",
    "amount": 30,
    "secret": "",
    "lock_secret_hash": "0x2fb55cec26a26d0212cf6bd6022aaa7426410916de09133be3b353ac1a91d843",
    "expiration": 131,
    "fee": null,
    "is_direct": false
}
```

## GET /api/1/transferstatus/*(token_address)*/*(locksecrethash)*
Query transaction status       
**Example Request :**    
`GET /api/1/transferstatus/0xD82E6be96a1457d33B35CdED7e9326E1A40c565D/0xdb0d663a82d04fedf4f558f75d7be801ab6707ea765662919063bad93cd71c82`  
**Example Response :**     
```json
{
    "LockSecretHash": "0xdb0d663a82d04fedf4f558f75d7be801ab6707ea765662919063bad93cd71c82",
    "Status": 0,
    "StatusMessage": "MediatedTransfer 正在发送 target=151e\nMediatedTransfer 发送成功\n"
}
```
**Response JSON Array of Objects :**  
- `Status`
  - 0 - Revocable transaction status(Before sending the secret)  
  - 1 - Irrevocable transaction status(After sending the secret)  
  - 2 - Transaction Successed  
  -  3 - Transaction revocation    

## POST /api/1/registersecret  
Register `secret`, after which `MediatedTransfer` can be successfully unlocked.    
**PAYLOAD :**  
```json
{
	"secret":"0xad96e0d02aa2f4db096e3acdba0831f95bb09d876a5c6f44bc3f7325a0a45ea1",
	"token_address":"0xF2747ea1AEE15D23F3a49E37A146d3967e2Ea4E5"
}
```
**Status Codes :**  
- `200 OK` - Transfer Success  
- `400 Bad Request` - Invalid Parameter  
- `409 Conflict` - No Valid Router    

## PUT /api/1/token_swaps/*(target_address)*/*(lock_secret_hash)*      
Token Swap can be used to exchange within two types of tokens. Under the circumstances that valid routing strategies are existed, first invoke `taker` then `maker`, and with `/api/1/secret/` channel participants can receive a `lock_secret_hash` / `secret` pair.    
**Example Request :**      
`PUT /api/1/token_swaps/0x31DdaC67e610c22d19E887fB1937BEE3079B56Cd/0x8e90b850fdc5475efb04600615a1619f0194be97a6c394848008f33823a7ee03`        
**PAYLOAD :**  
```json
{
    "role": "taker",
    "sending_amount": 10,
    "sending_token": "0x7B874444681F7AEF18D48f330a0Ba093d3d0fDD2",
    "receiving_amount": 100,
    "receiving_token": "0x9E7c6C6bf3A60751df8AAee9DEB406f037279C2a"
}
```

`PUT /api/1/token_swaps/0x69C5621db8093ee9a26cc2e253f929316E6E5b92/0x8e90b850fdc5475efb04600615a1619f0194be97a6c394848008f33823a7ee03`      

**PAYLOAD :**   
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
**Status Codes :**  
- `201 Created` - Success  
- `400 Bad Request` - Invalid Parameter       
## GET /api/1/secret
Receive `lock_secret_hash` / `secret` pair.    
**Example Response :**   
```json
{
    "lock_secret_hash": "0x8e90b850fdc5475efb04600615a1619f0194be97a6c394848008f33823a7ee03",
    "secret": "0x40a6994181d0b98efcf80431ff38f9bae6fefda303f483e7cf5b7de7e341502a"
}
```

## Post /api/1/transfercancel/*(token)*/*(locksecrethash)*
To revoke a transaction according to token and locksecrethash, only the initiator can invoke it, and the transaction must be revocable.    
**Example Request :**   
`POST /api/1/transfercancel/0xD82E6be96a1457d33B35CdED7e9326E1A40c565D/0xe0f8d65ddb4f70899b97f36795925a97c1b286582f58f56a041f141d345acdca`     
**Example Response :**    
**200 OK**    
The transaction status can be querying through the`/api/1/transferstatus`  
## GET /api/1/switch/*(Boolean)*
Switch to no net state  
- Boolean  
  - `true` - Switch to nonetwork  
  - `false` - Switch to network  
When switching to no net state, only direct transactions can be accepted.    

##  POST /api/1/updatenodes
Update node information,It is necessary to update node information in order to ensure normal transaction without network conditions.    
**PAYLOAD :**     
```json
[{
   "address":"0x151E62a787d0d8d9EfFac182Eae06C559d1B68C2",
   "ip_port":"127.0.0.1:60002"
},
{
   "address":"0x10b256b3C83904D524210958FA4E7F9cAFFB76c6",
   "ip_port":"127.0.0.1:60001",
   "device_type":"mobile"
}]
```
**Example Response :**    
**200 OK**  


