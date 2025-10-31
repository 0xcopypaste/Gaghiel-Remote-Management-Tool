package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

type Client struct {
	conn     net.Conn
	addr     string
	user     string
	host     string
	ip       string
	respChan chan string
}

var (
	clients      = make(map[string]*Client)
	clientsMutex sync.Mutex
)

func acceptClients(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("[SERVER] Error accepting connection:", err)
			continue
		}

		go handleClientConnection(conn)
	}
}

func handleClientConnection(c net.Conn) {
	buffer := make([]byte, 1024)
	n, err := c.Read(buffer)
	if err != nil {
		c.Close()
		return
	}

	info := strings.Split(strings.TrimSpace(string(buffer[:n])), "|")
	if len(info) < 3 {
		c.Close()
		return
	}

	client := &Client{
		conn:     c,
		user:     info[0],
		host:     info[1],
		ip:       info[2],
		addr:     c.RemoteAddr().String(),
		respChan: make(chan string, 100),
	}

	clientsMutex.Lock()
	sessionID := fmt.Sprintf("%d", len(clients)+1)
	clients[sessionID] = client
	clientsMutex.Unlock()

	fmt.Println("[SERVER] New client connected:", sessionID, "-", client.user+"@"+client.host, client.addr)

	go receiveClientMessages(client, sessionID)
}

func receiveClientMessages(client *Client, sessionID string) {
	reader := bufio.NewReader(client.conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			clientsMutex.Lock()
			delete(clients, sessionID)
			clientsMutex.Unlock()
			fmt.Println("[SERVER] Client disconnected:", sessionID)
			return
		}
		client.respChan <- strings.TrimRight(msg, "\r\n")
	}
}
