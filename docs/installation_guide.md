# System Requirements and Installation Guide
## Introduction
Smartraiden is an off-chain scaling solution, enabling instant, low-fee and scalable payments. It’s complementary to the Ethereum blockchain and Spectrum blockchain and works with ERC20 compatible token and ERC223 compatible token. SmartRaiden currently can works on Windows, Linux ,Android, iOS etc. smartraiden adds some new functions in version 0.9, such as , cooperative settlement,widraw without closing the channel, and more perfect third-party services(SmartRaiden-Monitoring and SmartRaiden-Path-Finder). In order to improve the user experience and to better fit for the mobile network, Smartraiden adopts the Matrix communication mechanism and supports the crash recovery and channel charging function.
## Installation
The preferred way to install SmartRaiden is downloading a self contained application bundle from the [GitHub release page](https://github.com/SmartMeshFoundation/SmartRaiden/releases)
### Linux
Download the latest smartraiden-<version>-[linux.tar.gz](https://github.com/SmartMeshFoundation/SmartRaiden/releases), and extract it:

```
tar -xvzf smartraiden-<version>-linux.tar.gz
```

The SmartRaiden binary should work on most 64bit GNU/Linux distributions without any specific system dependencies, other than an Ethereum client installed in your system (see below). The Raiden binary takes the same command line arguments as the raiden script.
### macOS
Download the latest smartraiden-<version>-[macOS.zip](https://github.com/SmartMeshFoundation/SmartRaiden/releases), and extract it:
```
unzip smartraiden-<version>-macOS.zip
```
The resulting binary will work on any version of macOS from 10.12 onwards without any other dependencies. An Ethereum client is required.

### mobile
If you want to develop smartraiden on mobile, you can execute the following commands

```
go get github.com/SmartMeshFoundation/SmartRaiden/
cd mobile
gomobile bind -target=android
gomobile bind -target=ios
```

### Dependencies
You will need to have an Ethereum client installed in your system.

- Check [this link](https://github.com/ethereum/go-ethereum/wiki/Building-Ethereum) for instructions on the go-ethereum client.

## For developers
If you plan to develop on the SmartRaiden source code, or the binary distributions do not work for your system, you can follow these steps to install a development version.

### Preliminaries
In order to work with  [`SmartRaiden`](https://github.com/SmartMeshFoundation/SmartRaiden), the following build dependencies are required:  

- **Go:**  `SmartRaiden`  is written in Go. To install, run one of the following commands:

	**Note**: The minimum version of Go supported is Go 1.9. We recommend that users use the latest version of Go, which at the time of writing is  [`1.10`](https://blog.golang.org/go1.10).
	On Linux:
	```
	sudo apt-get install golang-1.10-go
	```
	On Mac OS X:
	```
	brew install go
	```
	Alternatively, one can download the pre-compiled binaries hosted on the [golang download page](https://golang.org/dl/). If one seeks to install from source, then more detailed installation instructions can be found [here](http://golang.org/doc/install).
	At this point, you should set your  `$GOPATH`  environment variable, which represents the path to your workspace. By default,  `$GOPATH`  is set to  `~/go`. You will also need to add  `$GOPATH/bin`  to your  `PATH`. This ensures that your shell will be able to detect the binaries you install.
	```
	export GOPATH=~/gocode
	export PATH=$PATH:$GOPATH/bin
	```	
### Installing SmartRaiden
With the preliminary steps completed, to install `SmartRaiden`, and all related dependencies run the following commands:
```
go get github.com/SmartMeshFoundation/SmartRaiden/ 
cd cmd/smartraiden
go install
```
**Updating**
```
git pull 
cd cmd/smartraiden
go install
```

### Requirements for Safe Usage
In order to use SmartRaiden correctly and safely there are some things that need to be taken care of by the user:

- **Layer 1 works reliably:** That means that you have a local public chain node, or connect to a stable server which own the public chain node , that is always synced and working reliably. If there are any problems or bugs on the client then SmartRaiden can not work reliably.   
- **Unique account for SmartRaiden:** We need to have a specific public chain account dedicated to SmartRaiden. Creating any manual transaction with the account that SmartRaiden uses, while the SmartRaiden client is running, can result in undefined behaviour.  
- **SmartRaiden account need sufficient Public Chain Coin:** SmartRaiden will not to warn you whether there is enough Public chain coin in smartraiden account or not, so you need to transfer sufficient public chain coin to the account in order to maintain your current open chanels and go through their entire cycle.    
- **Persistency of local DB:** Your local state database is located at ~/.smartraiden. This data should not be deleted by the user or tampered with in any way. Frequent backups are also recommended. Deleting this directory could mean losing funds.    
- **SmartRaiden node can be offline or online according to usage :** In smartraiden,some nodes ,such as mobile nodes(just for payment requirment), can be offline after delegating evidence to the SmartMonitor; Another nodes, such as meshbox nodes,which provide the mediate transfer service, need to confirm that the inside smartraiden node is always running, If it crashes for whatever reason you are responsible to restart it.

## Firing it up
### Ethereum
Run the Ethereum client and let it sync with the Ropsten testnet:
```
geth --testnet --fast --rpc --rpcapi eth,net,web3 --bootnodes "enode://20c9ad97c081d63397d7b685a412227a40e23c8bdc6688c6f37e97cfbc22d2b4d1db1510d8f61e6a8866ad7f0e17c02b14182d37ea7c3c8b9c2683aeb6b733a1@52.169.14.227:30303,enode://6ce05930c72abc632c58e2e4324f7c7ea478cec0ed4fa2528982cf34483094e9cbc9216e7aa349691242576d552a2a56aaeae426c5303ded677ce455ba1acd9d@13.84.180.240:30303"
```
Unless you already have an account you can also create one in the console by invoking `personal.newAccount()`.

If problems arise for above method, please see [the Ropsten README](https://github.com/ethereum/ropsten) for further instructions.

Then launch SmartRaiden with the default testnet keystore path:
```
smartraiden --keystore-path  ~/.ethereum/testnet/keystore
```

### Spectrum
Run SmartRaiden nodes on the Spectrum testnet
#### Installing Spectrum
For prerequisites and detailed build instructions please read the  [Installation Instructions](https://github.com/SmartMeshFoundation/Spectrum/wiki/Building-Specturm)  on the wiki.

Building geth requires both a Go (version 1.9.2 or later) and a C compiler. You can install them using your favourite package manager.
#### Starting Spectrum testnet
Run boot script
```sh
geth  --datadir=. --testnet --syncmode full     --ws --wsapi  "eth,admin,web3,net,debug,personal"   --rpc  --rpccorsdomain "*" --rpcapi "eth,admin,web3,net,debug,personal"   --wsaddr "0.0.0.0" --rpcaddr "0.0.0.0"   --wsorigins "*"
```
