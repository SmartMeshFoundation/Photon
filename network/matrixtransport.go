package network

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/SmartMeshFoundation/Photon/network/wakeuphandler"

	"github.com/SmartMeshFoundation/Photon/network/gomatrix"

	"time"

	"github.com/SmartMeshFoundation/Photon/channel/channeltype"
	"github.com/SmartMeshFoundation/Photon/encoding"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/network/netshare"
	"github.com/SmartMeshFoundation/Photon/network/xmpptransport"
	"github.com/SmartMeshFoundation/Photon/params"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	// ONLINE network state -online
	ONLINE = "online"
	// UNAVAILABLE state -unavailable
	UNAVAILABLE = "unavailable"
	// OFFLINE state -offline
	OFFLINE = "offline"
	// UNKNOWN or other state -unknown
	UNKNOWN = "unknown"
	// ROOMPREFIX room prefix
	ROOMPREFIX = "photon"
	// ROOMSEP with ',' to separate room name's part
	ROOMSEP = "_"
	// PATHPREFIX0 the lastest matrix client api version
	PATHPREFIX0 = "/_matrix/client/r0"
	// LOGINTYPE login type we used
	LOGINTYPE = "m.login.password"
	// CHATPRESET the type of chat=public
	CHATPRESET = "public_chat"
	//EventAddressRoom is user defined event type
	EventAddressRoom = "network.photon.rooms"
)

type jobType int

const (
	jobSendMessage = iota
	jobPresence
	jobMessage
	jobMemberShip
	jobAccountData
)

/*
move process of matrix event to one goroutine
*/
type matrixJob struct {
	jobType jobType
	Data1   interface{}
	Data2   interface{}
}

// MatrixTransport represents a matrix transport Instantiation
type MatrixTransport struct {
	matrixcli      *gomatrix.MatrixClient //the instantiated matrix
	servername     string                 //the homeserver's name
	serverURL      string                 //http://transport01.smartmesh.cn
	running        bool                   //running status
	stopreceiving  bool                   //Whether to stop accepting(data)
	key            *ecdsa.PrivateKey      //key
	NodeAddress    common.Address
	protocol       ProtocolReceiver
	Peers          map[common.Address]*MatrixPeer //这里面保存的是与我有通道的节点
	temporaryPeers *matrixTemporaryPeers          //与我没有通道,但是需要临时通信的节点
	UserID         string                         //the current user's ID(@kitty:thisserver)
	NodeDeviceType string
	log            log.Logger
	statusChan     chan netshare.Status
	removePeerChan chan common.Address
	status         netshare.Status
	servers        map[string]string
	trustServers   map[string]bool // these server's userID can be trusted
	db             xmpptransport.XMPPDb
	/*
		标记是否已经完成了基本的启动,在该标志为true之前
		1. 是不能发送消息
		2. 收到消息会忽略
		3. 收到的memshipchange,presence等消息会立即处理因为还没有进入主循环
	*/
	hasDoneStartCheck        bool
	quitChan                 chan struct{}
	jobChan                  chan *matrixJob
	lock                     sync.RWMutex
	discoveryroom            string
	wakeUpChanListMapLock    sync.Mutex
	isAlreadyCollectedDbInfo bool

	// 挂起/唤醒服务
	*wakeuphandler.WakeUpHandler
}

var (
	// ValidUserIDRegex user ID 's format
	ValidUserIDRegex = regexp.MustCompile(`^@(0x[0-9a-f]{40})(?:\.[0-9a-f]{8})?(?::.+)?$`) //(`^[0-9a-z_\-./]+$`)
	//NETWORKNAME which network is used
	NETWORKNAME = params.NETWORKNAME
	//ALIASFRAGMENT the terminal part of alias
	ALIASFRAGMENT = params.AliasFragment
	//DISCOVERYROOMSERVER discovery room server name
	DISCOVERYROOMSERVER   = params.DiscoveryServer
	networkPartHasChannel = "y"
	networkPartNoChannel  = "n"
)

// NewMatrixTransport init matrix
func NewMatrixTransport(logname string, key *ecdsa.PrivateKey, devicetype string, servers map[string]string, dao xmpptransport.XMPPDb) *MatrixTransport {
	mtr := &MatrixTransport{
		running:        false,
		stopreceiving:  false,
		NodeAddress:    crypto.PubkeyToAddress(key.PublicKey),
		key:            key,
		Peers:          make(map[common.Address]*MatrixPeer),
		temporaryPeers: newMatrixTemporaryPeers(),
		NodeDeviceType: devicetype,
		log:            log.New("matrix", logname),
		statusChan:     make(chan netshare.Status, 10),
		status:         netshare.Disconnected,
		servers:        servers,
		quitChan:       make(chan struct{}),
		jobChan:        make(chan *matrixJob, 100),
		trustServers:   make(map[string]bool),
		WakeUpHandler:  wakeuphandler.NewWakeupHandler("matrix"),
		db:             dao,
	}
	var serverNames []string
	for s := range servers {
		serverNames = append(serverNames, s)
	}
	mtr.setTrustServers(serverNames)
	return mtr
}

func (m *MatrixTransport) changeStatus(newStatus netshare.Status) {
	if m.status == newStatus {
		return
	}
	m.log.Info(fmt.Sprintf("changeStatus from %d to %d", m.status, newStatus))
	m.status = newStatus
	select {
	case m.statusChan <- newStatus:
	default:
	}
}

//todo fixme 不要用set,应该在创建的时候指定,或者从指定服务器上下载才行
func (m *MatrixTransport) setTrustServers(servers []string) {
	for _, s := range servers {
		m.trustServers[s] = true
	}
}

func (m *MatrixTransport) addPeerIfNotExist(peer common.Address) bool {

	m.lock.Lock()
	defer m.lock.Unlock()
	p, ok := m.Peers[peer]
	if ok {
		p.increaseChannelCount()
		return false
	}
	m.Peers[peer] = NewMatrixPeer(peer)
	return true
}
func (m *MatrixTransport) isTrustedServerUser(userID string) bool {
	_, domain, err := extractUserInfo(userID)
	if err != nil {
		return false
	}
	return m.trustServers[domain]
}

// collectChannelInfo subscribe status change
func (m *MatrixTransport) collectChannelInfo(db xmpptransport.XMPPDb) error {
	//cs =[nil] when getchannellist on first\second or third time
	if m.isAlreadyCollectedDbInfo {
		return nil
	}
	cs, err := db.GetChannelList(utils.EmptyAddress, utils.EmptyAddress)
	if err != nil {
		return err
	}
	m.isAlreadyCollectedDbInfo = true
	for _, c := range cs {
		m.addPeerIfNotExist(c.PartnerAddress())
	}
	db.RegisterNewChannelCallback(func(c *channeltype.Serialization) (remove bool) { //create new channel->add participant for matrix peer
		if m.addPeerIfNotExist(c.PartnerAddress()) {
			m.lock.RLock()
			p := m.Peers[c.PartnerAddress()]
			m.lock.RUnlock()
			go func() {
				err := m.startupCheckOneParticipant(p)
				if err != nil {
					log.Error(fmt.Sprintf("handleNewPartner for %s,err %s", utils.APex2(c.PartnerAddress()), err))
				}
			}()
		}
		return false
	})
	db.RegisterChannelSettleCallback(func(c *channeltype.Serialization) (remove bool) { //settle channel->remove participant for matrix peer
		m.lock.RLock()
		p := m.Peers[c.PartnerAddress()]
		m.lock.RUnlock()
		if p.decreaseChannelCount() {
			//todo mark peer to delete, delete on next restart
			log.Info(fmt.Sprintf("matrix User %s should be removed", utils.APex2(c.PartnerAddress())))
		}
		return false
	})
	return nil
}

// HandleMessage regist the interface of call receive(func)
func (m *MatrixTransport) HandleMessage(from common.Address, data []byte) {
	if !m.running || m.stopreceiving {
		return
	}
	if m.protocol != nil {
		m.protocol.receive(data)
	}
}

// RegisterProtocol regist the interface of call RegisterProtocol(func)
func (m *MatrixTransport) RegisterProtocol(protcol ProtocolReceiver) {
	m.protocol = protcol
}

// Stop Does Stop need to destroy matrix resource ?
func (m *MatrixTransport) Stop() {
	if m.running == false {
		return
	}
	m.running = false
	m.changeStatus(netshare.Closed)
	close(m.quitChan)
	if m.matrixcli != nil {
		m.log.Info("matrix will offline")
		err := m.matrixcli.SetPresenceState(&gomatrix.ReqPresenceUser{
			Presence: OFFLINE,
		})
		if err != nil {
			m.log.Error(fmt.Sprintf("[Matrix] SetPresenceState failed : %s", err.Error()))
		}
		m.matrixcli.StopSync()
		if _, err := m.matrixcli.Logout(); err != nil {
			m.log.Error("[Matrix] Logout failed")
		}
	}
}

// StopAccepting stop receive message and wait
func (m *MatrixTransport) StopAccepting() {
	m.stopreceiving = true
}

func (m *MatrixTransport) nodeStatusInternal(addr common.Address) (deviceType string, isOnline bool, isPartner bool) {
	if m.matrixcli == nil {
		return "", false, false
	}
	m.lock.RLock()
	u, ok := m.Peers[addr]
	//m.log.Debug(fmt.Sprintf("m.Peers1:%s",utils.StringInterface(m.Peers,5)))
	m.lock.RUnlock()
	if !ok {
		return "", false, false
	}
	return u.deviceType, u.status == peerStatusOnline, true
}

// NodeStatus gets Node states of network, if check self node, `status` is not always be true instead it switches according to server handshake signal.
func (m *MatrixTransport) NodeStatus(addr common.Address) (deviceType string, isOnline bool) {
	deviceType, isOnline, _ = m.nodeStatusInternal(addr)
	return
}

// Send send message
func (m *MatrixTransport) Send(receiverAddr common.Address, data []byte) error {
	if !m.running || len(data) == 0 {
		return fmt.Errorf("[Matrix]Send failed,matrix not running or send data is null")
	}
	if m.matrixcli == nil {
		return errors.New("no matrix connection")
	}
	if !m.hasDoneStartCheck {
		return errors.New("ignore message when not startup complete")
	}
	//_, isOnline, isPartner := m.nodeStatusInternal(receiverAddr)
	/*if !isOnline && isPartner {
		//如果接收方不在线,会重复不停的发送,造成不必要的网络浪费.
		return fmt.Errorf("message receiver %s  not online ", receiverAddr.String())
	}*/
	//m.log.Trace(fmt.Sprintf("sendmsg  %s", utils.StringInterface(m.Peers, 7)))
	//send should not block
	select {
	case m.jobChan <- &matrixJob{
		jobType: jobSendMessage,
		Data1:   receiverAddr,
		Data2:   data,
	}:
	default:
	}

	return nil

}
func (m *MatrixTransport) doSend(job *matrixJob) {
	var err error
	receiverAddr := job.Data1.(common.Address)
	data := job.Data2.([]byte)
	//m.log.Trace(fmt.Sprintf("send msg %s", string(data)))
	m.lock.RLock()
	p := m.Peers[receiverAddr]
	m.lock.RUnlock()
	var roomID string
	if p == nil {
		roomID = m.temporaryPeers.getRoomID(receiverAddr)
		if roomID == "" {
			var users []*gomatrix.UserInfo
			roomID, users, err = m.findOrCreateRoomByAddress(receiverAddr, false)
			if err != nil || roomID == "" {
				m.log.Error(fmt.Sprintf("[Matrix]Send failed,cann't find the peer address findOrCreateRoomByAddress err %s", err))
				return
			}
			m.temporaryPeers.addPeer(receiverAddr, roomID)
			//whether these users are in this room or not ,invite them. maybe dupclicate.
			for _, u := range users {
				_, err = m.matrixcli.InviteUser(roomID, &gomatrix.ReqInviteUser{
					UserID: u.UserID,
				})
				if err != nil {
					httpErr, ok := err.(gomatrix.HTTPError)
					//can ignore, it's possible duplicate
					if ok && httpErr.Code == http.StatusForbidden {
						continue
					}
					m.log.Error(fmt.Sprintf("InviteUser %s to room %s err %s", u.UserID, roomID, err))
				}
			}
		}
	} else {
		roomID = p.defaultMessageRoomID
		if roomID == "" {
			var users []*gomatrix.UserInfo
			roomID, users, err = m.findOrCreateRoomByAddress(receiverAddr, false)
			if err != nil || roomID == "" {
				m.log.Error(fmt.Sprintf("[Matrix]Send failed,cann't find the peer address findOrCreateRoomByAddress err %s", err))
				return
			}
			m.temporaryPeers.addPeer(receiverAddr, roomID)
			//whether these users are in this room or not ,invite them. maybe dupclicate.
			for _, u := range users {
				_, err = m.matrixcli.InviteUser(roomID, &gomatrix.ReqInviteUser{
					UserID: u.UserID,
				})
				if err != nil {
					httpErr, ok := err.(gomatrix.HTTPError)
					//can ignore, it's possible duplicate
					if ok && httpErr.Code == http.StatusForbidden {
						continue
					}
					m.log.Error(fmt.Sprintf("InviteUser %s to room %s err %s", u.UserID, roomID, err))
				}
			}
		}
	}
	_data := base64.StdEncoding.EncodeToString(data)
	m.log.Trace(fmt.Sprintf("send to %s[%s] message=%s", utils.APex2(receiverAddr), roomID, encoding.MessageType(data[0])))
	_, err = m.matrixcli.SendText(roomID, _data)
	if err != nil {
		m.log.Error(fmt.Sprintf("[matrix]send failed to %s, message=%s err=%s", utils.APex2(receiverAddr), encoding.MessageType(data[0]), err))
		return
	}
	m.log.Trace(fmt.Sprintf("[Matrix]Send to %s success, message=%s", utils.APex2(receiverAddr), encoding.MessageType(data[0])))
	return
}

/*
Start matrix
后台不断重试登陆，注册，并初始化相关信息
如果网络连接正常的话，会保证登陆，初始化完成以后再返回。
如果网络连接异常，那么会立即返回，然后后台不断尝试
*/
func (m *MatrixTransport) Start() {
	if m.running {
		return
	}
	m.running = true
	// 2019.06.10 启动时主线程不等待,优化无网情况下的启动速度
	// 就算如果matrix连接不上,而主线程正常启动完成开始发送消息,也会被Send方法拒绝,重要消息会进入重发阶段,没有影响
	//wg := sync.WaitGroup{}
	//wg.Add(1)
	//firstStart := true
	go func() {
		for {
			var err error
			var store gomatrix.Storer
			var syncer *gomatrix.DefaultSyncer
			var homeServerValid = ""
			var homeServerURLValid = ""
			var matrixClientValid *gomatrix.MatrixClient
			firstSync := make(chan struct{}, 5)
			isFirstSynced := false
			for name, url := range m.servers {
				//如果用户制定了server,那么就用用户指定的,仅用于调试
				if len(params.UserSpecifiedMatrixServer) > 0 && url != params.UserSpecifiedMatrixServer {
					continue
				}
				var mcli *gomatrix.MatrixClient
				mcli, err = gomatrix.NewClient(url, "", "", PATHPREFIX0, m.log)
				if err != nil {
					continue
				}
				_, err = mcli.Versions()
				if err != nil {
					m.log.Warn(fmt.Sprintf("connect to martrix server err %s", err))
					continue
				}
				homeServerValid = name
				homeServerURLValid = url
				matrixClientValid = mcli
				break
			}
			if homeServerValid == "" || matrixClientValid == nil {
				m.log.Error("unable to find any reachable Matrix server")
				goto tryNext
			}
			m.servername = homeServerValid
			m.serverURL = homeServerURLValid
			m.matrixcli = matrixClientValid
			m.changeStatus(netshare.Connected)
			m.log.Debug(fmt.Sprintf("m.servername = %s", homeServerValid))
			// log in
			if err = m.loginOrRegister(); err != nil {
				m.log.Error(fmt.Sprintf("loginOrRegister err %s", err))
				goto tryNext
			}
			//initialize Filters/NextBatch/Rooms
			store = gomatrix.NewInMemoryStore()
			m.matrixcli.Store = store

			//handle the issue of discoveryroom,FOR TEST,temporarily retain this room
			if err = m.joinDiscoveryRoom(); err != nil {
				m.log.Error(fmt.Sprintf("joinDiscoveryRoom err %s", err))
				goto tryNext
			}
			//notify to server i am online（include the other participating servers）
			if err = m.matrixcli.SetPresenceState(&gomatrix.ReqPresenceUser{
				Presence:  ONLINE,
				StatusMsg: m.NodeDeviceType, //register device type to server
			}); err != nil {
				m.log.Error(fmt.Sprintf("SetPresenceState err %s", err))
				goto tryNext
			}
			//register receive-datahandle or other message received
			m.matrixcli.Store = store
			m.matrixcli.Syncer = gomatrix.NewDefaultSyncer(m.UserID, store)
			syncer = m.matrixcli.Syncer.(*gomatrix.DefaultSyncer)

			syncer.OnEventType(EventAddressRoom, m.onHandleAccountData)

			syncer.OnEventType("m.room.message", m.onHandleReceiveMessage)

			syncer.OnEventType("m.presence", m.onHandlePresenceChange)

			syncer.OnEventType("m.room.member", m.onHandleMemberShipChange)
			m.hasDoneStartCheck = false
			err = m.collectChannelInfo(m.db)
			if err != nil {
				m.log.Error("collectChannelInfo err %s", err)
			}
			go func() {
				for {
					err2 := m.matrixcli.Sync()

					if !isFirstSynced {
						isFirstSynced = true
						firstSync <- struct{}{}
					}
					if !m.running {
						return
					}
					if err2 != nil {
						m.log.Error(fmt.Sprintf("Matrix Sync return,err=%s ,will try agin..", err2))
						m.changeStatus(netshare.Reconnecting) //不能切换，sync终究会来结果,但在此期间会运行路由判断，因为下面sleep 5秒
						time.Sleep(time.Second * 1)
					} else {
						m.changeStatus(netshare.Connected)
					}
				}
			}()
			//wait for first sync complete
			<-firstSync
			isFirstSynced = true
			m.hasDoneStartCheck = true
			err = m.startupCheckAllParticipants()
			if err != nil {
				m.log.Error(fmt.Sprintf("startupCheckAllParticipants error %s", err))
			}
			//在启动的时候检测是否加入了一些不必要的聊天室,然后主动leave
			m.leaveUselessRoom()
			//if firstStart {
			//	firstStart = false
			//	wg.Done()
			//}
			return
		tryNext:
			//if firstStart {
			//	firstStart = false
			//	wg.Done()
			//}
			time.Sleep(time.Second * 1)
		}
	}()
	m.log.Trace(fmt.Sprintf("[Matrix] transport started peers=%s", utils.StringInterface(m.Peers, 7)))
	//wg.Wait()
	go m.loop()
}

func (m *MatrixTransport) loop() {
	for {
		select {
		case <-m.quitChan:
			return
		case j := <-m.jobChan:
			switch j.jobType {
			case jobSendMessage:
				m.doSend(j)
			case jobMessage:
				m.doHandleReceiveMessage(j)
			case jobPresence:
				m.doHandlePresenceChange(j)
			case jobMemberShip:
				m.doHandleMemberShipChange(j)
			case jobAccountData:
				m.doHandleAccountData(j)
			}
		}
	}
}

/*
------------------------------------------------------------------------------------------------------------------------
*/
//onHandleReceiveMessage push the message of some one send "account_data"
func (m *MatrixTransport) onHandleAccountData(event *gomatrix.Event) {
	//log.Trace(fmt.Sprintf("onHandleAccountData %s", utils.StringInterface(event, 5)))
	if m.stopreceiving || event.Type != EventAddressRoom {
		return
	}
	job := &matrixJob{
		jobType: jobAccountData,
		Data1:   event,
	}
	if !m.hasDoneStartCheck {
		m.doHandleAccountData(job)
	}
}
func (m *MatrixTransport) doHandleAccountData(job *matrixJob) {
	event := job.Data1.(*gomatrix.Event)
	//todo chen Dev版本没有任何AccountData消息，注意添加m.addPeerIfNotExist(address, true/false)
	//我关注的 peer 所在的聊天室
	for addrHex, roomIDInterface := range event.Content {
		roomID := roomIDInterface.(string)
		addr := common.HexToAddress(addrHex)
		m.lock.RLock()
		p := m.Peers[addr]
		m.lock.RUnlock()
		m.log.Debug(fmt.Sprintf("m.peers.doHandleAccountData = %s,peer=%s,p=%s", utils.StringInterface(m.Peers, 5), addrHex, utils.StringInterface(p, 5)))
		if p != nil {
			p.defaultMessageRoomID = roomID
		}
	}
}

/*
onHandleReceiveMessage handle text messages sent to listening rooms
收到消息
必须保证对应的 UserID 是验证过的,否则就不能认定此 ID 的有效性.
*/
func (m *MatrixTransport) onHandleReceiveMessage(event *gomatrix.Event) {
	//m.log.Trace(fmt.Sprintf("discoveryroomid=%s", m.discoveryroomid))
	if m.stopreceiving {
		return
	}
	if !m.hasDoneStartCheck {
		return
	}
	m.jobChan <- &matrixJob{
		jobType: jobMessage,
		Data1:   event,
	}
}
func (m *MatrixTransport) doHandleReceiveMessage(job *matrixJob) {
	event := job.Data1.(*gomatrix.Event)
	m.log.Trace(fmt.Sprintf("onHandleReceiveMessage %s", utils.StringInterface(event, 7)))
	msgTime := time.Unix(event.Timestamp/1000, 0)
	/*
		ignore any history message
		we assume that client's time and server's time are all correct.
		there are three types of history message:
		1. my first login
		2. my new joined room
		3. message during my disconnection
		todo fixme use better message filter
	*/
	if time.Now().Sub(msgTime) > time.Second*15 { //测试发现最新的消息（确定）的时间戳最高延迟是42分钟
		m.log.Trace(fmt.Sprintf("ignore message because of it's too early, now=%s,msgtime=%s,event=%s", time.Now(), msgTime, utils.StringInterface(event, 5)))
		return
	}
	if m.stopreceiving || event.Type != "m.room.message" {
		fmt.Println("m.room.message stop receive", event.Type, event.Sender)
		return
	}
	msgtype, ok := event.MessageType()
	if ok == false || msgtype != "m.text" {
		m.log.Error(fmt.Sprintf("onHandleReceiveMessage unkown msg type %s", utils.StringInterface(event, 3)))
		return
	}
	senderID := event.Sender
	if senderID == m.UserID {
		return
	}
	/*
		this userID is validated,
	*/
	if !m.isTrustedServerUser(senderID) {
		m.log.Warn(fmt.Sprintf("onHandleReceiveMessage receive msg %s,but userId is never validate", utils.StringInterface(event, 3)))
		return
	}
	peerAddress := m.userIDToAddress(senderID)
	m.lock.RLock()
	peer := m.Peers[peerAddress]
	m.lock.RUnlock()
	if peer == nil {
		m.temporaryPeers.addPeer(peerAddress, event.RoomID)
	} else {
		err := m.inviteIfPossible(senderID, event.RoomID)
		if err != nil {
			m.log.Error(fmt.Sprintf("inviteIfPossible in handle message ,err %s", err))
		}
	}
	data, ok := event.Body()
	if !ok || len(data) < 2 {
		m.log.Error(fmt.Sprintf("onHandleReceiveMessage data=%s,ok=%v", string(data), ok))
	}
	dataContent, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		m.log.Error(fmt.Sprintf("[Matrix]Receive unkown message %s", utils.StringInterface(event, 5)))
	} else {
		m.HandleMessage(peerAddress, dataContent)
		m.log.Trace(fmt.Sprintf("[Matrix]Receive message %s from %s", encoding.MessageType(dataContent[0]), utils.APex2(peerAddress)))
	}
}

// onHandleMemberShipChange Handle message when eventType==m.room.member and join all invited rooms
/*
需要验证如何处理 invite,members 这些事件
214e 初次创建聊天室,并邀请214e
收到加入 room邀请,不能立即得到邀请人的 display name,可以得到 UserId,所以最好是立即加入,
加入以后获取DisplayName, 如果验证 UserID和 DisplayName 不匹配,可以选择退出.
invite event: 214e invite  214e
	"invite": {
			"!oydnNtPOxjZjQbdmQi:transport01.smartmesh.cn": {
				"invite_state": {
					"events": [{
						"content": {
							"membership": "join",
							"avatar_url": null,
							"displayname": "214e-118e7ade8cf61531f0d1629febf299ab939f314f04391b59b3567314525b4bee77220115d9bb0a808b6036857dad8bf242cbf5c09b18e7778c95022a77058bdf1c"
						},
						"type": "m.room.member",
						"sender": "@0x214e7247a2757696ed2986a8331a9e27a330c750:transport01.smartmesh.cn",
						"state_key": "@0x214e7247a2757696ed2986a8331a9e27a330c750:transport01.smartmesh.cn"
					}, {
						"content": {
							"join_rule": "invite"
						},
						"type": "m.room.join_rules",
						"sender": "@0x214e7247a2757696ed2986a8331a9e27a330c750:transport01.smartmesh.cn",
						"state_key": ""
					}, {
						"origin_server_ts": 1537849280000,
						"sender": "@0x214e7247a2757696ed2986a8331a9e27a330c750:transport01.smartmesh.cn",
						"event_id": "$1537849280199DBBSm:transport01.smartmesh.cn",
						"unsigned": {
							"age": 147
						},
						"state_key": "@0x2262fb3113baa4152286ae601e4278ff8a46d205:transport02.smartmesh.cn",
						"content": {
							"membership": "invite",
							"avatar_url": null,
							"displayname": "2262-c4f2ce4138f95f14f93a60fb2a366842a58eb62c0634fcfa4a7ff0dbfad3bf126f6ed6b45c379399290147bddfe181a8408d8e655a94526969e15ad1800008ec1c"
						},
						"membership": "invite",
						"type": "m.room.member"
					}]
				}
			}
		}
*/
func (m *MatrixTransport) onHandleMemberShipChange(event *gomatrix.Event) {
	//m.log.Trace(fmt.Sprintf("onHandleMemberShipChange %s ", utils.StringInterface(event, 10)))
	if m.stopreceiving {
		return
	}
	if event.Type != "m.room.member" {
		return
	}
	userid := event.Sender
	//ignore event that trigged by myself
	if userid == m.UserID {
		return
	}
	membership, exists := event.ViewContent("membership")
	if !exists || membership == "" {
		m.log.Warn(fmt.Sprintf("receive m.room.member,but don't have mermberyship %s", utils.StringInterface(event, 10)))
		return
	}
	job := &matrixJob{
		jobType: jobMemberShip,
		Data1:   event,
		Data2:   membership,
	}
	//startup stage, do the job right now
	if !m.hasDoneStartCheck {
		m.doHandleMemberShipChange(job)
	} else {
		m.jobChan <- job
	}
}
func (m *MatrixTransport) doHandleMemberShipChange(job *matrixJob) {
	event := job.Data1.(*gomatrix.Event)
	membership := job.Data2.(string)
	userid := event.Sender
	if !m.isTrustedServerUser(userid) {
		m.log.Warn(fmt.Sprintf("receive invite,but i don't know this user %s", utils.StringInterface(event, 5)))
		return
	}
	/*
		The following membership states are specified:
			invite - The user has been invited to join a room, but has not yet joined it. They may not participate in the room until they join.
			join - The user has joined the room (possibly after accepting an invite), and may participate in it.
			leave - The user was once joined to the room, but has since left (possibly by choice, or possibly by being kicked).
			ban - The user has been banned from the room, and is no longer allowed to join it until they are un-banned from the room (by having their membership state set to a value other than ban).
			knock - This is a reserved word, which currently has no meaning.
	*/
	if membership == "invite" {
		/*
			{
							"origin_server_ts": 1537849280000,
							"sender": "@0x214e7247a2757696ed2986a8331a9e27a330c750:transport01.smartmesh.cn",
							"event_id": "$1537849280199DBBSm:transport01.smartmesh.cn",
							"unsigned": {
								"age": 147
							},
							"state_key": "@0x2262fb3113baa4152286ae601e4278ff8a46d205:transport02.smartmesh.cn",
							"content": {
								"membership": "invite",
								"avatar_url": null,
								"displayname": "2262-c4f2ce4138f95f14f93a60fb2a366842a58eb62c0634fcfa4a7ff0dbfad3bf126f6ed6b45c379399290147bddfe181a8408d8e655a94526969e15ad1800008ec1c"
							},
							"membership": "invite",
							"type": "m.room.member"
						}
		*/
		if event.StateKey == nil || *event.StateKey != m.UserID {
			//ignore invite message not send to me
			return
		}
		m.log.Trace(fmt.Sprintf("receive invite, event=%s", utils.StringInterface(event, 5)))
		go func() {
			//one must join to be able to get room alias
			var err error
			for i := 0; i < 3; i++ {
				serverName := getServerFromRoomID(event.RoomID, 0)
				m.log.Debug(fmt.Sprintf("serverName %s", serverName))

				_, err = m.matrixcli.JoinRoom(event.RoomID, "", nil) //todo chen servername=""?为空不能跨服务器
				if err != nil {
					if strings.Index(err.Error(), "already in the room") > -1 {
						break
					}
					m.log.Info(fmt.Sprintf("JoinRoom %s ,err %s, sleep 5 seconds and retry,times=%d", event.RoomID, err, i))
					time.Sleep(time.Second)
					continue
				} else {
					break
				}

			}
			if err != nil {
				m.log.Error(fmt.Sprintf("JoinRoom %s ,err %s", event.RoomID, err))
				return
			}
			peerAddress := m.userIDToAddress(userid)
			m.lock.RLock()
			peer := m.Peers[peerAddress]
			m.lock.RUnlock()
			if peer == nil {
				//maybe a peer want send secret request to me
				m.temporaryPeers.addPeer(peerAddress, event.RoomID)
			} else {
				err = m.inviteIfPossible(userid, event.RoomID)
				if err != nil {
					m.log.Error(fmt.Sprintf("inviteIfPossible %s to default room err %s", userid, err))
				}
			}
		}()
	} else if membership == "join" {
		err := m.inviteIfPossible(userid, event.RoomID)
		if err != nil {
			m.log.Error(fmt.Sprintf("inviteIfPossible %s to default room err %s", userid, err))
		}
	} else {
		//todo fix me handle leave event
	}
}

func getServerFromRoomID(roomid string, index int) string {
	//#photon_ropsten_discovery:transport01.smartmesh.cn
	if strings.ContainsAny(roomid, ":") {
		return strings.Split(roomid, ":")[index]
	}
	return ""
}

/*
onHandlePresenceChange handle events in this message, about changes of nodes and update AddressToPresence
{
	"content": {
		"status_msg": "other",
		"currently_active": true,
		"last_active_ago": 13,
		"presence": "online"
	},
	"type": "m.presence",
	"sender": "@0xf156aba37a64767769a96a0083f02f540e7856ab:transport01.smartmesh.cn"
}
*/
func (m *MatrixTransport) onHandlePresenceChange(event *gomatrix.Event) {
	//m.log.Trace(fmt.Sprintf("onHandlePresenceChange %s", utils.StringInterface(event, 5)))
	//m.log.Trace(fmt.Sprintf("address i want to know: %s", utils.StringInterface(m.Peers, 3)))
	if m.stopreceiving == true {
		return
	}
	if event.Type != "m.presence" {
		m.log.Error(fmt.Sprintf("onHandlePresenceChange receive unkonw event %s", utils.StringInterface(event, 5)))
		return
	}
	job := &matrixJob{
		jobType: jobPresence,
		Data1:   event,
	}

	//startup stage,do the job right now
	//经测试，经常存在m.presence事件延后，但是不会丢事件
	if m.hasDoneStartCheck {
		m.doHandlePresenceChange(job)
	} else {
		m.jobChan <- job
	}

}
func (m *MatrixTransport) doHandlePresenceChange(job *matrixJob) {
	event := job.Data1.(*gomatrix.Event)
	// parse address of message sender
	presence, exists := event.ViewContent("presence") //newest network status
	if !exists {
		return
	}
	// message sender
	userid := event.Sender
	//my self status change
	if userid == m.UserID {
		return
	}
	address := m.userIDToAddress(userid)
	if !m.isTrustedServerUser(userid) {
		m.log.Info(fmt.Sprintf("receive presence %s", utils.StringInterface(event, 5)))
		return
	}
	m.lock.RLock()
	peer, ok := m.Peers[address]
	m.lock.RUnlock()
	if !ok {
		//m.log.Trace(fmt.Sprintf("receive presence,but peer is unkown %s", utils.StringInterface(event, 5)))
		if presence == OFFLINE { //历史数据可导致错误
			m.temporaryPeers.removePeer(address)
		}
		return
	}
	if peer.isValidUserID(userid) && peer.setStatus(userid, presence) {
		// 节点上线通知所有已经挂起的通道
		if presence == ONLINE {
			m.WakeUp(address)
		}
		//device type
		deviceType, _ := event.ViewContent("status_msg") //newest network status
		peer.deviceType = deviceType
		m.log.Trace(fmt.Sprintf("%s matrix status changed status=%s,devicetype=%s", userid, peer.status, peer.deviceType))
	}
	m.log.Trace(fmt.Sprintf("peer %s status=%s,deviceType=%s", utils.APex2(address), peer.status, peer.deviceType))
}

//register new user on homeserver using application service
func (m *MatrixTransport) register(username, password string) (userID string, err error) {
	type reg struct {
		LocalPart   string `json:"localpart"`   //@someone:matrix.org someone is localpoart,matrix.org is domain
		DisplayName string `json:"displayname"` // displayname of this user
		Password    string `json:"password,omitempty"`
	}
	type regResp struct {
		AccessToken string `json:"access_token"`
		HomeServer  string `json:"home_server"`
		UserID      string `json:"user_id"`
	}
	regurl := fmt.Sprintf("%s/regapp/1/register", m.serverURL)
	//regurl := fmt.Sprintf("http://127.0.0.1:8009/regapp/1/register")
	userID = fmt.Sprintf("@%s:%s", username, m.servername)
	log.Trace(fmt.Sprintf("register user userid=%s", userID))
	req := &reg{
		LocalPart:   username,
		Password:    password,
		DisplayName: m.getUserDisplayName(userID),
	}
	resp := &regResp{}
	_, err = m.matrixcli.MakeRequest(http.MethodPost, regurl, req, resp)
	if err != nil {
		return
	}
	if resp.UserID != userID {
		err = fmt.Errorf("expect userid=%s,got=%s", userID, resp.UserID)
	}
	return
}

// loginOrRegister node login, if failed, register again then try login,
// displayname of nodes as the signature of userID
func (m *MatrixTransport) loginOrRegister() (err error) {
	loginok := false
	baseAddress := crypto.PubkeyToAddress(m.key.PublicKey)
	baseUsername := strings.ToLower(baseAddress.String())

	username := baseUsername
	password := hex.EncodeToString(m.dataSign([]byte(m.servername)))
	//password := "12345678"
	for i := 0; i < 5; i++ {
		var resplogin *gomatrix.RespLogin
		m.matrixcli.AccessToken = ""
		resplogin, err = m.matrixcli.Login(&gomatrix.ReqLogin{
			Type:     LOGINTYPE,
			User:     username,
			Password: password,
			DeviceID: "",
		})
		if err != nil {
			httpErr, ok := err.(gomatrix.HTTPError)
			if !ok { // network error,try again
				continue
			}
			if httpErr.Code == 403 { //Invalid username or password
				if i > 0 {
					m.log.Trace(fmt.Sprintf("couldn't sign in for matrix,trying register %s", username))
				}
				userID, rerr := m.register(username, password)
				if rerr != nil {
					return rerr
				}
				m.matrixcli.UserID = userID
				continue
			}
		} else {
			//cache the node's and report the UserID and AccessToken to matrix
			m.matrixcli.SetCredentials(resplogin.UserID, resplogin.AccessToken)
			m.UserID = resplogin.UserID
			m.NodeAddress = baseAddress
			loginok = true
			break
		}
	}
	if !loginok {
		err = fmt.Errorf("could not register or login")
		return
	}
	return err
}

// makeRoomAlias name room's alias
func (m *MatrixTransport) makeRoomAlias(thepart string, network string) string {
	return ROOMPREFIX + ROOMSEP + network + ROOMSEP + thepart
}

func (m *MatrixTransport) getUserDisplayName(userID string) string {
	sig := m.dataSign([]byte(userID))
	return fmt.Sprintf("%s-%s", utils.APex2(m.NodeAddress), hex.EncodeToString(sig))
}

// dataSign 签名数据
// dataSign signature data
func (m *MatrixTransport) dataSign(data []byte) (signature []byte) {
	signature, err := utils.SignData(m.key, data)
	if err != nil {
		m.log.Error(fmt.Sprintf("SignData err %s", err))
		return nil
	}
	return
}

func (m *MatrixTransport) userIDToAddress(userID string) common.Address {
	//check grammar of user ID
	_match := ValidUserIDRegex.MatchString(userID)
	if _match == false {
		m.log.Warn(fmt.Sprintf("UserID %s, format error", userID))
		return utils.EmptyAddress
	}
	addressHex, err := extractUserLocalpart(userID) //"@myname:photon.org:cy"->"myname"
	if err != nil {
		m.log.Error(fmt.Sprintf("extractUserLocalpart err %s", err))
		return utils.EmptyAddress
	}
	var addrmuti = regexp.MustCompile(`^(0x[0-9a-f]{40})`)
	addrlocal := addrmuti.FindString(addressHex)
	if addrlocal == "" {
		err = fmt.Errorf("%s not match our userid rule", userID)
		return utils.EmptyAddress
	}
	address := common.HexToAddress(addrlocal)
	return address
}

// findOrCreateRoomByAddress Get the room(*object) info by the node address.
// If no communication(peer-to-peer) room found yet,then "SearchUserDirectory" and create a temporary(unnamed) room for communication,invite the node finally.
func (m *MatrixTransport) findOrCreateRoomByAddress(address common.Address, hasChannel bool) (roomID string, users []*gomatrix.UserInfo, err error) {
	if m.stopreceiving {
		err = errors.New("already topped")
		return
	}
	//The following is the case where peer-to-peer communication room does not exist.
	var addressOfPairs = ""
	strPairs := []string{hex.EncodeToString(m.NodeAddress.Bytes()), hex.EncodeToString(address.Bytes())}
	sort.Strings(strPairs)
	addressOfPairs = strings.Join(strPairs, "_") //format "0cccc_dddd"
	var roomName string
	if hasChannel {
		roomName = m.makeRoomAlias(addressOfPairs, networkPartHasChannel)
	} else {
		roomName = m.makeRoomAlias(addressOfPairs, networkPartNoChannel)
	}
	//try to get user-infos of communication with "account_data" from homeserver include the other participating servers.
	users = m.getAllPossibleUserID(address)

	//Join a room that connot be found by search_room_directory
	roomID, err = m.getUnlistedRoom(roomName, users)
	m.log.Info(fmt.Sprintf("CHANNEL ROOM,peer_address=%s room=%s", address.String(), roomID))
	return
}

// getUnlistedRoom get a conversation room that cannnot be found by search_room_directory.
// If the room is not exist and create a unnamed room for communication,invite the node finally.
// This process of join-create-join-room may be repeated 3 times(network delay)
func (m *MatrixTransport) getUnlistedRoom(roomname string, users []*gomatrix.UserInfo) (roomID string, err error) {
	roomNameFull := "#" + roomname + ":" + m.servername
	var inviteesUids []string
	for _, user := range users {
		inviteesUids = append(inviteesUids, user.UserID)
	}
	req := &gomatrix.ReqCreateRoom{
		Invite:     inviteesUids,
		Visibility: "public",
		Preset:     "public_chat",
	}
	if true {
		req.Visibility = "public"
		req.Preset = "public_chat"
	} else {
		req.Visibility = "private"
		req.Preset = "trusted_private_chat"
	}
	unlistedRoomid := ""
	for i := 0; i < 5; i++ {
		var respJoinRoom *gomatrix.RespJoinRoom
		respJoinRoom, err = m.matrixcli.JoinRoom(roomNameFull, m.servername, nil)
		if err != nil {
			m.log.Error(fmt.Sprintf("JoinRoom %s error: %s,respJoinRoom: %s", roomname, err, utils.StringInterface(respJoinRoom, 5)))
			req.RoomAliasName = roomname
			_, err = m.matrixcli.CreateRoom(req)
			if err != nil {
				m.log.Info(fmt.Sprintf("Room %s not found,trying to create it. but fail %s", roomname, err))
			}
			time.Sleep(200)
			continue
		} else {
			unlistedRoomid = respJoinRoom.RoomID
			m.log.Info(fmt.Sprintf("Room joined successfully,room=%s", unlistedRoomid))
			break
		}
	}
	//if can't join nor create, create an unnamed one
	if unlistedRoomid == "" {
		var createdRoom *gomatrix.RespCreateRoom
		req.RoomAliasName = ""
		createdRoom, err = m.matrixcli.CreateRoom(req)
		if err != nil {
			return
		}
		unlistedRoomid = createdRoom.RoomID
	}
	return unlistedRoomid, nil
}
func (m *MatrixTransport) getAllPossibleUserID(address common.Address) (users []*gomatrix.UserInfo) {
	for s := range m.trustServers {
		users = append(users, &gomatrix.UserInfo{
			UserID: fmt.Sprintf("@%s:%s", strings.ToLower(address.String()), s),
		})
	}
	return users
}
func (m *MatrixTransport) handleNewPartner(p *MatrixPeer) (err error) {
	roomID, _, err := m.findOrCreateRoomByAddress(p.address, true)
	if err != nil {
		return
	}
	p.defaultMessageRoomID = roomID
	address2Room := make(map[common.Address]string)
	m.lock.RLock()
	for addr, peer := range m.Peers {
		//有可能为所有的通道实现分配了Peers,但是还没有来得及创建聊天室
		if peer.defaultMessageRoomID == "" {
			continue
		}
		address2Room[addr] = peer.defaultMessageRoomID
	}
	m.lock.RUnlock()
	return m.matrixcli.SetAccountData(m.UserID, EventAddressRoom, address2Room)
}

/*
与我有通道的地址,是需要长期保持有聊天室的.
1. 如果这些地址没有相应的聊天室,应该创建
2. 对于那些没在缺省聊天室中的UserID, 发出邀请.
3. 如果这些地址在线状态未知,需要通过GET /_matrix/client/r0/presence/{userId}/status来手动查询
*/
func (m *MatrixTransport) startupCheckAllParticipants() error {
	errBuf := new(bytes.Buffer)
	for _, p := range m.Peers {
		err := m.startupCheckOneParticipant(p)
		if err != nil {
			fmt.Fprintf(errBuf, "startupCheckOneParticipant %s err %s\n", utils.APex2(p.address), err)
		}
	}
	errStr := string(errBuf.Bytes())
	if len(errStr) > 0 {
		return errors.New(errStr)
	}
	return nil
}

//startupCheckOneParticipant do a lot of things for a peer
func (m *MatrixTransport) startupCheckOneParticipant(p *MatrixPeer) error {
	errBuf := new(bytes.Buffer)
	//don't have room with the partner
	if p.defaultMessageRoomID == "" {
		err := m.handleNewPartner(p)
		if err != nil {
			fmt.Fprintf(errBuf, "handleNewPartner for %s,err=%s\n", utils.APex2(p.address), err)
			return err
		}
		m.log.Debug(fmt.Sprintf("handleNewPartner ok:%s", utils.StringInterface(p, 5)))
	}
	/*
		1.list default message room's member
		2. save to the MatrixPeer
	*/
	respJoinedMembers, err := m.matrixcli.JoinedMembers(p.defaultMessageRoomID)
	if err != nil {
		fmt.Fprintf(errBuf, "JoinedMembers for peer:%s,room:%s,err=%s", utils.APex2(p.address), p.defaultMessageRoomID, err)
	} else {
		var users []*gomatrix.UserInfo
		for userID, joined := range respJoinedMembers.Joined {
			if m.userIDToAddress(userID) == m.NodeAddress {
				continue
			}
			u := &gomatrix.UserInfo{
				UserID:      userID,
				DisplayName: *joined.DisplayName,
			}
			if joined.DisplayName != nil {
				u.DisplayName = *joined.DisplayName
			}
			if joined.AvatarURL != nil {
				u.AvatarURL = *joined.AvatarURL
			}
			if m.isTrustedServerUser(userID) {
				users = append(users, u)
			} else {
				fmt.Fprintf(errBuf, "JoinedMembers %s verify signature err %s,displayname=%s", u.UserID, err, u.DisplayName)
			}
		}
		p.updateUsers(users)
	}
	//update peer's presence status
	if p.status != peerStatusOnline {
		for _, u := range p.candidateUsers {
			var presenceResponse *gomatrix.RespPresenceUser
			presenceResponse, err = m.matrixcli.GetPresenceState(u.UserID)
			if err != nil {
				fmt.Fprintf(errBuf, "GetPresenceState for %s,err=%s\n", u.UserID, err)
			} else {
				p.setStatus(u.UserID, presenceResponse.Presence)
				//stop check if one online userid found
				if p.status == peerStatusOnline {
					break
				}
			}
		}
	}
	/*
		there maybe Matrix Users should in the default message room,but I don't know yet.
		for example.
		1. alice login in server1 and join the default message room.
		2. then alice log off and relogin to server2 using another userid
		3. alice won't join my default message room only if she receives my invites
	*/

	users := m.getAllPossibleUserID(p.address)

	for _, u := range users {
		//it's a user should in the default message room,but it isn't in now.
		if !p.isValidUserID(u.UserID) {
			_, err = m.matrixcli.InviteUser(p.defaultMessageRoomID, &gomatrix.ReqInviteUser{
				UserID: u.UserID,
			})
			if err != nil {
				fmt.Fprintf(errBuf, "InviteUser %s to room %s err %s", u.UserID, p.defaultMessageRoomID, err)
			}
		}
	}
	errStr := string(errBuf.Bytes())
	if len(errStr) > 0 {
		return errors.New(errStr)
	}
	return nil
}

func getSignatureFromDisplayName(displayName string) (signature []byte, err error) {
	ss := strings.Split(displayName, "-")
	//userAddr-Signature
	if len(ss) != 2 {
		err = fmt.Errorf("display name format error %s", displayName)
		return
	}
	//signature length is 130
	if len(ss[1]) != 130 {
		err = fmt.Errorf("signature error")
	}
	signature, err = hex.DecodeString(ss[1])
	return
}

func (m *MatrixTransport) inviteIfPossible(userID string, eventRoom string) error {
	peerAddress := m.userIDToAddress(userID)
	/*
		if this address has channel with me ,it may be login with another UserID,
		so I need update my info
	*/
	m.lock.RLock()
	peer := m.Peers[peerAddress]
	m.lock.RUnlock()
	if peer == nil {
		return nil
	}
	/*
			invite a user should meet the following pre requests:
		1. it's valid userID, when   call this function, the userID is already verified
		2.this user not in the default message room
			isValidUserID checks userID not in the default message room right now.
			hasDoneStartCheck: because of in startup ,I don't know who are in the default message room
	*/
	if peer.defaultMessageRoomID != "" && peer.defaultMessageRoomID != eventRoom && !peer.isValidUserID(userID) {
		//invite this user to the default room
		_, err := m.matrixcli.InviteUser(peer.defaultMessageRoomID, &gomatrix.ReqInviteUser{
			UserID: userID,
		})
		//if has already startup ,we should known all user status
		var i int
		for i = 0; i < 5; i++ {
			if err != nil && m.hasDoneStartCheck {
				if strings.Index(err.Error(), "already in the room") > -1 {
					break
				}
				m.log.Info(fmt.Sprintf("InviteUser %s ,err %s,and retry...", userID, err))
				time.Sleep(time.Second)
				continue
			} else {
				break
			}
		}
		if i == 5 {
			return fmt.Errorf("InviteUser perr=%s,room=%s err=%s",
				utils.APex2(peerAddress), peer.defaultMessageRoomID, err)
		}

	}
	if peer.defaultMessageRoomID == eventRoom {
		//this user is joinning in the default message room, a new user
		if !peer.isValidUserID(userID) {
			//update presence if possible
			peer.updateUsers([]*gomatrix.UserInfo{{UserID: userID}})
			presence, err := m.matrixcli.GetPresenceState(userID)
			if err == nil {
				peer.setStatus(userID, presence.Presence)
				peer.deviceType = presence.StatusMsg
			} else {
				return err
			}
		}
	}
	return nil
}

// extractUserLocalpart Extract user name from user ID
func extractUserLocalpart(userID string) (string, error) {
	if len(userID) == 0 || userID[0] != '@' {
		return "", fmt.Errorf("%s is not a valid user id", userID)
	}
	return strings.SplitN(userID[1:], ":", 2)[0], nil
}

// extractUserInfo Extract user name from user ID
func extractUserInfo(userID string) (localPart, domain string, err error) {
	if len(userID) == 0 || userID[0] != '@' {
		err = fmt.Errorf("%s is not a valid user id", userID)
		return
	}
	ss := strings.SplitN(userID[1:], ":", 2)
	if len(ss) != 2 {
		err = fmt.Errorf("%s is not a valid user id", userID)
		return
	}
	localPart = ss[0]
	domain = ss[1]
	return
}

// validate_userid_signature
func validateUseridSignature(user *gomatrix.UserInfo) (address common.Address, err error) {
	//displayname should be an address in the self._userid_re format
	_match := ValidUserIDRegex.MatchString(user.UserID)
	if _match == false {
		err = fmt.Errorf("validate user info failed")
		return
	}
	_address, err := extractUserLocalpart(user.UserID) //"@myname:photon.org:cy"->"myname"
	if err != nil {
		return
	}
	var addrmuti = regexp.MustCompile(`^(0x[0-9a-f]{40})`)
	addrlocal := addrmuti.FindString(_address)
	if addrlocal == "" {
		err = fmt.Errorf("%s not match our userid rule", user.UserID)
		return
	}
	address = common.HexToAddress(addrlocal)
	signature, err := getSignatureFromDisplayName(user.DisplayName)
	if err != nil {
		return
	}
	recovered, err := utils.Ecrecover(utils.Sha3([]byte(user.UserID)), signature)
	if err != nil {
		return
	}
	if !bytes.Equal(recovered[:], address[:]) {
		err = fmt.Errorf("validate %s failed", user.UserID)
		return
	}
	return
}

/* joinDiscoveryRoom : check discoveryroom if not exist, then create a new one.
client caches all memebers of this room, and invite nodes checked from this room again.
todo 需要找到一个可靠的方式来移除DiscoveryRoom,
目前不能移除DiscoveryRoom是因为PathFinder需要依赖DiscoveryRoom来发现节点的上线下线,正常的Matrix通信已经可以做到不依赖DiscoveryRoom了
发现聊天室设计目标主要是让节点之间能够找到对方,主要是通过Matrix的Search方式找到相关UserID以及指导这些UserID的上线下线状态.
但是目前来说这些都不再需要,
1. Search可以通过@<address>:domain方式自己生产所有可能的UserID
2. 节点上线下线通知,有Channel的节点直接创建私有的聊天室来解决上线下线状态问题
*/
func (m *MatrixTransport) joinDiscoveryRoom() (err error) {
	//read discovery room'name and fragment from "params-settings"
	// combine discovery room's alias
	discoveryRoomAlias := m.makeRoomAlias(ALIASFRAGMENT, NETWORKNAME)
	discoveryRoomAliasFull := "#" + discoveryRoomAlias + ":" + DISCOVERYROOMSERVER
	m.discoveryroom = ""
	// this node join the discovery room, if not exist, then create.
	for i := 0; i < 10; i++ {
		var respJoinRoom *gomatrix.RespJoinRoom
		var respCreateRoom *gomatrix.RespCreateRoom
		respJoinRoom, err = m.matrixcli.JoinRoom(discoveryRoomAliasFull, m.servername, nil)
		if err != nil {
			//if Room doesn't exist and then create the room(this is the node's resposibility)
			if m.servername != DISCOVERYROOMSERVER {
				break
			}
			//try to create the discovery room
			respCreateRoom, err = m.matrixcli.CreateRoom(&gomatrix.ReqCreateRoom{
				RoomAliasName: discoveryRoomAlias,
				Preset:        CHATPRESET,
				Visibility:    "public",
			})
			if err != nil {
				m.log.Error("can't create a discovery room,try again")
				time.Sleep(time.Second)
				continue
			}
			m.discoveryroom = respCreateRoom.RoomID
			time.Sleep(time.Second)
			continue
		} else {
			m.discoveryroom = respJoinRoom.RoomID
			break
		}
	}
	//exit if join room failed
	if m.discoveryroom == "" {
		err = fmt.Errorf("Discovery room {%s} not found and can't be created on a federated homeserver {%s}", discoveryRoomAliasFull, m.servername)
		m.log.Error(err.Error())
		return
	}
	m.log.Info(fmt.Sprintf("Join Discovery room %s success", discoveryRoomAliasFull))
	return nil
}

/*
Matrix运行一段时间以后,一个账户必定会累积不少无用的聊天室.
哪些聊天室可以移除?
1. 与我没有channel的账户,并且没有通信一天以上(考虑utc问题)的,基本上可以认为近期不会再用
2. 与我有channel,但是对方在好几个服务器上都登陆,并且创建了聊天室,应该只保留一个即可.//这个暂时保留
哪些聊天室必须保留?
1. 正在进行交易的 如何判断呢,最近一天只要有活动(有人加入,有人退出,发生过任何事件)
2. 与我有通道的节点,并且是我创建的聊天室,
3. 与我有通道的节点,聊天室是对方创建的,但是对方正处于活跃中.
假设我是A,我与10个不同的账户有通道,那么可能至少保有20个聊天室. 如果这十个账户和我在不同的matrix服务器上
*/
func (m *MatrixTransport) leaveUselessRoom() {
	rooms := m.matrixcli.Store.LoadRoomOfAll()
	for roomID, room := range rooms {
		//log.Trace(fmt.Sprintf("for leave room %s,%s", roomID, utils.StringInterface(room, 2)))
		//discovery room必须保留
		if roomID == m.discoveryroom {
			continue
		}
		if m.isUseLessRoom(room) {
			//m.log.Info(fmt.Sprintf("try to leave room %s,because of useless ", utils.StringInterface(room, 5)))
			_, err := m.matrixcli.LeaveRoom(roomID)
			if err != nil {
				m.log.Error(fmt.Sprintf("leave room %s err %s,room=%s", roomID, err, utils.StringInterface(room, 5)))
			} else {
				/* forget调用无用,经陈云测试
				_, err = m.matrixcli.ForgetRoom(roomID)
				if err != nil {
					m.log.Error(fmt.Sprintf("forget room %s err %s,room=%s", roomID, err, utils.StringInterface(room, 5)))
				}*/
			}
		}
	}
}

//	//#photon_y_37bd76c0187ebc21e3fd3d474d83810bb495a518_4533775cfd13a2b07bf910c04d2038fd028ff73c:transport02.smartmesh.cn"
func splitRoomAlias(alias string) (prefix, isChannel string, addr1, addr2 common.Address, err error) {
	ss := strings.Split(alias, ":")
	if len(ss) != 2 {
		err = fmt.Errorf("room alias %s format error", alias)
		return
	}
	s := ss[0]
	if len(s) < len(ROOMPREFIX) {
		err = fmt.Errorf("room alias %s local part length error", alias)
		return
	}
	ss = strings.Split(s[1:], ROOMSEP)
	if len(ss) != 4 {
		err = fmt.Errorf("room alias %s local part format error", alias)
		return
	}
	prefix = ss[0]
	isChannel = ss[1]
	b, err := hex.DecodeString(ss[2])
	if err != nil {
		err = fmt.Errorf("room alias %s addr1 decode error", alias)
		return
	}
	addr1 = common.BytesToAddress(b)
	b, err = hex.DecodeString(ss[3])
	if err != nil {
		err = fmt.Errorf("room alias %s addr2 decode error", alias)
		return
	}
	addr2 = common.BytesToAddress(b)
	return
}
func (m *MatrixTransport) isUseLessRoom(r *gomatrix.Room) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	for _, p := range m.Peers {
		if p.defaultMessageRoomID == r.ID {
			return false
		}
	}
	//#photon_y_37bd76c0187ebc21e3fd3d474d83810bb495a518_4533775cfd13a2b07bf910c04d2038fd028ff73c:transport02.smartmesh.cn"
	//"#photon_ropsten_discovery:transport01.smartmesh.cn"
	//	return ROOMPREFIX + ROOMSEP + network + ROOMSEP + thepart
	prefix, isChannel, addr1, addr2, err := splitRoomAlias(r.Alias)
	if err != nil {
		return true //格式不对
	}
	if prefix != ROOMPREFIX {
		return true
	}

	if addr1 != m.NodeAddress && addr2 != m.NodeAddress {
		//我没有参与的聊天室,肯定是有问题的
		return true
	}
	//临时通道还是固定通道
	if isChannel == networkPartNoChannel {
		//最近一天没有任何活动
		if getRoomLastestActiveTime(r).Add(time.Hour * 24).Before(time.Now()) {
			return true
		}
	}
	//其他的都下来.
	return false
}
func getRoomLastestActiveTime(r *gomatrix.Room) time.Time {
	t := time.Now().Add(time.Hour * (-24)) //最早认为是一天前,再早的活动就忽略
	for _, events := range r.State {
		for _, e := range events {
			msgTime := time.Unix(e.Timestamp/1000, 0)
			if msgTime.After(t) {
				t = msgTime
			}
		}
	}
	return t
}
