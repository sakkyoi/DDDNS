package store

import (
	"fmt"
	"time"
)

type Store interface {
	Register(domain string, sourceIp string, destIp string, ttl time.Duration) error
	Lookup(domain string, sourceIp string, mask *uint8) ([]string, error)
	Unregister(domain string, sourceIp string) error
}

type Record struct {
	Domain   string
	SourceIp string
	DestIp   string
	TTL      time.Duration
}

func MakeKey(domain string, sourceIp string) string {
	return fmt.Sprintf("dddns:%s:%s", domain, sourceIp)
}
