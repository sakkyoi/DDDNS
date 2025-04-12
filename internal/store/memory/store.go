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

func (ms *memoryStore) Register(domain string, sourceIp string, destIp string, ttl time.Duration) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	key := store.MakeKey(domain, sourceIp)

	ms.store[key] = store.Record{
		Domain:   domain,
		SourceIp: sourceIp,
		DestIp:   destIp,
		TTL:      ttl,
	}

	return nil
}

func (ms *memoryStore) Lookup(domain string, sourceIp string) (string, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	key := store.MakeKey(domain, sourceIp)
	if record, exists := ms.store[key]; exists {
		return record.DestIp, nil
	}

	return "", errors.New("no entry found")
}

func (ms *memoryStore) Unregister(domain string, sourceIp string) error {
	delete(ms.store, store.MakeKey(domain, sourceIp))

	return nil
}
