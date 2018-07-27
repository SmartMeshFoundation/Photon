# The Principle for Third-Party Settlement
## Objective 
Once functions to close payment channel has been invoked, the settlement window will pop up. Within this period, nodes, who  is not the one for closing payment channel, requires to invoke updateTransfer (to update transaction states of the counterpart) and withdraw (to unlock transfer and withdraw token), or a part of money that sender transfers to will return to who it belongs. It is quite difficult to keep nodes, especially mobile nodes, online all the time. When they go offline, we need a mechanism that offline nodes are able to delegate settlement services to another trusted online third-parties.

## Requirement

### All third-parties enable to receive messages of settlement

Third-parties need to update settlement channel on behalf of their delegators, so that it is necessary for them to get messages of transactions and unlock with signature from sender.

### All third-parties have access to relevant functions of settlement
Once our delegators go offline and payment channels have been closed, our third-parties can invoke updateTransfer and withdraw on behalf of their delegators.

###  All third-parties must be trustworthy

For the purpose of convenience in implement and usage, all third-parties require to be trustworthy and honest, to avoid any case that they are in collusion with the channel-closer and do nothing when fraud behaviors occur. The most serious situation is that nodes will not have any response once their counterparts close the payment channel on the condition that third-parties is dishonest.


## The Design Principle of Third-Party Settlement Channel

Once settlement window pops up, third-parties will invoke updateTransferDelegate on behalf of their offline delegators to update transaction states of channel-closers, then it will invoke withdrawDelegate to unlock transaction. It is permitted that any third-party can invoke updateTransferDelegate in multiple times.

![](/docs/images/third_party_settlement.png)

### updateTransferDelegate()

updateTransferDelegate() is designed to empower a third-party with the ability to update transaction state. This method can be invoked by the third-party in multiple times.

![](/docs/images/updateTransferDelegate.png)

### withdrawDelegate() 

withdrawDelegate() is designed to allow the third-party to unclock transaction. It requires signature from the delegator, for the reason that updateTransfer and updateTransferDelegate can be invoked, respectively, by our delegator and the third-party node. Both methods will update the value of transferred_amount in the transaction but the number of tokens unlocked by withdrawDelegate has to be added onto unlocked_amount, not transferred_amount mentioned above. When settle() is invoked, after settlement window close, we have:

`transferred_amount = transferred_amount + unlocked_amount
`

![](/docs/images/withdrawDelegate.png)

##  Empower / Revoke Third-Party Delegation
SmartRaiden provides new APIs so that any node on the network can delegate and revoke their jobs to the third-party node. Once delegated, third-party node is able to obtain messages fed by his delegator. When delegation has been revoked, third-party node has no access to those messages. 
## Feed Data to Third-Parties
Nodes requires to feed messages to their third-parties,these messages contains whatever needed by updateTransferDelegate and withdrawDelegate
## Payment of Third-Party Fees
SMT, as a unified digital monetary for SmartRaiden, is applied to pay for third-parties. Via the app layer, any node can transfer SMT to his third-party miners.

![](/docs/images/third_party_fees.png)

#  Start up SmartRaiden Monitoring(SM)
## Preparation 
### 1. Nodes
To start up a smartraiden monitoring, one need create at least 4 nodes. One of them is to charge fees for monitoring service, one to play as a monitoring node, and another two nodes to make transactions. 

For instance : 
* node 1 `0x69c5621db8093ee9a26cc2e253f929316e6e5b92`
* node 2 `0x31ddac67e610c22d19e887fb1937bee3079b56cd`
* node 3 `0xf0f6e53d6bbb9debf35da6531ec9f1141cd549d5`
* node 4 `0x6b9e4d89ee3828e7a477ea9aa7b62810260e27e9`

node 1 : charge node
node 2 : sender node
node 3 : receiver node 
node 4 : monitoring node

Check up the balances of node 2 and node 3, in case that later we can verify whether our smartraiden monitoring operates successfully, and whether balances of nodes alter.  
**`GET: /api/<version>/debug/balance/<token>/<address>`**

Check-up Balance

 **Example Request**:
 `GET:http://localhost:5002/api/1/debug/balance/0x745D52e50cd1b19563D3a3B7B6d2eB60b17E6bAE/0xf0f6E53d6bbB9Debf35Da6531eC9f1141cd549d5`  

**Example Response**:  
*`200 OK`* and 
```js
4999795
```

 **Example Request**:  
  `GET:http://localhost:5002/api/1/debug/balance/0x745D52e50cd1b19563D3a3B7B6d2eB60b17E6bAE/0x31DdaC67e610c22d19E887fB1937BEE3079B56Cd`  

**Example Response**:  
*`200 OK`* and 
```js
4989425
```

Hence, there is the current state of node 2 and node 3.  
node 2 balance ：4989425  
node 3 balance ：4999795  

### 2. Install scripts of SM nodes

**Install SmartRaiden Monitoring :**  
```
go get github.com/SmartMeshFoundation/SmartRaiden-Monitoring
cd cmd/smartraidenmonitoring
go install
```

**Install charge node :**  
```sh
#!/bin/sh
#rm -rf .smartraiden
#echo privnet
#echo install ...
#go install
#set-title 1·
echo run smartraiden...
smartraiden --datadir=~/niexin/smartraiden --api-address=0.0.0.0:5001 --listen-address=127.0.0.1:40001 --address="0x69c5621db8093ee9a26cc2e253f929316e6e5b92" --keystore-path ~/.ethereum/keystore --registry-contract-address 0xA5a0E448ded405d86291D37A5561e91F72601751 --password-file 123  --eth-rpc-endpoint ws://127.0.0.1:5555  --conditionquit "{\"QuitEvent\":\"1EventSendMediatedTransferAfter\"}"  --debug --debugcrash   --verbosity 5  --ignore-mediatednode-request
#--enable-health-check  #--signal-server 182.254.155.208:5222
echo "quit ok"
```

Node 1 is necessary to add `ignore-mediatednode-request`, so that it would not work as a route node for raiden, preventing cases like smartraiden monitoring nodes are charged mistaknly.  

**Install smartraiden monitoring node :**   
```sh
#!/bin/sh
echo run smartraiden-monitoring
smartraidenmonitoring  --datadir=/home/niexin/niexin/.smartraidenmonitoring --eth-rpc-endpoint ws://127.0.0.1:5555  --address="0x6b9e4d89ee3828e7a477ea9aa7b62810260e27e9" --keystore-path ~/.ethereum/keystore --registry-contract-address 0xA5a0E448ded405d86291D37A5561e91F72601751  --password-file 123 --verbosity 5  --debug   --smt 0x745D52e50cd1b19563D3a3B7B6d2eB60b17E6bAE
echo "quit ok"
```
`--smt` : tokens to be transferred.
***

## Procedures for Test
### 1. Node 2 sends 20 tokens to node 3

Via API offered below   
**`POST  /api/<version>/transfers/<token_address>/<target_address>`**  

**Example Request**:  
`
http://{{ip2}}/api/1/transfers/0x745D52e50cd1b19563D3a3B7B6d2eB60b17E6bAE/0xf0f6E53d6bbB9Debf35Da6531eC9f1141cd549d5
`  

with payload:
```json
{
    "amount":100,
    "fee":0,
    "is_direct":true
}
```

   **Example Response**:  
*`200 OK`* and 
```json
{
    "initiator_address": "0x31DdaC67e610c22d19E887fB1937BEE3079B56Cd",
    "target_address": "0xf0f6E53d6bbB9Debf35Da6531eC9f1141cd549d5",
    "token_address": "0x745D52e50cd1b19563D3a3B7B6d2eB60b17E6bAE",
    "amount": 100,
    "identifier": 13450803429907647519,
    "fee": 0,
    "is_direct": true
}
```

### 2. Check Delegation Messages
Node 3 checks any message from other nodes who require SM monitoring service.  

**`GET: /api/<version1>/thirdparty/<channel_address>/<thirdparty_address>`**  

**Example Request**:  
`GET:http:// localhost:5003/api/1/thirdparty/0x05c468707fBdf56f944d6292fa16234167c704f0/0x6b9e4d89ee3828e7a477ea9aa7b62810260e27e9`  

**Example Response**:  
*`200 OK`* and 
```js
{
    "channel_address": "0xd1102D7a78B6f92De1ed3C7a182788DA3a630DDA",
    "update_transfer": {
        "nonce": 1,
        "transfer_amount": 100,
        "locksroot": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "extra_hash": "0x8cb02ef8009377267b03aaa1793cbc3c6b2fbb06f58f59bf2ee919108999f192",
        "closing_signature": "a664e8c5f09d32700e2f0110a4aa230c7d3af5a62fa86b01f4e0210c53a7627676e1d4317f39642fcb755ab469ccf43f2d0f90c3decc724d3f37aa7ba43e400b1c",
        "non_closing_signature": "abf712eb853f3b0016ee3af14c5902e0ff72efca03299b10bf50b4b2cda763993033c934fbd01f460ced62ba748b3df75de6ae6be1d07ae99a869b59f35bc84d1c"
    },
    "withdraws": null
}
```
### 3. Commit Proofs of Delegation to SM

Commit proofs of delegation to smartraiden monitoring nodes :  

**`POST: http://192.168.124.13:6000/delegate/<address_of_delegator>`**  
> comment : `192.168.124.13:6000` is a smartraiden monitoring port (default to be local 6000 port).

**Example Request**:  
`POST:http://192.168.124.13:6000/delegate/0xf0f6E53d6bbB9Debf35Da6531eC9f1141cd549d5`  

with payload:
```json
{
    "channel_address": "0xd1102D7a78B6f92De1ed3C7a182788DA3a630DDA",
    "update_transfer": {
        "nonce": 2,
        "transfer_amount": 100,
        "locksroot": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "extra_hash": "0x72173e7a2c119683375e2f78037f9dcf6eaeefb518260abb7bb96edbea256eb1",
        "closing_signature": "f118f7d1ae1db9048c42b1897f7862676143269873ca9f6d9439fde9fdf3747404f3a0cb2a8913d9cebed04604c923726c9b9074fc880a746bcc54561f2cbc591b",
        "non_closing_signature": "281eb8f59a48759b0090c3e76bb88bae70ed18a2dd707de9504483a04eb08fc128e45bbf0dfc2323b4ca92177e501d16af8cb9332ca3ac3d2060a4fe34cbdf331b"
    },
    "withdraws": null
}
```

  **Example Response**:  
  *`200 OK`* and 
  ```json
  {
    "Status": 3,
    "Error": ""
}
```
Status tag has 3 values : 
- status = 1 : delegation failure
- status = 2 : delegation success but lack of fund
- status = 3 : delegation success with sufficient fund

### 4. Verify whether charge node has sufficient fund for node 3 to operate. 

To check whether delegator node (node 3) deposits enough tokens in charge node (node 1) :  

**`GET: http://192.168.124.13:6000/fee/<委托方地址>`**  

**Example Request**:  
`GET http://192.168.124.13:6000/fee/0xf0f6E53d6bbB9Debf35Da6531eC9f1141cd549d5`  

**Example Response**:  
*`200 OK`* and 
```JSON
{
    "Available": 50,
    "NeedSmt": 3
}
```

### 5. Node 3 disconnect  


### 6. Node 2 close up payment channel.

**Example Request**:  
` PATCH http:// localhost:5002/api/1/channels/0xd1102D7a78B6f92De1ed3C7a182788DA3a630DDA`  
 
with payload
```js
{"state":"closed"}
```
**Example Response**:  
*`200 OK`* and 
```json
{
    "channel_address": "0xd1102D7a78B6f92De1ed3C7a182788DA3a630DDA",
    "partner_address": "0xf0f6E53d6bbB9Debf35Da6531eC9f1141cd549d5",
    "balance": 0,
    "partner_balance": 200,
    "locked_amount": 0,
    "partner_locked_amount": 0,
    "token_address": "0x745D52e50cd1b19563D3a3B7B6d2eB60b17E6bAE",
    "state": "closed",
    "settle_timeout": 100,
    "reveal_timeout": 0
}
```

### 7. Check on-chain account balance

`http://localhost:5002/api/1/debug/balance/0x745D52e50cd1b19563D3a3B7B6d2eB60b17E6bAE/0xf0f6E53d6bbB9Debf35Da6531eC9f1141cd549d5`  
```json
4999795
```
We can conclude that node 3 can submit its delegation message, and when it goes offline, then node 4, as a third-party service node, will autonomously enforce UpdateTransfer & WithDraw. After node 2, as the counterpart of node 3, closes payment channels and settles,
node 3 could safely get all the refund in this channel.
