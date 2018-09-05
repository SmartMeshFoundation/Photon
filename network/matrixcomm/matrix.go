package matrixcomm

import (
	"net/url"
	"net/http"
	"fmt"
	"time"
	"encoding/json"
	"bytes"
	"io/ioutil"
	"strconv"
	"io"
	"sync"
	"path"
	"strings"
	"net"
)

const (
	ONLINE = "online"
	//"private_chat", "public_chat", "trusted_private_chat"
)

var MatrixHttpClient = &http.Client{
	Transport: &http.Transport{
		Dial: func(netw, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(netw, addr, time.Second*2)
			if err != nil {
				fmt.Println("dail timeout", err)
				return nil, err
			}
			return c, nil
		},
		MaxIdleConnsPerHost:   2,
		ResponseHeaderTimeout: time.Second * 5,
	},
}

type MatrixClient struct {
	HomeserverURL    *url.URL
	Prefix           string
	UserID           string
	AccessToken      string
	Client           *http.Client
	Store            Storer //store rooms/tokens/ids
	Syncer           Syncer //process /sync responses
	AppServiceUserID string
	syncingMutex     sync.Mutex
	syncingID        uint32
}

//get user's track info
func (mcli *MatrixClient) GetWhois(userid string)(xagent interface{},err error){
	urlPath := mcli.BuildURL( "admin", "whois",userid)
	_, err = mcli.MakeRequest("GET", urlPath, nil, &xagent)
	return
}

//direct get nodes status
func (mcli *MatrixClient) NodesStatus(address[] string) (map[string]string,error) {
	numaddrs := len(address)
	if numaddrs > 10 {
		return nil, nil
	}
	var askwho []string
	for _, dwho := range address {
		userID := "@" + dwho + ":cy"
		askwho = append(askwho, userID)
	}
	err0 := mcli.PostPresenceList(&ReqPresenceList{
		Drop: askwho,
	})
	if err0 != nil {
	}
	time.Sleep(time.Millisecond * 5)
	//DROP all
	err1 := mcli.PostPresenceList(&ReqPresenceList{
		Invite: askwho,
	})
	if err1 != nil {
		//return nil, err1
	}
	respl, err := mcli.GetPresenceList(mcli.UserID)
	if err != nil {
		return nil, err
	}
	rx := make(map[string]string)
	for i, x := range respl {
		user_id := x.UserID
		user_id = strings.Split(user_id, "@")[1]
		user_id = strings.Split(user_id, ":")[0]
		rtnpresence := x.Presence
		rtnaddr := user_id
		rx[rtnaddr] = rtnpresence
		println(i, ":", user_id, "---", rtnpresence)
	}
	return rx, nil
}

func (mcli *MatrixClient) Sync() error {
	syncingID := mcli.incrementSyncingID()
	nextBatch := mcli.Store.LoadNextBatch(mcli.UserID)
	filterID := mcli.Store.LoadFilterID(mcli.UserID)
	if filterID == "" {
		filterJSON := mcli.Syncer.GetFilterJSON(mcli.UserID)
		resFilter, err := mcli.CreateFilter(filterJSON)
		if err != nil {
			return err
		}
		filterID = resFilter.FilterID
		mcli.Store.SaveFilterID(mcli.UserID, filterID)
	}
	for {
		resSync, err := mcli.SyncRequest(2000, nextBatch, filterID, false, "online")
		if err != nil {
			duration, err2 := mcli.Syncer.OnFailedSync(resSync, err)
			if err2 != nil {
				return err2
			}
			time.Sleep(duration)
			continue
		}
		//fmt.Println(syncingID,"\t",filterID)
		if mcli.getSyncingID() != syncingID {
			return nil
		}
		mcli.Store.SaveNextBatch(mcli.UserID, resSync.NextBatch)
		if err = mcli.Syncer.ProcessResponse(resSync, nextBatch); err != nil {
			return err
		}
		nextBatch = resSync.NextBatch
		/*errpu := mcli.SetPresenceState(&ReqPresenceUser{
			Presence: ONLINE,
		})
		if errpu != nil {
			return errpu
		}*/
	}
}

//as allDoneMRE=gcnew ManualResetEvent(false);allDoneMRE->Reset();allDoneMRE->WaitOne();allDoneMRE->Set();
func (mcli *MatrixClient) incrementSyncingID() uint32 {
	mcli.syncingMutex.Lock()
	defer mcli.syncingMutex.Unlock()
	mcli.syncingID++
	return mcli.syncingID
}

//get syncingID
func (mcli *MatrixClient) getSyncingID() uint32 {
	mcli.syncingMutex.Lock()
	defer mcli.syncingMutex.Unlock()
	return mcli.syncingID
}

func (mcli *MatrixClient) StopSync() {
	mcli.incrementSyncingID()
}

//Presence/list/{userId} invite!=remove）
func (mcli *MatrixClient) PostPresenceList(req *ReqPresenceList)(err error) {
	urlPath := mcli.BuildURL( "presence", "list", mcli.UserID)
	_, err = mcli.MakeRequest("POST", urlPath, req, nil)
	return
}

//Presence/list/{userId}
func (mcli *MatrixClient) GetPresenceList(userid string)(resp[] *RespPresenceList,err error) {
	urlPath := mcli.BuildURL( "presence", "list", mcli.UserID)
	_, err = mcli.MakeRequest("GET", urlPath, nil, &resp)
	return
}

//Presence/{userId}/status
func (mcli *MatrixClient) SetPresenceState(req *ReqPresenceUser)(err error) {
	urlPath := mcli.BuildURL( "presence", mcli.UserID, "status")
	_, err = mcli.MakeRequest("PUT", urlPath, req, nil)
	return
}

//Presence/{userId}/status
func (cli *MatrixClient) GetPresenceState(userid string)(resp *RespPresenceUser,err error) {
	urlPath := cli.BuildURL( "presence", userid, "status")
	_, err = cli.MakeRequest("GET", urlPath, nil, &resp)
	return
}

//user/{userId}/account_data/{type}
func (mcli *MatrixClient) SetAccountData(userid,xtype string,req *ReqAccountData)(err error) {
	urlPath := mcli.BuildURL("user", userid, "account_data", xtype) //"network0.smatrraiden.rooms"
	_, err = mcli.MakeRequest("PUT", urlPath, req, nil)
	return
}

func (mcli *MatrixClient) BuildURL(urlPath ...string) string {
	ps := []string{mcli.Prefix}
	for _, p := range urlPath {
		ps = append(ps, p)
	}
	return mcli.BuildBaseURL(ps...)
}

func (mcli *MatrixClient) BuildBaseURL(urlPath ...string) string {
	hsURL, _ := url.Parse(mcli.HomeserverURL.String())
	parts := []string{hsURL.Path}
	parts = append(parts, urlPath...)
	hsURL.Path = path.Join(parts...)
	query := hsURL.Query()
	if mcli.AccessToken != "" {
		query.Set("access_token", mcli.AccessToken)
	}
	if mcli.AppServiceUserID != "" {
		query.Set("user_id", mcli.AppServiceUserID)
	}
	hsURL.RawQuery = query.Encode()
	return hsURL.String()
}

func (mcli *MatrixClient) BuildURLWithQuery(urlPath []string, urlQuery map[string]string) string {
	u, _ := url.Parse(mcli.BuildURL(urlPath...))
	q := u.Query()
	for k, v := range urlQuery {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func (mcli *MatrixClient) SetCredentials(userID, accessToken string) {
	mcli.UserID = userID
	mcli.AccessToken = accessToken
}

func (mcli *MatrixClient) ClearCredentials() {
	mcli.AccessToken = ""
	mcli.UserID = ""
}

//MarshalByRefObject] [Serializable]
func (mcli *MatrixClient) MakeRequest(method string, httpURL string, reqBody interface{}, resBody interface{}) ([]byte, error) {
	var req *http.Request
	var err error
	if reqBody != nil {
		var jsonStr []byte
		jsonStr, err = json.Marshal(reqBody)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest(method, httpURL, bytes.NewBuffer(jsonStr))
	} else {
		req, err = http.NewRequest(method, httpURL, nil)
	}

	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := mcli.Client.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	contents, err := ioutil.ReadAll(res.Body)
	if res.StatusCode!=http.StatusOK{
	//if res.StatusCode/100 != 2 { // not 2xx
		var wrap error
		var respErr RespError
		if _ = json.Unmarshal(contents, &respErr); respErr.ErrCode != "" {
			wrap = respErr
		}
		msg := "Failed to " + method + " JSON to " + req.URL.Path
		if wrap == nil {
			msg = msg + ": " + string(contents)
		}
		return contents, HTTPError{
			Code:         res.StatusCode,
			Message:      msg,
			WrappedError: wrap,
		}
	}
	if err != nil {
		return nil, err
	}
	if resBody != nil {
		if err = json.Unmarshal(contents, &resBody); err != nil {
			return nil, err
		}
	}
	return contents, nil
}

func (mcli *MatrixClient) CreateFilter(filter json.RawMessage) (resp *RespCreateFilter, err error) {
	urlPath := mcli.BuildURL("user", mcli.UserID, "filter")
	_, err = mcli.MakeRequest("POST", urlPath, &filter, &resp)
	return
}

func (mcli *MatrixClient) SyncRequest(timeout int, since, filterID string, fullState bool, setPresence string) (resp *RespSync, err error) {
	query := map[string]string{
		"timeout": strconv.Itoa(timeout),
	}
	if since != "" {
		query["since"] = since
	}
	if filterID != "" {
		query["filter"] = filterID
	}
	if setPresence != "" {
		query["set_presence"] = setPresence
	}
	if fullState {
		query["full_state"] = "true"
	}
	urlPath := mcli.BuildURLWithQuery([]string{"sync"}, query)
	_, err = mcli.MakeRequest("GET", urlPath, nil, &resp)
	return
}

func (mcli *MatrixClient) register(u string, req *ReqRegister) (resp *RespRegister, uiaResp *RespUserInteractive, err error) {
	var bodyBytes []byte
	bodyBytes, err = mcli.MakeRequest("POST", u, req, nil)
	if err != nil {
		httpErr, ok := err.(HTTPError)
		if !ok {
			return
		}
		if httpErr.Code == 401 {
			err = json.Unmarshal(bodyBytes, &uiaResp)
			return
		}
		return
	}
	err = json.Unmarshal(bodyBytes, &resp)
	return
}

func (mcli *MatrixClient) Register(req *ReqRegister) (*RespRegister, *RespUserInteractive, error) {
	u := mcli.BuildURL("register")
	return mcli.register(u, req)
}

func (mcli *MatrixClient) RegisterGuest(req *ReqRegister) (*RespRegister, *RespUserInteractive, error) {
	query := map[string]string{
		"kind": "guest",
	}
	u := mcli.BuildURLWithQuery([]string{"register"}, query)
	return mcli.register(u, req)
}

func (mcli *MatrixClient) Login(req *ReqLogin) (resp *RespLogin, err error) {
	urlPath := mcli.BuildURL("login")
	_, err = mcli.MakeRequest("POST", urlPath, req, &resp)
	return
}

func (mcli *MatrixClient) Logout() (resp *RespLogout, err error) {
	urlPath := mcli.BuildURL("logout")
	_, err = mcli.MakeRequest("POST", urlPath, nil, &resp)
	return
}

func (mcli *MatrixClient) GetDisplayName(mxid string) (resp *RespUserDisplayName, err error) {
	urlPath := mcli.BuildURL("profile", mxid, "displayname")
	_, err = mcli.MakeRequest("GET", urlPath, nil, &resp)
	return
}

func (mcli *MatrixClient) GetOwnDisplayName() (resp *RespUserDisplayName, err error) {
	urlPath := mcli.BuildURL("profile", mcli.UserID, "displayname")
	_, err = mcli.MakeRequest("GET", urlPath, nil, &resp)
	return
}

func (mcli *MatrixClient) SetDisplayName(displayName string) (err error) {
	urlPath := mcli.BuildURL("profile", mcli.UserID, "displayname")
	s := struct {
		DisplayName string `json:"displayname"`
	}{displayName}
	_, err = mcli.MakeRequest("PUT", urlPath, &s, nil)
	return
}

func (mcli *MatrixClient) GetAvatarURL() (url string, err error) {
	urlPath := mcli.BuildURL("profile", mcli.UserID, "avatar_url")
	s := struct {
		AvatarURL string `json:"avatar_url"`
	}{}

	_, err = mcli.MakeRequest("GET", urlPath, nil, &s)
	if err != nil {
		return "", err
	}

	return s.AvatarURL, nil
}

func (mcli *MatrixClient) SetAvatarURL(url string) (err error) {
	urlPath := mcli.BuildURL("profile", mcli.UserID, "avatar_url")
	s := struct {
		AvatarURL string `json:"avatar_url"`
	}{url}
	_, err = mcli.MakeRequest("PUT", urlPath, &s, nil)
	if err != nil {
		return err
	}

	return nil
}

func (mcli *MatrixClient) CreateRoom(req *ReqCreateRoom) (resp *RespCreateRoom, err error) {
	urlPath := mcli.BuildURL("createRoom")
	_, err = mcli.MakeRequest("POST", urlPath, req, &resp)
	return
}

func (mcli *MatrixClient) JoinRoom(roomIDorAlias, serverName string, content interface{}) (resp *RespJoinRoom, err error) {
	var urlPath string
	if serverName != "" {
		urlPath = mcli.BuildURLWithQuery([]string{"join", roomIDorAlias}, map[string]string{
			"server_name": serverName,
		})
	} else {
		urlPath = mcli.BuildURL("join", roomIDorAlias)
	}
	_, err = mcli.MakeRequest("POST", urlPath, content, &resp)
	return
}

func (mcli *MatrixClient) LeaveRoom(roomID string) (resp *RespLeaveRoom, err error) {
	u := mcli.BuildURL("rooms", roomID, "leave")
	_, err = mcli.MakeRequest("POST", u, struct{}{}, &resp)
	return
}

func (mcli *MatrixClient) SendMessageEvent(roomID string, eventType string, contentJSON interface{}) (resp *RespSendEvent, err error) {
	txnID := txnID()
	urlPath := mcli.BuildURL("rooms", roomID, "send", eventType, txnID)
	_, err = mcli.MakeRequest("PUT", urlPath, contentJSON, &resp)
	//
	return
}

func (mcli *MatrixClient) SendStateEvent(roomID, eventType, stateKey string, contentJSON interface{}) (resp *RespSendEvent, err error) {
	urlPath := mcli.BuildURL("rooms", roomID, "state", eventType, stateKey)
	_, err = mcli.MakeRequest("PUT", urlPath, contentJSON, &resp)
	return
}

func (mcli *MatrixClient) SendText(roomID, text string) (*RespSendEvent, error) {
	return mcli.SendMessageEvent(roomID, "m.room.message",
		TextMessage{"m.text", text})
}

func (mcli *MatrixClient) SendImage(roomID, body, url string) (*RespSendEvent, error) {
	return mcli.SendMessageEvent(roomID, "m.room.message",
		ImageMessage{
			MsgType: "m.image",
			Body:    body,
			URL:     url,
		})
}

func (mcli *MatrixClient) SendVideo(roomID, body, url string) (*RespSendEvent, error) {
	return mcli.SendMessageEvent(roomID, "m.room.message",
		VideoMessage{
			MsgType: "m.video",
			Body:    body,
			URL:     url,
		})
}

func (mcli *MatrixClient) SendNotice(roomID, text string) (*RespSendEvent, error) {
	return mcli.SendMessageEvent(roomID, "m.room.message",
		TextMessage{"m.notice", text})
}

func (mcli *MatrixClient) RedactEvent(roomID, eventID string, req *ReqRedact) (resp *RespSendEvent, err error) {
	txnID := txnID()
	urlPath := mcli.BuildURL("rooms", roomID, "redact", eventID, txnID)
	_, err = mcli.MakeRequest("PUT", urlPath, req, &resp)
	return
}

func (mcli *MatrixClient) ForgetRoom(roomID string) (resp *RespForgetRoom, err error) {
	u := mcli.BuildURL("rooms", roomID, "forget")
	_, err = mcli.MakeRequest("POST", u, struct{}{}, &resp)
	return
}

func (mcli *MatrixClient) InviteUser(roomID string, req *ReqInviteUser) (resp *RespInviteUser, err error) {
	u := mcli.BuildURL("rooms", roomID, "invite")
	_, err = mcli.MakeRequest("POST", u, req, &resp)
	return
}

func (mcli *MatrixClient) InviteUserByThirdParty(roomID string, req *ReqInvite3PID) (resp *RespInviteUser, err error) {
	u := mcli.BuildURL("rooms", roomID, "invite")
	_, err = mcli.MakeRequest("POST", u, req, &resp)
	return
}

func (mcli *MatrixClient) KickUser(roomID string, req *ReqKickUser) (resp *RespKickUser, err error) {
	u := mcli.BuildURL("rooms", roomID, "kick")
	_, err = mcli.MakeRequest("POST", u, req, &resp)
	return
}

func (mcli *MatrixClient) BanUser(roomID string, req *ReqBanUser) (resp *RespBanUser, err error) {
	u := mcli.BuildURL("rooms", roomID, "ban")
	_, err = mcli.MakeRequest("POST", u, req, &resp)
	return
}

func (mcli *MatrixClient) UnbanUser(roomID string, req *ReqUnbanUser) (resp *RespUnbanUser, err error) {
	u := mcli.BuildURL("rooms", roomID, "unban")
	_, err = mcli.MakeRequest("POST", u, req, &resp)
	return
}

func (mcli *MatrixClient) UserTyping(roomID string, typing bool, timeout int64) (resp *RespTyping, err error) {
	req := ReqTyping{Typing: typing, Timeout: timeout}
	u := mcli.BuildURL("rooms", roomID, "typing", mcli.UserID)
	_, err = mcli.MakeRequest("PUT", u, req, &resp)
	return
}

func (mcli *MatrixClient) StateEvent(roomID, eventType, stateKey string, outContent interface{}) (err error) {
	u := mcli.BuildURL("rooms", roomID, "state", eventType, stateKey)
	_, err = mcli.MakeRequest("GET", u, nil, outContent)
	return
}

func (mcli *MatrixClient) UploadLink(link string) (*RespMediaUpload, error) {
	res, err := mcli.Client.Get(link)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	return mcli.UploadToContentRepo(res.Body, res.Header.Get("Content-Type"), res.ContentLength)
}

func (mcli *MatrixClient) UploadToContentRepo(content io.Reader, contentType string, contentLength int64) (*RespMediaUpload, error) {
	req, err := http.NewRequest("POST", mcli.BuildBaseURL("_matrix/media/r0/upload"), content)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	req.ContentLength = contentLength
	res, err := mcli.Client.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		contents, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, HTTPError{
				Message: "Upload request failed - Failed to read response body: " + err.Error(),
				Code:    res.StatusCode,
			}
		}
		return nil, HTTPError{
			Message: "Upload request failed: " + string(contents),
			Code:    res.StatusCode,
		}
	}
	var m RespMediaUpload
	if err := json.NewDecoder(res.Body).Decode(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (mcli *MatrixClient) JoinedMembers(roomID string) (resp *RespJoinedMembers, err error) {
	u := mcli.BuildURL("rooms", roomID, "joined_members")
	_, err = mcli.MakeRequest("GET", u, nil, &resp)
	return
}

func (mcli *MatrixClient) JoinedRooms() (resp *RespJoinedRooms, err error) {
	u := mcli.BuildURL("joined_rooms")
	_, err = mcli.MakeRequest("GET", u, nil, &resp)
	return
}

//room's history-info-page
func (mcli *MatrixClient) Messages(roomID, from, to string, dir rune, limit int) (resp *RespMessages, err error) {
	query := map[string]string{
		"from": from,
		"dir":  string(dir),
	}
	if to != "" {
		query["to"] = to
	}
	if limit != 0 {
		query["limit"] = strconv.Itoa(limit)
	}

	urlPath := mcli.BuildURLWithQuery([]string{"rooms", roomID, "messages"}, query)
	_, err = mcli.MakeRequest("GET", urlPath, nil, &resp)
	return
}

func (mcli *MatrixClient) TurnServer() (resp *RespTurnServer, err error) {
	urlPath := mcli.BuildURL("voip", "turnServer")
	_, err = mcli.MakeRequest("GET", urlPath, nil, &resp)
	return
}

func (mcli *MatrixClient) Versions() (resp *RespVersions, err error) {
	urlPath := mcli.BuildBaseURL("_matrix", "client", "versions")
	_, err = mcli.MakeRequest("GET", urlPath, nil, &resp)
	return
}


func NewClient(homeserverURL, userID, accessToken,pathPrefix string) (*MatrixClient, error) {
	hsURL, err := url.Parse(homeserverURL)
	if err != nil {
		return nil, err
	}
	cli := MatrixClient{
		AccessToken:   accessToken,
		HomeserverURL: hsURL,
		UserID:        userID,
		Prefix:        pathPrefix,
	}
	cli.Client = MatrixHttpClient
	return &cli, nil
}

type HTTPError struct {
	WrappedError error
	Message      string
	Code         int
}

func (e HTTPError) Error() string {
	var wrappedErrMsg string
	if e.WrappedError != nil {
		wrappedErrMsg = e.WrappedError.Error()
	}
	return fmt.Sprintf("msg=%s code=%d wrapped=%s", e.Message, e.Code, wrappedErrMsg)
}

func txnID() string {
	return "go" + strconv.FormatInt(time.Now().UnixNano(), 10)
}




