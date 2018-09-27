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
func (peer *MatrixPeer) stop() {
	close(peer.quitChan)
}
func (peer *MatrixPeer) loop() {
	for {
		select {
		case <-peer.quitChan:
			return
		case <-peer.receiveMessage:
			continue
		/*
			dont receive any message in ten minutes,this peer should be removed.
		*/
		case <-time.After(time.Minute * 10):
			peer.removeChan <- peer.address
		}
	}
}

func (peer *MatrixPeer) isValidUserID(userID string) bool {
	for _, u := range peer.candidateUsers {
		if u.UserID == userID {
			return true
		}
	}
	return false
}

func (peer *MatrixPeer) updateUsers(users []*gomatrix.UserInfo) {
	for _, lu := range users {
		peer.candidateUsers[lu.UserID] = lu
	}
	return
}

func (peer *MatrixPeer) addRoom(roomID string) {
	if peer.rooms[roomID] {
		log.Warn(fmt.Sprintf("roomID %s already exists for %s", roomID, utils.APex(peer.address)))
	}
	peer.rooms[roomID] = true
}
func (peer *MatrixPeer) setStatus(userID string, presence string) bool {
	var status peerStatus
	switch presence {
	case ONLINE:
		status = peerStatusOnline
	case OFFLINE:
		status = peerStatusOffline
	case UNAVAILABLE:
		status = peerStatusUnkown
	}
	user := peer.candidateUsers[userID]
	if user == nil {
		log.Error(fmt.Sprintf("peer %s ,userid %s set status %s ,but i don't kown this userid. MatrixPeer=%s",
			utils.APex2(peer.address), userID, status, utils.StringInterface(peer, 7)))
		return false
	}
	peer.candidateUsersStatus[userID] = status
	for _, s := range peer.candidateUsersStatus {
		if s > status {
			return false
		}
	}
	peer.status = status
	return true
}
