# Photon Pathfinding Service Specification

A centralized path finding service that has a global view on a token network and provides suitable payment paths for Photon nodes.

# Environment Construction 
## Build pfs 
You can refer to [installation documentation](./PFSinstallation.md).
Then execute it in the same directory of configuration file (pathfinder.yaml).
```sh
git clone https://github.com/SmartMeshFoundation/Photon-Path-Finder.git

cd /cmd/photon-pathfinding-service
go install
```
**Pfs default listening port is 9001**
Start service
```sh
photon-pathfinding-service  --eth-rpc-endpoint ws://192.168.124.13:5555  --registry-contract-address  0x3400aa968662Cfb6f7ea911Cd18254350e0C3d21  --port  9001  --verbosity 5
```
Parameters: 

- `eth-rpc-endpoint` Full node
- `registry-contract-address` Contract address
- `port` listening port

When using the pfs service, photon startup must be added:
- `fee` charge the fee
- `pfs` pfs service address
### Default address

- Spectrum mainnet http://transport01.smartmesh.cn:7000

- Spectrum testnet http://transport01.smartmesh.cn:7001
## Public Interface

###  Query the balance proof

    GET /api/1/pfs/*(channel_identifier)*

This interface is the query balance proof interface provided by photon for pfs(from photon node).

Example Request: 

`GET http://127.0.0.1:5001/api/1/pfs/0x622924d11071238ac70c39b508c37216d1a392097a80b26f5299a8d8f4bc0b7a`

Example Response：

```json
{
    "balance_proof": {
        "nonce": 4,
        "transfer_amount": 100000000000000000000,
        "locks_root": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "channel_identifier": "0x622924d11071238ac70c39b508c37216d1a392097a80b26f5299a8d8f4bc0b7a",
        "open_block_number": 7218036,
        "addition_hash": "0x158661e64ada377ec4164f91426b5c93dec9d90b9ded2944d9a82c55ec292022",
        "signature": "rupUzBuIRtUj2dEbvjMlUX7gm+6P1ZMOSHD+KMHTZegDtPLUK53XhaxXhvpcXgH48nmgCgkFmyBThaTzgPEalxw="
    },
    "balance_signature": "ADLsvd1iaRzfihYZkZeyNdjf3Grh6auZK/6vNQms95FWh7zKaiT6Rtzl39LVubRpQMPrlei5SEqfFsDWlfxUGhs=",
    "lock_amount": 0
}
```

### Update balanceproof 
PUT /pfs/1/*(peer_address)*/balance 

Update balance_proof from PFS.

Example Request: 

`PUT http://127.0.0.1:9001/pfs/1/0x10b256b3C83904D524210958FA4E7F9cAFFB76c6/balance`

PayLoad:  

```json
{
    "balance_proof": {
        "nonce": 4,
        "transfer_amount": 100000000000000000000,
        "locks_root": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "channel_identifier": "0x622924d11071238ac70c39b508c37216d1a392097a80b26f5299a8d8f4bc0b7a",
        "open_block_number": 7218036,
        "addition_hash": "0x158661e64ada377ec4164f91426b5c93dec9d90b9ded2944d9a82c55ec292022",
        "signature": "rupUzBuIRtUj2dEbvjMlUX7gm+6P1ZMOSHD+KMHTZegDtPLUK53XhaxXhvpcXgH48nmgCgkFmyBThaTzgPEalxw="
    },
    "balance_signature": "ADLsvd1iaRzfihYZkZeyNdjf3Grh6auZK/6vNQms95FWh7zKaiT6Rtzl39LVubRpQMPrlei5SEqfFsDWlfxUGhs=",
    "lock_amount": 0
}	
```
Example Response：

**200 OK**

If submitting an invalid  balance_proof, the response are as follows：

**400 Bad Request** and 

```json
{
    "Error": "illegal signature of balance message, for participant"
}
```


###  Query account charging rate
GET /pfs/1/account_rate/*(peer_address)*

Example Request: 
`GET http://127.0.0.1:9001/pfs/1/account_rate/0x6d946D646879d31a45bCE89a68B24cab165E9A2A`

Example Response：

```json
{
    "fee_policy": 2,
    "fee_constant": 5,
    "fee_percent": 10000
}
```

### Query the charging rate of a node on a certain token 
GET /pfs/1/token_rate/*(token_address)*/*(peer_address)*

Example Request: 

`GET http://127.0.0.1:9001/pfs/1/token_rate/0x83073FCD20b9D31C6c6B3aAE1dEE0a539458d0c5/0x6d946D646879d31a45bCE89a68B24cab165E9A2A`

Example Response：

```json
{
    "fee_policy": 2,
    "fee_constant": 5,
    "fee_percent": 10000
}
```
### Query the charging rate of a node on a certain channel 
GET /pfs/1/channel_rate/*(channel_identifier)*/*(peer_address)*  

Example Request: 
`GET http://127.0.0.1:9001/pfs/1/channel_rate/0x24bab913507cc9fcaa2c1efc4966ab35246448f19a7f0d44db21b8b3601db654/0x6d946D646879d31a45bCE89a68B24cab165E9A2A`

Example Response： 

```json
{
    "fee_policy": 2,
    "fee_constant": 7,
    "fee_percent": 10000
}
```
### Query the lowest cost path
POST /pfs/1/paths

Query the transfer routing, which will return all the lowest cost path.

Example Request: 
`POST http://127.0.0.1:9001/pfs/1/paths`

PayLoad: 

```json
{
	"peer_from":"0x201B20123b3C489b47Fde27ce5b451a0fA55FD60",
	"peer_to":"0x0D0EFCcda4f079C0dD1B728297A43eE54d7170Cd",
	"token_address":"0x37346b78de60f4F5C6f6dF6f0d2b4C0425087a06",
	"send_amount":20000000000000000000000
}
```

Example Response：

```json
[
    {
        "path_id": 0,
        "path_hop": 2,
        "fee": 4,
        "result": [
            "0x10b256b3c83904d524210958fa4e7f9caffb76c6",
            "0x151e62a787d0d8d9effac182eae06c559d1b68c2"
        ]
    },
    {
        "path_id": 1,
        "path_hop": 2,
        "fee": 4,
        "result": [
            "0x10b256b3c83904d524210958fa4e7f9caffb76c6",
            "0xce92bddda9de3806e4f4b55f47d20ea82973f2d7"
        ]
    }
]
```

 *tips：*
- fee_constant: Fixed charge 
- fee_percent: fee rate
  
  Where `fee_constant` is the fixed rate, for example, 5 means that the fixed fee is 5 tokens, and setting it to 0 means no charge. `fee_percent` is the proportional rate, calculated as the transaction amount/`fee_percent`, such as transaction amount 50000000000000000000000, `fee_percent`=10000, then the commission ratio part = 50000000000000000000000/10000=5000000000000000000, set to 0 means no charge.
 Charge rule fee = `fee_constant` + amount/`fee_percent`

 There are three charging modes for a node:
- account_fee    Node charging
- token_fee      Node charging  on Specific token
- channel_fee    Node charging at a certain channel
 The priority of the three charging modes is：`channel_fee`>`token_fee`>`account_fee`

When using pfs, node startup does not require the `--pfs` ` --fee` parameter, because they are default setting.If you want to change the PFS，you can add the `--pfs` and PFS address to the script code,and if you do not want to charge the fee, you can add the `--disable-fee` to the startup script to use the p2p path finding.









