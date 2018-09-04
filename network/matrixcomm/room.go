package matrixcomm

import "fmt"

type Room struct {
	ID    string
	State map[string]map[string]*Event
}

func (room Room) UpdateState(event *Event) {
	_, exists := room.State[event.Type]
	if !exists {
		room.State[event.Type] = make(map[string]*Event)
	}
	room.State[event.Type][*event.StateKey] = event
}

func (room Room) GetStateEvent(eventType string, stateKey string) *Event {
	stateEventMap, _ := room.State[eventType]
	event, _ := stateEventMap[stateKey]
	return event
}

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

func NewRoom(roomID string) *Room {
	return &Room{
		ID:    roomID,
		State: make(map[string]map[string]*Event),
	}
}
