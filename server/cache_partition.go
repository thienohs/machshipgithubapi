package server

import "sync"

type ICacheable interface {
	String() string
}

// ServerCachePartition a partition inside the cache
type ServerCachePartition[T ICacheable] struct {
	id        string
	cache     map[string]*T
	cacheLock *sync.RWMutex
}

// NewServerCachePartition return new cache partition
func NewServerCachePartition[T ICacheable](id string) *ServerCachePartition[T] {
	return &ServerCachePartition[T]{
		id:        id,
		cache:     make(map[string]*T),
		cacheLock: &sync.RWMutex{},
	}
}

// Get get cache value by key
func (scp *ServerCachePartition[T]) Get(key string) *T {
	scp.cacheLock.RLock()
	defer scp.cacheLock.RUnlock()
	value := scp.cache[key]
	return value
}

// Set set cache value by key
func (scp *ServerCachePartition[T]) Set(key string, value *T) {
	scp.cacheLock.Lock()
	defer scp.cacheLock.Unlock()
	scp.cache[key] = value
}
