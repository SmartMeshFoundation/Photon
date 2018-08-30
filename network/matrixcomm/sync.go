package matrixcomm

import (
"encoding/json"
"fmt"
"runtime/debug"
"time"
)

type Syncer interface {
	ProcessResponse(resp *RespSync, since string) error
	OnFailedSync(res *RespSync, err error) (time.Duration, error)
	GetFilterJSON(userID string) json.RawMessage
}

type DefaultSyncer struct {
	UserID    string
	Store     Storer
	listeners map[string][]OnEventListener
}

type OnEventListener func(*Event)

func NewDefaultSyncer(userID string, store Storer) *DefaultSyncer {
	return &DefaultSyncer{
		UserID:    userID,
		Store:     store,
		listeners: make(map[string][]OnEventListener),
	}
}

//callback result on listen
func (s *DefaultSyncer) ProcessResponse(res *RespSync, since string) (err error) {
	if !s.shouldProcessResponse(res, since) {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("ProcessResponse panicked! userID=%s since=%s panic=%s\n%s", s.UserID, since, r, debug.Stack())
		}
		fmt.Println("receive 2")
	}()
	for roomID, roomData := range res.Rooms.Join {
		room := s.getOrCreateRoom(roomID)
		for _, event := range roomData.State.Events {
			event.RoomID = roomID
			room.UpdateState(&event)
			s.notifyListeners(&event)
		}
		for _, event := range roomData.Timeline.Events {
			event.RoomID = roomID
			s.notifyListeners(&event)
		}
	}
	for roomID, roomData := range res.Rooms.Invite {
		room := s.getOrCreateRoom(roomID)
		for _, event := range roomData.State.Events {
			event.RoomID = roomID
			room.UpdateState(&event)
			s.notifyListeners(&event)
		}
	}
	for roomID, roomData := range res.Rooms.Leave {
		room := s.getOrCreateRoom(roomID)
		for _, event := range roomData.Timeline.Events {
			if event.StateKey != nil {
				event.RoomID = roomID
				room.UpdateState(&event)
				s.notifyListeners(&event)
			}
		}
	}
	return
}

func (s *DefaultSyncer) OnEventType(eventType string, callback OnEventListener) {
	_, exists := s.listeners[eventType]
	if !exists {
		s.listeners[eventType] = []OnEventListener{}
	}
	s.listeners[eventType] = append(s.listeners[eventType], callback)
}

func (s *DefaultSyncer) shouldProcessResponse(resp *RespSync, since string) bool {
	if since == "" {
		return false
	}

	for roomID, roomData := range resp.Rooms.Join {
		for i := len(roomData.Timeline.Events) - 1; i >= 0; i-- {
			e := roomData.Timeline.Events[i]
			if e.Type == "m.room.member" && e.StateKey != nil && *e.StateKey == s.UserID {
				m := e.Content["membership"]
				mship, ok := m.(string)
				if !ok {
					continue
				}
				if mship == "join" {
					_, ok := resp.Rooms.Join[roomID]
					if !ok {
						continue
					}
					delete(resp.Rooms.Join, roomID)
					delete(resp.Rooms.Invite, roomID)
					break
				}
			}
		}
	}
	return true
}

func (s *DefaultSyncer) getOrCreateRoom(roomID string) *Room {
	room := s.Store.LoadRoom(roomID)
	if room == nil {
		room = NewRoom(roomID)//new Room
		s.Store.SaveRoom(room)
	}
	return room
}

func (s *DefaultSyncer) notifyListeners(event *Event) {
	listeners, exists := s.listeners[event.Type]
	if !exists {
		return
	}
	for _, fn := range listeners {
		fn(event)
	}
}

func (s *DefaultSyncer) OnFailedSync(res *RespSync, err error) (time.Duration, error) {
	return 2 * time.Second, nil
}

func (s *DefaultSyncer) GetFilterJSON(userID string) json.RawMessage {
	return json.RawMessage(`{"room":{"timeline":{"limit":50}}}`)
}
