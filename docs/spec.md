# SmartRaiden Specification
## Introduction
This document, just like the SmartRaiden protocol is a constant work in progress

SmartRaiden is a payment network built on top of the ethereum network. The goal of the SmartRaiden project is to provide an easy to use conduit for off-chain payments without the need of trust among the involved parties.

While there are plans to extend SmartRaiden to generalized state channels and channels with multiple parties, this documentation concerns only off-chain payment channels.

#### How does the SmartRaiden Network provide safety without trust?
To achieve safety, all value transfers done off-chain must be backed up by value stored in the blockchain. Off-chain payments would be susceptible to double spending if that was not the case. The payment channel is represented in the blockchain by a smart contract which:
- Provides shared rules, agreed up-front by both parties, for the channel's operation.
- Holds the token value in escrow to back the off-chain payments.
- Arbitrates disputes using rules that cannot be abused by one party.

Given the above properties SmartRaiden can safely do off-chain value transfers, knowing that once a dispute happens the smart contract can be used to settle and withdraw the token.

## The Netting Channel Smart Contract
The netting channel smart contract is the executable code that contains the shared rules for operating an off-chain payment channel. These rules are implicitly agreed upon by each participant whenever a channel is used. The netting channel allows for:

- A large number of bidirectional value transfers among the channel participants.
- Conditional value transfers that have an expiration and predefined rules to withdraw.
- Rules to determine ordering of transfers.

Each netting channel backs a bi-directional off-chain payment channel. They deal with a predetermined token and each has its own settlement period configuration. Any of the two participants may deposit any number of times, any amount of the specified token.

Transfers may be conditionally finalized, meaning that at any given point in time there may be multiple in-flight transfers waiting to be completed. These transfers are represented by lock structures that contain a token amount, expiration, and hashlock. The set of all pending transfers is encoded in a merkle tree and represented in each transfer by its root.

The channel capacity is equal to the total deposits by both participants. The capacity is both the largest value a transfer may have and the total amount of token in pending transfers. The capacity is divided as available and locked balance to each participant/direction. The available balances vary during the lifetime of the channel depending on the direction and value of the completed transfers. It can be increased either by a participant’s deposit or by a counterparty’s payment. The locked balance depends on the direction and value of the pending locked transfers. It is increased with each locked transfer and decreased when the transfer is finalized, successfully or otherwise.

### A channel’s life cycle
1. Deployment
2. Funding / Usage  
3. Close
4. Settle

After being deployed the channel may receive multiple deposits from either participant. Once the [counterparty]() acknowledges it, the depositor may do transfers with the available balance.

Once either party wants to withdraw their tokens or a dispute arises the channel must be closed. After the close function is called the [settlement window]() opens. Within the settlement window both participants must update the counterparty state and withdraw the unlocked locks. A party can not perform a partial withdrawal.

The `updateTransfer()` function call receives a signed balance proof which contains an envelope with channel specific data. These are the [merkletree root](), the [transferred amount](), and a nonce. Since a node can only provide a signed message from the counterparty we know the data wasn’t tampered with and that it is valid. To disincentivize a node from providing an older message, withdraw balances are netted from the transferred amount, a monotonically increasing value. As a consequence there are no negative value transfers and if a participant provides an older message the wrongdoer’s netted balance will end up being smaller.

Another netting channel operation is the lock withdrawal. It receives an unlock proof composed of the lock data structure, a proof that this lock was contained in the merkle tree and the secret that unlocks it. The channel validates the lock, checks the containment proof by recomputing the merkle tree root and checks the secret. If all checks pass the transferred amount of the counterparty is increased.

### Balance Proofs
The netting channel requires a [balance proof]() containing the information to properly settle. These are:
- A nonce
- The transferred amount
- The root node of the pending locks merkle tree
- A signature containing all the above

For this reason each transfer must be encoded as a balance proof, this follows from the fact that transfer messages change the node balance and must be provable to the netting channel.

## SmartRaiden Transfers
Transfers in SmartRaiden come in three different flavors.
### Direct Transfers
A [DirectTransfer]() does not rely on locks to complete. It is automatically completed once the network packet is sent off. Since SmartRaiden runs on top of an asynchronous network that can not guarantee delivery, transfers can not be completed atomically. The main points to consider about direct transfers are the following:

- The messages are not locked, meaning the envelope transferred_amount is incremented and the message may be used to withdraw the token. This means that a [payer]() is unconditionally transferring the token, regardless of getting a service or not. Trust is assumed among the payer/[payee]() to complete the goods transaction.
- The sender must assume the transfer is completed once the message is sent to the network, there is no workaround. The acknowledgement in this case is only used as a synchronization primitive, the payer will only know about the transfer once the message is received.

A succesfull direct transfer involves only 2 messages. The direct transfer message and an `PROCESSED`. For an Alice - Bob example:
- Alice wants to transfer n tokens to Bob.
- **Alice creates a new transfer with**.
    - transferred_amount = `current_value + n`
    - [locksroot]() = `current_locksroot_value`
    - nonce = `current_value + 1` 

- Alice signs the transfer and sends it to Bob and at this point should consider the transfer complete.

### Mediated Transfers
A [MediatedTransfer]() is a hashlocked transfer. Currently SmartRaiden supports only one type of lock. The lock has an amount that is being transferred, a [hashlock]() used to verify the secret that unlocks it, and a [lock expiration]() to determine its validity.

Mediated transfers have an [initiator]() and a [target]() and a number of hops in between. The number of hops can also be zero as these transfers can also be sent to a direct partner. Assuming `N` number of hops a mediated transfer will require `6N + 8` messages to complete. These are:
- N + 1 mediated or refund messages
- 1 secret request
- N + 1 secret reveal
- N + 1 secret
- 3N + 4 PROCESSED

For the simplest Alice - Bob example:
- Alice wants to transfer n tokens to Bob.
- **Alice creates a new transfer with:**
    - transferred_amount = `current_value `
    - lock = `Lock(n, hash(secret), expiration)`
    - locksroot = `updated value containing  the lock`
    - nonce = `current_value + 1`

- Alice signs the transfer and sends it to Bob
- Bob requests the secret that can be used for withdrawing the transfer by sending a [SecretRequest]() message.
- Alice sends the [RevealSecret]() to Bob and at this point she must assume the transfer is complete.
- Bob receives the secret and at this point has effectively secured the transfer of n tokens to his side.
- Bob sends a [secret message]() back to Alice to inform her that the secret is known and acts as a request for off-chain synchronization.
- Finally Alice sends a secret message to Bob. This acts also as a synchronization message informing Bob that the lock will be removed from the merkle tree and that the transferred_amount and locksroot values are updated.

### Refund Transfers
A [RefundTransfer]() is a mediated transfer used in the special circumstance of when a node cannot make forward progress, and a routing backtrack must be done.

## Third parties
Third parties are required to provide for safe operation. Since a single node cannot be expected to have 100% up-time, third parties are required to operate the netting channels for the period of time the node is offline.

The purpose of a third party is to update the netting channel during settlement on behalf of a participant. For this reason SmartRaiden must be configured to keep the third party up-to-date with its received transfers, locks, and secrets. If a channel is closed while this node is offline then the third party must be capable of calling updateTransfer/withdraw on its behalf.

In order to avoid collusion among third parties and the channel counterparty and protect from DoS attacks, a node cannot rely on only a single third party. Given that a node relies on more than one third party, and that these services won’t have 100% uptime, out-of-sync third parties, that not always have the latest known balance proof, must be handled. The smart contract must have logic to totally order a stream of balance proofs from a single participant to do conflict resolution.

Additional problems can arise with the usage of multiple third parties. Because third parties can easily be impersonated, penalization is not an option for netting channels. If a user is using too many third parties, once a channel is closed there could be a [thundering herd]() problem.

## Mediating Transfers
SmartRaiden cannot rely on direct channels for most of its operations, especially if the majority of them are for target nodes that will only receive a transfer once. Mediated transfers are a form of value transfer that allows trustless cooperation among SmartRaiden nodes to facilitate movement of value.

Mediated transfers rely on locks for safety. Locks can be unlocked only by knowledge of the secret behind it. This information is used to determine whether a transfer was complete and is shared among all nodes in a mediation chain. The lock operation allows each participant to safely finalize their transfers without requiring trust.

For a mediated transfer to work a number of nodes need to collaborate. Which nodes that would be is determined by the path, detailed in the [transfer routing section]().

Let’s assume a path of Alice <-> Bob <-> Charlie. Alice does not have a direct channel with Charlie, therefore Alice can either open a new channel or mediate the transfer through other nodes. In our example Bob is a hop to whom Alice has an open channel and is considered good for routing.

The role of Bob is to mediate the transfer between Alice, the initiator, and Charlie, the target. The number of nodes in a given path may vary but the roles and guarantees work the same way.

Bob will first receive a mediated transfer `t1` from Alice. This transfer is conditionally locked with a secret generated by Alice so it’s on Alice’s hand whether to finalize the transfer or not. The transfer sent by Alice is a valid balance proof, that may be used by Bob on-chain at any time to reclaim the current received amount and if the unlocking secret is learned, to also withdraw the pending transfer `t1`. Because Bob knows that he may claim the transfer value, and Bob has properly done all the validation checks to guarantee that the balance proof contains the correct transferred amount, the merkle root effectively represents all pending transfers, and Alice does have the available balance to use with the given transfer, Bob can safely forward the transfer.

Bob will then create a new transfer `t2`, on the channel Bob <-> Charlie. `t2` has its value backed up, since Bob is not the payer, by transfer `t1`. In the given example Alice is the payer to Bob, and Charlie is the payee to Bob. Note that in the case of an increased number of hops there will be more payer/payee pairs, one for each mediator. The transfer `t2` will be another conditionally locked transfer, and the mediator is responsible to use the same lock amount and hashlock for this transfer.

Once the transfer target has received the mediated transfer it will request from the initiator the secret to unlock the transfer. At this point in time the initiator knows that some node will pay the target. The target informs the initiator about the received lock’s amount, token, and hashlock, so it can be sure that it’s the correct transfer. Now the initiator Alice is at the position of completing the transfer by revealing the secret to Charlie.

Once the secret is known by the target the payments flow from the back to the front of the payment chain. That means they start at Charlie who will request a withdrawal from Bob, informing Bob about the known secret, allowing Bob to request a withdrawal from Alice.

**Alternative Protocol Implementation**
There is nothing about the way that locks operate that forbids transfer splitting, i.e.: a mediator doesn’t have enough capacity on a channel but it can forward the transfer to two or more channels that on aggregate have the correct amount. Although that scheme is possible, it is not currently considered because of some added complexity. The target node would either need to know the transfer id and amount prior to its start, or it would need to make multiple secret requests to the initiator, as new transfers with the same hashlock arrives, until the correct transfer amount is reached.

## Locks
A lock has two parts, an amount used to track how much token is being locked, and rules to define how it may be unlocked. The lock itself is independent from the channel or token associated with it. What binds the lock to a specific channel is the balance proof’s merkle tree.

SmartRaiden currently relies on hash time locks heavily. They are the essential ingredient for safe trustless mediated transfers. This kind of lock has two additional data attributes, a hash image and a expiration. The lock is unlocked if the [preimage]() of the hash is revealed prior to its expiration. In Smartraiden the preimage is called secret and its hash is called the hashlock. The [secret]() is just 32 bytes of cryptographically secure random data. The hashlock can be the result of any cryptographically secure hash function but SmartRaiden currently relies on the keccak hash function.

With this lock construct it is possible to:
- Mediate token transfers, by relying on the same hashlock but different expiration times.
- Perform token swaps. Two mediated transfers for different tokens are made with the same hashlock and once the secret is revealed we end up having an atomic swap of the tokens. (Token swaps are not part of the Red Eyes release.)

**Alternative Protocol Implementation**

The preimage could be a hash of another structure, e.g. a written contract. This would bind the action of unlocking a lock release to a document.

The lock could require two hashlocks to unlock. This construct if used in a proper order, would allow for receipts to be generated.

The lock could require either of two hashlocks to unlock. This construct allows for safe refunds that don’t need to wait for the lock expiration.

## Safety of Mediated Transfers
The safety of mediated transfers relies on two rules:
- For a node to withdraw a lock it must reveal the secret.
- The mediator must have time to withdraw the payer’s transfer after the payee withdrew the mediator’s transfer.

The first is trivially achieved by allowing two forms of withdraw. A node may withdraw off-chain by exchanging the secret message and receiving a balance proof or on-chain by revealing the secret.

The second is the mediator’s responsibility to choose a lock expiration for the payer transfer that in the worst case would allow him enough time to withdraw. The worst case is a withdraw on-chain that requires:

- Learning about the secret from the payee withdrawal on-chain.
- Closing the channel.
- Updating the counter party transfer.
- Withdrawing the lock on the closed channel.

The number of blocks for the above is named [reveal timeout]().

**Alternative Protocol Implementation**   
The reveal timeout is large because the blockchain can be congested due to large amount of traffic. This delays the processing of closing/withdraw transactions enough that token loss is possible. At the same time it is impossible to predict how long congestion would last. Ideally the smart contract would be able compute the unlock operations that could have been executed and count lock expiration to the available [gas slots]() of the mined blocks.

## Failed Mediated Transfers
Failed mediated transfers are defined as transfers for which the initiator does not reveal the secret making it impossible to withdraw the lock. This may happen for two reasons. Either the initiator didn’t receive a [SecretRequest](), or the initiator discarded the secret to retry the transfer with a different route.

The initiator might not have received the SecretRequest for yet another set of reasons:

- Connectivity problems between the initiator and the target.
- The maximum number of hops was reached, the lock expiration cannot be further decremented so the last node is not willing to make progress.
- Some byzantine node along the path is not proceeding with the protocol.

For any of the above scenarios, each hop must hold the lock and wait until it expires before unlocking the token and letting the payer add it back to its available balance.

**Alternative Protocol Implementation**  
Use a new lock type that can be withdrawn if any of two secrets is revealed. Each mediator sends the payee transfer with a controlled refund secret. If the next hop cannot proceed with the transfer it sends back a mediated transfer using the same refund hashlock. This allows the mediator controlling the refund secret to release both locks without a risk of double spending.


## Channel Closing and Settlement
There are multiple reasons for which a channel might need to be closed:

- The partner node might be misbehaving.
- The channel owner might want to withdraw its token.
- The partner node might become unresponsive and an unlocked transfer might be at risk of expiring.

At any point in time any of the participants may close the smart contract. From the point the channel is closed and onwards transfers can not be done using the channel.

Once the channel enters the settlement window the partner state can be updated by calling `updateTransfer`. After the partner state is updated by the participant, locks may be withdrawn. A `withdraw` checks the lock and updates the partner’s current transferred amount. This is safe since a participant is allowed to provide the partner state only once and neither the transferred amount nor the locksroot will change after that call.

With third parties the process changes slightly. Since third parties are allowed to call `updateTransfer` multiple times, the transferred amount and locksroot must be reset each time a new transfer is provided and locks that have been withdrawn must be withdrawn again.

**Alternative Protocol Implementation**  

The current implementation has a local unlock, meaning that the same hashlock may be provided multiple times, once for each mediator that is closing the channel. The [sprites approach]() uses a global registry of known secrets and requires the secret to be unlocked only once. This saves the computation of the hash function for each additional withdraw.

- Nodes don’t need to close the channel to unlock, since the secret can be registered with the secret manager.
- Nodes doen’t need to care about learning the secret through the blockchain and reapplying it in their own channel.
- It really simplifies thinking about lock expiration for refunds since the expiration has a fixed lower bound.

## Transfer Routing
Routing is a hard problem and because of the lack of a global view SmartRaiden has a graph search strategy. The packet routing may be looked at as an `A*` search, using the sorted path with capacity as an heuristic to do the packet routing.

![channelsgraph](/docs/images/channelsgraph.png)

Consider the above graph where each graph node represents a SmartRaiden node, each edge an existing channel, arrows represent the direction of the transfers, solid lines the current searched space of the graph, dashed lines the rest of the path and the red line an exhausted/closed channel.

The transfer initiator is `A`, the transfer target is `G`. `A` decides locally the first hop of the path. `A`’s choice is determined by what it thinks will be the path that can complete the transfer using the shortest path and sending the transfer to the first node in that path.

`B` will mediate the transfer and do its own local path routing. It chose `C`, which in turn chose `T`. It turns out that both `B` and `C` made a suboptimal choice. `T` was not able to complete the transfer with its channel and will route the transfer through `D`. This will continue until either the transfer expires or the target is reached. Note that the transfer’s lock expiration is not the same as a protocol level TTL. This behaviour could improve if we add a TTL to protocol messages so that we can inform mediators that a transfer was discarded by the initiator and further tries will be futile.

Each of these hops forwarded a MediatedTransfer paying fees and sending the transfer value to the next hop to mediate the transfer.

**Alternative Protocol Implementation**  

Path finding services: Nodes may choose routing services to update with their current available balance, the routing services will charge a fee to the users to provide routes.

Onion encryption: To improve anonymity, encryption may be used. The initiator will choose a path that cannot be changed during the transfer and onion encrypt the hops. Garbage of a variable length must be added to the end of the onion encrypted path to hide the path length.

## Merkle Tree
![merkletree](/docs/images/merkletree.png)

The [merkle tree]() data blocks are composed of the hashes of the locks. The unique purpose of the merkle tree is to have an `O(log N)` proof of containment and a constant `O(1)` storage requirement for the signed messages. The alternative is to have linear space `O(n)` for the signed messages by having a list of all the pending locks in each message.

The merkle tree must have a deterministic order, that can be computed by any participant or the channel contract. The leaf nodes are defined to be in lexicographical order of the elements (lock hashes). For the other levels the interior nodes are also computed from the lexicographical order.

**Alternative Protocol Implementation**  

Use time order for the leaves and lexicographical for the intermediary nodes. This will greatly improve insertion performance since only the rightmost side of the tree must be recomputed. It may also improve removals since the nodes to the left don’t need to be recomputed.

## SmartRaiden Design Choices

### One Contract per Channel
At the beginning SmartRaiden was designed with simplicity in mind and the one contract per channel made the code simpler. We have plans to change to a single contract per token.

### Network Protocol Messages Must not have Inherited Trust

We don’t support informational messages like `TransferTimeout`, `TransferCancelled`, nor messages that can lead to change of the channel state without some mechanism backed by the smart contract as this would imply trust between participants and open attack vectors.

We also don’t support messages from nodes saying that something happened on the blockchain because that is both redundant and imply trust in the message.

### Invalid Messages can Happen

SmartRaiden is built on top of an asynchronous network and one may not trivially assume that things are globally ordered, so invalid messages cannot be naively assumed as an attack.

The lack of synchronization messages is a security measure as seen in the above section. As a consequence there are race conditions. For example picture a fresh channel between Alice and Bob.

- Alice deposits 10.
- Alice sends a transfer of 5 to Bob.
- Bob received the transfer, checks if Alice may spend this amount and it fails.
- Bob polls the blockchain and learns about the ChannelDeposit event.

This is fixed in the protocol layer with a retry mechanism. It assumes the partner node has a properly working ethereum client and that he is polling for events from the block.

### Hashlocks are not Transfer Identifiers

Even though hashlocks are unique, this value is not used as an identifier because routing via a specific path may fail. The initiator is at a position where he may choose to discard a path and retry with a different first hop. For this reason hashlocks can change for the same payment and an additional field just for transfer identification is used.

### Fixed charges
We can incentivize nodes to retain more balance in payment channels via a method to take a charge for them, in the end to maximize the rate of successful transactions. Hence, charging fees to nodes is highly relevant to the route choosing. Right now, SmartRaiden adopts a route algorithm that is based on the sender node, in which sender and routers individually enforce this algorithm. Sender has no idea either which route will finally connect to the receiver node, or how many gas he will be charged, so to avoid that, our sender need to deposit extra gas in the channel. To be clear, in SmartRaiden 3.0 version we make every node in the network to charge the same gas, in which case it is quite easy to calculate gases spent in each transaction. This is NOT the most efficient way for each node is able to charge differently. In future version, we plan to keep a record of charge list for each node in the server, and choose a route that only requires a minimum amount of gases to connect to the receiver node, and maximizes the benefit for all the nodes in this route.

### Internet-free Payment
Internet-free Payment : Internet-free transaction is a off-chain transaction that has no necessity for Internet to be enforced. However, if there is no connection to public chain for quite a long period, we can not ensure SmartRaiden is secure when transferring tokens. So we need somehow find a method, like App, to keep the record of period of how long nodes are off-line from public chain, before any Internet-free transaction and any node receiving a transaction message under such circumstances to check whether such a transaction is in safe state.

Let’s still take Alice & Bob as our example. When Alice proposes to give 30 tokens to Bob, they need to create a bilateral channel with the condition of Internet. Assume that Alice creates a direct payment channel to Bob in the spectrum, and deposits 100 tokens of hers, then in such channel, Alice has her total balance of 100 tokens, as to Bob, it’s 0 token. With Internet connection, actually they are able to make their transaction out of chain. Suppose it is the case, if Alice is transferring her part of deposit to Bob when there is something wrong with Internet connection, and there is a meshbox near this payment channel, and both Alice and Bob have registered on this meshbox, such as :

Alice :` { ‘address’ :  “0x1a9ec3b0b807464e6d3398a59d6b0a369bf422fa”, ‘ip_port’ : “192.168.0.5:40001”}`

Bob :`{ ‘address’ : “0x8c1b2E9e838e2Bf510eC7Ff49CC607b718Ce8401”,
‘ip_port’ : “192.168.0.7:40001”}`

then the payment channel between Alice & Bob is still effective.

Right now, Alice can invoke SwitchNetwork to alter to Internet-free state, and with resort to meshbox, Alice can give her 30 tokens directly to Bob. While Internet-free state has its limitation, that is, it does not support any indirect transaction with regard to sending, transferring, and receiving. Besides, nodes can only create direct channel to do transactions, in Internet-free state.

### Node synchronization state description
It is quite significant for transaction security to keep states of participants synchronized. At the time SmartRaiden adopts the way of state machine to maintain this synchronization among nodes. Assume that we have two participants, Alice and Bob, who want to use SmartRaiden as a system for off-chain transactions. For example, if Alice plans to transfer 30 tokens to Bob, by MediatedTransfer, then what will it be that synchronization states of Alice and Bob.

#### 1.Node synchronization when sending MediatedTransfer messages
In the event that Alice sends MediatedTransfer message to Bob, like to transfer 20 tokens, to make sure this message is the most recent one, Alice has to store data of MediatedTransfer into local storage. After Bob has received this message, the first thing to do is that, he extracts the history of records from local memory to check whether the nounce of this message is consecutive with the one before it, and to verify the balance proofs. If no faulty event, write this newest message into local storage and return a ACK to Alice. When Alice gets messages from Bob, the process of synchronization completes. If Alice goes offline before she gets any ACK from Bob, then both parties revoke all the state changes back to the one prior to Alice sending MediatedTranfer to Bob, which is atomic. After Alice goes online again, she will gets another ACK resent from Bob, then both of them starts into synchronization again. Then communication of messages continues till all the transactions complete.

#### 2. Node synchronization when receiving MediatedTransfer messages.
Alice sends MediatedTransfer messages to Bob, like to transfer 20 tokens, then Bob received them and verify with history record, after that Bob sends ACK message back to Alice. If Bob goes offline before sending ACK, operations of local record storage fail in order to secure atomicity of relevant operations by Bob, then all the states go back to states before Bob receiving any MediatedTransfer message. At the moment Alice has no certainty that Bob has already received messages sent by her, so that tokens in this transaction will lock up. When Bob is online, he will send another ACK message to Alice, and the state alters to one representing the lock-up state of tokens. After Alice gets ACK from Bob, both of their states start to synchronize. Then communication of messages goes on till all the transaction s complete.
