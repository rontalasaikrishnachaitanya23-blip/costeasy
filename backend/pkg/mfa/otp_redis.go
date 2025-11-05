package mfa

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisOTPStore struct {
	client *redis.Client
}

func NewRedisOTPStore(client *redis.Client) *RedisOTPStore {
	return &RedisOTPStore{client: client}
}

func (r *RedisOTPStore) Set(key, code string, ttl time.Duration) error {
	ctx := context.Background()
	return r.client.Set(ctx, key, code, ttl).Err()
}

func (r *RedisOTPStore) Get(key string) (string, bool) {
	ctx := context.Background()
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", false
	}
	if err != nil {
		return "", false
	}
	return val, true
}

func (r *RedisOTPStore) Delete(key string) error {
	ctx := context.Background()
	return r.client.Del(ctx, key).Err()
}
