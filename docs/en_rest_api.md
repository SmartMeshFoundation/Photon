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
- `channel_address` : address for a channel   
- `open_block_number` : block height when a channel opens   
- `partner_address` : address for your channel partner   
- `balance` : your token balance in this channel  
- `partner_balance` : token balance for your channel partner  
- `locked_amount` : the amount of token you locked in this channel  
- `partner_locked_amount` : the amount of token your partner locked in this channel  
- `token_address` : address for tokens in this channel  
- `state` : digits denoting transaction states   
- `StateString` : String literal for Channel States  
- `settle_timeout` : some amount of block denoting time period for transaction settlement  
- `reveal_timeout` : block height at which nodes registering `secret`  


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
**Status Codes:**    
- `200 OK` - Close/Settle Success    
- `400 Bad Request` - Invalid Parameter   
- `409 Conflict` - State Conflicts    


## POST /api/1/transfer/*(token_address)*/*(target_address)*  
When channel state is `open` with sufficient funds, participants can make transfers in it.   
**Example Request:**    
`POST /api/1/transfers/0x7B874444681F7AEF18D48f330a0Ba093d3d0fDD2/0xf2234A51c827196ea779a440df610F9091ffd570`    
**PAYLOAD**  
```json
{
    "amount":20,
    "fee":0, // fee for transfer routing 
    "is_direct":false // whether it is a direct transfer
   
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
Send transfers with specified `secret`.  

**Example Request**    
`http://{{ip1}}/api/1/transfers/0xF2747ea1AEE15D23F3a49E37A146d3967e2Ea4E5/0xf0f6E53d6bbB9Debf35Da6531eC9f1141cd549d5`    
**PAYLOAD:**  
```json
{
    "amount":20,
    "fee":0,
    "is_direct":false,
    "secret":"0xad96e0d02aa2f4db096e3acdba0831f95bb09d876a5c6f44bc3f7325a0a45ea1"
}
```
## GET /api/1/getunfinishedreceivedtransfer/*(token_address)*/*(locksecrethash)*  
Check unfinished transfers   
**Example Request:**    
`/api/1/getunfinishedreceivedtransfer/0xF2747ea1AEE15D23F3a49E37A146d3967e2Ea4E5/0x992a8b9751180ef5363184bd4af54b7d5bc66f99e4239250c6ef23840ee5464c`  

**Example Response:**
```json
{
    "initiator_address": "0x69C5621db8093ee9a26cc2e253f929316E6E5b92",
    "target_address": "0xf0f6E53d6bbB9Debf35Da6531eC9f1141cd549d5",
    "token_address": "0xF2747ea1AEE15D23F3a49E37A146d3967e2Ea4E5",
    "amount": 20,
    "secret": "",
    "lock_secret_hash": "0x992a8b9751180ef5363184bd4af54b7d5bc66f99e4239250c6ef23840ee5464c",
    "expiration": 132,
    "fee": null,
    "is_direct": false
}
```
## POST /api/1/registersecret  
Register `secret`, after which `MediatedTransfer` can be successfully unlocked.  
**PAYLOAD:**
```json
{
	"secret":"0xad96e0d02aa2f4db096e3acdba0831f95bb09d876a5c6f44bc3f7325a0a45ea1",
	"token_address":"0xF2747ea1AEE15D23F3a49E37A146d3967e2Ea4E5"
}
```
**Status Codes:**    
- `200 OK` - Transfer Success  
- `400 Bad Request` - Invalid Parameter  
- `409 Conflict` - No Valid Router  

## PUT /api/1/token_swaps/*(target_address)*/*(lock_secret_hash)*      
Token Swap can be used to exchange within two types of tokens. Under the circumstances that valid routing strategies are existed, first invoke `taker` then `maker`, and with `/api/1/secret/` channel participants can receive a `lock_secret_hash` / `secret` pair.  
**Example Request:**    
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
**Status Codes:**      
- `201 Created` - Success   
- `400 Bad Request` - Invalid Parameter     
## GET /api/1/secret
Receive `lock_secret_hash` / `secret` pair.  
**Example Response:** 
```json
{
    "lock_secret_hash": "0x8e90b850fdc5475efb04600615a1619f0194be97a6c394848008f33823a7ee03",
    "secret": "0x40a6994181d0b98efcf80431ff38f9bae6fefda303f483e7cf5b7de7e341502a"
}
```
