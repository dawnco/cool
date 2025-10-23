package utils

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

var client *http.Client

func init() {
	transport := &http.Transport{
		MaxIdleConns:        100,              // 最大空闲连接数
		MaxIdleConnsPerHost: 10,               // 每个 host 的最大空闲连接数
		IdleConnTimeout:     90 * time.Second, // 空闲连接超时
	}
	client = &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second, // 请求超时
	}
}

func Get(url string) (code int, respBody []byte, err error) {

	resp, err := client.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	code = resp.StatusCode

	respBody, err = io.ReadAll(resp.Body)
	return
}

func Post(url string, reqBody []byte, headers map[string]string) (StatusCode int, respBody []byte, err error) {

	req, err := http.NewRequest("POST", url, bytes.NewReader(reqBody))

	if err != nil {
		return
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	StatusCode = resp.StatusCode

	respBody, err = io.ReadAll(resp.Body)

	return
}
