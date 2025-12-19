package wredis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisInstance = make(map[string]*redis.Client)

func Init(cfg []Cfg) {
	for _, v := range cfg {
		redisInstance[v.Name] = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", v.Host, v.Port),
			Password: v.Password,
			DB:       v.Db,
		})
	}
}

func DefaultClient() *redis.Client {
	return Client("default")
}

func Client(name string) *redis.Client {
	r := redisInstance[name]
	return r
}

func Set(key string, value any, ttl int) error {
	r := DefaultClient()
	cmd := r.Set(context.Background(), key, value, time.Duration(ttl)*time.Second)
	return cmd.Err()
}

func GetString(key string) (string, error) {
	r := DefaultClient()
	return r.Get(context.Background(), key).Result()
}
