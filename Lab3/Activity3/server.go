package main

import (
    "bufio"
    "fmt"
    "net"
	// "strconv"
    "sync"
    "time"
)

var clients = make(map[net.Conn]bool)
var mu sync.Mutex
var clientID = 0

func main() {
    listener, err := net.Listen("tcp", ":8080")

	if err != nil {
 		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is read to assign task...")

    for {
        conn, err := listener.Accept()

		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

        mu.Lock()
		clientID++
		id := clientID
        clients[conn] = true
        mu.Unlock()

		fmt.Printf("Client %d connected.\n", id)

        go handleClient(conn, id)
    }
}

func handleClient(conn net.Conn, id int) {
    defer func() {
        mu.Lock()
        delete(clients, conn)
        mu.Unlock()
        conn.Close()
		fmt.Printf("Client %d disconnected.\n", id)

    }()

    for {
        task := time.Now().Unix() % 100
        fmt.Fprintf(conn, "%d\n", task)

        response, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Printf("Received result %s from Client %d\n", response, id)		

        time.Sleep(5 * time.Second)
    }
}
