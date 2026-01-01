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
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()
	go receiveMessages(conn)
 
	for {
		message, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		fmt.Fprintf(conn, message)
	}
}

func receiveMessages(conn net.Conn) {
	for {
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("\nServer disconnected. Closing client...")
		os.Exit(0)
    }
	fmt.Print("Message from server: ", message)
	}
}