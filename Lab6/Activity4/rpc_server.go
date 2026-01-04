package main

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"sync"
)

type Args struct {
	A, B int
}

type Calculator struct {
	lastResult int
	mu         sync.Mutex
}

func (c *Calculator) Add(args *Args, reply *int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.lastResult = args.A + args.B
	*reply = c.lastResult
	fmt.Println("Add called:", args.A, "+", args.B, "=", *reply)
	return nil
}

func (c *Calculator) Divide(args *Args, reply *int) error {
	if args.B == 0 {
		return errors.New("division by zero")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.lastResult = args.A / args.B
	*reply = c.lastResult
	fmt.Println("Divide called:", args.A, "/", args.B, "=", *reply)
	return nil
}

func (c *Calculator) GetLastResult(args *Args, reply *int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	*reply = c.lastResult
	fmt.Println("GetLastResult called:", *reply)
	return nil
}

func main() {
	calculator := new(Calculator)
	rpc.Register(calculator)

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}

	fmt.Println("RPC Server listening on port 1234")
	rpc.Accept(listener)
}
