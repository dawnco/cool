package utils

import (
	"net/http"
	"time"
)

type ClientWithLimit struct {
	client       *http.Client
	transport    *http.Transport
	requestCount int
	maxRequests  int
}

func NewClientWithLimit(maxRequests int) *ClientWithLimit {
	transport := &http.Transport{
		MaxIdleConns:        100,              // 最大空闲连接数
		MaxIdleConnsPerHost: 10,               // 每个 host 的最大空闲连接数
		IdleConnTimeout:     90 * time.Second, // 空闲连接超时
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second, // 请求超时
	}
	return &ClientWithLimit{
		client:       client,
		transport:    transport,
		requestCount: 0,
		maxRequests:  maxRequests,
	}
}

func (c *ClientWithLimit) Get(url string) (*http.Response, error) {
	// 如果请求次数达到了最大值，重新创建客户端
	if c.requestCount >= c.maxRequests {
		// 释放之前的连接池
		c.transport.CloseIdleConnections()

		// 重新创建新的客户端
		c.client = &http.Client{
			Transport: c.transport,
			Timeout:   10 * time.Second,
		}
		// 重置请求计数
		c.requestCount = 0
	}

	// 发起请求
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}

	// 增加请求次数
	c.requestCount++
	return resp, nil
}
