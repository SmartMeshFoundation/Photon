package mdns

import (
	"context"
	"fmt"
	"io"
	"net"
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
const ServiceTag = "_photon"

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
		log.Debug("mdns query complete")

		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			log.Debug("mdns service halting")
			return
		}
	}
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
			IP:   e.AddrV4,
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
