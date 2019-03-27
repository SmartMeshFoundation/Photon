# Photon REST API Reference  

Hey guys, welcome to Photon REST API Reference page. This is an API Spec for Photon version 1.1, which adds a lot more new features, such as, support multi-token functions, support SMT mortgage,use mDNS to solve node discovery, use PFS to support channel charging,etc. Please note that this reference is still updating. If any problem, feel free to submit at our Issue.

##  Channel Structure  

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

- `error_code`:  Error code

- `error_message`: Error Code description

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

-  `closed_block`: The block height at channel closure

-  `settled_block`: The block height at channel settlement

-  `our_balance_proof`: The balance proof data of the participant

-  `partner_balance_proof`: The balance proof data of the partner

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
11|StatePartnerCooperativeSettling|After the user receives and agrees the CooperativeSettle request from the other party, the channel is set to this state.
12|StatePartnerWithdrawing|After the user receives and agrees the withdraw request from the other party, the channel is set to this state.

Among them, App can see only 1-9 states, other states can not be directly observed, which is internal use.  **prepareForWithdraw and prepareForCooperativeSettle will not appear on the mobile phone** , only appear when the meshbox  be used as intermediate node of the transaction.

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

##  Query node address

 `GET /api/1/address`

 Return the address of photon node.

**Example Request :** 

`GET  http://{{ip1}}/api/1/address`

**Example Response :** 

```json
{
    "error_code": 0,
    "error_message": "SUCCESS",
    "data": "0x97Cd7291f93F9582Ddb8E9885bF7E77e3f34Be40"
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
**Status Codes:**  
- `200 OK` 


## Get all the channel partners of this token
 `GET /api/1/tokens/*(token_address)*/partners`

   Return all channels of this node under this token.

**Example Request :**  

`GET http://{{ip2}}/api/1/tokens/0xB31567308AD3c42D864FB41684bB40d3A2c57E1b/partners`

**Example Response :** 

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
**Status Codes:**  
- `200 OK` 

 If the node has not created the channel with the token, then respond the message:
 
```json
{
    "error_code": 0,
    "error_message": "SUCCESS",
    "data": null
}
```

 
## Query all the channels of the node
   `GET /api/1/channels`  

Return all the unsettled channels of the node.

**Example Request :**  

 `GET  http://{{ip1}}/api/1/channels`

**Example Response :**  

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
**Status Codes:**  
- `200 OK` 

## Query specific channel of the node
  `GET /api/1/channels/*(channel_identifier)* `

Query the specific channel and return all the information about the channel.

**Example Request :**  

`GET http://{{ip1}}/api/1/channels/0xfe738aa39610416e4100036130af7ae00930021d5a51be60b55b96c12b1f4af5`

**Example Response :**  

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

 1.Create channel and deposit
 
   - `new_channel`sets`true`, which means open a new channel and deposit to the channel;if there is no channel between the participants, `false`is no meaning for `new_channel`，which will return the error message "There is no channel". 
   - `settle_timeout`represent the settlement window for new channel, for example,settle_timeout：100; if the `settle_timeout` set to 0,the default window period is used which is 600 block.

 2.Only deposit:
 
   - `new_channel`must set to `false`，which means the channel has been existed；If the channel has been existed ,there is Meaningless to set the `new_channel`statue as `ture`，which will response the error message "The channel has already existed". 
   - `settle_timeout`must set to 0,because the channel has already existed.

 **Example Request :**  

  `PUT http://{{ip1}}/api/1/deposit`

**PAYLOAD:** 
```json 
{
    "partner_address": "0x3bC7726c489E617571792aC0Cd8b70dF8A5D0e22",
    "token_address": "0xB31567308AD3c42D864FB41684bB40d3A2c57E1b",
    "balance": 10000000000000000000000,
    "settle_timeout":0,
     "new_channel":false
   
}
```

**Example Response :**  

```json
{
    "error_code": 0,
    "error_message": "SUCCESS",
    "data": {
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

`PUT http://{{ip2}}/api/1/withdraw/0xfe738aa39610416e4100036130af7ae00930021d5a51be60b55b96c12b1f4af5`

**PAYLOAD:**  
```json
{
		"op":"preparewithdraw"
}
```
**Example Response :** 

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

**PAYLOAD:**   
```json
{
		"op":"cancelprepare"
}
```
**Example Response :** 

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

**PAYLOAD:**     
```json
{
	"amount":50000000000000000000,

}
```
**Example Response :**  
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
If the withdrawn amount is larger than the available balance of the channel, an error message will be returned.
```json
{
    "error_code": 1,
    "error_message": "ArgumentError:errorCode: 1, errorMsg ArgumentError:invalid withdraw amount, availabe=19900000000000000000000,want=1000000000000000000000000"
}
```


##  Close the channel
`PATCH /api/1/channels/*(channel_identifier)* `

Close the channel, which includes the unilateral close the channel and cooperative settle the channel.
set `force` default to `false`, meaning that channel participants cooperate settle the channel.When no new block is received from the connection point for more than one minute, an error message will be given when calling to close the channel."call smc SyncProgress err, client is closed”,which means that the connection point need to synchronize new blocks.

**Example Request :**    

`PATCH http://{{ip2}}/api/1/channels/0xe4c61eac5f3f45ea62c7f021cc0aa6a774feb14fed3eaa28af16b512f7fec966` 
   
**PAYLOAD:**      
```json
{"state":"closed"，
  "force":false
}
```
**Example Response :**    
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
Once channel partner is offline or has the locks, the cooperate settle can't be carried out.
```json 
{
    "error_code": 1,
    "error_message": "ArgumentError:errorCode: 1016, errorMsg NodeOffline:node 0xC445a8C326A8fD5a3e250C7dc0EFc566eDcB263B is not online"
}
```

The participant should alter the`force` to `true`, wait for settle_timeout and unilateral settle the channel.

**PAYLOAD:**     
```json 
{"state":"closed",
  "force":true
}
```
**Example Response :**    
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
##  Settle the Channel
`PATCH /api/1/channels/(channel_identifier)`

The interface of unilaterally settling channel is reused with closing channel, which the parameters are different.
 
 After unilaterally closing the channel, it is necessary to call the settlement channel to settle the closed channel.Once the half of the settle_timeout block has passed,the PMS can submit the balanceproof of the delegate and unlock the registered transaction,and the channel participants can submit the balanceproof to undate the channel and unlock the registered transaction at any time during the settlement window period.
 Tips:When no new block is received from the connection point for more than one minute, an error message will be given when calling to settle the channel."call smc SyncProgress err, client is closed”,which means that the connection point need to synchronize new blocks.

Note: Since settle_timeout does not include the penalty period (in spectrum, which is 257  block, about an hour), the actual settlement time is about 410 block.

**Example Request :**  

`PATCH http://{{ip2}}/api/1/channels/0xe4c61eac5f3f45ea62c7f021cc0aa6a774feb14fed3eaa28af16b512f7fec966`   

**PAYLOAD:**   
```json
{
    "state":"settled"
}
```

**Example Response :**  
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
- `200 OK` - close/settle success
- `409 Conflict` - State conflicts, such as, "failed to estimate gas needed: gas required exceeds allowance or always failing transaction",or "channel is still open".

## Initiate the payment
`POST /api/1/transfers/*(token_address)*/*(target_address)*`  

This interface is used to initiate a transfer transaction, which is currently associated with PFS by default.

**Example Request :**   
`POST http://{{ip1}}/api/1/transfers/0xB31567308AD3c42D864FB41684bB40d3A2c57E1b/0xd5dC7504e0b448b1c62D86306AE8e4a5836Fc1A1`

**PAYLOAD:**     
```json
{ "amount":10000000000,
    "is_direct":false,
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
]
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
    "error_code": 0,
    "error_message": "SUCCESS",
    "data": {
        "initiator_address": "0x97Cd7291f93F9582Ddb8E9885bF7E77e3f34Be40",
        "target_address": "0xd5dC7504e0b448b1c62D86306AE8e4a5836Fc1A1",
        "token_address": "0xB31567308AD3c42D864FB41684bB40d3A2c57E1b",
        "amount": 10000000000,
        "lockSecretHash": "0x14c97ba1f3a6850d5ddec5c486d673ada87cc3a9de7f4b1a6050b61e598a2ec9",
        "data": "",
        "route_info": [
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
        ]
    }
}
```
Note: The new version makes the designated routing transfer. If the local photon node does not update the rate to PFS in time, there may be inconsistency between the charge and the calculation of PFS, the actual charges shall prevail.


## Initiate the transfer with specified secret

The normal transfer secret is automatically generated by photon. If the user wants to precisely control the success or failure of the transaction, he can use the transfer of the specified `secret`. Currently a major application scenario is tokenswap.

**Example Request :**   

`POST http://{{ip1}}/api/1/transfers/0xB31567308AD3c42D864FB41684bB40d3A2c57E1b/0xd5dC7504e0b448b1c62D86306AE8e4a5836Fc1A1` 

```json
{
    "amount":10000000000,
    "is_direct":false,
    "secret":"0x9a01a92aebd7419a5645d05eb344896e25d9c919ef67efa0521996127adbc07d",
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
]
      }
```
Note: The specified secret is obtained by the interface  `/api/1/secret`.

**Example Response :**   
```json
{
    "error_code": 0,
    "error_message": "SUCCESS",
    "data": {
        "initiator_address": "0x97Cd7291f93F9582Ddb8E9885bF7E77e3f34Be40",
        "target_address": "0xd5dC7504e0b448b1c62D86306AE8e4a5836Fc1A1",
        "token_address": "0xB31567308AD3c42D864FB41684bB40d3A2c57E1b",
        "amount": 10000000000,
        "secret": "0x9a01a92aebd7419a5645d05eb344896e25d9c919ef67efa0521996127adbc07d",
        "lockSecretHash": "0x1dacb8dcdc1088dad043d07aff1b812760b8ea04525d1f76961a0d46765ec9e0",
        "data": "",
        "route_info": [
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
        ]
    }
}
```
 The specified secret transfer is locked.

```json
{
            "channel_identifier": "0xfe738aa39610416e4100036130af7ae00930021d5a51be60b55b96c12b1f4af5",
            "open_block_number": 1932436,
            "partner_address": "0x3bC7726c489E617571792aC0Cd8b70dF8A5D0e22",
            "balance": 999989976388879,
            "partner_balance": 10023611121,
            "locked_amount": 10023611121,
            "partner_locked_amount": 0,
            "token_address": "0xB31567308AD3c42D864FB41684bB40d3A2c57E1b",
            "state": 1,
            "state_string": "opened",
            "settle_timeout": 100,
            "reveal_timeout": 30
        }
```

After registering secret with the `allow disclosure secret`interface, the lock is unlocked.
```json
{
            "channel_identifier": "0xfe738aa39610416e4100036130af7ae00930021d5a51be60b55b96c12b1f4af5",
            "open_block_number": 1932436,
            "partner_address": "0x3bC7726c489E617571792aC0Cd8b70dF8A5D0e22",
            "balance": 999979952777758,
            "partner_balance": 20047222242,
            "locked_amount": 0,
            "partner_locked_amount": 0,
            "token_address": "0xB31567308AD3c42D864FB41684bB40d3A2c57E1b",
            "state": 1,
            "state_string": "opened",
            "settle_timeout": 100,
            "reveal_timeout": 30
        }
```

## Get the secret
` GET /api/1/secret `

Through calling the interface,the caller will Get `lock_secret_hash` / `secret` pair,which can be used in Specified secret transaction or tokenswap.

 **Example Request :**   

 `GET  http://{{ip1}}/api/1/secret`

**Example Response :** 
```json
{
    "error_code": 0,
    "error_message": "SUCCESS",
    "data": {
        "lock_secret_hash": "0x4e7a5c8043a9faa93d3b094146b2ea2a65ec466e8cb3dbf7986779f802edf024",
        "secret": "0xd01a3ee8f92664426245099d14435cf93d47feedc9bfee2908e648c7e47d60b7"
    }
}
```

## Allow disclosure of secret
 `Post /api/1/transfers/allowrevealsecret`

This interface is used in combination with the specified secret transfer interface. When a transfer with the special secret is sent, if the interface is not be called to unlock the Secret, the initiator will not accept the SecretRequest from the recipient. So when sending the transaction with the specified secret, the sender must actively call this interface ,then the transaction can be successfully continued.
  
  **Example Request :**   

`Post  http://{{ip2}}/api/1/transfers/allowrevealsecret`

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
```json
{
    "error_code": 0,
    "error_message": "SUCCESS",
    "data": null
}

```
Query the channel,the locked amount has been unlocked.

```json
 {
            "channel_identifier": "0xa628d9ee19415c574bc6861a2cf17c0269cb37436aa90d9f2da59d72217a14da",
            "open_block_number": 1694480,
            "partner_address": "0xC445a8C326A8fD5a3e250C7dc0EFc566eDcB263B",
            "balance": 1.000000000002e+22,
            "partner_balance": 9.99999999998e+21,
            "locked_amount": 0,
            "partner_locked_amount": 0,
            "token_address": "0xB31567308AD3c42D864FB41684bB40d3A2c57E1b",
            "state": 1,
            "state_string": "opened",
            "settle_timeout": 100,
            "reveal_timeout": 30
        }

```

## Query the sent successful transfer 
  `GET /api/1/querysenttransfer` 

For the sender of the transfer, the interface can be used to query the history information of all successful transfer which sent from itself, so that the user can accurately master the situation of the transfered funds. If there is too much history, you can use block filtering.Such as:".../querysenttransfer?from_block=3000&to_block=5000"

**Example Request :**   

`GET http://{{ip2}}/api/1/querysenttransfer`

**Example Response :** 
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
## Query the received successful transfer 
   `GET /api/1/queryreceivedtransfer`
   
For the receiver of the transfer, the interface can be used to query the history information of all successful transfer which received from other partners, so that the user can accurately master the situation of the received funds. If there is too much history, you can use block filtering.Such as:".../queryreceivedtransfer?from_block=3000&to_block=5000"

**Example Request：**

`GET http://{{ip1}}/api/1/queryreceivedtransfer`

**Example Response :** 
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
##  Query the transaction that have not yet been received
   ` GET /api/1/getunfinishedreceivedtransfer/*(tokenaddress)*/*(locksecrethash)* `  

 This interface is called by the receiver, also used to specify the secret transaction scenario. The  receiver can find out that  a transaction has been received through the interface, but there is no secret. The receiver can request the sender to call allowrevealsecret to complete the transaction, otherwise the transaction will be returned after expiration.
 
 **Example Request :**    
`GET {{ip1}}/api/1/getunfinishedreceivedtransfer/0xB31567308AD3c42D864FB41684bB40d3A2c57E1b/0xd8875761c93aa9b804c42855601326cf722ced2be5d84fdee36c52ced95ba587`   

 **Example Response :** 
```json
{
  "error_code": 0,
  "error_message": "SUCCESS",
  "data": {
    "initiator_address": "0x3bC7726c489E617571792aC0Cd8b70dF8A5D0e22",
    "target_address": "0x97Cd7291f93F9582Ddb8E9885bF7E77e3f34Be40",
    "token_address": "0xB31567308AD3c42D864FB41684bB40d3A2c57E1b",
    "amount": 1000000000000000000000,
    "secret": "",
    "lock_secret_hash": "0xd8875761c93aa9b804c42855601326cf722ced2be5d84fdee36c52ced95ba587",
    "expiration": 65,
    "fee": null,
    "is_direct": false
  }
}
```

## Query the transaction status
`GET /api/1/transferstatus/*(token_address)*/*(locksecrethash)* ` 

There are two ways for users to send and receive transactions, that is, synchronous and asynchronous. If the asynchronous mode is used (sync is false, that is, the default mode), the interface can be called to query the status information of the current transaction. Among them, locksecrethash is obtained from the message returned by the asynchronous transfer transaction.

**Example Request :**  

`GET http://{{ip2}}/api/1/transferstatus/0xB31567308AD3c42D864FB41684bB40d3A2c57E1b/0xd8875761c93aa9b804c42855601326cf722ced2be5d84fdee36c52ced95ba587` 

**Example Response :**  
```json
{
    "error_code": 0,
    "error_message": "SUCCESS",
    "data": {
        "key": "0x232cad7f5e8f51f48bd0d15753ac0fb0996923ff42bd99c50a9060de00411d3a",
        "lock_secret_hash": "0xd8875761c93aa9b804c42855601326cf722ced2be5d84fdee36c52ced95ba587",
        "token_address": "0xb31567308ad3c42d864fb41684bb40d3a2c57e1b",
        "status": 1,
        "status_message": "MediatedTransfer is sending target=97cd\nMediatedTransfer sending success\n"
    }
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

## Cancel the transaction
  ` Post /api/1/transfercancel/*(token)*/*(locksecrethash)*`

This interface is for cancellation of a transaction, which the transaction is in a cancelable state.

In the asynchronous transaction transfer process, if the current transaction is in the cancelable state (status code is 1) through the transaction status query, and the waiting time is too long, the interface can be used to cancel the transaction.

**Example Request :**  

`POST /api/1/transfercancel/0xD82E6be96a1457d33B35CdED7e9326E1A40c565D/0xe0f8d65ddb4f70899b97f36795925a97c1b286582f58f56a041f141d345acdca`

**Example Response :**  
**200 OK**
```json
{
    "error_code": 0,
    "error_message": "SUCCESS",
    "data": null
}
```
Note: Before using this interface, you need to query the corresponding transaction status through the interface `/api/1/transferstatus`. If it is not in the cancelable state, the interface will return an Error:"can not found transfer".

## Token exchange
  ` PUT /api/1/token_swaps/*(target_address)*/*(lock_secret_hash)*`

Token Swap can be used to atomic exchange within two types of tokens.


 Under the circumstances that valid routing strategies are existed, first invoke `taker` then `maker`,  It should be noted that the preimage of the `lock_secret_hash` which the maker introduced must be equal to the `secret` when taker adopted. Note that both taker and maker request the lock_secret_hash, the secret was given in the  maker's request parameters.

Note: With help of the interface  `/api/1/secret`,  channel participants can receive a `lock_secret_hash` / `secret` pair.

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

Note: If `tokenswap`is exchanged through direct channels between the two parties, no `route_info` information is needed; otherwise, as with transfer, the route information of the destination should be introduced into `taker`and `maker` requests respectively, and the indirect channel `tokenswap`should be assigned routes and charges.

## Switch to no network

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

**Example Response :**  

  ```json
{
    "error_code": 0,
    "error_message": "SUCCESS",
    "data": null
}
```

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

```json
{
    "error_code": 0,
    "error_message": "SUCCESS",
    "data": "ok"
}
```

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

**Example Response :**  

 ```json 
{
    "error_code": 0,
    "error_message": "SUCCESS",
    "data": "ok"
}
```

## Query the fee policy
`  GET /api/1/fee_policy `

Query the node charging information, which connect to default PFS server. If the rate has been set, return the fee rate information, otherwise, return the default information.

**Example Request :**  

`GET：http://{{ip1}}//api/1/fee_policy`

**Example Response :**  

**200 OK**   

```json  
{
    "error_code": 0,
    "error_message": "SUCCESS",
    "data": {
        "Key": "feePolicy",
        "account_fee": {
            "fee_constant": 5,
            "fee_percent": 10000,
            "signature": "r8WxYRc/Jei7vpy4wGMm2hAikM8enlZibeWy8FQEJut0CRH9gx/ZYA80gfesYYiXYpAl1IMci+UcfT79E9zyARs="
        },
        "token_fee_map": {
            "0xb31567308ad3c42d864fb41684bb40d3a2c57e1b": {
                "fee_constant": 5,
                "fee_percent": 10000,
                "signature": "r8WxYRc/Jei7vpy4wGMm2hAikM8enlZibeWy8FQEJut0CRH9gx/ZYA80gfesYYiXYpAl1IMci+UcfT79E9zyARs="
            }
        },
        "channel_fee_map": {
            "0xfe738aa39610416e4100036130af7ae00930021d5a51be60b55b96c12b1f4af5": {
                "fee_constant": 5,
                "fee_percent": 1000,
                "signature": "EAo6sV0d665BNTrQSWJC8fnO15POkc+sbWYKVV5VQbBf5+o9kPlNbag0InYCJ/FVhTtlYtVGXL5U5WBaGVGEpBs="
            }
        }
    }
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

`GET：http://{{ip1}}/api/1/path/0xC445a8C326A8fD5a3e250C7dc0EFc566eDcB263B/0xB31567308AD3c42D864FB41684bB40d3A2c57E1b/1000000000000000000000`
  
**Example Response :**  

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

### Revenue Detail Query
Post /api/1/income/details

Detailed revenue of query nodes, including fee revenue and direct revenue.

 **Example Request :**  
 
 ```json
{
     "token_address":"0x0000", // Filter by token query
    "from_time":1552901182, // Filtering by the server time of the photon node where the transaction occurred
    "to_time":1552901182, //Filtering by the server time of the photon node where the transaction occurred
    "limit":100, // The maximum number of returned items for this query,which is not limited by default
}
```

**Example Response :**

```json
            {
                "error_code": 0,
                "error_message": "SUCCESS",
                "data": [
                    {
                        "amount": 1,
                        "data": "",
                        "type": "1",
                        "time_stamp": 1552969089
                    },
                    {
                        "amount": 2,
                        "data": "",
                        "type": "1",
                        "time_stamp": 1552969089
                    },
                    {
                        "amount": 5,
                        "data": "",
                        "type": "1",
                        "time_stamp": 1552969090
                    },
                    {
                        "amount": 19999,
                        "data": "",
                        "type": "1",
                        "time_stamp": 1552969091
                    }
                ]
            }
```
  Note: type   0=transfer revenue 1-fee revenue

### N-day Revenue Query

Post /api/1/fee/query

 Detailed revenue of past N-day of query nodes

 **Example Request :**

 ```json
            {
                    "token_address":"0x0000", // Filter by token query
                    "days":7 // Query the revenue statistics of the past 7 days, excluding the current day
            }
```
**Example Response :**

```json
              {
    "error_code": 0,
    "error_message": "SUCCESS",
    "data": [
        {
            "token_address": "0x8fb0e62caa6ec21a6920b769bb35a07e62a0f8bc",
            "total_amount": 20007,
            "days": 3,
            "details": [
                {
                    "amount": 0,
                    "time_stamp": 1553184000
                },
                {
                    "amount": 0,
                    "time_stamp": 1553270400
                },
                {
                    "amount": 0,
                    "time_stamp": 1553356800
                }
            ]
        }
    ]
}
```
### Version query
Get /api/1/version

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



