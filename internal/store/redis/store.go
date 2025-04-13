package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sakkyoi/DDDNS/internal/store"
	"net"
	"time"
)

type redisStore struct {
	client *redis.Client
}

func New(host string, port int, db int, username string, password string) (store.Store, error) {
	addr := fmt.Sprintf("%s:%d", host, port)

	client := redis.NewClient(&redis.Options{Addr: addr, DB: db, Username: username, Password: password})

	// ping the Redis server to check if it's available
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &redisStore{
		client: client,
	}, nil
}

func (rs *redisStore) Register(domain string, sourceIp string, destIp string, ttl time.Duration) error {
	key := store.MakeKey(domain, sourceIp)

	return rs.client.Set(context.Background(), key, destIp, ttl).Err()
}

func (rs *redisStore) Lookup(domain string, sourceIp string, mask *uint8) ([]string, error) {
	// build the search key
	// TODO redis' scan is not efficient, consider using a sorted set
	var key string
	if mask == nil {
		key = store.MakeKey(domain, sourceIp)
	} else {
		ip, ipNet, err := net.ParseCIDR(fmt.Sprintf("%s/%d", sourceIp, *mask))
		if err != nil {
			return nil, fmt.Errorf("failed to parse CIDR: %w", err)
		}

		ipMask := ip.Mask(ipNet.Mask)

		ones, _ := ipNet.Mask.Size()
		byteCount := ones / 8

		// build the prefix for the scan
		prefix := ""
		for i := 0; i < byteCount; i++ {
			if i > 0 {
				prefix += "."
			}
			prefix += fmt.Sprintf("%d", ipMask[i])
		}

		key = store.MakeKey(domain, fmt.Sprintf("%s.*", prefix))
	}

	// scan the keys
	var ips []string

	iter := rs.client.Scan(context.Background(), 0, key, 1000).Iterator()
	for iter.Next(context.Background()) {
		key := iter.Val()
		val, err := rs.client.Get(context.Background(), key).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to get value for key %s: %w", key, err)
		}

		ips = append(ips, val)
	}

	if len(ips) == 0 {
		return nil, fmt.Errorf("no entry found for %s:%s", domain, sourceIp)
	}

	return ips, nil
}

func (rs *redisStore) Unregister(domain string, sourceIp string) error {
	key := store.MakeKey(domain, sourceIp)

	return rs.client.Del(context.Background(), key).Err()
}
