package main

import (
	"fmt"
	"time"
)

func producer(ch chan<- string, quit <-chan bool) {
	for i := 1; ; i++ {
		select {
		case <-quit:
			fmt.Println("Producer shutting down")
			return
		case ch <- fmt.Sprintf("Message %d", i):
			fmt.Printf("Produced Message %d\n", i)
			time.Sleep(1 * time.Second)
		}
	}
}

func consumer(ch <-chan string, quit chan<- bool) {
	for i := 0; i < 10; i++ {
		fmt.Println("Consumed:", <-ch)
	}
	quit <- true // shutdown signal
}

func main() {
	ch := make(chan string)
	quit := make(chan bool)

	go producer(ch, quit)
	go consumer(ch, quit)

	<-quit // wait shutdown
	fmt.Println("Main shutting down")
}
