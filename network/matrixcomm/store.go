package matrixcomm

type InMemoryStore struct {
	Filters   map[string]string
	NextBatch map[string]string
	Rooms     map[string]*Room
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		Filters:   make(map[string]string),
		NextBatch: make(map[string]string),
		Rooms:     make(map[string]*Room),
	}
}

type Storer interface {
	SaveFilterID(userID, filterID string)
	LoadFilterID(userID string) string
	SaveNextBatch(userID, nextBatchToken string)
	LoadNextBatch(userID string) string
	SaveRoom(room *Room)
	LoadRoom(roomID string) *Room
	LoadRoomOfAll() map[string]*Room
}

func (s *InMemoryStore) SaveFilterID(userID, filterID string) {
	s.Filters[userID] = filterID
}

func (s *InMemoryStore) LoadFilterID(userID string) string {
	return s.Filters[userID]
}

func (s *InMemoryStore) SaveNextBatch(userID, nextBatchToken string) {
	s.NextBatch[userID] = nextBatchToken
}

func (s *InMemoryStore) LoadNextBatch(userID string) string {
	return s.NextBatch[userID]
}

func (s *InMemoryStore) SaveRoom(room *Room) {
	s.Rooms[room.ID] = room
}

func (s *InMemoryStore) LoadRoom(roomID string) *Room {
	return s.Rooms[roomID]
}

func (s *InMemoryStore) LoadRoomOfAll() map[string]*Room  {
	return s.Rooms
}

