package network

import (
	"crypto/ecdsa"
	"strings"
	"fmt"
	"github.com/SmartMeshFoundation/SmartRaiden/network/matrixcomm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"regexp"
	"math/rand"
	"time"
	"encoding/hex"
	"bytes"
	"math"
	"encoding/binary"
	"github.com/ethereum/go-ethereum/common"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"strconv"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"encoding/base64"
)

type MatrixTransport struct {
	matrixcli          *matrixcomm.MatrixClient //the instantiated matrix
	servername         string                   //the homeserver's name
	running            bool                     //running status
	stopreceiving      bool                     //Whether to stop accepting(data)
	key                *ecdsa.PrivateKey        //key
	protocol           ProtocolReceiver
	discoveryroomalias string                  //the room's alias of sys pre-configured ("#[RoomNameLocalpart]:[ServerName]")
	discoveryroomid    string                  //the room's ID of sys pre-configured ("![RoomIdData]:[ServerName]")
	Users              map[string]*User        //cache user's base-infos("userID{userID,displayname}")
	AddressToRoomid    map[string]string       //all rooms with we knows this service
	UseridToPresence   map[string]string       //cache user's real-time presence by userID("userID{presence}")
	AddressToPresence  map[common.Address]bool //cache user's real-time presence by node's address("userID{presence}")
	UserId             string                  //the current user's ID(@kitty:thisserver)
	UseDeviceType      string
}
/*
------------------------------------------------------------------------------------------------------------------------
*/
func (mtr *MatrixTransport) HandleMessage(from common.Address, data []byte) {
	log.Trace(fmt.Sprintf("received from %s, message=%s", utils.APex2(from), encoding.MessageType(data[0])))
	if !mtr.running || mtr.stopreceiving {
		return
	}
	if mtr.protocol != nil {
		mtr.protocol.receive(data)
	}
}

func (mtr *MatrixTransport) RegisterProtocol(protcol ProtocolReceiver) {
	mtr.protocol = protcol
}

func (mtr *MatrixTransport) Send(receiverAddr common.Address, data []byte) {
	if !mtr.running {
		return
	}
	roomID := mtr.get_room_id_for_address(receiverAddr)
	if roomID == "" {
		return
	}
	if len(data) == 0 {
		return
	}
	_data := base64.StdEncoding.EncodeToString(data)
	_, err := mtr.matrixcli.SendText(roomID, _data)
	if err != nil {
		log.Trace(fmt.Sprintf("[matrix]send to %s, message=%s [error]", utils.APex2(receiverAddr), encoding.MessageType(data[0])))
	} else {
		log.Info(fmt.Sprintf("[Matrix]Send to %s, message=%s [succeed]", utils.APex2(receiverAddr), encoding.MessageType(data[0])))
	}
}

func (mtr *MatrixTransport) Start() {
	mtr.running = true
	if err := mtr.login_or_register(); err != nil {
		return
	}
	//init cache store
	store := matrixcomm.NewInMemoryStore()
	mtr.matrixcli.Store = store
	//
	if err := mtr.join_discovery_room(); err != nil {
		return
	}
	//
	if err := mtr.inventory_rooms(); err != nil {
		return
	}
	//notify my Presence State
	err := mtr.matrixcli.SetPresenceState(&matrixcomm.ReqPresenceUser{
		Presence: ONLINE,
	})
	if err != nil {
		return
	}
	mtr.node_health_check(crypto.PubkeyToAddress(mtr.key.PublicKey))
	//register receive-datahandle
	mtr.matrixcli.Store = store
	mtr.matrixcli.Syncer = matrixcomm.NewDefaultSyncer(mtr.UserId, store)
	syncer := mtr.matrixcli.Syncer.(*matrixcomm.DefaultSyncer)
	syncer.OnEventType("m.room.message", func(evt *matrixcomm.Event) {
		_msgSender := evt.Sender
		msgSender, _ := matrixcomm.ExtractUserLocalpart(_msgSender)
		var addrmuti = regexp.MustCompile(`^(0x[0-9a-f]{40})`)
		addrlocal := addrmuti.FindString(msgSender)
		if addrlocal == "" {
			return
		}
		if _, err := hexutil.Decode(addrlocal); err != nil {
			return
		}
		msgData, ok := evt.Body() //text msgSender="0xc67f23ce04ca5e8dd9f2e1b5ed4fad877f79267a"
		if ok {
			dataContent, err := base64.StdEncoding.DecodeString(msgData)
			if err != nil {
				log.Error(fmt.Sprintf("[Matrix]Receive unkown message %s", utils.StringInterface(evt, 0)))
			} else {
				//log.Info(fmt.Sprintf("[Matrix]Receive message=%s,from=%s",msgData,msgSender))
				mtr.HandleMessage(common.HexToAddress(msgSender), dataContent)
			}
		}
	})
	//(gorountime)run type with no blocking
	go func() {
		for {
			if err := mtr.matrixcli.Sync(); err != nil {
				fmt.Println("Sync() returned ", err)
			}
			time.Sleep(time.Second * 5)
		}
	}()
}

func (mtr *MatrixTransport) Stop() {
	if mtr.running==false {
		return
	}
	mtr.running=false
	mtr.matrixcli.SetPresenceState(&matrixcomm.ReqPresenceUser{
		Presence: OFFLINE,
	})
	mtr.matrixcli.StopSync()
	if _,err:=mtr.matrixcli.Logout();err!=nil{
		fmt.Println(err)
	}
}

func (mtr *MatrixTransport) StopAccepting()  {
	mtr.stopreceiving=true
}

func (mtr *MatrixTransport) NodeStatus(addr common.Address) (deviceType string, isOnline bool) {
	if mtr.matrixcli == nil {
		return "", false
	}
	ret, ok := mtr.AddressToPresence[addr]
	if !ok{
		isOnline=false
		return
	}
	deviceType=mtr.UseDeviceType
	isOnline=ret
	return
}
/*
------------------------------------------------------------------------------------------------------------------------
*/
func (mtr *MatrixTransport) get_user_presence(userid string) string {
	presence := UNKNOWN
	if _, ok := mtr.UseridToPresence[userid]; !ok {
		resp, err := mtr.matrixcli.GetPresenceState(userid)
		if err != nil {
			presence = UNAVAILABLE
		}else {
			presence = resp.Presence
			mtr.UseridToPresence[userid] = presence
		}
	}
	return mtr.UseridToPresence[userid]
}

func (mtr *MatrixTransport) update_address_presence(address []byte) {

}

func (mtr *MatrixTransport) handle_presence_change(event matrixcomm.Event){
	if mtr.running==false{
		return
	}
	userid:=event.Sender
	if event.Type!="m.presence" || userid==mtr.UserId{
		return
	}
	user,err:=mtr.get_user(userid,"")
	if err!=nil{
		return
	}
	displayname_value,exists:=event.ViewContent("displayname")
	if !exists{
		return
	}
	user.display_name=displayname_value
	adderss,err:=validate_userid_signature(*user)
	if err!=nil{
		return
	}
	mtr.maybe_invite_user(*user)
	presence_value,exists:=event.ViewContent("presence")
	if !exists{
		return
	}
	newstate:=presence_value
	if newstate==mtr.UseridToPresence[userid]{
		return
	}
	//change status
	mtr.UseridToPresence[userid]=newstate
	mtr.update_address_presence(adderss)
}

func (mtr *MatrixTransport) login_or_register() (_err error) {
	regok := false
	loginok:=false
	baseUsername := strings.ToLower(crypto.PubkeyToAddress(mtr.key.PublicKey).String())
	username := baseUsername
	password := hexutil.Encode(mtr._sign([]byte(mtr.servername)))
	for i := 0; i < 5; i++ {
		if regok==false {
			rand.Seed(time.Now().UnixNano())
			rnd := Int32ToBytes(rand.Int31n(math.MaxInt32))
			username = baseUsername + "." + hex.EncodeToString(rnd)
			//username="0xc67f23ce04ca5e8dd9f2e1b5ed4fad877f79267a.59a2bb27"//test data
		}
		mtr.matrixcli.AccessToken = ""
		resplogin, err := mtr.matrixcli.Login(&matrixcomm.ReqLogin{
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
				log.Trace(fmt.Sprintf("couldn't sign in for matrix,trying register %s", username))
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
				_, uia, err := mtr.matrixcli.Register(req)
				if err != nil && uia == nil {
					rhttpErr, _ := err.(matrixcomm.HTTPError)
					if rhttpErr.Code == 400 { //M_USER_IN_USE,M_INVALID_USERNAME,M_EXCLUSIVE
						log.Trace("username taken,continuing")
						continue
					}
				}
				log.Trace(fmt.Sprintf("register ok,Username=%s,Password=", username,password))
				regok = true;
				mtr.matrixcli.UserID = username
				continue
			}
		}else {
			//cache the node's and report the UserID and AccessToken to matrix
			mtr.matrixcli.SetCredentials(resplogin.UserID, resplogin.AccessToken)
			mtr.UserId=resplogin.UserID
			loginok=true
			break
		}
	}
	if(!loginok){
		_err=fmt.Errorf("could not register or login")
		return
	}
	//set displayname as publicly visible(=0x......)
	dispname:=hexutil.Encode(mtr._sign([]byte(mtr.matrixcli.UserID)))
	if err:=mtr.matrixcli.SetDisplayName(dispname);err!=nil{
		_err=fmt.Errorf("could set the node's displayname and quit as well")
		mtr.matrixcli.ClearCredentials()
		return
	}
	mtr.get_user(mtr.UserId,dispname)

	return _err
}

func (mtr *MatrixTransport) inventory_rooms() (err error) {
	for _, room := range mtr.matrixcli.Store.LoadRoomOfAll() {
		if _, ok := mtr.Users[mtr.discoveryroomalias]; ok {
			//room.
			fmt.Println("room:", room)
		}
	}
	return nil
}

func (mtr *MatrixTransport) make_room_alias(thepart string) (string) {
	return ROOMPREFIX + ROOMSEP + NETWORKNAME + ROOMSEP + thepart
}

func (mtr *MatrixTransport) _sign(data []byte)(signature []byte) {
	hash := crypto.Keccak256(data)
	signature, _ = crypto.Sign(hash[:], mtr.key);
	return
}

func (mtr *MatrixTransport) join_discovery_room() (err error) {
	discovery_room_alias := mtr.make_room_alias(ALIASFRAGMENT)
	//e.g:"#smartraiden_mainnet_discovery:transport01.smartraiden.network"
	discovery_room_alias_full := "#" + discovery_room_alias + ":" + DISCOVERYROOMSERVER
	respj, err := mtr.matrixcli.JoinRoom(discovery_room_alias_full, mtr.servername, nil)
	httpErr, _ := err.(matrixcomm.HTTPError)
	if httpErr.Code == 404 {
		//Room doesn't exist and create the room(this is the node's resposibility)
		if mtr.servername != DISCOVERYROOMSERVER {
			log.Error(fmt.Sprintf("discovery room {%s} not found and can't be created on a federated homeserver {%s}", discovery_room_alias_full, mtr.servername))
			err = fmt.Errorf("can't find or create a discovery room,config error?")
			return
		}
		respc, errc := mtr.matrixcli.CreateRoom(&matrixcomm.ReqCreateRoom{
			RoomAliasName: discovery_room_alias,
			Preset:        CHATPRESET,
		}) //after create room but you were in there(server bug?)
		if errc != nil {
			err = errc
			log.Error("can't create a discovery room")
			return
		}
		mtr.discoveryroomid = respc.RoomID
	} else {
		if err != nil {
			mtr.discoveryroomid = ""
		} else {
			mtr.discoveryroomid = respj.RoomID
		}
	}
	if mtr.discoveryroomid == "" {
		errinfo := "an error about discovery room occurred"
		err = fmt.Errorf(errinfo)
		log.Error(errinfo)
		return
	}
	//mtr.discoveryroomalias
	mtr.discoveryroomalias = discovery_room_alias
	//mtr.discoveryroomid
	//mtr.discoveryroomid=respj.RoomID||respc.RoomID

	//mtr.matrixcli.Store.Rooms
	theroom := &matrixcomm.Room{
		ID: mtr.discoveryroomid,
		//State:nil,
	}
	mtr.matrixcli.Store.SaveRoom(theroom)

	respin, err := mtr.matrixcli.JoinedMembers(mtr.discoveryroomid)
	if err != nil {
		log.Error("The node can't join room ", mtr.discoveryroomalias)
		return
	}
	for userid, userdata := range respin.Joined { //cache known users
		mtr.get_user(userid, *userdata.DisplayName)
		//invite to room
		usr := User{
			user_id:      userid,
			display_name: *userdata.DisplayName,
		}
		mtr.maybe_invite_user(usr)
	}

	return
}

func (mtr *MatrixTransport) maybe_invite_user(user User) (err error) {
	address, err := validate_userid_signature(user)
	if err != nil {
		return fmt.Errorf("validate user info failed")
	}
	roomid := mtr.get_room_id_for_address(common.BytesToAddress(address))
	if roomid == "" {
		_, _ = mtr.matrixcli.InviteUser("!DRmVhWXqYqetqiihpY:cy", &matrixcomm.ReqInviteUser{
			UserID: user.user_id,
		})
	}
	return nil
}

func (mtr *MatrixTransport) get_user(userid,displayname string) (user *User,err error) {
	_match := ValidUserIDRegex.MatchString(userid)
	if _match == false {
		user = nil
		err = fmt.Errorf("%s is not a valid user id", userid)
		return
	}
	if _, ok := mtr.Users[userid]; !ok {
		usr := &User{
			user_id:      userid,
			display_name: displayname,
		}
		mtr.Users[userid] = usr
	}
	user = mtr.Users[userid]
	err = nil
	return
}

func (mtr *MatrixTransport) set_room_id_for_address(address common.Address,roomid string) (err error) {
	addressHex := ChecksumAddress(hexutil.Encode(address[:]))
	//address_to_room_id := mtr.OfficialRooms
	if _, ok := mtr.AddressToRoomid[addressHex]; !ok {
		if roomid != "" {
			mtr.AddressToRoomid[addressHex] = roomid
		} else {
			delete(mtr.AddressToRoomid, addressHex)
		}
	}
	//report user's rooms
	mtr.matrixcli.SetAccountData(mtr.UserId, "network0.smatrraiden.rooms", &matrixcomm.ReqAccountData{
		Addresshex: addressHex,
		Roomid:     roomid,
	})

	theroom := &matrixcomm.Room{
		ID: roomid,
		//State:nil,
	}
	mtr.matrixcli.Store.SaveRoom(theroom)
	return
}

func (mtr *MatrixTransport) get_room_id_for_address(address common.Address) (roomid string) {
	addressHex := ChecksumAddress(hexutil.Encode(address[:]))
	roomid = mtr.AddressToRoomid[addressHex]
	if mtr.matrixcli.Store.LoadRoom(roomid) == nil { //Store-Rooms is null
		mtr.set_room_id_for_address(address, roomid)
		roomid = ""
	}
	return
}

func (mtr *MatrixTransport) node_health_check(nodeAddress common.Address) (err error) {
	if mtr.running == false {
		return
	}
	//nodeAddrHex := hexutil.Encode(nodeAddress.Bytes())
	return nil
}

func (mtr *MatrixTransport) get_use_devicetype() (rtn string) {
	deviceType := []string{"mobile", "meshbox", "pc", "other"}
	rtn = deviceType[3]
	resp, err := mtr.matrixcli.GetWhois(mtr.UserId);
	if err != nil {
		return
	}
	ismobile := 0
	tmpRtn := fmt.Sprint(resp)
	for _, x := range mobilefeature {
		if _index := strings.Index(tmpRtn, x); _index != -1 {
			ismobile++
		}
	}
	if ismobile == 0 {
		rtn = deviceType[0]
	}
	return
}
/*
------------------------------------------------------------------------------------------------------------------------
*/
/*"discovery_room": {
	"alias_fragment": "discovery",
	"server": "transport01.smartraiden.network",
},
servername ==matrix.local.raiden
"available_servers": [
	"https://transport01.smartraiden.network",
	"https://transport02.smartraiden.network",
	"https://transport03.smartraiden.network"
],*/
type User struct {
	user_id string
	display_name string
}

const(
	ONLINE 		= "online"
	OFFLINE 	= "offline"
	UNAVAILABLE = "unavailable"
	UNKNOWN 	= "unknown"
	ROOMPREFIX	= "smartraiden"
	ROOMSEP		= "_"
	PATHPREFIX0	= "/_matrix/client/r0"
	AUTHTYPE    = "m.login.dummy"
	nameSuffix	= "@smartraiden"
	LOGINTYPE 	= "m.login.password"
	ALIASFRAGMENT="discovery"
	DISCOVERYROOMALIASFULL	= "#matrix.local.smartraiden0:cy"
	//DISCOVERYROOMSERVER		= "transport01.smartraiden.network"
	DISCOVERYROOMSERVER		= "cy"//test
	CHATPRESET				= "public_chat"
)

var(
	ValidUserIDRegex = regexp.MustCompile(`^@(0x[0-9a-f]{40})(?:\.[0-9a-f]{8})?(?::.+)?$`)//(`^[0-9a-z_\-./]+$`)
	NETWORKNAME="mainnet"
	SERVERNAME	= params.DeFaultMatrixServerName
)

var mobilefeature=[]string{}

func InitMatrixTransport(name,matrixHomeServerURL string,key *ecdsa.PrivateKey,devicetype string)(*MatrixTransport,error) {
	mtr := &MatrixTransport{
		servername:        SERVERNAME,
		running:           false,
		stopreceiving:     false,
		key:               key,
		Users:             make(map[string]*User),
		UseridToPresence:  make(map[string]string),
		AddressToPresence: make(map[common.Address]bool),
		UseDeviceType:     devicetype,
	}
	mcli, err := matrixcomm.NewClient(matrixHomeServerURL, "", "", PATHPREFIX0)
	if err != nil {
		return nil, err
	}
	_, errchk := mcli.Versions()
	if errchk != nil {
		log.Error(fmt.Sprintf("matrix communication failed,cannot connect to server %s", SERVERNAME))
		return nil, nil
	}
	mtr.matrixcli = mcli
	return mtr, nil
}

func validate_userid_signature(user User) (address []byte,err error) {
	//displayname should be an address in the self._userid_re format
	err=fmt.Errorf("validate %s failed");
	_match := ValidUserIDRegex.MatchString(user.user_id)
	if _match == false {
		return
	}
	_address, err := matrixcomm.ExtractUserLocalpart(user.user_id) //"@myname:smartraiden.org:cy"->"myname"
	if err != nil {
		return
	}
	var addrmuti=regexp.MustCompile(`^(0x[0-9a-f]{40})`)
	addrlocal:=addrmuti.FindString(_address)
	if addrlocal==""{
		return
	}
	if len(_address)!=51 || len(user.display_name)!=132{
		return
	}
	if _,err0:=hexutil.Decode(addrlocal);err0!=nil{
		return
	}
	if _,err0:=hexutil.Decode(user.display_name);err0!=nil{
		return
	}
	address = hexutil.MustDecode(addrlocal)
	useridtmp:=utils.Sha3([]byte(user.user_id))//userid 格式:  @0x....:xx
	displaynametmp:=hexutil.MustDecode(user.display_name)//去掉0x转byte[]
	recovered,err:= _recover(useridtmp[:], displaynametmp) //或者临时读取服务器上的GetDisplayName（）
	if err!=nil{
		return
	}
	if !bytes.Equal(recovered, address) {
		address = nil
		err = fmt.Errorf("validate %s failed", user.user_id)
		return
	}
	err=nil
	return
}

func Int32ToBytes(i int32) []byte {
	var buf= make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(i))
	return buf
}

func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

func _recover(data ,signature []byte)(address []byte,err error) {
	recoverPub, err := crypto.Ecrecover(data, signature)
	if err != nil {
		return
	}
	//address = crypto.Keccak256(recoverPub)[1:][12:]
	address = utils.PubkeyToAddress(recoverPub).Bytes()

	return
}

func ChecksumAddress (address string) string {
	address = strings.Replace(strings.ToLower(address), "0x", "", 1)
	addressHash := hex.EncodeToString(crypto.Keccak256([]byte(address)))
	checksumAddress := "0x"
	for i := 0; i < len(address); i++ {
		l, _ := strconv.ParseInt(string(addressHash[i]), 16, 16)
		if (l > 7) {
			checksumAddress += strings.ToUpper(string(address[i]))
		} else {
			checksumAddress += string(address[i])
		}
	}
	return checksumAddress
}