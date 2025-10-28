package object_cache

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// 模拟资源，实现 io.Closer（可选）
type mockObj struct {
	val    string
	closed int32
}

func (m *mockObj) Close() error {
	atomic.StoreInt32(&m.closed, 1)
	return nil
}

func (m *mockObj) IsClosed() bool {
	return atomic.LoadInt32(&m.closed) == 1
}

func TestCacheManager_GetOrCreate_Basic(t *testing.T) {
	mgr := NewCacheManager[*mockObj](5, 1)
	defer mgr.Close()

	key := "a"

	calls := 0
	factory := func() (*mockObj, error) {
		calls++
		return &mockObj{val: "hello"}, nil
	}

	obj1, err := mgr.GetOrCreate(key, factory)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if obj1 == nil || obj1.val != "hello" {
		t.Fatal("object not created as expected")
	}

	obj2, _ := mgr.GetOrCreate(key, factory)
	if obj1 != obj2 {
		t.Fatal("expected cached object, got different instances")
	}
	if calls != 1 {
		t.Fatalf("factory called %d times, expected 1", calls)
	}
}

func TestCacheManager_GetOrCreate_Concurrent(t *testing.T) {
	mgr := NewCacheManager[*mockObj](5, 1)
	defer mgr.Close()

	key := "concurrent"
	var callCount int32

	factory := func() (*mockObj, error) {
		atomic.AddInt32(&callCount, 1)
		time.Sleep(100 * time.Millisecond)
		return &mockObj{val: "concurrent"}, nil
	}

	const n = 10
	var wg sync.WaitGroup
	wg.Add(n)

	results := make([]*mockObj, n)

	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			obj, err := mgr.GetOrCreate(key, factory)
			if err != nil {
				t.Errorf("goroutine %d error: %v", i, err)
			}
			results[i] = obj
		}(i)
	}

	wg.Wait()

	ref := results[0]
	for i := 1; i < n; i++ {
		if results[i] != ref {
			t.Fatalf("not all goroutines got the same instance")
		}
	}

	if callCount != 1 {
		t.Fatalf("expected factory to be called once, got %d", callCount)
	}
}

func TestCacheManager_GetOrCreate_FactoryError(t *testing.T) {
	mgr := NewCacheManager[int](5, 1)
	defer mgr.Close()

	key := "err-key"

	factory := func() (int, error) {
		return 0, errors.New("mock failure")
	}

	_, err := mgr.GetOrCreate(key, factory)
	if err == nil || err.Error() != "mock failure" {
		t.Fatalf("expected factory error, got %v", err)
	}

	// 第二次调用也应再次触发 factory，因为没缓存
	_, err2 := mgr.GetOrCreate(key, factory)
	if err2 == nil {
		t.Fatal("expected error again")
	}
}

func TestCacheManager_GetOrCreate_ConcurrentAccess(t *testing.T) {
	mgr := NewCacheManager[*mockObj](10, 2)
	defer mgr.Close()

	var factoryCalled int32
	key := "shared"

	factory := func() (*mockObj, error) {
		atomic.AddInt32(&factoryCalled, 1)
		time.Sleep(100 * time.Millisecond) // 模拟创建开销
		return &mockObj{val: "shared-object"}, nil
	}

	const goroutineCount = 50
	var wg sync.WaitGroup
	wg.Add(goroutineCount)

	results := make([]*mockObj, goroutineCount)

	for i := 0; i < goroutineCount; i++ {
		go func(idx int) {
			defer wg.Done()
			obj, err := mgr.GetOrCreate(key, factory)
			if err != nil {
				t.Errorf("goroutine %d error: %v", idx, err)
				return
			}
			results[idx] = obj
		}(i)
	}

	wg.Wait()

	// 所有结果都应是同一个对象
	ref := results[0]
	for i := 1; i < goroutineCount; i++ {
		if results[i] != ref {
			t.Errorf("result[%d] != result[0]", i)
		}
	}

	if n := atomic.LoadInt32(&factoryCalled); n != 1 {
		t.Errorf("expected factory to be called once, but was called %d times", n)
	}
}
