package store

import (
	"fmt"
	"time"
)

type Store interface {
	Register(domain string, sourceIP string, destIP string, ttl time.Duration) error
	Lookup(domain string, sourceIP string) (string, error)
	Unregister(domain string, sourceIP string) error
}

type Record struct {
	Domain   string
	SourceIP string
	DestIP   string
	TTL      time.Duration
}

func MakeKey(domain string, sourceIP string) string {
	return fmt.Sprintf("dddns:%s:%s", domain, sourceIP)
}
