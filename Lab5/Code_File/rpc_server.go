package main

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"sync"
)

type Args struct{
	A, B int
}

type Calculator struct {
	lastResult  int
	mu			sync.Mutex
}

func (c *Calculator) store(result int) {
	c.mu.Lock()
	c.lastResult = result
	c.mu.Unlock()
}
func (c *Calculator) Add(args *Args, reply *int) error {
	*reply = args.A + args.B
	c.store(*reply)
	return nil
}

func (c * Calculator) Subtract(args *Args, reply *int) error {
	*reply = args.A - args.B
	c.store(*reply)
	return nil
}

func (c * Calculator) Multiply(args *Args, reply *int) error {
	if args.A == 0 || args.B == 0 {
		return errors.New("multiplication by zero is not allowed")
	}
	*reply = args.A * args.B
	c.store(*reply)
	return nil
}

func (c * Calculator) Divide(args *Args, reply *int) error {
	if args.B == 0 {
		return errors.New("division by zero is not allowed")
	}
	*reply = args.A / args.B
	c.store(*reply)
	return nil
}

func (c *Calculator) GetLastResult(args *Args, reply *int) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	*reply = c.lastResult
	return nil
}

func main() {
	calc := new(Calculator)
	rpc.Register(calc)

	listener, err:= net.Listen("tcp", ":1234")

	if err != nil {
		fmt.Println("Error starting RPC server:", err)
		return 
	}
	

	fmt.Println("Stateful RPC Calculator running...")
	rpc.Accept(listener)
}