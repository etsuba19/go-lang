package main

import (
	"bufio"
	"fmt"
	"net"
	// "os"
)

func main(){
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server is listening on port 8080...")

	for {
		conn, err := listener.Accept()
		if err != nil {
		fmt.Println("Error accepting connection:", err)
		continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn){
	defer conn.Close()

	message, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Print("Message received:", string(message))

	conn.Write([]byte("Message received:" + message))
}