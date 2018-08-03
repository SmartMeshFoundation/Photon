# SmartRaiden’s API Documentation

## Introduction
SmartRaiden has a Restful API with URL endpoints corresponding to user-facing interaction allowed by a SmartRaiden node. The endpoints accept and return JSON encoded objects. The api url path always contains the api version in order to differentiate queries to different API versions. All queries start with:  `/api/<version>/`.
## JSON Object Encoding
The objects that are sent to and received from the API are JSON-encoded. Following are the common objects used in the API.
### Channel Object
```json
{
        "channel_address": "0xc4327c664D9c47230Be07436980Ea633cA3265e4",
        "partner_address": "0x31DdaC67e610c22d19E887fB1937BEE3079B56Cd",
        "balance": 200,
        "partner_balance": 100,
        "locked_amount": 0,
        "partner_locked_amount": 0,
        "token_address": "0x745D52e50cd1b19563D3a3B7B6d2eB60b17E6bAE",
        "state": "opened",
        "settle_timeout": 100,
        "reveal_timeout": 3
    }
```
A channel object consists of a:
-   `channel_address`  should be a  `string`  containing the hexadecimal address of the channel
 -   `partner_address`  should be a  `string`  containing the hexadecimal address of the partner with whom we have opened a channel
 - `balance` should be an integer of the amount of the `token_address` token we have available for transferring
 - `partner_balance` should be an integer of the amount of the `token_address` token partner have available for transferring
 - `locked_amount` should be an integer of the amount of the `token_address` token we have locked amount
 - `partner_locked_amount` should be an integer of the amount of the `token_address` token partner have locked amount
 - `token_address` should be a `string` containing the hexadecimal address of the token we are trading in the channel
 -   `state`  should be the current state of the channel represented by a string. Possible value are: -  `opened`: The channel is open and tokens are tradeable -  `closed`: The channel has been closed by a participant -  `settled`: The channel has been closed by a participant and also settled
 -  `settle_timeout`: The number of blocks that are required to be mined from the time that  `close()`  is called until the channel can be settled with a call to  `settle()`
 - `reveal_timeout`: The maximum number of blocks allowed between the setting of a hashlock and the revealing of the related secret
## Endpoints
Following are the available API endpoints with which you can interact with SmartRaiden.
### Querying Information About Your SmartRaiden Node
**`GET /api/<version>/address`**  
Query your address. When SmartRaiden starts, you choose an ethereum/Spectrum address which will also be your SmartRaiden address  
**Example Request**:  
`GET http://localhost:5001/api/1/address`  
**Example Response**:  
*`200 OK`* and 
```json
{
    "our_address": "0x69C5621db8093ee9a26cc2e253f929316E6E5b92"
}
```
### Deploying
**`PUT/api/<version>/tokens/<token_address>`**  
Registers a token. If a token is not registered yet (i.e.: A token network for that token does not exist in the registry), we need to register it by deploying a token network contract for that token.  
**Example Request**:  
`PUT http://localhost:5001/api/1/tokens/0xB0159439B496b8cebd54f232Ae06d61d0bE1Fe45`  
**Example Response**:  
*`200 OK`* and 
```json
{
    "channel_manager_address": "0x0aa88934bc3B0E9623d9555ceA48ab60FF3f2869"
}
```
Status Codes:
* `200 Created`– A token network for the token has been successfully created
* `409 Conflict` – The token was already registered before or  The registering transaction failed.

Response JSON Object:
- **channel_manager_address** Channel management contract address

### Querying Information About Channels and Tokens
**`GET/api/<version>/channels`**  
Querying all channels  
**Example Request**:  
`GET http://localhost:5004/api/1/channels`  
**Example Response**:  
*`200 OK`* and   
```json
 {
        "channel_address": "0xd5CF2248292e75531d314B118a0390132bc7a5F0",
        "partner_address": "0x088da4d932A716946B3542A10a7E84edc98F72d8",
        "balance": 100,
        "partner_balance": 100,
        "locked_amount": 0,
        "partner_locked_amount": 0,
        "token_address": "0x745D52e50cd1b19563D3a3B7B6d2eB60b17E6bAE",
        "state": "opened",
        "settle_timeout": 100,
        "reveal_timeout": 3
    },
    {
        "channel_address": "0xdF474bBc5802bFadc4A25cf46ad9a06589D5AF7D",
        "partner_address": "0x31DdaC67e610c22d19E887fB1937BEE3079B56Cd",
        "balance": 100,
        "partner_balance": 200,
        "locked_amount": 0,
        "partner_locked_amount": 0,
        "token_address": "0x745D52e50cd1b19563D3a3B7B6d2eB60b17E6bAE",
        "state": "opened",
        "settle_timeout": 100,
        "reveal_timeout": 3
    }
```
Status Codes:
* `200 OK`-Successful query

Querying a specific channel  
**Example Request**:  
`GET http://localhost:5004/api/1/channels/0xd5CF2248292e75531d314B118a0390132bc7a5F0`  
**Example Response**:  
*`200 OK`* and 
```json
{
    "channel_address": "0xd5CF2248292e75531d314B118a0390132bc7a5F0",
    "partner_address": "0x088da4d932A716946B3542A10a7E84edc98F72d8",
    "balance": 100,
    "patner_balance": 100,
    "locked_amount": 0,
    "partner_locked_amount": 0,
    "token_address": "0x745D52e50cd1b19563D3a3B7B6d2eB60b17E6bAE",
    "state": "opened",
    "settle_timeout": 100,
    "reveal_timeout": 0,
    "ClosedBlock": 0,
    "SettledBlock": 0,
    "OurUnkownSecretLocks": {},
    "OurKnownSecretLocks": {},
    "PartnerUnkownSecretLocks": {},
    "PartnerKnownSecretLocks": {},
    "OurLeaves": null,
    "PartnerLeaves": null,
    "OurBalanceProof": null,
    "PartnerBalanceProof": null,
    "Signature": null
}
```
Status Codes:
* `200 OK`-Successful query
* `404 Not Found` -If the channel does not exist

Querying all registered Tokens.Returns  a list of addresses of all registered tokens.  
**Example Request**:  
`GET http://localhost:5004/api/1/tokens`  
**Example Response**:  
*`200 OK`* and 
```json
[
    "0xb0159439b496b8cebd54f232ae06d61d0be1fe45",
    "0x541eefe890a10d27d947190ea976cb6dcbba650f",
    "0xf3db2928689cdbd9938d1e1ffc2c4980a96f299e",
    "0x745d52e50cd1b19563d3a3b7b6d2eb60b17e6bae"
]
```
Status Codes:
* `200 OK`-Successful query
* `404 Not Found` -If the token does not exist


 Querying all partners for a Token,Returns a list of all partners with whom you have non-settled channels for a certain token.  
 **Example Request**:  
 `GET http://localhost:5004/api/1/tokens/0x745D52e50cd1b19563D3a3B7B6d2eB60b17E6bAE/partners`  
 **Example Response**:      
*`200 OK`* and  
```json
[
    {
        "partner_address": "0x088da4d932A716946B3542A10a7E84edc98F72d8",
        "channel": "api/1/channles/0xd5CF2248292e75531d314B118a0390132bc7a5F0"
    },
    {
        "partner_address": "0x31DdaC67e610c22d19E887fB1937BEE3079B56Cd",
        "channel": "api/1/channles/0xdF474bBc5802bFadc4A25cf46ad9a06589D5AF7D"
    }
]
```
Status Codes:
* `200 OK`-Successful query
* `404 Not Found` -If the token does not exist

Response JSON Array of Objects:
-   **partner_address**  (_address_) – The partner we have a channel with
-   **channel**  (_link_) – A link to the channel resource
  
Token Swaps  

**`PUT /api/<version>/token_swaps/<target_address>/<identifier>`**

You can perform a token swap by using the  `token_swaps`  endpoint. A swap consists of two users agreeing on atomically exchanging two different tokens at a particular exchange rate.
tips：
* The parties involved in Swaps have an effective channel
* Call *taker* first and then call *maker*

**Example Request**:  
*the taker:*  `PUT http://localhost:5001/api/1/token_swaps/0x31DdaC67e610c22d19E887fB1937BEE3079B56Cd/3333`  
with payload:
```json
{
    "role": "taker",
    "sending_amount": 50,
    "sending_token": "0x745d52e50cd1b19563d3a3b7b6d2eb60b17e6bae",
    "receiving_amount": 5,
    "receiving_token": "0x541eefe890a10d27d947190ea976cb6dcbba650f"
}
```
*the maker:*
`PUT http:// localhost:5002/api/1/token_swaps/0x69C5621db8093ee9a26cc2e253f929316E6E5b92/3333`  
with payload:
```json
{
    "role": "maker",
    "sending_amount": 5,
    "sending_token": "0x541eefe890a10d27d947190ea976cb6dcbba650f",
    "receiving_amount": 50,
    "receiving_token": "0x745d52e50cd1b19563d3a3b7b6d2eb60b17e6bae"
}
```
Status Codes:
* `201 Created`-Successful query
* `400  Bad Request` -no available route

### Channel Management
**`PUT/api/<version>/channels`**  
Opens channel.  
**Example Request**:  
`PUT http:// localhost:5001/api/1/channels`  
with payload:
```json
{
    "partner_address": "0x31DdaC67e610c22d19E887fB1937BEE3079B56Cd",
    "token_address": "0x541eefe890a10d27d947190ea976cb6dcbba650f",
    "balance": 200,
    "settle_timeout": 100
}
```
The  `balance`  field will signify the initial deposit you wish to make to the channel.

The request to the endpoint should later return the fully created channel object from which we can find the address of the channel.
**Example Response**:  
*`200 OK`* and 
```json
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
Status Codes:
* `200 OK`-Channel created successfully
* `404 Not Found` -If the token does not exist
* `409  Conflict` -NewChannel tx execution failed

**`PATCH/api/<version>/channels/<channel_address>`**  
 Close Channel  
**Example Request**:  
`PATCH http:// localhost:5001/api/1/channels/0xD955A1BA24058BFbFfD98dF78253a861e5B029b9`  
with payload:
```json
{"state":"closed"}
```
**Example Response**:  
*`200 OK`* and 
```json
{
    "channel_address": "0xD955A1BA24058BFbFfD98dF78253a861e5B029b9",
    "partner_address": "0x1DdaC67E610c22d19e887FB1937bEe3079B56CD1",
    "balance": 200,
    "partner_balance": 0,
    "locked_amount": 0,
    "partner_locked_amount": 0,
    "token_address": "0x541eeFe890A10D27d947190EA976CB6DCBba650f",
    "state": "closed",
    "settle_timeout": 100,
    "reveal_timeout": 0
}
```
Settle Channel  
**Example Request**:  
`PATCH http:// localhost:5001/api/1/channels/0xD955A1BA24058BFbFfD98dF78253a861e5B029b9`  
with payload:
```json
{"state":"settled"}
```
**Example Response**:  
*`200 OK`* and 
```json
"channel_address": "0xD955A1BA24058BFbFfD98dF78253a861e5B029b9",
    "partner_address": "0x1DdaC67E610c22d19e887FB1937bEe3079B56CD1",
    "balance": 200,
    "partner_balance": 0,
    "locked_amount": 0,
    "partner_locked_amount": 0,
    "token_address": "0x541eeFe890A10D27d947190EA976CB6DCBba650f",
    "state": "settled",
    "settle_timeout": 100,
    "reveal_timeout": 0
```
**`PATCH  /api/<version>/channels/<channel_address>`**  
 Deposit to a Channel    
 You can deposit more of a particular token to a channel by updating the `balance` field of the channel in the corresponding endpoint with a `PATCH` http request.  
 **Example Request**:  
 `PATCH http://localhost:5002/api/1/channels/0x7f9bc53F7b3e08a3A9De564740f7FAf9Decb16B9`  
 with payload:
```json
{
    "balance": 100
}
```
**Example Response**:  
*`200 OK`* and 
```json
{
    "channel_address": "0x7f9bc53F7b3e08a3A9De564740f7FAf9Decb16B9",
    "partner_address": "0x69C5621db8093ee9a26cc2e253f929316E6E5b92",
    "balance": 100,
    "partner_balance": 200,
    "locked_amount": 0,
    "partner_locked_amount": 0,
    "token_address": "0x541eeFe890A10D27d947190EA976CB6DCBba650f",
    "state": "opened",
    "settle_timeout": 100,
    "reveal_timeout": 0
}
```
Status Codes:
* `200 OK`-For successful Deposit
* `400 Bad Request` -If the provided json is in some way malformed
### Connection Management

**`GET  /api/<version>/connections`**  
 Querying connections details  
You can query for details of previously joined token networks by making a GET request to the connection endpoint.  
 **Example Request**:  
 `GET http:// localhost:5003/api/1/connections`  
 **Example Response**:  
*`200 OK`* and 
```json
{
    "0x745D52e50cd1b19563D3a3B7B6d2eB60b17E6bAE": {
        "funds": 0,
        "sum_deposits": 200,
        "channels": 1
    }
}
```
**Example Response**:  
*`200 OK`* and 
```json
{
    "0x745D52e50cd1b19563D3a3B7B6d2eB60b17E6bAE": {
        "funds": 0,
        "sum_deposits": 280,
        "channels": 3
    }
}
```
Response JSON Array of Objects:
-   **funds**  (_int_) – Funds from last connect request
-   **sum_deposits**  (_int_) – Sum of deposits of all currently open channels
-   **channels**  (_int_) – Number of channels currently open for that token 

Status Codes:
* `200 OK`-For a successful query

**`PUT  /api/<version>/connections/<token_address>`**  
Automatically join a token network. The request will only return once all blockchain calls for opening and/or depositing to a channel have completed.  
 **Example Request**:  
 `PUT http://localhost:5001/api/1/connections/0xf1b0964f1e19ecf07ddd3bd8e20138c82680395d`  
 **Example Response**:  
*`201 Created`*   
Status Codes:
* `201 Created`-For a successful connection creation
* `500  Internal Server Error`-Internal SmartRaiden node error

**`DELETE  /api/<version>/connections/<token_address>`**  
The request will only return once all blockchain calls for closing/settling a channel have completed.  

Important note. If no arguments are given then SmartRaiden will only close and settle channels where your node has received transfers. This is safe from an accounting point of view since deposits can’t be lost and provides for the fastest and cheapest way to leave a token network when you want to shut down your node.

If the default behaviour is not desired and the goal is to leave all channels irrespective of having received transfers or not then you should provide as payload to the request  `only_receiving_channels=false`

A list with the addresses of all the closed channels will be returned.  
 **Example Request**:  
 `DELETE http://localhost:5003/api/1/connections/0x541eefe890a10d27d947190ea976cb6dcbba650f`  
  with payload:
 ```js
 {
  "only_receiving_channels":false
}
```
 **Example Response**:  
*`200 OK`* and 
```js
[
    "0x08Bb272f51c8974ACe71648d01afE933384A762e",
    "0x68f9390554789c2D658540C1d3A450fb858a849e",
    "0xc0cd666a125F9bbf4dAeFcbd24F2dc26f7BC9f8D"
]
```
The response is a list with the addresses of all closed channels.

Request JSON Object:

-   **only_receiving_channels**  (_boolean_) – Only close and settle channels where your node has received transfers. Defaults to  `true`.

Status Codes:
- `200 OK`-For successfully leaving a token network
- `500  Internal Server Error`-Internal SmartRaiden node error
### Transfers
**`POST  /api/<version>/transfers/<token_address>/<target_address>`**

 Initiating a Transfer
 You can create a new transfer by making a  `POST`  request to the following endpoint along with a json payload containing the transfer details such as amount and identifier. Identifier is optional.
 
The request will only return once the transfer either succeeded or failed. A transfer can fail due to the expiration of a lock, the target being offline, channels on the path to the target not having enough `settle_timeout` and `reveal_timeout` in order to allow the transfer to be propagated safely e.t.c  
 **Example Request**:  
 `POST http://localhost:5002/api/1/transfers/0x745D52e50cd1b19563D3a3B7B6d2eB60b17E6bAE/0x69C5621db8093ee9a26cc2e253f929316E6E5b92`  
with payload:
```json
{
    "amount":10,
    "fee":0,
    "is_direct":false
}
```
 **Example Response**:  
*`200 OK`* and 
```json
{
    "initiator_address": "0x31DdaC67e610c22d19E887fB1937BEE3079B56Cd",
    "target_address": "0x69C5621db8093ee9a26cc2e253f929316E6E5b92",
    "token_address": "0x745D52e50cd1b19563D3a3B7B6d2eB60b17E6bAE",
    "amount": 10,
    "identifier": 5018140839335492878,
    "fee": 0,
    "is_direct": false
}
```
Request JSON Object:

-   **amount**  (_int_) – Amount to be transferred
-   **fee**  (_int_) –  incentivize nodes to retain more balance in payment channels via a method to take a charge for them(default:0)
- **is_direct"**(_boolean_)–  If it is set to true, it can only satisfy the two parties who have direct access to the transaction. If the two sides do not have direct access, they will give up the transaction.

Status Codes:
- `200 OK` – Successful transfer
- `409 Conflict`– If the address or the amount is invalid or if there is no path to the target
-  `500  Internal Server Error`-Internal SmartRaiden node error
### Querying Events

Events are kept by the node. Once an event endpoint is queried the relevant events from either the beginning of time or the given block are returned.

Events are queried by two different endpoints depending on whether they are related to a specific channel or not.

All events can be filtered down by providing the query string argument  `from_block`  to signify the block from which you would like the events to be returned.  
**`GET  /api/<version>/events/network`**  
Query for registry network events.  
 **Example Request**:  
 `GET http://localhost:5001/api/1/events/network`  
 **Example Response**:  
*`200 OK`* and 
```json
[
    {
        "event_type": "TokenAdded",
        "token_address": "0x745D52e50cd1b19563D3a3B7B6d2eB60b17E6bAE",
        "channel_manager_address": "0x48fA4f2230DB0dEEA3989014CD21857DF6210B33"
    },
    {
        "event_type": "TokenAdded",
        "token_address": "0x541eeFe890A10D27d947190EA976CB6DCBba650f",
        "channel_manager_address": "0xdC47DF3eAc0E9a8373A258ccf4838bE3540a9D4E"
    },
    {
        "event_type": "TokenAdded",
        "token_address": "0xF3DB2928689Cdbd9938d1e1Ffc2c4980a96f299E",
        "channel_manager_address": "0x78494a9F7278F5eE5a3faB445685EABA7add547a"
    },
    {
        "event_type": "TokenAdded",
        "token_address": "0xB0159439B496b8cebd54f232Ae06d61d0bE1Fe45",
        "channel_manager_address": "0x0aa88934bc3B0E9623d9555ceA48ab60FF3f2869"
    }
]
```
Status Codes:
- `200 OK` – For successful Query
- `404  Not Found`–If the provided query string is malformed

**`GET  /api/<version>/events/tokens/<token_address>`**  
Querying token network events  
 **Example Request**:  
 `GET http://localhost:5001/api/1/events/tokens/0x745D52e50cd1b19563D3a3B7B6d2eB60b17E6bAE`  
**Example Response**:  
*`200 OK`* and   
```json
[
    {
        "event_type": "ChannelNew",
        "settle_timeout": 40,
        "netting_channel": "0x5629954B107E1889516E0CEC046432aa20f70778",
        "participant1": "0x69C5621db8093ee9a26cc2e253f929316E6E5b92",
        "participant2": "0x31DdaC67e610c22d19E887fB1937BEE3079B56Cd",
        "token_address": "0x745D52e50cd1b19563D3a3B7B6d2eB60b17E6bAE"
    },
    {
        "event_type": "ChannelNew",
        "settle_timeout": 100,
        "netting_channel": "0x8898917d1d2DF53595DA560f5b162BC7e6BCBDa0",
        "participant1": "0xf0f6E53d6bbB9Debf35Da6531eC9f1141cd549d5",
        "participant2": "0x62e68b5c745Fa25B439BF613dC1Fcd262277fa93",
        "token_address": "0x745D52e50cd1b1
        9563D3a3B7B6d2eB60b17E6bAE"
    }
]
```
Status Codes:
- `200 OK` – For successful Query
- `404  Not Found`–If the provided query string is malformed

**`GET  /api/<version>/events/channels/<channel_registry_address>`**  
 Querying channel events  
  **Example Request**:  
  `GET http://localhost:5002/api/1/events/channels/0xd1102D7a78B6f92De1ed3C7a182788DA3a630DDA`   
  **Example Response**:  
*`200 OK`* and   
```json
[
    {
        "balance": 100,
        "block_number": 2469154,
        "event_type": "ChannelNewBalance",
        "participant": "0x31DdaC67e610c22d19E887fB1937BEE3079B56Cd"
    },
    {
        "balance": 100,
        "block_number": 2727927,
        "event_type": "ChannelNewBalance",
        "participant": "0xf0f6E53d6bbB9Debf35Da6531eC9f1141cd549d5"
    },
    {
        "amount": 10,
        "block_number": 3417198,
        "event_type": "EventTransferSentSuccess",
        "identifier": 5018140839335492878,
        "target": "0x69c5621db8093ee9a26cc2e253f929316e6e5b92"
    }
]
```
Status Codes:
- `200 OK` – For successful Query
- `400  Bad Request`–If the channel does not exist

