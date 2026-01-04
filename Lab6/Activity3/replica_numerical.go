package main

import (
	"fmt"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Replica struct {
	value float64
	mu    sync.Mutex
	peers []string
}

func (r *Replica) Update(newValue, delta float64) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if math.Abs(newValue-r.value) <= delta {
		old := r.value
		r.value = newValue
		fmt.Printf("ACCEPTED | old=%.1f new=%.1f δ=%.1f\n", old, newValue, delta)
		return true
	}

	fmt.Printf("REJECTED | current=%.1f new=%.1f δ=%.1f\n", r.value, newValue, delta)
	return false
}

func (r *Replica) propagate() {
	r.mu.Lock()
	msg := fmt.Sprintf("%.1f\n", r.value)
	r.mu.Unlock()

	for _, peer := range r.peers {
		go func(p string) {
			conn, err := net.Dial("tcp", p)
			if err != nil {
				return
			}
			defer conn.Close()
			conn.Write([]byte(msg))
		}(peer)
	}
}

func handleConnection(conn net.Conn, replica *Replica, delta float64) {
	defer conn.Close()

	buffer := make([]byte, 64)
	n, err := conn.Read(buffer)
	if err != nil {
		return
	}

	valueStr := strings.TrimSpace(string(buffer[:n]))
	var received float64
	fmt.Sscanf(valueStr, "%f", &received)

	replica.Update(received, delta)
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage:")
		fmt.Println("go run replica.go <addr> <delta> <peer1> [peer2]")
		return
	}

	addr := os.Args[1]
	delta, _ := strconv.ParseFloat(os.Args[2], 64)
	peers := os.Args[3:]

	
	portStr := addr[strings.LastIndex(addr, ":")+1:]
	port, _ := strconv.Atoi(portStr)

	replica := &Replica{
		value: float64(port % 100), 
		peers: peers,
	}

	listener, _ := net.Listen("tcp", addr)
	defer listener.Close()

	fmt.Println("Replica started at", addr)
	fmt.Println("Initial value:", replica.value, "δ =", delta)

	go func() {
		for {
			conn, err := listener.Accept()
			if err == nil {
				go handleConnection(conn, replica, delta)
			}
		}
	}()

	time.Sleep(5 * time.Second)

	for i := 0; i < 4; i++ {
		replica.mu.Lock()
		replica.value += float64(i + 2) 
		fmt.Println("LOCAL UPDATE ->", replica.value)
		replica.mu.Unlock()

		replica.propagate()
		time.Sleep(2 * time.Second)
	}

	fmt.Println("Final value:", replica.value)
	select {}
}
