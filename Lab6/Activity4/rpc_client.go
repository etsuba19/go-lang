package main

import (
	"fmt"
	"log"
	"net/rpc"
	"time"
)

type Args struct {
	A, B int
}

func main() {
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Error connecting to server:", err)
	}

	addArgs := Args{A: 10, B: 5}
	var addReply int

	addCall := client.Go("Calculator.Add", &addArgs, &addReply, nil)

	select {
	case <-addCall.Done:
		if addCall.Error != nil {
			log.Println("Add RPC error:", addCall.Error)
		} else {
			fmt.Println("Add Result:", addReply)
		}
	case <-time.After(2 * time.Second):
		log.Println("Add RPC call timed out")
	}

	divArgs := Args{A: 10, B: 2}
	var divReply int

	divCall := client.Go("Calculator.Divide", &divArgs, &divReply, nil)

	select {
	case <-divCall.Done:
		if divCall.Error != nil {
			log.Println("Divide RPC error:", divCall.Error)
		} else {
			fmt.Println("Divide Result:", divReply)
		}
	case <-time.After(2 * time.Second):
		log.Println("Divide RPC call timed out")
	}

	var lastResult int
	err = client.Call("Calculator.GetLastResult", &Args{}, &lastResult)
	if err != nil {
		log.Println("GetLastResult RPC error:", err)
	} else {
		fmt.Println("Last Result from Server:", lastResult)
	}
}
