package ice

import (
	/* #nosec */
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"net"
	"sort"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
)

const defaultPreference = 0xffff

// Gatherer is source for addresses.
//
// See https://tools.ietf.org/html/rfc5245#section-2.1 for details
// on gathering.
type Gatherer interface {
	Gather() ([]Addr, error)
}

const precedencesCount = 11

var precedences [precedencesCount]precedenceConfig

type precedenceConfig struct {
	ipNet *net.IPNet
	value int
}

func init() {
	// Initializing precedences for IP.
	/*
	   ::1/128               50     0
	   ::/0                  40     1
	   ::ffff:0:0/96         35     4
	   2002::/16             30     2
	   2001::/32              5     5
	   fc00::/7               3    13
	   ::/96                  1     3
	   fec0::/10              1    11
	   3ffe::/16              1    12
	*/
	for i, p := range [precedencesCount]struct {
		cidr  string
		value int
		label int
	}{
		{"::1/128", 50, 0},
		{"127.0.0.1/8", 45, 0},
		{"::/0", 40, 1},
		{"::ffff:0:0/96", 35, 4},
		{"fe80::/10", 33, 1},
		{"2002::/16", 30, 2},
		{"2001::/32", 5, 5},
		{"fc00::/7", 3, 13},
		{"::/96", 1, 3},
		{"fec0::/10", 1, 11},
		{"3ffe::/16", 1, 12},
	} {
		_, ipNet, err := net.ParseCIDR(p.cidr)
		if err != nil {
			panic(err)
		}
		precedences[i] = precedenceConfig{
			ipNet: ipNet,
			value: p.value,
		}
	}
}

// addr represents gathered address from interface.
type Addr struct {
	IP         net.IP
	Zone       string
	Precedence int
}

//Addrs represents addr List
type Addrs []Addr

func (s Addrs) Less(i, j int) bool {
	return s[i].Precedence > s[j].Precedence
}

func (s Addrs) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Addrs) Len() int {
	return len(s)
}

func (a Addr) String() string {
	if len(a.Zone) > 0 {
		return fmt.Sprintf("%s (zone %s) [%d]",
			a.IP, a.Zone, a.Precedence,
		)
	}
	return fmt.Sprintf("%s [%d]", a.IP, a.Precedence)
}

// ZeroPortAddr return address with "0" port.
func (a Addr) ZeroPortAddr() string {
	host := a.IP.String()
	if len(a.Zone) > 0 {
		host += "%" + a.Zone
	}
	return net.JoinHostPort(host, "")
}

type defaultGatherer struct{}

func (defaultGatherer) precedence(ip net.IP) int {
	for _, p := range precedences {
		if p.ipNet.Contains(ip) {
			return p.value
		}
	}
	return 0
}

func (g defaultGatherer) Gather() ([]Addr, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	addrs := make([]Addr, 0, 10)
	for _, iface := range interfaces {
		log.Trace(fmt.Sprintf("%s flag=%s ", iface.Name, iface.Flags.String()))
		if iface.Flags&net.FlagUp != net.FlagUp {
			//skip invalid interface. shutdown now .
			continue
		}
		iAddrs, err := iface.Addrs()
		if err != nil {
			return addrs, err
		}
		for _, a := range iAddrs {
			ip, _, err := net.ParseCIDR(a.String())
			if err != nil {
				return addrs, err
			}
			//ipv4 only
			if len(ip.To4()) != net.IPv4len {
				continue //just support ipv4 now
			}
			if ip.IsLoopback() {
				continue
			}
			addr := Addr{
				IP:         ip,
				Precedence: g.precedence(ip),
			}
			if ip.IsLinkLocalUnicast() {
				// Zone must be set for link-local addresses.
				addr.Zone = iface.Name
			}
			addrs = append(addrs, addr)
		}
	}
	sort.Sort(Addrs(addrs))
	return addrs, nil
}

// DefaultGatherer uses net.Interfaces to gather addresses.
var DefaultGatherer Gatherer = defaultGatherer{}

/*
返回所有可能的
*/
const maxCandidates = 8 //candidate 列表中最多有多少个,太多了可能是攻击
func getLocalCandidates(primaryAddress string) (candidates []*Candidate, err error) {
	_, port, err := net.SplitHostPort(primaryAddress)
	if err != nil {
		return
	}
	addrs, err := DefaultGatherer.Gather()
	if err != nil {
		return
	}
	primaryFound := false
	for _, a := range addrs {
		c := new(Candidate)
		c.Type = CandidateHost
		c.addr = fmt.Sprintf("%s:%s", a.IP.String(), port)
		c.baseAddr = c.addr
		c.Foundation = calcFoundation(c.baseAddr)
		duplicate := false
		for _, c2 := range candidates {
			if c2.Equal(c) {
				duplicate = true
				break
			}
		}
		if duplicate {
			log.Trace(fmt.Sprintf("host %s:%s is duplicate", c.addr, c.baseAddr))
			continue
		}
		if c.addr == primaryAddress {
			primaryFound = true
			if len(candidates) != 0 {
				//保证候选列表中的第一个是我们的主要地址,也就是连接 stun server 的地址.
				t := candidates[0]
				l := len(candidates)
				candidates[0] = c
				candidates[l-1] = t
				primaryFound = true
			} else {
				candidates = append(candidates, c)
			}
		} else {
			candidates = append(candidates, c)
		}

	}
	if !primaryFound {
		log.Error(fmt.Sprintf("primaryaddress not found %s", primaryAddress))
	}
	if len(candidates) > maxCandidates-1 {
		candidates = candidates[:maxCandidates-1]
	}
	return
}

func calcFoundation(baseAddr string) int {
	/* #nosec */
	hash := md5.Sum([]byte(baseAddr))
	tmp := binary.BigEndian.Uint32(hash[:4])
	return int(tmp)
}

func addCandidates(candidates []*Candidate, new *Candidate) []*Candidate {
	for _, c := range candidates {
		if c.Equal(new) {
			return candidates
		}
	}
	candidates = append(candidates, new)
	return candidates
}
