package utils

import (
	"bytes"
	"compress/gzip"
	"io"
)

// CompressGz compressGz 压缩输入的 []byte，返回压缩后的 []byte
func CompressGz(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)

	_, err := gzWriter.Write(data)
	if err != nil {
		gzWriter.Close()
		return nil, err
	}

	// 关闭以确保所有数据都写入 buf
	err = gzWriter.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// DecompressGz decompressGz 解压缩输入的 gzip 格式 []byte，返回解压后的 []byte
func DecompressGz(data []byte) ([]byte, error) {
	buf := bytes.NewReader(data)
	gzReader, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}
	defer gzReader.Close()

	var out bytes.Buffer
	_, err = io.Copy(&out, gzReader)
	if err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}
