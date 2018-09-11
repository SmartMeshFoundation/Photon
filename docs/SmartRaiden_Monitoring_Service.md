# SmartRaiden Monitoring Service

## Objective

Once payment channels are closed, settlement window pops up. During this period, any participant of non channel-closing requires to invoke updateTransfer, to update states of transaction, and withdraw, to release locks on transactions, or those participants will not get all the tokens sent from their counterparts. Nodes, especially mobile nodes, hardly keep connected to the network all the time, whereby there is required any trusted third-party node to facilitate the process of settlement on behalf of offline nodes. 

Please note that  for all third-party nodes in this article, we name them as **SmartRaiden Monitoring node**, short as **SM node**. 
For all nodes that delegate or entrust their work to SM nodes, we call them **delegator**.

## Requirement

- Any settlement message can be grabbed by SM nodes.  


SM nodes are responsible for updating settlement channel, during period of settlment, on behalf of any participant.For which, any SM node must be capable in the ability of receiving messages of transferring or unlocking from any delegator with signature. 

- Any SM node enables to enforce the process of settlement. 


At the moment that any channel is closed once delegator goes offline, SM node can invoke `updateTransferDelegate` or `withdraw`, to continue settlement process.

-  Any single SM node is **trustworthy**.  


In order to reduce the cost, A single SM node was choosen which has to be honest and without any fraudulent intention. It's presumed that any SM node will not be in coalition with participants that closes the channel and do nothing with fraudulent behavior. If SM nodes act fraudulently, the maximum cost is that delegator will not commit to the underlying blockchain.

## SM Settlement Process

When settlement process occurs, the SM node update balance proofs of one closing payment channel on behalf of delegator, by invoking `updateTransferDelegate`, and to unlock the locks by `withdraw`. Any SM node is permitted with invocation of `updateTransferDelegate` with only one time. Presumed that Alice and Bob are APPs for both sides of the payment channel, SM node is the trusted node for delegation service, and carl is a node for charge (a smartraiden node), SMT as a token for fee to pay for any transaction. Bob attempts to disconnect the network, before which he calls for a SM node to help him with the following procedures. 

SM Settlement Process works as follow.

1. Starup the SmartRaiden Monitoring Service.
2. Bob checks the set of smartraiden nodes, to obtain committed data for SM node.
3. Bob commits to a SM node with a proof of transaction.
4. The SM node accepts requests from delegators,and Bob enquires whether sufficient balances are held in SM nodes.
5. If there is deficient balance in SM accounts, then Bob transfers some amount of SMT to carl, in order to pay for delegation fee. This is the regular smartraiden transfer.
6. Alice closes the payment channel.
7. SM node waits for half the period of settle timeout.
8. SM node commits to the chain with proofs of delegation.
9. After the settle timeout window, Alice invokes `settle`, so that Alice and Bob get their SMT, and Carl gets his fee.

## Data Feed of SM 
One delegator requires to feed his signature to his trusted SM node, in which all the information of `updateTransferDelegate` and `withdraw` gets contained.

## Payment to SM 
Right now, each payment for SM node is to use the SMT. Delegator has to check balance proofs of SM nodes. If there is no sufficient balance for any transaction, delegator is required to transfer enough SMT to SM node. Once SM node has done his job, SM node takes his part of charge and get some amount of SMT as his fee.

## Example on SM Service
Next, we are going to demonstrate an entire functionality for SmartRaiden Monitoring Service 

### Preparation 

#### Install Scripts of SM nodes
Install SmartRaiden Monitoring :
```
go get github.com/SmartMeshFoundation/SmartRaiden-Monitoring
cd cmd/smartraidenmonitoring
go install
```

#### Start Up Nodes 
To start up a SmartRaiden Service, there are at least two nodes required, one is Carl - the charge node for smartraiden, the other is the SM node.

For Example :   

- **Carl**    0x69c5621db8093ee9a26cc2e253f929316e6e5b92  
- **SM Node** 0x6b9e4d89ee3828e7a477ea9aa7b62810260e27e9

**start charge node Carl :**

```sh
#!/bin/sh
echo run smartraiden...
smartraiden   --address="0x69c5621db8093ee9a26cc2e253f929316e6e5b92"  --password-file /password-file-path  --eth-rpc-endpoint ws://127.0.0.1:18546 --ignore-mediatednode-request
echo "quit ok"
```
Carl monitors port 5001, to receive charge fees. `ignore-mediatednode-request` is required here, so that Carl would not be viewed as a route node, in the event that any SM node is mis-charged.

**start up sm node :**

```sh
#!/bin/sh
echo run smartraiden-monitoring...
smartraidenmonitoring  --password-file /password-file-path   --eth-rpc-endpoint ws://127.0.0.1:18546  --address="0x6b9e4d89ee3828e7a477ea9aa7b62810260e27e9"
echo "quit ok"
```
`--smt` : If SMT token is not used, then add flag `--smt`. After SM node started, it will automatically monitor port 6000.

#### SmartRaiden Monitoring Service in Use
Once SM service has launched, APP of smartraiden has the ability to enforce delegation service. Take Alice & Bob as our example here, we briefly explain how it works.

- Alice 0x31ddac67e610c22d19e887fb1937bee3079b56cd
- Bob 0xf0f6e53d6bbb9debf35da6531ec9f1141cd549d5

At the first step, we can check balance proofs for Alice and Bob, to verify whether smartraiden monitoring service has worked, or whether tokens in channels has already been settled, or that is there any change for balances on chain.

### Delegation
In our example, Bob is a **delegator** which submits delegation proofs to SM node. To verify that data of delegation from Bob is in the newest version, first we check the balances in the payment channel of Alice & Bob (it is not conpulsory).

**1. Bob checks local chennel information of smartraiden nodes.**

Via API offered below  
**`GET  /api/<version>/channels/<channels_address>`**  

Example Request :   
`
GET:http://localhost:5003/api/1/channels/0x8e537C30913A76C33a3A890a6aFc644f62F97B98
`  

Example Response :   
*`200 OK`* and 

```json
{
    "channel_address": "0x8e537C30913A76C33a3A890a6aFc644f62F97B98",
    "partner_address": "0xf0f6E53d6bbB9Debf35Da6531eC9f1141cd549d5",
    "balance": 200,
    "patner_balance": 100,
    "locked_amount": 0,
    "partner_locked_amount": 0,
    "token_address": "0x745D52e50cd1b19563D3a3B7B6d2eB60b17E6bAE",
    "state": "settled",
    "settle_timeout": 100,
    "reveal_timeout": 10,
    "ClosedBlock": 3850465,
    "SettledBlock": 3850645,
    "OurUnkownSecretLocks": {},
    "OurKnownSecretLocks": {},
    "PartnerUnkownSecretLocks": {},
    "PartnerKnownSecretLocks": {},
    "OurLeaves": null,
    "PartnerLeaves": null,
    "OurBalanceProof": {
        "Nonce": 2,
        "TransferAmount": 100,
        "LocksRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "ChannelAddress": "0x8e537c30913a76c33a3a890a6afc644f62f97b98",
        "MessageHash": "0xcaed88275255e3b02e77c113b15403dc4b70d94bf4ec0cd8293b22610008a801",
        "Signature": "1HWG94xxMNF4yKt9Crur7HgsYE85dfI93ReegbYqSW9Csp/0cdzjEQK58NL7zzLZ/pdqm4Pwk7jAgjhsLoQIhhs="
    },
    "PartnerBalanceProof": null,
    "Signature": null
}

```

**2. Bob checks local delegation information of smartraiden nodes.**

Via API below :   
**`GET: /api/<version1>/thirdparty/<channel_address>/<thirdparty_address>`**

Example Request :   
`GET:http:// localhost:5003/api/1/thirdparty/0x05c468707fBdf56f944d6292fa16234167c704f0/0x6b9e4d89ee3828e7a477ea9aa7b62810260e27e9`

Example Response :   
*`200 OK`* and   
```json
{
    "channel_address": "0x8e537C30913A76C33a3A890a6aFc644f62F97B98",
    "update_transfer": {
        "nonce": 2,
        "transfer_amount": 100,
        "locksroot": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "extra_hash": "0xcaed88275255e3b02e77c113b15403dc4b70d94bf4ec0cd8293b22610008a801",
        "closing_signature": "d47586f78c7130d178c8ab7d0abbabec782c604f3975f23ddd179e81b62a496f42b29ff471dce31102b9f0d2fbcf32d9fe976a9b83f093b8c082386c2e8408861b",
        "non_closing_signature": "54e486d55ef9830d1251217be1a9c3057318fff930e5d28066a5b617ec244c4244bc2b8fc4ecea244493c824a7234e8267e77a1a10e633729ed6f276f14859f81b"
    },
    "withdraws": null
}
```

**3. Bob sumbits delegation proofs to SM nodes.**

Via API below :   
**`POST: http://<sm_server_address:port>/delegate/<address_of_delegator>`**

> comment : `<sm_server_address:port>` is a smartraiden monitoring port (default to be local 6000 port).

**Example Request** :   
`POST:http://192.168.124.13:6000/delegate/0xf0f6E53d6bbB9Debf35Da6531eC9f1141cd549d5`

with payload : 
```json
{
    "channel_address": "0x8e537C30913A76C33a3A890a6aFc644f62F97B98",
    "update_transfer": {
        "nonce": 2,
        "transfer_amount": 100,
        "locksroot": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "extra_hash": "0xcaed88275255e3b02e77c113b15403dc4b70d94bf4ec0cd8293b22610008a801",
        "closing_signature": "d47586f78c7130d178c8ab7d0abbabec782c604f3975f23ddd179e81b62a496f42b29ff471dce31102b9f0d2fbcf32d9fe976a9b83f093b8c082386c2e8408861b",
        "non_closing_signature": "54e486d55ef9830d1251217be1a9c3057318fff930e5d28066a5b617ec244c4244bc2b8fc4ecea244493c824a7234e8267e77a1a10e633729ed6f276f14859f81b"
    },
    "withdraws": null
}
```

**Example Response** :   
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

**4. Bob checks whether sufficient balances are in the SM node.**  

Via API below :   
**`GET: http://<sm_server_address:port>/fee/<address_of_delegator>`**

**Example Request** :   
`GET http://192.168.124.13:6000/fee/0xf0f6E53d6bbB9Debf35Da6531eC9f1141cd549d5`  

**Example Response** :   
*`200 OK`* and
```JSON
{
    "Available": 100,
    "NeedSmt": 3
}
```

Only does Bob have sufficient balance that he can make a successful delegate and make payments. If Bob has deficient balance, then SM node may accept delegation from Bob but will not attempt to invoke `updateTransferDelegate`. Hence, Bob has to transfer funds to the charge node, in this example it's carl. This is the normal smartraiden transfer.

**5. Bob disconnects to the network.**  

Bob gets disconnected from the network, then he is in offline state. When using delegation service, Bob must be unconnected to the network, if not, he will automatically submit balance proofs on chain while Alice closes the payment channel, which eventually will lead to SM node having no chance to submit proofs.

**6. Alice closes payment channel.**  

Alice closes the payment channel, and submits her evidence.

**7. SM waits for half the period of settle timeout.**  

SM nodes are responsible for monitoring the event of `updateTransfer`. If Bob goes online and invokes `updateTransfer` within T blocks (half of settle timeout), then SM nodes will submit no delegation proof. If there is no trace that Bob invokes `updateTransfer` on his own at the moment of T blocks coming in, and funds in the entrusted SM node are enough for delegation service, Then SM node will invoke `updateTransferDelegate` and `withdraw` to commit to delegation proof on chain.

**8. Alice settles the channel.**  

Alice waits for the time of settle timeout, then she starts to settle the channel. After the process of channel settlement, Alice and Bob both get their token. According to verify on-chain balances of Alice and Bob, they can check whether tokens are correctly transferred. (one can also check balances on his wallet, and compares the value with on-chain balance.)

Original Balance :   
Alice = 4989125  
Bob = 4999795

Balance after trasaction :   
Alice = 4989325  
Bob = 4999895

According to our instance, Bob submits delegation proofs to SM nodes, and gets disconnected. Once Alice closes payment channel, SM nodes successfully commit balance proofs. After channel settlement, Bob gets tokens without any fault, and funds are secured by this mechanism.