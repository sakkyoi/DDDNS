package memory

import (
	"errors"
	"github.com/sakkyoi/DDDNS/internal/store"
	"sync"
	"time"
)

type memoryStore struct {
	mu    sync.RWMutex
	store map[string]store.Record
}

func New() store.Store {
	return &memoryStore{
		store: make(map[string]store.Record),
	}
}

func (ms *memoryStore) Register(domain string, sourceIP string, destIP string, ttl time.Duration) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	key := store.MakeKey(domain, sourceIP)

	ms.store[key] = store.Record{
		Domain:   domain,
		SourceIP: sourceIP,
		DestIP:   destIP,
		TTL:      ttl,
	}

	return nil
}

func (ms *memoryStore) Lookup(domain string, sourceIP string) (string, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	key := store.MakeKey(domain, sourceIP)
	if record, exists := ms.store[key]; exists {
		return record.DestIP, nil
	}

	return "", errors.New("no entry found")
}

func (ms *memoryStore) Unregister(domain string, sourceIP string) error {
	delete(ms.store, store.MakeKey(domain, sourceIP))

	return nil
}
