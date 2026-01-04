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

func main(){
	client, err := rpc.Dial("tcp", "localhost:1234")

	if err != nil {
		log.Fatal("Connection Error:", err)
	}
	
	// Addition
	addArgs := Args{A:3, B:5}
	var addReply int

	err = client.Call("Calculator.Add", &addArgs, &addReply)
	if err != nil {
		log.Fatal("Error calling RPC Add:", err)
	}
	fmt.Printf("Result of %d + %d = %d\n", addArgs.A, addArgs.B, addReply)
	
	var last int 
	err = client.Call("Calculator.GetLastResult", &Args{}, &last)
	fmt.Println("Last Result:", last)


	

	// Subtraction
	subArgs := Args{A:3, B:5}
	var subReply int

	err = client.Call("Calculator.Subtract", &subArgs, &subReply)
	if err != nil {
		log.Fatal("Error calling RPC Subtract:", err)
	}

	fmt.Printf("Result of %d - %d = %d\n", subArgs.A, subArgs.B, subReply)

	err = client.Call("Calculator.GetLastResult", &Args{}, &last)
	fmt.Println("Last Result:", last)

	// Multiplication
	multArgs := Args{A:3, B:5}
	var multReply int

	err = client.Call("Calculator.Multiply", &multArgs, &multReply)
	if err != nil {
		log.Fatal("Error calling RPC Multiply:", err)
	}

	fmt.Printf("Result of %d * %d = %d\n", multArgs.A, multArgs.B, multReply)

	err = client.Call("Calculator.GetLastResult", &Args{}, &last)
	fmt.Println("Last Result:", last)
	

	// Division
	divArgs := Args{A:10, B:5}
	var divReply int

	call := client.Go("Calculator.Divide", &divArgs, &divReply, nil)

	select {
	case <- call.Done:
		if call.Error != nil {
			log.Println("RPC error:", call.Error)
		} else {
			fmt.Printf("Result of %d / %d = %d\n", divArgs.A, divArgs.B, divReply)
		}
	case <- time.After(2 * time.Second):
		log.Panicln("RPC call timed out!")
		
	}

	err = client.Call("Calculator.GetLastResult", &Args{}, &last)
	fmt.Println("Last Result:", last)


	/*
	err = client.Call("Calculator.Divide", &divArgs, &divReply)
	if err != nil {
		log.Fatal("Error calling RPC Divide:", err)
	}

	fmt.Printf("Result of %d / %d = %d\n", divArgs.A, divArgs.B, divReply)
	*/
}
