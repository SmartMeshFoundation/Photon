package main

import (
	"time"
	"bytes"
	"net"
	"github.com/SmartMeshFoundation/SmartRaiden/network/nat/goice/ice"
	"github.com/nkbai/log"
)
const (
	typHost = 1
	typStun = 2
	typTurn = 3
)

type icecb struct {
	data      chan []byte
	iceresult chan error
	name      string
}

func Newicecb(name string) *icecb {
	return &icecb{
		name:      name,
		data:      make(chan []byte, 1),
		iceresult: make(chan error, 1),
	}
}
func (c *icecb) OnReceiveData(data []byte, from net.Addr) {
	c.data <- data
}

/*
	Callback to report status of various ICE operations.
*/
func (c *icecb) OnIceComplete(result error) {
	c.iceresult <- result
	log.Trace("%s negotiation complete", c.name)
}
func setupIcePair(typ int) (s1, s2 *ice.IceStreamTransport, err error) {
	var cfg *ice.TransportConfig
	switch typ {
	case typHost:
		cfg = ice.NewTransportConfigHostonly()
	case typStun:
		cfg = ice.NewTransportConfigWithStun("182.254.155.208:3478")
	case typTurn:
		cfg = ice.NewTransportConfigWithTurn("182.254.155.208:3478", "bai", "bai")
	}
	s1, err = ice.NewIceStreamTransport(cfg, "s1")
	if err != nil {
		return
	}
	s2, err = ice.NewIceStreamTransport(cfg, "s2")
	log.Trace("-----------------------------------------")
	return
}
func main(){
	s1, s2, err := setupIcePair(typTurn)
	if err != nil {
		log.Crit(err.Error())
		return
	}
	cb1 := Newicecb("s1")
	cb2 := Newicecb("s2")
	s1.SetCallBack(cb1)
	s2.SetCallBack(cb2)
	err = s1.InitIce(ice.SessionRoleControlling)
	if err != nil {
		log.Crit(err.Error())
		return
	}
	err = s2.InitIce(ice.SessionRoleControlled)
	if err != nil {
		log.Crit(err.Error())
		return
	}
	rsdp, err := s2.EncodeSession()
	if err != nil {
		log.Crit(err.Error())
		return
	}
	err = s1.StartNegotiation(rsdp)
	if err != nil {
		log.Crit(err.Error())
		return
	}
	lsdp, err := s1.EncodeSession()
	if err != nil {
		log.Crit(err.Error())
		return
	}
	err = s2.StartNegotiation(lsdp)
	if err != nil {
		log.Crit(err.Error())
		return
	}
	select {
	case <-time.After(10 * time.Second):
		log.Error("s1 negotiation timeout")
		return
	case err = <-cb1.iceresult:
		if err != nil {
			log.Error("s1 negotiation failed ", err)
			return
		}
	}
	select {
	case <-time.After(10 * time.Second):
		log.Error("s2 negotiation timeout")
		return
	case err = <-cb2.iceresult:
		if err != nil {
			log.Error("s2 negotiation failed", err)
			return
		}
	}
	s1data := []byte("hello,s2")
	s2data := []byte("hello,s1")
	err = s1.SendData(s1data)
	if err != nil {
		log.Crit(err.Error())
		return
	}
	err = s2.SendData(s2data)
	if err != nil {
		log.Crit(err.Error())
		return
	}
	select {
	case <-time.After(10 * time.Second):
		log.Error("s2 recevied timeout")
		return
	case data := <-cb2.data:
		if !bytes.Equal(data, s1data) {
			log.Error("s2 recevied error ,got ", string(data))
			return
		}
	}
	select {
	case <-time.After(10 * time.Second):
		log.Error("s1 recevied timeout")
		return
	case data := <-cb1.data:
		if !bytes.Equal(data, s2data) {
			log.Error("s1 recevied error ,got ", string(data))
			return
		}
	}
}
