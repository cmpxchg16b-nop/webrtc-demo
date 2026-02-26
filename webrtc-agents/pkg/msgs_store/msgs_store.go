package msgs_store

import (
	"errors"
	"sync/atomic"

	"github.com/google/uuid"
)

// Index by sender id
type IndexedMsgsCollection struct {
	store map[string][]interface{}
}

func NewIndexedMsgsCollection() *IndexedMsgsCollection {
	return &IndexedMsgsCollection{
		store: make(map[string][]interface{}),
	}
}

type IdentifiableMessage interface {
	GetSenderId() string
}

func (indexColl *IndexedMsgsCollection) DeepClone() MsgsCollection {
	newMap := make(map[string][]interface{})
	for senderId, li := range indexColl.store {
		newList := make([]interface{}, len(li))
		copy(newList, li)
		newMap[senderId] = newList
	}
	newIndexColl := new(IndexedMsgsCollection)
	newIndexColl.store = newMap
	return newIndexColl
}

func (indexColl *IndexedMsgsCollection) Append(msg interface{}) {
	senderId := ""
	if sender, ok := msg.(IdentifiableMessage); ok {
		senderId = sender.GetSenderId()
	}
	indexColl.store[senderId] = append(indexColl.store[senderId], msg)
}

// GetMessagesBySenderId returns messages for a specific sender
func (indexColl *IndexedMsgsCollection) GetMessagesBySenderId(senderId string) []interface{} {
	return indexColl.store[senderId]
}

type MsgsCollection interface {
	DeepClone() MsgsCollection
	Append(msg interface{})
}

type MsgsStore struct {
	ResourceVersion string
	collection      MsgsCollection
}

func NewMsgsStore(backend MsgsCollection) *MsgsStore {
	return &MsgsStore{
		ResourceVersion: uuid.NewString(),
		collection:      backend,
	}
}

func (store *MsgsStore) DeepClone() *MsgsStore {
	newMsgsStore := new(MsgsStore)
	newMsgsStore.ResourceVersion = uuid.NewString()
	newMsgsStore.collection = store.collection.DeepClone()
	return newMsgsStore
}

func (store *MsgsStore) Append(msg interface{}) {
	store.collection.Append(msg)
}

func (store *MsgsStore) Load() MsgsCollection {
	return store.collection
}

type BackendFactory func() MsgsCollection

type SyncMsgsStore struct {
	ptr                   atomic.Pointer[MsgsStore]
	defaultBackendFactory BackendFactory
}

func NewSyncMsgsStore(backendFactory BackendFactory) *SyncMsgsStore {
	return &SyncMsgsStore{
		defaultBackendFactory: backendFactory,
	}
}

func (syncMsgsStore *SyncMsgsStore) Append(msg interface{}) error {
	maxRetries := 10000
	for maxRetries > 0 {
		old := syncMsgsStore.ptr.Load()
		var updated *MsgsStore
		if old == nil {
			updated = NewMsgsStore(syncMsgsStore.defaultBackendFactory())
		} else {
			updated = old.DeepClone()
		}
		updated.Append(msg)

		swapped := syncMsgsStore.ptr.CompareAndSwap(old, updated)
		if swapped {
			break
		}
		maxRetries--
	}

	if maxRetries <= 0 {
		return errors.New("max retries reached, CAS lock timeout")
	}

	return nil
}

func (syncMsgsStore *SyncMsgsStore) Load() *MsgsStore {
	return syncMsgsStore.ptr.Load()
}
