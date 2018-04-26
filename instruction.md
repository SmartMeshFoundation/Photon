# Smartraiden Instruction
 In order to facilitate the enthusiasts to better test and understand the functions of the SmartRaiden project, this document gives a brief instruction of the program directories, and the detailed functional usage can be found from the notes in the code files.
## Test network construction 
* Build the private chain  
  
&emsp;&emsp;1.Download  the geth.

&emsp;&emsp;2.Creat the genesis.json：You will need to use the baipoatestnet.json to create the creation block of the private chain directly.

&emsp;&emsp;3.Initialize  the block chain：you can proceed the following command to carry out Initialization.

 &emsp; &emsp;&emsp;geth --datadir ~/privnet/ init baipoatestnet.json

&emsp;&emsp;The home directory of the private chain is ~/privnet/, which can be modified according to the circumstances.

&emsp;&emsp;4.Runn the Ethereum node：You will also need to start an Ethereum Node.

&emsp;&emsp;geth --datadir=~/privnet --networkid 8888（private chain ID）

* start the smartraiden nodes. 
 
 &emsp;&emsp;1.install [pjsip](http://www.pjsip.org/)
```
wget http://www.pjsip.org/release/2.7.2/pjproject-2.7.2.tar.bz2
tar xjf pjproject-2.7.2.tar.bz2
cd pjproject-2.7.2
./configure --disable-sound --disable-video --disable-ssl
make dep && make && make install
```
&emsp;&emsp;2.build smartraiden  
```
  go get github.com/SmartMeshFoundation/SmartRaiden/cmd/smartraiden
```

## The structure of the Project
The core of the whole program is RaidenServcie.go, which controls the state of the smartraiden node. And the smartraiden node can be considered as a state machine, which is driven by user requests, chain events and other node’s messages.


![](https://i.imgur.com/5xpaAuW.jpg)

 &ensp;&emsp;&emsp;![](https://i.imgur.com/oovbYQL.jpg)

## The description of the directory

###Abi directory

This directory describes how to call the interface of the smart contract by using Golang, which were introduced from the Ethereum and made some modifications.

![](https://i.imgur.com/yNDYmqZ.jpg)




###Blockchain directory
![](https://i.imgur.com/D5paxo8.jpg)

alarmtask.go is to receive the events for new blocks producing. In fact, the basic unit of time for smartraiden network is block number, that is, this is the time generator of the smartraiden network.

eventslistener.go is used to analyze the smartraiden network events on the block chain. There are the following events:

 * params.NameTokenAdded  
 * params.NameChannelNew  
 * params.NameChannelNewBalance  
 * params.NameChannelClosed  
 * params.NameChannelSettled  
 * params.NameChannelSecretRevealed 

 
TokenAdd: A new registration of token.  
ChannelNew: New Channel is created, including addresses for both parties.  
ChannleNewBalance: Description of which channel has a deposit event, including the depositor and the deposit amount.  
ChannelClosed: Description of  which channel should be tag of closed channel.This event is very critical,if it was missed,token may be lost.  
ChannelSettled: Channel was completely shut down and destroyed, and both paries took back their respective money.  
ChannelSecretRevealed: It represents a secret disclosure on some channel, which generally means that a transaction ends on the block chain instead of off-chain, and the node related to the transaction can withdraw his token on the chain with the secret.  


BlockChainEvents.go: Encapsulate the events on the chain and report it to the RaidenService.  

###Channel directory
![](https://i.imgur.com/bqi99eJ.jpg)

This is the location where the channel information is truly preserved, maintained and verified.

Channel.go: this file is to manage and maintain a channel. The sending transaction information is from here, and the receiving transaction information should be verified here. The information from the channel participants are two ChannelEndState. One channel must keep both participants' latest information. The counterparty also has a channel structure, The stored information between the participants should be identical theoretically.

ChannelEndState.go:The information belongs to one of the channel participants, which include how much token you have paid , how many unfinished transactions you should pend, etc.

Channelexternal.go: The functions that need to be invoked for interacting with the block chain, such as close channel,settle channel, deposit,withdraw.


###CMD directory
![](https://i.imgur.com/iH5YqpA.jpg)

![](https://i.imgur.com/fUGf4Wk.jpg)

Api test directory, the main Api and instructions are as follows:

* QueryingNodeAddress 　　　　　　 　Query a node address 
* QueryingNodeAllChannels　　　　　　Query all channels of a node
* QueryingNodeSpecificChannel　　　&ensp;&ensp;Query a specified channel for a node
* QueryingRegisteredTokens　　　　　&ensp;Query system registration Token
* QueryingAllPartnersForOneTokens　　Query the Partner address in the channel of special Token
* RegisteringOneToken　　　　　　　　Registration of new Token to Raiden Network
* TokenSwaps　　　　　　　　　　　　Exchange the token
* OpenChannel　　　　　　　　　　　&ensp;Establish a channel
* CloseChannel　　　　　　　　　　　&ensp;Close the specified channel for the node
* SettleChannel　　　　　　　　　　　&ensp;Settle the specified channel for the node
* Deposit2Channel　　　　　　　　　　Deposit to the specified channel
* Connecting2TokenNetwork　　　　　&ensp;Connect to a TokenNetwork
* LeavingTokenNetwork　　　　　　　&ensp;Leave the TokenNetwork
* QueryingConnectionsDetails　　　　&ensp;Query the details of the Token network connection 
* QueryingGeneralNetworkEvents　　　Query network events
* QueryingTokenNetworkEvents　　　　Query Token network events
* QueryingChannelEvents　　　　　　&ensp;Query channel event

###Encoding directory
![](https://i.imgur.com/s5ZIU95.jpg)

It mainly defines the format of exchanging messages. If the message structure contains SignedMessage, it is the signature that the message need to send.
There are the following types of messages:  

    const ACK_CMDID = 0  
    const PING_CMDID = 1  
    const SECRETREQUEST_CMDID = 3  
    const SECRET_CMDID = 4  
    const DIRECTTRANSFER_CMDID = 5  
    const MEDIATEDTRANSFER_CMDID = 7  
    const REFUNDTRANSFER_CMDID = 8  
    const REVEALSECRET_CMDID = 11  
    const REMOVEEXPIREDHASHLOCK_CMDID=13  
    
Among them, the mainly important messages are as follows:  
1.SecretRequest message：The receiver of the transaction request the corresponding secret to the initiator.  

    type SecretRequest struct {
       SignedMessage
       Identifier uint64
       HashLock   common.Hash
       Amount     *big.Int
    }

 Hashlock is the corresponding lock of the secret. Amount is the amount of the transaction.  
2.RevealSecret message：This message is to reveal the secret which the hashlock corresponding to.  

    type RevealSecret struct {
    SignedMessage
    Secret   common.Hash
     hashLock common.Hash
    }  
The message originally came from the initiator, but the initiator, the receiver and the intermediate nodes all can send and receive the message.  
 3.EnvelopMessage message:Basic message encapsulation.

    type EnvelopMessage struct {
       SignedMessage
       Nonce          int64
       Channel        common.Address
       TransferAmount *big.Int //The number has been transferred to the other party
       Locksroot      common.Hash
       Identifier     uint64
    }  
This message will not be sent directly, but other messages,such as,Secret, RemoveExpiredHashlock, Transfer,DirectTransfer, MediatedTransfer, RefundTransfer, all of them need to add the basis.   
This message includes the latest state of the current channel participant.   
TransferAmountrepresents how much token which has been transferred.  
Locksroot represents the transaction collection which being carried out.   
Nonce represents the total number of transaction times which the participants carried out. Nonce starts from 1,a monotonically increasing value. TransferAmount starts from 0,also a monotonically increasing value.   
4.DirectTransfer：Information needed for direct transfer.

    type DirectTransfer struct {
       EnvelopMessage
       Token     common.Address  
       Recipient common.Address  
    }  
Parameter Description:  
Toke:The token which direct transfer to the recipient.   
Recipient:The message which does not pass through the intermediate nodes and does not have the HTLCs.  
5.MediatedTransfer:Information needed for Indirect transfer.

    type MediatedTransfer struct {
       EnvelopMessage
       Expiration int64
       Token      common.Address
       Recipient  common.Address
       Target     common.Address
       Initiator  common.Address
       HashLock   common.Hash
       Amount     *big.Int //The number transferred to party
       Fee        *big.Int
       }  
Parameter description:  
Initiator:The sender of the transaction.  
Target:The final receiver of the transaction.  
Recipient:The intermediate receiver.  During the mediated transfer, the Recipient and the Target (token receiver) are most-time different .  
Expiration,hashlock,Amount and Fee are all related to HTLC. Please look up the principles in detail.  
This message is used for token transfer, for example, A-C has a channel, A gives C 10 token using this message, A-B-C has indirect channels, and A can also give C 10 token through this message.  
6.RefundTransfer

    type RefundTransfer struct {
       MediatedTransfer
    }  
 The RefundTransfer is same as the mediatedtransfer.The difference lies in the use of the scenes.  
A gives C 10 token passed B. B receives the mediated message from A, and finds he is unable to give C 10 token, then B needs to return A a RefundTransfer to assured that he mortgage 10 token to A. This will ensure the security of A.  
7.Secret mesage: The secret used to unlock the HTLC.


    type Secret struct {
       EnvelopMessage
       Secret common.Hash
    }  
This message is basically an EnvelopMessage, which only adds a secret. In fact, through this message, the so-called HTLC is replaced by unconditional giving. In order to complete a transaction, the recipient must withdraw the token before the time limit, once the recipient get the secret, the time lock and the hash lock will be cancelled.  
8.RemoveExpiredHashlockTransfer：Remove unfinished HTLC lock.

    type RemoveExpiredHashlockTransfer struct {
       EnvelopMessage
       HashLock  common.Hash
       }  
RemoveExpiredHashlockTransfer: The message is used when a HTLC is not completed, but the recipient and the sender all save the locked transaction, so the sender needs to send the message to remove the lock. And the receiver can get the signature of the new state after the lock is removed.

###Mobile directory  
![](https://i.imgur.com/0zUymXi.jpg)  
This directory is the API interface for mobile platforms ,such as Android and IOS.  

###Models directory  
![](https://i.imgur.com/Zz7rEMn.jpg)    
This directory realizes the related functions of database and stores related information of the transaction.  
It consists of three parts.  
1.Local mapping of on-chain information, such as, channels.go tokens.go, etc.  
2.transaction results,such as, channels.go.  
3.The crash recovery information, which need to preserve state changes caused by almost every step in the transaction.  

###Network directory  
![](https://i.imgur.com/emnpaR3.jpg)  
This directory has many contents, and it is also an important part of the project.  
1.UDP Realization  
Because UDP communication is not reliable, SmartRaiden uses the Message Request and Answer Pattern to confirm the reception of messages. Each message has a corresponding ACK. After receiving the ACK, the sender can confirm that the message is sent successfully.  
protocol.go  is the core to realize reliable communication of UDP.  
The related transporter is based on UDP's communication connection, in which icetransporter is to achieve NAT traversal.The following signal/nat directories are both services for this.  
2.discovery.go :Node discovery on chain.  
The main idea to realize the node discovery is that each node stores its own IP and port information in an smart contract, which are convenient for other nodes to query and communicate with them.  
3.Channelgraph: To solve the routing problem.  
At present, the shortest path is adopted for routing transfer, which requires to know the topology information of the whole nodes. Then the Dijkstra algorithm is used to find the shortest path.

###params  directory  
![](https://i.imgur.com/aWizJut.jpg)   
The directory describes the setting of the parameter type,the naming of the events, the declaration and definition of variables and constants.

###rerr directory  
![](https://i.imgur.com/r4STu26.jpg)  
error.go  defines common error processing function.

###Restful directory  
![](https://i.imgur.com/UhFZMZN.jpg)   
The directory uses the rest API to realize the HTTP interface. It can be viewed according to the API instructions.

###testdata directory  
![](https://i.imgur.com/FmjimPd.jpg)  
The test data are generated by the project test.

###Transfer directory  
![](https://i.imgur.com/nT4A8Gx.jpg)  
This is the core of the project, the system is a large state machine, which divides the state manager into three categories according to the different participants’ roles, that is,  initiators, receivers, and intermediate nodes of the route.  
In contrast, the intermediate nodes in the route are most complex.  
This is mainly about how to handle the message received from the party, as well as the events on the block chain, such as the emergence of a new block.

###Ui directory  
![](https://i.imgur.com/890UVoV.jpg)   
Web UI directory , currently unused but convenient to function extension.

###Utils directory   
![](https://i.imgur.com/LYS2gkv.jpg)  
These are some auxiliary functions which are easy to infer functions based on the names.

###vendor directory   
![](https://i.imgur.com/6uOQ4Ok.jpg)  
These are some tools provided by the third parties.

###home directory  
![](https://i.imgur.com/lPvYwc1.jpg)   
Accounts.go is the management of local store account information which is mainly private key. This file is actually a simple encapsulation of the account management of the Ethereum.  
Connectionmanager.go realize the related functions of Connection in API.  
RaidenApi.go is the encapsulation of the raiden service for external use.  
Snapshot.go is mainly to restore the data after the crash, which includes three aspects:  
The first is the recovery of channel state,  the second is the recovery of state manager, which is a very core issue.The third is the recovery of some unimportant messages to the state manager.  
Raidenservice.go ,eventhandler.go and messagehandler.go are the other three core files.  
The smartraiden node is also driven by message, all the users' requests, on-chain events, and the messages sent by the other party will eventually be converted to events or statedchange.  
Event is generated by statemanager after processing StateChange. It is driven by StateChange.  StateChange mainly comes from users, block chain and messages from the party.  
Among them, the messagehandler receives and processes the messages sent by the other party, the EventHandler handles on-chain events as well as the events triggered by the statechange.  
Raidenservice.go mainly deals with the requests from users. In fact, there are many codes which is no suitable locations to put were laid out here. So it seems more bigger than others.
## Raiden Contract  
![](https://i.imgur.com/dvNPPD3.jpg)  

These are the currently deployed raiden contract addresses for the ethereum testnet:  

* Netting Channel Library: [0xad5cb8fa8813f3106f3ab216176b6457ab08eb75](https://ropsten.etherscan.io/address/0xad5cb8fa8813f3106f3ab216176b6457ab08eb75#code)
* Channel Manager Library: [0xdb3a4dbae2b761ed2751f867ce197c531911382a](https://ropsten.etherscan.io/address/0xdb3a4dbae2b761ed2751f867ce197c531911382a#code)
* Registry Contract: [0x68e1b6ed7d2670e2211a585d68acfa8b60ccb828](https://ropsten.etherscan.io/address/0x68e1b6ed7d2670e2211a585d68acfa8b60ccb828#code)
* Discovery Contract: [0x1e3941d8c05fffa7466216480209240cc26ea577](https://ropsten.etherscan.io/address/0x1e3941d8c05fffa7466216480209240cc26ea577#code)
## Summary
The most difficult part of the whole system is the channel directory,the transfer directory,and some special files,that is, raidenservice.go, messagehandler.go, eventhandler.go.The directories and the contracts need to be cross referenced.So the project is actually a DAPP based on contracts.

