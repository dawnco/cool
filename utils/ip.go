package utils

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

// 获取当前机器外网 IPv4

var serverIP string
var serverLocalIP string

var lastUpdateTime time.Time

func init() {
	lastUpdateTime = time.Now()
	serverIP = getServerIP()
	serverLocalIP, _ = getServerLocalIP()
}

func GetServerIP() string {
	tn := time.Now()
	if tn.Sub(lastUpdateTime) < time.Hour {
		return serverIP
	}
	lastUpdateTime = tn
	serverIP = getServerIP()
	return serverIP
}

func GetServerLocalIP() string {
	tn := time.Now()
	if tn.Sub(lastUpdateTime) < time.Hour {
		return serverLocalIP
	}

	serverLocalIP, _ = getServerLocalIP()
	lastUpdateTime = tn

	return serverLocalIP
}

// GetServerIp 获取当前机器外网 IPv4
func getServerIP() string {

	req, err := http.NewRequest("GET", "https://ifconfig.me/", nil)
	if err != nil {
		return "0.0.0.0"
	}
	req.Header.Set("User-Agent", "curl/7.29.0")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "0.0.0.0"
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "0.0.0.0"
	}
	ipv4 := string(body)

	if isIPv4(ipv4) {
		return ipv4
	}

	if isIPV6(ipv4) {
		// 是ipv6 获取本地IP
		ip := GetServerLocalIP()
		if ip == "" {
			return "0.0.0.0"
		}
		if isIPv4(ip) {
			return ip
		}
	}

	return "0.0.0.0"
}

func getServerLocalIP() (string, error) {
	// 获取所有的网络接口
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	// 遍历每个网络接口
	for _, iface := range interfaces {
		// 忽略回环接口（lo0）
		if iface.Flags&net.FlagUp == 0 || iface.Name == "lo" {
			continue
		}

		// 获取接口的地址列表
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}

		// 遍历接口的所有地址
		for _, addr := range addrs {
			// 检查地址类型是否是 IPv4
			if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.To4() != nil {
				// 返回 IPv4 地址
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no local IP found")
}

func GetIpFromDomain(domain string) string {
	ips, err := net.LookupIP(domain)
	if err != nil {
		return ""
	}

	if len(ips) > 0 {
		return ips[0].String()
	}
	return ""
}

// 判断一个字符串是否是有效的IPv4地址
func isIPv4(ip string) bool {
	parsedIP := net.ParseIP(ip)
	// 使用 To4() 方法判断是否是IPv4地址
	return parsedIP != nil && parsedIP.To4() != nil
}

func isIPV6(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil && parsedIP.To4() == nil
}
