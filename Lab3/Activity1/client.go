package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	serverReader := bufio.NewReader(conn)

	userReader := bufio.NewScanner(os.Stdin)

	fmt.Println("Connected to server. Type a message:")

	for {
		fmt.Print("> ")

		if !userReader.Scan() {
			break
		}
		message := userReader.Text()

		fmt.Fprintf(conn, message+"\n")

		reply, err := serverReader.ReadString('\n')
		if err != nil {
			fmt.Println("Server closed:", err)
			return
		}

		fmt.Println("Server:", reply)
	}
}
