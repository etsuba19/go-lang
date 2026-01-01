package main

import (
    "bufio"
    "fmt"
    "net"
    "strconv"
    "strings"
)

func main() {
    conn, err := net.Dial("tcp", "localhost:8080")

	if err != nil {
 		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

    for {
        task, _ := bufio.NewReader(conn).ReadString('\n')
		task = strings.TrimSpace(task)

        num, _ := strconv.Atoi(task)
        result := num * num
		
        fmt.Fprintf(conn, "%d\n", result)
    }
}
