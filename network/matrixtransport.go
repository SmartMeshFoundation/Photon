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
)

// MatrixTransport represents a matrix transport Instantiation
type MatrixTransport struct {
	matrixcli          *matrixcomm.MatrixClient //the instantiated matrix
	servername         string                   //the homeserver's name
	running            bool                     //running status
	stopreceiving      bool                     //Whether to stop accepting(data)
	key                *ecdsa.PrivateKey        //key
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
	Address2Room       map[string]string //all rooms with we knows,just
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
		return fmt.Errorf("[Matrix]Send failed,cann't find the obj addr")
	}
	_data := base64.StdEncoding.EncodeToString(data)
	resp,err := mtr.matrixcli.SendText(room.ID, _data)
	if err != nil {
		log.Trace(fmt.Sprintf("[matrix]send failed to %s, message=%s", utils.APex2(receiverAddr), encoding.MessageType(data[0])))
		fmt.Println(resp)
	} else {
		log.Info(fmt.Sprintf("[Matrix]Send to %s, message=%s", utils.APex2(receiverAddr), encoding.MessageType(data[0])))
	}
	return nil
}

// Start matrix
func (mtr *MatrixTransport) Start() {
	if mtr.running {
		return
	}
	//health-check,功能之一即是寻找本节点曾经加入的room（非公开room的流程?）

	//登录
	if err := mtr.loginOrRegister(); err != nil {
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
	syncer.OnEventType("m.room.message", mtr.onHandleReceiveMessage)

	syncer.OnEventType("m.presence", mtr.onHandlePresenceChange)

	go func() {
		/*for {*/
		if err := mtr.matrixcli.Sync(); err != nil {
			log.Error("[Matrix] transport failed")
		}
		/*	time.Sleep(time.Second * 5)
		}*/
	}()

	mtr.running = true
	log.Trace("[Matrix] transport started")
	/*//test code
	go func() {
		for {
			sdata,_:=base64.StdEncoding.DecodeString("EQAAAIIUOo4Q4ck63CN6xdsHD0yAGoxPQ65Z0QMujNykGqzlAAAAAAAAAB4FhbWJbOJl3FIhxt+EWLjGZ2htMhTgi4XAflrhlNJvXAAAAAAAMJuDAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAGJYT+hStqK3NAnZFgqukPTzxFNy1+CbYtR014/+6sGnfoxumX8Eer2XqlwxRx2pwF02ItKOL3koK3k22hmZJrQb")
			mtr.Send(common.HexToAddress("0xc67f23ce04ca5e8dd9f2e1b5ed4fad877f79267a"), sdata)
			time.Sleep(time.Second * 10)
		}
	}()*/
}

/*
------------------------------------------------------------------------------------------------------------------------
*/

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
	user, err := mtr.getUser(tmpuser)
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
	adderss, err := validateUseridSignature(*user)
	if err != nil {
		return
	}
	//not a user we've started healthcheck, skip--??
	if _, ok := mtr.Address2User[adderss]; !ok {
		return
	}
	mtr.Address2User[adderss] = append(mtr.Address2User[adderss], user)

	//maybe inviting user used to also possibly invite user's from discovery presence changes
	mtr.maybeInviteUser(*user)

	presenceValue, exists := event.ViewContent("presence")
	if !exists {
		return
	}
	newstate := presenceValue
	//presence status unchanged
	if newstate == mtr.Userid2Presence[userid].Presence {
		return
	}
	//change status
	mtr.Userid2Presence[userid].Presence = newstate
	mtr.updateAddressPresence(adderss)
}

// onHandleReceiveMessage 处理接收到room的消息
func (mtr *MatrixTransport) onHandleReceiveMessage(event *matrixcomm.Event) {
	_msgSender := event.Sender
	if _msgSender == mtr.UserID {
		return
	}
	msgSender, err := extractUserLocalpart(_msgSender)
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
		}
	}
}

// getUserPresence get the presence state from Userid2Presence
func (mtr *MatrixTransport) getUserPresence(userid string) (presence *matrixcomm.RespPresenceUser) {
	//如果user id 不存在与cache的UseridToPresence，则临时向服务器请求
	if _, ok := mtr.Userid2Presence[userid]; !ok {
		resp, err := mtr.matrixcli.GetPresenceState(userid)
		if err != nil {
			presence.Presence = UNKNOWN
		} else {//全文尽在此处获取StatusMsg(deveceType)
			presence.Presence = resp.Presence
			presence.StatusMsg=resp.StatusMsg
			//presence=resp

			//更新此user id 的presence->UseridToPresence
			mtr.Userid2Presence[userid] = presence
		}
	}
	presence=mtr.Userid2Presence[userid]
	return
}

// updateAddressPresence Update synthesized address presence state from user presence state
func (mtr *MatrixTransport) updateAddressPresence(address common.Address) {
	//一个address可能对应多个userid即多个presence
	compositepresence:=[]string{}
	for _, xuser := range mtr.Address2User[address] {//同时可能获取到StatusMsg
		compositepresence=append(compositepresence,mtr.getUserPresence(xuser.UserID).Presence)
	}

	//按照online、unavailable、offline、unknown顺序核对presence state
	presencestates:=[]string{ONLINE,UNAVAILABLE,OFFLINE,UNKNOWN}
	newState := UNKNOWN
	for _,xstate:=range presencestates{
		for _, xpresence := range compositepresence{
			if xpresence==xstate{
				newState=xpresence
				break
			}
		}
	}

	//update AddressToPresence
	if newState == mtr.AddressToPresence[address].Presence {
		return
	}
	mtr.AddressToPresence[address].Presence = newState
}

// loginOrRegister 节点登录（如果不成功，新注册再尝试登录），节点的displayname为user ID的签名
func (mtr *MatrixTransport) loginOrRegister() (err error) {
	regok := false
	loginok := false
	baseUsername := strings.ToLower(crypto.PubkeyToAddress(mtr.key.PublicKey).String())

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
	_, err = mtr.getUser(thisUser)

	return err
}

// inventoryRooms 整理被侦听的room，discovery room 不放入listening object（暂时的，维护时可用于不同room类别的处理）
func (mtr *MatrixTransport) inventoryRooms() (err error) {
	for _, value := range mtr.matrixcli.Store.LoadRoomOfAll() {
		if value.Alias == mtr.discoveryroomalias {
			continue
		}
		if mtr.matrixcli.Store.LoadRoom(value.ID) == nil {
			mtr.matrixcli.Store.SaveRoom(value)
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
			var _visibility= "private"
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
	//try to commit the discovery room to memory
	mtr.discoveryroomalias = discoveryRoomAlias

	/*//add discovery room to listening object
	theroom := &matrixcomm.Room{
		ID:    mtr.discoveryroomid,
		Alias: discovery_room_alias,
		//State:nil,
	}
	mtr.matrixcli.Store.SaveRoom(theroom)*/

	//把discovery room放入RoomID2Address
	userAddr := crypto.PubkeyToAddress(mtr.key.PublicKey)
	err=mtr.setRoomID2Address(userAddr, mtr.discoveryroomid)

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
		_,xerr:=mtr.getUser(&usr)
		if xerr!=nil{}

		//invite them to the discovery room
		mtr.maybeInviteUser(usr)
	}
	return
}

// maybeInviteUser 邀请节点到其所在的room
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
	/*//room
	resp, err := mtr.matrixcli.JoinedMembers(room.ID)
	if err!=nil{
		for
	}
	for userid, userdata := range resp.Joined {

	}*/
	if room.ID != "" {
		_, err = mtr.matrixcli.InviteUser(roomid, &matrixcomm.ReqInviteUser{
			UserID: user.UserID,
		})
	}
	return
}

// getUser 把user的元素转换成UserInfo格式，先从cache的User读取
func (mtr *MatrixTransport) getUser(user0 *matrixcomm.UserInfo) (user1 *matrixcomm.UserInfo,err error) {
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

// getRoom2Address 通过节点地址获取*room对象，与getRoomID2Address差异性的是，通讯双方在不存在已建立的room，则需要临时组建起二者用于通讯的匿名room
func (mtr *MatrixTransport) getRoom2Address(address common.Address) (room *matrixcomm.Room, err error) {
	addressHex := hexutil.Encode(address.Bytes())
	if !common.IsHexAddress(addressHex) {
		return
	}
	roomid := mtr.getRoomID2Address(address)
	if roomid != "" { //查询的对象在监听room内
		room = mtr.matrixcli.Store.LoadRoom(roomid)
		return
	}
	//以下是两两通讯room不存在的情况
	var addressOfPairs = "" //room_name由两个节点的地址组合，地址大的在前
	addrint1 := BytesToInt32(crypto.PubkeyToAddress(mtr.key.PublicKey).Bytes())
	addrint2 := BytesToInt32(address.Bytes())
	addrstr1 := hexutil.Encode(crypto.PubkeyToAddress(mtr.key.PublicKey).Bytes())
	addrstr2 := hexutil.Encode(address.Bytes())
	if addrint1 == addrint2 {
		return
	}
	if addrint1 < addrint2 {
		addressOfPairs = addrstr1 + "_" + addrstr2
	} else {
		addressOfPairs = addrstr2 + "_" + addrstr1
	}
	//format 0xaaaa_0xbbbb
	tmpRoomName := mtr.makeRoomAlias(addressOfPairs)

	var tmpUserInfos []*matrixcomm.UserInfo
	//从服务器（include the other participating servers）检索包含节点addressHex的UserInfo
	//模糊查询,通过对方的地址查询user info
	respusers, err := mtr.matrixcli.SearchUserDirectory(&matrixcomm.ReqUserSearch{
		SearchTerm: addressHex,
		Limit:10,
	})
	if err != nil {
		return
	}
	for _, resultx := range respusers.Results {
		/*xaddr,err:=validateUseridSignature(resultx)
		if err!=nil{
			continue
		}
		if xaddr!=address{
			continue
		}*/
		//更新Users
		_,xerr:=mtr.getUser(&resultx)
		if xerr!=nil{}

		tmpUserInfos=append(tmpUserInfos, &resultx)
	}
	//没有对方任何踪迹,不允许pees为空
	if len(tmpUserInfos)==0{
		return
	}
	//获取或创建（invite them）当前会话的需要的匿名room（join or create and join or a unnamed-room）并invite the peers
	room, err = mtr.getUnlistedRoom(tmpRoomName, tmpUserInfos)

	err=mtr.setRoomID2Address(address,room.ID)
	//刷新map(address->userids)AddressToUserids
	mtr.Address2User[address]=tmpUserInfos

	//再次invite users,确保在新建的room里(非public room呢)
	for _,xuser:=range tmpUserInfos{
		mtr.maybeInviteUser(*xuser)
	}
	//确保此room存在侦听任务中
	if mtr.matrixcli.Store.LoadRoom(room.ID) == nil {
		mtr.matrixcli.Store.SaveRoom(room)
	}
	log.Info("channel room ")
	fmt.Println(addressOfPairs)

	return
}

// getUnlistedRoom 获取两两会话的room并invite users，如果不存在就创建一个匿名的room,且invite users
func (mtr *MatrixTransport) getUnlistedRoom(roomname string, invitees []*matrixcomm.UserInfo) (room *matrixcomm.Room, err error) {
	roomNameFull := "#"+roomname + ":" + DISCOVERYROOMSERVER
	var inviteesUids []string
	for _,xuser:= range invitees{
		inviteesUids=append(inviteesUids, xuser.UserID)
	}
	unlistedRoomid := ""
	for i := 0; i < 3; i++ {
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
				continue
			}
			continue
		} else {
			unlistedRoomid = respj.RoomID
			break
		}
	}
	//如果创建或者join room 失败，则创建一个匿名room并invite peers,尝试多次？
	if unlistedRoomid == "" {
		respc, err := mtr.matrixcli.CreateRoom(&matrixcomm.ReqCreateRoom{
			Preset: CHATPRESET, //debug
			Invite: inviteesUids,
		})
		if err == nil {
			unlistedRoomid = respc.RoomID
			log.Info("Could not create or join a named room. Successfuly created an unnamed one")
		}
	}
	room = &matrixcomm.Room{
		ID: unlistedRoomid,
	}

	return
}

// setRoomIDForAddress,更新addresses->rooms设置AccountData的内容，具体是map["mark"]map[address][roomids]
func (mtr *MatrixTransport) setRoomID2Address(address common.Address, roomid string) (err error) {
	addressHex := address.String()
	if roomid!=mtr.Address2Room[addressHex]{
		if roomid!=""{
			mtr.Address2Room[addressHex] = roomid
		}else {
			delete(mtr.Address2Room, addressHex)
		}
	}
	//report node's rooms(address_hex->rooms)
	err = mtr.matrixcli.SetAccountData(mtr.UserID, "network.raiden.rooms", &matrixcomm.ReqAccountData{
		Addresshex: addressHex,
		Roomid:     []string{roomid},
	})
	return
}

// getRoomID2Address 获取节点所在的room id
func (mtr *MatrixTransport) getRoomID2Address(address common.Address) (roomid string) {
	addressHex := address.String()
	roomid = mtr.Address2Room[addressHex]
	if roomid!="" && mtr.matrixcli.Store.LoadRoom(roomid) == nil { //(cache)Store-Rooms is null
		err:=mtr.setRoomID2Address(address, roomid)
		if err!=nil{
			return
		}
		roomid = ""
	}
	//roomid="!OOMYBnlndieRuzkXtt:transport01.smartraiden.network"
	return
}

// nodeHealthCheck Legality check
func (mtr *MatrixTransport) nodeHealthCheck(nodeAddress common.Address) (err error) {
	if mtr.running == false {
		return
	}
	//nodeAddrHex := hexutil.Encode(nodeAddress.Bytes())
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
		stopreceiving:     false,
		key:               key,
		Users:             make(map[string]*matrixcomm.UserInfo),
		Address2Room:      make(map[string]string),
		Userid2Presence:   make(map[string]*matrixcomm.RespPresenceUser),
		AddressToPresence: make(map[common.Address]*matrixcomm.RespPresenceUser),
		Address2User:   make(map[common.Address][]*matrixcomm.UserInfo),
		UseDeviceType:     "user_device_type",
		log:               log.New("name", logname),
		avatarurl:         "",//收费规则
	}
	mtr.matrixcli = matrixclieValid
	return mtr, nil
}

// validate_userid_signature
func validateUseridSignature(user matrixcomm.UserInfo) (address common.Address, err error) {
	//displayname should be an address in the self._userid_re format
	err=fmt.Errorf("validate user info failed");
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
	useridtmp := utils.Sha3([]byte(user.UserID))            //userid 格式:  @0x....:xx
	displaynametmp := hexutil.MustDecode(user.DisplayName)  //去掉0x转byte[]
	recovered, err := recoverData(useridtmp[:], displaynametmp) //或者临时读取服务器上的GetDisplayName（）
	if err != nil {
		return
	}
	if !bytes.Equal(recovered, addressBytes) {
		addressBytes = nil
		err = fmt.Errorf("validate %s failed", user.UserID)
		return
	}
	address=common.BytesToAddress(addressBytes)
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
