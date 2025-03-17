package utils

import (
	"net"
	"net/url"
)

// IsValidURL 判断是否是一个合法的 URL
func IsValidURL(urlStr string) bool {
	_, err := url.ParseRequestURI(urlStr)
	return err == nil
}

// GetOutboundIP 获取 IP
func GetOutboundIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String(), nil
}
