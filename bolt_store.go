package sm

import (
	"encoding/binary"
	"os"
	"path"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/blevesearch/bleve"
	"github.com/coreos/bbolt"

	"github.com/prologic/sm/codec"
	"github.com/prologic/sm/codec/json"
)

func idToBytes(id ID) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(id))
	return b
}

type BoltStore struct {
	db     *bolt.DB
	nextid *IdGenerator
	index  bleve.Index
	codec  codec.MarshalUnmarshaler
}

func (store *BoltStore) Close() error {
	return store.db.Close()
}

func (store *BoltStore) NextId() ID {
	return ID(0)
}

func (store *BoltStore) Save(event *Event) error {
	err := store.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("events"))
		if err != nil {
			log.Errorf("error creating events bucket: %s", err)
			return err
		}

		if event.ID == ID(0) {
			id, err := b.NextSequence()
			if err != nil {
				log.Errorf("error generating new events sequence: %s", err)
				return err
			}
			event.ID = ID(id)
		}

		buf, err := store.codec.Marshal(event)
		if err != nil {
			log.Errorf("error serializing event: %s", err)
			return err
		}

		key := idToBytes(event.ID)
		return b.Put(key, buf)
	})

	if err != nil {
		log.Errorf("error saving event: %s", err)
		return err
	}

	t := time.Now()
	store.index.Index(event.ID.String(), event)
	metrics.Summary("event", "index").Observe(time.Now().Sub(t).Seconds())

	return nil
}

func (store *BoltStore) Get(id ID) (*Event, error) {
	var event Event

	err := store.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("events"))
		if b == nil {
			return nil
		}

		key := idToBytes(id)
		buf := b.Get(key)
		if buf == nil {
			log.Errorf("event #%d not found", id)
			err := &KeyError{id, ErrNotExist}
			return err
		}

		err := store.codec.Unmarshal(buf, &event)
		if err != nil {
			log.Errorf("error deserializing event #%s: %s", id, err)
			return err
		}

		return nil
	})

	return &event, err
}

func (store *BoltStore) Find(ids ...ID) (events []*Event, err error) {
	err = store.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("events"))
		if b == nil {
			return nil
		}

		for _, id := range ids {
			var event Event
			key := idToBytes(id)
			buf := b.Get(key)
			err := store.codec.Unmarshal(buf, &event)
			if err != nil {
				log.Errorf("error deserializing event #%s: %s", id, err)
				return err
			}

			events = append(events, &event)
		}
		return nil
	})

	return
}

func (store *BoltStore) All() (events []*Event, err error) {
	err = store.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("events"))
		if b == nil {
			return nil
		}

		b.ForEach(func(k, v []byte) error {
			var event Event
			err := store.codec.Unmarshal(v, &event)
			if err != nil {
				log.Errorf("error deserializing events: %s", err)
				return err
			}

			events = append(events, &event)
			return nil
		})
		return nil
	})

	return
}

func (store *BoltStore) Search(q string) (events []*Event, err error) {
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

	var ids []ID
	for _, hit := range res.Hits {
		ids = append(ids, ParseId(hit.ID))
	}

	events, err = store.Find(ids...)

	return
}

func NewBoltStore(dbpath string) (Store, error) {
	db, err := bolt.Open(dbpath, 0644, &bolt.Options{})
	if err != nil {
		log.Errorf("error opening store %s: %s", dbpath, err)
		return nil, err
	}

	var index bleve.Index
	indexpath := path.Join(path.Dir(dbpath), "index.db")
	if _, err = os.Stat(indexpath); err == nil {
		index, err = bleve.Open(indexpath)
	} else {
		index, err = bleve.New(indexpath, bleve.NewIndexMapping())
	}
	if err != nil {
		log.Errorf("error creating index: %s", err)
		return nil, err
	}

	return &BoltStore{
		db:     db,
		nextid: &IdGenerator{},
		index:  index,
		codec:  json.Codec,
	}, nil
}
