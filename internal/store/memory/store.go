package memory

import (
	"fmt"
	"github.com/sakkyoi/DDDNS/internal/store"
	"net"
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

func (ms *memoryStore) Lookup(domain string, sourceIp string, mask *uint8) ([]string, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	var ips []string
	if mask == nil {
		key := store.MakeKey(domain, sourceIp)
		if record, exists := ms.store[key]; exists {
			ips = append(ips, record.DestIp)
		}

	} else {
		_, cidr, err := net.ParseCIDR(fmt.Sprintf("%s/%d", sourceIp, *mask))
		if err != nil {
			return nil, fmt.Errorf("failed to parse CIDR: %w", err)
		}

		for _, record := range ms.store {
			if record.Domain == domain && cidr.Contains(net.ParseIP(record.SourceIp)) {
				ips = append(ips, record.DestIp)
			}
		}
	}

	if len(ips) == 0 {
		return nil, fmt.Errorf("no entry found for %s:%s", domain, sourceIp)
	}

	return ips, nil
}

func (ms *memoryStore) Unregister(domain string, sourceIp string) error {
	delete(ms.store, store.MakeKey(domain, sourceIp))

	return nil
}
