缓存对象 使用方式

```
var mgr CacheManager[T]
init(){
    mgr = NewCacheManager[T](1, 2) // live=1s, cleanInterval=2s
}

func [T]GetInstance(key string)T{
    mgr.GetOrCreate(key, factory func() (T, error)) (T, error){
            // 如果 key不存在 返回T实例
            // 存在返回缓存的T
    })
}

	
    
    
    
    
```