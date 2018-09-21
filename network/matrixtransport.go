package network

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/channel/channeltype"
	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/matrixcomm"
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
	ROOMPREFIX = "raiden"
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
)

// MatrixTransport represents a matrix transport Instantiation
type MatrixTransport struct {
	matrixcli          *matrixcomm.MatrixClient //the instantiated matrix
	servername         string                   //the homeserver's name
	running            bool                     //running status
	stopreceiving      bool                     //Whether to stop accepting(data)
	key                *ecdsa.PrivateKey        //key
	NodeAddress        common.Address
	protocol           ProtocolReceiver
	discoveryroomalias string                          //the room's alias of sys pre-configured ("#[RoomNameLocalpart]:[ServerName]")
	discoveryroomid    string                          //the room's ID of sys pre-configured ("![RoomIdData]:[ServerName]")
	Users              map[string]*matrixcomm.UserInfo //cache user's base-infos("userID{userID,displayname,avatarurl}")
	Address2User       map[common.Address][]*matrixcomm.UserInfo
	AddressToPresence  map[common.Address]*matrixcomm.RespPresenceUser //cache user's real-time presence by node's address("userID{presence}")
	Userid2Presence    map[string]*matrixcomm.RespPresenceUser         //cache user's real-time presence by userID(one address maybe have many userIDs)
	UserID             string                                          //the current user's ID(@kitty:thisserver)
	NodeDeviceType     string
	UserDeviceType     map[common.Address]string
	avatarurl          string
	Address2Room       map[string]string //all rooms with we knows,just
	log                log.Logger
	ChargeRegulation   string
	statusChan         chan netshare.Status
	status             netshare.Status
}

var (
	// ValidUserIDRegex user ID 's format
	ValidUserIDRegex = regexp.MustCompile(`^@(0x[0-9a-f]{40})(?:\.[0-9a-f]{8})?(?::.+)?$`) //(`^[0-9a-z_\-./]+$`)
	//NETWORKNAME which network is used
	NETWORKNAME = params.NETWORKNAME
	//ALIASFRAGMENT the terminal part of alias
	ALIASFRAGMENT = ""
	//DISCOVERYROOMSERVER discovery room server name
	DISCOVERYROOMSERVER = ""
)

func (mtr *MatrixTransport) changeStatus(newStatus netshare.Status) {
	log.Info(fmt.Sprintf("changeStatus from %d to %d", mtr.status, newStatus))
	mtr.status = newStatus
	select {
	case mtr.statusChan <- newStatus:
	default:
	}
}

// CollectNeighbors subscribe status change
func (mtr *MatrixTransport) CollectNeighbors(db xmpptransport.XMPPDb) error {
	cs, err := db.GetChannelList(utils.EmptyAddress, utils.EmptyAddress)
	if err != nil {
		return err
	}
	for _, c := range cs {
		err := mtr.nodeHealthCheck(c.PartnerAddress())
		if err != nil {
		}
	}
	db.RegisterNewChannellCallback(func(c *channeltype.Serialization) (remove bool) {
		err := mtr.nodeHealthCheck(c.PartnerAddress())
		if err != nil {
			return false
		}
		return true
	})
	db.RegisterChannelStateCallback(func(c *channeltype.Serialization) (remove bool) {
		err := mtr.nodeHealthCheck(c.PartnerAddress())
		if err != nil {
			return false
		}
		return true
	})
	return nil
}

/*
------------------------------------------------------------------------------------------------------------------------
*/

// HandleMessage regist the interface of call receive(func)
func (mtr *MatrixTransport) HandleMessage(from common.Address, data []byte) {
	if !mtr.running || mtr.stopreceiving {
		return
	}
	if mtr.protocol != nil {
		mtr.protocol.receive(data)
	}
}

// RegisterProtocol regist the interface of call RegisterProtocol(func)
func (mtr *MatrixTransport) RegisterProtocol(protcol ProtocolReceiver) {
	mtr.protocol = protcol
}

// Stop Does Stop need to destroy matrix resource ?
func (mtr *MatrixTransport) Stop() {
	if mtr.running == false {
		return
	}
	mtr.running = false
	mtr.changeStatus(netshare.Closed)
	go func() {
		mtr.matrixcli.SetPresenceState(&matrixcomm.ReqPresenceUser{
			Presence: OFFLINE,
		})
	}()
	mtr.matrixcli.StopSync()
	if _, err := mtr.matrixcli.Logout(); err != nil {
		log.Error("[Matrix] Logout failed")
	}
}

// StopAccepting stop receive message and wait
func (mtr *MatrixTransport) StopAccepting() {
	mtr.stopreceiving = true
}

// NodeStatus gets Node states of network, if check self node, `isOnline` is not always be true instead it switches according to server handshake signal.
func (mtr *MatrixTransport) NodeStatus(addr common.Address) (deviceType string, isOnline bool) {
	if mtr.matrixcli == nil {
		return "", false
	}
	_, isexist := mtr.AddressToPresence[addr]
	if !isexist {
		isOnline = false
		return "", isOnline
	}

	if mtr.AddressToPresence[addr].Presence != ONLINE {
		isOnline = false
		return "", isOnline
	}
	isOnline = true //unless node is online cann't get status_msg
	deviceType = mtr.AddressToPresence[addr].StatusMsg

	/*deviceType = mtr.NodeDeviceType //just test
	isOnline = true*/
	// we can check user presence on any partner servers when invite node is in the presence list.
	// code above can't get devideType, which handles via /presence/list (if nodes online, it returns one more status_msg(with devideType)
	// via {userid}/presence we can check and expand states of multiple nodes.
	return
}

// Send send message
func (mtr *MatrixTransport) Send(receiverAddr common.Address, data []byte) error {
	if !mtr.running || len(data) == 0 {
		return fmt.Errorf("[Matrix]Send failed,matrix not running or send data is null")
	}
	room, err := mtr.getRoom2Address(receiverAddr)
	if err != nil || room == nil {
		return fmt.Errorf("[Matrix]Send failed,cann't find the peer address")
	}
	_data := base64.StdEncoding.EncodeToString(data)
	_, err = mtr.matrixcli.SendText(room.ID, _data)
	if err != nil {
		log.Error(fmt.Sprintf("[matrix]send failed to %s, message=%s", utils.APex2(receiverAddr), encoding.MessageType(data[0])))
		return err
	}
	log.Trace(fmt.Sprintf("[Matrix]Send to %s, message=%s", utils.APex2(receiverAddr), encoding.MessageType(data[0])))
	return nil
}

// Start matrix
func (mtr *MatrixTransport) Start() {
	if mtr.running {
		return
	}

	// log in
	if err := mtr.loginOrRegister(); err != nil {
		return
	}
	mtr.running = true
	mtr.stopreceiving = false
	mtr.changeStatus(netshare.Connected)

	// health-check, used to find history rooms this node ever joined.
	err := mtr.nodeHealthCheck(mtr.NodeAddress)
	if err != nil {
		return
	}

	//initialize Filters/NextBatch/Rooms
	store := matrixcomm.NewInMemoryStore()
	mtr.matrixcli.Store = store

	//handle the issue of discoveryroom,FOR TEST,temporarily retain this room
	if err := mtr.joinDiscoveryRoom(); err != nil {
		return
	}
	//search store->room，isn't it in listening room
	if err := mtr.inventoryRooms(); err != nil {
		return
	}
	//notify to server i am online（include the other participating servers）
	if err := mtr.matrixcli.SetPresenceState(&matrixcomm.ReqPresenceUser{
		Presence:  ONLINE,
		StatusMsg: mtr.NodeDeviceType, //register device type to server
	}); err != nil {
		return
	}
	//register receive-datahandle or other message received
	mtr.matrixcli.Store = store
	mtr.matrixcli.Syncer = matrixcomm.NewDefaultSyncer(mtr.UserID, store)
	syncer := mtr.matrixcli.Syncer.(*matrixcomm.DefaultSyncer)

	syncer.OnEventType("network.raiden.rooms", mtr.onHandleAccountData)

	syncer.OnEventType("m.room.message", mtr.onHandleReceiveMessage)

	syncer.OnEventType("m.presence", mtr.onHandlePresenceChange)

	syncer.OnEventType("m.room.member", mtr.onHandleMemberShipChange)

	go func() {
		for {
			if err := mtr.matrixcli.Sync(); err != nil {
				log.Error(fmt.Sprintf("Matrix Sync return,err=%s ,will try agin..", err))
				mtr.changeStatus(netshare.Reconnecting)
			}
			time.Sleep(time.Second * 5)
			mtr.changeStatus(netshare.Connected)
		}
	}()

	log.Trace("[Matrix] transport started")
	/*//test code
	go func() {
		for {
			sdata:="testhellohellohellohellotesthellohellohellohello"
			xbyte,err:=base64.StdEncoding.DecodeString(sdata)
			if err!=nil{
				//fmt.Println("ERROR XXX")
				panic("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
			}
			err0:=mtr.Send(common.HexToAddress("0xc67f23ce04ca5e8dd9f2e1b5ed4fad877f79267a"), xbyte)
			if err0!=nil{}
			time.Sleep(time.Second* 3)
		}
	}()*/
}

/*
------------------------------------------------------------------------------------------------------------------------
*/
//onHandleReceiveMessage push the message of some one send "account_data"
func (mtr *MatrixTransport) onHandleAccountData(event *matrixcomm.Event) {
	if mtr.stopreceiving || event.Type != "network.raiden.rooms" {
		return
	}
	userid := event.Sender
	if userid == mtr.UserID {
		return
	}
	value, exist := event.ViewContent("account_data")
	if exist && value != "" {

	}
}

// onHandleReceiveMessage handle text messages sent to listening rooms
func (mtr *MatrixTransport) onHandleReceiveMessage(event *matrixcomm.Event) {
	if mtr.stopreceiving || event.Type != "m.room.message" {
		return
	}
	msgtype, ok := event.MessageType()
	if ok == false || msgtype != "m.text" {
		return
	}

	senderID := event.Sender
	roomID := event.RoomID
	if senderID == mtr.UserID {
		return
	}

	tmpuser := &matrixcomm.UserInfo{
		UserID: senderID,
	}
	user, err := mtr.verifyAndUpdateUserCache(tmpuser)
	if err != nil {
		return
	}
	peerAddress, err := validateUseridSignature(*user)
	if err != nil {
		log.Warn(fmt.Sprintf("Receive message from a user without legala displayName signature,peer_user=%s,room=%s", user.UserID, roomID))
		return
	}

	oldRoomID := mtr.getRoomID2Address(peerAddress)
	if oldRoomID != roomID {
		log.Warn(fmt.Sprintf("Receive message triggered new room for peer,"+
			" peer_user=%s,peer_address=%s,old_room=%s,room=%s", user.UserID, peerAddress.String(), oldRoomID, roomID))
	}
	err = mtr.setRoomID2Address(peerAddress, roomID)

	if _, ok = mtr.Address2User[peerAddress]; !ok {
		//return
	}

	data, ok := event.Body()
	if !ok || len(data) < 2 {
		//return
	}
	//message :=[]byte{}
	/*if data[0:len(data)-2] == "0x" {
		_, err = hexutil.Decode(data)
		if err != nil {
			log.Warn(fmt.Sprintf("Receive message binary data is not a valid message,message_data=%s,peer_address=%s", data, peerAddress.String()))
			return
		}

	} else {*/
	/*//解析json数据
	message,err=hexutil.Decode(data)
	if err!=nil{}
	if !json.Valid(message){
		log.Warn(fmt.Sprintf("Message data JSON are not a valid message,message_data=%s,peer_address=%s",data,peerAddress.String()))
		return
	}
	//ping
	bytesHead:=encoding.MessageType(message[0])
	if bytesHead==1{
		if senderID!=hexutil.Encode(peerAddress.Bytes()){
			log.Warn(fmt.Sprintf("Not required Ping received,message=%s",data))
			return
		}
	}
	//or signedmessage
	if senderID!=hexutil.Encode(peerAddress.Bytes()){
		log.Warn(fmt.Sprintf("Receive message from a user without legala displayName signature,peer_user=%s,room=%s",user.UserID,roomID))
		return
	}*/

	msgSender, err := extractUserLocalpart(senderID)
	if err != nil {
		return
	}
	dataContent, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		log.Error(fmt.Sprintf("[Matrix]Receive unkown message %s", utils.StringInterface(event, 5)))
	} else {
		mtr.HandleMessage(common.HexToAddress(msgSender), dataContent)
		log.Info(fmt.Sprintf("[Matrix]Receive message %s from %s", encoding.MessageType(dataContent[0]), utils.APex2(common.HexToAddress(msgSender))))

	}
}

// onHandleMemberShipChange Handle message when eventType==m.room.member and join all invited rooms
func (mtr *MatrixTransport) onHandleMemberShipChange(event *matrixcomm.Event) {
	if mtr.stopreceiving || event.Type != "m.room.member" {
		return
	}
	value, exists := event.ViewContent("membership")
	if !exists || value == "" || value != "invite" {
		return
	}
	userid := event.Sender
	tmpuser := &matrixcomm.UserInfo{
		UserID: userid,
	}
	user, err := mtr.verifyAndUpdateUserCache(tmpuser)
	peerAddress, err := validateUseridSignature(*user)
	if err != nil {
		log.Warn("Got invited to a room by invalid signed user - ignoring")
		return
	}
	//one must join to be able to get room alias
	_, err = mtr.matrixcli.JoinRoom(event.RoomID, "", nil)
	if err != nil {
		return
	}
	if mtr.matrixcli.Store.LoadRoom(event.RoomID) == nil {
		theroom := &matrixcomm.Room{
			ID: event.RoomID,
		}
		mtr.matrixcli.Store.SaveRoom(theroom)
	}
	//cache RooID2ADDRESS and notify to servers
	err = mtr.setRoomID2Address(peerAddress, event.RoomID)
	if err != nil {
		return
	}
}

// onHandlePresenceChange handle events in this message, about changes of nodes and update AddressToPresence
func (mtr *MatrixTransport) onHandlePresenceChange(event *matrixcomm.Event) {
	if mtr.stopreceiving == true {
		return
	}
	// message sender
	userid := event.Sender
	if event.Type != "m.presence" || userid == mtr.UserID {
		return
	}
	var userDisplayname = ""
	// read sender's displayname from user cache.
	tmpuser := &matrixcomm.UserInfo{
		UserID: userid,
	}
	user, err := mtr.verifyAndUpdateUserCache(tmpuser)
	if err == nil {
		userDisplayname = user.DisplayName
	}
	// get sender's displayname from message
	value, exists := event.ViewContent("displayname") //no displayname return form event sometimes
	if exists && value != "" {
		userDisplayname = value
	}

	// both sources above are ok
	if userDisplayname == "" {
		return
	}
	user.DisplayName = userDisplayname

	// parse address of message sender
	peerAdderss, err := validateUseridSignature(*user)
	if err != nil {
		return
	}
	//not a user we've started healthcheck, skip--??
	if _, ok := mtr.Address2User[peerAdderss]; !ok {
		return
	}
	mtr.Address2User[peerAdderss] = append(mtr.Address2User[peerAdderss], user)

	//maybe inviting user used to also possibly invite user's from discovery presence changes
	mtr.maybeInviteUser(*user)

	vValue, exists := event.ViewContent("presence") //newest network status
	if !exists {
		return
	}
	newstate := vValue
	//presence status unchanged
	if _, ok := mtr.Userid2Presence[userid]; !ok {
		return
	}
	oldstate := mtr.Userid2Presence[userid].Presence
	if newstate == oldstate {
		return
	}
	//device type
	dValue, exists := event.ViewContent("status_msg") //newest network status
	if !exists {

	} else {
		nodeDeviceType := dValue
		mtr.UserDeviceType[peerAdderss] = nodeDeviceType

	}
	//presence status hava changed
	mtr.Userid2Presence[userid].Presence = newstate
	mtr.Userid2Presence[userid].StatusMsg = dValue
	mtr.updateAddressPresence(peerAdderss, dValue)
}

// getUserPresence get the presence state from Userid2Presence
func (mtr *MatrixTransport) getUserPresence(userid string) (presence *matrixcomm.RespPresenceUser, err error) {
	//if this user does not exist on cache of UseridToPresence，then request the server temporarily
	pUser := &matrixcomm.RespPresenceUser{
		UserID:    userid,
		Presence:  "",
		StatusMsg: "",
	}
	if _, ok := mtr.Userid2Presence[userid]; !ok {
		resp, err := mtr.matrixcli.GetPresenceState(userid)
		if err != nil {
			pUser.Presence = UNKNOWN
			mtr.Userid2Presence[userid] = pUser
		} else {
			pUser = resp
			//update userID's presence->UseridToPresence
			mtr.Userid2Presence[userid] = pUser
		}
	}

	pUser = mtr.Userid2Presence[userid]
	presence = pUser
	return
}

// updateAddressPresence Update synthesized address presence state from user presence state
func (mtr *MatrixTransport) updateAddressPresence(address common.Address, msgstatus string) {
	// an address can match multiple userid / presence
	compositepresence := []string{}
	var tmpUserInfos []*matrixcomm.UserInfo
	if _, ok := mtr.Address2User[address]; !ok {
		return
	}
	tmpUserInfos = mtr.Address2User[address]
	for _, v := range tmpUserInfos {
		resp, err := mtr.getUserPresence(v.UserID)
		if err != nil {
			continue
		}
		compositepresence = append(compositepresence, resp.Presence)
	}

	// check presence state by the order of online, unavailable, offlien, unknown
	// if it's not just a userid,calculate netstaus according to this rule
	presencestates := []string{ONLINE, UNAVAILABLE, OFFLINE, UNKNOWN}
	newState := UNKNOWN
	for _, xstate := range presencestates {
		for _, xpresence := range compositepresence {
			if xpresence == xstate {
				newState = xpresence
				break
			}
		}
	}
	//update AddressToPresence
	if _, ok := mtr.AddressToPresence[address]; ok {
		if mtr.AddressToPresence[address].Presence == newState {
			return
		}
	}
	tmpuserp := &matrixcomm.RespPresenceUser{
		Presence:  newState,
		UserID:    tmpUserInfos[0].UserID,
		StatusMsg: msgstatus,
	}
	mtr.AddressToPresence[address] = tmpuserp
}

// loginOrRegister node login, if failed, register again then try login,
// displayname of nodes as the signature of userID
func (mtr *MatrixTransport) loginOrRegister() (err error) {
	//TODO:Consider the risk of being registered maliciously
	regok := false
	loginok := false
	baseAddress := crypto.PubkeyToAddress(mtr.key.PublicKey)
	baseUsername := strings.ToLower(baseAddress.String())

	username := baseUsername
	password := hexutil.Encode(mtr.dataSign([]byte(mtr.servername)))
	//password := "12345678"
	for i := 0; i < 5; i++ {
		var resplogin *matrixcomm.RespLogin
		if regok == false {
			//rand.Seed(time.Now().UnixNano())
			//rnd := Int32ToBytes(rand.Int31n(math.MaxInt32))
			//username = baseUsername + "." + hex.EncodeToString(rnd)
		}
		mtr.matrixcli.AccessToken = ""
		resplogin, err = mtr.matrixcli.Login(&matrixcomm.ReqLogin{
			Type:     LOGINTYPE,
			User:     username,
			Password: password,
			DeviceID: "",
		})
		if err != nil {
			httpErr, ok := err.(matrixcomm.HTTPError)
			if !ok { // network error,try again
				continue
			}
			if httpErr.Code == 403 { //Invalid username or password
				if i > 0 {
					log.Trace(fmt.Sprintf("couldn't sign in for matrix,trying register %s", username))
				}
				authDict := &matrixcomm.AuthDict{
					Type: AUTHTYPE,
				}
				req := &matrixcomm.ReqRegister{
					DeviceID: "",
					Auth:     *authDict,
					Username: username,
					Password: password,
					Type:     LOGINTYPE,
				}
				_, uia, rerr := mtr.matrixcli.Register(req)
				if rerr != nil && uia == nil {
					rhttpErr, _ := err.(matrixcomm.HTTPError)
					if rhttpErr.Code == 400 { //M_USER_IN_USE,M_INVALID_USERNAME,M_EXCLUSIVE
						log.Trace("username is in use or invalid,try again")
						continue
					}
				}
				regok = true
				mtr.matrixcli.UserID = username
				continue
			}
		} else {
			//cache the node's and report the UserID and AccessToken to matrix
			mtr.matrixcli.SetCredentials(resplogin.UserID, resplogin.AccessToken)
			mtr.UserID = resplogin.UserID
			mtr.NodeAddress = baseAddress
			loginok = true
			break
		}
	}
	if !loginok {
		err = fmt.Errorf("could not register or login")
		return
	}
	//set displayname as publicly visible
	dispname := hexutil.Encode(mtr.dataSign([]byte(mtr.matrixcli.UserID)))
	if err = mtr.matrixcli.SetDisplayName(dispname); err != nil {
		err = fmt.Errorf("could set the node's displayname and quit as well")
		mtr.matrixcli.ClearCredentials()
		return
	}
	// Add nodes info into Users
	thisUser := &matrixcomm.UserInfo{
		UserID:      mtr.UserID,
		DisplayName: dispname,
		AvatarURL:   mtr.avatarurl,
	}
	_, err = mtr.verifyAndUpdateUserCache(thisUser)

	return err
}

// inventoryRooms : collect monitored room, discovery room are not put inside listening object.
func (mtr *MatrixTransport) inventoryRooms() (err error) {
	for _, value := range mtr.matrixcli.Store.LoadRoomOfAll() {
		if value.Alias == mtr.discoveryroomalias {
			continue
		}
		theroom := &matrixcomm.Room{
			ID:    mtr.discoveryroomid,
			Alias: mtr.discoveryroomalias,
		}
		mtr.matrixcli.Store.SaveRoom(theroom)
	}
	return nil
}

// makeRoomAlias name room's alias
func (mtr *MatrixTransport) makeRoomAlias(thepart string) string {
	return ROOMPREFIX + ROOMSEP + NETWORKNAME + ROOMSEP + thepart
}

// dataSign signature data
func (mtr *MatrixTransport) dataSign(data []byte) (signature []byte) {
	hash := crypto.Keccak256(data)
	signature, err := crypto.Sign(hash[:], mtr.key)
	if err != nil {
		return nil
	}
	return
}

// joinDiscoveryRoom : check discoveryroom if not exist, then create a new one.
// client caches all memebers of this room, and invite nodes checked from this room again.
func (mtr *MatrixTransport) joinDiscoveryRoom() (err error) {
	//read discovery room'name and fragment from "params-settings"
	discoveryRoomList := params.MatrixDiscoveryRoomConfig
	for _, value := range discoveryRoomList {
		itemname := value[0]
		itemvalue := value[1]
		if itemname == "aliassegment" {
			ALIASFRAGMENT = itemvalue
		}
		if itemname == "server" {
			DISCOVERYROOMSERVER = itemvalue
		}
	}
	// combine discovery room's alias
	discoveryRoomAlias := mtr.makeRoomAlias(ALIASFRAGMENT)
	discoveryRoomAliasFull := "#" + discoveryRoomAlias + ":" + DISCOVERYROOMSERVER
	mtr.discoveryroomid = ""
	// this node join the discovery room, if not exist, then create.
	for i := 0; i < 5; i++ {
		respj, errj := mtr.matrixcli.JoinRoom(discoveryRoomAliasFull, mtr.servername, nil)
		if errj != nil {
			//if Room doesn't exist and then create the room(this is the node's resposibility)
			if mtr.servername != DISCOVERYROOMSERVER {
				break
			}
			//try to create the discovery room
			var _visibility = "private"
			if CHATPRESET == "public_chat" {
				_visibility = "public"
			}
			respc, errc := mtr.matrixcli.CreateRoom(&matrixcomm.ReqCreateRoom{
				RoomAliasName: discoveryRoomAlias,
				Preset:        CHATPRESET,
				Visibility:    _visibility,
			})
			if errc != nil {
				log.Error("can't create a discovery room,try again")
				continue
			}
			mtr.discoveryroomid = respc.RoomID
			continue
		} else {
			mtr.discoveryroomid = respj.RoomID
			break
		}
	}
	//exit if join room failed
	if mtr.discoveryroomid == "" {
		errinfo := fmt.Sprintf("Discovery room {%s} not found and can't be created on a federated homeserver {%s}", discoveryRoomAliasFull, mtr.servername)
		err = fmt.Errorf(errinfo)
		log.Error(errinfo)
		return
	}
	/*//try to commit the discovery room to memory
	mtr.discoveryroomalias = discoveryRoomAlias

	//add discovery room to listening object
	theroom := &matrixcomm.Room{
		ID:    mtr.discoveryroomid,
		Alias: discoveryRoomAlias,
		//State:nil,
	}
	mtr.matrixcli.Store.SaveRoom(theroom)

	//把discovery room放入RoomID2Address
	userAddr := mtr.NodeAddress
	err = mtr.setRoomID2Address(userAddr, mtr.discoveryroomid)*/

	//get the members which were joined the discovery room
	respin, err := mtr.matrixcli.JoinedMembers(mtr.discoveryroomid)
	if err != nil {
		log.Error("The node can't join room ", mtr.discoveryroomalias)
		return
	}
	for userid, userdata := range respin.Joined {
		//invite users
		usr := matrixcomm.UserInfo{
			UserID:      userid,
			DisplayName: *userdata.DisplayName,
		}
		//cache known users to Users
		_, xerr := mtr.verifyAndUpdateUserCache(&usr)
		if xerr != nil {
		}
		//invite them to the discovery room
		mtr.maybeInviteUser(usr)
	}

	return
}

// maybeInviteUser invite nodes to their rooms via Address2Room(search by "Address2Room").
func (mtr *MatrixTransport) maybeInviteUser(user matrixcomm.UserInfo) {
	address, err := validateUseridSignature(user)
	if err != nil {
		return
	}
	roomid := mtr.getRoomID2Address(address)
	if roomid == "" {
		return
	}
	room := mtr.matrixcli.Store.LoadRoom(roomid)
	if room == nil {
		theroom := &matrixcomm.Room{
			ID: roomid,
		}
		mtr.matrixcli.Store.SaveRoom(theroom)
	}
	//room already found the invite the user
	resp, err := mtr.matrixcli.JoinedMembers(roomid)
	if err != nil {
		return
	}
	//invite the user when it not in Address2Room
	if _, exist := resp.Joined[user.UserID]; !exist {
		_, err = mtr.matrixcli.InviteUser(roomid, &matrixcomm.ReqInviteUser{
			UserID: user.UserID,
		})
	}
	return
}

// verifyAndUpdateUserCache Verify user and standardized user to user-info,cache as "Users"
func (mtr *MatrixTransport) verifyAndUpdateUserCache(user0 *matrixcomm.UserInfo) (user1 *matrixcomm.UserInfo, err error) {
	//check grammar of user ID
	_match := ValidUserIDRegex.MatchString(user0.UserID)
	if _match == false {
		user1 = nil
		err = fmt.Errorf("User ID is illegal")
		return
	}
	if _, ok := mtr.Users[user0.UserID]; !ok {
		mtr.Users[user0.UserID] = user0
	}
	user1 = mtr.Users[user0.UserID]
	err = nil

	return
}

// getRoom2Address Get the room(*object) info by the node address.
// If no communication(peer-to-peer) room found yet,then "SearchUserDirectory" and create a temporary(unnamed) room for communication,invite the node finally.
func (mtr *MatrixTransport) getRoom2Address(address common.Address) (room *matrixcomm.Room, err error) {
	if mtr.stopreceiving {
		return
	}

	addressHex := hexutil.Encode(address.Bytes())

	//Well,I know where the peer is.
	roomid := mtr.getRoomID2Address(address)
	if roomid != "" {
		room = mtr.matrixcli.Store.LoadRoom(roomid)
		return
	}

	//The following is the case where peer-to-peer communication room does not exist.
	var addressOfPairs = ""
	if mtr.NodeAddress == address {
		return
	}
	strPairs := []string{hexutil.Encode(mtr.NodeAddress.Bytes()), hexutil.Encode(address.Bytes())}
	sort.Strings(strPairs)
	addressOfPairs = strings.Join(strPairs, "_") //format "0cccc_0xdddd"
	tmpRoomName := mtr.makeRoomAlias(addressOfPairs)

	//try to get user-infos of communication with "account_data" from homeserver include the other participating servers.
	var tmpUserInfos []*matrixcomm.UserInfo
	respusers, err := mtr.matrixcli.SearchUserDirectory(&matrixcomm.ReqUserSearch{
		SearchTerm: addressHex,
		//Limit:1024,
	})
	if err != nil {
		return
	}
	for _, resultx := range respusers.Results {
		xaddr, xerr := validateUseridSignature(resultx)
		if xerr != nil {
			continue
		}
		if xaddr != address {
			continue
		}
		_, cerr := mtr.verifyAndUpdateUserCache(&resultx)
		if cerr != nil {
		}
		tmpUserInfos = append(tmpUserInfos, &resultx)
	}
	//Shoot! I don't know where the node is
	if len(tmpUserInfos) == 0 {
		return
	}

	//update cache as Address2Usermap
	mtr.Address2User[address] = tmpUserInfos

	//Join a room that connot be found by search_room_directory
	room, err = mtr.getUnlistedRoom(tmpRoomName, tmpUserInfos)

	//update user account_data,also update cache as "RoomID2Address"
	err = mtr.setRoomID2Address(address, room.ID)

	//Make sure the users(one node more than one account) invited,the users may be on different servers.
	for _, xuser := range tmpUserInfos {
		mtr.maybeInviteUser(*xuser)
	}

	//Ensure that this room exists in my listening task
	if mtr.matrixcli.Store.LoadRoom(room.ID) == nil {
		mtr.matrixcli.Store.SaveRoom(room)
	}

	log.Info(fmt.Sprintf("CHANNEL ROOM,peer_address=%s room=%s", addressHex, room.ID))

	//fmt.Println(addressOfPairs)
	if _, ok := mtr.Address2User[address]; !ok {
		log.Info(fmt.Sprintf("Address not health checked:me=%s peer_address=%s", mtr.UserID, addressHex))
	}

	return
}

// getUnlistedRoom get a conversation room that cannnot be found by search_room_directory.
// If the room is not exist and create a unnamed room for communication,invite the node finally.
// This process of join-create-join-room may be repeated 3 times(network delay)
func (mtr *MatrixTransport) getUnlistedRoom(roomname string, invitees []*matrixcomm.UserInfo) (room *matrixcomm.Room, err error) {
	roomNameFull := "#" + roomname + ":" + DISCOVERYROOMSERVER
	var inviteesUids []string
	for _, xuser := range invitees {
		inviteesUids = append(inviteesUids, xuser.UserID)
	}
	unlistedRoomid := ""
	for i := 0; i < 6; i++ {
		respj, err := mtr.matrixcli.JoinRoom(roomNameFull, mtr.servername, nil)
		if err != nil {
			_, errc := mtr.matrixcli.CreateRoom(&matrixcomm.ReqCreateRoom{
				RoomAliasName: roomname,
				Preset:        CHATPRESET,
				Invite:        inviteesUids,
			})
			if errc != nil {
				log.Info(fmt.Sprintf("Room %s not found,trying to create it.", roomname))

				continue
			}
			continue
		} else {
			unlistedRoomid = respj.RoomID
			log.Info(fmt.Sprintf("Room joined successfully,room=%s", unlistedRoomid))
			break
		}
	}
	//if can't join nor create, create an unnamed one
	if unlistedRoomid == "" {
		respc, err := mtr.matrixcli.CreateRoom(&matrixcomm.ReqCreateRoom{
			Preset: CHATPRESET, //TODO: debug only
			Invite: inviteesUids,
		})
		if err == nil {
			unlistedRoomid = respc.RoomID
			log.Info("Could not create or join a named room. Successfuly created an unnamed one")
		}
	}
	room = &matrixcomm.Room{
		ID: unlistedRoomid,
		//（create or join just retrun roomID）what do other's content of this room do,how did matrix homeserver handle it
	}

	return
}

// setRoomIDForAddress update addresses->rooms, which is map["mark"]map[address][roomids]
func (mtr *MatrixTransport) setRoomID2Address(address common.Address, roomid string) (err error) {
	addressHex := address.String()
	if roomid != mtr.Address2Room[addressHex] {
		if roomid != "" {
			mtr.Address2Room[addressHex] = roomid
		} else {
			delete(mtr.Address2Room, addressHex)
		}
		err = mtr.matrixcli.SetAccountData(mtr.UserID, "network.raiden.rooms", mtr.Address2Room)
	}
	return
}

// getRoomID2Address : get room id of nodes from cache.
func (mtr *MatrixTransport) getRoomID2Address(address common.Address) (roomid string) {
	addressHex := address.String()
	if _, ok := mtr.Address2Room[addressHex]; !ok {
		err := mtr.setRoomID2Address(address, "")
		if err != nil {
		}
		return ""
	}
	roomid = mtr.Address2Room[addressHex]
	if roomid != "" {
		err := mtr.setRoomID2Address(address, roomid)
		if err != nil {
			log.Error(fmt.Sprintf("set room id to address err %s", err))
		}
	}
	//roomid="!OOMYBnlndieRuzkXtt:transport01.smartraiden.network"//test
	return
}

// nodeHealthCheck The purpose is to: 1、bob login from a other homeserver,i can't find him,unless bob publish his userinfo to all servers
func (mtr *MatrixTransport) nodeHealthCheck(nodeAddress common.Address) (err error) {
	if mtr.running == false {
		return
	}
	if _, ok := mtr.Address2User[nodeAddress]; ok {
		return //already healthchecked
	}
	nodeAddrHex := hexutil.Encode(nodeAddress.Bytes())
	log.Info(fmt.Sprintf("HealthCheck,peer_address=%s", nodeAddrHex))

	// check UserInfo of addressHex from server
	// fuzz check user info via partner's address.
	var tmpUserInfos []*matrixcomm.UserInfo
	respusers, err := mtr.matrixcli.SearchUserDirectory(&matrixcomm.ReqUserSearch{
		SearchTerm: nodeAddrHex,
		//Limit:10,
	})
	if err != nil {
		return
	}
	if len(respusers.Results) == 0 {
		return fmt.Errorf("%s cannot found", nodeAddress.String())
	}
	for _, resultx := range respusers.Results {
		xaddr, xerr := validateUseridSignature(resultx)
		//validate failed
		if xerr != nil {
			continue
		}
		//validate failed
		if xaddr != nodeAddress {
			continue
		}
		tmpUserInfos = append(tmpUserInfos, &resultx)
		_, verr := mtr.verifyAndUpdateUserCache(&resultx)
		if verr != nil {
		}
	}

	//cache as "Address2User"
	mtr.Address2User[nodeAddress] = tmpUserInfos

	//Ensure network state is updated in case we already know about the user presences representing the target node
	mtr.updateAddressPresence(nodeAddress, mtr.NodeDeviceType)
	return nil
}

/*
------------------------------------------------------------------------------------------------------------------------
*/

// InitMatrixTransport init matrix
func InitMatrixTransport(logname string, key *ecdsa.PrivateKey, devicetype string) (*MatrixTransport, error) {
	serverList := params.MatrixServerConfig
	var homeserverValid = ""
	var matrixclieValid = &matrixcomm.MatrixClient{}
	for _, value := range serverList {
		homeserverurl := value[0]
		homeservername := value[1]
		mcli, err := matrixcomm.NewClient(homeserverurl, "", "", PATHPREFIX0)
		if err != nil {
			continue
		}
		_, errchk := mcli.Versions()
		if errchk != nil {
			log.Error(fmt.Sprintf("Could not connect to requested server %s,and retrying", homeserverurl))
			continue
		}
		homeserverValid = homeservername
		matrixclieValid = mcli
		break
	}
	if homeserverValid == "" {
		errinfo := "Unable to find any reachable Matrix server"
		log.Error(errinfo)
		return nil, fmt.Errorf(errinfo)
	}
	mtr := &MatrixTransport{
		servername:        homeserverValid,
		running:           false,
		stopreceiving:     true,
		NodeAddress:       crypto.PubkeyToAddress(key.PublicKey),
		key:               key,
		Users:             make(map[string]*matrixcomm.UserInfo),
		Address2Room:      make(map[string]string),
		Userid2Presence:   make(map[string]*matrixcomm.RespPresenceUser),
		AddressToPresence: make(map[common.Address]*matrixcomm.RespPresenceUser),
		Address2User:      make(map[common.Address][]*matrixcomm.UserInfo),
		NodeDeviceType:    devicetype,
		UserDeviceType:    make(map[common.Address]string),
		log:               log.New("name", logname),
		avatarurl:         "", // charge rule
		statusChan:        make(chan netshare.Status, 10),
		status:            netshare.Disconnected,
	}
	mtr.matrixcli = matrixclieValid
	log.Warn(fmt.Sprintf("-->%s", homeserverValid))
	return mtr, nil
}

// validate_userid_signature
func validateUseridSignature(user matrixcomm.UserInfo) (address common.Address, err error) {
	//displayname should be an address in the self._userid_re format
	err = fmt.Errorf("validate user info failed")
	_match := ValidUserIDRegex.MatchString(user.UserID)
	if _match == false {
		return
	}
	_address, err := extractUserLocalpart(user.UserID) //"@myname:smartraiden.org:cy"->"myname"
	if err != nil {
		return
	}
	var addrmuti = regexp.MustCompile(`^(0x[0-9a-f]{40})`)
	addrlocal := addrmuti.FindString(_address)
	if addrlocal == "" {
		return
	}
	/*if len(_address) != 42 || len(user.DisplayName) != 132 {
		return
	}*/
	if _, err0 := hexutil.Decode(addrlocal); err0 != nil {
		return
	}
	if _, err0 := hexutil.Decode(user.DisplayName); err0 != nil {
		return
	}
	addressBytes := hexutil.MustDecode(addrlocal)
	useridtmp := utils.Sha3([]byte(user.UserID))                //userID's formart:  @0x....:xx
	displaynametmp := hexutil.MustDecode(user.DisplayName)      //delete "0x",to byte[]
	recovered, err := recoverData(useridtmp[:], displaynametmp) //or GetDisplayName() from server
	if err != nil {
		return
	}
	if !bytes.Equal(recovered, addressBytes) {
		addressBytes = nil
		err = fmt.Errorf("validate %s failed", user.UserID)
		return
	}
	address = common.BytesToAddress(addressBytes)
	err = nil
	return
}

// Int32ToBytes int32 to bytes
func Int32ToBytes(i int32) []byte {
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(i))
	return buf
}

// BytesToInt64 byte int64
func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

// BytesToInt32 bytes to int32
func BytesToInt32(buf []byte) int32 {
	return int32(binary.BigEndian.Uint32(buf))
}

// recover recover the node's address
func recoverData(data, signature []byte) (address []byte, err error) {
	recoverPub, err := crypto.Ecrecover(data, signature)
	if err != nil {
		return
	}
	address = utils.PubkeyToAddress(recoverPub).Bytes()
	/*addr,err:=utils.Ecrecover(utils.Sha3(data),signature)
	if err!=nil{
		return
	}
	address=addr.Bytes()*/
	return
}

// ChecksumAddress Use common.address.String() instead
func ChecksumAddress(address string) string {
	address = strings.Replace(strings.ToLower(address), "0x", "", 1)
	addressHash := hex.EncodeToString(crypto.Keccak256([]byte(address)))
	checksumAddress := "0x"
	for i := 0; i < len(address); i++ {
		l, err := strconv.ParseInt(string(addressHash[i]), 16, 16)
		if err != nil {
			return ""
		}
		if l > 7 {
			checksumAddress += strings.ToUpper(string(address[i]))
		} else {
			checksumAddress += string(address[i])
		}
	}
	return checksumAddress
}

// ExtractUserLocalpart Extract user name from user ID
func extractUserLocalpart(userID string) (string, error) {
	if len(userID) == 0 || userID[0] != '@' {
		return "", fmt.Errorf("%s is not a valid user id", userID)
	}
	return strings.TrimPrefix(
		strings.SplitN(userID, ":", 2)[0],
		"@", // remove "@" prefix
	), nil
}
