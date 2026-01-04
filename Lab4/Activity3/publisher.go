package main

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

func main() {
    nc, err := nats.Connect(nats.DefaultURL)
    if err != nil {
        log.Fatal(err)
    }
    defer nc.Close()

    // subject := "updates"
    // message := "Hello, NATS!"

    // if err := nc.Publish(subject, []byte(message)); err != nil {
    //     log.Fatal(err)
    // }

	if err := nc.Publish("updates.info", []byte("Info: All systems operational")); err != nil {
        log.Fatal(err)
    }    
	if err := nc.Publish("updates.error", []byte("Error: Something went wrong")); err != nil {
        log.Fatal(err)
    }

    fmt.Println("Sent messages to updates.info and updates.error")
}
