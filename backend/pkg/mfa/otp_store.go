package mfa

import (
	"time"

	"github.com/redis/go-redis/v9"
)

// OTPStore defines the standard interface for storing OTPs
type OTPStore interface {
	Set(key, code string, ttl time.Duration) error
	Get(key string) (string, bool)
	Delete(key string) error
}

// NewOTPStore returns a Redis-based or in-memory OTP store
func NewOTPStore(redisClient *redis.Client) OTPStore {
	if redisClient != nil {
		return &RedisOTPStore{client: redisClient}
	}
	return NewInMemoryOTPStore()
}
