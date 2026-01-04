package main

import (
	"fmt"
	"time"
)

func producer(ch chan<- string){
	for i:= 1; i <= 5;i++ {
		fmt.Printf("Sending: Message %d\n", i)
		ch<- fmt.Sprintf("Message %d", i)
	}
	close(ch)
}

func consumer(ch <- chan string){
	for msg := range ch {
		fmt.Println("Recieved:", msg)
		time.Sleep(2 * time.Second)
	}
}

func main() {
	ch := make(chan string, 3) // buffer size of 3

	go producer(ch)
	consumer(ch)
}