package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

var (
	store = make(map[string]string)
	mu 		sync.Mutex
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		command, err := reader.ReadString('\n')

		if err != nil {
			fmt.Println("Client disconnected.")
			return
		}

		command = strings.TrimSpace(command)
		parts := strings.Fields(command)
		
		switch strings.ToUpper(parts[0]){
			case "PUT": 
				if len(parts) < 3 {
					conn.Write([]byte("Error: invalid PUT format\n"))
					continue
				}
				key := parts[1]
				value := parts[2]

				mu.Lock()
				store[key] = value
				mu.Unlock()
				conn.Write([]byte("OK\n"))
			
			
			case "GET":
				if len(parts) < 2 {
					conn.Write([]byte("Error: invalid GET format\n"))
					continue
				}
				key := parts[1]
				
				mu.Lock()
				value, exists := store[key]
				mu.Unlock()

				if !exists{
					conn.Write([]byte("Error: key not found\n"))
				}else{
					conn.Write([]byte(value + "\n"))
				}

			 
			case "DELETE":
				if len(parts) < 2 {
					conn.Write([]byte("Error: invalid DELETE format\n"))
					continue
				}

				key := parts[1]

				mu.Lock()
				_, exists := store[key]
				if exists{
					delete(store, key)
				}
				mu.Unlock()

				if exists{
					conn.Write([]byte("OK\n"))
				}else{
					conn.Write([]byte("Error: key not found\n"))
				}

			case "LIST":
				mu.Lock()
				if len(store) == 0 {
					conn.Write([]byte("Empty Store\n"))
				} else {
					var response strings.Builder
					for k, v := range store{
						response.WriteString(fmt.Sprintf("%s: %s, ",k, v))
					}
					str := strings.TrimSuffix(response.String(), ", ")
					conn.Write([]byte(str + "\n"))
				}
				mu.Unlock()

			default:
				conn.Write([]byte("Error: unknown command\n"))
			
		}

	}
}

func main() {
	ln, err := net.Listen("tcp", ":8080")

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Server running on port 8080...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		go handleConnection(conn)

	}

}