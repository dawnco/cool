package client

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

type HttpClient struct {
	http      *http.Client
	timeout   time.Duration
	url       string
	method    string
	header    map[string]string
	reqBody   []byte
	resBody   []byte
	resStatus int
	err       error
}

func (s *HttpClient) SetTimeout(timeout time.Duration) *HttpClient {
	s.timeout = timeout
	return s
}

func (s *HttpClient) SetUrl(url string) *HttpClient {
	s.url = url
	return s
}

func (s *HttpClient) SetHeader(k, v string) *HttpClient {
	s.header[k] = v
	return s
}

func (s *HttpClient) SetHeaderJson() *HttpClient {
	s.header["Content-Type"] = "application/json"
	return s
}

func (s *HttpClient) SetMethod(method string) *HttpClient {
	s.method = method
	return s
}

func (s *HttpClient) SetBody(body []byte) *HttpClient {
	s.reqBody = body
	return s
}

func (s *HttpClient) Do() *HttpClient {

	var req *http.Request
	var err error

	if len(s.resBody) == 0 {
		req, err = http.NewRequest(s.method, s.url, nil)
	} else {
		req, err = http.NewRequest(s.method, s.url, bytes.NewBuffer(s.reqBody))
	}

	if err != nil {
		s.err = err
		return s
	}

	// 设置请求头
	for k, v := range s.header {
		req.Header.Set(k, v)
	}

	innerClient := &http.Client{
		Timeout: s.timeout,
	}

	resp, err := innerClient.Do(req)
	if err != nil {
		s.err = err
		return s
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.err = err
		return s
	}
	s.resStatus = resp.StatusCode
	s.resBody = body

	return s
}

func (s *HttpClient) GetErr() error {
	return s.err
}

func (s *HttpClient) GetBody() []byte {
	return s.resBody
}
func (s *HttpClient) GetStatus() int {
	return s.resStatus
}

func NewHttpClient() *HttpClient {
	return &HttpClient{
		timeout: 10 * time.Second,
		header:  make(map[string]string),
	}
}
