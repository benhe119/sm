package sm

import (
	"errors"
	"fmt"
	"sync"
)

var (
	ErrNotExist = errors.New("key does not exist")
)

type KeyError struct {
	Key ID
	Err error
}

func (e *KeyError) Error() string {
	return fmt.Sprintf("%s: %d", e.Err, e.Key)
}

type ID uint64

func (id ID) String() string {
	return fmt.Sprintf("%d", id)
}

func ParseId(s string) ID {
	return ID(SafeParseUint64(s, 0))
}

type IdGenerator struct {
	sync.Mutex
	next ID
}

func (id *IdGenerator) Next() ID {
	id.Lock()
	defer id.Unlock()

	id.next++
	return id.next
}

type Store interface {
	Close() error
	NextId() ID
	Save(*Event) error
	Get(id ID) (*Event, error)
	Find(id ...ID) ([]*Event, error)
	All() ([]*Event, error)
	Search(q string) ([]*Event, error)
}
