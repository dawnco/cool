package env

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

// Get 获取环境变量
func Get[T any](key string, defaultValue T) T {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	var result T
	switch any(result).(type) {
	case string:
		return any(value).(T)
	case int:
		if v, err := strconv.Atoi(value); err == nil {
			return any(v).(T)
		}
	case int64:
		if v, err := strconv.ParseInt(value, 10, 64); err == nil {
			return any(v).(T)
		}
	case float64:
		if v, err := strconv.ParseFloat(value, 64); err == nil {
			return any(v).(T)
		}
	case bool:
		if v, err := strconv.ParseBool(value); err == nil {
			return any(v).(T)
		}
	default:
		panic(fmt.Sprintf("env unsupported type: %v", reflect.TypeOf(defaultValue)))
	}

	// 转换失败，返回默认值
	return defaultValue
}
