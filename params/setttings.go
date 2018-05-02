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

var ROPSTEN_REGISTRY_ADDRESS = common.HexToAddress("0xd01Ca23F2B84AF393550271bFCC2A8b48d6f65b8")
var ROPSTEN_DISCOVERY_ADDRESS = common.HexToAddress("0x4CDAF98516490d42E1E6F050bcfBD143dCb58CcD")

const NETTINGCHANNEL_SETTLE_TIMEOUT_MIN = 6

/*
 The maximum settle timeout is chosen as something above
 1 year with the assumption of very fast block times of 12 seconds.
 There is a maximum to avoidpotential overflows as described here:
 https://github.com/SmartRaiden/raiden/issues/1038
*/
const NETTINGCHANNEL_SETTLE_TIMEOUT_MAX = 2700000

const UDP_MAX_MESSAGE_SIZE = 1200
const DefaultSignalServer = "139.199.6.114:5222"

func init() {
	GAS_LIMIT_HEX = fmt.Sprintf("0x%x", GAS_LIMIT)
}
