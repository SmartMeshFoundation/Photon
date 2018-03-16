package gopjnath

import (
	"net"
	"testing"
)

func TestSockAddr(t *testing.T) {
	s, err := NewSockAddr(AfIP, "192.168.1.55", 5000)
	if err != nil {
		t.Fatalf("NewSockAddr error: %s", err)
	}
	if ip := s.IP().String(); ip != "192.168.1.55" {
		t.Fatalf("IP should be 192.168.1.55, instead: %s", ip)
	}
	if port := s.Port(); port != 5000 {
		t.Fatalf("Port should be 5000, instead: %d", port)
	}
	s.SetIP(net.IPv4(122, 121, 120, 119))
	if ip := s.IP().String(); ip != "122.121.120.119" {
		t.Fatalf("IP should be 122.121.120.119, instead: %s", ip)
	}
	s.SetIP(net.IPv6loopback)
	if ip := s.IP().String(); ip != "::1" {
		t.Fatalf("IP should be ::1, instead: %s", ip)
	}
	s.Destroy()
}
