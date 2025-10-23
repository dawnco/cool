// ip_test.go
package utils

import (
	"net"
	"testing"
)

func TestGetServerIP(t *testing.T) {
	// 测试 GetServerIP 函数
	ip := GetServerIP()

	// 验证返回非空IP地址
	if ip == "" {
		t.Error("GetServerIP() should not return empty string")
	}

	// 验证返回格式为有效的IPv4地址
	if ip != "0.0.0.0" && net.ParseIP(ip) == nil {
		t.Errorf("GetServerIP() should return valid IP address, got %s", ip)
	}

	// 验证是IPv4地址
	if ip != "0.0.0.0" && net.ParseIP(ip).To4() == nil {
		t.Errorf("GetServerIP() should return IPv4 address, got %s", ip)
	}
}

func TestGetServerLocalIP(t *testing.T) {
	// 测试 GetServerLocalIP 函数
	ip := GetServerLocalIP()

	// 验证能成功获取本地IP地址
	if ip == "" {
		t.Error("GetServerLocalIP() should not return empty string")
	}

	// 验证返回格式为有效的IPv4地址
	if net.ParseIP(ip) == nil {
		t.Errorf("GetServerLocalIP() should return valid IP address, got %s", ip)
	}

	// 验证是IPv4地址
	if net.ParseIP(ip).To4() == nil {
		t.Errorf("GetServerLocalIP() should return IPv4 address, got %s", ip)
	}
}

func TestGetIpFromDomain(t *testing.T) {
	// 测试 GetIpFromDomain 函数
	tests := []struct {
		name   string
		domain string
		wantIP bool
	}{
		{
			name:   "正常域名解析",
			domain: "google.com",
			wantIP: true,
		},
		{
			name:   "无效域名",
			domain: "this-domain-should-not-exist.invalid",
			wantIP: false,
		},
		{
			name:   "空域名",
			domain: "",
			wantIP: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetIpFromDomain(tt.domain)

			if tt.wantIP {
				// 期望获取到IP地址
				if result == "" {
					t.Errorf("GetIpFromDomain(%s) should return IP address, got empty string", tt.domain)
				}
				// 验证返回的是有效IP地址
				if result != "" && net.ParseIP(result) == nil {
					t.Errorf("GetIpFromDomain(%s) should return valid IP address, got %s", tt.domain, result)
				}
			} else {
				// 不期望获取到IP地址
				if result != "" {
					t.Errorf("GetIpFromDomain(%s) should return empty string, got %s", tt.domain, result)
				}
			}
		})
	}
}
