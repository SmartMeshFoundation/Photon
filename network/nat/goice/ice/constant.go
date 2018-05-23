package ice

import "time"

/**
 * Maximum number of STUN transmission count.
 *
 * Default: 7 (as per RFC 3489-bis)
 */
const maxRetryBindingRequest = 7

/**
 * The STUN transaction timeout value, in miliseconds.
 * After the last retransmission is sent and if no response is received
 * after this time, the STUN transaction will be considered to have failed.
 *
 * The default value is 16x RTO (as per RFC 3489-bis).
 */
const stunTimeoutValue = 1600 * time.Millisecond

/**
* The default initial STUN round-trip time estimation (the RTO value
* in RFC 3489-bis), in miliseconds.
* This value is used to control the STUN request
* retransmit time. The initial value of retransmission interval
* would be set to this value, and will be doubled after each
* retransmission.
 */
const defaultRTOValue = time.Millisecond * 100

/**
 * The TURN permission lifetime setting. This value should be taken from the
 * TURN protocol specification.
 */
const turnPermissionTimeout = time.Second * 300

/**
 * The TURN channel binding lifetime. This value should be taken from the
 * TURN protocol specification.
 */
const turnChannelTimeout = time.Second * 600

/**
 * Number of seconds to refresh the permission/channel binding before the
 * permission/channel binding expires. This value should be greater than
 * PJ_TURN_PERM_TIMEOUT setting.
 */
const turnRefreshSecondsBefore = time.Second * 60

/**
 * The TURN session timer heart beat interval. When this timer occurs, the
 * TURN session will scan all the permissions/channel bindings to see which
 * need to be refreshed.
 */
const turnKeepAliveSecond = time.Second * 15

/**
 * Duration to keep response in the cache, in msec.
 *
 * Default: 10000 (as per RFC 3489-bis)
 */
const stunResponseCacheDuration = time.Second * 10
