package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dawnco/cool/env"
	"github.com/zeromicro/go-zero/core/logx"
)

var lastAlertTime time.Time

// AlertMessage 发送消息到飞书
func AlertMessage(msg string, projectName string) {
	go clientSend(msg, "Notice", projectName)
}

func AlertMessageAndTitle(msg, projectName string, title string) {
	go clientSend(msg, projectName, title)
}

func clientSend(msg, projectName, title string) error {

	tn := time.Now()

	if tn.Sub(lastAlertTime) < time.Second*10 {
		// 10s 内只发送一条
		return fmt.Errorf("请求频繁了")
	}
	lastAlertTime = tn

	// 创建一个 map 并为每个键赋值
	n := time.Now().Format("2006-01-02 15:04:05")
	data := map[string]any{
		"key":      "",            // 消息key 用于区别消息类型和推送
		"title":    title,         // 消息内容
		"summary":  msg,           // 消息内容
		"receive":  "cqkaifa",     // 发送给谁 那个群多个用逗号分开
		"at":       "",            // at谁 多个逗号分割
		"ip":       GetServerIP(), // 警告或者恢复   默认告警
		"wr":       "w",           // 发生的ip
		"service":  projectName,   // 发生的服务
		"datetime": n,             // 发生时间 格式 2024-01-22 23:20:10
		"debug":    "",            // 调试信息
		"HideExt":  0,             // 隐藏调试信息
		"MsgId":    "",            // 编辑消息的消息ID
	}

	// 将 map 转换为 JSON 格式
	reqJson, err := json.Marshal(data)
	if err != nil {
		logx.Errorf("发送告警失败 %s : 告警消息 %s",
			err.Error(),
			msg,
		)
		return err
	}

	hostAndPort := env.Get("API_ADDR_FEISHU_ALERT", "http://feishu.message.api.com:8081/feishu/message")
	req, err := http.NewRequest("POST", hostAndPort, bytes.NewReader(reqJson))
	if err != nil {
		logx.Errorf("发送告警失败 %s : 告警消息 %s",
			err.Error(),
			msg,
		)
		return err
	}

	req.Header.Set("content-type", "application/json")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logx.Errorf("发送告警失败 %s : 告警消息 %s",
			err.Error(),
			msg,
		)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {

		body, _ := io.ReadAll(resp.Body)

		err = fmt.Errorf("发送告警响应异常 statusCode %d body %s : 告警消息 %s",
			resp.StatusCode,
			string(body),
			msg,
		)
		logx.Error(err.Error())
		return err

	}

	return nil
}
