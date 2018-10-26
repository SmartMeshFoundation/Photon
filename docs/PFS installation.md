# PFS(SmartRaiden Path Finding Service)
PFS will be a smartraiden supporting function server written in go.

#Installing PFS

 * A cluster of individual components, dealing with different aspects of the
     SmartRaiden protocol.
     
## Requirements
 - Go 1.8+
 - Postgres v10.0+
 
## Setting up a development environment
```bash
# Get the code
git clone https://github.com/SmartMeshFoundation/SmartRaiden-Path-Finder.git
cd SmartRaiden-Path-Finder

# Build it
go get github.com/constabulary/gb/...
gb build
```
## Congiguration

### Postgres database setup
PFS requires a postgres database engine,version 10.0 or later.
* Postgres downl
  * "https://www.postgresql.org/download/"
* Create role:
 ```bash
 sudo -u postgres createuser -P pfs     # prompts for password
 ```
* Create database for PFS server:
 ```bash
 sudo -u postgres createdb -O pfs pfs_nodeinfos
 ```
(On macOs omit `sudo -u postgres` from the above commands.)

### PFS running parameters

* Configuration parameters are stored at the file named "pathfinder.yaml",you cannnot rename it,
    this file is in the same directory as the executable file.
* The parameters that need to be configured according to the actual application environment are:
   * registory_address :hex encoded address of the registry contract(e.g:0xd66d3719E89358e0790636b8586b539467EDa596)
   * address :The ethereum address you would link pfs to use sign transaction on ethereum
   * keystore_path :Path for the ethereum keystore directory
   * eth_rpc_endpoint :"host:port" address of ethereum JSON-RPC server,
                          also accepts a protocol prefix (ws:// or ipc channel) with optional port
   * password-file :Text file containing password for provided account
   * chain_id :Chain id identifies the current chain and is used for replay protection
   
   
   * pfs->server_name:PFS name(default:"localhost")
   * ratelimited->stationary_feerate_default:Default fee rate,if you do not fill in,default as "0.0001"
   
   * logging->params->path :The path of log file 
   
   * database->nodeinfos :The database connection string in postgres,that is data source name.
   

## Starting a PFS server

```bash
./bin/smartraiden-pathfinding-service
```

## Todo
It's still very much a work in progress.