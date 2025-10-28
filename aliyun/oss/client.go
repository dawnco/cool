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
func NewClient(c Cfg) *Client {
	client := &Client{
		accessId:  c.AccessId,
		accessKey: c.AccessKey,
		bucket:    c.Bucket,
		region:    c.Region,
		internal:  c.Internal,
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

// Put Parameters:
//
//	objectName - objectName 不能以 / 开头
//	content - 数据内容
//	expire  过期时间  YYYY-MM-DD 格式  不含当天  utc0时区
func (s *Client) Put(objectName string, content []byte, expire *time.Time) (string, error) {

	client := s.client

	// expire 过期日期
	request := &oss.PutObjectRequest{
		Bucket: oss.Ptr(s.bucket),
		Key:    oss.Ptr(strings.TrimLeft(objectName, "/")),
		Body:   bytes.NewReader(content),
	}

	if expire != nil {
		request.Expires = oss.Ptr(expire.Format("2006-01-02T15:04:05.000Z"))
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

// List 指定列表文件的列表
func (s *Client) List(prefix string, limit int) ([]ObjectItem, error) {

	bucketName := s.bucket
	client := s.client

	// Create the Paginator for the ListObjectsV2 operation.
	p := client.NewListObjectsV2Paginator(&oss.ListObjectsV2Request{
		Bucket: oss.Ptr(bucketName),
		Prefix: oss.Ptr(prefix),
	})

	ret := make([]ObjectItem, 10)
	// Iterate through the object pages
	var i int
	for p.HasNext() {
		i++
		page, err := p.NextPage(context.TODO())
		if err != nil {
			return nil, fmt.Errorf("failed to get page %v, %v", i, err)
		}

		// Print the objects found
		for _, obj := range page.Contents {
			//fmt.Printf("Object:%v, %v, %v\n", oss.ToString(obj.Key), obj.Size, oss.ToTime(obj.LastModified))
			ret = append(ret, ObjectItem{
				LastUpdate: *obj.LastModified,
				Key:        *obj.Key,
				Size:       obj.Size,
			})
		}
		if i >= limit {
			break
		}
	}

	return ret, nil

}
