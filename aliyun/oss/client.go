package oss

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

type Client struct {
	client    *oss.Client
	internal  bool
	accessId  string
	accessKey string
	bucket    string
	region    string
}

func (s *Client) init() {
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(s.accessId, s.accessKey)).
		WithRegion(s.region).WithUseInternalEndpoint(s.internal)

	s.client = oss.NewClient(cfg)
}

// NewClient 初始化
// 例如  region = ap-southeast-1 参考 https://www.alibabacloud.com/help/zh/oss/user-guide/regions-and-endpoints
// 是否使用内网  internal
func NewClient(accessId, accessKey, region, bucket string, internal bool) *Client {
	client := &Client{
		accessId:  accessId,
		accessKey: accessKey,
		bucket:    bucket,
		region:    region,
		internal:  internal,
	}
	client.init()
	return client
}

// AuthByUrl 通过 url 签名
func (s *Client) AuthByUrl(urlString string, expire int64) (string, error) {
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return "", fmt.Errorf("error parsing URL: %s %w", parsedURL, err)
	}
	host := parsedURL.Host
	objectName := strings.TrimLeft(parsedURL.Path, "/")

	parts := strings.Split(host, ".")
	bucketName := parts[0]
	return s.AuthByObjectName(objectName, bucketName, expire)
}

// AuthByObjectName 通过 objectName 签名
func (s *Client) AuthByObjectName(objectName string, bucketName string, expire int64) (string, error) {

	client := s.client

	if bucketName == "" {
		bucketName = s.bucket
	}

	result, err := client.Presign(context.TODO(), &oss.GetObjectRequest{
		Bucket: oss.Ptr(bucketName),
		Key:    oss.Ptr(objectName),
	},
		oss.PresignExpires(time.Duration(expire)*time.Second),
	)

	if err != nil {
		return "", fmt.Errorf("签名 %s 错误 %w", objectName, err)
	}

	return result.URL, nil
}

// Put 写oss objectName 不能以 / 开头
func (s *Client) Put(objectName string, content []byte) (string, error) {

	client := s.client

	request := &oss.PutObjectRequest{
		Bucket: oss.Ptr(s.bucket),
		Key:    oss.Ptr(objectName),
		Body:   bytes.NewReader(content),
		//Body:   strings.NewReader(content),
	}

	result, err := client.PutObject(context.TODO(), request)
	if err != nil {
		return "", fmt.Errorf("OSS 上传内容错误 %w", err)
	}

	if result.ResultCommon.StatusCode != 200 {
		return "", fmt.Errorf("OSS 上传内容错误 %s", result.ResultCommon.Status)
	}

	fullUrl := fmt.Sprintf("https://%s.oss-%s.aliyuncs.com/%s",
		s.bucket,
		s.region,
		objectName,
	)
	return fullUrl, nil
}

// Delete 删除oss
func (s *Client) Delete(objectName string) error {

	client := s.client
	bucketName := s.bucket

	request := &oss.DeleteObjectRequest{
		Bucket: oss.Ptr(bucketName),
		Key:    oss.Ptr(objectName),
	}
	_, err := client.DeleteObject(context.TODO(), request)
	if err != nil {
		return err
	}
	return nil
}

// Get 获取oss
func (s *Client) Get(objectName string) ([]byte, error) {

	bucketName := s.bucket

	client := s.client

	request := &oss.GetObjectRequest{
		Bucket: oss.Ptr(bucketName),
		Key:    oss.Ptr(objectName),
	}
	result, err := client.GetObject(context.TODO(), request)
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	data, _ := io.ReadAll(result.Body)
	return data, nil
}

// Exist 是否存在
func (s *Client) Exist(objectName string) (bool, error) {

	bucketName := s.bucket
	client := s.client

	return client.IsObjectExist(context.TODO(), bucketName, objectName)

}
