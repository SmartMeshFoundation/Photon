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
)

type MatrixTransport struct {
	matrixcli          *matrixcomm.MatrixClient //the instantiated matrix
	servername         string                   //the homeserver's name
	running            bool                     //running status
	stopreceiving      bool                     //Whether to stop accepting(data)
	key                *ecdsa.PrivateKey        //key
	protocol           ProtocolReceiver
	log                log.Logger
	discoveryroomalias string                  //the room's alias of sys pre-configured ("#[RoomNameLocalpart]:[ServerName]")
	discoveryroomid    string                  //the room's ID of sys pre-configured ("![RoomIdData]:[ServerName]")
	Users              map[string]*User        //cache user's base-infos("userID{userID,displayname}")
	OfficialRooms      map[string]string       //all rooms with we knows this service
	UseridToPresence   map[string]string       //cache user's real-time presence by userID("userID{presence}")
	AddressToPresence  map[common.Address]bool //cache user's real-time presence by node's address("userID{presence}")
	UserId             string                  //the current user's ID(@kitty:thisserver)
	UseDeviceType      string
}
/*
------------------------------------------------------------------------------------------------------------------------
*/
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
	_data := utils.BytesToString(data)
	_, err := mtr.matrixcli.SendText(roomID, _data)
	if err != nil {
		mtr.log.Trace(fmt.Sprintf("send to %s, message=%s [succeed]", receiverAddr.String(), _data))
	}
	mtr.log.Trace(fmt.Sprintf("send to %s, message=%s [succeed]", receiverAddr.String(), _data))
}

func (mtr *MatrixTransport) Start() {
	mtr.running = true
	if err := mtr.login_or_register(); err != nil {
		return
	}
	//the node's device type
	//join discovery room
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
	//SaveRoom on other place
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
func (mtr *MatrixTransport) DataHandler(from common.Address, data []byte) {
	mtr.log.Trace(fmt.Sprintf("received from %s, message=%s", utils.APex2(from), encoding.MessageType(data[0])))
	if !mtr.running || mtr.stopreceiving {
		return
	}
	if mtr.protocol != nil {
		mtr.protocol.receive(data)
	}
}
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
	//由于节点可能在多个主服务器上使用帐户，因此从缓存的单个用户存在状态合成复合地址状态。
	//从存在事件更新节点网络的可达性
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
	//状态改变
	mtr.UseridToPresence[userid]=newstate
	mtr.update_address_presence(adderss)
}

func (mtr *MatrixTransport) login_or_register() (_err error) {
	password := hexutil.Encode(mtr._sign([]byte(mtr.servername)))
	regok := false
	loginok:=false
	baseUsername := strings.ToLower(crypto.PubkeyToAddress(mtr.key.PublicKey).String())
	username := baseUsername
	for i := 0; i < 5; i++ {
		if regok==false {
			rand.Seed(time.Now().UnixNano())
			rnd := Int32ToBytes(rand.Int31n(math.MaxInt32))
			username = baseUsername + "." + hex.EncodeToString(rnd)
			username="0xc67f23ce04ca5e8dd9f2e1b5ed4fad877f79267a.59a2bb27"//test data
		}
		//mtr.matrixcli.UserID
		mtr.matrixcli.AccessToken = ""
		//try login
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
				fmt.Println("Could not login. Trying register")
				//register a new account
				authDict := &matrixcomm.AuthDict{
					Type: AUTHTYPE,
				}
				req := &matrixcomm.ReqRegister{
					Auth:     *authDict,
					Username: username,
					Password: password,
					Type:     LOGINTYPE,
					DeviceID: "",
				}
				_, uia, err := mtr.matrixcli.Register(req)
				if err != nil && uia == nil {
					rhttpErr, _ := err.(matrixcomm.HTTPError)
					if rhttpErr.Code == 400 { //M_USER_IN_USE,M_INVALID_USERNAME,M_EXCLUSIVE
						fmt.Println("Username taken,continuing")
						continue
					}
				}
				//register ok
				fmt.Println("register ok,Username=",username,"Password=", password)
				regok = true;
				mtr.matrixcli.UserID = username
				continue
			}
		}else {
			//cache the node's UserID and AccessToken
			mtr.matrixcli.SetCredentials(resplogin.UserID, resplogin.AccessToken)
			mtr.UserId=resplogin.UserID
			loginok=true
			break
		}
	}
	//set displayname as publicly visible(=0x......)
	name:=hexutil.Encode(mtr._sign([]byte(mtr.matrixcli.UserID)))
	if err:=mtr.matrixcli.SetDisplayName(name);err!=nil{

	}
	fmt.Println("the node in matrix named(userID):",mtr.matrixcli.UserID," the node's displayname is:",name)
	if !loginok{
		_err=fmt.Errorf("Could not register or login!")
	}
	return _err
}

func (mtr *MatrixTransport) inventory_rooms() (err error) {
	for _,room:=range mtr.matrixcli.Store.LoadRoomOfAll(){
		if _, ok := mtr.Users[mtr.discoveryroomalias];ok {
			//room.
			fmt.Println("room:",room)
		}
	}
	return nil
}

func (mtr *MatrixTransport) make_room_alias(networkname string) (string) {
	return ROOMPREFIX + ROOMSEP + networkname + ROOMSEP + "discovery"
}

func (mtr *MatrixTransport) _sign(data []byte)(signature []byte) {
	hash := crypto.Keccak256(data)//私钥地址长度32字节，公钥地址长度65
	signature, _ = crypto.Sign(hash[:], mtr.key);
	return
}

func (mtr *MatrixTransport) join_discovery_room() (err error) {
	//定义room_alias
	discovery_room_alias:=mtr.make_room_alias(SERVERNAME)
	//定义格式化的room_alias(#roomalias:servername) e.g:	#smartraiden_cy_discovery:matrix.local.smartraiden e.g:	#smartraiden_networkname_discovery:cy
	discovery_room_alias_full:="#"+discovery_room_alias+":"+SERVERNAME
	//room alias赋值
	mtr.discoveryroomalias=discovery_room_alias_full
	//加入此服务器提供的room(discoveryRoomAliasFull02)
	discovery_room_alias_full=DISCOVERYROOMALIASFULL//测试用
	resp,err:=mtr.matrixcli.JoinRoom(discovery_room_alias_full,mtr.servername,nil)
	httpErr, _ := err.(matrixcomm.HTTPError)
	if httpErr.Code == 500 {
		//Room doesn't exist and create the room(this is the node's resposibility)
		resp, errc := mtr.matrixcli.CreateRoom(&matrixcomm.ReqCreateRoom{
			RoomAliasName: discovery_room_alias,
			Preset:        CHATPRESET,
		}) //创建后自己已经在里面?服务器bug
		if errc != nil {
			err = errc
			fmt.Println(fmt.Errorf("Discovery room s% not found and can't be created on a federated homeserver s%.", discovery_room_alias_full, SERVERNAME))
			return
		}
		mtr.discoveryroomalias=discovery_room_alias
		mtr.discoveryroomid=resp.RoomID
	}
	//room ID赋值
	mtr.discoveryroomid=resp.RoomID
	//填充初始成员,缓存user,参数只能是roomID
	respin, err:= mtr.matrixcli.JoinedMembers(mtr.discoveryroomid)
	if err != nil {
		fmt.Println("The node can't join room ",mtr.discoveryroomalias)
		return
	}
	fmt.Println("room alias:",mtr.discoveryroomalias)
	fmt.Println("room id:",mtr.discoveryroomid)
	for userid,userdata:=range respin.Joined{
		//cache known users
		mtr.get_user(userid,*userdata.DisplayName)
		//邀请到room内(错误)
		usr:=User{
			user_id:      userid,
			display_name: *userdata.DisplayName,
		}
		mtr.maybe_invite_user(usr)
	}
	fmt.Println("join_discovery_room has done")
	return
}

func (mtr *MatrixTransport) maybe_invite_user(user User) (err error) {
	address,err:=validate_userid_signature(user)
	if err!=nil{
		return fmt.Errorf("illegal user id")
	}

	roomid:=mtr.get_room_id_for_address(common.BytesToAddress(address))
	if roomid!=""{
		_,err=mtr.matrixcli.InviteUser(roomid,&matrixcomm.ReqInviteUser{
			UserID:user.user_id,
		})
	}
	return nil
}

func (mtr *MatrixTransport) get_user(userid,displayname string) (user *User,err error) {
	//通过user_id创建User,如果存在，则获取缓存的用户
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
	user =mtr.Users[userid]
	err = nil
	return
}

func (mtr *MatrixTransport) get_room_id_for_address(address common.Address) (roomid string) {
	addressHex:=ChecksumAddress(hexutil.Encode(address[:]))
	roomid=mtr.OfficialRooms[addressHex]
	return
}

func (mtr *MatrixTransport) set_room_id_for_address(address common.Address,roomid string) (err error) {
	addressHex := ChecksumAddress(hexutil.Encode(address[:]))
	address_to_room_id := mtr.OfficialRooms
	if _, ok := address_to_room_id[addressHex]; !ok {
		mtr.OfficialRooms[addressHex] = roomid
	}
	return
}

func (mtr *MatrixTransport) node_health_check(nodeAddress common.Address) (err error) {
	if mtr.running==false{
		return
	}
	fmt.Println(mtr.running)
	nodeAddrHex := hexutil.Encode(nodeAddress.Bytes())
	fmt.Println(nodeAddrHex)
	return
}

func (mtr *MatrixTransport) get_use_devicetype() (rtn string) {
	deviceType := []string{"mobile", "meshbox", "pc", "other"}
	rtn = deviceType[3]
	resp,err:=mtr.matrixcli.GetWhois(mtr.UserId);if err!=nil{
		return
	}
	ismobile:=0
	tmpRtn:=fmt.Sprint(resp)
	for _,x:=range mobilefeature{
		if _index := strings.Index(tmpRtn,x); _index!=-1{
			ismobile++
		}
	}
	if ismobile==0{
		rtn=deviceType[0]
	}
	return
}
/*
------------------------------------------------------------------------------------------------------------------------
*/
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
	DISCOVERYROOMALIASFULL	= "#matrix.local.smartraiden0:cy"
	DISCOVERYROOMSERVERNAME	= "matrix.local.smartraiden0"
	CHATPRESET				= "public_chat"
)

var(
	ValidUserIDRegex = regexp.MustCompile(`^@(0x[0-9a-f]{40})(?:\.[0-9a-f]{8})?(?::.+)?$`)//(`^[0-9a-z_\-./]+$`)
	ID_TO_NETWORKNAME="NXin"
	SERVERNAME	= params.DeFaultMatrixServerName

)

type User struct {
	user_id string
	display_name string
}

var mobilefeature=[]string{}

func InitMatrixTransport(name,matrixHomeServerURL string,key *ecdsa.PrivateKey,devicetype string)(*MatrixTransport,error) {
	mtr := &MatrixTransport{
		servername:        SERVERNAME,
		running:           false,
		stopreceiving:     true,
		key:               key,
		Users:             make(map[string]*User),
		OfficialRooms:     make(map[string]string),
		UseridToPresence:  make(map[string]string),
		AddressToPresence: make(map[common.Address]bool),
		UseDeviceType:     devicetype,
	}
	//baseUsername := strings.ToLower(crypto.PubkeyToAddress(key.PublicKey).String())
	/*baseUsername := hexutil.Encode(crypto.PubkeyToAddress(key.PublicKey).Bytes())
	fmt.Println("username(string):", baseUsername)
	fmt.Println("username(string-->[]byte):", hexutil.MustDecode(baseUsername))
	fmt.Println("username([]byte-->string):", hexutil.Encode(hexutil.MustDecode(baseUsername)))*/
	//there is no access-token and userID equal as string(address)
	mcli, err := matrixcomm.NewClient(matrixHomeServerURL, "", "", PATHPREFIX0)
	if err != nil {
		return nil, err
	}
	//try to check the state of communication
	_, errchk := mcli.Versions()
	if errchk != nil {
		mtr.log.Error(fmt.Sprintf("matrix communication failed,cannot connect to server %s", SERVERNAME))
		return nil,nil
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
	//var addrid=validUsernameRegex.FindString(userid)
	_address, err := matrixcomm.ExtractUserLocalpart(user.user_id) //"@myname:smartraiden.org:cy"->"myname"
	if err != nil {
		return
	}
	//读出地址（点前面的）
	var addrmuti=regexp.MustCompile(`^(0x[0-9a-f]{40})`)
	addrlocal:=addrmuti.FindString(_address)
	if addrlocal==""{//解不出来
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
	recovered,err:= _recover(useridtmp[:], displaynametmp) //或者必须临时读取服务器上的GetDisplayName（）
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
	recoverPub, err := crypto.Ecrecover(data, signature) //从签名中提取公钥
	if err!=nil{
		return
	}
	//address = crypto.Keccak256(recoverPub)[1:][12:]      //账户地址(后20字节)
	address=utils.PubkeyToAddress(recoverPub).Bytes()

	return
}

func ChecksumAddress (address string) string {
	address = strings.Replace(strings.ToLower(address),"0x","",1)
	addressHash := hex.EncodeToString(crypto.Keccak256([]byte(address)))
	checksumAddress := "0x"
	for i := 0; i < len(address); i++ {
		// If ith character is 8 to f then make it uppercase
		l, _ := strconv.ParseInt(string(addressHash[i]), 16, 16)
		if (l > 7) {
			checksumAddress += strings.ToUpper(string(address[i]))
		} else {
			checksumAddress += string(address[i])
		}
	}
	return checksumAddress
}

