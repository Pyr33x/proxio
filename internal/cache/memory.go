package cache

import (
	"context"
	"sync"
	"time"
)

type memoryEntry struct {
	value      []byte
	expiration time.Time
}

type memoryStore struct {
	data sync.Map
}

func NewMemoryStore() *memoryStore {
	return &memoryStore{}
}

func (m *memoryStore) Get(ctx context.Context, key string) ([]byte, error) {
	raw, ok := m.data.Load(key)
	if !ok {
		return nil, nil
	}

	entry := raw.(memoryEntry)

	if !entry.expiration.IsZero() && time.Now().After(entry.expiration) {
		m.data.Delete(key)
		return nil, nil
	}

	return entry.value, nil
}

func (m *memoryStore) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	var exp time.Time
	if ttl > 0 {
		exp = time.Now().Add(ttl)
	}

	m.data.Store(key, memoryEntry{
		value:      value,
		expiration: exp,
	})

	return nil
}

func (m *memoryStore) Clear(ctx context.Context) error {
	m.data.Range(func(key, _ any) bool {
		m.data.Delete(key)
		return true
	})
	return nil
}
