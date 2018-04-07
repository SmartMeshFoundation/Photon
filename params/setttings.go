package params

import (
	"fmt"

	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
)

const INITIAL_PORT = 40001

const CACHE_TTL = 60
const ESTIMATED_BLOCK_TIME = 7
const GAS_LIMIT = 3141592 //den's gasLimit.

const GAS_PRICE = params.Shannon * 20

const DEFAULT_PROTOCOL_RETRIES_BEFORE_BACKOFF = 5
const DEFAULT_PROTOCOL_THROTTLE_CAPACITY = 10.
const DEFAULT_PROTOCOL_THROTTLE_FILL_RATE = 10.
const DEFAULT_PROTOCOL_RETRY_INTERVAL = 1.

const DEFAULT_REVEAL_TIMEOUT = 10
const DEFAULT_SETTLE_TIMEOUT = DEFAULT_REVEAL_TIMEOUT * 9
const DEFAULT_EVENTS_POLL_TIMEOUT = time.Second
const DEFAULT_POLL_TIMEOUT = 180 * time.Second
const DEFAULT_JOINABLE_FUNDS_TARGET = 0.4
const DEFAULT_INITIAL_CHANNEL_TARGET = 3
const DEFAULT_WAIT_FOR_SETTLE = true

const DEFAULT_NAT_KEEPALIVE_RETRIES = 5
const DEFAULT_NAT_KEEPALIVE_TIMEOUT = 500
const DEFAULT_NAT_INVITATION_TIMEOUT = 15000
const Default_Tx_Timeout = 5 * time.Minute //15seconds for one block,it may take sever minutes
const MaxRequestTimeout = 20 * time.Minute //longest time for a request ,for example ,settle all channles?

var GAS_LIMIT_HEX string

//for private
var ROPSTEN_REGISTRY_ADDRESS = common.HexToAddress("0x1BB1437d4e387Be1E8C04762536217B3240f2323")

//private clean
//var ROPSTEN_REGISTRY_ADDRESS = common.HexToAddress("0x7F167440F1aB963ddbDe19F5B355e1889E9DA187")
var ROPSTEN_DISCOVERY_ADDRESS = common.HexToAddress("0x95A4e1251B87DCEf6B0cD18D3356CdA8cFB8f6CC")

//for testnet
//var ROPSTEN_REGISTRY_ADDRESS = common.HexToAddress("66eea3159a01d134dd64bfe36fde4be9ed9c1695")
//var ROPSTEN_DISCOVERY_ADDRESS = common.HexToAddress("1e3941d8c05fffa7466216480209240cc26ea577")

const NETTINGCHANNEL_SETTLE_TIMEOUT_MIN = 6

/*
 The maximum settle timeout is chosen as something above
 1 year with the assumption of very fast block times of 12 seconds.
 There is a maximum to avoidpotential overflows as described here:
 https://github.com/raiden-network/raiden/issues/1038
*/
const NETTINGCHANNEL_SETTLE_TIMEOUT_MAX = 2700000

const UDP_MAX_MESSAGE_SIZE = 1200
const DefaultSignalServer = "119.28.43.121:5222"

func init() {
	GAS_LIMIT_HEX = fmt.Sprintf("0x%x", GAS_LIMIT)
}
