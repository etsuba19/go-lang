package main 

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main(){
	conn, err := net.Dial("tcp", "localhost:8080")

	if err != nil {
		fmt.Println("Unable to connect to server.")
		return
	}

	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	serverReader := bufio.NewReader(conn)

	fmt.Println("Connected to Key-Value store server.")

	for{
		fmt.Println("Client: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		_, err := conn.Write([]byte(input + "\n"))
		if err != nil {
			fmt.Println("Server disconnected.")
			return
		}

		resp, err := serverReader.ReadString('\n')
		if err != nil{
			fmt.Println("Server stopped responding.")
			return
		}

		fmt.Print("Server: " + resp)
	}

}