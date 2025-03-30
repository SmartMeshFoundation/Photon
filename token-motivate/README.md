# Metalife Server installation instructions and functional description

  Metalife Server is a P2P communication service built using the Scuttlebutt protocol. Its main purpose is to provide remote Peer node interactive connections and to collect like data that the authors publicly released in Pub ,which will report to super nodes for incentives. Since Pub is used as a node for Mediate Transfer, it is necessary to install the corresponding Photon node for payment functions at the same time. The following are installation instructions and related functional descriptions.

## Server Installation  

### Server environment requirements

   Golang version 1.16.6 or higher
  
  The available running memory needs to reach 2G and above, and the available disk space needs to reach 50G and above.

The above configuration is the basic configuration, and the hardware configuration can be further upgraded for Pubs with storage requirements in the future.

### Install the Server

 This version of Server is written in Go. We recommend that users pre-install the latest version of Go. And further execute the following script to install Server.

```bash
git clone https://github.com/MetaLife-Foundation/MetalifePub
cd MetalifePub/cmd/metalifeserver
export GO111MODULE=on
./build.sh
```

 **Instruction**
    Just pull the code on the master branch. This branch supports ebt-mutlifromat by default. It should be noted that a photon node needs to be run on the server where Server is installed, and please enter the parameter configuration file **/restful/params/config.go**. The parameter of PhotonAddress(in line 46)  need to be replaced by the address of the photon node running on the same machine (this account needs to ensure that the number of SMT is sufficient, the pub node will automatically establish the channel with each connected client account, and deposit 0.1smt to the channel)

 **Make sure the  go-sbot and sbotcli are in the $PATH**
 
1.We need to manually create a folder **. ssb-go** on the user's home directory of the server.

```bash
    mkdir -p $HOME/.ssb-go
```
2.create a file named 'secret' and input the ssb-server's private key file,like this:

```bash
{
  "curve": "ed25519",
  "public": "M2nU6wM+F17J0RSLXP05x3Lag8jGv3F3LzHMjh72coE=.ed25519",
  "private": "6t1JnzJz0M4imTUUeoQuYdNnFPcZ78IwwRsjgQN1kMcdmdTrAz4XXsnRFItc/TnHcreDaMa/cXbtMcyOHvZygQ==.ed25519",
  "id": "@M2nU6wM+F17J0RSLXP05x3Lag8jGv3F3LzHMjh72coE=.ed25519"
}
chmod +600 ./secret

```

3.Run(The database will be automatically created after the program runs. The type is sqlite)

```bash
#see: 
# metalifeserver --help
#   --addr value                    tcp address of the sbot to connect to (or listen on) (default: "54.179.3.93:8008")
#   --remoteKey value               the remote pubkey you are connecting to (by default the local key)
#   --datadir value                 directory for storing pub's parsing data (default: "$HOME/.ssb-go/pubdata")
#   --token-address value           which token is used in metalife app,if set,the default will be replaced (default: "0x6601F810eaF2fa749EEa10533Fd4CC23B8C791dc")
#   --photon-host value             host:port link to the photon service. (default: "127.0.0.1:11001")
#   --pub-eth-address value         ethereum address the pub 's address is bound for reward.
#   --settle-timeout value          set settle timeout on photon. (default: 40000)
#   --service-port value            port' for the metalife service to listen on. (default: 10008)
#   --message-scan-interval value   the time interval at which messages are scanned and calculated (unit:second). (default: 60)
#   --min-balance-inchannel value   minimum balance in photon channel between this pub and ssb client (unit: 1e18 wei). (default: 1)
#   --report-rewarding value        pub will reward the person who provides the report (if the report is true). (unit: 1e15 wei) (default: 0)
#   --registration-rewarding value  pub will reward the person who provides ethereum address for his ssb client. (unit: 1e15 wei) (default: 0)
#   --sensitive-words-file value    the path of the sensitive-words file (default: "$HOME/.ssb-go/sensitive.txt")


nohup metalifeserver \
 --pub-eth-address 0xBaBaeafB77585472531D3E8E6f3C3bCF4c04cBE4 \
 --addr 127.0.0.1:8008 \
 --token-address 0x6601F810eaF2fa749EEa10533Fd4CC23B8C791dc \
 --photon-host 127.0.0.1:11001 \
 --settle-timeout 40000 \
 --service-port 10008 \
 --message-scan-interval 120 \
 --min-balance-inchannel 1 \
 --report-rewarding 1 \
 --registration-rewarding 1 \
 > log &
```

The running log will be saved in **log**.

4.Some key operating parameters are located in /MetalifePub/restful/params/config.go

```bash
/MetalifePub/restful/params/config.go
```

### Function Description

**metalifeserver**
The MetaLife server program ,which will provide client access and message log synchronization services that follow the ssb protocol.
The pub message monitoring service provides:

1.Get all ssb-client information on a pub:

```bash
GET http://{ssb-server-public-ip}:18008/ssb/api/node-info

```
Response e.g:
```json
	{
        "error_code": 0,
        "error_message": "SUCCESS",
        "data": [
            {
                "client_id": "@C49GskstTGIrvYPqvTk+Vjyj23tD0wbCSkvX7A4zoHw=.ed25519",
                "client_Name": "beefi",
                "client_alias": "9527",
                "client_bio": "SG",
                "client_eth_address": "0xce92bddda9de3806e4f4b55f47d20ea82973f2d7"
            },
            {
                "client_id": "@eVs235wBX5aRoyUwWyZRbo9r1oZ9a7+V+wEvf+F/MCw=.ed25519",
                "client_Name": "an-Pub1",
                "client_alias": "",
                "client_bio": "SG",
                "client_eth_address": ""
            }
        ]
    }
```

2.Get someone-ssb-client information on a pub:

```bash
GET http://{ssb-server-public-ip}:18008/ssb/api/node-info
```
Body:
```json
{
    "client_id":"@C49GskstTGIrvYPqvTk+Vjyj23tD0wbCSkvX7A4zoHw=.ed25519"
}
```
Response e.g:
```json
	{
        "error_code": 0,
        "error_message": "SUCCESS",
        "data": [
            {
                "client_id": "@C49GskstTGIrvYPqvTk+Vjyj23tD0wbCSkvX7A4zoHw=.ed25519",
                "client_Name": "beefi",
                "client_alias": "",
                "client_bio": "SG",
                "client_eth_address": "0xce92bddda9de3806e4f4b55f47d20ea82973f2d7"
            }
        ]
    }
```

3.The ssb-client register with metalifeserver the ETH address used to receive MetaLife's reward:
(This call is synchronous. It may take a few seconds to return the request result. Pub will establish an Photon Channel  with this client_eth_address)

```bash
Post http://{ssb-server-public-ip}:18008/ssb/api/id2eth
```
Body:
```json
{
    "client_id":"@C49GskstTGIrvYPqvTk+Vjyj23tD0wbCSkvX7A4zoHw=.ed25519",
    "client_eth_address":"0xce92bddda9de3806e4f4b55f47d20ea82973f2d7"
}
```
Response e.g:
```json
{
    "error_code": 0,
    "error_message": "SUCCESS",
    "data": "success"
}
```

4.Get 'Like' Statistics of all on pub

```bash
GET http://{ssb-server-public-ip}:18008/ssb/api/likes
```
Response e.g:
```json
	{
        "error_code": 0,
        "error_message": "SUCCESS",
        "data": {
            "@C49GskstTGIrvYPqvTk+Vjyj23tD0wbCSkvX7A4zoHw=.ed25519": {
                "client_id": "@C49GskstTGIrvYPqvTk+Vjyj23tD0wbCSkvX7A4zoHw=.ed25519",
                "laster_like_num": 7,
                "client_name": "beefi",
                "client_eth_address": "0xce92bddda9de3806e4f4b55f47d20ea82973f2d7",
                "message_from_pub": "@HZnU6wM+F17J0RSLXP05x3Lag2jGv3F3LzHMjh72coE=.ed25519"
            },
            "@eVs235wBX5aRoyUwWyZRbo9r1oZ9a7+V+wEvf+F/MCw=.ed25519": {
                "client_id": "@eVs235wBX5aRoyUwWyZRbo9r1oZ9a7+V+wEvf+F/MCw=.ed25519",
                "laster_like_num": 0,
                "client_name": "an-Pub1",
                "client_eth_address": "",
                "message_from_pub": "@HZnU6wM+F17J0RSLXP05x3Lag2jGv3F3LzHMjh72coE=.ed25519"
            }
        }
    }
```

5.Get 'Like' Statistics of someone-ssb-client on pub


```bash
GET http://{ssb-server-public-ip}:18008/ssb/api/likes
```
Body:
```json
{
    "client_id":"@C49GskstTGIrvYPqvTk+Vjyj23tD0wbCSkvX7A4zoHw=.ed25519"
}
```
Response e.g:
```json
	{
        "error_code": 0,
        "error_message": "SUCCESS",
        "data": {
            "@C49GskstTGIrvYPqvTk+Vjyj23tD0wbCSkvX7A4zoHw=.ed25519": {
                "client_id": "@C49GskstTGIrvYPqvTk+Vjyj23tD0wbCSkvX7A4zoHw=.ed25519",
                "laster_like_num": 7,
                "client_name": "beefi",
                "client_eth_address": "0xce92bddda9de3806e4f4b55f47d20ea82973f2d7",
                "message_from_pub": "@HZnU6wM+F17J0RSLXP05x3Lag2jGv3F3LzHMjh72coE=.ed25519"
            }
        }
    }
```

6.Get pub's profile

```bash
GET http://{ssb-server-public-ip}:18008/ssb/api/pub-whoami
```
Response e.g:
```json
{
        "error_code": 0,
        "error_message": "SUCCESS",
        "data": {
            "pub_id": "@HZnU6wM+F17J0RSLXP05x3Lag2jGv3F3LzHMjh72coE=.ed25519",
            "pub_eth_address": "0xb05Feb81fB4BF6d8B2eB5A5Ae883BAA9E7530cB7"
        }
    }
```

7.The ssb-client post a message tip-off to pub administrator

```bash
POST http://{ssb-server-public-ip}:18008/ssb/api/***
```
Body:
```json
{
    "plaintiff":"@C49GskstTGIrvYPqvTk+Vjyj23tD0wbCSkvX7A4zoHw=.ed25519",  
    "defendant":"@Sg5b3BjZH8XWyJ7mGpH3txrDJmIQtSGxV6MbH6CgeCw=.ed25519",  
    "messagekey":"%w5S3q0eVkTzcfpIKdIR3tJueFTMIOQP1lwcsQkhWSMs=.sha256",   
    "reasons":"sex"                            
}
```
Response e.g:
```json
{
    "error_code": 0,
    "error_message": "SUCCESS",
    "data": "Success, the pub administrator will verify as soon as possible, thank you for your report👍"
}
```

8.Get all message about tip-off, the pub administrator will use this information, this is a combination query

```bash
POST http://{ssb-server-public-ip}:18008/ssb/api/***
```
Body:
```json
{
    "plaintiff":"@C49GskstTGIrvYPqvTk+Vjyj23tD0wbCSkvX7A4zoHw=.ed25519",  
    "defendant":"@Sg5b3BjZH8XWyJ7mGpH3txrDJmIQtSGxV6MbH6CgeCw=.ed25519",  
    "messagekey":"%w5S3q0eVkTzcfpIKdIR3tJueFTMIOQP1lwcsQkhWSMs=.sha256"                        
}
```
Response e.g: (dealtag:0-init, 1-affirm the statement to be true, 2-things didn't turn out like that)
```json
{
    "error_code": 0,
    "error_message": "SUCCESS",
    "data": [
        {
            "plaintiff": "@C49GskstTGIrvYPqvTk+Vjyj23tD0wbCSkvX7A4zoHw=.ed25519",
            "defendant": "@Sg5b3BjZH8XWyJ7mGpH3txrDJmIQtSGxV6MbH6CgeCw=.ed25519",
            "messagekey": "%w5S3q0eVkTzcfpIKdIR3tJueFTMIOQP1lwcsQkhWSMs=.sha256",
            "reasons": "sex",
            "dealtag": "0",
            "recordtime": 1649821765131,
            "dealtime": 1649821765131,
            "dealreward": ""
        }
    ]
}
```

9.the pub administrator handle the data about tip-off

```bash
POST http://{ssb-server-public-ip}:18008/ssb/api/***
```
Body: (dealtag:0-init, 1-affirm the statement to be true, 2-things didn't turn out like that)
```json
{
    "plaintiff":"@C49GskstTGIrvYPqvTk+Vjyj23tD0wbCSkvX7A4zoHw=.ed25519",  
    "defendant":"@Sg5b3BjZH8XWyJ7mGpH3txrDJmIQtSGxV6MbH6CgeCw=.ed25519",  
    "messagekey":"%w5S3q0eVkTzcfpIKdIR3tJueFTMIOQP1lwcsQkhWSMs=.sha256",   
    "dealtag":"1"                            
}
```
Response e.g: 
```json
{
    "error_code": 0,
    "error_message": "SUCCESS",
    "data": "success, [@Sg5b3BjZH8XWyJ7mGpH3txrDJmIQtSGxV6MbH6CgeCw=.ed25519] has been block by [pub administrator], and pub will award token to [@C49GskstTGIrvYPqvTk+Vjyj23tD0wbCSkvX7A4zoHw=.ed25519]"
}
```

10.get the message containing sensitive words preliminarily detected by pub

```bash
POST http://{ssb-server-public-ip}:18008/ssb/api/sensitive-word-events
```
Body: (dealtag:0-init, 1-message contain some sensitive word, 2-things didn't turn out like that)
```json
{
    "deal_tag":"0"
}
```
Response e.g: 
```json
{
    "error_code": 0,
    "error_message": "SUCCESS",
    "data": [
            {
                "pub_id": "@HZnU6wM+F17J0RSLXP05x3Lag2jGv3F3LzHMjh72coE=.ed25519",
                "message_scan_time": 1654676435575,
                "message_text": " Fu +ck",
                "message_key": "%wVdd0jphnuOBOsFQesV+xvoj7SrFcgN3SfEFShsOKrA=.sha256",
                "message_author": "@9I5SiHMp4uEFrev7FyG9G2fgGAamZlqstzjA8OiVY6k=.ed25519",
                "deal_tag": "0",
                "deal_time": 0
            },
            {
                "pub_id": "@HZnU6wM+F17J0RSLXP05x3Lag2jGv3F3LzHMjh72coE=.ed25519",
                "message_scan_time": 1654678260253,
                "message_text": "Fuc::;;k",
                "message_key": "%EbPIJmQT1RI0fIIbNm9InG7K7RqFTOMI4hX58JqgtRY=.sha256",
                "message_author": "@9I5SiHMp4uEFrev7FyG9G2fgGAamZlqstzjA8OiVY6k=.ed25519",
                "deal_tag": "0",
                "deal_time": 0
            }
    ]
}
```
11.the pub administrator makes a judgment on the message suspected of violating the provisions of sensitive words

```bash
POST http://{ssb-server-public-ip}:18008/ssb/api/sensitive-word-deal
```
Body: (dealtag:0-init, 1-message contain some sensitive word, 2-things didn't turn out like that, if dealtag="2", pub will block the author)
```json
{
    "message_key": "%Pz53G0pJENKF7sDShpQQ5Mf6DNIQYPpkf2ZibcVyP0Q=.sha256",
    "deal_tag": "2"
}
```
Response e.g: 
```json
{
    "error_code": 0,
    "error_message": "SUCCESS",
    "data": "success"
}
```
12.ssb client notify pub the login infomation, pub will collect through this interface

```bash
POST http://{ssb-server-public-ip}:18008/ssb/api/notify-login
```
Body: 
```json
{
    "client_id": "@9I5SiHMp4uEFrev7FyG9G2fgGAamZlqstzjA8OiVY6k=.ed25519",
    "login_time": 1654678270241
}
```
Response e.g: 
```json
{
    "error_code": 0,
    "error_message": "SUCCESS",
    "data": "success"
}
```
13.[temporary scheme] notify the pub that user have created a NFT in metalife app

```bash
POST http://{ssb-server-public-ip}:18008/ssb/api/notify-created-nft
```
Body: 
```json
{
    "client_id": "@9I5SiHMp4uEFrev7FyG9G2fgGAamZlqstzjA8OiVY6k=.ed25519",
    "nft_created_time": 1654678310631,
    "nft_tx_hash": "0x909e1b8e2ab3af6435a7f4276c05699e59783fed4214a3e65df8e45ed26f744f",
    "nft_token_id": "1",
    "nft_store_url": "http://106.52.171.12:11111/ipfs/QmW8Ta2dKXZUBVGfVkaiVpTTmyu5JDcijXxpzZYgrYRH1T"
}
```
Response e.g: 
```json
{
    "error_code": 0,
    "error_message": "SUCCESS",
    "data": "success"
}
```
14.get some user daily task infos from pub, a message may appear in multiple pubs, and the client removes redundant data through messagekey and pub id

```bash
POST http://{ssb-server-public-ip}:18008/ssb/api/get-user-daily-task
```
Body: (1-login 2-post 3-comment 4-create NFT)
```json
{
    "author": "@9I5SiHMp4uEFrev7FyG9G2fgGAamZlqstzjA8OiVY6k=.ed25519",
    "message_type": 4,
    "start_time": 1654678319525,
    "start_time": 1654678360051
}
```
Response e.g: 
```json
{
    "error_code": 0,
    "error_message": "SUCCESS",
    "data": [
            {
                "pub_id": "@HZnU6wM+F17J0RSLXP05x3Lag2jGv3F3LzHMjh72coE=.ed25519",
                "author": "@9I5SiHMp4uEFrev7FyG9G2fgGAamZlqstzjA8OiVY6k=.ed25519",
                "message_key": "%wVdd0jphnuOBOsFQesV+xvoj7SrFcgN3SfEFShsOKrA=.sha256",
                "message_type": 4,
                "message_root": "%avPIJmQT1RI0fIIbNm9InG7K7RqFTOMI4hX58Jqgt2g=.sha256",
                "message_time": 1654676435575,
                "nft_tx_hash": "0x909e1b8e2ab3af6435a7f4276c05699e59783fed4214a3e65df8e45ed26f744f",
                "nft_token_id": "1",
                "nft_store_url": "http://106.52.171.12:11111/ipfs/QmW8Ta2dKXZUBVGfVkaiVpTTmyu5JDcijXxpzZYgrYRH1T"
            }
    ]
}
```

15.Get the statistics of users' likes, if the body parms is null, you will get all users' likes

```bash
POST http://{ssb-server-public-ip}:18008/ssb/api/set-like-info
```
Body: (or null)
```json
{
    "client_id":"@P2AR780TWII9tJXYfarlqAlU74hcU11XQ6ZdkPuv19A=.ed25519"
}
```
Response e.g: 
```json
{
    "error_code": 0,
    "error_message": "SUCCESS",
    "data": {
        "@P2AR780TWII9tJXYfarlqAlU74hcU11XQ6ZdkPuv19A=.ed25519": {
            "client_id": "@P2AR780TWII9tJXYfarlqAlU74hcU11XQ6ZdkPuv19A=.ed25519",
            "laster_like_num": 2,
            "client_name": "",
            "client_eth_address": "",
            "message_from_pub": "@HZnU6wM+F17J0RSLXP05x3Lag2jGv3F3LzHMjh72coE=.ed25519"
        }
    }
}
```

16.get all or someones' reward information order by PUB RULE

```bash
POST http://{ssb-server-public-ip}:18008/ssb/api/get-reward-info
```
Body: (or null)
```json
{
    "client_id":"@P2AR780TWII9tJXYfarlqAlU74hcU11XQ6ZdkPuv19A=.ed25519",
    "time_from":0,
    "time_to":1656863999000
}
```
Response e.g: 
```json
{
    "error_code": 0,
    "error_message": "SUCCESS",
    "data": [
        {
            "client_id": "@P2AR780TWII9tJXYfarlqAlU74hcU11XQ6ZdkPuv19A=.ed25519",
            "client_eth_address": "0xea753b41854D37bAD352Cd7464F104421d325BD1",
            "grant_success": "success",
            "grant_token_amount": 1000000000000000000,
            "reward_reason": "like a post",
            "message_key": "%O7FkxpIFdX2W8RF3YpBmQBsQ26K2k376vavvwIsDRWk=.sha256",
            "message_time": 1656800926053,
            "reward_time": 1656800982331
        },
        {
            "client_id": "@P2AR780TWII9tJXYfarlqAlU74hcU11XQ6ZdkPuv19A=.ed25519",
            "client_eth_address": "0xea753b41854D37bAD352Cd7464F104421d325BD1",
            "grant_success": "success",
            "grant_token_amount": 10000000000000000000,
            "reward_reason": "daily login",
            "message_key": "",
            "message_time": 1656801371198,
            "reward_time": 1656801433901
        },
        {
            "client_id": "@P2AR780TWII9tJXYfarlqAlU74hcU11XQ6ZdkPuv19A=.ed25519",
            "client_eth_address": "0xea753b41854D37bAD352Cd7464F104421d325BD1",
            "grant_success": "success",
            "grant_token_amount": 5000000000000000000,
            "reward_reason": "post message",
            "message_key": "%7TNo6zaiYsYQgpB5E3cIvvV21XeRRMd6qaDP6+xsfw4=.sha256",
            "message_time": 1656801880962,
            "reward_time": 1656801991149
        },
        {
            "client_id": "@P2AR780TWII9tJXYfarlqAlU74hcU11XQ6ZdkPuv19A=.ed25519",
            "client_eth_address": "0xea753b41854D37bAD352Cd7464F104421d325BD1",
            "grant_success": "success",
            "grant_token_amount": 2000000000000000000,
            "reward_reason": "post comment",
            "message_key": "%kEATQayHCKeq9kQlSFq4QB+DZSjgQxebeV5U+5/VBRo=.sha256",
            "message_time": 1656802265942,
            "reward_time": 1656802480744
        },
        {
            "client_id": "@P2AR780TWII9tJXYfarlqAlU74hcU11XQ6ZdkPuv19A=.ed25519",
            "client_eth_address": "0xea753b41854D37bAD352Cd7464F104421d325BD1",
            "grant_success": "success",
            "grant_token_amount": 10000000000000000000,
            "reward_reason": "mint a nft",
            "message_key": "",
            "message_time": 1656802714262,
            "reward_time": 1656802774163
        }
    ]
}
```

3.Channel establishment and pre-deposit service  
After receiving the ETH address registration message, the  MetaLife server will actively establish a channel with the client to obtain rewards , on Spectrum Main Chain.

**Notice**Here is the spectrum ,TokenAddress=”0x6601F810eaF2fa749EEa10533Fd4CC23B8C791dc”

4.Extension: The monitoring service has mastered all the likes and other types of messages, and will provide details of the like-source of the specific likes in the follow-up.
