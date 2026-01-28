package utils

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/dawnco/cool/env"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	wEventCenterConn *net.UDPConn
	wEventCenterOnce sync.Once
)

func ApiEc(requestId string, name string, projectName string, params any) int {
	return ApiEventCenter(
		requestId,
		name,
		"",
		projectName,
		time.Now().UnixMilli(),
		params,
	)
}

func ApiEventCenter(
	requestId string,
	name string,
	topic string,
	from string,
	timeMs int64,
	params any) int {

	// 确保连接只初始化一次
	var initErr error

	wEventCenterOnce.Do(func() {
		hostAndPort := env.Get("API_ADDR_EVENT_CENTER", "center.stat.com:9820")
		addr, err := net.ResolveUDPAddr("udp", hostAndPort)
		if err != nil {
			initErr = fmt.Errorf("failed to resolve address: %v", err)
		}
		wEventCenterConn, err = net.DialUDP("udp", nil, addr)
		if err != nil {
			initErr = fmt.Errorf("failed to initialize UDP connection: %v", err)
		}
	})

	// 如果初始化失败，返回错误
	if initErr != nil {
		logx.Error(fmt.Sprintf("WEventCenter init error: %s", initErr.Error()))
		return 0
	}

	if wEventCenterConn == nil {
		logx.Error("WEventCenter conn is nil")
		return 0
	}

	data := map[string]any{
		"name":      name,
		"_topic_":   topic,
		"requestId": requestId,
		"from":      from,
		"country":   "bd",
		"timestamp": timeMs,
		"params":    params,
	}

	message, err := json.Marshal(data)

	if err != nil {
		logx.Error(fmt.Sprintf("WEventCenter data error: %s", err.Error()))
		return 0
	}

	if len(message) > 65535 {
		logx.Error(fmt.Sprintf("WEventCenter message too long: %d", len(message)))
		return 0
	}

	// 初始化一个切片，包含总共 5 个字节
	prefix := make([]byte, 5)
	// 前两个字节表示整数  55 表示
	binary.BigEndian.PutUint16(prefix[0:2], 55)
	// 第三个字节固定为 0
	prefix[2] = 0
	// 第 4 到第 5 个 字节表示内容长度
	binary.BigEndian.PutUint16(prefix[3:5], uint16(len(message)))

	message = append(prefix, message...)
	// 发送数据
	sendLen, err := wEventCenterConn.Write(message)
	if err != nil {
		logx.Error(fmt.Sprintf("WEventCenter send error: %s", err.Error()))
		return 0
	}
	return sendLen

}
