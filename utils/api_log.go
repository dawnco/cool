package utils

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"sync"

	"github.com/dawnco/cool/env"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	wApiLogConn *net.UDPConn
	wApiLogOnce sync.Once
)

type kv struct {
	Key string `json:"key"`
	Val string `json:"val"`
}

// ApiLog 错误日志
// t 秒时间戳
func ApiLog(t int, topic string, store string, data map[string]string) {

	// 确保连接只初始化一次
	var initErr error

	wApiLogOnce.Do(func() {

		hostAndPort := env.Get("API_ADDR_LOG", "global.log.stat.com:8844")
		addr, err := net.ResolveUDPAddr("udp", hostAndPort)
		if err != nil {
			initErr = fmt.Errorf("failed to resolve address: %v", err)
		}
		wApiLogErrorConn, err = net.DialUDP("udp", nil, addr)
		if err != nil {
			initErr = fmt.Errorf("failed to initialize UDP connection: %v", err)
		}
	})

	// 如果初始化失败，返回错误
	if initErr != nil {
		logx.Error(fmt.Sprintf("WApiLog init error: %s", initErr.Error()))
		return
	}

	// "t":         strconv.FormatInt(eTime, 10),
	//			"date":      time.Now().Format("2006-01-02"),

	row := map[string]any{}
	row["t"] = t
	row["topic"] = topic
	row["store"] = store

	kvs := []kv{}

	for k, v := range data {
		kvs = append(kvs, kv{Key: k, Val: fmt.Sprintf("%v", v)})
	}

	row["kv"] = kvs

	message, err := json.Marshal(row)

	if err != nil {
		logx.Error(fmt.Sprintf("WApiLogError data error: %s", err.Error()))
		return
	}

	// 初始化一个切片，包含总共 7 个字节
	prefix := make([]byte, 7)
	// 前两个字节表示整数 61（大端字节序）
	binary.BigEndian.PutUint16(prefix[0:2], 60)
	// 第三个字节固定为 0
	prefix[2] = 0
	// 第 4 到第 7 个字节表示整数 230（大端字节序） 取了第 4 个字节到第 7 个字节（包括 3，不包括 7），长度正好是 4 个字节。
	binary.BigEndian.PutUint32(prefix[3:7], uint32(len(message)))

	message = append(prefix, message...)
	// 发送数据
	_, err = wApiLogConn.Write(message)
	if err != nil {
		logx.Error(fmt.Sprintf("wApiLogConn send error: %s", err.Error()))
	}
}
