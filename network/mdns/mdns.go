package mdns

import (
	"context"
	"fmt"
	"io"
	"net"
	"sort"
	"sync"
	"time"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/whyrusleeping/mdns"
)

func init() {
	// don't let mdns use logging...
	mdns.DisableLogging = true
}

var logger = log.New("pgk", "mdns")

//ServiceTag 服务类型
var ServiceTag = "_photon"

//Service interface for mdns
type Service interface {
	io.Closer
	RegisterNotifee(Notifee)
	UnregisterNotifee(Notifee)
}

//Notifee notify handler fro mdns
type Notifee interface {
	HandlePeerFound(id string, addr *net.UDPAddr)
}

type mdnsService struct {
	server   *mdns.Server
	service  *mdns.MDNSService
	myid     string
	lk       sync.Mutex
	notifees []Notifee
	interval time.Duration
}

//NewMdnsService create mdns service
func NewMdnsService(ctx context.Context, port int, myid string, interval time.Duration) (Service, error) {

	info := []string{myid}
	ips := mdns.GetLocalIP()
	log.Info(fmt.Sprintf("NewMDNSService ips=%s", ips))
	service, err := mdns.NewMDNSService(myid, ServiceTag, "", "", port, ips, info)
	if err != nil {
		return nil, err
	}

	// Create the mDNS server, defer shutdown
	server, err := mdns.NewServer(&mdns.Config{Zone: service})
	if err != nil {
		return nil, err
	}

	s := &mdnsService{
		server:   server,
		service:  service,
		interval: interval,
		myid:     myid,
	}

	go s.pollForEntries(ctx)

	return s, nil
}

func (m *mdnsService) Close() error {
	return m.server.Shutdown()
}

func (m *mdnsService) pollForEntries(ctx context.Context) {

	ticker := time.NewTicker(m.interval)
	for {
		//execute mdns query right away[ at method call and then with every tick
		entriesCh := make(chan *mdns.ServiceEntry, 16)
		go func() {
			for entry := range entriesCh {
				m.handleEntry(entry)
			}
		}()

		//log.Debug("starting mdns query")
		qp := &mdns.QueryParam{
			Domain:  "local",
			Entries: entriesCh,
			Service: ServiceTag,
			Timeout: m.interval * 5,
		}

		err := mdns.Query(qp)
		if err != nil {
			log.Error("mdns lookup error: ", err)
		}
		close(entriesCh)
		//log.Debug("mdns query complete")

		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			log.Debug("mdns service halting")
			return
		}
	}
}

//111110xxx中前面1的个数,第一个0出现之前1的个数
func simpleMaskLength(mask []byte) int {
	var n int
	for _, v := range mask {
		if v == 0xff {
			n += 8
			continue
		}
		// found non-ff byte
		// count 1 bits
		for v&0x80 != 0 {
			n++
			v <<= 1
		}
		break
	}
	return n
}

//根据对方发送过来的候选ipv4列表和我自己本地的ip进行比较,找出最匹配的那个
//比如对方告诉我自己的ip地址是[192.168.124.13,192.168.122.1],
//我的ip是[192.168.124.2,10.0.0.17],那么最匹配的应该就是192.168.124.13
//这里只处理ipv4,暂不考虑ipv6
func (m *mdnsService) getBestMatchIP(remotes []net.IP) net.IP {
	remotes2 := make([]net.IP, len(remotes))
	for i, r := range remotes {
		if r.To4() == nil {
			return net.IPv4zero
		}
		//先强制转换成4个自己的字节形式
		remotes2[i] = r.To4()
	}
	remotes = remotes2
	if len(remotes) <= 0 {
		return net.IPv4zero
	}

	//本机的ipv4地址,确保长度是4个字节而不是16个字节
	locals := m.service.IPs
	//统计最长前缀匹配
	maskLen := make(map[string]int) //每个ip地址最长前缀匹配
	for _, r := range remotes {
		//忽略非ipv4地址
		if r.To4() == nil {
			continue
		}
		max := 0
		for _, l := range locals {
			ipmask := make(net.IPMask, len(l))
			for i := 0; i < len(l) && i < len(r); i++ {
				ipmask[i] = ^(l[i] ^ r[i])
			}
			a := simpleMaskLength(ipmask)
			if max < a {
				max = a
			}
			//log.Trace(fmt.Sprintf("l=%s,%s,r=%s,%s,masksize=%d,mask=%s", l.String(),
			//	hex.EncodeToString(l), r.String(), hex.EncodeToString(r), a, hex.EncodeToString(ipmask)))
		}
		maskLen[r.String()] = max
	}
	//按照最长匹配排序,
	sort.Slice(remotes, func(i, j int) bool {
		return maskLen[remotes[i].String()] > maskLen[remotes[j].String()]
	})
	return remotes[0]
}
func (m *mdnsService) handleEntry(e *mdns.ServiceEntry) {
	log.Debug(fmt.Sprintf("Handling MDNS entry: %s:%d %s", e.AddrV4, e.Port, e.Info))

	if e.Info == m.myid {
		//log.Debug("got our own mdns entry, skipping")
		return
	}

	m.lk.Lock()
	for _, n := range m.notifees {
		n.HandlePeerFound(e.Info, &net.UDPAddr{
			IP:   m.getBestMatchIP(e.AddrV4),
			Port: e.Port,
		})
	}
	m.lk.Unlock()
}

func (m *mdnsService) RegisterNotifee(n Notifee) {
	m.lk.Lock()
	m.notifees = append(m.notifees, n)
	m.lk.Unlock()
}

func (m *mdnsService) UnregisterNotifee(n Notifee) {
	m.lk.Lock()
	found := -1
	for i, notif := range m.notifees {
		if notif == n {
			found = i
			break
		}
	}
	if found != -1 {
		m.notifees = append(m.notifees[:found], m.notifees[found+1:]...)
	}
	m.lk.Unlock()
}
