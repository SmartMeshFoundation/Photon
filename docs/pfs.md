# SmartRaiden Pathfinding Service Specification

A centralized path finding service that has a global view on a token network and provides suitable payment paths for SmartRaiden nodes.


# Environment Construction 
## build pfs 
You can refer to [installation documentation](./PFSinstallation.md).
Then execute it in the same directory of configuration file (pathfinder.yaml).
```sh
./smartraiden-pathfinding-service
```
**pfs Default listening port is 9001**

## 
## Public Interface

### GET /api/1/pfs/<channel_identifier>
Query the latest balance_proof  
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


### PUT /pathfinder/<peer_address>/balance
Update balance_proof 
Via API offered below:  
`PUT http://127.0.0.1:9001/pathfinder/0x10b256b3C83904D524210958FA4E7F9cAFFB76c6/balance`

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
```json
{
    "errcode": "M_OK",
    "error": "true"
}
```
If submitted, it is not the latest balance_proof：
```json
{
    "errcode": "M_INVALID_ARGUMENT_VALUE",
    "error": "Outdated balance proof"
}
```

### PUT /pathfinder/<peer_address>/set_fee_rate
Set smartraiden node charging rules and update to update to PFS

`PUT http://127.0.0.1:9001/pathfinder/0x3607806E038fED0985567992188E919802486bf3/set_fee_rate`

Example Request：

```json
{

	"channel_id":"0x6bebe91a40c39fc3ffcd6adc8dbc46052a02ba6912e45b025e058a07c5f2f0dd",
	"fee_rate": "1",
	"signature":"xC487RM+e1TITDeEjvNq3UMLPc2mGHyiw4T7k6TG1WYRT9GMmptX/8NPyYLdpxfwhpC+TO6dhPMBs57rsZiOIxs="

}
```
Example Response：

```json
{
    "errcode": "M_OK",
    "error": "true"
}
```
### POST /pathfinder/<peer_address>/paths
Query routing, return to the lowest cost path.
`POST http://127.0.0.1:9001/pathfinder/0x201B20123b3C489b47Fde27ce5b451a0fA55FD60/paths`
Example Request：
```json
{
	"peer_from": "0x201B20123b3C489b47Fde27ce5b451a0fA55FD60",
	"peer_to": "0x0D0EFCcda4f079C0dD1B728297A43eE54d7170Cd",
	"limit_paths": 6,
	"send_amount": 100,
	"signature": "KUmUDRbyJzrt5CR0Tlgvlh2PZ+Q8c8m4rdaHU+Cu9yIxfQ9drHw99qiWs/qXbtr/ok8m7N0ZUvUOX3ldxhcSXBw="
}
```
Example Response：
```json
[
    {
        "path_id": 0,
        "path_hop": 2,
        "fee": 0.03999999898951501,
        "result": [
            "0x10b256b3c83904d524210958fa4e7f9caffb76c6",
            "0xce92bddda9de3806e4f4b55f47d20ea82973f2d7"
        ]
    }
]
```

### POST /pathfinder/<peer_address>/get_fee_rate
Query node charging information

`POST http:127.0.0.1:9001/pathfinder/0x3607806E038fED0985567992188E919802486bf3/get_fee_rate`
Example Request：
```json
{
	"obtain_obj": "0x3607806E038fED0985567992188E919802486bf3",
	"channel_id": "0x6bebe91a40c39fc3ffcd6adc8dbc46052a02ba6912e45b025e058a07c5f2f0dd",
	"signature": "A4SqcEM4+z8oUGdsoBf0Kj8T+JfrRZi2uBQTpSCwdHRNswm8nVvj/s7JNN3eeerQZxdpfsMLzkrTN6K2J2NWVBw="
	
}
```
Example Response
```json
{
    "result": {
        "0x6bebe91a40c39fc3ffcd6adc8dbc46052a02ba6912e45b025e058a07c5f2f0dd": {
            "channel_id": "0x6bebe91a40c39fc3ffcd6adc8dbc46052a02ba6912e45b025e058a07c5f2f0dd",
            "peer_address": "0x3607806e038fed0985567992188e919802486bf3",
            "fee_rate": "1",
            "effective_time": 1540518875801
        }
    }
}
```




