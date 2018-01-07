# raiden-network
Go implementation of the Ethereum yellowpaper

## About 
raiden-network is a pure Go  implementation of the raiden network,it's goal is to cooperate with python raiden and gives users another choice for mobile device.
## Build
```
go get github.com/SmartMeshFoundation/raiden-network/tree/master/cmd/goraiden
```
## Usage

```                                                                                                                                                    
VERSION:                                                                                                                                                              
   0.1                                                                                                                                                                
                                                                                                                                                                      
COMMANDS:                                                                                                                                                             
     help, h  Shows a list of commands or help for one command                                                                                                        
                                                                                                                                                                      
GLOBAL OPTIONS:                                                                                                                                                       
   --address value                                                             The ethereum address you would like raiden to use and fo                               
r which a keystore file exists in your local system.                                                                                                                  
   --keystore-path "C:\Users\Administrator\AppData\Roaming\Ethereum\keystore"  If you have a non-standard path for the ethereum keystor                               
e directory provide it using this argument.                                                                                                                           
   --eth-rpc-endpoint value                                                    "host:port" address of ethereum JSON-RPC server.\n'                                    
                                                                                          'Also accepts a protocol prefix (ws:// or ipc                               
 channel) with optional port', (default: "\\\\.\\pipe\\geth.ipc")                                                                                                     
   --registry-contract-address value                                           hex encoded address of the registry contract. (default:                                
"0xCf3C7400C227be86FcdB2c9Be7DEf5c671087620")                                                                                                                         
   --discovery-contract-address value                                          hex encoded address of the discovery contract. (default:                               
 "0x5a93A5E5b754898f06F7A0f4abac419547600B25")                                                                                                                        
   --listen-address value                                                      "host:port" for the raiden service to listen on. (defaul                               
t: "0.0.0.0:40001")                                                                                                                                                   
   --rpccorsdomain value                                                       Comma separated list of domains to accept cross origin r                               
equests.                                                                                                                                                              
                                                                                     (localhost enabled by default) (default: "http://l                               
ocalhost:* /*")                                                                                                                                                       
   --logging value                                                             ethereum.slogging config-string{trace,debug,info,warn,er                               
ror,critical  (default: "trace")                                                                                                                                      
   --logfile value                                                             file path for logging to file                                                          
   --max-unresponsive-time value                                               Max time in seconds for which an address can send no pac                               
kets and                                                                                                                                                              
                                                                                              still be considered healthy. (default: 12                               
0)                                                                                                                                                                    
   --send-ping-time value                                                      Time in seconds after which if we have received no messa                               
ge from a                                                                                                                                                             
                                                                                              node we have a connection with, we are go                               
ing to send a PING message (default: 60)                                                                                                                              
   --rpc                                                                       Start with or without the RPC server. Default is to star                               
t                                                                                                                                                                     
                                                                                              the RPC server                                                          
   --api-address value                                                         host:port" for the RPC server to listen on. (default: "0                               
.0.0.0:5001")                                                                                                                                                         
   --datadir "C:\Users\Administrator\AppData\Roaming\GoRaiden"                 Directory for storing raiden data.                                                     
   --password-file value                                                       Text file containing password for provided account                                     
   --help, -h                                                                  show help                                                                              
   --version, -v                                                               print the version                                                                      
                                                                                                                                                                      
```

## Requirements

geth >=1.7.3