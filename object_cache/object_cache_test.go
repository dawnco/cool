package object_cache

import (
	"sync/atomic"
	"testing"
	"time"
)

// 测试用的模拟资源，实现 io.Closer
type mockResource struct {
	closed int32
}

func (r *mockResource) Close() error {
	atomic.StoreInt32(&r.closed, 1)
	return nil
}

func (r *mockResource) IsClosed() bool {
	return atomic.LoadInt32(&r.closed) == 1
}

func TestCacheBasic(t *testing.T) {
	cache := NewCache[*mockResource](1, 1) // live=1s, cleanInterval=1s
	defer cache.Close()

	key := "res1"
	res := &mockResource{}

	cache.Set(key, res)

	// 取出应该存在
	got, ok := cache.Get(key)
	if !ok || got != res {
		t.Fatalf("expected to get resource, got %v, ok=%v", got, ok)
	}

	// 刷新访问时间，防止清理
	time.Sleep(500 * time.Millisecond)
	got2, ok2 := cache.Get(key)
	if !ok2 || got2 != res {
		t.Fatalf("expected to get resource again, got %v, ok=%v", got2, ok2)
	}

	// 删除
	cache.Delete(key)
	_, ok3 := cache.Get(key)
	if ok3 {
		t.Fatal("expected resource to be deleted")
	}
	if !res.IsClosed() {
		t.Fatal("expected resource to be closed after delete")
	}
}

func TestCacheAutoCleanup(t *testing.T) {
	cache := NewCache[*mockResource](1, 3) // live=1s, cleanInterval=10s
	defer cache.Close()

	key := "res2"
	res := &mockResource{}
	cache.Set(key, res)

	// 不访问，等待清理
	time.Sleep(7 * time.Second) // > live+buffer

	_, ok := cache.Get(key)
	if ok {
		t.Fatal("expected resource to be auto deleted")
	}

	if !res.IsClosed() {
		t.Fatal("expected resource to be closed after auto cleanup")
	}
}

func TestCacheNoPanicOnNonStringKey(t *testing.T) {
	cache := NewCache[int](1, 1)
	defer cache.Close()

	// 直接存储非字符串key（通过sync.Map接口存），模拟异常情况
	cache.data.Store(123, &cacheItem[int]{value: 42, lastAccess: time.Now().Add(-time.Hour)})

	// startCleaner 会 Range 并跳过非string key
	// 为了测试，可以直接调用一次清理函数

	cache.loopClear()

	// key 123 应该还在，因为我们没有删除它
	if _, ok := cache.data.Load(123); !ok {
		t.Fatal("expected key 123 to still be present")
	}
}
