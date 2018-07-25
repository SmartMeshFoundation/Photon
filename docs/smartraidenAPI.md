 Smartraiden API使用文档
=========

## 前言
Smartraiden可以通过访问URL端点来执行API进行通道操作。端点接受和返回JSON编码的对象。 API URL路径中包含了API版本。所有的查询以/api/<version>/ 开始，<version>是一个整数，表示当前API版本。参与Smartraiden Token网络进行转账，可以通过一些必要的步骤及不同的场景。如加入一个已经存在的token网络、注册一个新的toke网络、打开通道、关闭通道、结算通道等。  以下为smartfaiden的相关场景及API使用。
## 场景
以下为一系列不同的场景，用户可以通过Smart raiden API与端点进行交互。
### 1.自启动一个token网络  
首先假定用户持有某种*ERC20 token*还没有注册到Smart raiden网络。用户可以建立一个通道管理注册这个*token*，每一个注册的token有一个相应的通道管理。通道管理负责在两个节点之间打开新的支付通道（后面注册token中具体实现）。
首先，对于一个节点来说，需要知道自己地址。smartraiden开始选择一个以太坊地址，也是smartraiden地址，你可以通过一个`Get`请求到 `/api/<version>/address`的端点查询地址。
```
GET /api/1/address
```
返回**200 OK** 以及
```json
{"our_address": "0x2a65aca4d5fc5b5c859090a6c34d164135398226"}
```
本文档预设7个节点。
示例：
```
GET http://{{ip2}}/api/1/address
```
返回结果：
```json
{"our_address": "0x33Df901ABc22DcB7F33c2a77aD43CC98FbFa0790"}
```
通过查询获得所有节点地址：
  * 节点1："0x1a9eC3b0b807464e6D3398a59d6b0a369Bf422fA"
  * 节点2："0x33Df901ABc22DcB7F33c2a77aD43CC98FbFa0790"
  * 节点3："0x8c1b2E9e838e2Bf510eC7Ff49CC607b718Ce8401" 
  * 节点4："0xc4c08f9227BE0F1750F5D5467EeD462Ec133B15E"
  * 节点5："0x215c0D259AC31571a43295f2E411A697CD30748C"
  * 节点6："0x543Fc024CDD1F0d346a306f5E99ec0D8FE392920"
  * 节点7："0x920A90ACC9164272Ede4Ae1E9C33841F019f53A4"   
### 2.检查token是否注册
 可以通过查询得到所有注册*token*列表来检查*token*是否已经注册。如果列表中存在想要交互的*token*地址，则*token*已注册。
 ```
 GET /api/1/tokens
 ```
 如果token地址已经存在，则可以进行通道操作。如果没有，需要先注册这个token。
示例：
```
GET  http://{{ip1}}/api/1/tokens
```
返回结果：
```
["0x590511b52a46384e1bacf29fa937c6332fe60858",

"0xdc7ff7683c883f3ebd6cb8814583017e50df280d",

"0x65494231d3046df617b75270c817e77c6f15bcee",

"0x1986e04955c9b76e1b19ddd77e94782c2f12c81a",

"0x58e90b62f518dfd49fa70d0c16a48cc3f6157c26",

"0x13800624f22ecc4e3e22afff98c315892aa32db8",

"0x0b51cef630c850a3cd72d673015752a0c191b63a"
]
```
### 3.注册一个token
为了注册一个token,只需要这个token的地址。当新token注册时，一个通道管理合约被部署。
通过一个PUT请求到端点
 `/api/<version>/tokens/<token_address>`
这个请求将返回部署的通道管理地址
`PUT /api/1/tokens/0xea674fdde714fd979de3edf0f56aa9716b898ec8`
**201 Created**
```json
"channel_manager_address": "0xC4F8393fb7971E8B299bC1b302F85BfFB3a1275a"
```
示例：
```
PUT http://{{ip1}}/api/1/tokens/0x6112cd0e03e0fb88f0bd5ad6c64355d391c1fcfd
```
如果成功，这个调用将返回新建立的通道管理地址：
```json
{"channel_manager_address": "0xC4F8393fb7971E8B299bC1b302F85BfFB3a1275a"}
```
此时，token已经注册。但此时，token刚注册，还没有其他节点连接到token网络，因此没有节点相连。因此，这个特定的token网络需要自启动。如果持有这个token的其他smartraiden节点地址已知，或者想单向转账给另一个smartraiden节点，可以通过向这个节点打开一个通道来实现（自启动）。不管这个节点是否持有这个token，都可以打开通道。
### 4.打开一个通道
为与另一个smartraiden节点打开一个通道，有四个东西要需要：token地址，对方节点地址，想要存款的token数量，结算超时的时间。
```
PUT /api/1/channels
```
示例：
```
PUT http://{{ip1}}/api/1/channels
```
请求参数：
```json
{
"partner_address": "0x61c808d82a3ac53231750dadc13c777b59310bd9",
"token_address": "0x9aBa529db3FF2D8409A1da4C9eB148879b046700",
"balance": 1337,
"settle_timeout": 600
}
```
响应:
```json
{
"channel_address": "0x2a65aca4d5fc5b5c859090a6c34d164135398226",
"partner_address": "0x61c808d82a3ac53231750dadc13c777b59310bd9",
"token_address": "0x9aBa529db3FF2D8409A1da4C9eB148879b046700",
"balance": 1337,
"state": "opened",
"settle_timeout": 600
}
```
注意到通道地址已经建立。这意味着一个净通道合约已经部署到区块链上。也表示两个节点间的一个特定token支付通道地址。详细如下：
```
PUT /api/<version>/channels
```
```
PUT /api/1/channel
```
请求参数：
```json
{
"partner_address": "0x61c808d82a3ac53231750dadc13c777b59310bd9",
"token_address": "0xea674fdde714fd979de3edf0f56aa9716b898ec8",
"balance": 35000000,
"settle_timeout": 100
}
```
响应：
**201 Created**
```json
{
"channel_address": "0x2a65aca4d5fc5b5c859090a6c34d164135398226",
"partner_address": "0x61c808d82a3ac53231750dadc13c777b59310bd9",
"token_address": "0xea674fdde714fd979de3edf0f56aa9716b898ec8",
"balance": 35000000,
"state": "open",
"settle_timeout": 100,
"reveal_timeout": 30
}
```
### 5.存钱进一个通道
支付通道打开后，因为只有一个节点有token，所以只能这个节点向对方转账。可以通知另一个节点有一个通道已经向他打开，他也可以存钱进这个通道。
```
PATCH/api/1/channels/0x2a65aca4d5fc5b5c859090a6c34d164135398226
```
示例：
```
PATCH http://{{ip1}}/api/1/channels/0x09468c9F787dD4316aE6404bb86B9F1Ac6E501E3
```
请求参数：
```json
{
"balance": 7331
}
```
当对方存了token,通道可以进行查询
```
GET /api/1/events/channels/0x2a65aca4d5fc5b5c859090a6c34d164135398226?
```
**from_block:** 表示从哪个区块开始
```
GET /api/1/events/channels/0x2a65aca4d5fc5b5c859090a6c34d164135398226?from_block=1337
```
返回结果：
```json
{
"event_type": "ChannelNewBalance",
"participant": "0x61c808d82a3ac53231750dadc13c777b59310bd9",
"balance": 7331,
"block_number": 54388
}
```
从上面的事件可以得出对方存钱进了通道。节点双方也可以查询特定通道的状态：
```
GET /api/1/channels/0x2a65aca4d5fc5b5c859090a6c34d164135398226
```
返回结果与打开通道类似。
示例：
```
GET http://{{ip1}}/api/1/events/channels/0x09468c9F787dD4316aE6404bb86B9F1Ac6E501E3?from_block=1
```
```json
[
{
"balance": 300,
"block_number": 232223,
"event_type": "ChannelNewBalance",
"participant": "0x1a9eC3b0b807464e6D3398a59d6b0a369Bf422fA"
},
]
```
完整描述如下：
```
PATCH /api/<version>/channels/<channel_address>
```
```
PATCH /api/1/channels/0x2a65aca4d5fc5b5c859090a6c34d164135398226
```
请求参数：
```
{"balance": 100}
```
返回结果：
**200OK**
```json
{
"channel_address": "0x2a65aca4d5fc5b5c859090a6c34d164135398226",
"partner_address": "0x61c808d82a3ac53231750dadc13c777b59310bd9",
"token_address": "0xea674fdde714fd979de3edf0f56aa9716b898ec8",
"balance": 100,
"state": "open",
"settle_timeout": 100
}
```
### 6.**加入一个已经存在的****token****网络**
在上面的场景中，已经展示了对一个未注册的token如何注册和自启动一个token网络。在这一节，介绍最普通的加入一个token网络方式。大部分情况，用户不需要建立一个新的token网络，如果他们经持有的ERC20token，他们可以加入一个已经存在的token网络。
* 1.  连接到一个token网络

连接到一个已经存在的token网络很简单。所需要的是，要加入的token网络地址，相应的打算存进通道里的token数量。
```
PUT /api/1/connections/0xc9d55C7bbd80C0c2AEd865e9CA13D015096ce671
```
body：
```json
{
"funds": 2000
}
```
示例：
```
PUT http://{{ip1}}/api/1/connections/0xf1b0964f1e19ecf07ddd3bd8e20138c82680395d
```
body:
```json
{
"funds": 1337
}
```
这将自动连接并打开通道和三个随机的对等方在token网络里, 每个通道存20%资金,40%未分配, 允许新的节点加入这个网络和这个节点打开双向资金支付通道。对这个AETtoken，用户节点现在连接到了token网络，将有一个路径到所有已经加入这个toke网络的其他节点，因此它可以传输token给所有参与这个网络的节点。
*  2.  离开一个token网络 

如果一个节点想离开一个token网络，可以通过以下方式：
```
DELETE /api/1/connections/0xc9d55C7bbd80C0c2AEd865e9CA13D015096ce671
```
示例：
```
DELETE http://{{ip1}}/api/1/connections/0xcCD137FF778083B0C32737Da6BB2eDAc3c3Ba98E
```
对token网络里的特定节点关闭和结算所有打开的通道。
根据  `settlement_timeout`  ，调用将花费一段时间。Leave只关闭和结算已经收到转账的通道。

### 7.转账token
假定一个节点用AET token连接到通道网络，在这种情况下，节点连接到5个节点。
转账token给另一个节点很容易，想转账Token的地址已知，如  `0xc9d55C7bbd80C0c2AEd865e9CA13D015096ce671`.剩下的是需要想转账的目标节点地址。假定如下：`0x61c808d82a3ac53231750dadc13c777b59310bd9`
```
POST /api/<version>/transfers/<token_address>/<target_address>
```
这个请求只返回转账或者成功或者失败。转账失败可能由于**时间锁过期**，**目标下线**，路由上的通道到目的地**没有足够的** **settle_timeout** 和**reveal_timeout** 。
```
POST /api/1/transfers/0xc9d55C7bbd80C0c2AEd865e9CA13D015096ce671/0x61c808d82a3ac53231750dadc13c777b59310bd9
```
转账数量记录在负载里:
```json
{
"amount": 42
}
```
`"identifier": some_integer` 也可以加入负载，是可选的，为了提供一种方式标记转账。
如果有一条网络路径有足够的容量，这个地址发送转账的有足够的token，这个转账将成功。
响应：
**200 OK** with payload
```json
{
"initiator_address": "0xea674fdde714fd979de3edf0f56aa9716b898ec8",
"target_address": "0x61c808d82a3ac53231750dadc13c777b59310bd9",
"token_address": "0x2a65aca4d5fc5b5c859090a6c34d164135398226",
"amount": 200,
"identifier": 42
}
```
接收节点通过查询所有他打开的通道将能够看到进入的转账。
```
GET /api/1/events/channels/_0x000397DFD32aFAAE870E6b5FB44154FD43e43224_?from_block=1337
```
将返回事件列表。可以通过过滤进入转账的列表得到。请注意smartraiden最有力的特征是用户可以发送转账到任何连接到网络的节点，只要有路径，有足够的容量，不需要用户直接相连。称之为 mediated transfers。
在进行mediated transfers时，需要考虑路由情况。以本文预设的7个节点为例。

![smartraiden](/docs/images/smartraidenAPI.png)
#### 场景一：转账给直接通道对方
节点A同时向节点B和节点C、节点D进行转账，假定当前余额足够。
示例：节点2向节点1，节点3，节点4，节点7同时转账20 token，由于节点2与1、3、4、7均有直接通道，因此，转账不通过中转节点。转账成功。 
#### 场景二：给没有直接通道的对手方转账
一个中间节点：节点1向节点3转账，需要通过节点2，此时，需要节点2 有足够余额，settletimeout时间足够，否则会发生转账失败，锁定转账的token.多个中间节点：节点2向节点6转账：
此时，由于拓扑图可知：
- 2-3-6
- 2-4--5-6
- 2-7-3-6
根据最短路径，转账选择路由为2-3-6。可通过中间节点的通道查看余额变化情况。
#### 场景三： 无路由转账
在节点6上用token 7打开一个通道，建立与节点5的连接。由于没有路由（直接或间接)，返回结果
```
{	
	"Error": "no available route"
}
```
#### 场景四：连续转账
节点1连续向节点2多次转账，此时，如果可用余额足够，则转账成功，看转账结果。
#### 场景五：余额不足转账
余额不足的时候，报错，转账失败。
```
{
"Error": "no available route"
}
```
#### 场景六：中间节点不在线转账
节点1给节点6转账，通过1，2，3，6 或者1，2，4，5，6 或者1，2，7，3，6
如果此时3节点不在线，转账不成功，锁定转账金额。原因是按照最短路线，首选1，2，3，6。**只有3在线，告诉2不行，2才回去尝试4，3不在线，就转不了账**。
#### 场景七：中间节点余额不足转账
此时使用refundtransfer
![smartraiden](/docs/images/smartraidenAPI.png)
节点1向节点6转账 45token。
如果此前，节点3已向节点6转账270 token,向节点7转账30token,显然，节点3收到转账请求后，发现123 6 走不通，退回2，走12 7 3 6，也走不通，最后选择走12456，成功转账。由于refundtransfer时间较长，此时，需要确定settletimeout时间是否足够，否则很容易失败，锁定转账金额。
#### 场景八：奔溃恢复后转账
此场景由于崩溃情况众多，独立进行说明。
### 8.关闭通道
如果在任何时刻想要关闭一个特定的通道，可以用close。
```
PATCH /api/1/channels/0x000397DFD32aFAAE870E6b5FB44154FD43e43224
```
负载:
```jsson
{
"state":"closed"
}
```
关闭成功，返回状态为closed:
```
{
"channel_address": "0x000397DFD32aFAAE870E6b5FB44154FD43e43224",
"partner_address": "0x61c808d82a3ac53231750dadc13c777b59310bd9",
"token_address": "0xc9d55C7bbd80C0c2AEd865e9CA13D015096ce671",
"balance": 350,
"state": "closed",
"settle_timeout": 600
}
```
状态由opened变为closed.
详细如下：
```PATCH /api/1/channels/0x2a65aca4d5fc5b5c859090a6c34d164135398226```
请求参数：
```json
{
"state":"closed"
}
```
返回结果：
**200 OK**
```json
{
"channel_address": "0x2a65aca4d5fc5b5c859090a6c34d164135398226",
"partner_address": "0x61c808d82a3ac53231750dadc13c777b59310bd9",
"token_address": "0xea674fdde714fd979de3edf0f56aa9716b898ec8",
"balance": 35000000,
"state": "closed",
"settle_timeout": 100
}
```
### 9.结算通道
一旦关闭通道调用，settle 超时开始计算。在这个期间，对方节点提供最新接收的消息。当结算超时过后，通道最终被结算。
```
PATCH  /api/1/channels/0x000397DFD32aFAAE870E6b5FB44154FD43e43224
```
负载：
```
{
"state":"settled"
}
```
 这时，将触发settle()函数。一旦结算成功，将返回：
 ```json
 {
"channel_address": "0x000397DFD32aFAAE870E6b5FB44154FD43e43224",
"partner_address": "0x61c808d82a3ac53231750dadc13c777b59310bd9",
"token_address": "0xc9d55C7bbd80C0c2AEd865e9CA13D015096ce671",
"balance": 0,
"state": "settled",
"settle_timeout": 600
}
```
此时，通道余额为0，状态变成settled.  意味着两个参与者的余额已经被转移到区块链，支付通道的生命周期结束。这时，区块链合约已经自毁.
完整描述如下：
```
PATCH /api/<version>/channels/<channel_address>
```
```
PATCH /api/1/channels/0x2a65aca4d5fc5b5c859090a6c34d164135398226
```
请求参数：
```
{"state":"settled"}
```
返回结果：
**200 OK**
```json
{
"channel_address": "0x2a65aca4d5fc5b5c859090a6c34d164135398226",
"partner_address": "0x61c808d82a3ac53231750dadc13c777b59310bd9",
"token_address": "0xea674fdde714fd979de3edf0f56aa9716b898ec8",
"balance": 0,
"state": "settled",
"settle_timeout": 100
}
```
### 10.token swap
token互换允许Alice和Bob交换两种不同的token。这意味着如果Alice和Bob参与tokenA和tokenB网络，那么他们能够原子的交换一些数量的tokenA对一些数量的tokenB。假定Alice想要10个tokenB换Bob 2个tokenA，如果Bob同意这个条款，互换能够被执行. 在上面例子情况下，Alice是maker,Bob是taker.
```
PUT/api/1/token_swaps/0x61c808d82a3ac53231750dadc13c777b59310bd9/1337
```
负载：
```json
{
"role": "maker",
"sending_amount":42,
"sending_token": "0xea674fdde714fd979de3edf0f56aa9716b898ec8",
"receiving_amount": 76,
"receiving_token": "0x2a65aca4d5fc5b5c859090a6c34d164135398226"
}
```
这里有一些有趣的参数：role定义报文地址的角色是maker还是taker。taker调用必须被执行在maker调用被执行之前，  sending_amount  和  sending_token  表示maker想发送的token数量换回receiving_token  和  receiving_amount。上面的例子，Alice想用42个tokenA(  0xea674fdde714fd979de3edf0f56aa9716b898ec8),换Bob76个tokenB(0x2a65aca4d5fc5b5c859090a6c34d164135398226).  现在需要Bob接收这个offer,  因为Alice和Bob私下交换这个swap，因此，Alice简单告诉Bob的 identifier。Bob能take这个offer。
```
PUT /api/1/token_swaps/0xbbc5ee8be95683983df67260b0ab033c237bde60/1337
```
这里地址是Alice的地址，注意identifier在请求里相同的。负载如下：
```json
{
"role": "taker",
"sending_amount": 76,
"sending_token": "0x2a65aca4d5fc5b5c859090a6c34d164135398226",
"receiving_amount": 42,
"receiving_token": "0xea674fdde714fd979de3edf0f56aa9716b898ec8"
}
```
这里角色从maker变成了taker.进一步，发送和接收参数已经颠倒了。这是因为从Bob的观点，swap已经看见了。此时此刻，Alice和Bob的余额在互换后将反映状态。
节点1和节点2在token5 和token7上各有余额，按照4：12的比例，进行token swap
节点1是maker,节点2是taker,
节点1在通道0x09468c9F787dD4316aE6404bb86B9F1Ac6E501E3上  用0x0b51CEf630c850A3CD72d673015752A0C191B63A （t5）发送4个token
节点2在通道0x26e3b56A9A121cEbB514Ff50015447ba435d3e83上用0x58e90B62F518dfD49Fa70d0C16A48CC3f6157C26（t7）发送12个token
通过token swap后，成功实现token互换。执行顺序为先执行taker,后执行maker,此时，一定要记住：标识一定要相同，否则转账会失败。另外，如果，只进行maker或taker，因为不能进行原子操作，则转账不成功，锁定转账金额，等超时。
maker:
```
PUT /api/1/token_swaps/0x61c808d82a3ac53231750dadc13c777b59310bd9/1337
```
```json
{
"role": "maker",
"sending_amount": 42,
"sending_token": "0xea674fdde714fd979de3edf0f56aa9716b898ec8",
"receiving_amount": 76,
"receiving_token": "0x2a65aca4d5fc5b5c859090a6c34d164135398226"
}
```
taker:
```
PUT /api/1/token_swaps/0xbbc5ee8be95683983df67260b0ab033c237bde60/1337
```
```json
{
"role": "taker",
"sending_amount": 76,
"sending_token": "0x2a65aca4d5fc5b5c859090a6c34d164135398226",
"receiving_amount": 42,
"receiving_token": "0xea674fdde714fd979de3edf0f56aa9716b898ec8"
}
```
**201 CREATED**
### 11.查询关于通道和token的信息
查询一个特定的通道
```
GET /api/1/channels/0x2a65aca4d5fc5b5c859090a6c34d164135398226
```
返回结果:
```json
{
"channel_address": "0x2a65aca4d5fc5b5c859090a6c34d164135398226",
"partner_address": "0x61c808d82a3ac53231750dadc13c777b59310bd9",
"token_address": "0xea674fdde714fd979de3edf0f56aa9716b898ec8",
"balance": 35000000,
"state": "open",
"settle_timeout": 100
}
```
查询所有通道
```
GET /api/1/channels
```
**200 OK**
```JSON
[
	{
	"channel_address": "0x2a65aca4d5fc5b5c859090a6c34d16413539822",
	"partner_address": "0x61c808d82a3ac53231750dadc13c777b59310bd9",
	"token_address":  "0xea674fdde714fd979de3edf0f56aa9716b898ec8",
	"balance": 35000000,
	"state": "open",
	"settle_timeout": 100
	}
]
```
查询所有注册的tokens
```
GET /api/1/tokens
```
**200 OK**
```json
[

"0xea674fdde714fd979de3edf0f56aa9716b898ec8",

"0x61bb630d3b2e8eda0fc1d50f9f958ec02e3969f6"

]
```
查询一个token所有的对方
```
GET /api/1/tokens/0x61bb630d3b2e8eda0fc1d50f9f958ec02e3969f6/partners
```
**200 OK**
```json
[

{

"partner_address":"0x61c808d82a3ac53231750dadc13c777b59310bd9",

"channel": "/api/<version>/channels/0x2a65aca4d5fc5b5c859090a6c34d164135398226"

}

]
```
 ### 12.连接管理
- 连接到一个token网络

你可以自动加入一个token网络通过一个PUT请求到下面的端点用一个JSON负载，包含连接细节如资金你打算放到网络里，初始通道目标数量业建立和目标可加入的资金。
```
PUT /api/<version>/connections/<token_address>
```
```
PUT /api/1/connections/0x2a65aca4d5fc5b5c859090a6c34d164135398226
```
请求参数：
```json
{
"funds": 1337
}
```
返回结果：
**204 NO CONTENT**
- 离开一个token网络

你可以离开一个token网络通过制造一个delete请求到以下的端点随着一个JSON负载，包含细节关于你想离开的网络。
```json
DELETE /api/<version>/connections/<token_address>
```
这个请求将返回，一旦所有区块链调用关闭/结算一个通道已经完成。
```
DELETE /api/1/connections/0x2a65aca4d5fc5b5c859090a6c34d164135398226
```
请求参数：
```json
{
"only_receiving_channels": false
}
```
返回结果：
**200 OK**
```json
[

"0x41bcbc2fd72a731bcc136cf6f7442e9c19e9f313",

"0x5a5f458f6c1a034930e45dc9a64b99d7def06d7e",

"0x8942c06faa74cebff7d55b79f9989adfc85c6b85"

]
```
- 查询连接细节

你可以查询以前加入的token网络细节通过制造一个get请求到一个连接端点。
```
GET /api/<version>/connections
```
这个请求将返回一个JSON对象，每个键是一个你打开的通道的token地址。这个值是一个JSON对象包含数值值funds从上一个连接请求。sum_deposits 是所有当前打开的通道和通道数对那个token.
```
GET /api/1/connections
```
**200 OK**
```json
{

"0x2a65aca4d5fc5b5c859090a6c34d164135398226": {

"funds": 100,

"sum_deposits": 67,

"channels": 3

},

"0x0f114a1e9db192502e7856309cc899952b3db1ed": {

"funds": 49

"sum_deposits": 31,

"channels": 1

}

}
```
### 13.查询事件
- 查询整个网络事件
```
GET /api/<version>/events/network
```
网络注册是默认的注册。默认的注册是预先设定的能被雷电配置文件编辑。你可以查询注册网络事件通过一个GET请求。
```
GET /api/1/events/network
```
**200 OK**
```json
[

{

"event_type": "TokenAdded",

"token_address": "0xea674fdde714fd979de3edf0f56aa9716b898ec8",

"channel_manager_address": "0xc0ea08a2d404d3172d2add29a45be56da40e2949"

}, {

"event_type": "TokenAdded",

"token_address": "0x91337a300e0361bddb2e377dd4e88ccb7796663d",

"channel_manager_address": "0xc0ea08a2d404d3172d2add29a45be56da40e2949"

}

]
```
- 查询token网络事件
```
GET /api/<version>/events/tokens/<token_address>
```
查询为一个token对所有打开的新通道
```
GET /api/1/events/tokens/0x61c808d82a3ac53231750dadc13c777b59310bd9
```
返回结果：
**200 OK** 
```json
[

	{

	"event_type": "ChannelNew",

	"settle_timeout": 10,

	"netting_channel": "0xc0ea08a2d404d3172d2add29a45be56da40e2949",

	"participant1": "0x4894a542053248e0c504e3def2048c08f73e1ca6",

	"participant2": "0x356857Cd22CBEFccDa4e96AF13b408623473237A"

	}, {

	"event_type": "ChannelNew",

	"settle_timeout": 15,

	"netting_channel": "0x61c808d82a3ac53231750dadc13c777b59310bd9",

	"participant1": "0xea674fdde714fd979de3edf0f56aa9716b898ec8",

	"participant2": "0xc7262f1447fcb2f75ab14b2a28deed6006eea95b"

	}

]
```
- 查询通道网络事件
```
GET /api/<version>/events/channels/<channel_registry_address>
```
你可以查询事件联系于一个特定的通道，通过产生一个get请求  对事件端点用他的地址。
```
GET /api/1/events/channels/0x2a65aca4d5fc5b5c859090a6c34d164135398226?from_block=1337
```
返回结果：
**200 OK**
```json
[

	{

	"event_type": "ChannelNewBalance",

	"participant": "0xea674fdde714fd979de3edf0f56aa9716b898ec8",

	"balance": 150000,

	"block_number": 54388

	}, {

	"event_type": "TransferUpdated",

	"token_address": "0x91337a300e0361bddb2e377dd4e88ccb7796663d",

	"channel_manager_address": "0xc0ea08a2d404d3172d2add29a45be56da40e2949"

	}, {

	"event_type": "EventTransferSentSuccess",

	"identifier": 14909067296492875713,

	"block_number": 2226,

	"amount": 7,

	"target": "0xc7262f1447fcb2f75ab14b2a28deed6006eea95b"

	}

]
```
### 14.查询token余额
为了查询通道对方`token`中还剩下的余额，可以通过Get请求到对方地址。
```http
http://{{ip2}}/api/1/debug/balance/0x0b51CEf630c850A3CD72d673015752A0C191B63A/0x543Fc024CDD1F0d346a306f5E99ec0D8FE392920
```
显示余额：`4999870`
Ip4新建通道，"0x60b4DF99906D81CDB9a1C75C120F0a4504618011",存款120
查询得余额5000000
现在节点5进行存款，120后，再次查询余额，4999880 settled 后，恢复5000000。
```http
http://{{ip3}}/api/1/debug/balance/0x58e90b62f518dfd49fa70d0c16a48cc3f6157c26/0xc4c08f9227BE0F1750F5D5467EeD462Ec133B15E
```
显示余额：`4999910`
创建通道并转账后
```json
{

"initiator_address": "0x8c1b2E9e838e2Bf510eC7Ff49CC607b718Ce8401",

"target_address": "0xc4c08f9227BE0F1750F5D5467EeD462Ec133B15E",

"token_address": "0x58e90b62f518dfd49fa70d0c16a48cc3f6157c26",

"amount": 120,

"identifier": 4080314417071621426,

"fee": 0

}
```
关闭并settled后，链上余额变为：`5000030`
对方这个token剩下可用的余额。(某个账户有多少token，与雷电网络无关)
### 15.查询转账的token
为了查询你转给通道对方对应token转账金额（你想转多少给对方），可以通过GET请求到对方地址。
```http
http://{{ip2}}/api/1/debug/transfer/0x0b51CEf630c850A3CD72d673015752A0C191B63A/0x543Fc024CDD1F0d346a306f5E99ec0D8FE392920/45
```
返回结果：
**OK**
（你想转多少token给谁）

### 16.崩溃恢复
![smartraiden](/docs/images/smartraidenAPI.png)
#### 1.发送消息崩溃
- 场景一：EventSendMediatedTransferAfter

发送中转转账后崩溃
节点1向节点2发送MTR后，节点1崩溃，此时，节点2默认收到MTR，但由于没有ACK确认，没发生转账，余额不变。节点2没收到转账token.
重启节点1后，继续转账，转账成功。 
- 场景二：EventSendRevealSecretAfter

节点2向节点6转账20token,发送revealsecret后，节点2崩，路由走2-3-6，查询节点6，节点3，交易未完成，锁定节点3 20个token,节点2 20个token；重启节点2后，锁定的token解锁，节点3和节点6的交易完成，节点2和节点3交易完成。转账继续完成。


- 场景三：EventSendBalanceProofAfter

发送余额证明后崩溃（发送方崩）

节点2向节点6转账20 token,发送balanceProof后，节点2崩，路由走2-3-6，查询节点3，节点6，节点3和6之间交易完成。节点2、3交易未完成，节点2锁定20token。重启节点2后，节点2、3交易完成，实现转账继续。
- 场景四：EventSendSecretRequestAfter  

发送Secretrequest后崩溃

节点2向节点6转账20 token,节点6发送Secretrequest后，节点6崩。查询节点2，节点3，节点2锁定20 token,节点3锁定20token,交易未完成。重启节点6后，交易完成，实现转账继续。
- 场景五：EventSendRefundTransferAfter

发送refundtransfer交易崩溃

节点2发送45token给节点6 ，发送refundtransfer后节点3崩，节点2锁定45token，其余节点无锁定;重启节点3后，节点2，3各锁定 45，节点2，4、节点4、5，节点5、6交易成功，转账成功，但节点2、3各锁定45token.存在问题：锁定45个token未解锁。虽然转账完成，但交易未完成。正常。

（如果选择节点7，则转账失败）

说明：refundtransfer成功率与选择路由有一定关系，如果路由回退过程中选择了一条余额不足的新路，则转账失败。  由于最短路径算法的局限性，造成在节点7赌死，转账失败。

选择节点1发送相同的转账，结果崩溃恢复后按12456路线走，转账成功。
存在问题：路由选择会导致转账失败
#### 2. 收到消息崩溃
- 场景一：ActionInitTargetStateChange

收到mtr后崩,它是接收方
从节点2向节点6发送45个token，节点6崩后，节点2 锁定45token，节点3锁定45token，转帐失败；重启后，转账继续。

- 场景二：ReceiveSecretRequestStateChange

收到Secretrequest后崩
节点1向节点6发送20个token,节点6向节点1发送secretrequest请求，节点1收到崩,
节点1、节点2、节点3各锁定20个token；重启节点1后，节点锁定token解锁，转账成功。
- 场景三：ReceiveTransferRefundStateChange

收到refundtransfer崩
节点1向节点6发送45个token，（提前进行两次转账，降低部分余额，新余额分配为节点3和节点6 余额：30， 320；节点3和节点7余额： 30 90），因此，节点3要回退节点2，节点2崩；节点1锁定45，节点2，节点3锁定45，节点6未锁定；重启节点2后，重启转账成功，锁定token解锁。

测试节点2向节点6发送45个token.崩溃前情况如上；恢复节点2后，节点2和节点3锁定45token,转账成功。由于路由问题造成中间节点锁定token正常，下一版本进行优化。
- 场景四：ReceiveBalanceProofStateChange

节点1向节点6发送45个token （节点6收到balance），节点6崩。

节点1扣钱，节点2扣钱，节点3扣钱，节点6收到钱，转账成功。

此时，转账继续完成。重启后，再次转账，转账成功。
- 场景五：ActionInitMediatorStateChange

（收到mtr,它是中间节点）

节点1向节点6发送45个token,路由1，2，3，6。先设计节点2崩，再设计节点3崩。

节点2崩后，节点1锁定45token; 节点3、节点6均未锁定token;重启节点2，锁定token解锁，转账消失；

节点3崩后，节点1锁定45，节点2锁定45，节点6未锁定；重启节点3，节点1锁定45，节点2锁定45，节点3未锁定（与节点2锁定不一致，数据不同步，对方未锁定token），节点6未锁定。

总之，此种情况，转账未成功，越往后的崩溃，锁定的token越多。发生数据不同步的情况。

这种情况，原因在于新版本将原子性问题分成两步，造成数据不同步，使用老版本恢复功能。
- 场景六：BeforeSendRevealSecret

（发送secret之前）

节点1向节点6发送20 token,节点1崩。节点1锁定20 token, 节点2 成功，节点3成功，

节点6成功。重启节点1后，节点1锁定解锁，转账成功。

此种情况下，转账继续，不影响使用。
#### 3.收到ack崩溃
- 场景一：SecretRequestRecevieAck

节点2向节点6发送20个token，发送成功，节点6崩。
此种情况下，转账成功，崩溃不影响交易。
继续转账，转账成功。
- 场景二：SecretRecevieAck

（#balanceproof）
节点2向节点6发送45个token，发送成功，节点2崩。
转账成功，没有锁定token,重启后，节点2扣钱。
此种情况下，崩溃不影响交易。
- 场景三：MediatedTransferRecevieAck 

节点2向节点6发送45个token，节点2崩，节点2，3各锁定45 token
重启后，节点2、3token解锁，成功转账节点6。
- 场景四：RefundTransferRecevieAck

节点2向节点6发送45个token,节点3崩。节点2、节点3各锁定45，走路由2，4，5，6成功；
转账成功;重启后，2，3节点锁定45未解除。未完成。正常。
- 场景五：RevealSecretRecevieAck 

节点2向节点6发送20个token，节点2崩，节点2和节点3之间通道节点2锁定20 token，节点3和节点6之间转账完成；重启后，节点2锁定20解除，完成转账。
