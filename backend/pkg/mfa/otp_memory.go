package mfa

import (
	"sync"
	"time"
)

type otpEntry struct {
	Code      string
	ExpiresAt time.Time
}

type InMemoryOTPStore struct {
	mu   sync.Mutex
	data map[string]otpEntry
}

func NewInMemoryOTPStore() *InMemoryOTPStore {
	store := &InMemoryOTPStore{
		data: make(map[string]otpEntry),
	}
	go store.cleanupLoop()
	return store
}

func (s *InMemoryOTPStore) Set(key, code string, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = otpEntry{Code: code, ExpiresAt: time.Now().Add(ttl)}
	return nil
}

func (s *InMemoryOTPStore) Get(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, ok := s.data[key]
	if !ok || time.Now().After(entry.ExpiresAt) {
		delete(s.data, key)
		return "", false
	}
	return entry.Code, true
}

func (s *InMemoryOTPStore) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
	return nil
}

func (s *InMemoryOTPStore) cleanupLoop() {
	for range time.Tick(1 * time.Minute) {
		now := time.Now()
		s.mu.Lock()
		for k, v := range s.data {
			if now.After(v.ExpiresAt) {
				delete(s.data, k)
			}
		}
		s.mu.Unlock()
	}
}
