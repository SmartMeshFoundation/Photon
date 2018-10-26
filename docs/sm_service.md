# SmartRaiden Monitoring Service
SmartRaiden Monitoring Service, as SM, primarily focuses on mobile platforms. As user's mobile devides disconnected from SmartRaiden, they would be automatically execute `UpdateTransfer` and `Punish` according to delegation, to secure user's assets.  
## Build and Installation
### Build
```bash
cd cmd/smartraidenmonitoring
go install
```
### Run as Service
Smartraidenmonitoring requires a SmartRaiden node to charge the service. The specific Token and the charge can be configured. The SmartRaiden node must be running on port 127.0.0.1:5001, otherwise it will not work.
an example run script:
```
smartraidenmonitoring --datadir=.smartraidenmonitoring --eth-rpc-endpoint ws://192.168.124.13:5555  --address="0x292650fee408320D888e06ed89D938294Ea42f99" --keystore-path=/Users/bai/privnet3/keystore  --registry-contract-address 0x7B319fB135811caeED9969E6a97544f74E312A65  --password-file 123 --verbosity 5  --debug   --smt 0x40db17463AD4A00cc824a37d851725aC7eA4E0B6
```
The field descriptions are shown in SmartRaiden Monitoring Start.
## How to make a use of SmartRaiden Monitoring Service ?
Well, how to use this SM service ? In-depth Tutorials will be presented below.   

## Environment Construction
A complete workflow of a SmartRaiden Monitoring Service, at least three nodes and a SmartRaiden Monitoring Service are required.  

Node Names|Node Address  
--|--
SmartRaiden Monitoring|0x6b9e4d89ee3828e7a477ea9aa7b62810260e27e9
Delegated Charge Node|0x69c5621db8093ee9a26cc2e253f929316e6e5b92
Channel Participant `Alice`|0x31ddac67e610c22d19e887fb1937bee3079b56cd
Channel Participant `Bob`|0xf0f6e53d6bbb9debf35da6531ec9f1141cd549d5

**SmartRaiden Monitoring Start**  
Run script code below to kickstart.  
```sh
smartraidenmonitoring  --datadir=.smartraidenmonitoring --eth-rpc-endpoint ws://127.0.0.1:5555  --address="0x6b9e4d89ee3828e7a477ea9aa7b62810260e27e9" --keystore-path ~/.ethereum/keystore --registry-contract-address 0x4dc3388E72e45E99061Ec4Fe17Db2ebfe3B4341f  --password-file /home/niexin/niexin/data.txt  --smt  0xc0dfdD7821c762eF38F86225BD45ff4e912fFA20
echo "quit ok"
```

**Parameter Description**  
- `datadir` : local storage path of data for SmartRaiden Monitoring.  
- `eth-rpc-endpoint` : Full node.  
- `address` : the address of SmartRaiden Monitoring Service.  
- `keystore-path` : the path of keystore  
- `registry-contract-address` : the address of contract  
- `password-file` : the path of password  
- `smt` : the address of token  

*Default port for SmartRaiden Monitoring : 6000*  

**Start-up Delegated Charge Node**
Run script below to operate delegated charged node.  
```sh
echo run smartraiden...
smartraiden --datadir=.smartraiden --api-address=0.0.0.0:5001 --listen-address=127.0.0.1:40011 --address="0x69c5621db8093ee9a26cc2e253f929316e6e5b92" --keystore-path ~/.ethereum/keystore --registry-contract-address 0x4dc3388E72e45E99061Ec4Fe17Db2ebfe3B4341f  --password-file 123  --eth-rpc-endpoint ws://127.0.0.1:5555  --ignore-mediatednode-request   
echo "quit ok"
```

*`ignore-mediatednode-request` must be added when start delegated charge node, in case it works as a mediated node of smartraiden and make mistakes at fee-charging.*  

**Run Channel Participant Nodes:**  
Same as Start-up Delegated Charge Node, you can run codes below to achieve that.  

**Alice:**      
```sh
echo privnet
echo run smartraiden...
smartraiden --datadir=.smartraiden --api-address=0.0.0.0:5002 --listen-address=127.0.0.1:40002 --address="0x31ddac67e610c22d19e887fb1937bee3079b56cd" --keystore-path ~/.ethereum/keystore --registry-contract-address 0x4dc3388E72e45E99061Ec4Fe17Db2ebfe3B4341f --password-file 123  --eth-rpc-endpoint ws://127.0.0.1:5555  
echo "quit ok"
```
**Bob:**    
```sh
echo privnet
echo run smartraiden...
smartraiden --datadir=.smartraiden --api-address=0.0.0.0:5003 --listen-address=127.0.0.1:40003 --address="0xf0f6e53d6bbb9debf35da6531ec9f1141cd549d5" --keystore-path ~/.ethereum/keystore  --registry-contract-address 0x4dc3388E72e45E99061Ec4Fe17Db2ebfe3B4341f --password-file 123  --eth-rpc-endpoint ws://127.0.0.1:5555 
echo "quit ok"
```
Till now, you've completed works for environment construction and nodes have been started. Let me take you to next steps to see how they actually work.  

## SmartRaiden Monitoring Showcase

**Scenario :**  
There is a payment channel between Alice and Bob. For some specific reasons, in next few days Bob is about to go offline. when Bob is during the offline period, if Alice closes payment channel, how can Bob update the newest version of `BalanceProof` of Alice, so that he can ensure his Security of funds.  

This is exactly what our SmartRaiden Monitoring Service need to do.   

When Bob realizes that he's going to get off from network, he can delegate SmartRaiden Monitoring to update `BalanceProof` . Once Bob does disconnect from SmartRaiden, the SM node will update `BalanceProof` on behalf of Bob to retrieve the deserved token from the channel.  

There are other cases,such as, a channel participant is a fraudulent actor and attempts to steal tokens from this channel. If he unlocks a abandon transfer which he has declared to dispose, then Our SmartRaiden provides `punish` feature to prevent this case to happen and fraudulent actors will be punished. The SmartRaiden Monitoring also provides`punish` services.Latter we will give you an invidual showcase for this.  

### Normal Delegation, not including punish
#### 1. Alice makes a transfer to Bob
**Via API offered below**  
`POST  /api/<version>/transfer/<token_address>/target_address`
**Example Request:**    
`POST http://127.0.0.1:5002/api/1/transfers/0xc0dfdD7821c762eF38F86225BD45ff4e912fFA20/0x31DdaC67e610c22d19E887fB1937BEE3079B56Cd`  

**PAYLOAD**    
```json
{
    "amount":200,
    "fee":0,
    "is_direct":false
}
```

#### 2. Bob queries delegation info which need to submit to SmartRaiden Monitoring Node
**Via API offered below**   
`GET: /api/<version>/thirdparty/<channel_address>/<thirdparty_address>`    
**Example Request :**    
`GET http://127.0.0.1:5003/api/1/thirdparty/0x2f6418b01422de6cc84fd52e4378fc4449436aeadf69dd543e79e87ee38b6dc8/0x6B9E4D89EE3828e7a477eA9AA7B62810260e27E9`  

**Example Response:**    
```json
{
    "channel_identifier": "0x2f6418b01422de6cc84fd52e4378fc4449436aeadf69dd543e79e87ee38b6dc8",
    "open_block_number": 4392735,
    "token_network_address": "0x445b92d89e21bea3510597a12fefa037718017c1",
    "partner_address": "0x31ddac67e610c22d19e887fb1937bee3079b56cd",
    "update_transfer": {
        "nonce": 2,
        "transfer_amount": 200,
        "locksroot": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "extra_hash": "0xb07638a54d4cbbec93d687d840f13c472df2addae79944dfa00aa2b66ec29b1c",
        "closing_signature": "vwA2lAU414nPajnk5s7dXS/8LHzgCXZ/gVt3lbrzbOhbLkbwcmEW3XtZAjN0qqNYOC8F2opBu7N05FLckrQ5kRw=",
        "non_closing_signature": "aJeOYwRShmkXUBdVTUEKoi5HKKElxn3NfcM8zCRGPcdYJPMv/u+tz9gS3vkb6ypGyvfn1yzKOU1vax3C5/vGMRw="
    },
    "unlocks": null,
    "punishes": null
}
```

#### 3. SmartRaiden Monitoring Service queries whether Tokens are sufficient in delegated charge node.
**Via API offered below**    
`GET /fee/<delegater_address>`    
**Example Request :**    
`http://127.0.0.1:6000/fee/0xf0f6E53d6bbB9Debf35Da6531eC9f1141cd549d5`   

**Example Response:**    
```json
{
    "Available": 0,
    "NeedSmt": 0
}
```
- `Available` - available balance  
-  `NeedSmt` - fees SM node need to pay  

If our delegator has no sufficient fund deposited in Delegated charge node, then he needs to make enough deposit, just  normal smartraiden transfers. In our case, Bob deposits/transfers 20 tokens into Delegated charge node.

**Example Request :**    
`POST http://127.0.0.1:5003/api/1/transfers/0xc0dfdD7821c762eF38F86225BD45ff4e912fFA20/0x69C5621db8093ee9a26cc2e253f929316E6E5b92`    
**PAYLOAD**     
```json
{
    "amount":20,
    "fee":0,
    "is_direct":false
}
```
Then Bob check whether there is sufficient fund in delegated charge node :  

```json
{
    "Available": 20,
    "NeedSmt": 0
}
```

#### 4. Bob submit the  proofs of delegation to SmartRaiden Monitoring
**Via API below :**     
`POST: /delegate/<delegater_address>`  

**Example Request :**     
`http://127.0.0.1:6000/delegate/0xf0f6E53d6bbB9Debf35Da6531eC9f1141cd549d5`  

**Example Response :**    
```json
{
    "Status": 3,
    "Error": ""
}
```
- `status = 1 ` - delegation failure    
- `status = 2 ` - delegation success without enough balance    
- `status = 3` - delegation success with enough balance  

#### 5. Bob Goes Offline
Bob gets disconnected from Internet. If Bob has used SmartRaiden Delegation Service, then he has to disconnect before channel settle, otherwise when Alice attempts to close this payment channel, Bob will automatically submit `balance proof` on his own, which can lead to failture that SmartRaiden Monitoring Service can not submit valid `balance proof`.  

#### 6. Alice Closes Payment Channel
 Alice  close the payment channel.  

#### 7. SmartRaiden Monitoring Service waits till half past settletimeout, then update BalanceProof on behalf of delegator.
After our delegator disconnected, if Alice has closed the payment channel, SM Service will monitor the closing events and wait half past the settletimeout, then it will update `BalanceProof` on behalf of the delegator. In this phase, we can query delegation status via API below.  

**Via API below :**      
`GET /tx/<delegater_address>/<channel_address>`    

**Example Request :**     
`GET http://127.0.0.1:6000/tx/0xf0f6E53d6bbB9Debf35Da6531eC9f1141cd549d5/0x2f6418b01422de6cc84fd52e4378fc4449436aeadf69dd543e79e87ee38b6dc8`  

**Example Response :**    

```json
{
    "Status": 3,
    "Error": "",
    "Delegate": {
        "Key": "L2QYsBQi3mzIT9UuQ3j8RElDaurfad1UPnnofuOLbcjw9uU9a7ud6/NdplMeyfEUHNVJ1Q==",
        "Address": "0xf0f6e53d6bbb9debf35da6531ec9f1141cd549d5",
        "PartnerAddress": "0x31ddac67e610c22d19e887fb1937bee3079b56cd",
        "ChannelIdentifier": "L2QYsBQi3mzIT9UuQ3j8RElDaurfad1UPnnofuOLbcg=",
        "OpenBlockNumber": 4392735,
        "SettleBlockNumber": 0,
        "TokenNetworkAddress": "0x445b92d89e21bea3510597a12fefa037718017c1",
        "Time": "2018-09-20T16:42:40.411800993+08:00",
        "TxTime": "0001-01-01T00:00:00Z",
        "TxBlockNumber": 0,
        "MinBlockNumber": 0,
        "MaxBlockNumber": 0,
        "Status": 0,
        "Error": "",
        "Content": {
            "channel_identifier": "0x2f6418b01422de6cc84fd52e4378fc4449436aeadf69dd543e79e87ee38b6dc8",
            "open_block_number": 4392735,
            "token_network_address": "0x445b92d89e21bea3510597a12fefa037718017c1",
            "partner_address": "0x31ddac67e610c22d19e887fb1937bee3079b56cd",
            "update_transfer": {
                "nonce": 2,
                "transfer_amount": 200,
                "locksroot": "0x0000000000000000000000000000000000000000000000000000000000000000",
                "extra_hash": "0xb07638a54d4cbbec93d687d840f13c472df2addae79944dfa00aa2b66ec29b1c",
                "closing_signature": "vwA2lAU414nPajnk5s7dXS/8LHzgCXZ/gVt3lbrzbOhbLkbwcmEW3XtZAjN0qqNYOC8F2opBu7N05FLckrQ5kRw=",
                "non_closing_signature": "aJeOYwRShmkXUBdVTUEKoi5HKKElxn3NfcM8zCRGPcdYJPMv/u+tz9gS3vkb6ypGyvfn1yzKOU1vax3C5/vGMRw=",
                "TxStatus": 0,
                "TxError": "",
                "TxHash": "0x0000000000000000000000000000000000000000000000000000000000000000"
            },
            "unlocks": null,
            "punishes": null
        }
    }
}
```

#### 8. Channel Settle
Alice settle the channel. Once Bob reconnects, he can verify his token amount on chain so that he makes sure SmartRaiden Monitoring Service actually help him update the most recent BalanceProof.  

### Punish Delegation 
This feature is provided to prevent fradulent behaviors, like when one participant goes off, his partner tries to steal tokens via updating abandoned lock.   

**Scenario :**  
There is a transfer lock that has been claimed abandoned. Now, lets change the delegator different from above scenario. Alice  plans to get off the network, and she has submitted proofs to SmartRaiden Monitoring Service, which includes `punish`. Once Alice actually gets disconnected, Bob attempts to unlock the abandoned lock via fraudulent behaviors, to steal the assets of Alice. For the reason that Alice has submitted proofs of `punish`, which has the ability to verify and punish Bob's fraudulent behaviors. As a consequence, channel assets of Bob will be transferred to Alice.  

#### 1. Alice submits proof of fraud to SmartRaiden Monitoring Service
An abandoned transfer lock exists in payment channel of Alice and Bob. Alice submits proof of `punish` to SmartRaiden Monitoring Service.   

**Example Response :**   
```json
{
    "channel_identifier": "0x2f6418b01422de6cc84fd52e4378fc4449436aeadf69dd543e79e87ee38b6dc8",
    "open_block_number": 4399083,
    "token_network_address": "0x445b92d89e21bea3510597a12fefa037718017c1",
    "partner_address": "0xf0f6e53d6bbb9debf35da6531ec9f1141cd549d5",
    "update_transfer": {
        "nonce": 0,
        "transfer_amount": null,
        "locksroot": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "extra_hash": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "closing_signature": null,
        "non_closing_signature": null
    },
    "unlocks": null,
    "punishes": [
        {
            "lock_hash": "0x948443120c7e0c16dac5a0f1df1b3a37c9f19746dd6aca3ad4c28ec201be02c2",
            "additional_hash": "0x6518299834b647647ad4b57dc25ec8bd97b8c3ddd17017c8685065f98b3c0331",
            "signature": "6EemZ91/zqNF7s0wWfgGiPcbm69jOBhmHxfuJ3AQl1wb/NEa7XVrh/o+BDedTBnyA7cubQhssEJ4nHUUE58pMBs="
        }
    ]
}
```

#### 2. Bob attempts to steal tokens 
Once Alice gets off, Bob takes a trial to operate fraudulent behavior, and attempts to steal channel assets of Alice.  
#### 3. SmartRaiden Monitoring Service submits proof of punish
When Bob closes payment channel, and unlock the abandoned lock, our SmartRaiden Monitoring Service will help Alice submit proof of `punish` to the  public chain, which has the ability to approve that the lock which Bob tries to unlock has been abandoned by Alice. After verification check, Bob will get punished and all his channel assets will be transferred to Alice. This process ensures asset security of Alice.  
