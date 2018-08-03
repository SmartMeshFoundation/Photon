# Getting started with the SmartRaiden API
## Introduction
SmartRaiden has a Restful API with URL endpoints corresponding to actions that users can perform with their channels. The endpoints accept and return JSON encoded objects. The API URL path always contains the API version in order to differentiate queries to different API versions. All queries start with: /api/<version>/ where <version> is an integer representing the current API version.

This section will walk through the steps necessary to participate in a SmartRaiden Token Network. Some different scenarios such as joining an already existing token network, registering a new token network, together with opening, closing and settling channels, will be provided.

Before getting started with below guides, please see [Overview and Guide](https://github.com/SmartMeshFoundation/SmartRaiden/blob/master/docs/overview.md), to make sure that a proper connection to SmartRaiden is established.

Furthermore, to see all available endpoints, please see [REST API Endpoints](https://github.com/SmartMeshFoundation/SmartRaiden/blob/master/docs/rest_api.md).
## Scenarios
Below is a series of different scenarios showing different ways a user can interact with the SmartRaiden API.

A good way to check that SmartRaiden was started correctly before proceeding is to check that the SmartRaiden address is the same address as the Ethereum address chosen, when starting the SmartRaiden node:
```
GET /api/1/address
```
If this returns the same address, we know that the SmartRaiden node is up and running correctly.

## Bootstrapping a token network
In this scenario it is assumed that a user holds some ERC20 token, with address `0x745D52e50cd1b19563D3a3B7B6d2eB60b17E6bAE`, which has not yet been registered with SmartRaiden.

The user wants to register the token, which will create a [Channel Manager](https://github.com/SmartMeshFoundation/SmartRaiden/blob/master/network/rpc/contracts/ChannelManagerContract.sol). For each registered token there is a corresponding channel manager. Channel managers are responsible for opening new payment channels between two parties.

### Checking if a token is already registered
One way of checking if a token is already registered is to get the list of all registered tokens and check if the address of the token wanted for interaction exists in the list:

```
GET /api/1/tokens
```
If the address of the token exists in the list, see the next scenario. If it does not exist in the list, it is desired to register the token.

### Registering a token
In order to register a token only its address is needed. When a new token is registered a Channel Manager contract is deployed, which makes it quite an expensive thing to do in terms of gas usage (costs about 1.8 million gas).

To register a token simply use the endpoint listed below:
```
PUT /api/1/tokens/0x745D52e50cd1b19563D3a3B7B6d2eB60b17E6bAE

```

If successful this call will return the address of the freshly created Channel Manager like this:

```
{"channel_manager_address": "0xC4F8393fb7971E8B299bC1b302F85BfFB3a1275a"}
```
The token is now registered. However, since the token was just registered, there will be no other SmartRaiden nodes connected to the token network and hence no nodes to connect to. This means that the network for this specific token needs to be bootstrapped. If the address of some other SmartRaiden node that holds some of the tokens is known or it’s simply desired to transfer some tokens to another SmartRaiden node in a one-way-channel, it can be done by simply opening a channel with this node. The way to open a channel with another SmartRaiden node is the same whether the partner already holds some tokens or not.

### Opening a channel
To open a channel with another SmartRaiden node four things are needed: the address of the token, the address of the partner node, the amount of tokens desired for deposit, and the settlement timeout period. With these things ready a channel can be opened:

```
PUT /api/1/channels
```
With the payload

```
{
    "partner_address": "0x31DdaC67e610c22d19E887fB1937BEE3079B56Cd",
    "token_address": "0x541eefe890a10d27d947190ea976cb6dcbba650f",
    "balance": 200,
    "settle_timeout": 100
}

```
At this point the specific value of the balance field isn’t too important, since it’s always possible to deposit more tokens to a channel if need be.

Successfully opening a channel will return the following information:
```
{
    "channel_address": "0x7f9bc53F7b3e08a3A9De564740f7FAf9Decb16B9",
    "partner_address": "0x31DdaC67e610c22d19E887fB1937BEE3079B56Cd",
    "balance": 200,
    "partner_balance": 0,
    "locked_amount": 0,
    "partner_locked_amount": 0,
    "token_address": "0x541eeFe890A10D27d947190EA976CB6DCBba650f",
    "state": "opened",
    "settle_timeout": 100,
    "reveal_timeout": 0
}
```
Here it’s interesting to notice that a channel_address has been generated. This means that a [Netting Channel contract](https://github.com/SmartMeshFoundation/SmartRaiden/blob/v0.3/network/rpc/contracts/NettingChannelContract.sol) has been deployed to the blockchain. Furthermore it also represents the address of the payment channel between two parties for a specific token.

### Depositing to a channel

A payment channel is now open between the user’s node and a counterparty with the address `0x61c808d82a3ac53231750dadc13c777b59310bd9`. However, since only one of the nodes has deposited to the channel, only that node can make transfers at this point in time. Now would be a good time to notify the counterparty that a channel has been opened with it, so that it can also deposit to the channel. All the counterparty needs in order to do this is the address of the payment channel:
```
PATCH /api/1/channels/0x2a65aca4d5fc5b5c859090a6c34d164135398226
``` 
With the payload
```
{
    "balance": 7331
}
```
To see if and when the counterparty deposited tokens, the channel can be queried for the corresponding events. The `from_block` parameter in the request represents the block number to query from:

```
GET /api/1/events/channels/0x2a65aca4d5fc5b5c859090a6c34d164135398226?from_block=1337
```
This will return a list of events that has happened in the specific payment channel. The relevant event in this case is:
```
{
    "event_type": "ChannelNewBalance",
    "participant": "0x61c808d82a3ac53231750dadc13c777b59310bd9",
    "balance": 7331,
    "block_number": 54388
}
```
From above event it can be deducted that the counterparty deposited to the channel. It is possible for both parties to query the state of the specific payment channel by calling:

```
GET /api/1/channels/0x2a65aca4d5fc5b5c859090a6c34d164135398226
```
This will give us result similar to those in Opening a Channel that represents the current state of the payment channel.

A new token resulting in a new token network has now been registered. A channel between two SmartRaiden nodes has been opened, and both nodes have deposited to the channel. From here on the two nodes can start transferring tokens between each other.

The above is not how a user would normally join an already existing token network. It is only included here to show how it works under the hood.

In the next scenario it will be explained how to join already bootstrapped token networks.

## Joining an already existing token network

In above scenario it was shown how to bootstrap a token network for an unregistered token. In this section the most common way of joining a token network will be explained. In most cases users don’t want to create a new token network, but they want to join an already existing token network for an ERC20 token that they already hold.


The main focus of this section will be the usage of the `connect` and the `leave` endpoints. The `connect` endpoint allows users to automatically connect to a token network and open channels with other nodes. Furthermore the `leave` endpoint allows users to leave a token network by automatically closing and settling all of their open channels.

It’s assumed that a user holds 2000 of some awesome ERC20 token (AET). The user knows that a SmartRaiden based token network already exists for this token.

### Connect
Connecting to an already existing token network is quite simple. All that is needed, is as mentioned above, the address of the token network to join and the amount of the corresponding token that the user is willing to deposit in channels:

```
PUT /api/1/connections/0x68d94665787c85016c7db16e0be2eae78e2a1032
```
With the payload
```
{
    "funds": 3000
}
```
This will automatically connect to and open channels with three random peers in the token network, with 20% of the funds deposited to each channel. Furthermore it will leave 40% of the funds initially unassigned. This will allow new nodes joining the network to open bi-directionally funded payment channels with this node in the same way that it just opened channels with random nodes already in the network. The default behaviour of opening three channels and leaving 40% of the tokens for new nodes to connect with, can be changed by adding "initial_channel_target": 3 and "joinable_funds_target": 0.4 to the payload and adjusting the default value.

The user node is now connected to the token network for the AET token, and should have a path to all other nodes that have joined this token network, so that it can transfer tokens to all nodes participating in this network. See the Transferring tokens section for instructions on how to transfer tokens to other nodes.

### Leave
If at some point it is desired to leave the token network, the leave endpoint is available. This endpoint will take care of closing and settling all open channels for a specific in the token network:
```
DELETE /api/1/connections/0x68d94665787c85016c7db16e0be2eae78e2a1032
```


This call will take some time to finalize, due to the nature of the way that settlement of payment channels work. For instance there is a `settlement_timeout `period after calling close that needs to expire before settle can be called.

For reasons of speed and financial efficiency the leave call will only `close` and `settle` channels for which the node has received a transfer.

To override the default behaviour and `leave` all open channels add the following payload:

```
{
    "only_receiving_channels": false
}
```

## Transferring tokens
For the token transfer example it is assumed a node is connected to the token network of the AET token mentioned above. In this case the node is connected to five peers, since the standard connect() parameters were used.

### Transfer
Transferring tokens to another node is quite easy. The address of the token desired for transfer is known `0xc9d55C7bbd80C0c2AEd865e9CA13D015096ce671`. All that then remains is the address of the target node. Assume the address of the transfer node is `0x61c808d82a3ac53231750dadc13c777b59310bd9`:
```
POST /api/1/transfers/0xc9d55C7bbd80C0c2AEd865e9CA13D015096ce671/0x61c808d82a3ac53231750dadc13c777b59310bd9
```
The amount of the transfer is specified in the payload:

```
{
    "amount": 42
}
```
An "identifier": some_integer can also be added to the payload, but it’s optional. The purpose of the identifier is solely for the benefit of the Dapps built on top of SmartRaiden in order to provide a way to tag transfers.

If there is a path in the network with enough capacity and the address sending the transfer holds enough tokens to transfer the amount in the payload, the transfer will go through. The receiving node should then be able to see incoming transfers by querying all its open channels. This is done by doing the following for all addresses of open channels:

```
GET /api/1/events/channels/0x000397DFD32aFAAE870E6b5FB44154FD43e43224?from_block=1337

```
Which will return a list of events. All that then needs to be done is to filter for incoming transfers.

Please note that one of the most powerful features of SmartRaiden is that users can send transfers to anyone connected to the network as long as there is a path to them with enough capacity, and not just to the nodes that a user is directly connected to. This is called `mediated transfers`.

### Close
If at any point in time it is desired to close a specific channel it can be done with the close endpoint:

```
PATCH /api/1/channels/0x000397DFD32aFAAE870E6b5FB44154FD43e43224
```
with the payload:
```
{
    "state":"closed"
}

```
When successful this will give a response with a channel object where the state is set to `"closed"`:
```
{
    "channel_address": "0x000397DFD32aFAAE870E6b5FB44154FD43e43224",
    "partner_address": "0x61c808d82a3ac53231750dadc13c777b59310bd9",
    "token_address": "0xc9d55C7bbd80C0c2AEd865e9CA13D015096ce671",
    "balance": 350,
    "state": "closed",
    "settle_timeout": 600
}
```
Notice how the `state `is now set to `"closed"` compared to the previous channel objects where it was `"opened"`.
### Settle
Once `close` has been called, the settle timeout period starts. During this period the counterparty of the node who closed the channel has to provide its last received message. When the settlement timeout period is over, the channel can finally be settled by doing:

```
PATCH /api/1/channels/0x000397DFD32aFAAE870E6b5FB44154FD43e43224

```
with the payload:
```
{
    "state":"settled"
}
```
this will trigger the `settle()` function in the [Netting Channel contract](https://github.com/SmartMeshFoundation/SmartRaiden/blob/v0.3/network/rpc/contracts/NettingChannelContract.sol). Once settlement is successful a channel object will be returned:


```
{
    "channel_address": "0x000397DFD32aFAAE870E6b5FB44154FD43e43224",
    "partner_address": "0x61c808d82a3ac53231750dadc13c777b59310bd9",
    "token_address": "0xc9d55C7bbd80C0c2AEd865e9CA13D015096ce671",
    "balance": 0,
    "state": "settled",
    "settle_timeout": 600
}
```
Here it’s interesting to notice that the balance of the channel is now `0` and that the state is set to `"settled"`. This means that the netted balances that the two parties participating in the channel owe each other has now been transferred on the blockchain and that the life cycle of the payment channel has ended. At this point the blockchain contract has also self-destructed.

## Token Swaps
Something that has not yet been mentioned in this guide is the functionality of token swaps. A token swap allows Alice and Bob to exchange tokenA for tokenB. This means that if both Alice and Bob participate in the token networks for `tokenA` and `tokenB`, then they’re able to atomically swap some amount of tokenA for some amount of tokenB. Let’s say Alice wants to buy 5 `tokenB` for 50 `tokenA`. If Bob agrees to these terms a swap can be carried out using the `token_swaps` endpoint. In the case of the example above, Alice would be the `maker` and Bob would be the `taker`:
```
PUT /api/1/token_swaps/0x61c808d82a3ac53231750dadc13c777b59310bd9/1337
```
Where the first part after `token_swaps` is the address of Bob and the second part is an identifier for the token swap. Furthermore the following payload is needed:
```
{
    "role": "maker",
    "sending_amount": 50,
    "sending_token": "0xea674fdde714fd979de3edf0f56aa9716b898ec8",
    "receiving_amount": 5,
    "receiving_token": "0x2a65aca4d5f5b5c859090a6c34d164135398226"
}
```

There are some interesting parameters to note here. The role defines whether the address sending the message is the maker or the taker. The `taker` call must be carried out before the `maker` call can be carried out. In our design,`taker` just reigster the token swap info of which node will accept the swap, the `maker`really implement the token swap. The `sending_amount` and the `sending_token` represent the token for which the maker wants to send some amount in return for a `receiving_token` and a `receiving_amount`. So Alice is making an offer of 50 of tokenA with the address `0xea674fdde714fd979de3edf0f56aa9716b898ec8` for 5 of tokenB with the address `0x2a65aca4d5fc5b5c859090a6c34d164135398226`.
The `taker` is someone to take the offer. It could be that Alice and Bob have decided on the swap in private and thus Alice simply tells Bob the identifier. Bob can take the offer by using the same endpoint as above, but with some changes:
```
PUT /api/1/token_swaps/0xbbc5ee8be95683983df67260b0ab033c237bde60/1337
```
Here the address is the address of Alice and note that the identifier is the same as that Alice used to enforce the swap. As above, a payload is needed:
```
{
    "role": "taker",
    "sending_amount": 5,
    "sending_token": "0x2a65aca4d5fc5b5c859090a6c34d164135398226",
    "receiving_amount": 50,
    "receiving_token": "0xea674fdde714fd979de3edf0f56aa9716b898ec8"
}
```
Note that the role is changed from `maker` to `taker` . Furthermore the sending and receiving parameters have been reversed. This is because the swap is now seen from Bob’s perspective.

At this point Alice’s and Bob’s balances should reflect the state after the swap.




