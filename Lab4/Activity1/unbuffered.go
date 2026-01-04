package main

import (
	"fmt"
	"time"
)

func producer(ch chan<- string){
	for i:= 1; i <= 5;i++ {
		ch<- fmt.Sprintf("Message %d", i)
		time.Sleep(1 * time.Second)
	}
	close(ch)
}

func consumer(ch <- chan string){
	for msg := range ch {
		fmt.Println("Recieved:", msg)
	}
}

func main() {
	ch := make(chan string) 

	go producer(ch)
	consumer(ch)
}