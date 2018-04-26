package main

import (
	"bytes"
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/nat/goice/stun"
	"github.com/SmartMeshFoundation/SmartRaiden/network/nat/goice/turn"
)

var (
	server = flag.String("server",
		fmt.Sprintf("182.254.155.208:3478"),
		"turn server address",
	)
	peer = flag.String("peer",
		"182.254.155.208:3333", //test echo server
		"peer addres",
	)
	username = flag.String("username", "bai", "username")
	password = flag.String("password", "bai", "password")
)

const (
	udp = "udp"
)

func isErr(m *stun.Message) bool {
	return m.Type.Class == stun.ClassErrorResponse
}

func do(req, res *stun.Message, c *net.UDPConn, attrs ...stun.Setter) error {
	start := time.Now()
	if err := req.Build(attrs...); err != nil {
		log.Error("failed to build %s", err)
		return err
	}
	if _, err := req.WriteTo(c); err != nil {
		log.Error("failed to write %s m:", err, req)
		return err
	}
	//log.Info("sent message m:%s", req)
	if cap(res.Raw) < 800 {
		res.Raw = make([]byte, 0, 1024)
	}
	c.SetReadDeadline(time.Now().Add(time.Second * 20))
	_, err := res.ReadFrom(c)
	if err != nil {
		log.Error("failed to read  err:%s message:%s", err, req)
	}
	log.Info("got message m:%s, rtt:%s", res, time.Since(start))
	return err
}

func main() {
	flag.Parse()
	var (
		req = new(stun.Message)
		res = new(stun.Message)
	)
	if flag.Arg(0) == "peer" {
		_, port, err := net.SplitHostPort(*peer)
		log.Info("running in peer mode")
		if err != nil {
			log.Crit("failed to find port in peer address %s", err)
		}
		laddr, err := net.ResolveUDPAddr(udp, ":"+port)
		if err != nil {
			log.Crit("failed to resolve UDP addr  %s", err)
		}
		c, err := net.ListenUDP(udp, laddr)
		if err != nil {
			log.Crit("failed to listen  %s", err)
		}
		log.Info("listening as echo server laddr:%s", c.LocalAddr())
		for {
			// Starting echo server.
			buf := make([]byte, 1024)
			n, addr, err := c.ReadFromUDP(buf)
			if err != nil {
				log.Crit("failed to read  %s", err)
			}
			log.Info("got message body:%s raddr:%s", string(buf[:]), addr)
			// Echoing back.
			if _, err := c.WriteToUDP(buf[:n], addr); err != nil {
				log.Crit("failed to write back %s", err)
			}
			log.Info("echoed back raddr:%s ", addr)
		}
	}
	if len(*password) == 0 {
		fmt.Fprintln(os.Stderr, "No password set, auth is required.")
		flag.Usage()
		os.Exit(2)
	}

	// Resolving to TURN server.
	raddr, err := net.ResolveUDPAddr(udp, *server)
	if err != nil {
		log.Crit("failed to resolve TURN server %s", err)
	}
	c, err := net.DialUDP(udp, nil, raddr)
	if err != nil {
		log.Crit("failed to dial to TURN server %s", err)
	}
	log.Info("dial server laddr:%s raddr:%s", c.LocalAddr(), c.RemoteAddr())

	// Crafting allocation request.
	if err := do(req, res, c,
		stun.TransactionIDSetter,
		turn.AllocateRequest,
		turn.RequestedTransportUDP,
	); err != nil {
		log.Crit("do failed %s", err)
	}
	var (
		code  stun.ErrorCodeAttribute
		nonce stun.Nonce
		realm stun.Realm
	)
	if res.Type.Class != stun.ClassErrorResponse {
		log.Crit("expected error class, got " + res.Type.Class.String())
	}
	if err := code.GetFrom(res); err != nil {
		log.Crit("failed to get error code from message %s",
			err,
		)
	}
	if code.Code != stun.CodeUnauthorised {
		log.Crit("unexpected code of error err:%s", code)
	}
	if err := nonce.GetFrom(res); err != nil {
		log.Crit("failed to nonce from message", err)
	}
	if err := realm.GetFrom(res); err != nil {
		log.Crit("failed to get realm from message", err)
	}
	realmStr := realm.String()
	nonceStr := nonce.String()
	log.Info("got credentials  nonce:%s,realm:%s", nonce, realm)
	var (
		credentials = stun.NewLongTermIntegrity(*username, realm.String(), *password)
	)
	log.Info("using integrity i:%s", credentials)

	// Constructing allocate request with integrity
	req = new(stun.Message)
	if err := do(req, res, c, stun.TransactionIDSetter, turn.AllocateRequest,
		turn.RequestedTransportUDP, realm,
		stun.NewUsername(*username), nonce, credentials,
	); err != nil {
		log.Crit("failed to do request %s", err)
	}
	if isErr(res) {
		code.GetFrom(res)
		log.Crit("got error response %s ", code)
	}
	// Decoding relayed and mapped address.
	var (
		reladdr turn.RelayedAddress
		maddr   stun.XORMappedAddress
	)
	if err := reladdr.GetFrom(res); err != nil {
		log.Crit("failed to get relayed address %s", err)
	}
	log.Info("relayed address addr:%s", reladdr)
	if err := maddr.GetFrom(res); err != nil && err != stun.ErrAttributeNotFound {
		log.Crit("failed to decode relayed address %s", err)
	} else {
		log.Info("mapped address %s", maddr)
	}

	//test sendindication
	//if err := do(req, res, c, stun.TransactionIDSetter, turn.SendIndication); err != nil {
	//	log.Crit("failed to sendindication %s", err)
	//}

	// Creating permission request.
	echoAddr, err := net.ResolveUDPAddr(udp, *peer)
	if err != nil {
		log.Crit("failed to resonve addr %s", err)
	}
	peerAddr := turn.PeerAddress{
		IP:   echoAddr.IP,
		Port: echoAddr.Port,
	}
	log.Info("peer address addr:%s", peerAddr)
	if err := do(req, res, c, stun.TransactionIDSetter,
		turn.CreatePermissionRequest,
		peerAddr,
		stun.Realm(realmStr),
		stun.Nonce(nonceStr),
		stun.Username(*username),
		credentials,
	); err != nil {
		log.Crit("failed to do request %s", err)
	}
	if isErr(res) {
		code.GetFrom(res)
		log.Crit("failed to allocate %s ", code)
	}
	if err := credentials.Check(res); err != nil {
		log.Error("failed to check integrity %s", err)
	}
	var (
		sentData = turn.Data("Hello world!")
	)
	// Allocation succeed.
	// Sending data to echo server.
	// can be as resetTo(type, attrs)?
	if err := do(req, res, c, stun.TransactionIDSetter,
		turn.SendIndication,
		sentData,
		peerAddr,
		stun.Fingerprint,
	); err != nil {
		log.Crit("failed to build %s", err)
	}
	log.Info("sent data %s", string(sentData))
	if isErr(res) {
		code.GetFrom(res)
		log.Crit("got error response %s", code)
	}
	var data turn.Data
	if err := data.GetFrom(res); err != nil {
		log.Crit("failed to get DATA attribute %s", err)
	}
	log.Info("got data v:%s", string(data))
	var peer turn.PeerAddress
	if err := peer.GetFrom(res); err != nil {
		log.Crit("failed to get peer addr %s", err)
	}
	log.Info("peer is :%s", peer.String())
	if bytes.Equal(data, sentData) {
		log.Info("OK")
	} else {
		log.Info("DATA missmatch")
	}

	if true {
		//for channel data
		var (
			sentData = turn.Data("Hello world, channel!")
		)
		// Allocation succeed.
		// Sending data to echo server.
		// can be as resetTo(type, attrs)?
		if err := do(req, res, c, stun.TransactionIDSetter,
			turn.ChannelBindRequest,
			turn.ChannelNumber(0x4000),
			peerAddr,
			stun.Username(*username),
			stun.Realm(realmStr),
			stun.Nonce(nonceStr),
			credentials,
		); err != nil {
			log.Crit("failed to build %s", err)
		}
		log.Info("sent data %s", string(sentData))
		if isErr(res) {
			code.GetFrom(res)
			log.Crit("got error response %s", code)
		}
		log.Info("channel bind success")
		var channelNumber uint16 = 0x4000
		//send data on channel data
		cdata := &turn.ChannelData{
			ChannelNumber: channelNumber,
			Data:          []byte("hello,data from channel"),
		}
		if err := do(req, res, c, turn.ChannelDataRequest, cdata); err != nil {
			log.Crit("failed to build %s", err)
		}
		if isErr(res) {
			code.GetFrom(res)
			log.Crit("got error response %s", code)
		}
		var cdata2 = &turn.ChannelData{}
		if err := cdata2.GetFrom(res); err != nil {
			log.Crit("failed to get channel data %s", err)
		}
		if cdata2.ChannelNumber != channelNumber {
			log.Crit("channel number not equal expect=%d,got %d", channelNumber, cdata2.ChannelNumber)
		}
		if !bytes.Equal(cdata2.Data, cdata.Data) {
			log.Crit("data not equal")
		}
		log.Info("received channel data :%s", string(cdata2.Data))
	}
	//return
	// De-allocating.
	if err := do(req, res, c, stun.TransactionIDSetter,
		turn.RefreshRequest,
		stun.Realm(realmStr),
		stun.Username(*username),
		stun.Nonce(nonceStr),
		turn.ZeroLifetime,
		credentials,
	); err != nil {
		log.Crit("failed to do %s", err)
	}
	if isErr(res) {
		code.GetFrom(res)
		log.Crit("got error response %s", code)
	}
	log.Info("closing")
}
