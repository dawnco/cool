package wredis

import (
	"fmt"

	"github.com/redis/go-redis/v9"

	"sync"
)

var instance = sync.Map{}

// Init 初始化
// name 配置名称, 后面通过 GetClient 获取这个配置的客户端
// cfg 配置参数
func Init(name string, cfg Cfg) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.Db,
	})
	instance.Store(name, client)

}

func Get(name string) *redis.Client {
	conn, ok := instance.Load(name)
	if !ok {
		panic(fmt.Errorf("redis connection %s not found", name))
	}
	return conn.(*redis.Client)
}
