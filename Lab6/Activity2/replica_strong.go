package main

import (
	"fmt"
	"net"
	"net/rpc"
	"os"
	"sync"
)

type Replica struct {
	data  map[string]string
	mu    sync.Mutex
	acks  map[string]int
	peers []string
}

type Args struct {
	Key   string
	Value string
}

func (r *Replica) Update(args *Args, reply *bool) error {
	r.mu.Lock()
	r.data[args.Key] = args.Value
	r.mu.Unlock()
	*reply = true
	fmt.Println("Update received:", args.Key, args.Value)
	return nil
}

func (r *Replica) propagateUpdates(key, value string) {
	r.acks[key] = 0
	for _, peer := range r.peers {
		go func(peer string) {
			client, err := rpc.Dial("tcp", peer)
			if err != nil {
				return
			}
			defer client.Close()

			var reply bool
			client.Call("Replica.Update", &Args{key, value}, &reply)

			if reply {
				r.mu.Lock()
				r.acks[key]++
				fmt.Println("ACK from", peer)
				r.mu.Unlock()
			}
		}(peer)
	}
}

func (r *Replica) waitForAcks(key string, quorum int) {
	for {
		r.mu.Lock()
		if r.acks[key] >= quorum {
			r.mu.Unlock()
			break
		}
		r.mu.Unlock()
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run replica_strong.go <addr> <peer1> [peer2...]")
		return
	}

	addr := os.Args[1]
	peers := os.Args[2:]

	replica := &Replica{
		data:  make(map[string]string),
		acks:  make(map[string]int),
		peers: peers,
	}

	rpc.Register(replica)
	listener, _ := net.Listen("tcp", addr)
	fmt.Println("RPC Replica listening on", addr)

	go func() {
		for {
			conn, _ := listener.Accept()
			go rpc.ServeConn(conn)
		}
	}()

	key := "key1"
	value := "value1"

	replica.Update(&Args{key, value}, new(bool))
	replica.propagateUpdates(key, value)

	quorum := (len(peers) / 2) + 1
	replica.waitForAcks(key, quorum)

	fmt.Println("Update COMMITTED with quorum:", quorum)
}
