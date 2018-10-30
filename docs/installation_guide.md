# System Requirements and Installation Guide
## Introduction
Photon is an off-chain scaling solution, enabling instant, low-fee and scalable payments. It’s complementary to the Ethereum blockchain and Spectrum blockchain and works with ERC20 compatible token and ERC223 compatible token. Photon currently can works on Windows, Linux ,Android, iOS etc. Photon adds some new functions in version 0.9, such as , cooperative settlement,widraw without closing the channel, and more perfect third-party services(Photon-Monitoring and Photon-Path-Finder). In order to improve the user experience and to better fit for the mobile network, Photon adopts the Matrix communication mechanism and supports the crash recovery and channel charging function.
## Installation
The preferred way to install Photon is downloading a self contained application bundle from the [GitHub release page](https://github.com/SmartMeshFoundation/Photon/releases)
### Linux
Download the latest photon-<version>-[linux.tar.gz](https://github.com/SmartMeshFoundation/Photon/releases), and extract it:

```
tar -xvzf photon-<version>-linux.tar.gz
```

The Photon binary should work on most 64bit GNU/Linux distributions without any specific system dependencies, other than an Ethereum client installed in your system (see below). The Photon binary takes the same command line arguments as the photon script.
### macOS
Download the latest photon-<version>-[macOS.zip](https://github.com/SmartMeshFoundation/Photon/releases), and extract it:
```
unzip photon-<version>-macOS.zip
```
The resulting binary will work on any version of macOS from 10.12 onwards without any other dependencies. An Ethereum client is required.

### mobile
If you want to develop photon on mobile, you can execute the following commands

```
go get github.com/SmartMeshFoundation/Photon/
cd mobile
gomobile bind -target=android
gomobile bind -target=ios
```

### Dependencies
You will need to have an Spectrum client installed in your system.

- Check [this link](https://github.com/SmartMeshFoundation/Spectrum/wiki/Building-Specturm) for instructions on the smc client.

## For developers
If you plan to develop on the Photon source code, or the binary distributions do not work for your system, you can follow these steps to install a development version.

### Preliminaries
In order to work with  [`Photon`](https://github.com/SmartMeshFoundation/Photon), the following build dependencies are required:  

- **Go:**  `Photon`  is written in Go. To install, run one of the following commands:

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
### Installing Photon
With the preliminary steps completed, to install `Photon`, and all related dependencies run the following commands:
```
go get github.com/SmartMeshFoundation/Photon/ 
cd cmd/photon
./build.sh
```
**Updating**
```
git pull 
cd cmd/photon
./build.sh
```

### Requirements for Safe Usage
In order to use Photon correctly and safely there are some things that need to be taken care of by the user:

- **Layer 1 works reliably:** That means that you have a local public chain node, or connect to a stable server which own the public chain node , that is always synced and working reliably. If there are any problems or bugs on the client then Photon can not work reliably.   
- **Unique account for Photon:** We need to have a specific public chain account dedicated to Photon. Creating any manual transaction with the account that Photon uses, while the Photon client is running, can result in undefined behaviour.  
- **Photon account need sufficient Public Chain Coin:** Photon will not to warn you whether there is enough Public chain coin in photon account or not, so you need to transfer sufficient public chain coin to the account in order to maintain your current open chanels and go through their entire cycle.    
- **Persistency of local DB:** Your local state database is located at ~/.photon. This data should not be deleted by the user or tampered with in any way. Frequent backups are also recommended. Deleting this directory could mean losing funds.    
- **Photon node can be offline or online according to usage :** In photon,some nodes ,such as mobile nodes(just for payment requirment), can be offline after delegating evidence to the PhotonMonitor; Another nodes, such as meshbox nodes,which provide the mediate transfer service, need to confirm that the inside photon node is always running, If it crashes for whatever reason you are responsible to restart it.

## Firing it up
### Spectrum
Run Photon nodes on the Spectrum testnet
#### Installing Spectrum
For prerequisites and detailed build instructions please read the  [Installation Instructions](https://github.com/SmartMeshFoundation/Spectrum/wiki/Building-Specturm)  on the wiki.

Building smc requires both a Go (version 1.9.2 or later) and a C compiler. You can install them using your favourite package manager.
#### Starting Spectrum testnet
Run boot script
```sh
smc  --datadir=. --testnet --syncmode full     --ws --wsapi  "eth,admin,web3,net,debug,personal"   --rpc  --rpccorsdomain "*" --rpcapi "eth,admin,web3,net,debug,personal"   --wsaddr "0.0.0.0" --rpcaddr "0.0.0.0"   --wsorigins "*"
```
#### Starting Photon node
Before start the Photon, you need to creat accounts on spectrum,such as "0x97cd7291f93f9582ddb8e9885bf7e77e3f34be40". Also you need to set the keystore path and password files or password in the boot script.
```sh
photon  --datadir=.photon  --address="0x97cd7291f93f9582ddb8e9885bf7e77e3f34be40"  --keystore-path ./keystore --registry-contract-address 0xb3aE919aB595f5844cba80499ee6423688E06F89 --password-file pass.txt --eth-rpc-endpoint ws://192.168.124.13:28546
```
After you start the photon node,you can register the token in the photonnetwork and use the various functions provided by photon.
#### Deployed contract address
- Specrum  Mainnet:RegistryAddress=0xCf3C7400C227be86FcdB2c9Be7DEf5c671087620
- Specrum  Testnet:RegistryAddress=0xb3aE919aB595f5844cba80499ee6423688E06F89 
