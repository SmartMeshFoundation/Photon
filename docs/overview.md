# Smartraiden Overview
## Introduction
Smartraiden is a standard compliant implementation of the Raiden Network protocol using Golang, which enforces small-amount off-chain transactions on the mobile platform, and is able to run without Internet connection. For users can have a in-depth understanding about Smartraiden and feel free using it, we offer an outline briefly introducing relative concepts and procedures. At the beginning, there is a review of the Raiden Network.Next, we introduce primary functions and characteristics of Smartraiden. Then, we will give a showcase of fund-transferring transaction in the Smartraiden. In further detail, we have API specification and standard.

## Review Raiden Network
Raiden Network is an off-chain scalability solution enforcing erc-20 compliant token transferring on the Ethereal. It allows for secure token transactions among participants without any global consensus mechanism, which is implemented through pre-set on-chain deposits transferring with digital signature and lock hash. We still rely on several on-chain processes to open and close a payment channel within a pair of nodes, so that it is incredibly hard for every pair of nodes on the network to create channels. However, if there exists one channel (at least), connecting two nodes through other nodes in the network, then we have no need to create another individual channel for these two nodes. This network is named as the Raiden Network, with all the contracts as to route algorithms and interlock channel communications. 

Figure Payment Channel Network:  
![](/docs/images/Smartraiden_network.png)

## Primary Functionalities and Characteristics in SmartRaiden
The primary goal for SmartRaiden aims to construct a structure to enforce an off-chain scalability solution for SmartRaiden Network, which improves usability, compatibility, and security. 
 
Conventional functions include queries, registrations, channel dependencies, and transfers in different scenarios, as detailed in [rest_api](https://github.com/SmartMeshFoundation/SmartRaiden/blob/master/docs/rest_api.md).

Additional functions include :
- Multiplatform & Mobile Adaptability  

SmartRaiden network will be available on multiple platforms and decentralized micropayment on smart mobile devices can be realized. SmartRaiden currently can work on Windows, Linux, Android, iOS etc. SmartRaiden builds its own messaging mechanism on XMPP, not P2P, and separate nodes and start-up processes, making sure that it is capable of running on multiple platform with correct operations.  

- Nodes State Synchronization

To make transaction secure, SmartRaiden adopts state machine to the design of nodes, ensuring relevant operations are atomic. For instance, it must be consistent with information of received unlock record of data and information sending out in the ACK message, both or neither are successful, no medium state existed. In the process of transactions, if any faulty condition occurs, ensure transaction state of both parties are consistent, and after crash recovery, either transaction continues or transaction fails, without any token loss.

- Internet-free Payment 

It is a special functionality added in the SmartRaiden. Via network construction functions in meshbox, SmartRaiden is able to enforce off-chain fund transferring without reliance on the Internet.

- Third-party Delegation 

 Third-party delegation service, also known as SmartRaiden Monitoring, mainly used to facilitate mobile devices to enforce UpdateTransferDelegate & WithDraw on the blockchain by third-party delegation when they’re offline. Third-party service interacts with three parts in the system of out service, App, SmartRaiden, and spectrum. 

- Fixed-rate Charge 

 Similar to Lightning Network, we have an additional fixed-rate charge function in the process of transferring tokens. Incentivized by this charge, all the nodes on this route will retain channel balance to improve the efficiency and higher the rate of successful transactions.

## Showcases of Transactions in SmartRaiden
 Assume that we have one node using AET token connected to our channel network, in which case, this node connects to another 5 nodes, and easy to transfer tokens to direct nodes. If this channel network gets complicated, then we have our tokens transferred though several nodes, and the state of nodes of this channel will alter successively. 

Showcase is as follows :

```
                              

+ --- + (200)          (100)   + --- +  (100)           (100)  + --- +
|node1| <----- channel ----->  |node2|  <----- channel ----->  |node3|
+ --- + (150)          (150)   + --- +  (100)           (100)  + --- +
                             (150)|(100)
                                  |
                                  |
                                  |
                                  |
                                  |
                             (150)|(100)
                               + --- +  (50)            (150) + --- +
                               |node4|  <----- channel -----> |node5|
                               + --- +  (100)           (100) + --- +

```                               

In this diagram, the addresses of each node are

- node 1 : `0x69C5621db8093ee9a26cc2e253f929316E6E5b92`  
- node 2 : `0x31DdaC67e610c22d19E887fB1937BEE3079B56Cd ` 
- node 3 : `0xf0f6E53d6bbB9Debf35Da6531eC9f1141cd549d5 ` 
- node 4 : `0x6B9E4D89EE3828e7a477eA9AA7B62810260e27E9 ` 
- node 5 : `0x088da4d932A716946B3542A10a7E84edc98F72d8`

And our transaction process starts from node 1, ends at node 5. However, we have only one route to 5, in the diagram, namely 1 to 2 to 4 to 5. After transaction completes, the alteration of balance tokens in this channel is

node 1 to node 2 : `0xc4327c664D9c47230Be07436980Ea633cA3265e4`  
**node 1 initial deposit** : `200 `  
**node 2 initial deposit** : `100`  
**node 1 balance** : `150 `   
**node 2 balance** : `150`  

node 2 to node 3 : `0xd1102D7a78B6f92De1ed3C7a182788DA3a630DDA`  
**node 2 initial deposit** : `100`  
**node 3 initial deposit** : `100`  
**node 2 balance** : `100`    
**node 3 balance** : `100`

node 2 to node 4 : `0xdF474bBc5802bFadc4A25cf46ad9a06589D5AF7D`  
**node 2 initial deposit** : `200`  
**node 4 initial deposit** : `100`   
**node 2 balance** : `150`    
**node 4 balance** : `150`  

node 4 to node 5 : `0xd5CF2248292e75531d314B118a0390132bc7a5F0`  
**node 4 initial deposit** : `100`  
**node 5 initial deposit** : `100`  
**node 4 balance** : `50`  
**node 5 balance** : `150`  

## SmartRaiden Contract and Channel lifecycle 
SmartRaiden contract includes :   
- Netting Channel Library : `0xad5cb8fa8813f3106f3ab216176b6457ab08eb75`  
- Channel Manager Library : `0xdb3a4dbae2b761ed2751f867ce197c531911382a`  
- Registry Contract : `0x68e1b6ed7d2670e2211a585d68acfa8b60ccb828`  
- Discovery Contract : `0x1e3941d8c05fffa7466216480209240cc26ea577`

Spectrum contract registry address  = `0x41Df0be8c4e4917f9Fc5F6F5F32e03F226E2410B`

### Channel lifecycle
- Channel nonexistence 

There are two cases for channel nonexistence : one is our channel never exists, the other is we have already settled our transaction, so that all the data of channel and participants have been removed. Under both cases, we can not verify transactions among nodes, except that we create channels for our transactions.

- Channel open

There is channel creation between a node and its directly connected counterpart, channel creation operator has the right to indicate addresses of tokens, counterpart, and the number of token to deposit, and time period for settlement. Once channel opens, whereby participants can make their transactions.

- Channel deposit

We have only one node made deposits, after payment channel opens, so that only this node can transfer his tokens to his counterpart. Then this node can send a message informing that there is a payment channel opened for transaction to the counterpart, after which the counterpart is also able to deposit its tokens.

- Channel transfer

Once a node has connected to the payment channel network, by AET token, under which it has access to another 5 nodes. It is quite easy for this node to transfer its token to another directly-connected node, but if it wants to transfer to the intermediate nodes among them, they both need to construct channels to the intermediate nodes, and if tokens in these nodes is sufficient for this transaction, then transaction occurs.


- Channel close

If any node wants to shut down a certain channel connected to it, he can invoke close function. After that, channel close operator and its counterpart need to submit their balance proofs during the settlement. 

- Channel settle

Once payment channel close is invoked, settle timeout starts to count. During this period, both nodes submit the most recent message. After timeout, channel finishes the settlement.

## Conclusion 

At here, you have finished your learning about all the concepts and functionality specifications about SmartRaiden. For further usage, please go through installation [instruction](https://github.com/SmartMeshFoundation/SmartRaiden/blob/master/docs/installation_guide.md) and [tutorials](https://github.com/SmartMeshFoundation/SmartRaiden/blob/master/docs/api_walkthrough.md) or [SmartRaiden Specification](https://github.com/SmartMeshFoundation/SmartRaiden/blob/master/docs/spec.md)




