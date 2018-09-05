# SmartRaiden 
  SmartRaiden is a standard compliant implementation of the [Raiden Network](http://raiden-network.readthedocs.io/en/stable/index.html) protocol in Golang. In this project, optimized raiden network will be available on multiple platforms, and decentralized micropayments on smart mobile devices can be realized.SmartRaiden currently can works on Windows, Linux ,Android, iOS etc. In order to better fit for the mobile network, SmartRaiden also has some special functions, the internet message communication based on XMPP , crash recovery and channel charging function.
## Project Status 
   This project is still very much a work in progress. It can be used for testing, but it should not be used for real funds. We are doing our best to identify and fix problems, and implement missing features. Any help testing the implementation, reporting bugs, or helping with outstanding issues is very welcome.
## Build
```
  go get github.com/SmartMeshFoundation/SmartRaiden/
  cd cmd/smartraiden
  go install
```
## Listing API
  SmartRaiden API has the main functions of Raiden Network API. Detailed function description please refer to[ Getting started with the SmartRaiden API](https://smartraiden.readthedocs.io/en/latest/api_walkthrough/) and [SmartRaiden’s API Documentation ](https://smartraiden.readthedocs.io/en/latest/api_walkthrough/#getting-started-with-the-smartraiden-api) in the official documentation of Raiden Network. The primary API list and description are as follows：

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
## Raiden Contract
These are the currently deployed raiden contract addresses for the ethereum testnet:
* Netting Channel Library: [0xad5cb8fa8813f3106f3ab216176b6457ab08eb75](https://ropsten.etherscan.io/address/0xad5cb8fa8813f3106f3ab216176b6457ab08eb75#code)
* Channel Manager Library: [0xdb3a4dbae2b761ed2751f867ce197c531911382a](https://ropsten.etherscan.io/address/0xdb3a4dbae2b761ed2751f867ce197c531911382a#code)
* Registry Contract: [0x68e1b6ed7d2670e2211a585d68acfa8b60ccb828](https://ropsten.etherscan.io/address/0x68e1b6ed7d2670e2211a585d68acfa8b60ccb828#code)
* Discovery Contract: [0x1e3941d8c05fffa7466216480209240cc26ea577](https://ropsten.etherscan.io/address/0x1e3941d8c05fffa7466216480209240cc26ea577#code)
## Usage
```
smartraiden [global options] command [command options] [arguments...]

VERSION:
   0.8

COMMANDS:
   help  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --address value                                            The ethereum address you would like raiden to use and for which a keystore file exists in your local system.
   --keystore-path "/Users/dognie/Library/Ethereum/keystore"  If you have a non-standard path for the ethereum keystore directory provide it using this argument.
   --eth-rpc-endpoint value                                   "host:port" address of ethereum JSON-RPC server.\n'
                                                                         'Also accepts a protocol prefix (ws:// or ipc channel) with optional port', (default: "/Users/dognie/Library/Ethereum/geth.ipc")
   --registry-contract-address value                          hex encoded address of the registry contract. (default: "0x52d7167FAD53835a2356C7A872BfbC17C03aD758")
   --listen-address value                                     "host:port" for the raiden service to listen on. (default: "0.0.0.0:40001")
   --api-address value                                        host:port" for the RPC server to listen on. (default: "127.0.0.1:5001")
   --datadir "/Users/dognie/Library/smartraiden"              Directory for storing raiden data.
   --password-file value                                      Text file containing password for provided account
   --debugcrash                                               enable debug crash feature
   --conditionquit value                                      quit at specified point for test
   --nonetwork                                                disable network, for example ,when we want to settle all channels
   --fee                                                      enable mediation fee
   --xmpp-server value                                        use another xmpp server  (default: "193.112.248.133:5222")
   --ignore-mediatednode-request                              this node doesn't work as a mediated node, only work as sender or receiver
   --enable-health-check                                      enable health check
   --verbosity value                                          Logging verbosity: 0=silent, 1=error, 2=warn, 3=info, 4=debug, 5=trace (default: 3)
   --vmodule value                                            Per-module verbosity: comma-separated list of <pattern>=<level> (e.g. eth/*=5,p2p=4)
   --backtrace value                                          Request a stack trace at a specific logging statement (e.g. "block.go:271")
   --debug                                                    Prepends log messages with call-site location (file and line number)
   --pprof                                                    Enable the pprof HTTP server
   --pprofaddr value                                          pprof HTTP server listening interface (default: "127.0.0.1")
   --pprofport value                                          pprof HTTP server listening port (default: 6060)
   --memprofilerate value                                     Turn on memory profiling with the given rate (default: 524288)
   --blockprofilerate value                                   Turn on block profiling with the given rate (default: 0)
   --cpuprofile value                                         Write CPU profile to the given file
   --trace value                                              Write execution trace to the given file
   --logfile value                                            redirect log to this the given file
   --help, -h                                                 show help
   --version, -v                                              print the version
   ```
## Requirements
geth >=1.7.3
