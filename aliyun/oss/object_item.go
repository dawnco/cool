package oss

import "time"

type ObjectItem struct {
	Key        string // 文件 含路径
	LastUpdate time.Time
	Size       int64
}
