package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"paxosLab/paxos"
)

var (
	acceptors = []*paxos.Acceptor{
		{Alive: true},
		{Alive: true},
		{Alive: true},
	}
	mu sync.Mutex
)

func proposeHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	var body struct {
		ProposalNumber int
		Value          string
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	done := make(chan interface{})

	go func() {
		proposer := paxos.Proposer{
			ProposalNumber: body.ProposalNumber,
			Value:          body.Value,
		}

		// mu.Lock()
		var value interface{}
		maxRetries := 3

		for i := 0; i < maxRetries; i++ {

			mu.Lock()
			value = proposer.Propose(body.Value, acceptors)
			mu.Unlock()

			if value != nil {
				done <- value
				return
			}
			log.Println("WARNING: Consensus not reached, retrying attempt", i+1)
			time.Sleep(300 * time.Millisecond)
		}

		// mu.Unlock()

		done <- nil
	}()

	select {
	case value := <-done:
		if value != nil {
			fmt.Fprintf(w, "Consensus reached: %s\n", value)
		} else {
			log.Println("ERROR: Consensus failed after retries")
			http.Error(w, "Consensus not reached", http.StatusConflict)
		}

	case <-ctx.Done():
		log.Println("ERROR: Proposal timed out")
		http.Error(w, "Request timed out", http.StatusRequestTimeout)
	}
}

func main() {
	http.HandleFunc("/propose", proposeHandler)

	log.Println("Paxos web service running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}