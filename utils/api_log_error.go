package utils

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/dawnco/cool/env"
	"github.com/zeromicro/go-zero/core/logx"
)

// 全局变量
var (
	wApiLogConn     *net.UDPConn
	wApiLogConnOnce sync.Once
)

// ApiLogError 参数 https://ek8l1y505u.feishu.cn/wiki/OkvrwCLmWiRiBpkiC6dc6yjgn7d
func ApiLogError(data map[string]any, ProjectName string) {

	// 确保连接只初始化一次
	var initErr error

	wApiLogConnOnce.Do(func() {

		hostAndPort := env.Get("API_ADDR_LOG_ERROR", "logerror.stat.com:9823")
		addr, err := net.ResolveUDPAddr("udp", hostAndPort)
		if err != nil {
			initErr = fmt.Errorf("failed to resolve address: %v", err)
		}
		wApiLogConn, err = net.DialUDP("udp", nil, addr)
		if err != nil {
			initErr = fmt.Errorf("failed to initialize UDP connection: %v", err)
		}
	})

	// 如果初始化失败，返回错误
	if initErr != nil {
		logx.Error(fmt.Sprintf("WApiLogError init error: %s", initErr.Error()))
		return
	}

	// "t":         strconv.FormatInt(eTime, 10),
	//			"date":      time.Now().Format("2006-01-02"),

	if _, exists := data["t"]; !exists {
		data["t"] = strconv.FormatInt(time.Now().UnixMilli(), 10)
	}

	if _, exists := data["service"]; !exists {
		data["service"] = ProjectName
	}

	if _, exists := data["date"]; !exists {
		data["date"] = time.Now().Format("2006-01-02")
	}

	if _, exists := data["ip"]; !exists {
		data["ip"] = GetServerIP()
	}

	if _, exists := data["file"]; !exists {
		_, file, line, ok := runtime.Caller(2)
		if ok {
			data["file"] = file
			data["line"] = line
		}
	}

	message, err := json.Marshal(data)

	if err != nil {
		logx.Error(fmt.Sprintf("WApiLogError data error: %s", err.Error()))
		return
	}

	// 初始化一个切片，包含总共 7 个字节
	prefix := make([]byte, 7)
	// 前两个字节表示整数 61（大端字节序）
	binary.BigEndian.PutUint16(prefix[0:2], 61)
	// 第三个字节固定为 0
	prefix[2] = 0
	// 第 4 到第 7 个字节表示整数 230（大端字节序） 取了第 4 个字节到第 7 个字节（包括 3，不包括 7），长度正好是 4 个字节。
	binary.BigEndian.PutUint32(prefix[3:7], uint32(len(message)))

	message = append(prefix, message...)
	// 发送数据
	_, err = wApiLogConn.Write(message)
	if err != nil {
		logx.Error(fmt.Sprintf("WApiLogError send error: %s", err.Error()))
	}
}
