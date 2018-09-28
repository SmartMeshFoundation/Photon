package network

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
)

//DefaultTemporaryPeerTimeout is the time when to remove a peer without receiving new message
const DefaultTemporaryPeerTimeout = time.Minute * 5

type temporaryPeerRoomInfo struct {
	roomID          string    //id of this room
	lastMessageTime time.Time //when the latest message received
}

/*
process peers which I don't have channel with them.
*/
type matrixTemporaryPeers struct {
	Address2Room map[common.Address]*temporaryPeerRoomInfo
}

func newMatrixTemporaryPeers() *matrixTemporaryPeers {
	return &matrixTemporaryPeers{
		Address2Room: make(map[common.Address]*temporaryPeerRoomInfo),
	}
}
func (p *matrixTemporaryPeers) addPeer(peerAddress common.Address, roomID string) {
	p.Address2Room[peerAddress] = &temporaryPeerRoomInfo{
		roomID:          roomID,
		lastMessageTime: time.Now(),
	}
}
func (p *matrixTemporaryPeers) removePeer(peer common.Address) {
	delete(p.Address2Room, peer)
}

func (p *matrixTemporaryPeers) getRoomID(peer common.Address) string {
	r := p.Address2Room[peer]
	if r == nil {
		return ""
	}
	if time.Now().Sub(r.lastMessageTime) > DefaultTemporaryPeerTimeout {
		delete(p.Address2Room, peer)
		return ""
	}
	return r.roomID
}
