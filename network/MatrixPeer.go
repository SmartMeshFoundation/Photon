package network

import (
	"fmt"
	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/utils"

	"github.com/SmartMeshFoundation/SmartRaiden/log"

	"github.com/SmartMeshFoundation/SmartRaiden/network/gomatrix"
	"github.com/ethereum/go-ethereum/common"
)

type peerStatus int

const (
	peerStatusUnkown = iota
	peerStatusOffline
	peerStatusOnline
)

func (s peerStatus) String() string {
	switch s {
	case peerStatusUnkown:
		return UNAVAILABLE
	case peerStatusOffline:
		return OFFLINE
	case peerStatusOnline:
		return ONLINE
	}
	return "error status"
}

//MatrixPeer is the  raiden node on matrix server
type MatrixPeer struct {
	address common.Address //需要通信的对象
	//address 对应的所有可能的 User
	candidateUsers       map[string]*gomatrix.UserInfo
	candidateUsersStatus map[string]peerStatus
	//确定对方在线的那个聊天室
	defaultMessageRoomID string
	rooms                map[string]bool //roomID exists?
	status               peerStatus
	deviceType           string
	hasChannelWith       bool
	removeChan           chan<- common.Address
	quitChan             chan struct{}
	receiveMessage       chan struct{}
}

//NewMatrixPeer create matrix user
func NewMatrixPeer(address common.Address, hasChannel bool, removeChan chan<- common.Address) *MatrixPeer {
	u := &MatrixPeer{
		address:              address,
		hasChannelWith:       hasChannel,
		rooms:                make(map[string]bool),
		candidateUsers:       make(map[string]*gomatrix.UserInfo),
		candidateUsersStatus: make(map[string]peerStatus),
		removeChan:           removeChan,
		quitChan:             make(chan struct{}),
	}
	if !u.hasChannelWith {
		go u.loop()
	}
	return u
}
func (u *MatrixPeer) stop() {
	close(u.quitChan)
}
func (u *MatrixPeer) loop() {
	for {
		select {
		case <-u.quitChan:
			return
		case <-u.receiveMessage:
			continue
		/*
			dont receive any message in ten minutes,this peer should be removed.
		*/
		case <-time.After(time.Minute * 10):
			u.removeChan <- u.address
		}
	}
}

func (u *MatrixPeer) isValidUserID(userID string) bool {
	for _, u := range u.candidateUsers {
		if u.UserID == userID {
			return true
		}
	}
	return false
}

func (u *MatrixPeer) updateUsers(users []*gomatrix.UserInfo) {
	for _, lu := range users {
		u.candidateUsers[lu.UserID] = lu
	}
	return
}

func (u *MatrixPeer) addRoom(roomID string) {
	if u.rooms[roomID] {
		log.Warn(fmt.Sprintf("roomID %s already exists for %s", roomID, utils.APex(u.address)))
	}
	u.rooms[roomID] = true
}
func (u *MatrixPeer) setStatus(userID string, presence string) bool {
	var status peerStatus
	switch presence {
	case ONLINE:
		status = peerStatusOnline
	case OFFLINE:
		status = peerStatusOffline
	case UNAVAILABLE:
		status = peerStatusUnkown
	}
	user := u.candidateUsers[userID]
	if user == nil {
		log.Error(fmt.Sprintf("peer %s ,userid %s set status %s ,but i don't kown this userid",
			utils.APex2(u.address), userID, status))
		return false
	}
	u.candidateUsersStatus[userID] = status
	for _, s := range u.candidateUsersStatus {
		if s > status {
			return false
		}
	}
	u.status = status
	return true
}
