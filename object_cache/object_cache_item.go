package object_cache

import (
	"sync"
	"time"
)

type cacheItem[V any] struct {
	mu         sync.Mutex
	value      V
	lastAccess time.Time
}

func (i *cacheItem[V]) touch() {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.lastAccess = time.Now()
}

func (i *cacheItem[V]) expired(live time.Duration) bool {
	i.mu.Lock()
	defer i.mu.Unlock()
	return time.Since(i.lastAccess) > live+minLiveBuffer
}
