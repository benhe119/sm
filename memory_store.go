package sm

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/blevesearch/bleve"
)

type MemoryStore struct {
	sync.RWMutex

	nextid ID
	data   map[ID]*Event
	index  bleve.Index
}

func (store *MemoryStore) Close() error {
	return nil
}

func (store *MemoryStore) NextId() ID {
	store.Lock()
	defer store.Unlock()

	store.nextid++
	return store.nextid
}

func (store *MemoryStore) Save(event *Event) error {
	store.Lock()
	store.data[event.ID] = event
	store.Unlock()

	t := time.Now()
	store.index.Index(event.ID.String(), event)
	metrics.Summary("event", "index").Observe(time.Now().Sub(t).Seconds())

	return nil
}

func (store *MemoryStore) Get(id ID) (event *Event, err error) {
	var ok bool

	store.RLock()
	event, ok = store.data[id]
	store.RUnlock()

	if !ok {
		err = &KeyError{id, ErrNotExist}
	}

	return
}

func (store *MemoryStore) Find(ids ...ID) (events []*Event, err error) {
	store.RLock()
	for _, id := range ids {
		event, ok := store.data[id]
		if ok {
			events = append(events, event)
		}
	}
	store.RUnlock()

	return
}

func (store *MemoryStore) All() (events []*Event, err error) {
	store.RLock()
	for _, event := range store.data {
		events = append(events, event)
	}
	store.RUnlock()

	return
}

func (store *MemoryStore) Search(q string) (events []*Event, err error) {
	size, err := store.index.DocCount()
	if err != nil {
		log.Errorf("error getting index size: %s", err)
		return
	}

	query := bleve.NewQueryStringQuery(q)
	req := bleve.NewSearchRequestOptions(query, int(size), 0, false)
	res, err := store.index.Search(req)
	if err != nil {
		log.Errorf("error performing index search %s: %s", q, err)
		return
	}

	for _, hit := range res.Hits {
		store.RLock()
		event, ok := store.data[ParseId(hit.ID)]
		store.RUnlock()
		if !ok {
			log.Warnf("event #%s missing from store but exists in index!", hit.ID)
			continue
		}
		events = append(events, event)
	}

	return
}

func NewMemoryStore() (Store, error) {
	index, err := bleve.NewMemOnly(bleve.NewIndexMapping())
	if err != nil {
		log.Errorf("error creating index: %s", err)
		return nil, err
	}

	return &MemoryStore{
		data:  make(map[ID]*Event),
		index: index,
	}, nil
}
