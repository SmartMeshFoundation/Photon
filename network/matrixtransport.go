package network

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/SmartMeshFoundation/SmartRaiden/network/gomatrix"

	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/channel/channeltype"
	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/netshare"
	"github.com/SmartMeshFoundation/SmartRaiden/network/xmpptransport"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
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
	ROOMPREFIX = "smartraiden"
	// ROOMSEP with ',' to separate room name's part
	ROOMSEP = "_"
	// PATHPREFIX0 the lastest matrix client api version
	PATHPREFIX0 = "/_matrix/client/r0"
	// AUTHTYPE login identity as dummy
	AUTHTYPE = "m.login.dummy"
	// LOGINTYPE login type we used
	LOGINTYPE = "m.login.password"
	// CHATPRESET the type of chat=public
	CHATPRESET = "public_chat"
	//EventAddressRoom is user defined event type
	EventAddressRoom = "network.smartraiden.rooms"
)

// MatrixTransport represents a matrix transport Instantiation
type MatrixTransport struct {
	matrixcli             *gomatrix.MatrixClient //the instantiated matrix
	servername            string                 //the homeserver's name
	running               bool                   //running status
	stopreceiving         bool                   //Whether to stop accepting(data)
	key                   *ecdsa.PrivateKey      //key
	NodeAddress           common.Address
	protocol              ProtocolReceiver
	discoveryroomid       string //the room's ID of sys pre-configured ("![RoomIdData]:[ServerName]")
	Peers                 map[common.Address]*MatrixPeer
	temporaryAddress2Room map[common.Address]string     //temporary room to
	validatedUsers        map[string]*gomatrix.UserInfo //this userId and display name is validated
	UserID                string                        //the current user's ID(@kitty:thisserver)
	NodeDeviceType        string
	avatarurl             string
	log                   log.Logger
	statusChan            chan netshare.Status
	removePeerChan        chan common.Address
	status                netshare.Status
	servers               map[string]string
	db                    xmpptransport.XMPPDb
	hasDoneStartCheck     bool
}

var (
	// ValidUserIDRegex user ID 's format
	ValidUserIDRegex = regexp.MustCompile(`^@(0x[0-9a-f]{40})(?:\.[0-9a-f]{8})?(?::.+)?$`) //(`^[0-9a-z_\-./]+$`)
	//NETWORKNAME which network is used
	NETWORKNAME = params.NETWORKNAME
	//ALIASFRAGMENT the terminal part of alias
	ALIASFRAGMENT = params.AliasFragment
	//DISCOVERYROOMSERVER discovery room server name
	DISCOVERYROOMSERVER = params.DiscoveryServer
)

// NewMatrixTransport init matrix
func NewMatrixTransport(logname string, key *ecdsa.PrivateKey, devicetype string, servers map[string]string) *MatrixTransport {
	mtr := &MatrixTransport{
		running:               false,
		stopreceiving:         false,
		NodeAddress:           crypto.PubkeyToAddress(key.PublicKey),
		key:                   key,
		Peers:                 make(map[common.Address]*MatrixPeer),
		validatedUsers:        make(map[string]*gomatrix.UserInfo),
		temporaryAddress2Room: make(map[common.Address]string),
		NodeDeviceType:        devicetype,
		log:                   log.New("matrix", logname),
		avatarurl:             "", // charge rule
		statusChan:            make(chan netshare.Status, 10),
		removePeerChan:        make(chan common.Address, 10),
		status:                netshare.Disconnected,
		servers:               servers,
	}
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

//todo should refactor to make db set in constructor
func (m *MatrixTransport) setDB(db xmpptransport.XMPPDb) error {
	m.db = db
	return nil
}

func (m *MatrixTransport) addPeerIfNotExist(peer common.Address, hasChannel bool) bool {
	_, ok := m.Peers[peer]
	if ok {
		return false
	}
	m.Peers[peer] = NewMatrixPeer(peer, hasChannel, m.removePeerChan)
	return true
}
func (m *MatrixTransport) isThePeerIWantMessage(peer common.Address, userID string) bool {
	u, ok := m.Peers[peer]
	if !ok {
		return false
	}
	return u.isValidUserID(userID)
}
func (m *MatrixTransport) isUserIDValidated(userID string) bool {
	return m.validatedUsers[userID] != nil
}
func (m *MatrixTransport) validateAndUpdateUser(user *gomatrix.UserInfo) error {
	oldUser := m.validatedUsers[user.UserID]
	//user display name is signature of user id
	if oldUser != nil {
		//exists user ,but it doesn't have display name or display name exactly match.
		if len(oldUser.DisplayName) == 0 {
			oldUser.DisplayName = user.DisplayName
			return nil
		}
		if oldUser.DisplayName == user.DisplayName {
			return nil
		}
		return fmt.Errorf("displayname already exists, old=%s,new=%s, validatedusers=%s",
			oldUser.DisplayName,
			user.DisplayName,
			utils.StringInterface(m.validatedUsers, 3),
		)
	}
	//a new user ,validate user id is valid or not
	if len(user.DisplayName) == 0 {
		return fmt.Errorf("validateAndUpdateUser")
	}
	_, err := validateUseridSignature(user)
	if err != nil {
		return err
	}
	m.validatedUsers[user.UserID] = user
	return nil
}
func (m *MatrixTransport) knownThisPeer(peer common.Address) bool {
	_, ok := m.Peers[peer]
	return ok
}

// collectChannelInfo subscribe status change
func (m *MatrixTransport) collectChannelInfo(db xmpptransport.XMPPDb) error {
	cs, err := db.GetChannelList(utils.EmptyAddress, utils.EmptyAddress)
	if err != nil {
		return err
	}
	for _, c := range cs {
		m.addPeerIfNotExist(c.PartnerAddress(), true)
	}
	db.RegisterNewChannellCallback(func(c *channeltype.Serialization) (remove bool) {
		if m.addPeerIfNotExist(c.PartnerAddress(), true) {
			err := m.startupCheckOneParticipant(m.Peers[c.PartnerAddress()])
			if err != nil {
				log.Error(fmt.Sprintf("handleNewPartner for %s,err %s", utils.APex2(c.PartnerAddress()), err))
			}
		}
		return false
	})
	db.RegisterChannelStateCallback(func(c *channeltype.Serialization) (remove bool) {
		//todo mark peer to delete
		log.Info(fmt.Sprintf("matrix User should be removed"))
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
	m.matrixcli.SetPresenceState(&gomatrix.ReqPresenceUser{
		Presence: OFFLINE,
	})
	m.matrixcli.StopSync()
	if _, err := m.matrixcli.Logout(); err != nil {
		m.log.Error("[Matrix] Logout failed")
	}
}

// StopAccepting stop receive message and wait
func (m *MatrixTransport) StopAccepting() {
	m.stopreceiving = true
}

// NodeStatus gets Node states of network, if check self node, `status` is not always be true instead it switches according to server handshake signal.
func (m *MatrixTransport) NodeStatus(addr common.Address) (deviceType string, isOnline bool) {
	if m.matrixcli == nil {
		return "", false
	}
	u, ok := m.Peers[addr]
	if !ok {
		return "", false
	}
	return u.deviceType, u.status == peerStatusOnline
}

// Send send message
func (m *MatrixTransport) Send(receiverAddr common.Address, data []byte) error {
	var err error
	if !m.running || len(data) == 0 {
		return fmt.Errorf("[Matrix]Send failed,matrix not running or send data is null")
	}
	m.log.Trace(fmt.Sprintf("sendmsg  %s", utils.StringInterface(m.Peers, 7)))
	p := m.Peers[receiverAddr]
	var roomID string
	if p == nil {
		roomID = m.temporaryAddress2Room[receiverAddr]
		if roomID == "" {
			roomID, _, err = m.findOrCreateRoomByAddress(receiverAddr, true)
			if err != nil || roomID == "" {
				return fmt.Errorf("[Matrix]Send failed,cann't find the peer address")
			}
			m.temporaryAddress2Room[receiverAddr] = roomID
		}
	} else {
		roomID = p.defaultMessageRoomID
	}
	_data := base64.StdEncoding.EncodeToString(data)
	_, err = m.matrixcli.SendText(roomID, _data)
	if err != nil {
		m.log.Error(fmt.Sprintf("[matrix]send failed to %s, message=%s", utils.APex2(receiverAddr), encoding.MessageType(data[0])))
		return err
	}
	m.log.Trace(fmt.Sprintf("[Matrix]Send to %s, message=%s", utils.APex2(receiverAddr), encoding.MessageType(data[0])))
	return nil
}

// Start matrix
func (m *MatrixTransport) Start() {
	if m.running {
		return
	}
	m.running = true
	wg := sync.WaitGroup{}
	wg.Add(1)
	firstStart := true
	go func() {
		for {
			var err error
			var store gomatrix.Storer
			var syncer *gomatrix.DefaultSyncer
			var homeServerValid = ""
			var matrixClientValid *gomatrix.MatrixClient
			firstSync := make(chan struct{}, 5)
			isFirstSynced := false
			for name, url := range m.servers {
				var mcli *gomatrix.MatrixClient
				homeserverurl := url
				homeservername := name
				mcli, err = gomatrix.NewClient(homeserverurl, "", "", PATHPREFIX0, m.log)
				if err != nil {
					continue
				}
				_, err = mcli.Versions()
				if err != nil {
					m.log.Error(fmt.Sprintf("Could not connect to requested server %s,and retrying,err %s", homeserverurl, err))
					continue
				}
				homeServerValid = homeservername
				matrixClientValid = mcli
				break
			}
			if homeServerValid == "" || matrixClientValid == nil {
				errinfo := "unable to find any reachable Matrix server"
				m.log.Error(errinfo)
				goto tryNext
			}
			m.servername = homeServerValid
			m.matrixcli = matrixClientValid
			m.changeStatus(netshare.Connected)

			err = m.collectChannelInfo(m.db)
			if err != nil {
				m.log.Warn("collectChannelInfo err %s", err)
			}
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
						m.log.Error(fmt.Sprintf("Matrix Sync return,err=%s ,will try agin..", err))
						m.changeStatus(netshare.Reconnecting)
						time.Sleep(time.Second * 5)
					} else {
						m.changeStatus(netshare.Connected)
					}
				}
			}()
			//wait for first sync complete
			<-firstSync
			isFirstSynced = true
			if m.status == netshare.Connected {
				if !m.hasDoneStartCheck {
					m.hasDoneStartCheck = true
					err = m.startupCheckAllParticipants()
					if err != nil {
						m.log.Error(fmt.Sprintf("startupCheckAllParticipants error %s", err))
					}
				}
			}
			if firstStart {
				firstStart = false
				wg.Done()
			}
			return
		tryNext:
			if firstStart {
				firstStart = false
				wg.Done()
			}
			time.Sleep(time.Second * 5)
		}
	}()
	m.log.Trace(fmt.Sprintf("[Matrix] transport started peers=%s", utils.StringInterface(m.Peers, 7)))
	wg.Wait()
}

/*
------------------------------------------------------------------------------------------------------------------------
*/
//onHandleReceiveMessage push the message of some one send "account_data"
func (m *MatrixTransport) onHandleAccountData(event *gomatrix.Event) {
	log.Trace(fmt.Sprintf("onHandleAccountData %s", utils.StringInterface(event, 5)))
	if m.stopreceiving || event.Type != EventAddressRoom {
		return
	}
	if !m.hasDoneStartCheck {
		//我关注的 peer 所在的聊天室
		for addrHex, roomIDInterface := range event.Content {
			roomID := roomIDInterface.(string)
			addr := common.HexToAddress(addrHex)
			p := m.Peers[addr]
			if p != nil {
				p.defaultMessageRoomID = roomID
			}
		}
	}
}

/*
onHandleReceiveMessage handle text messages sent to listening rooms
收到消息
必须保证对应的 UserID 是验证过的,否则就不能认定此 ID 的有效性.
*/
func (m *MatrixTransport) onHandleReceiveMessage(event *gomatrix.Event) {
	m.log.Trace(fmt.Sprintf("discoveryroomid=%s", m.discoveryroomid))
	if event.RoomID == m.discoveryroomid {
		//ignore any message sent to discovery room.
		return
	}
	if !m.hasDoneStartCheck {
		return
	}
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
	if time.Now().Sub(msgTime) > time.Second*10 {
		m.log.Trace(fmt.Sprintf("ignore message because of it's too early, now=%s,msgtime=%s", time.Now(), msgTime))
		return
	}
	if m.stopreceiving || event.Type != "m.room.message" {
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
	if !m.isUserIDValidated(senderID) {
		m.log.Warn(fmt.Sprintf("onHandleReceiveMessage receive msg %s,but userId is never validate", utils.StringInterface(event, 3)))
		//return
	}
	peerAddress := m.userIDToAddress(senderID)
	peer := m.Peers[peerAddress]
	if peer == nil {
		m.temporaryAddress2Room[peerAddress] = event.RoomID
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
	m.log.Trace(fmt.Sprintf("onHandleMemberShipChange %s ", utils.StringInterface(event, 10)))
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
		//cannot verify this user
		if !m.isUserIDValidated(userid) {
			m.log.Warn(fmt.Sprintf("receive invite,but i don't know this user %s", utils.StringInterface(event, 5)))
			return
		}
		go func() {
			//todo fixme why need sleep, otherwise join will faile because of forbidden
			time.Sleep(time.Second)
			//one must join to be able to get room alias
			_, err := m.matrixcli.JoinRoom(event.RoomID, "", nil)
			if err != nil {
				m.log.Error(fmt.Sprintf("JoinRoom %s ,err %s", event.RoomID, err))
				return
			}
			peerAddress := m.userIDToAddress(userid)
			peer := m.Peers[peerAddress]
			if peer == nil {
				//maybe a peer want send secret request to me
				m.temporaryAddress2Room[peerAddress] = event.RoomID
			} else {
				err = m.inviteIfPossible(userid, event.RoomID)
				if err != nil {
					m.log.Error(fmt.Sprintf("inviteIfPossible %s to default room err %s", userid, err))
				}
			}
		}()
	} else if membership == "join" {
		/*
			{
							"content": {
								"membership": "join",
								"avatar_url": null,
								"displayname": "214e-118e7ade8cf61531f0d1629febf299ab939f314f04391b59b3567314525b4bee77220115d9bb0a808b6036857dad8bf242cbf5c09b18e7778c95022a77058bdf1c"
							},
							"type": "m.room.member",
							"sender": "@0x214e7247a2757696ed2986a8331a9e27a330c750:transport01.smartmesh.cn",
							"state_key": "@0x214e7247a2757696ed2986a8331a9e27a330c750:transport01.smartmesh.cn"
						}
		*/
		displayname, ok := event.ViewContent("displayname")
		if !ok {
			m.log.Warn(fmt.Sprintf("receive join,but has no display name %s", utils.StringInterface(event, 5)))
			return
		}
		avatar, _ := event.ViewContent("avatar_url")
		user := &gomatrix.UserInfo{
			DisplayName: displayname,
			UserID:      event.Sender,
			AvatarURL:   avatar,
		}
		//user can not be verified
		if err := m.validateAndUpdateUser(user); err != nil {
			m.log.Warn(fmt.Sprintf("receive join ,but user cannot be verified %s err %s", utils.StringInterface(event, 5), err))
			return
		}
		err := m.inviteIfPossible(userid, event.RoomID)
		if err != nil {
			m.log.Error(fmt.Sprintf("inviteIfPossible %s to default room err %s", userid, err))
		}
	} else {
		//todo fix me handle leave event
	}

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
	m.log.Trace(fmt.Sprintf("onHandlePresenceChange %s", utils.StringInterface(event, 5)))
	m.log.Trace(fmt.Sprintf("address i want to know: %s", utils.StringInterface(m.Peers, 3)))
	if m.stopreceiving == true {
		return
	}
	// message sender
	userid := event.Sender
	if event.Type != "m.presence" {
		m.log.Error(fmt.Sprintf("onHandlePresenceChange receive unkonw event %s", utils.StringInterface(event, 5)))
		return
	}
	//my self status change
	if userid == m.UserID {
		return
	}
	address := m.userIDToAddress(userid)
	peer, ok := m.Peers[address]
	if !ok {
		m.log.Trace(fmt.Sprintf("receive presence,but peer is unkown %s", utils.StringInterface(event, 5)))
		return
	}
	if !m.isUserIDValidated(userid) {
		m.log.Info(fmt.Sprintf("receive presence %s", utils.StringInterface(event, 5)))
		return
	}
	// parse address of message sender
	presence, exists := event.ViewContent("presence") //newest network status
	if !exists {
		return
	}
	if peer.isValidUserID(userid) && peer.setStatus(userid, presence) {
		//device type
		deviceType, _ := event.ViewContent("status_msg") //newest network status
		peer.deviceType = deviceType
	}
	m.log.Trace(fmt.Sprintf("peer %s status=%s,deviceType=%s", utils.APex2(address), peer.status, peer.deviceType))
}

// loginOrRegister node login, if failed, register again then try login,
// displayname of nodes as the signature of userID
func (m *MatrixTransport) loginOrRegister() (err error) {
	//TODO:Consider the risk of being registered maliciously
	loginok := false
	baseAddress := crypto.PubkeyToAddress(m.key.PublicKey)
	baseUsername := strings.ToLower(baseAddress.String())

	username := baseUsername
	password := hexutil.Encode(m.dataSign([]byte(m.servername)))
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
				authDict := &gomatrix.AuthDict{
					Type: AUTHTYPE,
				}
				req := &gomatrix.ReqRegister{
					DeviceID: "",
					Auth:     *authDict,
					Username: username,
					Password: password,
					Type:     LOGINTYPE,
				}
				_, uia, rerr := m.matrixcli.Register(req)
				if rerr != nil && uia == nil {
					rhttpErr, _ := err.(gomatrix.HTTPError)
					if rhttpErr.Code == 400 { //M_USER_IN_USE,M_INVALID_USERNAME,M_EXCLUSIVE
						m.log.Trace("username is in use or invalid,try again")
						continue
					}
				}
				m.matrixcli.UserID = username
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
	//set displayname as publicly visible
	dispname := m.getUserDisplayName(m.matrixcli.UserID)
	if err = m.matrixcli.SetDisplayName(dispname); err != nil {
		err = fmt.Errorf("could set the node's displayname and quit as well")
		m.matrixcli.ClearCredentials()
		return
	}
	m.log.Trace(fmt.Sprintf("userdisplayname=%s", dispname))
	return err
}

// makeRoomAlias name room's alias
func (m *MatrixTransport) makeRoomAlias(thepart string) string {
	return ROOMPREFIX + ROOMSEP + NETWORKNAME + ROOMSEP + thepart
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

// joinDiscoveryRoom : check discoveryroom if not exist, then create a new one.
// client caches all memebers of this room, and invite nodes checked from this room again.
func (m *MatrixTransport) joinDiscoveryRoom() (err error) {
	//read discovery room'name and fragment from "params-settings"
	// combine discovery room's alias
	discoveryRoomAlias := m.makeRoomAlias(ALIASFRAGMENT)
	discoveryRoomAliasFull := "#" + discoveryRoomAlias + ":" + DISCOVERYROOMSERVER
	m.discoveryroomid = ""
	// this node join the discovery room, if not exist, then create.
	for i := 0; i < 5; i++ {
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
				continue
			}
			m.discoveryroomid = respCreateRoom.RoomID
			continue
		} else {
			m.discoveryroomid = respJoinRoom.RoomID
			break
		}
	}
	//exit if join room failed
	if m.discoveryroomid == "" {
		errinfo := fmt.Sprintf("Discovery room {%s} not found and can't be created on a federated homeserver {%s}", discoveryRoomAliasFull, m.servername)
		err = fmt.Errorf(errinfo)
		m.log.Error(errinfo)
		return
	}
	return nil
}

func (m *MatrixTransport) userIDToAddress(userID string) common.Address {
	//check grammar of user ID
	_match := ValidUserIDRegex.MatchString(userID)
	if _match == false {
		m.log.Warn(fmt.Sprintf("UserID %s, format error", userID))
		return utils.EmptyAddress
	}
	addressHex, err := extractUserLocalpart(userID) //"@myname:smartraiden.org:cy"->"myname"
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
func (m *MatrixTransport) findOrCreateRoomByAddress(address common.Address, isPublic bool) (roomID string, users []*gomatrix.UserInfo, err error) {
	if m.stopreceiving {
		err = errors.New("already topped")
		return
	}
	//The following is the case where peer-to-peer communication room does not exist.
	var addressOfPairs = ""
	strPairs := []string{hexutil.Encode(m.NodeAddress.Bytes()), hexutil.Encode(address.Bytes())}
	sort.Strings(strPairs)
	addressOfPairs = strings.Join(strPairs, "_") //format "0cccc_0xdddd"
	tmpRoomName := m.makeRoomAlias(addressOfPairs)

	//try to get user-infos of communication with "account_data" from homeserver include the other participating servers.
	users, err = m.searchNode(address)
	if err != nil {
		return
	}

	//Join a room that connot be found by search_room_directory
	roomID, err = m.getUnlistedRoom(tmpRoomName, users, isPublic)
	m.log.Info(fmt.Sprintf("CHANNEL ROOM,peer_address=%s room=%s", address.String(), roomID))
	return
}

// getUnlistedRoom get a conversation room that cannnot be found by search_room_directory.
// If the room is not exist and create a unnamed room for communication,invite the node finally.
// This process of join-create-join-room may be repeated 3 times(network delay)
func (m *MatrixTransport) getUnlistedRoom(roomname string, users []*gomatrix.UserInfo, isPublic bool) (roomID string, err error) {
	roomNameFull := "#" + roomname + ":" + m.servername
	var inviteesUids []string
	for _, user := range users {
		inviteesUids = append(inviteesUids, user.UserID)
	}
	req := &gomatrix.ReqCreateRoom{
		Invite:     inviteesUids,
		Visibility: "private",
		Preset:     "trusted_private_chat",
	}
	if isPublic {
		req.Visibility = "public"
		req.Preset = "public_chat"
	} else {
		req.Visibility = "private"
		req.Preset = "trusted_private_chat"
	}
	unlistedRoomid := ""
	for i := 0; i < 6; i++ {
		var respJoinRoom *gomatrix.RespJoinRoom
		respJoinRoom, err = m.matrixcli.JoinRoom(roomNameFull, m.servername, nil)
		if err != nil {
			req.RoomAliasName = roomname
			_, err = m.matrixcli.CreateRoom(req)
			if err != nil {
				m.log.Info(fmt.Sprintf("Room %s not found,trying to create it. but fail %s", roomname, err))
			}
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
func (m *MatrixTransport) getPeerRoomID(address common.Address) string {
	u := m.Peers[address]
	if u != nil {
		return u.defaultMessageRoomID
	}
	return ""
}

func (m *MatrixTransport) searchNode(address common.Address) (users []*gomatrix.UserInfo, err error) {
	respusers, err := m.matrixcli.SearchUserDirectory(&gomatrix.ReqUserSearch{
		SearchTerm: strings.ToLower(address.String()),
		Limit:      10,
	})
	if err != nil {
		return
	}
	if len(respusers.Results) == 0 {
		return
	}
	for _, user := range respusers.Results {
		var xaddr common.Address
		xaddr, err = validateUseridSignature(&user)
		//validate failed
		if err != nil {
			m.log.Error(fmt.Sprintf("validateUseridSignature for %s err %s", user.UserID, err))
			continue
		}
		//validate failed
		if xaddr != address {
			continue
		}
		//save this validated user, should copy of  `user`, otherwise `user` will be changed later
		user2 := user
		m.validatedUsers[user.UserID] = &user2
		users = append(users, &user2)
	}
	return
}
func (m *MatrixTransport) handleNewPartner(p *MatrixPeer) (err error) {
	roomID, users, err := m.findOrCreateRoomByAddress(p.address, false)
	if err != nil {
		return
	}
	for _, u := range users {
		err = m.validateAndUpdateUser(u)
		if err != nil {
			return
		}
	}
	p.defaultMessageRoomID = roomID
	address2Room := make(map[common.Address]string)
	for addr, peer := range m.Peers {
		address2Room[addr] = peer.defaultMessageRoomID
	}
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
		}
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
			if err = m.validateAndUpdateUser(u); err == nil {
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
				if p.setStatus(u.UserID, presenceResponse.Presence) {
					p.deviceType = presenceResponse.Presence
				}
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

	users, err := m.searchNode(p.address)
	if err != nil {
		fmt.Fprintf(errBuf, "findValidUsersOnServer for %s,err=%s\n", utils.APex2(p.address), err)
	} else {
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
	peer := m.Peers[peerAddress]
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
	if peer.defaultMessageRoomID != eventRoom && !peer.isValidUserID(userID) {
		//invite this user to the default room
		_, err := m.matrixcli.InviteUser(peer.defaultMessageRoomID, &gomatrix.ReqInviteUser{
			UserID: userID,
		})
		//if has already startup ,we should known all user status
		if err != nil && m.hasDoneStartCheck {
			return fmt.Errorf("InviteUser perr=%s,room=%s err=%s",
				utils.APex2(peerAddress), peer.defaultMessageRoomID, err)
		}
	}
	if peer.defaultMessageRoomID == eventRoom {
		//this user is joinning in the default message room, a new user
		if !peer.isValidUserID(userID) {
			//update presence if possible
			peer.updateUsers([]*gomatrix.UserInfo{m.validatedUsers[userID]})
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

// ExtractUserLocalpart Extract user name from user ID
func extractUserLocalpart(userID string) (string, error) {
	if len(userID) == 0 || userID[0] != '@' {
		return "", fmt.Errorf("%s is not a valid user id", userID)
	}
	return strings.SplitN(userID[1:], ":", 2)[0], nil
}

// validate_userid_signature
func validateUseridSignature(user *gomatrix.UserInfo) (address common.Address, err error) {
	//displayname should be an address in the self._userid_re format
	_match := ValidUserIDRegex.MatchString(user.UserID)
	if _match == false {
		err = fmt.Errorf("validate user info failed")
		return
	}
	_address, err := extractUserLocalpart(user.UserID) //"@myname:smartraiden.org:cy"->"myname"
	if err != nil {
		return
	}
	var addrmuti = regexp.MustCompile(`^(0x[0-9a-f]{40})`)
	addrlocal := addrmuti.FindString(_address)
	if addrlocal == "" {
		err = fmt.Errorf("%s not match our userid rule", user.UserID)
		return
	}
	addressBytes, err := hexutil.Decode(addrlocal)
	if err != nil {
		return
	}
	signature, err := getSignatureFromDisplayName(user.DisplayName)
	if err != nil {
		return
	}
	recovered, err := utils.Ecrecover(utils.Sha3([]byte(user.UserID)), signature)
	if err != nil {
		return
	}
	if !bytes.Equal(recovered[:], addressBytes) {
		err = fmt.Errorf("validate %s failed", user.UserID)
		return
	}
	address = common.BytesToAddress(addressBytes)
	return
}
