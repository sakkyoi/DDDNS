package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/sakkyoi/DDDNS/internal/store"
	"time"
)

type redisStore struct {
	client *redis.Client
}

func New(addr string) store.Store {
	client := redis.NewClient(&redis.Options{Addr: addr})

	return &redisStore{
		client: client,
	}
}

func (rs *redisStore) Register(domain string, sourceIP string, destIP string, ttl time.Duration) error {
	key := store.MakeKey(domain, sourceIP)

	return rs.client.Set(context.Background(), key, destIP, ttl).Err()
}

func (rs *redisStore) Lookup(domain string, sourceIP string) (string, error) {
	key := store.MakeKey(domain, sourceIP)

	return rs.client.Get(context.Background(), key).Result()
}

func (rs *redisStore) Unregister(domain string, sourceIP string) error {
	key := store.MakeKey(domain, sourceIP)

	return rs.client.Del(context.Background(), key).Err()
}
