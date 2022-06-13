package util

import (
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetRealIp(c *gin.Context) string {
	r := c.Request
	ip := c.ClientIP()
	if ip == "" {
		ip = ClientPublicIP(r)
	}
	if ip == "" {
		ip = ClientIp(r)
	}
	return ip
}

func ClientIp(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}

	ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil {
		return ip
	}
	return ""
}

// ClientPublicIP 尽最大努力实现获取客户端公网 IP 的算法。
// 解析 X-Real-IP 和 X-Forwarded-For 以便于反向代理（nginx 或 haproxy）可以正常工作。
func ClientPublicIP(r *http.Request) string {
	var ip string
	for _, ip = range strings.Split(r.Header.Get("X-Forwarded-For"), ",") {
		ip = strings.TrimSpace(ip)
		if ip != "" && !IsLocalIp(ip) {
			return ip
		}
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" && !IsLocalIp(ip) {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		if !IsLocalIp(ip) {
			return ip
		}
	}

	return ""
}

func IsLocalIp(ip string) bool {
	ipAddr := strings.Split(ip, ".")

	if strings.EqualFold(ipAddr[0], "10") {
		return true
	} else if strings.EqualFold(ipAddr[0], "172") {
		addr, _ := strconv.Atoi(ipAddr[1])
		if addr >= 16 && addr < 31 {
			return true
		}
	} else if strings.EqualFold(ipAddr[0], "192") && strings.EqualFold(ipAddr[1], "168") {
		return true
	}
	return false
}
