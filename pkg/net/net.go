package net

import (
	"net"
	"strings"
)

func LocalIP() ([]string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	ips := make([]string, 0)
	for _, i := range interfaces {
		if i.Flags&net.FlagLoopback != 0 || i.Flags&net.FlagUp == 0 {
			continue
		}
		addrs, err := i.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip4 := addr.(*net.IPNet).IP.To4()
			if ip4 != nil {
				ips = append(ips, ip4.String())
			}
		}
	}
	return ips, nil
}

func GetOutBoundIP() (string, error) {
	// 向任一外网地址发送一个包，从本地连接中获取自己对外的IP地址
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		return "", nil
	}
	ip := conn.LocalAddr().String()
	return strings.Split(ip, ":")[0], nil
}
