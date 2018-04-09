# Raiden-Network 
   The Raiden Network is an off-chain scaling solution, enabling near-instant, low-fee and scalable payments. Our Raiden Network is a standard compliant implementation of the Raiden Network protocol in Golang. In this project, optimized raiden network will be available on multiple platforms, and decentralized micropayment on smart mobile devices can be realized.

## Project Status 
   This implementation is still very much a work in progress. It can be used for testing, but it should not be used for real funds. We do our best to identify and fix problems, and implement missing features. Any help testing the implementation, reporting bugs, or helping with outstanding issues is very welcome.
## Getting Started
   Raiden Network currently can works on Windows, Linux ,Android, iOS  etc, and requires a locally running "Raiden node" that is fully caught up with the network you're testing on.

## Build
  go get github.com/SmartMeshFoundation/raiden-mobile/tree/master/cmd/goraiden
Please refer to Installation section in raiden-network official docs the Usage in this doc for detailed instructions.
## Listing API
  Our Raiden API has the main functions of Raiden Network API, which includes：Querying Information About Channels and Tokens(Querying a specific channel,Querying all channels,Querying all registered Tokens),Registering a token,Channel Management(Open Channel,Close Channel, Settle Channel, Deposit to a Channel),Connection Management(Connecting to a token network, Leaving a token network）,Transfers（Initiating a Transfer）and Querying Events , Detailed function description please refer to Getting started with the Raiden API and Raiden’s API Documentation  in raiden-network official docs. Besides, Our Raiden also has some special function, that is, P2P communication under ICE framework, crash recovery and channel charging function. The primary API list and description are as follows：
* QueryingNodeAddress                    Query a node address
* QueryingNodeAllChannels                Query all channels of a node
* QueryingNodeSpecificChannel            Query a specified channel for a node
* QueryingRegisteredTokens               Query system registration Token
* QueryingAllPartnersForOneTokens        Query the Partner address in the channel of special Token
* RegisteringOneToken                    Registration of new Token to Raiden Network
* TokenSwaps                             Exchange the token
* OpenChannel                            Establish a channel
* CloseChannel                           Close the specified channel for the node
* SettleChannel                          Settle the specified channel for the node
* Deposit2Channel                        Deposit to the specified channel
* Connecting2TokenNetwork                Connect to a TokenNetwork
* LeavingTokenNetwork                    Leave the TokenNetwork
* QueryingConnectionsDetails             Query the details of the Token network connection 
* QueryingGeneralNetworkEvents           Query network events
* QueryingTokenNetworkEvents             Query Token network events
* QueryingChannelEvents                  Query channel event
## Raiden Contract
These are the currently deployed raiden contract addresses for the spectrum testnet:
* Netting Channel Library: [0xad5cb8fa8813f3106f3ab216176b6457ab08eb75](https://ropsten.etherscan.io/address/0xad5cb8fa8813f3106f3ab216176b6457ab08eb75#code)
* Channel Manager Library: [0xdb3a4dbae2b761ed2751f867ce197c531911382a](https://ropsten.etherscan.io/address/0xdb3a4dbae2b761ed2751f867ce197c531911382a#code)
* Registry Contract: [0x68e1b6ed7d2670e2211a585d68acfa8b60ccb828](https://ropsten.etherscan.io/address/0x68e1b6ed7d2670e2211a585d68acfa8b60ccb828#code)
* Discovery Contract: [0x1e3941d8c05fffa7466216480209240cc26ea577](https://ropsten.etherscan.io/address/0x1e3941d8c05fffa7466216480209240cc26ea577#code)

## Usage

```                                                                                                                                                    
 goraiden [global options] command [command options] [arguments...]
VERSION: 0.1
 COMMANDS:
  help, h                                                   Shows a list of commands or help for one command
 GLOBAL OPTIONS:
 --address value                                            The ethereum address you would like raiden to use and for
                                                            which a keystore file exists in your local system.
 --keystore-path                                           "/Users/your name/Library/Ethereum/keystore"  If you have a 
                                                            non-standard path for the ethereum keystore directory 
                                                            provide it using this argument.
--eth-rpc-endpoint value                                   "host:port" address of ethereum JSON-RPC server.\n'
                                                           'Also accepts a protocol prefix (ws:// or ipc channel)
                                                            with optional port', (default: 
                                                            "/Users/your name/Library/Ethereum/geth.ipc")
--registry-contract-address value                           hex encoded address of the registry contract. (default:
                                                            "0x1BB1437d4e387Be1E8C04762536217B3240f2323")
--discovery-contract-address value                          hex encoded address of the discovery contract. (default: 
                                                            "0x95A4e1251B87DCEf6B0cD18D3356CdA8cFB8f6CC")
--listen-address value                                     "host:port" for the raiden service to listen on. (default:
                                                            "0.0.0.0:40001")
--rpccorsdomain value                                        Comma separated list of domains to accept cross origin
                                                             requests. (localhost enabled by default) (default: 
                                                            "http://localhost:* /*")
--logging value                                              ethereum.slogging config-string{trace,debug,info,warn,
                                                             error,critical  (default: "trace")
--logfile value                                              file path for logging to file
--max-unresponsive-time value                                Max time in seconds for which an address can send no
                                                             packets and still be considered healthy. (default: 120)
--send-ping-time value                                       Time in seconds after which if we have received no 
                                                             message from a node we have a connection with, we are 
                                                             going to send a PING message (default: 60)
--rpc                                                        Start with or without the RPC server. Default is to
                                                              start the RPC server
--api-address value                                          host:port" for the RPC server to listen on. (default: 
                                                             "0.0.0.0:5001")
--datadir "/Users/bai/Library/GoRaiden"                       Directory for storing raiden data.
--password-file value                                         Text file containing password for provided account
--nat value                                                   [auto|upnp|stun|none] Manually specify method to use 
                                                              for determining public IP / NAT traversal.
                                                             "auto" - Try UPnP, then STUN, fallback to none "upnp"
                                                             - Try UPnP,fallback to none "stun" - Try STUN, fallback 
                                                                to none
                                                              "none" - Use the local interface
                                                              address (this will likely cause connectivity issues)
                                                              "ice"- Use ice framework for nat punching
                                                              [default: ice] (default: "ice")
--debug                                                       enable debug feature
--conditionquit value                                        quit at specified point for test
--turn-server value                                          turn server for ice (default: "182.254.155.208:3478")
--turn-user value                                            turn username for turn server 
                                                             (default: "bai")  
--turn-pass value                                            turn password for turn server 
                                                             (default: "bai")
--nonetwork                                                  disable network, for example ,when we 
                                                             want to settle all channels
--fee                                                        enable mediation fee
--help, -h                                                   how help
--version,-v                                                 print the version
                                                                                                                                                                                                                                     
```

## Requirements

geth >=1.7.3