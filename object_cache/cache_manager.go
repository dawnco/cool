package object_cache

import (
	"sync"
)

// CacheManager 是基于 Cache 封装的带“读写互斥 + 创建工厂”的缓存管理器
type CacheManager[V any] struct {
	cache *Cache[V]
	mu    sync.Mutex
}

func NewCacheManager[V any](liveSeconds, cleanIntervalSeconds int) *CacheManager[V] {
	return &CacheManager[V]{
		cache: NewCache[V](liveSeconds, cleanIntervalSeconds),
	}
}

// GetOrCreate 获取指定 key 的缓存对象，若不存在则调用 factory 创建并存入缓存
func (m *CacheManager[V]) GetOrCreate(key string, factory func() (V, error)) (V, error) {
	// 先快速查找
	if val, ok := m.cache.Get(key); ok {
		return val, nil
	}

	// 加锁保证同一时刻只有一个 goroutine 执行 factory
	m.mu.Lock()
	defer m.mu.Unlock()

	// 再次检查缓存，避免重复创建
	if val, ok := m.cache.Get(key); ok {
		return val, nil
	}

	// 创建新对象
	val, err := factory()
	if err != nil {
		var zero V
		return zero, err
	}

	m.cache.Set(key, val)
	return val, nil
}

// Close 停止内部缓存的清理 goroutine
func (m *CacheManager[V]) Close() {
	m.cache.Close()
}
