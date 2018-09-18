package matrixcomm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
)

// MatrixHTTPClient is a custom http client
var MatrixHTTPClient = &http.Client{
	Transport: &http.Transport{
		Dial: func(netw, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(netw, addr, time.Second*30)
			if err != nil {
				//fmt.Println("dail timeout", err)
				return nil, err
			}
			return c, nil
		},
		MaxIdleConnsPerHost:   100,
		ResponseHeaderTimeout: time.Second * 30,
	},
}

// MatrixClient matrix client
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

// GetWhois get user's track info
func (mcli *MatrixClient) GetWhois(userid string) (xagent interface{}, err error) {
	urlPath := mcli.BuildURL("admin", "whois", userid)
	_, err = mcli.MakeRequest("GET", urlPath, nil, &xagent)
	return
}

// NodesStatus direct get nodes status
func (mcli *MatrixClient) NodesStatus(address []string) (map[string]string, error) {
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
		xuserID := x.UserID
		xuserID = strings.Split(xuserID, "@")[1]
		xuserID = strings.Split(xuserID, ":")[0]
		rtnpresence := x.Presence
		rtnaddr := xuserID
		rx[rtnaddr] = rtnpresence
		println(i, ":", xuserID, "---", rtnpresence)
	}
	return rx, nil
}

// Sync starts syncing with the provided Homeserver.
func (mcli *MatrixClient) Sync() error {
	// Mark the client as syncing.
	// We will keep syncing until the syncing state changes. Either because
	// Sync is called or StopSync is called.
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
		resSync, err := mcli.SyncRequest(20000, nextBatch, filterID, false, "online")
		//fmt.Println(time.Now().Format("2006/1/2 15:04:05"),filterID,  "\t",nextBatch)
		if err != nil {
			duration, err2 := mcli.Syncer.OnFailedSync(resSync, err)
			if err2 != nil {
				return err2
			}
			time.Sleep(duration)
			continue
		}

		// Check that the syncing state hasn't changed
		// Either because we've stopped syncing or another sync has been started.
		// We discard the response from our sync.
		if mcli.getSyncingID() != syncingID {
			return nil
		}

		// Save the token now *before* processing it. This means it's possible
		// to not process some events, but it means that we won't get constantly stuck processing
		// a malformed/buggy event which keeps making us panic.
		mcli.Store.SaveNextBatch(mcli.UserID, resSync.NextBatch)
		if err = mcli.Syncer.ProcessResponse(resSync, nextBatch); err != nil {
			return err
		}
		nextBatch = resSync.NextBatch

		/*//send heartbeat to homeserver and the other participating servers
		errpu := mcli.SetPresenceState(&ReqPresenceUser{
			Presence: ONLINE,
		})
		if errpu != nil {
			return errpu
		}*/
	}
}

//i ncrementSyncingID If Sync is called twice then the first sync will be stopped
func (mcli *MatrixClient) incrementSyncingID() uint32 {
	mcli.syncingMutex.Lock()
	defer mcli.syncingMutex.Unlock()
	mcli.syncingID++
	return mcli.syncingID
}

// getSyncingID get syncingID
func (mcli *MatrixClient) getSyncingID() uint32 {
	mcli.syncingMutex.Lock()
	defer mcli.syncingMutex.Unlock()
	return mcli.syncingID
}

// StopSync stop sync
func (mcli *MatrixClient) StopSync() {
	mcli.incrementSyncingID()
}

// SearchUserDirectory Performs a search for users on the homeserver
func (mcli *MatrixClient) SearchUserDirectory(req *ReqUserSearch) (resp *RespUserSearch, err error) {
	urlPath := mcli.BuildURL("user_directory", "search")
	_, err = mcli.MakeRequest("POST", urlPath, req, &resp)
	return
}

// PostPresenceList Presence/list/{userId} invite!=remove）
func (mcli *MatrixClient) PostPresenceList(req *ReqPresenceList) (err error) {
	urlPath := mcli.BuildURL("presence", "list", mcli.UserID)
	_, err = mcli.MakeRequest("POST", urlPath, req, nil)
	return
}

// GetPresenceList Presence/list/{userId}
func (mcli *MatrixClient) GetPresenceList(userid string) (resp []*RespPresenceList, err error) {
	urlPath := mcli.BuildURL("presence", "list", mcli.UserID)
	_, err = mcli.MakeRequest("GET", urlPath, nil, &resp)
	return
}

// SetPresenceState Presence/{userId}/status
func (mcli *MatrixClient) SetPresenceState(req *ReqPresenceUser) (err error) {
	urlPath := mcli.BuildURL("presence", mcli.UserID, "status")
	_, err = mcli.MakeRequest("PUT", urlPath, req, nil)
	return
}

// GetPresenceState Presence/{userId}/status
func (mcli *MatrixClient) GetPresenceState(userid string) (resp *RespPresenceUser, err error) {
	urlPath := mcli.BuildURL("presence", userid, "status")
	_, err = mcli.MakeRequest("GET", urlPath, nil, &resp)
	return
}

// SetAccountData user/{userId}/account_data/{type}
func (mcli *MatrixClient) SetAccountData(userid, xtype string, addr2room map[string]string) (err error) {
	/*
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
	*/
	/*	query:=map[string]string{

		}

		urlPath:=mcli.BuildURLWithQuery()*/

	urlaccountdata := mcli.BuildURLWithQuery([]string{"user", userid, "account_data,xtype"}, addr2room)
	_, err = mcli.MakeRequest("PUT", urlaccountdata, nil, nil)

	/*urlPath := mcli.BuildURL("user", userid, "account_data", xtype) //"network0.smatrraiden.rooms"
	_, err = mcli.MakeRequest("PUT", urlPath, req, nil)*/
	return
}

// BuildURL builds a URL with the Client's homserver/prefix/access_token set already.
func (mcli *MatrixClient) BuildURL(urlPath ...string) string {
	ps := []string{mcli.Prefix}
	for _, p := range urlPath {
		ps = append(ps, p)
	}
	return mcli.BuildBaseURL(ps...)
}

// BuildBaseURL builds a URL with the Client's homeserver/access_token set already. You must supply the prefix in the path.
func (mcli *MatrixClient) BuildBaseURL(urlPath ...string) string {
	hsURL, err := url.Parse(mcli.HomeserverURL.String())
	if err != nil {
		return ""
	}
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

// BuildURLWithQuery builds a URL with query parameters in addition
func (mcli *MatrixClient) BuildURLWithQuery(urlPath []string, urlQuery map[string]string) string {
	u, err := url.Parse(mcli.BuildURL(urlPath...))
	if err != nil {
		return ""
	}
	q := u.Query()
	for k, v := range urlQuery {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

// SetCredentials sets the user ID and access token on this client instance.
func (mcli *MatrixClient) SetCredentials(userID, accessToken string) {
	mcli.UserID = userID
	mcli.AccessToken = accessToken
}

// ClearCredentials removes the user ID and access token on this client instance.
func (mcli *MatrixClient) ClearCredentials() {
	mcli.AccessToken = ""
	mcli.UserID = ""
}

// MakeRequest makes a JSON HTTP request to the given URL
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
	log.Trace(fmt.Sprintf("matrix url:%s,req:%s", httpURL, reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	// 完成后断开连接
	req.Header.Set("Connection", "close")
	res, err := mcli.Client.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	contents, err := ioutil.ReadAll(res.Body)
	log.Trace(fmt.Sprintf("matrix response err=%s,contents=%s", err, string(contents)))
	if err != nil {
		return nil, err
	}
	if res.StatusCode/100 != 2 {
		var wrap error
		var respErr RespError
		err = json.Unmarshal(contents, &respErr)
		if err != nil {
			return nil, err
		}
		if respErr.ErrCode != "" {
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

// CreateFilter makes an HTTP request with filter
func (mcli *MatrixClient) CreateFilter(filter json.RawMessage) (resp *RespCreateFilter, err error) {
	urlPath := mcli.BuildURL("user", mcli.UserID, "filter")
	_, err = mcli.MakeRequest("POST", urlPath, &filter, &resp)
	return
}

// SyncRequest makes an sync HTTP request
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

// register registe a user on a homeserver
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

// Register Registe a user on a homeserver
func (mcli *MatrixClient) Register(req *ReqRegister) (*RespRegister, *RespUserInteractive, error) {
	u := mcli.BuildURL("register")
	return mcli.register(u, req)
}

// RegisterGuest register as a guest
func (mcli *MatrixClient) RegisterGuest(req *ReqRegister) (*RespRegister, *RespUserInteractive, error) {
	query := map[string]string{
		"kind": "guest",
	}
	u := mcli.BuildURLWithQuery([]string{"register"}, query)
	return mcli.register(u, req)
}

// Login sign in a homeserver
func (mcli *MatrixClient) Login(req *ReqLogin) (resp *RespLogin, err error) {
	urlPath := mcli.BuildURL("login")
	_, err = mcli.MakeRequest("POST", urlPath, req, &resp)
	return
}

// Logout sign out
func (mcli *MatrixClient) Logout() (resp *RespLogout, err error) {
	urlPath := mcli.BuildURL("logout")
	_, err = mcli.MakeRequest("POST", urlPath, nil, &resp)
	return
}

// GetDisplayName get a user's dislayname according to profile
func (mcli *MatrixClient) GetDisplayName(mxid string) (resp *RespUserDisplayName, err error) {
	urlPath := mcli.BuildURL("profile", mxid, "displayname")
	_, err = mcli.MakeRequest("GET", urlPath, nil, &resp)
	return
}

// GetOwnDisplayName get own's displayname
func (mcli *MatrixClient) GetOwnDisplayName() (resp *RespUserDisplayName, err error) {
	urlPath := mcli.BuildURL("profile", mcli.UserID, "displayname")
	_, err = mcli.MakeRequest("GET", urlPath, nil, &resp)
	return
}

// SetDisplayName user's displayname
func (mcli *MatrixClient) SetDisplayName(displayName string) (err error) {
	urlPath := mcli.BuildURL("profile", mcli.UserID, "displayname")
	s := struct {
		DisplayName string `json:"displayname"`
	}{displayName}
	_, err = mcli.MakeRequest("PUT", urlPath, &s, nil)
	return
}

// GetAvatarURL get some one's avatar_url
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

// SetAvatarURL set user's avatar_url
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

// CreateRoom create a room ,retrun  the room's id
func (mcli *MatrixClient) CreateRoom(req *ReqCreateRoom) (resp *RespCreateRoom, err error) {
	urlPath := mcli.BuildURL("createRoom")
	_, err = mcli.MakeRequest("POST", urlPath, req, &resp)
	return
}

// JoinRoom join a room ,the content is user's info
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

// LeaveRoom leave room
func (mcli *MatrixClient) LeaveRoom(roomID string) (resp *RespLeaveRoom, err error) {
	u := mcli.BuildURL("rooms", roomID, "leave")
	_, err = mcli.MakeRequest("POST", u, struct{}{}, &resp)
	return
}

//SendMessageEvent sends a message event into a room
func (mcli *MatrixClient) SendMessageEvent(roomID string, eventType string, contentJSON interface{}) (resp *RespSendEvent, err error) {
	txnID := txnID()
	urlPath := mcli.BuildURL("rooms", roomID, "send", eventType, txnID)
	_, err = mcli.MakeRequest("PUT", urlPath, contentJSON, &resp)
	//
	return
}

// SendStateEvent sends a state event into a room
func (mcli *MatrixClient) SendStateEvent(roomID, eventType, stateKey string, contentJSON interface{}) (resp *RespSendEvent, err error) {
	urlPath := mcli.BuildURL("rooms", roomID, "state", eventType, stateKey)
	_, err = mcli.MakeRequest("PUT", urlPath, contentJSON, &resp)
	return
}

// SendText sends an m.room.message event into the given room with a msgtype of m.text
func (mcli *MatrixClient) SendText(roomID, text string) (*RespSendEvent, error) {
	return mcli.SendMessageEvent(roomID, "m.room.message",
		TextMessage{"m.text", text})
}

// SendNotice sends an m.room.message event into the given room with a msgtype of m.notice
func (mcli *MatrixClient) SendNotice(roomID, text string) (*RespSendEvent, error) {
	return mcli.SendMessageEvent(roomID, "m.room.message",
		TextMessage{"m.notice", text})
}

// SendImage sends an m.room.message event into the given room with a msgtype of m.image
func (mcli *MatrixClient) SendImage(roomID, body, url string) (*RespSendEvent, error) {
	return mcli.SendMessageEvent(roomID, "m.room.message",
		ImageMessage{
			MsgType: "m.image",
			Body:    body,
			URL:     url,
		})
}

// SendVideo sends an m.room.message event into the given room with a msgtype of m.video
func (mcli *MatrixClient) SendVideo(roomID, body, url string) (*RespSendEvent, error) {
	return mcli.SendMessageEvent(roomID, "m.room.message",
		VideoMessage{
			MsgType: "m.video",
			Body:    body,
			URL:     url,
		})
}

// RedactEvent redacts the given event.
func (mcli *MatrixClient) RedactEvent(roomID, eventID string, req *ReqRedact) (resp *RespSendEvent, err error) {
	txnID := txnID()
	urlPath := mcli.BuildURL("rooms", roomID, "redact", eventID, txnID)
	_, err = mcli.MakeRequest("PUT", urlPath, req, &resp)
	return
}

// ForgetRoom forgets a room entirely.
func (mcli *MatrixClient) ForgetRoom(roomID string) (resp *RespForgetRoom, err error) {
	u := mcli.BuildURL("rooms", roomID, "forget")
	_, err = mcli.MakeRequest("POST", u, struct{}{}, &resp)
	return
}

// InviteUser invites a user to a room.
func (mcli *MatrixClient) InviteUser(roomID string, req *ReqInviteUser) (resp *RespInviteUser, err error) {
	u := mcli.BuildURL("rooms", roomID, "invite")
	_, err = mcli.MakeRequest("POST", u, req, &resp)
	return
}

// InviteUserByThirdParty invites a third-party identifier to a room.
func (mcli *MatrixClient) InviteUserByThirdParty(roomID string, req *ReqInvite3PID) (resp *RespInviteUser, err error) {
	u := mcli.BuildURL("rooms", roomID, "invite")
	_, err = mcli.MakeRequest("POST", u, req, &resp)
	return
}

//KickUser kicks a user from a room.
func (mcli *MatrixClient) KickUser(roomID string, req *ReqKickUser) (resp *RespKickUser, err error) {
	u := mcli.BuildURL("rooms", roomID, "kick")
	_, err = mcli.MakeRequest("POST", u, req, &resp)
	return
}

// BanUser bans a user from a room.
func (mcli *MatrixClient) BanUser(roomID string, req *ReqBanUser) (resp *RespBanUser, err error) {
	u := mcli.BuildURL("rooms", roomID, "ban")
	_, err = mcli.MakeRequest("POST", u, req, &resp)
	return
}

// UnbanUser unbans a user from a room.
func (mcli *MatrixClient) UnbanUser(roomID string, req *ReqUnbanUser) (resp *RespUnbanUser, err error) {
	u := mcli.BuildURL("rooms", roomID, "unban")
	_, err = mcli.MakeRequest("POST", u, req, &resp)
	return
}

// UserTyping sets the typing status of the user.
func (mcli *MatrixClient) UserTyping(roomID string, typing bool, timeout int64) (resp *RespTyping, err error) {
	req := ReqTyping{Typing: typing, Timeout: timeout}
	u := mcli.BuildURL("rooms", roomID, "typing", mcli.UserID)
	_, err = mcli.MakeRequest("PUT", u, req, &resp)
	return
}

// StateEvent gets a single state event in a room. It will attempt to JSON unmarshal
// into the given "outContent" struct with the HTTP response body, or return an error.
func (mcli *MatrixClient) StateEvent(roomID, eventType, stateKey string, outContent interface{}) (err error) {
	u := mcli.BuildURL("rooms", roomID, "state", eventType, stateKey)
	_, err = mcli.MakeRequest("GET", u, nil, outContent)
	return
}

// UploadLink uploads an HTTP URL and then returns an MXC URI.
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

// UploadToContentRepo uploads the given bytes to the content repository and returns an MXC URI.
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

// JoinedMembers returns a map of joined room members.
// In general, usage of this API is discouraged in favour of /sync, as calling this API can race with incoming membership changes.
func (mcli *MatrixClient) JoinedMembers(roomID string) (resp *RespJoinedMembers, err error) {
	u := mcli.BuildURL("rooms", roomID, "joined_members")
	_, err = mcli.MakeRequest("GET", u, nil, &resp)
	return
}

// JoinedRooms returns a list of rooms which the client is joined to.
// In general, usage of this API is discouraged in favour of /sync, as calling this API can race with incoming membership changes.
func (mcli *MatrixClient) JoinedRooms() (resp *RespJoinedRooms, err error) {
	u := mcli.BuildURL("joined_rooms")
	_, err = mcli.MakeRequest("GET", u, nil, &resp)
	return
}

// Messages list of message and state events for a room
// It uses pagination query parameters to paginate history in the room.
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

// TurnServer returns turn server details and credentials for the client to use when initiating calls.
func (mcli *MatrixClient) TurnServer() (resp *RespTurnServer, err error) {
	urlPath := mcli.BuildURL("voip", "turnServer")
	_, err = mcli.MakeRequest("GET", urlPath, nil, &resp)
	return
}

//Versions retruns the matrix client-server-api's version
func (mcli *MatrixClient) Versions() (resp *RespVersions, err error) {
	urlPath := mcli.BuildBaseURL("_matrix", "client", "versions")
	_, err = mcli.MakeRequest("GET", urlPath, nil, &resp)
	return
}

// NewClient creates a new Matrix Client ready for syncing
func NewClient(homeserverURL, userID, accessToken, pathPrefix string) (*MatrixClient, error) {
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
	cli.Client = MatrixHTTPClient
	return &cli, nil
}

// HTTPError An HTTP Error response, which may wrap an underlying native Go Error.
type HTTPError struct {
	WrappedError error
	Message      string
	Code         int
}

// Error retrun an error with this fixed format
func (e HTTPError) Error() string {
	var wrappedErrMsg string
	if e.WrappedError != nil {
		wrappedErrMsg = e.WrappedError.Error()
	}
	return fmt.Sprintf("msg=%s code=%d wrapped=%s", e.Message, e.Code, wrappedErrMsg)
}

// txnID sync a request with a different event request
func txnID() string {
	return "go" + strconv.FormatInt(time.Now().UnixNano(), 10)
}
