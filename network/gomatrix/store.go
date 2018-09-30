package gomatrix

// Storer is an interface which must be satisfied to store client data.
type Storer interface {
	SaveFilterID(userID, filterID string)
	LoadFilterID(userID string) string
	SaveNextBatch(userID, nextBatchToken string)
	LoadNextBatch(userID string) string
	SaveRoom(room *Room)
	LoadRoom(roomID string) *Room
	LoadRoomOfAll() map[string]*Room
}

// InMemoryStore implements the Storer interface.
type InMemoryStore struct {
	Filters   map[string]string
	NextBatch map[string]string
	Rooms     map[string]*Room
}

// NewInMemoryStore constructs a new InMemoryStore.
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		Filters:   make(map[string]string),
		NextBatch: make(map[string]string),
		Rooms:     make(map[string]*Room),
	}
}

// SaveFilterID to memory.
func (s *InMemoryStore) SaveFilterID(userID, filterID string) {
	s.Filters[userID] = filterID
}

// LoadFilterID from memory.
func (s *InMemoryStore) LoadFilterID(userID string) string {
	return s.Filters[userID]
}

// SaveNextBatch to memory.
func (s *InMemoryStore) SaveNextBatch(userID, nextBatchToken string) {
	s.NextBatch[userID] = nextBatchToken
}

// LoadNextBatch from memory.
func (s *InMemoryStore) LoadNextBatch(userID string) string {
	return s.NextBatch[userID]
}

// SaveRoom to memory.
func (s *InMemoryStore) SaveRoom(room *Room) {
	s.Rooms[room.ID] = room
}

// LoadRoom from memory.
func (s *InMemoryStore) LoadRoom(roomID string) *Room {
	return s.Rooms[roomID]
}

//LoadRoomOfAll get all rooms from cache memeory
func (s *InMemoryStore) LoadRoomOfAll() map[string]*Room {
	return s.Rooms
}
