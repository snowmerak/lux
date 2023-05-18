package util

import "net"

func GetIP(addr string) string {
	ip, _, _ := net.SplitHostPort(addr)
	return ip
}

func GetPort(addr string) string {
	_, port, _ := net.SplitHostPort(addr)
	return port
}
