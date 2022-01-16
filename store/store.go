package store

import (
	"github.com/raftgo/config"
	"github.com/sirupsen/logrus"
)

type Store interface {
	Put(key Entry) error

	Get(key []byte) ([]byte, error)

	Commit(term int32, lastId int64)

	Snapshot() error

	FetchEntries(term int32, lastId int64) ([]Entry, error)
}

type Entry struct {
	Meta       Meta
	Key, Value []byte
}

type Meta struct {
	Term        int32
	CommitIndex int64
	LastIndex   int64
}

type memoryStore struct {
	data       []Entry
	commitData []Entry
	config     *config.Config
}

func NewMemoryStore(config *config.Config) Store {
	return &memoryStore{data: make([]Entry, 0), commitData: make([]Entry, 0)}
}

func (store *memoryStore) Put(entry Entry) error {
	store.data = append(store.data, entry)
	return nil
}
func (store *memoryStore) Get(key []byte) ([]byte, error) {
	for _, d := range store.commitData {
		logrus.Infof("reqkey:%s==key:%s", string(d.Key), string(key))
		if string(d.Key) == string(key) {
			return d.Value, nil
		}
	}
	return nil, nil
}
func (store *memoryStore) Commit(term int32, lastId int64) {
	for _, d := range store.data {
		if d.Meta.LastIndex == lastId && d.Meta.Term == term {
			store.commitData = append(store.commitData, d)
			return
		}
	}

}
func (store *memoryStore) Snapshot() error {
	return nil
}
func (store *memoryStore) FetchEntries(term int32, lastId int64) ([]Entry, error) {
	var enties []Entry
	isAppend := false
	for _, d := range store.data {
		if d.Meta.LastIndex == lastId && d.Meta.Term == term {
			store.commitData = append(store.commitData, d)
			isAppend = true
		}
		if isAppend {
			enties = append(enties, d)
		}
	}
	return enties, nil
}
