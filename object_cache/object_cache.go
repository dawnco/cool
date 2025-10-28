package object_cache

import (
	"io"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

const minLiveBuffer = 100 * time.Millisecond

type Cache[V any] struct {
	data          sync.Map // map[string]*cacheItem[V]
	liveDuration  time.Duration
	cleanInterval time.Duration
	stopChan      chan struct{}
}

// NewCache 创建一个泛型缓存，live 是最大未访问时间（秒），cleanInterval 是清理周期（秒）
func NewCache[V any](liveSeconds, cleanIntervalSeconds int) *Cache[V] {
	cache := &Cache[V]{
		liveDuration:  time.Duration(liveSeconds) * time.Second,
		cleanInterval: time.Duration(cleanIntervalSeconds) * time.Second,
		stopChan:      make(chan struct{}),
	}
	go cache.startCleaner()
	return cache
}

// Set 设置键值
func (c *Cache[V]) Set(key string, value V) {
	item := &cacheItem[V]{
		value:      value,
		lastAccess: time.Now(),
	}
	c.data.Store(key, item)
}

// Get 获取键值，刷新 lastAccess
func (c *Cache[V]) Get(key string) (V, bool) {
	var zero V
	v, ok := c.data.Load(key)
	if !ok {
		return zero, false
	}
	item := v.(*cacheItem[V])
	item.touch()
	return item.value, true
}

// Delete 删除缓存项并调用 Close（如果实现了 io.Closer）
func (c *Cache[V]) Delete(key string) {
	v, ok := c.data.LoadAndDelete(key)
	if ok {
		item := v.(*cacheItem[V])
		if closer, ok := any(item.value).(io.Closer); ok {
			if err := closer.Close(); err != nil {
				logx.Errorf("Error closing cache item %v", err)
			}
		}
	}
}

// startCleaner 删除过期缓存
func (c *Cache[V]) startCleaner() {
	ticker := time.NewTicker(c.cleanInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.loopClear()
		case <-c.stopChan:
			return
		}
	}
}

func (c *Cache[V]) loopClear() {
	c.data.Range(func(key, value interface{}) bool {
		item := value.(*cacheItem[V])
		if item.expired(c.liveDuration) {
			k, ok := key.(string)
			if !ok {
				return true
			}
			c.Delete(k)
		}
		return true
	})
}

// Close 停止后台清理协程
func (c *Cache[V]) Close() {
	select {
	case <-c.stopChan:
		return
	default:
		close(c.stopChan)
	}
}
