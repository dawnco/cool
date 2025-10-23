package redis

import (
	"encoding/json"
	"fmt"
)

func Set(key string, value any, ttl int) error {

	r := GetDefaultClient()

	switch value.(type) {
	case int8, int16, int32, int, int64,
		uint8, uint16, uint32, uint, uint64:
		return r.Setex(key, fmt.Sprintf("%d", value), ttl)

	case float32, float64:
		return r.Setex(key, fmt.Sprintf("%f", value), ttl)

	case string:
		return r.Setex(key, value.(string), ttl)
	}

	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.Setex(key, string(jsonData), ttl)
}

func Get(key string) (string, error) {
	return GetDefaultClient().Get(key)
}
