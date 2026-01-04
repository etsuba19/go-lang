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
    // _, err = nc.Subscribe(subject, func(m *nats.Msg) {
    //     fmt.Printf("Received a message: %s\n", string(m.Data))
    // })
    // if err != nil {
    //     log.Fatal(err)
    // }
	_, err = nc.Subscribe("updates.info", func(m *nats.Msg) {
        fmt.Println("Info message:", string(m.Data))
    })
    if err != nil {
        log.Fatal(err)
    }
	_, err = nc.Subscribe("updates.error", func(m *nats.Msg) {
        fmt.Println("Error message:", string(m.Data))
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Subscribed to updates. Waiting for messages...")
    select {} 
}
