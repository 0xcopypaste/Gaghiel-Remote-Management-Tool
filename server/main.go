package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
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
