package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type Replica struct {
	data  map[string]string
	mu    sync.Mutex
	peers []string
}

func (r *Replica) Update(key, value string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[key] = value
}

func (r *Replica) propagateUpdates(key, value string) {
	for _, peer := range r.peers {
		time.Sleep(2 * time.Second)
		go func(peer string) {
			conn, err := net.Dial("tcp", peer)
			if err != nil {
				fmt.Println("Error connecting to peer:", peer, err)
				return
			}
			defer conn.Close()
			fmt.Fprintf(conn, "%s:%s\n", key, value)
		}(peer)
	}
}

func handleConnection(conn net.Conn, replica *Replica) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		parts := strings.Split(strings.TrimSpace(message), ":")
		if len(parts) == 2 {
			replica.Update(parts[0], parts[1])
		}
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run replica_eventual.go <machine_ip:port> <peer1_ip:port>[<peer2_ip:port>...]")
		return
	}

	machineAddr := os.Args[1]
	peers := os.Args[2:]

	replica := &Replica{
		data:  make(map[string]string),
		peers: peers,
	}

	listener, err := net.Listen("tcp", machineAddr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Printf("Replica listening on %s\n", machineAddr)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			go handleConnection(conn, replica)
		}
	}()

	
	replica.Update("key1", "value1")
	replica.propagateUpdates("key1", "value1")

	start := time.Now()
	replica.propagateUpdates("key1", "value1")
	fmt.Println("Propagation time:", time.Since(start))


	replica.mu.Lock()
	fmt.Println("Replica Data:", replica.data)
	replica.mu.Unlock()
}
