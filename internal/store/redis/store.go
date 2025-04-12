package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sakkyoi/DDDNS/internal/store"
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
