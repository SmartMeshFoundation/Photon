package network

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/matrixcomm"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"sort"
	"encoding/json"
	"time"
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
	Userid2Presence    map[string]*matrixcomm.RespPresenceUser         //cache user's real-time presence by userID("userID{presence}")
	UserID             string                                          //the current user's ID(@kitty:thisserver)
	UseDeviceType      string
	avatarurl          string
	Address2Room       map[string]map[string]string //all rooms with we knows,just
	log                log.Logger
	ChargeRegulation   string
}

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

// Stop 是否需要销毁即matrix资源？
func (mtr *MatrixTransport) Stop() {
	if mtr.running == false {
		return
	}
	mtr.running = false
	mtr.matrixcli.SetPresenceState(&matrixcomm.ReqPresenceUser{
		Presence: OFFLINE,
	})
	mtr.matrixcli.StopSync()
	if _, err := mtr.matrixcli.Logout(); err != nil {
		log.Error("[Matrix] i-node logout failed")
	}
}

// StopAccepting stop receive message and wait
func (mtr *MatrixTransport) StopAccepting() {
	mtr.stopreceiving = true
}

// NodeStatus 获取节点网络状态，如果查询自身节点，isOnline状态根据服务器握手信号来判断而非一直是true(可作为一个维护点)
func (mtr *MatrixTransport) NodeStatus(addr common.Address) (deviceType string, isOnline bool) {
	//matrix服务未启动不允许使用此接口
	if mtr.matrixcli == nil {
		return "", false
	}
	deviceType = mtr.UseDeviceType
	_, isexist := mtr.AddressToPresence[addr]
	if !isexist {
		isOnline = false
		return
	}
	if mtr.AddressToPresence[addr].Presence!=ONLINE{
		isOnline=false
	}else {
		isOnline=true//只有online的时候才会返回staus_msg(deviceType)
		deviceType=mtr.AddressToPresence[addr].StatusMsg
	}
	//在invite被查询节点到presence list的前提下可查询任何联盟服务器上的user presence
	//以上代码不能获取deviceType,通过/presence/list来处理（如果节点在线，会多返回一个status_msg（装载有deviceType）,
	//通过{userid}/presence也能查阅,可扩展（/presence/list）同时查多个节点的状态


	return
}

// Send send message
func (mtr *MatrixTransport) Send(receiverAddr common.Address, data []byte) error {
	if !mtr.running || len(data) == 0 {
		return fmt.Errorf("[Matrix]Send failed,matrix not running or send data is null")
	}
	room, err:= mtr.getRoom2Address(receiverAddr)
	if err !=nil || room==nil {
		return fmt.Errorf("[Matrix]Send failed,cann't find the object address")
	}
	_data := base64.StdEncoding.EncodeToString(data)
	resp,err := mtr.matrixcli.SendText(room.ID, _data)
	if err != nil {
		log.Trace(fmt.Sprintf("[matrix]send failed to %s, message=%s", utils.APex2(receiverAddr), encoding.MessageType(data[0])))
		fmt.Println(resp)
	} else {
		//log.Info(fmt.Sprintf("[Matrix]Send to %s, message=%s", utils.APex2(receiverAddr), encoding.MessageType(data[0])))
		log.Info(fmt.Sprintf("[Matrix]Send to %s, message=%s", utils.APex2(receiverAddr), _data))
	}
	return nil
}

// Start matrix
func (mtr *MatrixTransport) Start() {
	if mtr.running {
		return
	}
	mtr.running = true

	//登录
	if err := mtr.loginOrRegister(); err != nil {
		return
	}
	mtr.running = true
	mtr.stopreceiving=false
	//health-check,功能之一即是寻找本节点曾经加入的room（非公开room的流程?）
	err:=mtr.nodeHealthCheck(mtr.NodeAddress)
	if err!=nil{
		return
	}

	//初始化Filters/NextBatch/Rooms 均为空
	store := matrixcomm.NewInMemoryStore()
	mtr.matrixcli.Store = store

	//处理discoveryroom,测试用，暂时保留此room
	if err := mtr.joinDiscoveryRoom(); err != nil {
		return
	}
	//检索所有store-room，检查是否加入的listening room
	if err := mtr.inventoryRooms(); err != nil {
		return
	}
	//向服务器（include the other participating servers）提交本节点上线状态
	if err := mtr.matrixcli.SetPresenceState(&matrixcomm.ReqPresenceUser{
		Presence:  ONLINE,
		StatusMsg: mtr.UseDeviceType, //向服务器联盟注册使用的设备
	}); err != nil {
		return
	}
	//mtr.nodeHealthCheck(crypto.PubkeyToAddress(mtr.key.PublicKey))
	//register receive-datahandle
	mtr.matrixcli.Store = store
	mtr.matrixcli.Syncer = matrixcomm.NewDefaultSyncer(mtr.UserID, store)
	syncer := mtr.matrixcli.Syncer.(*matrixcomm.DefaultSyncer)

	syncer.OnEventType("network.raiden.rooms", mtr.onHandleAccountData)

	syncer.OnEventType("m.room.message", mtr.onHandleReceiveMessage)

	syncer.OnEventType("m.presence", mtr.onHandlePresenceChange)

	syncer.OnEventType("m.room.member",mtr.onHandleMemberShipChange)

	go func() {
		/*for {*/
		if err := mtr.matrixcli.Sync(); err != nil {
			log.Error("[Matrix] transport failed")
		}
		/*	time.Sleep(time.Second * 5)
		}*/
	}()

	log.Trace("[Matrix] transport started")
	//test code
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
	}()
}

/*
------------------------------------------------------------------------------------------------------------------------
*/
//onHandleReceiveMessage push the message of some one send "account_data"
func (mtr *MatrixTransport) onHandleAccountData(event *matrixcomm.Event) {
	if mtr.stopreceiving|| event.Type != "network.raiden.rooms"{
		return
	}
	userid := event.Sender
	if userid == mtr.UserID {
		return
	}
	value, exist:= event.ViewContent("account_data")
	if exist && value != "" {

	}
	fmt.Println("+++",event)
}
// onHandleReceiveMessage handle text messages sent to listening rooms
func (mtr *MatrixTransport) onHandleReceiveMessage(event *matrixcomm.Event) {
	if mtr.stopreceiving ||event.Type!="m.room.message"{
		return
	}
	msgtype,ok:=event.MessageType()
	if ok==false || msgtype!="m.text"{
		return
	}

	senderID := event.Sender
	roomID:=event.RoomID
	if senderID == mtr.UserID {
		return
	}

	tmpuser:=&matrixcomm.UserInfo{
		UserID:senderID,
	}
	user, err := mtr.standardizedUser(tmpuser)
	if err!=nil{
		return
	}
	peerAddress,err:=validateUseridSignature(*user)
	if err!=nil{
		log.Warn(fmt.Sprintf("Receive message from a user without legala displayName signature,peer_user=%s,room=%s",user.UserID,roomID))
		return
	}

	oldRoomID:=mtr.getRoomID2Address(peerAddress)
	if oldRoomID!=roomID{
		log.Warn(fmt.Sprintf("Receive message triggered new room for peer," +
			" peer_user=%s,peer_address=%s,old_room=%s,room=%s",user.UserID,peerAddress.String(),oldRoomID,roomID))
	}
	err=mtr.setRoomID2Address(peerAddress,roomID)

	if _,ok=mtr.Address2User[peerAddress];!ok{
		return
	}

	data,ok:=event.Body()
	if !ok|| len(data)<2{
		return
	}

	message :=[]byte{}
	if data[0:len(data)-2]=="0x"{
		message,err=hexutil.Decode(data)
		if err!=nil{
			log.Warn(fmt.Sprintf("Receive message binary data is not a valid message,message_data=%s,peer_address=%s",data,peerAddress.String()))
			return
		}

	}else {
		//解析json数据
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
		}

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
	/*msgSender, err := extractUserLocalpart(senderID)
	if err != nil {
		return
	}
	var addrmuti = regexp.MustCompile(`^(0x[0-9a-f]{40})`)
	addrlocal := addrmuti.FindString(msgSender)
	if addrlocal == "" {
		return
	}
	if _, err := hexutil.Decode(addrlocal); err != nil {
		return
	}
	msgData, ok := event.Body()
	if ok {
		dataContent, err := base64.StdEncoding.DecodeString(msgData)
		if err != nil {
			log.Error(fmt.Sprintf("[Matrix]Receive unkown message %s", utils.StringInterface(event, 0)))
		} else {
			mtr.HandleMessage(common.HexToAddress(addrlocal), dataContent)
			log.Info(fmt.Sprintf("[Matrix]Receive message %s from %s", encoding.MessageType(dataContent[0]), utils.APex2(common.HexToAddress(addrlocal))))
			log.Info(fmt.Sprintf("[Matrix]Receive message %s from %s", utils.StringInterface(event, 5), utils.APex2(common.HexToAddress(addrlocal))))

		}
	}*/
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
	user, err := mtr.standardizedUser(tmpuser)
	peerAddress, err := validateUseridSignature(*user)
	if err != nil {
		log.Info("Got invited to a room by invalid signed user - ignoring")
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

// onHandlePresenceChange 处理消息内的事件中(Content)关于节点状态的改变，刷新AddressToPresence
func (mtr *MatrixTransport) onHandlePresenceChange(event *matrixcomm.Event) {
	if mtr.stopreceiving == true {
		return
	}
	//此条消息的发送者
	userid := event.Sender
	if event.Type != "m.presence" || userid == mtr.UserID {
		return
	}
	var userDisplayname= ""
	//从Users（cache）中读取sender的sender's displayname
	tmpuser:=&matrixcomm.UserInfo{
		UserID:userid,
	}
	user, err := mtr.standardizedUser(tmpuser)
	if err == nil {
		userDisplayname = user.DisplayName
	}
	//从消息来源中获取sender's displayname
	value, exists := event.ViewContent("displayname")
	if exists && value != "" {
		userDisplayname = value
	}
	//从上述两种来源均可
	if userDisplayname == "" {
		return
	}
	user.DisplayName = userDisplayname

	//解出消息发出者的address
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

	vValue, exists := event.ViewContent("presence")
	if !exists {
		return
	}
	newstate := vValue
	//presence status unchanged
	if _,ok:=mtr.Userid2Presence[userid];!ok{
		return
	}
	oldstate:=mtr.Userid2Presence[userid].Presence
	if newstate == oldstate {
		return
	}
	//presence status hava changed
	mtr.Userid2Presence[userid].Presence = newstate
	mtr.updateAddressPresence(peerAdderss)
}

// getUserPresence get the presence state from Userid2Presence
func (mtr *MatrixTransport) getUserPresence(userid string) (presence *matrixcomm.RespPresenceUser,err error) {
	//如果user id 不存在与cache的UseridToPresence，则临时向服务器请求
	if _, ok := mtr.Userid2Presence[userid]; !ok {
		resp, err := mtr.matrixcli.GetPresenceState(userid)
		if err != nil {
			presence.Presence = UNKNOWN
		} else { //此处获取StatusMsg(deveceType)
			//presence.Presence = resp.Presence
			//presence.StatusMsg = resp.StatusMsg
			presence=resp

			//更新此user id 的presence->UseridToPresence
			mtr.Userid2Presence[userid] = presence
		}
	}
	presence = mtr.Userid2Presence[userid]
	return
}

// updateAddressPresence Update synthesized address presence state from user presence state
func (mtr *MatrixTransport) updateAddressPresence(address common.Address) {
	//一个address可能对应多个userid即多个presence
	compositepresence := []string{}
	tmpUserInfos := []*matrixcomm.UserInfo{}
	if _,ok:=mtr.Address2User[address];!ok{
		return
	}
	tmpUserInfos=mtr.Address2User[address]
	for _,v:=range tmpUserInfos{
		resp, err := mtr.getUserPresence(v.UserID)
		if err != nil {
			continue
		}
		compositepresence = append(compositepresence, resp.Presence)
	}

	//按照online、unavailable、offline、unknown顺序核对presence state
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
	if _, ok := mtr.AddressToPresence[address]; !ok {
		return
	}
	tmpuserp:=&matrixcomm.RespPresenceUser{
		Presence:newState,
		UserID:tmpUserInfos[0].UserID,
	}
	mtr.AddressToPresence[address]=tmpuserp
}

// loginOrRegister 节点登录（如果不成功，新注册再尝试登录），节点的displayname为user ID的签名
func (mtr *MatrixTransport) loginOrRegister() (err error) {
	//TODO:考虑被恶意注册的风险
	regok := false
	loginok := false
	baseAddress:=crypto.PubkeyToAddress(mtr.key.PublicKey)
	baseUsername := strings.ToLower(baseAddress.String())

	username := baseUsername
	password := hexutil.Encode(mtr.dataSign([]byte(mtr.servername)))
	for i := 0; i < 5; i++ {
		if regok == false {
			//rand.Seed(time.Now().UnixNano())
			//rnd := Int32ToBytes(rand.Int31n(math.MaxInt32))
			//username = baseUsername + "." + hex.EncodeToString(rnd)
		}
		mtr.matrixcli.AccessToken = ""
		resplogin, lerr := mtr.matrixcli.Login(&matrixcomm.ReqLogin{
			Type:     LOGINTYPE,
			User:     username,
			Password: password,
			DeviceID: "",
		})
		if lerr != nil {
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
			mtr.NodeAddress=baseAddress
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
	//把本节点的信息加入Users
	thisUser := &matrixcomm.UserInfo{
		UserID:      mtr.UserID,
		DisplayName: dispname,
		AvatarURL:   mtr.avatarurl,
	}
	_, err = mtr.standardizedUser(thisUser)

	return err
}

// inventoryRooms 整理被侦听的room，discovery room 不放入listening object（暂时的，维护时可用于不同room类别的处理）
// TODO:debug onlu
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

// makeRoomAlias
func (mtr *MatrixTransport) makeRoomAlias(thepart string) string {
	return ROOMPREFIX + ROOMSEP + NETWORKNAME + ROOMSEP + thepart
}

// dataSign 签名数据
func (mtr *MatrixTransport) dataSign(data []byte) (signature []byte) {
	hash := crypto.Keccak256(data)
	signature, err := crypto.Sign(hash[:], mtr.key)
	if err!=nil{
		return nil
	}
	return
}

// joinDiscoveryRoom 检查deicoveryroom(不存在则新建)，客户端暂存此room内的成员，再次invite从此room内检索出来的节点（如果已存在的话）
func (mtr *MatrixTransport) joinDiscoveryRoom() (err error) {
	//从配置文件中读取discovery room的名字和名字的结束标识
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
	//合成discovery room's alias
	discoveryRoomAlias := mtr.makeRoomAlias(ALIASFRAGMENT)
	discoveryRoomAliasFull := "#" + discoveryRoomAlias + ":" + DISCOVERYROOMSERVER
	mtr.discoveryroomid = ""
	//本节点加入此discovery room（不存在则创建）
	for i := 0; i < 5; i++ {
		respj, errj := mtr.matrixcli.JoinRoom(discoveryRoomAliasFull, mtr.servername, nil)
		if errj != nil {
			//if Room doesn't exist and then create the room(this is the node's resposibility)
			if mtr.servername != DISCOVERYROOMSERVER {
				log.Error(fmt.Sprintf("discovery room {%s} not found and can't be created on a federated homeserver {%s}", discoveryRoomAliasFull, mtr.servername))
				break
			}
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
		errinfo := "an error about discovery room occurred"
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
		_, xerr := mtr.standardizedUser(&usr)
		if xerr != nil {
		}
		//invite them to the discovery room
		mtr.maybeInviteUser(usr)
	}
	return
}

// maybeInviteUser 邀请节点到其所在的room(通过Address2Room搜索)
func (mtr *MatrixTransport) maybeInviteUser(user matrixcomm.UserInfo) {
	address, err := validateUseridSignature(user)
	if err != nil {
		return
	}
	roomid := mtr.getRoomID2Address(address)
	if roomid==""{
		return
	}
	room :=mtr.matrixcli.Store.LoadRoom(roomid)
	if room!=nil{
		return
	}
	if room.ID==""{
		return
	}
	//room already found the invite the user
	resp, err := mtr.matrixcli.JoinedMembers(room.ID)
	if err!=nil{
		return
	}
	//invite the user when it not in Address2Room
	if _,exist:=resp.Joined[user.UserID];!exist{
		_, err = mtr.matrixcli.InviteUser(roomid, &matrixcomm.ReqInviteUser{
			UserID: user.UserID,
		})
	}
	return
}

// getUser 把user的元素转换成UserInfo格式，先从cache的User读取 Standardized
func (mtr *MatrixTransport) standardizedUser(user0 *matrixcomm.UserInfo) (user1 *matrixcomm.UserInfo,err error) {
	//检查user ID是否合法
	_match := ValidUserIDRegex.MatchString(user0.UserID)
	if _match == false {
		user1 = nil
		err=fmt.Errorf("user id is illegal")
		return
	}
	if _, ok := mtr.Users[user0.UserID]; !ok {
		mtr.Users[user0.UserID] = user0
	}
	user1 = mtr.Users[user0.UserID]
	err=nil
	return
}

// getRoom2Address 通过节点地址获取*room对象，通讯双方在不存在已建立的room，则需要临时组建起二者用于通讯的匿名room
func (mtr *MatrixTransport) getRoom2Address(address common.Address) (room *matrixcomm.Room, err error) {
	if mtr.stopreceiving {
		return
	}
	addressHex := hexutil.Encode(address.Bytes())
	//try to get roomID from account_data from server include the other participating servers
	roomid := mtr.getRoomID2Address(address)
	if roomid != "" { //查询的对象在监听room内
		room = mtr.matrixcli.Store.LoadRoom(roomid)
		return
	}
	//以下是两两通讯room不存在的情况
	var addressOfPairs = "" //room_name由两个节点的地址组合，地址大的在前
	if mtr.NodeAddress == address {
		return
	}
	strPairs := []string{hexutil.Encode(mtr.NodeAddress.Bytes()), hexutil.Encode(address.Bytes())}
	sort.Strings(strPairs)
	addressOfPairs = strings.Join(strPairs, "_") //format 0xaaaa_0xbbbb
	tmpRoomName := mtr.makeRoomAlias(addressOfPairs)

	//从服务器（include the other participating servers）检索包含节点addressHex的UserInfo
	//模糊查询,通过对方的地址查询user info
	var tmpUserInfos []*matrixcomm.UserInfo
	respusers, err := mtr.matrixcli.SearchUserDirectory(&matrixcomm.ReqUserSearch{
		SearchTerm: addressHex,
		//Limit:10,
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
		/*//update Users
		_,xerr:=mtr.getUser(&resultx)
		if xerr!=nil{}*/
		tmpUserInfos = append(tmpUserInfos, &resultx)
	}
	//没有对方任何踪迹,不允许pees为空
	if len(tmpUserInfos) == 0 {
		return
	}

	//刷新map(address->userids)AddressToUserids
	mtr.Address2User[address] = tmpUserInfos

	//获取或创建（and invite them）当前会话的需要的匿名room（join or create and join or a unnamed-room）并invite the peers
	room, err = mtr.getUnlistedRoom(tmpRoomName, tmpUserInfos)

	//update my account_data and update RoomID2Address
	err = mtr.setRoomID2Address(address, room.ID)

	//invite users,把对方可能多个userID（分布在不同服务器上）invert 确保在新建的room里(非public room)
	for _, xuser := range tmpUserInfos {
		mtr.maybeInviteUser(*xuser)
	}
	//确保此room存在侦听任务中
	if mtr.matrixcli.Store.LoadRoom(room.ID) == nil {
		mtr.matrixcli.Store.SaveRoom(room)
	}
	log.Info(fmt.Sprintf("channel room,peer_address=%s room=%s", addressHex, room.ID))
	//fmt.Println(addressOfPairs)
	if _,ok:=mtr.Address2User[address];!ok{
		log.Info(fmt.Sprintf("address not health checked:me=%s peer_address=%s", mtr.UserID, addressHex))
	}
	return
}

// getUnlistedRoom 获取两两会话的room并invite users，如果不存在就创建一个匿名的room,且invite users
func (mtr *MatrixTransport) getUnlistedRoom(roomname string, invitees []*matrixcomm.UserInfo) (room *matrixcomm.Room, err error) {
	roomNameFull := "#" + roomname + ":" + DISCOVERYROOMSERVER
	var inviteesUids []string
	for _, xuser := range invitees {
		inviteesUids = append(inviteesUids, xuser.UserID)
	}
	unlistedRoomid := ""
	for i := 0; i < 2; i++ {
		respj, err := mtr.matrixcli.JoinRoom(roomNameFull, mtr.servername, nil)
		if err != nil {
			fmt.Println(err)
			_, errc := mtr.matrixcli.CreateRoom(&matrixcomm.ReqCreateRoom{
				RoomAliasName: roomname,
				Preset:        CHATPRESET,
				Invite:        inviteesUids,
			})
			if errc != nil {
				fmt.Println(err)
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

// setRoomIDForAddress,更新addresses->rooms设置AccountData的内容，具体是map["mark"]map[address][roomids]
func (mtr *MatrixTransport) setRoomID2Address(address common.Address, roomid string) (err error) {
	addressHex := address.String()
	//address2Room := make(map[string]string)
	xmap1:=make(map[string]interface{})
	_, ok := mtr.Address2Room["network.raiden.rooms"]
	if !ok {
		//return fmt.Errorf("can't find account_data in memory")
		//mtr.Address2Room["network.raiden.rooms"] = nil
		//xmap := make(map[string]interface{})
		//mtr.Address2Room["network.raiden.rooms"] = xmap
	}
	if roomid != "" {
		//xmap[addressHex] = roomid
	} else {
		//delete(xmap, addressHex)

		xmap1[addressHex] = roomid
		//report to server my nodes's rooms(addressehexs->rooms) and the Address2Room is latest
		/*	var reqdata interface{}
	reqdata = new(matrixcomm.ReqAccountData)
	xvalue := reflect.ValueOf(reqdata)
	if xvalue.Kind() == reflect.Ptr {
		elem := xvalue.Elem()
		accountdata := elem.FieldByName("account_data")
		if accountdata.Kind() == reflect.Map {
			*(*map[string]interface{})(unsafe.Pointer(accountdata.Addr().Pointer())) = map[string]interface{}{"addressHex":roomid,}}*/
		err = mtr.matrixcli.SetAccountData(mtr.UserID, "network.raiden.rooms", &matrixcomm.ReqAccountData{
			xmap1,
		})
	}

	return
}

// getRoomID2Address 从cache中获取节点所在room的room id
func (mtr *MatrixTransport) getRoomID2Address(address common.Address) (roomid string) {
	addressHex := address.String()
	value, exist := mtr.Address2Room["network.raiden.rooms"]
	if !exist {
		return ""
	}
	address2Room := make(map[string]string)
	address2Room = value
	if _, ok := address2Room[addressHex]; !ok {
		return ""
	}
	roomid = address2Room[addressHex]
	//check RoomID2Address 本节刷新此addressHex的信息
	//TODO: 不在监听对象中就没有任何意义
	if roomid != "" && mtr.matrixcli.Store.LoadRoom(roomid) == nil { //(cache)Store-Rooms is null
		err := mtr.setRoomID2Address(address, "")
		if err != nil {
		}
		return ""
	}
	//roomid="!OOMYBnlndieRuzkXtt:transport01.smartraiden.network"
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

	//从服务器（include the other participating servers）检索包含节点addressHex的UserInfo
	//模糊查询,通过对方的地址查询user info
	var tmpUserInfos []*matrixcomm.UserInfo
	respusers, err := mtr.matrixcli.SearchUserDirectory(&matrixcomm.ReqUserSearch{
		SearchTerm: nodeAddrHex,
		//Limit:10,
	})
	if err != nil {
		return
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
		/*//update Users
		_,xerr:=mtr.getUser(&resultx)
		if xerr!=nil{}*/
		tmpUserInfos = append(tmpUserInfos, &resultx)
		mtr.standardizedUser(&resultx)
	}

	//刷新map(address->userids)AddressToUserids
	mtr.Address2User[nodeAddress] = tmpUserInfos

	//Ensure network state is updated in case we already know about the user presences representing the target node
	mtr.updateAddressPresence(nodeAddress)
	return nil
}

/*
------------------------------------------------------------------------------------------------------------------------
*/

const (
	// ONLINE network state -online
	ONLINE      = "online"
	// UNAVAILABLE state -unavailable
	UNAVAILABLE = "unavailable"
	// OFFLINE state -offline
	OFFLINE     = "offline"
	// UNKNOWN or other state -unknown
	UNKNOWN     = "unknown"
	// ROOMPREFIX room prefix
	ROOMPREFIX  = "raiden"
	// ROOMSEP with ',' to separate room name's part
	ROOMSEP     = "_"
	// PATHPREFIX0 the lastest matrix client api version
	PATHPREFIX0 = "/_matrix/client/r0"
	// AUTHTYPE login identity as dummy
	AUTHTYPE    = "m.login.dummy"
	// LOGINTYPE login type we used
	LOGINTYPE   = "m.login.password"
	// CHATPRESET the type of chat=public
	CHATPRESET  = "public_chat"
)

var (
	// ValidUserIDRegex user ID 's format
	ValidUserIDRegex    = regexp.MustCompile(`^@(0x[0-9a-f]{40})(?:\.[0-9a-f]{8})?(?::.+)?$`) //(`^[0-9a-z_\-./]+$`)
	//NETWORKNAME which network is used
	NETWORKNAME         = "ropsten"
	//ALIASFRAGMENT the terminal part of alias
	ALIASFRAGMENT       = ""
	//DISCOVERYROOMSERVER discovery room server name
	DISCOVERYROOMSERVER = ""
)

// InitMatrixTransport init matrix
func InitMatrixTransport(logname string, key *ecdsa.PrivateKey, devicetype string) (*MatrixTransport, error) {
	serverList := params.MatrixServerConfig
	var homeserverValid= ""
	var matrixclieValid= &matrixcomm.MatrixClient{}
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
		Address2Room:      make(map[string]map[string]string),
		Userid2Presence:   make(map[string]*matrixcomm.RespPresenceUser),
		AddressToPresence: make(map[common.Address]*matrixcomm.RespPresenceUser),
		Address2User:      make(map[common.Address][]*matrixcomm.UserInfo),
		UseDeviceType:     "user_device_type",
		log:               log.New("name", logname),
		avatarurl:         "", //收费规则
	}
	mtr.matrixcli = matrixclieValid
	return mtr, nil
}

// validate_userid_signature
func validateUseridSignature(user matrixcomm.UserInfo) (address common.Address, err error) {
	//displayname should be an address in the self._userid_re format
	err = fmt.Errorf("validate user info failed");
	_match := ValidUserIDRegex.MatchString(user.UserID)
	if _match == false {
		return
	}
	_address, err := extractUserLocalpart(user.UserID) //"@myname:smartraiden.org:cy"->"myname"
	if err != nil {
		return
	}
	var addrmuti= regexp.MustCompile(`^(0x[0-9a-f]{40})`)
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
	useridtmp := utils.Sha3([]byte(user.UserID))                //userid 格式:  @0x....:xx
	displaynametmp := hexutil.MustDecode(user.DisplayName)      //去掉0x转byte[]
	recovered, err := recoverData(useridtmp[:], displaynametmp) //或者临时读取服务器上的GetDisplayName（）
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
		if err!=nil{
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

// ExtractUserLocalpart 从userID中提取username"@xxxx:"->"xxxx"
func extractUserLocalpart(userID string) (string, error) {
	if len(userID) == 0 || userID[0] != '@' {
		return "", fmt.Errorf("%s is not a valid user id", userID)
	}
	return strings.TrimPrefix(
		strings.SplitN(userID, ":", 2)[0],
		"@", // remove "@" prefix
	), nil
}
