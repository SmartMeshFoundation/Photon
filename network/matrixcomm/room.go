package matrixcomm

import "fmt"

type Room struct {
	ID    string
	State map[string]map[string]*Event
}

//更新room的状态
func (room Room) UpdateState(event *Event) {
	_, exists := room.State[event.Type]
	if !exists {
		room.State[event.Type] = make(map[string]*Event)
	}
	room.State[event.Type][*event.StateKey] = event
}

//获取最新的ROOMS EVENT
func (room Room) GetStateEvent(eventType string, stateKey string) *Event {
	stateEventMap, _ := room.State[eventType]
	event, _ := stateEventMap[stateKey]
	return event
}

//分出来写，计算XXX
func (room Room) GetMembershipState(userID string) string {
	state := "leave"
	event := room.GetStateEvent("m.room.member", userID)
	if event != nil {
		membershipState, found := event.Content["membership"]
		if found {
			mState, isString := membershipState.(string)
			if isString {
				state = mState
			}
		}
	}
	if state !="leave"{
		fmt.Println(userID,":")
	}
	return state
}

//新增一个room,初始化一个map(room)
func NewRoom(roomID string) *Room {
	return &Room{
		ID:    roomID,
		State: make(map[string]map[string]*Event),
	}
}
