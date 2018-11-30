package db

// Store :
type Store interface {
	Set(table string, key interface{}, value interface{}) error
	Remove(table string, key interface{}) error
	Save(v KeyGetter) error

	Get(table string, key interface{}, to interface{}) error
	All(table string, to interface{}) error
	Find(table string, fieldName string, value interface{}, to interface{}) error
	Range(table string, fieldName string, min, max, to interface{}) error
}

// DB :
type DB interface {
	Store
	Begin(writable bool) (TX, error)
	Close() error
}

// TX :
type TX interface {
	Store
	Commit() error
	Rollback() error
}

// KeyGetter :
type KeyGetter interface {
	GetKey() []byte
}
