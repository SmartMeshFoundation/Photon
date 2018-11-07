# Photon Pathfinding Service Specification

A centralized path finding service that has a global view on a token network and provides suitable payment paths for Photon nodes.


# Environment Construction 
## Build pfs 
You can refer to [installation documentation](./PFSinstallation.md).
Then execute it in the same directory of configuration file (pathfinder.yaml).
```sh
./photon-pathfinding-service
```
**Pfs default listening port is 9001**

## 
## Public Interface

### GET /api/1/pfs/*(channel_identifier)*
Query the latest balance_proof (from photon node) 
Via API offered below   

`GET http://127.0.0.1:5001/api/1/pfs/0x622924d11071238ac70c39b508c37216d1a392097a80b26f5299a8d8f4bc0b7a`

Example Response：   
```json
{
    "balance_proof": {
        "nonce": 4,
        "transfer_amount": 100,
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


### PUT /pfs/1/*(peer_address)*/balance
Update balance_proof from PFS
Via API offered below: 

`PUT http://127.0.0.1:9001/pfs/1/0x10b256b3C83904D524210958FA4E7F9cAFFB76c6/balance`

Example Request：
```json
{
    "balance_proof": {
        "nonce": 4,
        "transfer_amount": 100,
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

###  GET /pfs/1/account_rate/*(peer_address)*
Query account charging rate  
Example Request： 
`GET http://127.0.0.1:9001/pfs/1/account_rate/0x6d946D646879d31a45bCE89a68B24cab165E9A2A`
Example Response： 
```json
{
    "fee_policy": 2,
    "fee_constant": 5,
    "fee_percent": 10000
}
```

### GET /pfs/1/token_rate/*(token_address)*/*(peer_address)*
Query the charging rate of a node on a certain token 
Example Request： 
`GET http://127.0.0.1:9001/pfs/1/token_rate/0x83073FCD20b9D31C6c6B3aAE1dEE0a539458d0c5/0x6d946D646879d31a45bCE89a68B24cab165E9A2A`
Example Response： 
```json
{
    "fee_policy": 2,
    "fee_constant": 5,
    "fee_percent": 10000
}
```
### GET /pfs/1/channel_rate/*(channel_identifier)*/*(peer_address)*  
Query the charging rate of a node on a certain channel 
Example Request：  
`GET http://127.0.0.1:9001/pfs/1/channel_rate/0x24bab913507cc9fcaa2c1efc4966ab35246448f19a7f0d44db21b8b3601db654/0x6d946D646879d31a45bCE89a68B24cab165E9A2A`
Example Response： 
```json
{
    "fee_policy": 2,
    "fee_constant": 7,
    "fee_percent": 10000
}
```


### POST /pfs/1/paths
Query the transfer routing, return the lowest cost path.

Example Request： 
`POST http://127.0.0.1:9001/pfs/1/paths`

```json
{
	"peer_from":"0x201B20123b3C489b47Fde27ce5b451a0fA55FD60",
	"peer_to":"0x0D0EFCcda4f079C0dD1B728297A43eE54d7170Cd",
	"token_address":"0x37346b78de60f4F5C6f6dF6f0d2b4C0425087a06",
	"send_amount":20000
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

When using pfs, node startup requires the `--pfs` ` --fee` parameter.
 - `pfs` : pathfinder service host,The PFS main network and test network have been [deployed online](./pfs_online_bulletin.md). 
 - `fee` : enable mediation fee, After opening, you can query and set the rate.


For example, using pfs on the test network 

`POST http://transport01.smartmesh.cn:7001/pfs/1/paths`






