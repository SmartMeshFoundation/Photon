// Package ice implements RFC 5245
// Interactive Connectivity Establishment (ICE):
// A Protocol for Network Address Translator (NAT)
// Traversal for Offer/Answer Protocols.
package ice

import (
	"bytes"
	"fmt"
	"net"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

// AddressType is type for ConnectionAddress.
type AddressType byte

// Possible address types.
const (
	AddressIPv4 AddressType = iota
	AddressIPv6
	AddressFQDN
)

func (a AddressType) String() string {
	switch a {
	case AddressIPv4:
		return "IPv4"
	case AddressIPv6:
		return "IPv6"
	case AddressFQDN:
		return "FQDN"
	default:
		panic("unexpected address type")
	}
}

// ConnectionAddress represents address that can be ipv4/6 or FQDN.
type ConnectionAddress struct {
	Host []byte
	IP   net.IP
	Type AddressType
}

// reset sets all fields to zero values.
func (a *ConnectionAddress) reset() {
	a.Host = a.Host[:0]
	for i := range a.IP {
		a.IP[i] = 0
	}
	a.Type = AddressIPv4
}

// Equal returns true if b equals to a.
func (a ConnectionAddress) Equal(b ConnectionAddress) bool {
	if a.Type != b.Type {
		return false
	}
	switch a.Type {
	case AddressFQDN:
		return bytes.Equal(a.Host, b.Host)
	default:
		return a.IP.Equal(b.IP)
	}
}

func (a ConnectionAddress) str() string {
	switch a.Type {
	case AddressFQDN:
		return string(a.Host)
	default:
		return a.IP.String()
	}
}

func (a ConnectionAddress) String() string {
	return a.str()
}

// CandidateType encodes the type of candidate. This specification
// defines the values "host", "srflx", "prflx", and "relay" for host,
// server reflexive, peer reflexive, and relayed candidates,
// respectively. The set of candidate types is extensible for the
// future.
type CandidateType byte

// Set of candidate types.
const (
	CandidateUnknown         CandidateType = iota
	CandidateHost                          // "host"
	CandidateServerReflexive               // "srflx"
	CandidatePeerReflexive                 // "prflx"
	CandidateRelay                         // "relay"
)

func (c CandidateType) String() string {
	switch c {
	case CandidateHost:
		return candidateHost
	case CandidateServerReflexive:
		return candidateServerReflexive
	case CandidatePeerReflexive:
		return candidatePeerReflexive
	case CandidateRelay:
		return candidateRelay
	default:
		return "unknown"
	}
}

const (
	candidateHost            = "host"
	candidateServerReflexive = "srflx"
	candidatePeerReflexive   = "prflx"
	candidateRelay           = "relay"
)

// Candidate is ICE candidate defined in RFC 5245 Section 21.1.1.
//
// This attribute is used with Interactive Connectivity
// Establishment (ICE), and provides one of many possible candidate
// addresses for communication. These addresses are validated with
// an end-to-end connectivity check using Session Traversal Utilities
// for NAT (STUN)).
//
// The candidate attribute can itself be extended. The grammar allows
// for new name/value pairs to be added at the end of the attribute. An
// implementation MUST ignore any name/value pairs it doesn't
// understand.
type Candidate struct {
	/**
	 * IP address of this candidate. For host candidates, this represents
	 * the local address of the socket. For reflexive candidates, the value
	 * will be the public address allocated in NAT router for the host
	 * candidate and as reported in MAPPED-ADDRESS or XOR-MAPPED-ADDRESS
	 * attribute of STUN Binding request. For relayed candidate, the value
	 * will be the address allocated in the TURN server by STUN Allocate
	 * request.
	 */
	addr string //192.168.0.1:3301 or www.baidu.com:3301
	/* Base address of this candidate. "Base" refers to the address an agent
	 * sends from for a particular candidate.  For host candidates, the base
	 * is the same as the host candidate itself. For reflexive candidates,
	 * the base is the local IP address of the socket. For relayed candidates,
	 * the base address is the transport address allocated in the TURN server
	 * for this candidate.
	 */
	baseAddr string //192.168.0.1:3301 etc
	/*
	  raddr 10.1.22.220 rport 51941
	*/
	relatedAddr    string        //candidates provide:  raddr 10.1.22.220 rport 56024
	transport      TransportType //must be udp,
	transportValue []byte
	/**
	 * The foundation string, which is an identifier which value will be
	 * equivalent for two candidates that are of the same type, share the
	 * same base, and come from the same STUN server. The foundation is
	 * used to optimize ICE performance in the Frozen algorithm.
	 */
	Foundation  int //TODO 为了和 pjnath 兼容,按照 RFC 规定这个 Foundation 可以是任意类型,长度不超过36,应该来说最好是 pjnath 那种定义.
	ComponentID int
	/**
	 * The candidate's priority, a 32-bit unsigned value which value will be
	 * calculated by the ICE session when a candidate is registered to the
	 * ICE session.
	 */
	Priority int
	Type     CandidateType

	// Extended attributes
	NetworkCost int
	Generation  int

	// Other attributes
	Attributes Attributes
}

/*
example:
a=candidate:Sc0a82b28 1 UDP 1694498815 14.108.27.205 55501 typ srflx
a=candidate:Hc0a82b28 1 UDP 2130706431 192.168.43.40 49927 typ host
a=candidate:Rb6fe9bd0 1 UDP 16777215 182.254.155.208 55019 typ relay

*/
func (c *Candidate) String() string {
	host, port, err := net.SplitHostPort(c.addr)
	if err != nil {
		log.Error(fmt.Sprintf("SplitHostPort %s err %s", c.addr, err))
	}
	s := fmt.Sprintf("a=candidate:%d %d %s %d %s %s typ %s",
		c.Foundation, c.ComponentID, c.transport, c.Priority, host, port, c.Type)
	if c.Type == CandidateServerReflexive {
		rhost, rport, err := net.SplitHostPort(c.baseAddr)
		if err != nil {
			log.Error(fmt.Sprintf("SplitHostPort %s err %s", c.baseAddr, err))
		}
		s = fmt.Sprintf("%s raddr %s rport %s", s, rhost, rport)
	}
	return s
}

// reset sets all fields to zero values.
func (c *Candidate) reset() {
	c.addr = ""
	c.relatedAddr = ""
	c.NetworkCost = 0
	c.Generation = 0
	c.transport = TransportUnknown
	c.transportValue = c.transportValue[:0]
	c.Attributes = c.Attributes[:0]
}

// Equal returns true if b candidate is equal to c.
func (c Candidate) Equal(b *Candidate) bool {
	if c.addr != b.addr {
		return false
	}
	if c.transport != b.transport {
		return false
	}
	//if !bytes.Equal(c.transportValue, b.transportValue) {
	//	return false
	//}
	if c.Foundation != b.Foundation {
		return false
	}
	if c.ComponentID != b.ComponentID {
		return false
	}
	if c.Priority != b.Priority {
		return false
	}
	if c.Type != b.Type {
		return false
	}
	//if c.NetworkCost != b.NetworkCost {
	//	return false
	//}
	//if c.Generation != b.Generation {
	//	return false
	//}
	//if !c.Attributes.Equal(b.Attributes) {
	//	return false
	//}

	return true
}

// Attribute is key-value pair.
type Attribute struct {
	Key   []byte
	Value []byte
}

// Attributes is list of attributes.
type Attributes []Attribute

// Value returns first attribute value with key k or
// nil of none found.
func (a Attributes) Value(k []byte) []byte {
	for _, attribute := range a {
		if bytes.Equal(attribute.Key, k) {
			return attribute.Value
		}
	}
	return nil
}

func (a Attributes) Equal(b Attributes) bool {
	if len(a) != len(b) {
		return false
	}
	for _, attr := range a {
		v := b.Value(attr.Key)
		if !bytes.Equal(v, attr.Value) {
			return false
		}
	}
	for _, attr := range b {
		v := a.Value(attr.Key)
		if !bytes.Equal(v, attr.Value) {
			return false
		}
	}
	return true
}

func (a Attribute) String() string {
	return fmt.Sprintf("%s:%s", a.Key, a.Value)
}

// TransportType is transport type for candidate.
type TransportType byte

// Supported transport types.
const (
	TransportUDP TransportType = iota + 1
	TransportUnknown
)

func (t TransportType) String() string {
	switch t {
	case TransportUDP:
		return "UDP"
	default:
		return "Unknown"
	}
}

// candidateParser should parse []byte into Candidate.
//
// a=candidate:3862931549 1 udp 2113937151 192.168.220.128 56032 typ host generation 0 network-cost 50
//     foundation ---┘    |  |      |            |          |
//   component id --------┘  |      |            |          |
//      transport -----------┘      |            |          |
//       priority ------------------┘            |          |
//  conn. address -------------------------------┘          |
//           port ------------------------------------------┘
type candidateParser struct {
	buf []byte
	c   *Candidate
}

const sp = ' '

const (
	mandatoryElements = 6
)

func parseInt(v []byte) (int, error) {
	var i int
	//跳过开始的非数字,为了和 pjnath 兼容.
	for i = range v {
		if v[i] < '0' || v[i] > '9' {
			continue
		}
		break
	}
	v2 := v[i:]
	i, err := fasthttp.ParseUint(v2)
	if err != nil {
		//为了兼容 pjnath
		fmt.Sscanf(string(v2), "%x", &i)
	}
	return i, nil
}

func (p *candidateParser) parseFoundation(v []byte) error {
	i, err := parseInt(v)
	if err != nil {
		return errors.Wrap(err, "failed to parse foundation")
	}
	p.c.Foundation = i
	return nil
}

func (p *candidateParser) parseComponentID(v []byte) error {
	i, err := parseInt(v)
	if err != nil {
		return errors.Wrap(err, "failed to parse component ID")
	}
	p.c.ComponentID = i
	return nil
}

func (p *candidateParser) parsePriority(v []byte) error {
	i, err := parseInt(v)
	if err != nil {
		return errors.Wrap(err, "failed to parse priority")
	}
	p.c.Priority = i
	return nil
}

func (p *candidateParser) parsePort(v []byte) error {
	p.c.addr = fmt.Sprintf("%s:%s", p.c.addr, string(v))
	return nil
}

func (p *candidateParser) parseRelatedPort(v []byte) error {
	p.c.relatedAddr = fmt.Sprintf("%s:%s", p.c.relatedAddr, string(v))
	return nil
}

func (p *candidateParser) parseConnectionAddress(v []byte) error {
	p.c.addr = string(v)
	return nil
}

func (p *candidateParser) parseRelatedAddress(v []byte) error {
	p.c.relatedAddr = string(v)
	return nil
}

func (p *candidateParser) parseTransport(v []byte) error {
	if bytes.Equal(v, []byte("udp")) || bytes.Equal(v, []byte("UDP")) {
		p.c.transport = TransportUDP
	} else {
		p.c.transport = TransportUnknown
		p.c.transportValue = v
	}
	return nil
}

// possible attribute keys.
const (
	aGeneration     = "generation"
	aNetworkCost    = "network-cost"
	aType           = "typ"
	aRelatedAddress = "raddr"
	aRelatedPort    = "rport"
)

func (p *candidateParser) parseAttribute(a Attribute) error {
	switch string(a.Key) {
	case aGeneration:
		return p.parseGeneration(a.Value)
	case aNetworkCost:
		return p.parseNetworkCost(a.Value)
	case aType:
		return p.parseType(a.Value)
	case aRelatedAddress:
		return p.parseRelatedAddress(a.Value)
	case aRelatedPort:
		return p.parseRelatedPort(a.Value)
	default:
		p.c.Attributes = append(p.c.Attributes, a)
		return nil
	}
}

type parseFn func(v []byte) error

const (
	minBufLen = 10
)

// parse populates internal Candidate from buffer.
func (p *candidateParser) parse() error {
	if len(p.buf) < minBufLen {
		return errors.Errorf("buffer too small (%d < %d)", len(p.buf), minBufLen)
	}
	// special cases for raw value support:
	if p.buf[0] == 'a' {
		p.buf = bytes.TrimPrefix(p.buf, []byte("a="))
	}
	if p.buf[0] == 'c' {
		p.buf = bytes.TrimPrefix(p.buf, []byte("candidate:"))
	}
	// pos is current position
	// l is value length
	// last is last character offset
	// of mandatory elements
	var pos, l, last int
	fns := []parseFn{
		p.parseFoundation,        // 0
		p.parseComponentID,       // 1
		p.parseTransport,         // 2
		p.parsePriority,          // 3
		p.parseConnectionAddress, // 4
		p.parsePort,              // 5
	}
	for i, c := range p.buf {
		if pos > mandatoryElements-1 {
			// saving offset
			last = i
			break
		}
		if c != sp {
			// non-space character
			l++
			continue
		}
		// space character reached
		if err := fns[pos](p.buf[i-l : i]); err != nil {
			return errors.Wrapf(err, "failed to parse char %d, pos %d",
				i, pos,
			)
		}
		pos++ // next element
		l = 0 // reset length of element
	}
	if last == 0 {
		// no non-mandatory elements
		return nil
	}
	// offsets:
	var (
		start  int // key start
		end    int // key end
		vStart int // value start
	)
	// subslicing to simplify offset calculation
	buf := p.buf[last-1:]
	// saving every k:v pair ignoring spaces
	for i, c := range buf {
		if c != sp && i != len(buf)-1 {
			// char is non-space or end of buffer
			if start == 0 {
				// key not started
				start = i
				continue
			}
			if vStart == 0 && end != 0 {
				// value not started and key ended
				vStart = i
			}
			continue
		}
		// char is space or end of buf reached
		if start == 0 {
			// key not started, skipping
			continue
		}
		if end == 0 {
			// key ended, saving offset
			end = i
			continue
		}
		if vStart == 0 {
			// value not started, skipping
			continue
		}
		if i == len(buf)-1 && buf[len(buf)-1] != sp {
			// fix for end of buf
			i = len(buf)
		}
		// value ended, saving attribute
		a := Attribute{
			Key:   buf[start:end],
			Value: buf[vStart:i],
		}
		if err := p.parseAttribute(a); err != nil {
			return errors.Wrapf(err, "failed to parse attribute at char %d",
				i+last,
			)
		}
		// reset offset
		vStart = 0
		end = 0
		start = 0
	}
	return nil
}

func (p *candidateParser) parseNetworkCost(v []byte) error {
	i, err := parseInt(v)
	if err != nil {
		return errors.Wrap(err, "failed to parse network cost")
	}
	p.c.NetworkCost = i
	return nil
}

func (p *candidateParser) parseGeneration(v []byte) error {
	i, err := parseInt(v)
	if err != nil {
		return errors.Wrap(err, "failed to parse generation")
	}
	p.c.Generation = i
	return nil
}

func (p *candidateParser) parseType(v []byte) error {
	switch string(v) {
	case candidateHost:
		p.c.Type = CandidateHost
	case candidatePeerReflexive:
		p.c.Type = CandidatePeerReflexive
	case candidateRelay:
		p.c.Type = CandidateRelay
	case candidateServerReflexive:
		p.c.Type = CandidateServerReflexive
	default:
		return errors.Errorf("unknown candidate %q", v)
	}
	return nil
}

// ParseAttribute parses v into c and returns error if any.
func ParseAttribute(v []byte, c *Candidate) error {
	p := candidateParser{
		buf: v,
		c:   c,
	}
	err := p.parse()
	return err
}

//按照 RFC5245计算 Preference 的算法.
func calcCandidatePriority(candidateType CandidateType, localPreference int, componentId int) int {
	var typePreference int32
	switch candidateType {
	case CandidateHost:
		typePreference = 126
	case CandidateServerReflexive:
		typePreference = 100
	case CandidatePeerReflexive:
		typePreference = 110
	case CandidateRelay:
		typePreference = 0
	default:
		typePreference = 0
	}
	p := ((typePreference & 0xff) << 24) + ((int32(localPreference) & 0xffff) << 8) +
		((256 - int32(componentId)) & 0xff)
	return int(p)
}
