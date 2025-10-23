package redis

import (
	"time"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

var redisInstance = make(map[string]*redis.Redis)

func GetDefaultClient() *redis.Redis {
	return GetClient("default")
}

func GetClient(name string) *redis.Redis {
	r := redisInstance[name]
	return r
}

func Init(c []Cfg) {
	for _, v := range c {
		redisInstance[v.Name] = redis.MustNewRedis(redis.RedisConf{
			Host:        v.Host,
			Type:        "node",
			Pass:        v.Pass,
			Tls:         false,
			NonBlock:    true,
			PingTimeout: time.Second * 2,
		})
	}
}
