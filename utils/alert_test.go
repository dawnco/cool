// ip_test.go
package utils

import (
	"testing"
)

func TestAlertMessage(t *testing.T) {
	// 测试 GetServerIP 函数
	err := clientSend("测试body", "测试Project", "测试Title")

	// 验证返回非空IP地址
	if err != nil {
		t.Errorf("clientSend() should not return error get %v", err)
	}

}
